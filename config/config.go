package config

import (
	"os"
	"strings"

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
	LanguageCodeJapanese LanguageCode = "jp"
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
	ApiKey string `yaml:"api_key"`
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
