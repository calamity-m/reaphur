package mapping

import (
	"github.com/calamity-m/reaphur/central/internal/persistence"
	"github.com/calamity-m/reaphur/central/internal/util"
	"github.com/calamity-m/reaphur/pkg/errs"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/calamity-m/reaphur/proto/v1/domain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapCentralProtoFoodFilterToPersistenceFoodFilter(f *centralproto.GetFoodFilter, userId string) (persistence.FoodFilter, error) {
	if f == nil {
		return persistence.FoodFilter{}, errs.ErrBadUserId
	}

	uuidUser, err := uuid.Parse(userId)
	if err != nil {
		return persistence.FoodFilter{}, errs.ErrBadUserId
	}

	return persistence.FoodFilter{
		Id:          util.ParseUUIDRegardless(f.GetId()),
		UserId:      uuidUser,
		Name:        f.GetName(),
		Description: f.GetDescription(),
		BeforeTime:  util.ParseProtoTimestamp(f.GetBeforeTime()),
		AfterTime:   util.ParseProtoTimestamp(f.GetAfterTime()),
	}, nil
}

func MapPersistenceFoodRecordEntryToDomainFoodRecord(entry persistence.FoodRecordEntry) *domain.FoodRecord {
	record := &domain.FoodRecord{
		Id:          entry.Id.String(),
		UserId:      entry.UserId.String(),
		Name:        entry.Name,
		Description: entry.Description,
		Kj:          entry.KJ,
		Grams:       entry.Grams,
		Ml:          entry.ML,
		Calories:    kjToCals(entry.KJ),
		Oz:          gramsToOz(entry.Grams),
		FlOz:        mlToFLOz(entry.ML),
		Time:        timestamppb.New(entry.Created),
	}

	return record
}

func MapDomainFoodRecordToPersistenceFoodRecordEntry(record *domain.FoodRecord) (persistence.FoodRecordEntry, error) {
	if record == nil {
		return persistence.FoodRecordEntry{}, errs.ErrNilNotAllowed
	}

	if _, err := uuid.Parse(record.GetUserId()); err != nil {
		return persistence.FoodRecordEntry{}, errs.ErrBadUserId
	}

	entry := persistence.FoodRecordEntry{
		Name:        record.Name,
		Description: record.Description,
		KJ:          calsToKJ(record.Calories),
		ML:          flOzToML(record.FlOz),
		Grams:       ozToGrams(record.Oz),
		Created:     util.ParseProtoTimestamp(record.GetTime()),
	}

	// Yucky imperial system
	if record.Kj != 0 {
		entry.KJ = record.Kj
	}
	if record.Grams != 0 {
		entry.Grams = record.Grams
	}
	if record.Ml != 0 {
		entry.ML = record.Ml
	}

	entry.Id = util.ParseUUIDRegardless(record.GetId())
	entry.UserId = util.ParseUUIDRegardless(record.GetUserId())

	return entry, nil
}

func calsToKJ(cals float32) float32 {
	return cals * 4.184
}

func ozToGrams(oz float32) float32 {
	return oz * 28.35
}

func flOzToML(flOz float32) float32 {
	return flOz * 29.574
}

func kjToCals(kj float32) float32 {
	return kj / 4.184
}

func gramsToOz(grams float32) float32 {
	return grams / 28.35
}

func mlToFLOz(ml float32) float32 {
	return ml / 29.574
}
