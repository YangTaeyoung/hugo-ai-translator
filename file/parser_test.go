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
		want    MarkdownFiles
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
					SourceLanguage: config.LanguageCodeKorean,
				},
			},
			want: MarkdownFiles{
				{
					OriginDir: "hello",
					FileName:  "bar",
					Language:  config.LanguageCodeEnglish,
					Content:   Markdown(barMd),
				},
				{
					OriginDir: "hello",
					FileName:  "foo",
					Language:  config.LanguageCodeEnglish,
					Content:   Markdown(fooMd),
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
		want    MarkdownFiles
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
			want: MarkdownFiles{
				{
					FileName:  "test",
					Content:   "# Hello, World!",
					OriginDir: ".",
					Language:  config.LanguageCodeEnglish,
				},
				{
					FileName:  "test",
					Content:   "# Hello, World!",
					OriginDir: ".",
					Language:  config.LanguageCodeKorean,
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
		SourceLanguage:  config.LanguageCodeKorean,
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
