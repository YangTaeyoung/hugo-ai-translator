package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteMarkdownWithFrontmatter(t *testing.T) {
	updateTestFile, err := os.ReadFile("test_md/update_test.md")
	if err != nil {
		t.Fatal(err)
	}
	updateTestMd := string(updateTestFile)

	updateWantFile, err := os.ReadFile("test_md/update_want.md")
	if err != nil {
		t.Fatal(err)
	}

	createTestFile, err := os.ReadFile("test_md/create_test.md")
	if err != nil {
		t.Fatal(err)
	}

	createTestMd := string(createTestFile)

	createWantFile, err := os.ReadFile("test_md/create_want.md")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		path      string
		file      []byte
		perm      os.FileMode
		keyValues []any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "파라미터가 홀수인 경우",
			args: args{
				"test_md/odd_key_values.md",
				[]byte(""),
				os.ModePerm,
				[]any{"translated"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "key가 string이 아닌 경우",
			args: args{
				"test_md/not_string_key.md",
				[]byte(""),
				os.ModePerm,
				[]any{1, true},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "업데이트 케이스",
			args: args{
				path:      "test_md/update_result.md",
				file:      []byte(updateTestMd),
				perm:      os.ModePerm,
				keyValues: []any{"translated", true},
			},
			want:    string(updateWantFile),
			wantErr: false,
		},
		{
			name: "생성 케이스",
			args: args{
				path:      "test_md/create_result.md",
				file:      []byte(createTestMd),
				perm:      os.ModePerm,
				keyValues: []any{"translated", true},
			},
			want:    string(createWantFile),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotFile []byte
			err = WriteMarkdownWithFrontmatter(tt.args.path, tt.args.file, tt.args.perm, tt.args.keyValues...)
			assert.Equalf(t, tt.wantErr, err != nil, "writeMarkdownWithFrontmatter() error = %v, wantErr %v", err, tt.wantErr)

			if !tt.wantErr {
				gotFile, err = os.ReadFile(tt.args.path)
				if err != nil {
					t.Fatal(err)
				}

				got := string(gotFile)
				assert.Equalf(t, tt.want, got, "writeMarkdownWithFrontmatter() = %v, want %v", got, tt.want)
			}
		})
	}
}
