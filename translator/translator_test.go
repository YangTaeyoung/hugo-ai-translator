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
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/stretchr/testify/assert"
)

//go:embed test_posting.md
var testPostingMd string

func Test_translator_Translate(t *testing.T) {
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

	client := openai.NewClient(option.WithAPIKey(cfg.OpenAI.ApiKey))
	type fields struct {
		client *openai.Client
		cfg    *Config
	}

	type args struct {
		ctx    context.Context
		source file.ParsedMarkdownFile
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
				source: file.ParsedMarkdownFile{
					Markdown: file.Markdown(testPostingMd),
					Path:     "path/to/file",
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish, config.LanguageCodeFrench, config.LanguageCodeChinese,
					},
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
			got, err := tr.Translate(tt.args.ctx, tt.args.source)
			for _, g := range got {
				fileName := fmt.Sprintf("testing.%s.md", g.Language.String())
				if err = os.WriteFile(path.Join(".", "test_result", fileName), []byte(g.Markdown), 0644); err != nil {
					t.Fatal(err)
				}
			}
			slog.InfoContext(ctx, "translated", "got", got)
			assert.Equalf(t, tt.wantErr, err != nil, "translator.Translate() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equalf(t, len(got), tt.want, "translator.Translate() = %v, want %v", len(got), tt.want)
		})
	}
}
