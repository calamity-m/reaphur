package srv

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/calamity-m/reaphur/central/internal/mapping"
	"github.com/calamity-m/reaphur/pkg/errs"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/calamity-m/reaphur/proto/v1/domain"
	"github.com/google/uuid"
)

// Simple RPC
//
// Create some food record in the food diary/journal
func (s *CentralServiceServer) CreateFoodRecord(ctx context.Context, r *centralproto.CreateFoodRecordRequest) (*centralproto.CreateFoodRecordResponse, error) {
	s.logger.DebugContext(ctx, "received create food record request", slog.Any("request", r), slog.Any("request_record", r.GetRecord()))

	if err := s.commonServiceValidation(); err != nil {
		return nil, err
	}

	// Map inner record
	wanted, err := mapping.MapDomainFoodRecordToPersistenceFoodRecordEntry(r.GetRecord())
	if err != nil {
		return nil, err
	}

	// Validate description isn't empty
	if wanted.Description == "" {
		return nil, fmt.Errorf("description must not be empty - %w", errs.ErrBadRequest)
	}

	// Generate a UUID id
	if wanted.Id == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("failed to generate id - %w", err)
		}

		wanted.Id = id
	}

	// Ensure a created time is set
	if wanted.Created.IsZero() {
		wanted.Created = time.Now()
	}

	// Create the food item
	err = s.foodStore.CreateFood(wanted)
	if err != nil {
		return nil, err
	}

	// Fetch the recently created food item
	created, err := s.foodStore.GetFood(wanted.Id)
	if err != nil {
		return nil, err
	}

	// Send that single record off as a response
	return &centralproto.CreateFoodRecordResponse{
		Record: mapping.MapPersistenceFoodRecordEntryToDomainFoodRecord(created),
	}, nil
}

// Simple RPC
//
// Fetch some food records from the food diary/journal
func (s *CentralServiceServer) GetFoodRecords(ctx context.Context, r *centralproto.GetFoodRecordsRequest) (*centralproto.GetFoodRecordsResponse, error) {
	if err := s.commonServiceValidation(); err != nil {
		return nil, err
	}

	// Validate and map inner record
	filter, err := mapping.MapCentralProtoFoodFilterToPersistenceFoodFilter(r.GetFilter(), r.GetRequestUserId())
	if err != nil {
		return nil, err
	}

	found, err := s.foodStore.GetFoods(filter)
	if err != nil {
		return nil, err
	}

	records := make([]*domain.FoodRecord, 0, len(found))
	for _, entry := range found {
		records = append(records, mapping.MapPersistenceFoodRecordEntryToDomainFoodRecord(entry))
	}

	return &centralproto.GetFoodRecordsResponse{
		Records: records,
	}, nil
}
