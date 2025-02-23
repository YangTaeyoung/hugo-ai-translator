package translator

import (
	"bytes"
	"context"
	_ "embed"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/YangTaeyoung/hugo-ai-translator/file"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed instruction.md
	instructionMd string

	//go:embed prompt.md
	promptMd string
)

var (
	ErrorEmptyResult = errors.New("empty result")
)

type Config struct {
	SourceLanguage  config.LanguageCode
	TargetLanguages config.LanguageCodes
}

type Translator interface {
	Translate(ctx context.Context, source file.ParsedMarkdownFile) (Results, error)
}

type translator struct {
	client *openai.Client
	cfg    *Config
}

type promptProps struct {
	SourceLanguage config.LanguageCode
	TargetLanguage config.LanguageCode
	Source         string
}

type Result struct {
	Language  config.LanguageCode
	FileName  string
	OriginDir string
	Markdown  file.Markdown
}

type Results []Result

func (r Results) MarkdownFiles() file.MarkdownFiles {
	var files file.MarkdownFiles

	for _, result := range r {
		files = append(files, file.MarkdownFile{
			FileName:  result.FileName,
			OriginDir: result.OriginDir,
			Language:  result.Language,
			Content:   result.Markdown,
		})
	}

	return files
}

func (t translator) Translate(ctx context.Context, source file.ParsedMarkdownFile) (Results, error) {
	var (
		translated Results
		mu         sync.Mutex
		err        error
	)

	g, gctx := errgroup.WithContext(ctx)
	for _, targetLanguage := range t.cfg.TargetLanguages {
		g.Go(func() error {
			var (
				fileName string
				tmpl     *template.Template
				res      *openai.ChatCompletion
			)
			slog.DebugContext(ctx, "translating markdown file", "language", targetLanguage, "path", source.Path)

			tmpl, err = template.New("prompt").Parse(promptMd)
			if err != nil {
				return err
			}

			var buf bytes.Buffer

			if err = tmpl.Execute(&buf, promptProps{
				SourceLanguage: t.cfg.SourceLanguage,
				TargetLanguage: targetLanguage,
				Source:         source.Markdown.String(),
			}); err != nil {
				return err
			}

			prompt := buf.String()

			res, err = t.client.Chat.Completions.New(gctx, openai.ChatCompletionNewParams{
				Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
					openai.ChatCompletionDeveloperMessageParam{
						Role: openai.F(openai.ChatCompletionDeveloperMessageParamRoleDeveloper),
						Content: openai.F([]openai.ChatCompletionContentPartTextParam{
							{
								Text: openai.F(instructionMd),
								Type: openai.F(openai.ChatCompletionContentPartTextTypeText),
							},
						}),
					},
					openai.UserMessage(prompt),
				}),
				Model: openai.F(openai.ChatModelGPT3_5Turbo),
			})

			if err != nil {
				return err
			}

			if len(res.Choices) == 0 {
				return ErrorEmptyResult
			}

			if res.Choices[0].Message.Content == "" {
				return ErrorEmptyResult
			}

			content := res.Choices[0].Message.Content
			content = strings.Trim(content, " ")
			content = strings.TrimPrefix(content, "```markdown")
			content = strings.TrimSuffix(content, "\n")
			content = strings.TrimSuffix(content, "```")

			fileName, err = fileNameWithoutExtension(source.Path)
			if err != nil {
				return err
			}

			mu.Lock()
			translated = append(translated, Result{
				Language:  targetLanguage,
				FileName:  fileName,
				OriginDir: filepath.Dir(source.Path),
				Markdown:  file.Markdown(res.Choices[0].Message.Content),
			})
			mu.Unlock()

			slog.DebugContext(ctx, "translated markdown file", "language", targetLanguage, "fileName", fileName)

			return nil
		})
	}
	if err = g.Wait(); err != nil {
		return nil, err
	}

	return translated, nil
}

func New(client *openai.Client, cfg Config) Translator {
	return &translator{
		client: client,
		cfg:    &cfg,
	}
}
