package prompts

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
)

type FnCreateFoodParameters struct {
	// A generated description of the food that contains some helpful information to make the user happy
	Description string `json:"description" jsonschema:"required"`
	// Normalized name of the food being created, e.g. chicken parm and vegetables
	Name string `json:"name" jsonschema:"required"`
	// If provided by the user, the energy the food contained, e.g. 500
	Energy float32 `json:"energy" jsonschema:"required"`
	// The energy unit the user provided. If they provided no energy amount, this should be none
	EnegyUnit string `json:"energy_unit" jsonschema:"required,enum=calorie,enum=kilojule,enum=none"`
}

type Set struct {
	// Number of reps performed for this set
	Reps int `json:"reps" jsonschema:"required"`
	// Rest time in seconds for this set
	RestDuration int `json:"rest_duration" jsonschema:"required"`
}

type FnCreateWeightLiftingParameters struct {
	// Type of weight lifting activity, e.g. "squats" or "bench press"
	Activity string `json:"routine" jsonschema:"required"`
	// Weight the user is performing the activity with
	Weight float32 `json:"weight" jsonschema:"required"`
	// The weight unit the user provided
	WeightUnit string `json:"weight_unit" jsonschema:"required,enum=kilogram,enum=pound,enum=none"`
	// Recorded sets the user provided
	Sets []Set `json:"sets" jsonschema:"required"`
	// Notes the user might have about this weight lifting activity
	Notes string `json:"notes" jsonschema:"required"`
}

type FnCreateCardioParameters struct {
	// activity the user is performing, e.g. "walk" or "biycle"
	Activity string `json:"routine" jsonschema:"required"`
	// Duration of the cardio activity in seconds
	Duration int `json:"duration" jsonschema:"required"`
}

type FnGetFoodParameters struct {
	// Optional text match query the user wants
	Query string `json:"query" jsonschema:"required"`
	// Get all food records after this time
	AfterTime string `json:"after_time" jsonschema:"required"`
	// Get all food records before this time
	BeforeTime string `json:"before_time" jsonschema:"required"`
}

func generateMarshaledSchema[T any]() ([]byte, error) {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
		ExpandedStruct:            true,
	}

	if err := reflector.AddGoComments("github.com/calamity-m/reaphur", "./"); err != nil {
		// Panic and exit immediately to force translator to not run
		panic(err)
	}

	var v T
	schema := reflector.Reflect(v)

	return schema.MarshalJSON()
}

func GenerateSchemas() error {

	// Generate the create food parameters
	createFood, err := generateMarshaledSchema[FnCreateFoodParameters]()
	if err != nil {
		return fmt.Errorf("failed to write create food fn")
	}

	// Generate the create weight lifting parameters
	createWeightLifting, err := generateMarshaledSchema[FnCreateWeightLiftingParameters]()
	if err != nil {
		return fmt.Errorf("failed to write create weight lifting fn")
	}

	// Generate the create cardio parameters
	createCardio, err := generateMarshaledSchema[FnCreateCardioParameters]()
	if err != nil {
		return fmt.Errorf("failed to write create cardio fn")
	}

	// Generate the get food parameters
	getFood, err := generateMarshaledSchema[FnGetFoodParameters]()
	if err != nil {
		return fmt.Errorf("failed to write get food fn")
	}

	schemaMap := make(map[string][]byte, 4)
	schemaMap["createfood.json"] = createFood
	schemaMap["createweightlifting.json"] = createWeightLifting
	schemaMap["createcardio.json"] = createCardio
	schemaMap["getfood.json"] = getFood

	return writeArr(schemaMap)

}

func writeArr(bmap map[string][]byte) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	for filename, barr := range bmap {
		path := fmt.Sprintf("%s/central/internal/prompts/generated/%s", wd, filename)
		fmt.Printf(">>> writing %s to %s\n", filename, path)
		if err := os.WriteFile(path, barr, 0777); err != nil {
			return err
		}
	}

	return nil
}
