package translator

import (
	"context"
	"testing"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/YangTaeyoung/hugo-ai-translator/file"
	"github.com/YangTaeyoung/hugo-ai-translator/llm"
	"github.com/YangTaeyoung/hugo-ai-translator/mocks"
	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_translator_Translate(t *testing.T) {
	var testPrompt = `The source language and the target language are given, please translate the content from the source language to the target language.  

purpose is to correctly convert the source language to the target language.

The content that needs to be translated is as follows.
- title field in markdown front matter
  - do not ":" in the title field because it is used as a delimiter.
- Markdown content

## SourceLanguage
Korean

## TargetLanguage
English

## Source
"""
안녕, 세계!
"""`

	type fields struct {
		cfg *Config
	}
	type args struct {
		ctx    context.Context
		source file.ParsedMarkdownFile
	}
	tests := []struct {
		name       string
		fields     fields
		mockClient func() llm.OpenAIClient
		args       args
		want       file.MarkdownFiles
		wantErr    bool
	}{
		{
			name: "성공 시",
			fields: fields{
				cfg: &Config{
					SourceLanguage: config.LanguageCodeKorean,
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
					Model: openai.ChatModelGPT4oMini,
				},
			},
			mockClient: func() llm.OpenAIClient {
				m := mocks.NewOpenAIClient(t)
				m.EXPECT().New(mock.Anything, openai.ChatCompletionNewParams{
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
						openai.UserMessage(testPrompt),
					}),
					ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
						openai.ResponseFormatJSONSchemaParam{
							Type: openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
							JSONSchema: openai.F(openai.ResponseFormatJSONSchemaJSONSchemaParam{
								Name:        openai.F("markdown"),
								Description: openai.F("translated markdown"),
								Schema:      openai.F(TranslateMarkdownSchema()),
								Strict:      openai.Bool(true),
							}),
						}),
					Model: openai.F(openai.ChatModelGPT4oMini),
				}).Return(&openai.ChatCompletion{
					Choices: []openai.ChatCompletionChoice{
						{
							Message: openai.ChatCompletionMessage{
								Content: "{\"markdown\":\"Hello, world!\"}",
							},
						},
					},
				}, nil)

				return m
			},
			args: args{
				ctx: t.Context(),
				source: file.ParsedMarkdownFile{
					Path:     "hello/foo.md",
					Markdown: "안녕, 세계!",
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
				},
			},
			want: file.MarkdownFiles{
				{
					FileName:  "foo",
					OriginDir: "hello",
					Language:  "en",
					Content:   "Hello, world!",
				},
			},
			wantErr: false,
		},
		{
			name: "res.Choices가 비어있을 때",
			fields: fields{
				cfg: &Config{
					SourceLanguage: config.LanguageCodeKorean,
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
					Model: openai.ChatModelGPT4oMini,
				},
			},
			mockClient: func() llm.OpenAIClient {
				m := mocks.NewOpenAIClient(t)
				m.EXPECT().New(mock.Anything, openai.ChatCompletionNewParams{
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
						openai.UserMessage(testPrompt),
					}),
					ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
						openai.ResponseFormatJSONSchemaParam{
							Type: openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
							JSONSchema: openai.F(openai.ResponseFormatJSONSchemaJSONSchemaParam{
								Name:        openai.F("markdown"),
								Description: openai.F("translated markdown"),
								Schema:      openai.F(TranslateMarkdownSchema()),
								Strict:      openai.Bool(true),
							}),
						}),
					Model: openai.F(openai.ChatModelGPT4oMini),
				}).Return(&openai.ChatCompletion{
					Choices: nil,
				}, nil)

				return m
			},
			args: args{
				ctx: t.Context(),
				source: file.ParsedMarkdownFile{
					Path:     "hello/foo.md",
					Markdown: "안녕, 세계!",
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "res.Choices[0].Message.Content가 비어있을 때",
			fields: fields{
				cfg: &Config{
					SourceLanguage: config.LanguageCodeKorean,
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
					Model: openai.ChatModelGPT4oMini,
				},
			},
			mockClient: func() llm.OpenAIClient {
				m := mocks.NewOpenAIClient(t)
				m.EXPECT().New(mock.Anything, openai.ChatCompletionNewParams{
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
						openai.UserMessage(testPrompt),
					}),
					ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
						openai.ResponseFormatJSONSchemaParam{
							Type: openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
							JSONSchema: openai.F(openai.ResponseFormatJSONSchemaJSONSchemaParam{
								Name:        openai.F("markdown"),
								Description: openai.F("translated markdown"),
								Schema:      openai.F(TranslateMarkdownSchema()),
								Strict:      openai.Bool(true),
							}),
						}),
					Model: openai.F(openai.ChatModelGPT4oMini),
				}).Return(&openai.ChatCompletion{
					Choices: []openai.ChatCompletionChoice{
						{
							Message: openai.ChatCompletionMessage{
								Content: "",
							},
						},
					},
				}, nil)

				return m
			},
			args: args{
				ctx: t.Context(),
				source: file.ParsedMarkdownFile{
					Path:     "hello/foo.md",
					Markdown: "안녕, 세계!",
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := translator{
				client: tt.mockClient(),
				cfg:    tt.fields.cfg,
			}
			got, err := tr.Translate(tt.args.ctx, tt.args.source)
			assert.Equalf(t, tt.wantErr, err != nil, "Translate() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equalf(t, tt.want, got, "Translate(%v, %v)", tt.args.ctx, tt.args.source)
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		client llm.OpenAIClient
		cfg    Config
	}
	tests := []struct {
		name string
		args args
		want Translator
	}{
		{
			name: "성공",
			args: args{
				client: mocks.NewOpenAIClient(t),
				cfg: Config{
					SourceLanguage:  config.LanguageCodeKorean,
					TargetLanguages: config.LanguageCodes{config.LanguageCodeEnglish},
					Model:           openai.ChatModelGPT4oMini,
				},
			},
			want: &translator{
				client: mocks.NewOpenAIClient(t),
				cfg: &Config{
					SourceLanguage:  config.LanguageCodeKorean,
					TargetLanguages: config.LanguageCodes{config.LanguageCodeEnglish},
					Model:           openai.ChatModelGPT4oMini,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, New(tt.args.client, tt.args.cfg), "New(%v, %v)", tt.args.client, tt.args.cfg)
		})
	}
}
