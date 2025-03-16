package translator

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/YangTaeyoung/hugo-ai-translator/file"
	"github.com/YangTaeyoung/hugo-ai-translator/llm"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/stretchr/testify/assert"
)

//go:embed test_posting.md
var testPostingMd string

func Test_Integration_Translate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	ctx := context.Background()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	slog.InfoContext(ctx, "homeDir", "homeDir", homeDir)
	cfg, err := config.New(path.Join(homeDir, ".hugo_ai_translator", "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	client := llm.NewOpenAIClient(openai.NewClient(option.WithAPIKey(cfg.OpenAI.ApiKey)))
	type fields struct {
		client llm.OpenAIClient
		cfg    *Config
	}

	type args struct {
		ctx    context.Context
		source *file.MarkdownFile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    int
	}{
		{
			name: "통합 테스트",
			fields: fields{
				client: client,
				cfg: &Config{
					SourceLanguage: "ko",
					Model:          openai.ChatModelGPT4oMini,
				},
			},
			args: args{
				ctx: ctx,
				source: &file.MarkdownFile{
					Content:   file.Markdown(testPostingMd),
					OriginDir: "path/to",
					FileName:  "file",
					Language:  config.LanguageCodeEnglish,
				},
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := translator{
				client: tt.fields.client,
				cfg:    tt.fields.cfg,
			}
			err = tr.Translate(tt.args.ctx, tt.args.source)
			fileName := fmt.Sprintf("testing.%s.md", tt.args.source.Language.String())
			if err = os.WriteFile(path.Join(".", "test_result", fileName), []byte(tt.args.source.Translated), 0644); err != nil {
				t.Fatal(err)
			}

			slog.InfoContext(ctx, "translated", "content", tt.args.source.Translated)
			assert.Equalf(t, tt.wantErr, err != nil, "translator.Translate() error = %v, wantErr %v", err, tt.wantErr)
			assert.NotEmptyf(t, tt.args.source.Translated, "translated content is empty")
		})
	}
}
