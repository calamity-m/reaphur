package fncall

import (
	"fmt"

	"github.com/calamity-m/reaphur/central/internal/prompts"
	"github.com/openai/openai-go"
)

const (
	createFoodName          = "log_food"
	createWeightLiftingName = "log_weight_lifting"
	createCardioName        = "log_cardio"

	getFoodName = "get_food"

	failedToolCallMessage = `{"success":false, "message":"tool calling failed"}`
)

func CreateFoodParam() (openai.FunctionDefinitionParam, error) {
	properties, ok := prompts.CreateFoodParameters["properties"]
	if !ok {
		return openai.FunctionDefinitionParam{}, fmt.Errorf("could not find properties in schema")
	}

	required, ok := prompts.CreateFoodParameters["required"]
	if !ok {
		return openai.FunctionDefinitionParam{}, fmt.Errorf("could not find required in schema")
	}

	return openai.FunctionDefinitionParam{
		Name:        openai.String(createFoodName),
		Description: openai.String("log food entry in diary with supplied details"),
		Strict:      openai.Bool(true),
		Parameters: openai.F(openai.FunctionParameters{
			"type":                 "object",
			"properties":           properties,
			"required":             required,
			"additionalProperties": openai.Bool(false),
		}),
	}, nil
}

func CreateWeightLiftingParam() (openai.FunctionDefinitionParam, error) {
	properties, ok := prompts.CreateWeightLiftingParameters["properties"]
	if !ok {
		return openai.FunctionDefinitionParam{}, fmt.Errorf("could not find properties in schema")
	}

	required, ok := prompts.CreateWeightLiftingParameters["required"]
	if !ok {
		return openai.FunctionDefinitionParam{}, fmt.Errorf("could not find required in schema")
	}

	return openai.FunctionDefinitionParam{
		Name:        openai.String(createWeightLiftingName),
		Description: openai.String("log weight lifting session in diary with supplied details"),
		Strict:      openai.Bool(true),
		Parameters: openai.F(openai.FunctionParameters{
			"type":                 "object",
			"properties":           properties,
			"required":             required,
			"additionalProperties": openai.Bool(false),
		}),
	}, nil
}

func CreateCardioParam() (openai.FunctionDefinitionParam, error) {
	properties, ok := prompts.CreateCardioParameters["properties"]
	if !ok {
		return openai.FunctionDefinitionParam{}, fmt.Errorf("could not find properties in schema")
	}

	required, ok := prompts.CreateCardioParameters["required"]
	if !ok {
		return openai.FunctionDefinitionParam{}, fmt.Errorf("could not find required in schema")
	}

	return openai.FunctionDefinitionParam{
		Name:        openai.String(createCardioName),
		Description: openai.String("log cardio workout in diary with supplied details"),
		Strict:      openai.Bool(true),
		Parameters: openai.F(openai.FunctionParameters{
			"type":                 "object",
			"properties":           properties,
			"required":             required,
			"additionalProperties": openai.Bool(false),
		}),
	}, nil
}

func GetFoodParam() (openai.FunctionDefinitionParam, error) {
	properties, ok := prompts.GetFoodParameters["properties"]
	if !ok {
		return openai.FunctionDefinitionParam{}, fmt.Errorf("could not find properties in schema")
	}

	required, ok := prompts.GetFoodParameters["required"]
	if !ok {
		return openai.FunctionDefinitionParam{}, fmt.Errorf("could not find required in schema")
	}

	return openai.FunctionDefinitionParam{
		Name:        openai.String(getFoodName),
		Description: openai.String("retrieves food entries from the diary"),
		Strict:      openai.Bool(true),
		Parameters: openai.F(openai.FunctionParameters{
			"type":                 "object",
			"properties":           properties,
			"required":             required,
			"additionalProperties": openai.Bool(false),
		}),
	}, nil
}

func GetChatCompletionToolParamList() ([]openai.ChatCompletionToolParam, error) {

	createFoodFn, err := CreateFoodParam()
	if err != nil {
		return nil, err
	}

	getFoodFn, err := GetFoodParam()
	if err != nil {
		return nil, err
	}

	createWeightFn, err := CreateWeightLiftingParam()
	if err != nil {
		return nil, err
	}

	createCardioFn, err := CreateCardioParam()
	if err != nil {
		return nil, err
	}

	parr := []openai.ChatCompletionToolParam{
		{
			Type:     openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(createFoodFn),
		},
		{
			Type:     openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(getFoodFn),
		},
		{
			Type:     openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(createWeightFn),
		},
		{
			Type:     openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(createCardioFn),
		},
	}

	return parr, nil
}
