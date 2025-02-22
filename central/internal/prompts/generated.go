package prompts

import (
	_ "embed"

	"github.com/invopop/jsonschema"
)

var (
	//go:embed generated/createfood.json
	CreateFoodJson string

	CreateFoodProperties = initProperties(CreateFoodJson)
	CreateFoodRequired   = initRequired(CreateFoodJson)

	//go:embed generated/createweightlifting.json
	CreateWeightLiftingJson       string
	CreateWeightLiftingProperties = initProperties(CreateWeightLiftingJson)
	CreateWeightLiftingRequired   = initRequired(CreateWeightLiftingJson)

	//go:embed generated/createcardio.json
	CreateCardioJson       string
	CreateCardioProperties = initProperties(CreateCardioJson)
	CreateCardioRequired   = initRequired(CreateCardioJson)

	//go:embed generated/getfood.json
	GetFoodJson       string
	GetFoodProperties = initProperties(GetFoodJson)
	GetFoodRequired   = initRequired(GetFoodJson)
)

func initProperties(input string) interface{} {
	schema := initSchema(input)

	return schema.Properties
}

func initRequired(input string) interface{} {
	schema := initSchema(input)

	return schema.Required
}

func initSchema(input string) jsonschema.Schema {
	schema := jsonschema.Schema{}
	if err := schema.UnmarshalJSON([]byte(input)); err != nil {
		panic(err)
	}

	return schema
}
