package fncall

import (
	"context"
	"log/slog"

	"github.com/calamity-m/reaphur/pkg/errs"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/openai/openai-go"
)

type FnCallOutputRequest struct {
	Functions      []interface{}
	Name           string
	Description    string
	DeveloperInput string
	UserInput      string
}

type OpenAIFnCaller struct {
	logger *slog.Logger
	client *openai.Client
}

func CreateGenericFnCallOutputRequest(userInput string) FnCallOutputRequest {
	return FnCallOutputRequest{
		UserInput: userInput,
	}
}

func (oa *OpenAIFnCaller) EnactUserInput(ctx context.Context, r FnCallOutputRequest, food centralproto.CentralFoodServiceServer) ([]byte, error) {
	return nil, errs.ErrNotImplementedYet
}

func NewOpenAIFnCaller(logger *slog.Logger, client *openai.Client) *OpenAIFnCaller {
	return &OpenAIFnCaller{
		logger: logger,
		client: client,
	}
}
