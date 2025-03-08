package config

import (
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

type LanguageCode string

func (l LanguageCode) Name() Language {
	return LanguageCodeToLanguage[l]
}

func (l LanguageCode) String() string {
	return string(l)
}

type LanguageCodes []LanguageCode

func (l LanguageCodes) Strings() []string {
	strs := make([]string, 0, len(l))
	for _, v := range l {
		strs = append(strs, string(v))
	}

	return strs
}

type Language string

func (l Language) String() string {
	return string(l)
}

const (
	LanguageCodeKorean   LanguageCode = "ko"
	LanguageCodeEnglish  LanguageCode = "en"
	LanguageCodeJapanese LanguageCode = "ja"
	LanguageCodeChinese  LanguageCode = "cn"
	LanguageCodeSpanish  LanguageCode = "es"
	LanguageCodeFrench   LanguageCode = "fr"
	LanguageCodeGerman   LanguageCode = "de"
)

const (
	LanguageKorean   Language = "Korean"
	LanguageEnglish  Language = "English"
	LanguageJapanese Language = "Japanese"
	LanguageChinese  Language = "Chinese"
	LanguageSpanish  Language = "Spanish"
	LanguageFrench   Language = "French"
	LanguageGerman   Language = "German"
)

const (
	SimpleTargetPathRule = "{origin}/{fileName}.{language}.md"
)

type LanguageMap map[LanguageCode]Language

func (l LanguageMap) Keys() LanguageCodes {
	keys := make(LanguageCodes, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}

	return keys
}

var (
	LanguageCodeToLanguage = LanguageMap{
		LanguageCodeKorean:   LanguageKorean,
		LanguageCodeEnglish:  LanguageEnglish,
		LanguageCodeJapanese: LanguageJapanese,
		LanguageCodeChinese:  LanguageChinese,
		LanguageCodeSpanish:  LanguageSpanish,
		LanguageCodeFrench:   LanguageFrench,
		LanguageCodeGerman:   LanguageGerman,
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

	// Set default values
	if config.OpenAI.Model == "" {
		config.OpenAI.Model = openai.ChatModelGPT4o
	}

	return &config, nil
}

func (c *Config) validateSimple() error {
	if c.OpenAI.ApiKey == "" {
		return errors.New("api key is required")
	}

	if c.OpenAI.Model == "" {
		return errors.New("model is required")
	}

	if c.Translator.Source.SourceLanguage == "" {
		return errors.New("source language is required")
	}

	if len(c.Translator.Target.TargetLanguages) == 0 {
		return errors.New("target languages are required")
	}

	return nil
}

func bindOriginConfig(cfg *Config, originConfig *Config) {
	if cfg.OpenAI.ApiKey == "" {
		cfg.OpenAI.ApiKey = originConfig.OpenAI.ApiKey
	}

	if cfg.OpenAI.Model == "" {
		cfg.OpenAI.Model = originConfig.OpenAI.Model
	}

	if cfg.Translator.Source.SourceLanguage == "" {
		cfg.Translator.Source.SourceLanguage = originConfig.Translator.Source.SourceLanguage
	}

	if len(cfg.Translator.Target.TargetLanguages) == 0 {
		cfg.Translator.Target.TargetLanguages = originConfig.Translator.Target.TargetLanguages
	}
}

func Simple(cmd *cli.Command) (*Config, error) {
	var (
		cfg             Config
		cfgPath         = cmd.String("config")
		apiKey          = cmd.String("api-key")
		model           = cmd.String("model")
		sourceLanguage  = cmd.String("source-language")
		targetLanguages = cmd.StringSlice("target-languages")
	)
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg.OpenAI.ApiKey = apiKey
	cfg.OpenAI.Model = model
	cfg.Translator.Source.SourceLanguage = LanguageCode(sourceLanguage)
	for _, lang := range targetLanguages {
		if lang == "all" {
			cfg.Translator.Target.TargetLanguages = lo.Without(lo.Keys(LanguageCodeToLanguage), cfg.Translator.Source.SourceLanguage)
			break
		}
		cfg.Translator.Target.TargetLanguages = append(cfg.Translator.Target.TargetLanguages, LanguageCode(lang))
	}
	cfg.Translator.ContentDir = currentDir
	cfg.Translator.Target.TargetPathRule = SimpleTargetPathRule

	if err = cfg.validateSimple(); err != nil && cfgPath == "" {
		return nil, err
	} else if err != nil {
		var (
			originConfig *Config
		)
		originConfig, err = New(cfgPath)
		if err != nil {
			return nil, err
		}

		bindOriginConfig(&cfg, originConfig)
	}

	if err = cfg.validateSimple(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
