package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/YangTaeyoung/hugo-ai-translator/environment"
	"github.com/YangTaeyoung/hugo-ai-translator/file"
	"github.com/YangTaeyoung/hugo-ai-translator/llm"
	"github.com/YangTaeyoung/hugo-ai-translator/translator"
	"github.com/k0kubun/go-ansi"
	"github.com/manifoldco/promptui"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

var ChatModels = []string{
	openai.ChatModelO3Mini,
	openai.ChatModelO3Mini2025_01_31,
	openai.ChatModelO1,
	openai.ChatModelO1_2024_12_17,
	openai.ChatModelO1Preview,
	openai.ChatModelO1Preview2024_09_12,
	openai.ChatModelO1Mini,
	openai.ChatModelO1Mini2024_09_12,
	openai.ChatModelGPT4o,
	openai.ChatModelGPT4o2024_11_20,
	openai.ChatModelGPT4o2024_08_06,
	openai.ChatModelGPT4o2024_05_13,
	openai.ChatModelGPT4oAudioPreview,
	openai.ChatModelGPT4oAudioPreview2024_10_01,
	openai.ChatModelGPT4oAudioPreview2024_12_17,
	openai.ChatModelGPT4oMiniAudioPreview,
	openai.ChatModelGPT4oMiniAudioPreview2024_12_17,
	openai.ChatModelChatgpt4oLatest,
	openai.ChatModelGPT4oMini,
	openai.ChatModelGPT4oMini2024_07_18,
	openai.ChatModelGPT4Turbo,
	openai.ChatModelGPT4Turbo2024_04_09,
	openai.ChatModelGPT4_0125Preview,
	openai.ChatModelGPT4TurboPreview,
	openai.ChatModelGPT4_1106Preview,
	openai.ChatModelGPT4VisionPreview,
	openai.ChatModelGPT4,
	openai.ChatModelGPT4_0314,
	openai.ChatModelGPT4_0613,
	openai.ChatModelGPT4_32k,
	openai.ChatModelGPT4_32k0314,
	openai.ChatModelGPT4_32k0613,
	openai.ChatModelGPT3_5Turbo,
	openai.ChatModelGPT3_5Turbo16k,
	openai.ChatModelGPT3_5Turbo0301,
	openai.ChatModelGPT3_5Turbo0613,
	openai.ChatModelGPT3_5Turbo16k0613,
	openai.ChatModelGPT3_5Turbo1106,
	openai.ChatModelGPT3_5Turbo0125,
}

var progressbarOpts = []progressbar.Option{
	progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
	progressbar.OptionEnableColorCodes(true),
	progressbar.OptionSetWidth(15),
	progressbar.OptionSetDescription("Translating ..."),
	progressbar.OptionSetTheme(progressbar.Theme{
		Saucer:        "[green]=[reset]",
		SaucerHead:    "[green]>[reset]",
		SaucerPadding: " ",
		BarStart:      "[",
		BarEnd:        "]",
	}),
}

func TranslateAction(ctx context.Context, cmd *cli.Command) error {
	var (
		cfgPath = cmd.String("config")
	)

	cfg, err := config.New(cfgPath)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "config parsed", "path", cfgPath)

	env := environment.New(cfg)
	slog.InfoContext(ctx, "environment created")

	markdownFiles, err := env.Parser.Parse(ctx)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "markdown files parsed", "count", len(markdownFiles))

	bar := progressbar.NewOptions(len(markdownFiles), progressbarOpts...)

	for _, markdownFile := range markdownFiles {
		var results file.MarkdownFiles

		bar.Describe(fmt.Sprintf("Translating %s ...", markdownFile.Path))

		results, err = env.Translator.Translate(ctx, markdownFile)
		if err != nil {
			return err
		}

		if err = env.Writer.Write(ctx, results); err != nil {
			return err
		}

		if err = bar.Add(1); err != nil {
			return errors.Wrap(err, "failed to update progress bar")
		}
	}

	return nil
}

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrEmptyInput   = errors.New("empty input")
)

func ConfigureAction(_ context.Context, cmd *cli.Command) error {
	var (
		cfgPath = cmd.String("config")
		dryRun  = cmd.Bool("dry-run")
		cfg     config.Config
		p       promptui.Prompt
		err     error
	)
	if strings.HasPrefix(cfgPath, "~") {
		var homeDir string
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return err
		}

		cfgPath = filepath.Join(homeDir, cfgPath[1:])
	}

	if _, err = os.Stat(cfgPath); os.IsExist(err) {
		var answer string

		p = promptui.Prompt{
			Label: "Config file already exists. Do you want to overwrite it? (y/n)",
			Validate: func(s string) error {
				s = strings.ToLower(s)
				if !slices.Contains([]string{"y", "n"}, s) {
					return ErrInvalidInput
				}

				return nil
			},
		}

		answer, err = p.Run()
		if err != nil {
			return errors.Wrap(err, "failed to get answer")
		}

		if answer == "n" {
			return nil
		}
	}

	if err = openAIStep(&cfg); err != nil {
		return err
	}

	if err = contentDirStep(&cfg); err != nil {
		return err
	}

	if err = languageChoiceStep(&cfg); err != nil {
		return err
	}

	if err = ignoreRuleStep(&cfg); err != nil {
		return err
	}

	if err = targetPathRuleStep(&cfg); err != nil {
		return err
	}

	configFile, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	if dryRun {
		fmt.Println("Config files are: ")
		fmt.Println(string(configFile))
		return nil
	}

	if err = os.MkdirAll(filepath.Dir(cfgPath), os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create config directory")
	}

	if err = os.WriteFile(cfgPath, configFile, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write config file")
	}

	fmt.Println("Config file is created at ", cfgPath)

	return nil
}

func DebugModeAction(ctx context.Context, _ *cli.Command, debug bool) error {
	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.DebugContext(ctx, "debug mode on")

	return nil
}

func SimpleTranslateAction(ctx context.Context, cmd *cli.Command) error {
	var (
		sourceLanguage  config.LanguageCode
		targetLanguages config.LanguageCodes
		cfg             *config.Config
		cfgPath         = cmd.String("config")
		apiKey          = cmd.String("api-key")
		model           = cmd.String("model")
	)

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if cmd.String("source-language") == "" {
		if cfg, err = config.New(cfgPath); err != nil {
			return err
		}

		sourceLanguage = cfg.Translator.Source.SourceLanguage
	} else {
		sourceLanguage = config.LanguageCode(cmd.String("source-language"))
	}

	if len(cmd.StringSlice("target-languages")) == 0 {
		if cfg, err = config.New(cfgPath); err != nil {
			return err
		}
		targetLanguages = cfg.Translator.Target.TargetLanguages
	} else {
		for _, lang := range cmd.StringSlice("target-languages") {
			targetLanguages = append(targetLanguages, config.LanguageCode(lang))
		}
	}
	if apiKey == "" {
		if cfg, err = config.New(cfgPath); err != nil {
			return err
		}
		apiKey = cfg.OpenAI.ApiKey
	}
	if model == "" {
		if cfg, err = config.New(cfgPath); err != nil {
			return err
		}
		model = cfg.OpenAI.Model
	}
	slog.InfoContext(ctx, "config parsed",
		"sourceLanguage", sourceLanguage,
		"targetLanguages", targetLanguages,
		"apiKey", apiKey,
		"model", model,
	)

	client := llm.NewOpenAIClient(openai.NewClient(option.WithAPIKey(apiKey)))

	t := translator.New(client, translator.Config{
		SourceLanguage:  sourceLanguage,
		TargetLanguages: targetLanguages,
		Model:           model,
	})

	parser := file.NewParser(file.ParserConfig{
		ContentDir:      currentDir,
		TargetLanguages: targetLanguages,
	})

	parsedMarkdownFiles, err := parser.Simple(ctx)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "markdown files parsed",
		"count", len(parsedMarkdownFiles),
	)

	bar := progressbar.NewOptions(len(parsedMarkdownFiles),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Translating ..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	var results file.MarkdownFiles
	for _, markdownFile := range parsedMarkdownFiles {
		var translated file.MarkdownFiles

		bar.Describe(fmt.Sprintf("Translating %s ...", markdownFile.Path))

		translated, err = t.Translate(ctx, markdownFile)
		if err != nil {
			return err
		}

		results = append(results, translated...)

		if err = bar.Add(1); err != nil {
			return errors.Wrap(err, "failed to update progress bar")
		}
	}

	w := file.NewWriter(file.WriterConfig{
		ContentDir:     currentDir,
		TargetPathRule: "{fileName}.{language}.md",
	})

	if results == nil {
		return errors.New("no results")
	}

	if err = w.Write(ctx, results); err != nil {
		return err
	}
	slog.InfoContext(ctx, "markdown files written", "count", len(results))

	return nil
}
