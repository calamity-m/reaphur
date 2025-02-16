package prompts

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/invopop/jsonschema"
)

var (
	CreateFoodParameters          = generateParametersMap[FnCreateFoodParameters]()
	GetFoodParameters             = generateParametersMap[FnGetFoodParameters]()
	CreateWeightLiftingParameters = generateParametersMap[FnCreateWeightLiftingParameters]()
	CreateCardioParameters        = generateParametersMap[FnCreateCardioParameters]()
)

type FnCreateFoodParameters struct {
	// User description of the food being created, e.g. chicken parmi with veggies and some other stuff
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
	// Must be a ISO8601 formatted string. Get all food records after this time
	AfterTime string `json:"after_time" jsonschema:"required"`
	// Must be a ISO8601 formatted string. Get all food records before this time
	BeforeTime string `json:"before_time" jsonschema:"required"`
}

func generateParametersMap[T any]() map[string]interface{} {
	params := make(map[string]interface{})

	//schema := GenerateSchema[T]()
	marshaled, err := generateMarshaledSchema[T]()
	if err != nil {
		panic("failed to generate marshaled schema must exit")
	}

	json.Unmarshal(marshaled, &params)

	return params
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

// Debug function meant for writing out schemas for inspection
func WriteSchemas(filename string, marshaledSchema string) error {

	// Write out response schema
	var buf bytes.Buffer
	_, err := buf.WriteString(marshaledSchema)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, buf.Bytes(), 0777)
}
