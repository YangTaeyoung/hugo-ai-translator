package config

import (
	"os"
	"path"
	"testing"

	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	configPath := path.Join(currentDir, "test_config", "config.yaml")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "성공",
			args: args{
				configPath: configPath,
			},
			want: &Config{
				OpenAI: OpenAIConfig{
					Model:  openai.ChatModelGPT4oMini,
					ApiKey: "test-api-key",
				},
				Translator: TranslatorConfig{
					ContentDir: path.Join(homeDir, "hugo-home", "content"),
					Source: TranslatorSourceConfig{
						SourceLanguage: LanguageCodeKorean,
						IgnoreRules: []string{
							"some/ignore/path",
						},
					},
					Target: TranslatorTargetConfig{
						TargetLanguages: LanguageCodes{
							LanguageCodeEnglish,
							LanguageCodeJapanese,
							LanguageCodeFrench,
							LanguageCodeGerman,
						},
						TargetPathRule: "{origin}/{fileName}.{language}.md",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.configPath)
			assert.Equalf(t, tt.wantErr, err != nil, "New() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equalf(t, tt.want, got, "New(%v)", tt.args.configPath)
		})
	}
}
