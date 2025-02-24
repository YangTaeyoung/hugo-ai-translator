package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type LanguageCode string

func (l LanguageCode) Name() Language {
	return LanguageMap[l]
}

func (l LanguageCode) String() string {
	return string(l)
}

type LanguageCodes []LanguageCode

type Language string

const (
	LanguageCodeKorean   LanguageCode = "ko"
	LanguageCodeEnglish  LanguageCode = "en"
	LanguageCodeJapanese LanguageCode = "ja"
	LanguageCodeChinese  LanguageCode = "cn"
	LanguageCodeSpanish  LanguageCode = "es"
	LanguageCodeFrench   LanguageCode = "fr"
)

const (
	LanguageKorean   Language = "Korean"
	LanguageEnglish  Language = "English"
	LanguageJapanese Language = "Japanese"
	LanguageChinese  Language = "Chinese"
	LanguageSpanish  Language = "Spanish"
	LanguageFrench   Language = "French"
)

var (
	LanguageMap = map[LanguageCode]Language{
		LanguageCodeKorean:   LanguageKorean,
		LanguageCodeEnglish:  LanguageEnglish,
		LanguageCodeJapanese: LanguageJapanese,
		LanguageCodeChinese:  LanguageChinese,
		LanguageCodeSpanish:  LanguageSpanish,
		LanguageCodeFrench:   LanguageFrench,
	}
)

type TranslatorSourceConfig struct {
	SourceLanguage LanguageCode `yaml:"source_language"`
	IgnoreRules    []string     `yaml:"ignore_rules"`
}

type TranslatorTargetConfig struct {
	TargetLanguages LanguageCodes `yaml:"target_languages"`
	TargetPathRule  string        `yaml:"target_path_rule"`
}

type TranslatorConfig struct {
	ContentDir string                 `yaml:"content_dir"`
	Source     TranslatorSourceConfig `yaml:"source"`
	Target     TranslatorTargetConfig `yaml:"target"`
}

type Config struct {
	OpenAI     OpenAIConfig     `yaml:"openai"`
	Translator TranslatorConfig `yaml:"translator"`
}

type OpenAIConfig struct {
	Model  openai.ChatModel `yaml:"model"`
	ApiKey string           `yaml:"api_key"`
}

func replaceHomeDir(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}

		return homeDir + path[1:]
	}

	return path
}

func New(configPath string) (*Config, error) {
	var config Config

	configFile, err := os.ReadFile(replaceHomeDir(configPath))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	if err = yaml.Unmarshal(configFile, &config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config file")
	}

	config.Translator.ContentDir = replaceHomeDir(config.Translator.ContentDir)

	return &config, nil
}

func TranslatedPaths(historyPath string) ([]string, error) {
	file, err := os.ReadFile(replaceHomeDir(historyPath))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	return strings.Split(string(file), "\n"), nil
}

func WriteTranslatedPaths(historyPath string, translatedPaths ...string) error {
	if err := os.Mkdir(filepath.Dir(replaceHomeDir(historyPath)), os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	f, err := os.Create(replaceHomeDir(historyPath))
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer f.Close()

	for _, p := range translatedPaths {
		if _, err = f.WriteString(p + "\n"); err != nil {
			return errors.Wrap(err, "failed to write string")
		}
	}

	return nil
}

func AppendTranslatedPaths(historyPath string, translatedPaths ...string) error {
	f, err := os.OpenFile(replaceHomeDir(historyPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(historyPath)
		if err != nil {
			return errors.Wrap(err, "failed to create file")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	for _, p := range translatedPaths {
		if _, err = f.WriteString(p + "\n"); err != nil {
			return errors.Wrap(err, "failed to write string")
		}
	}

	return nil
}
