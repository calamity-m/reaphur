package srv

import (
	"context"
	"fmt"

	"github.com/calamity-m/reaphur/central/internal/mapping"
	"github.com/calamity-m/reaphur/central/internal/persistence"
	"github.com/calamity-m/reaphur/central/internal/util"
	"github.com/calamity-m/reaphur/pkg/errs"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/calamity-m/reaphur/proto/v1/domain"
	"github.com/google/uuid"
)

// Simple RPC
//
// Create some food record in the food diary/journal
func (s *CentralServiceServer) CreateFoodRecord(ctx context.Context, r *centralproto.CreateFoodRecordRequest) (*centralproto.CreateFoodRecordResponse, error) {
	if err := s.commonServiceValidation(); err != nil {
		return nil, err
	}

	// Validate and map inner record
	wanted, err := convertValidDomainFoodRecord(r.GetRecord())
	if err != nil {
		return nil, err
	}

	// Generate a UUID id
	if wanted.Id == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("failed to generate id: %w", err)
		}

		wanted.Id = id
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
		Records: mapping.MapEntryToRecord(created),
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
	filter, err := convertValidGetFoodFilter(r.GetFilter())
	if err != nil {
		return nil, err
	}

	found, err := s.foodStore.GetFoods(filter)
	if err != nil {
		return nil, err
	}

	records := make([]*domain.FoodRecord, 0, len(found))
	for _, entry := range found {
		records = append(records, mapping.MapEntryToRecord(entry))
	}

	return &centralproto.GetFoodRecordsResponse{
		Records: records,
	}, nil
}

func convertValidGetFoodFilter(f *centralproto.GetFoodFilter) (persistence.FoodFilter, error) {

	return persistence.FoodFilter{
		Id:          util.ParseUUIDRegardless(f.GetId()),
		UserId:      util.ParseUUIDRegardless(f.GetUserId()),
		Name:        f.GetName(),
		Description: f.GetDescription(),
		BeforeTime:  f.GetBeforeTime().AsTime(),
		AfterTime:   f.GetAfterTime().AsTime(),
	}, nil
}

func convertValidDomainFoodRecord(fr *domain.FoodRecord) (persistence.FoodRecordEntry, error) {
	conv := persistence.FoodRecordEntry{}

	if fr == nil {
		return conv, errs.ErrNilNotAllowed
	}

	if fr.Id != "" {

	}

	if fr.UserId != "" {

	}

	return conv, errs.ErrNotImplementedYet
}
