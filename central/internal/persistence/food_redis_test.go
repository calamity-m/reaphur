package persistence

import "testing"

func TestCreateFoodIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
