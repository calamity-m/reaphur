package parser

import (
	"context"
	"log/slog"

	"github.com/calamity-m/reaphur/pkg/errs"
	"github.com/openai/openai-go"
)

type StructuredOutputRequest struct {
	Schema         interface{}
	Name           string
	Description    string
	DeveloperInput string
	UserInput      string
}

type OpenAIParser struct {
	logger *slog.Logger
	client *openai.Client
}

func CreateGenericStructuredOutputRequest(userInput string) StructuredOutputRequest {
	return StructuredOutputRequest{
		UserInput: userInput,
	}
}

func (oa *OpenAIParser) ActionStructuredOutput(ctx context.Context, r StructuredOutputRequest) ([]byte, error) {
	return nil, errs.ErrNotImplementedYet
}

func NewOpenAIParser(logger *slog.Logger, client *openai.Client) *OpenAIParser {
	return &OpenAIParser{
		logger: logger,
		client: client,
	}
}
