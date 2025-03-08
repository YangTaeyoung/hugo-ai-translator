package file

import (
	"context"
	_ "embed"
	"testing"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed test_content/hello/foo.md
	fooMd string
	//go:embed test_content/hello/bar.md
	barMd string
)

func Test_parser_Parse(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		cfg ParserConfig
	}
	tests := []struct {
		name    string
		fields  fields
		want    ParsedMarkdownFiles
		wantErr bool
	}{
		{
			name: "정상 동작 시",
			fields: fields{
				cfg: ParserConfig{
					ContentDir: "./test_content",
					IgnoreRules: []string{
						"world/**",
					},
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
					TargetPathRule: "{{origin}}/{{filename}}.{{language}}.md",
				},
			},
			want: ParsedMarkdownFiles{
				{
					Path: "hello/bar.md",
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
					Markdown: Markdown(barMd),
				},
				{
					Path: "hello/foo.md",
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
					},
					Markdown: Markdown(fooMd),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser{
				cfg: tt.fields.cfg,
			}

			got, err := p.Parse(ctx)

			assert.Equalf(t, tt.want, got, "parser.Parse() = %v, want %v", got, tt.want)
			assert.Equalf(t, tt.wantErr, err != nil, "parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
