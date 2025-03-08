package file

import (
	"context"
	"os"
	"testing"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/stretchr/testify/assert"
)

func Test_writer_Write(t *testing.T) {
	type fields struct {
		cfg WriterConfig
	}
	type args struct {
		ctx   context.Context
		files MarkdownFiles
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "성공",
			fields: fields{
				cfg: WriterConfig{
					ContentDir:     "test_writer_content",
					TargetPathRule: "{origin}/{fileName}.{language}.md",
				},
			},
			args: args{
				ctx: t.Context(),
				files: MarkdownFiles{
					{
						FileName:  "test",
						OriginDir: "origin_dir",
						Language:  config.LanguageCodeKorean,
						Content:   "# Hello",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := writer{
				cfg: tt.fields.cfg,
			}
			err := w.Write(tt.args.ctx, tt.args.files)
			assert.Equalf(t, tt.wantErr, err != nil, "Write() error = %v, wantErr %v", err, tt.wantErr)

			file, err := os.ReadFile("test_writer_content/origin_dir/test.ko.md")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "---\ntranslated: true\n---\n# Hello", string(file))
		})
	}
}

func TestNewWriter(t *testing.T) {
	testConfig := WriterConfig{
		ContentDir:     "test_content",
		TargetPathRule: "{origin}/{fileName}.{language}.md",
	}
	type args struct {
		cfg WriterConfig
	}
	tests := []struct {
		name string
		args args
		want Writer
	}{
		{
			name: "성공",
			args: args{
				cfg: testConfig,
			},
			want: &writer{
				cfg: testConfig,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWriter(tt.args.cfg)
			assert.Equalf(t, tt.want, w, "NewWriter(%v)", tt.args.cfg)
		})
	}
}
