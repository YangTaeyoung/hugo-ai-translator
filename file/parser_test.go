package file

import (
	"context"
	_ "embed"
	"os"
	"path"
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

func Test_parser_Simple(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	contentDir := path.Join(currentDir, "test_simple_content_dir")

	type fields struct {
		cfg ParserConfig
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ParsedMarkdownFiles
		wantErr bool
	}{
		{
			name: "성공 시",
			fields: fields{
				cfg: ParserConfig{
					ContentDir: contentDir,
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
						config.LanguageCodeKorean,
					},
				},
			},
			args: args{
				ctx: t.Context(),
			},
			want: ParsedMarkdownFiles{
				{
					Path:     "test.md",
					Markdown: "# Hello, World!",
					TargetLanguages: config.LanguageCodes{
						config.LanguageCodeEnglish,
						config.LanguageCodeKorean,
					},
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
			got, err := p.Simple(tt.args.ctx)
			assert.Equalf(t, tt.wantErr, err != nil, "parser.Simple() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equalf(t, tt.want, got, "parser.Simple(%v)", tt.args.ctx)
		})
	}
}

func TestNewParser(t *testing.T) {
	testConfig := ParserConfig{
		ContentDir:      "test_content",
		IgnoreRules:     []string{"world/**"},
		TargetLanguages: config.LanguageCodes{config.LanguageCodeEnglish},
		TargetPathRule:  "{{origin}}/{{filename}}.{{language}}.md",
	}
	type args struct {
		cfg ParserConfig
	}
	tests := []struct {
		name string
		args args
		want Parser
	}{
		{
			name: "성공",
			args: args{
				cfg: testConfig,
			},
			want: &parser{
				cfg: testConfig,
			},
		},
	}
	for _, tt := range tests {
		p := NewParser(tt.args.cfg)
		assert.Equalf(t, tt.want, p, "NewParser() = %v, want %v", p, tt.want)
	}
}
