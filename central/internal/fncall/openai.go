package fncall

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/calamity-m/reaphur/central/internal/prompts"
	"github.com/calamity-m/reaphur/pkg/serr"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/calamity-m/reaphur/proto/v1/domain"
	"github.com/openai/openai-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FnCallOutputRequest struct {
	UserId    string `json:"user_id"`
	UserInput string `json:"user_input"`
}

type FnCallOutputResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type OpenAIFnCaller struct {
	logger *slog.Logger
	client *openai.Client
	seed   int64
	model  openai.ChatModel
}

func (oa *OpenAIFnCaller) handleCreateFood(ctx context.Context, fnReq FnCallOutputRequest, args prompts.FnCreateFoodParameters, food centralproto.CentralFoodServiceServer) FnCallOutputResponse {
	rec := &centralproto.CreateFoodRecordRequest{
		Record: &domain.FoodRecord{
			Name:        args.Name,
			Description: args.Description,
			UserId:      fnReq.UserId,
		},
	}

	if args.EnegyUnit == "calorie" {
		rec.Record.Calories = args.Energy
	}
	if args.EnegyUnit == "kilojule" {
		rec.Record.Kj = args.Energy
	}

	created, err := food.CreateFoodRecord(ctx, rec)
	if err != nil {
		return FnCallOutputResponse{
			Success: false,
			Message: "failed to create food record",
		}
	}

	oa.logger.InfoContext(ctx, "created food record", slog.Any("created", created))

	return FnCallOutputResponse{
		Success: true,
		Message: "successfully created food record",
	}
}

func (oa *OpenAIFnCaller) handleGetFood(ctx context.Context, fnReq FnCallOutputRequest, args prompts.FnGetFoodParameters, food centralproto.CentralFoodServiceServer) FnCallOutputResponse {

	before, err := time.Parse("2006-01-02T15:04:05-0700", args.BeforeTime)
	if err != nil {
		oa.logger.ErrorContext(ctx, "failed parsing before time arg", slog.Any("err", err), slog.Any("args", args))
		return FnCallOutputResponse{
			Success: false,
			Message: "failed to get food records",
		}
	}

	after, err := time.Parse("2006-01-02T15:04:05-0700", args.AfterTime)
	if err != nil {
		oa.logger.ErrorContext(ctx, "failed parsing after time arg", slog.Any("err", err), slog.Any("args", args))
		return FnCallOutputResponse{
			Success: false,
			Message: "failed to get food records",
		}
	}

	req := &centralproto.GetFoodRecordsRequest{
		RequestUserId: fnReq.UserId,
		Filter: &centralproto.GetFoodFilter{
			Name:       &args.Query,
			BeforeTime: timestamppb.New(before),
			AfterTime:  timestamppb.New(after),
		},
	}

	found, err := food.GetFoodRecords(ctx, req)
	if err != nil {
		return FnCallOutputResponse{
			Success: false,
			Message: "failed to get food records",
		}
	}

	if len(found.Records) == 0 {
		return FnCallOutputResponse{
			Success: true,
			Message: "no records found with given arguments",
		}
	}

	for _, record := range found.Records {
		oa.logger.InfoContext(ctx, "found record", slog.Any("record", record))
	}

	return FnCallOutputResponse{
		Success: true,
		Message: fmt.Sprintf("successfully found %d food records", len(found.Records)),
	}
}

func (oa *OpenAIFnCaller) EnactUserInput(ctx context.Context, r FnCallOutputRequest, food centralproto.CentralFoodServiceServer) (FnCallOutputResponse, error) {
	if oa.model == "" {
		return FnCallOutputResponse{}, fmt.Errorf("no model selected")
	}

	oa.logger.InfoContext(ctx, "received user input request", slog.Any("request", r))

	tools, err := GetChatCompletionToolParamList()
	if err != nil {
		return FnCallOutputResponse{}, err
	}

	// Create first chat interaction to discover which food input to use
	params := openai.ChatCompletionNewParams{
		Model: openai.F(oa.model),
		Tools: openai.F(tools),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionDeveloperMessageParam{
				Role: openai.F(openai.ChatCompletionDeveloperMessageParamRoleDeveloper),
				Content: openai.F([]openai.ChatCompletionContentPartTextParam{
					openai.TextPart(prompts.CENTRAL_PROMPT),
				}),
			},
			openai.UserMessage(r.UserInput),
		}),
	}

	completion, err := oa.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return FnCallOutputResponse{}, err
	}
	oa.logger.DebugContext(ctx, "completed first tool call completion", slog.Any("completion", completion))

	toolCalls := completion.Choices[0].Message.ToolCalls

	// Ensure we have some function calls
	if len(toolCalls) == 0 {
		oa.logger.ErrorContext(ctx, "no tools were called", slog.Any("completion", completion))
		return FnCallOutputResponse{
			Message: completion.Choices[0].Message.Content,
		}, nil
	}

	// Append the response message from openai to the chain of conversation
	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)

	// Evaluate the functions
	for _, call := range toolCalls {
		switch call.Function.Name {
		case createFoodName:
			// CREATE FOOD FUNCTION
			args, err := serr.DecodeJSONS[prompts.FnCreateFoodParameters](call.Function.Arguments)
			if err != nil {
				oa.logger.ErrorContext(ctx, "failed to decode arguments from tool call", slog.Any("call", call))
				return FnCallOutputResponse{}, err
			}

			record := oa.handleCreateFood(ctx, r, args, food)

			resp, err := serr.EncodeJSON(record)
			if err != nil {
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, failedToolCallMessage))
				break
			}

			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, resp))
		case getFoodName:
			// GET FOOD FUNCTION
			args, err := serr.DecodeJSONS[prompts.FnGetFoodParameters](call.Function.Arguments)
			if err != nil {
				oa.logger.ErrorContext(ctx, "failed to decode arguments from tool call", slog.Any("call", call))
				return FnCallOutputResponse{}, err
			}

			record := oa.handleGetFood(ctx, r, args, food)

			resp, err := serr.EncodeJSON(record)
			if err != nil {
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, failedToolCallMessage))
				break
			}

			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, resp))
		case createWeightLiftingName:
			resp, err := serr.EncodeJSON(FnCallOutputResponse{Success: false, Message: "weight lifting not yet completed, sorry"})
			if err != nil {
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, failedToolCallMessage))
				break
			}

			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, resp))
		case createCardioName:
			resp, err := serr.EncodeJSON(FnCallOutputResponse{Success: false, Message: "cardio logging not yet completed, sorry"})
			if err != nil {
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, failedToolCallMessage))
				break
			}

			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, resp))
		default:
			// UNMATCHED FUNCTIONS
			errResp, err := serr.EncodeJSON(FnCallOutputResponse{Success: false, Message: "unmatched"})
			if err != nil {
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, failedToolCallMessage))
				break
			}

			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(call.ID, errResp))
		}
	}

	oa.logger.DebugContext(ctx, "sending completed params", slog.Any("params", params))

	completion, err = oa.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return FnCallOutputResponse{}, err
	}
	oa.logger.DebugContext(ctx, "completed final completion", slog.Any("completion", completion))

	return FnCallOutputResponse{
		Message: completion.Choices[0].Message.Content,
	}, nil
}

func CreateGenericFnCallOutputRequest(userInput string, userId string) FnCallOutputRequest {

	var inputBuilder strings.Builder

	inputBuilder.WriteString("<extra>")
	inputBuilder.WriteString(fmt.Sprintf("current date: %s", time.Now()))
	inputBuilder.WriteString("</extra>")

	inputBuilder.WriteString("<input>")
	inputBuilder.WriteString(userInput)
	inputBuilder.WriteString("</input>")

	return FnCallOutputRequest{
		UserInput: inputBuilder.String(),
		UserId:    userId,
	}
}

func NewOpenAIFnCaller(logger *slog.Logger, client *openai.Client) *OpenAIFnCaller {
	return &OpenAIFnCaller{
		logger: logger,
		client: client,
		seed:   99,
		model:  openai.ChatModelGPT4oMini,
	}
}
