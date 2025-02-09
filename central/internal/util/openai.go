package util

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func CreateNewOpenAIClient(token string) *openai.Client {
	client := openai.NewClient(
		option.WithAPIKey(token),
	)

	return client
}
