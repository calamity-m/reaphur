package prompts

import (
	"testing"
)

func TestGeneratedParametersMap(t *testing.T) {

	t.Run("create food generated", func(t *testing.T) {
		found, ok := CreateFoodParameters["properties"]
		if !ok {
			t.Fatal("failed to find properties")
		}

		if properties, ok := found.(map[string]interface{}); !ok {
			t.Fatalf("properties is not of the correct type - got %#v but want map[string]interface{}", found)
		} else {
			_, ok := properties["description"]
			if !ok {
				t.Errorf("no found description in properties %#v ", properties)
			}
		}

	})
}
