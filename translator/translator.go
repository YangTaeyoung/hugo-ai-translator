package translator

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"log/slog"
	"path/filepath"
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
	Model           openai.ChatModel
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
	for _, targetLanguage := range source.TargetLanguages {
		targetLanguage := targetLanguage
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

			schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:        openai.F("markdown"),
				Description: openai.F("translated markdown"),
				Schema:      openai.F(TranslateMarkdownSchema()),
				Strict:      openai.Bool(true),
			}

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
				ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
					openai.ResponseFormatJSONSchemaParam{
						Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
						JSONSchema: openai.F(schemaParam),
					}),
				Model: openai.F(t.cfg.Model),
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

			var response TranslateResponse
			if err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &response); err != nil {
				return err
			}

			fileName, err = file.FileNameWithoutExtension(source.Path)
			if err != nil {
				return err
			}

			mu.Lock()
			translated = append(translated, Result{
				Language:  targetLanguage,
				FileName:  fileName,
				OriginDir: filepath.Dir(source.Path),
				Markdown:  file.Markdown(response.Markdown),
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
