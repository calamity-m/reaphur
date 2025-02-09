package persistence

import (
	"time"

	"github.com/google/uuid"
)

type FoodRecordEntry struct {
	DbId        int
	Id          uuid.UUID
	UserId      uuid.UUID
	Name        string
	Description string
	KJ          float32
	Grams       float32
	ML          float32
	Created     time.Time
}

type FoodPersistence interface {
	// Create a food record entry
	CreateFood(record FoodRecordEntry) error
	// Retrieve a single food record based on the
	// record's uuid. Internal DB primary key is ignored
	// by this call.
	GetFood(uuid uuid.UUID) (FoodRecordEntry, error)
	// Provided FoodRecordEntry is treated as a filter, allowing
	// the caller to retrieve multiple food records at will.
	GetFoods(filter FoodRecordEntry) ([]FoodRecordEntry, error)
	// Update the record in place
	UpdateFood(record FoodRecordEntry) error
	// Delete matching record
	DeleteFood(uuid uuid.UUID) error
}
