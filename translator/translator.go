package translator

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/YangTaeyoung/hugo-ai-translator/file"
	"github.com/YangTaeyoung/hugo-ai-translator/llm"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
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
	Translate(ctx context.Context, source *file.MarkdownFile) error
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

func (t *translator) Translate(ctx context.Context, source *file.MarkdownFile) error {
	var (
		tmpl *template.Template
		res  *openai.ChatCompletion
		err  error
	)
	slog.DebugContext(ctx, "translating markdown file", "language", source.Language, "originDir", source.OriginDir, "fileName", source.FileName)

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
		TargetLanguage: source.Language.Name().String(),
		Source:         source.Content.String(),
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

	res, err = t.client.New(ctx, openai.ChatCompletionNewParams{
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
		return errors.Wrap(err, "failed to translate markdown")
	}

	if len(res.Choices) == 0 {
		return ErrorEmptyResult
	}

	if res.Choices[0].Message.Content == "" {
		return ErrorEmptyResult
	}

	var response TranslateResponse
	if err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &response); err != nil {
		fmt.Println(res.Choices[0].Message.Content)
		return errors.Wrap(err, "failed to unmarshal translated markdown response")
	}

	source.Translated = file.Markdown(response.Markdown)

	slog.DebugContext(ctx, "translated markdown file", "language", source.Language, "fileName", source.FileName)

	return nil
}
