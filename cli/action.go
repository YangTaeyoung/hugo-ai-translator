package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/YangTaeyoung/hugo-ai-translator/environment"
	"github.com/YangTaeyoung/hugo-ai-translator/file"
	"github.com/YangTaeyoung/hugo-ai-translator/translator"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
)

func TranslateAction(ctx context.Context, cmd *cli.Command) error {
	var (
		translatedPaths []string
		cfgPath         = cmd.String("config")
		reTranslate     = cmd.Bool("--re-translate")
	)

	cfg, err := config.New(cfgPath)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "config parsed", "path", cfgPath)

	if !reTranslate {
		translatedPaths, err = config.TranslatedPaths(cmd.String("history-path"))
		if err != nil {
			return err
		}
	}

	env := environment.New(cfg)
	slog.InfoContext(ctx, "environment created")

	p := file.NewParser(file.ParserConfig{
		ContentDir:      cfg.Translator.ContentDir,
		TranslatedPaths: translatedPaths,
		TargetLanguages: cfg.Translator.Target.TargetLanguages,
		TargetPathRule:  cfg.Translator.Target.TargetPathRule,
		IgnoreRules:     cfg.Translator.Source.IgnoreRules,
	})

	markdownFiles, err := p.Parse(ctx)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "markdown files parsed", "count", len(markdownFiles))

	t := translator.New(env.Client, translator.Config{
		SourceLanguage:  cfg.Translator.Source.SourceLanguage,
		TargetLanguages: cfg.Translator.Target.TargetLanguages,
		Model:           cfg.OpenAI.Model,
	})

	w := file.NewWriter(file.WriterConfig{
		ContentDir:     cfg.Translator.ContentDir,
		TargetPathRule: cfg.Translator.Target.TargetPathRule,
	})

	bar := progressbar.NewOptions(len(markdownFiles),
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

	for _, markdownFile := range markdownFiles {
		var (
			results translator.Results
			saved   []string
		)
		bar.Describe(fmt.Sprintf("Translating %s ...", markdownFile.Path))

		results, err = t.Translate(ctx, markdownFile)
		if err != nil {
			return err
		}

		saved, err = w.Write(ctx, results.MarkdownFiles())
		if err != nil {
			return err
		}

		if reTranslate {
			if err = config.WriteTranslatedPaths(cmd.String("history-path"), saved...); err != nil {
				return err
			}
		} else {
			if err = config.AppendTranslatedPaths(cmd.String("history-path"), saved...); err != nil {
				return err
			}
		}

		if err = bar.Add(1); err != nil {
			return err
		}
	}

	return nil
}
