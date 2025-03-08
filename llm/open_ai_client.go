package llm

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type RequestOption struct {
	option.RequestOption
}

type OpenAIClient interface {
	New(ctx context.Context, body openai.ChatCompletionNewParams, opts ...RequestOption) (res *openai.ChatCompletion, err error)
}

type openAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(client *openai.Client) OpenAIClient {
	return &openAIClient{
		client: client,
	}
}

func (o openAIClient) New(ctx context.Context, body openai.ChatCompletionNewParams, opts ...RequestOption) (res *openai.ChatCompletion, err error) {
	var chatGptOpts []option.RequestOption

	for _, opt := range opts {
		chatGptOpts = append(chatGptOpts, opt.RequestOption)
	}

	return o.client.Chat.Completions.New(ctx, body, chatGptOpts...)
}
