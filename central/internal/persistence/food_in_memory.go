package persistence

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/calamity-m/reaphur/pkg/errs"
	"github.com/google/uuid"
)

type MemoryFoodStore struct {
	mux     sync.RWMutex
	entries map[string]FoodRecordEntry
	log     *slog.Logger
}

// Create a food record entry
func (s *MemoryFoodStore) CreateFood(record FoodRecordEntry) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.entries[record.Id.String()] = record

	if s.log != nil {
		// Safety debug logging :)
		s.log.Debug("updated in memory store with a creation", slog.Any("entries", s.entries))
	}

	return nil
}

// Retrieve a single food record based on the
// record's uuid. Internal DB primary key is ignored
// by this call.
func (s *MemoryFoodStore) GetFood(uuid uuid.UUID) (FoodRecordEntry, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	found, ok := s.entries[uuid.String()]

	if !ok {
		return FoodRecordEntry{}, errs.ErrNotFound
	}

	return found, nil
}

// Provided FoodRecordEntry is treated as a filter, allowing
// the caller to retrieve multiple food records at will.
func (s *MemoryFoodStore) GetFoods(filter FoodRecordEntry) ([]FoodRecordEntry, error) {
	entries := make([]FoodRecordEntry, 0)

	s.mux.RLock()
	defer s.mux.RUnlock()
	for _, entry := range s.entries {
		s.log.Debug("found entry", slog.Any("entry", entry))
		// Check ID first as an override
		if filter.Id != uuid.Nil && entry.Id == filter.Id {
			// Forcefully add if we have a matching ID
			entries = append(entries, entry)
			continue
		}

		// Skip non matching user ids
		if entry.UserId != filter.UserId {
			continue
		}

		// Check KJ/Grams/ML only if they're non zero
		if filter.KJ != 0 && filter.KJ != entry.KJ {
			continue
		}

		if filter.Grams != 0 && filter.Grams != entry.Grams {
			continue
		}

		if filter.ML != 0 && filter.Grams != entry.Grams {
			continue
		}

		if !strings.Contains(entry.Name, filter.Name) {
			continue
		}

		if !strings.Contains(entry.Description, filter.Description) {
			continue
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// Update the record in place
func (s *MemoryFoodStore) UpdateFood(record FoodRecordEntry) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.entries[record.Id.String()] = record

	if s.log != nil {
		// Safety debug logging :)
		s.log.Debug("updated in memory store with update", slog.Any("entries", s.entries))
	}

	return nil
}

// Delete matching record
func (s *MemoryFoodStore) DeleteFood(uuid uuid.UUID) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.entries, uuid.String())

	return nil
}

func NewMemoryFoodStore() *MemoryFoodStore {
	entries := make(map[string]FoodRecordEntry, 0)
	return &MemoryFoodStore{entries: entries}
}
