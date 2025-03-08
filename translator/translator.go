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
	"github.com/YangTaeyoung/hugo-ai-translator/llm"
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

var MaxWorkers = 4

var (
	ErrorEmptyResult = errors.New("empty result")
)

type Config struct {
	SourceLanguage  config.LanguageCode
	TargetLanguages config.LanguageCodes
	Model           openai.ChatModel
}

type Translator interface {
	Translate(ctx context.Context, source file.ParsedMarkdownFile) (file.MarkdownFiles, error)
}

type translator struct {
	client llm.OpenAIClient
	cfg    *Config
}

func New(client llm.OpenAIClient, cfg Config) Translator {
	return &translator{
		client: client,
		cfg:    &cfg,
	}
}

func (t *translator) Translate(ctx context.Context, source file.ParsedMarkdownFile) (file.MarkdownFiles, error) {
	var (
		translated file.MarkdownFiles
		mu         sync.Mutex
		err        error
	)

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(MaxWorkers)

	for _, language := range source.TargetLanguages {
		targetLanguage := language

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

			if err = tmpl.Execute(&buf, struct {
				SourceLanguage string
				TargetLanguage string
				Source         string
			}{
				SourceLanguage: t.cfg.SourceLanguage.Name().String(),
				TargetLanguage: targetLanguage.Name().String(),
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

			res, err = t.client.New(gctx, openai.ChatCompletionNewParams{
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
			translated = append(translated, file.MarkdownFile{
				Language:  targetLanguage,
				FileName:  fileName,
				OriginDir: filepath.Dir(source.Path),
				Content:   file.Markdown(response.Markdown),
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
