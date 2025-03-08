package environment

import (
	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/YangTaeyoung/hugo-ai-translator/file"
	"github.com/YangTaeyoung/hugo-ai-translator/llm"
	"github.com/YangTaeyoung/hugo-ai-translator/translator"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Environment struct {
	Translator translator.Translator
	Parser     file.Parser
	Writer     file.Writer
}

func New(cfg *config.Config) *Environment {
	var env Environment

	openaiClient := llm.NewOpenAIClient(openai.NewClient(option.WithAPIKey(cfg.OpenAI.ApiKey)))

	env.Translator = translator.New(openaiClient, translator.Config{
		SourceLanguage:  cfg.Translator.Source.SourceLanguage,
		TargetLanguages: cfg.Translator.Target.TargetLanguages,
		Model:           cfg.OpenAI.Model,
	})
	env.Parser = file.NewParser(file.ParserConfig{
		ContentDir:      cfg.Translator.ContentDir,
		TargetLanguages: cfg.Translator.Target.TargetLanguages,
		TargetPathRule:  cfg.Translator.Target.TargetPathRule,
		IgnoreRules:     cfg.Translator.Source.IgnoreRules,
	})

	env.Writer = file.NewWriter(file.WriterConfig{
		ContentDir:     cfg.Translator.ContentDir,
		TargetPathRule: cfg.Translator.Target.TargetPathRule,
	})

	return &env
}
