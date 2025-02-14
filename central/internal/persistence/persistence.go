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

type FoodFilter struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	Name        string
	Description string
	BeforeTime  time.Time
	AfterTime   time.Time
}

type FoodPersistence interface {
	// Create a food record entry
	CreateFood(record FoodRecordEntry) error
	// Retrieve a single food record based on the
	// record's uuid.
	GetFood(uuid uuid.UUID) (FoodRecordEntry, error)
	// Provided FoodRecordEntry is treated as a filter, allowing
	// the caller to retrieve multiple food records at will.
	GetFoods(filter FoodFilter) ([]FoodRecordEntry, error)
	// Update the record in place
	UpdateFood(record FoodRecordEntry) error
	// Delete matching record
	DeleteFood(uuid uuid.UUID) error
}
