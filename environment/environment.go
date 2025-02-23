package environment

import (
	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Environment struct {
	Client *openai.Client
}

func New(cfg *config.Config) *Environment {
	client := openai.NewClient(option.WithAPIKey(cfg.OpenAI.ApiKey))

	return &Environment{
		Client: client,
	}
}
