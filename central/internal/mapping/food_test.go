package mapping

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/calamity-m/reaphur/central/internal/persistence"
	"github.com/calamity-m/reaphur/pkg/errs"
	"github.com/calamity-m/reaphur/proto/v1/domain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func fakeTime() time.Time {
	tm, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 4, 2014 at 6:05pm (UTC)")
	return tm
}

func fakeTimestamp() *timestamppb.Timestamp {
	return timestamppb.New(fakeTime())
}

func TestMapDomainFoodRecordToPersistenceFoodRecordEntry(t *testing.T) {
	type args struct {
		Record *domain.FoodRecord
	}
	SuccessTests := []struct {
		Name string
		Args args
		Want persistence.FoodRecordEntry
	}{
		{
			Name: "Time is parsed",
			Args: args{&domain.FoodRecord{
				Time: fakeTimestamp(),
			}},
			Want: persistence.FoodRecordEntry{
				Created: fakeTime(),
			},
		},
		{
			Name: "No ID is zeroed",
			Args: args{&domain.FoodRecord{}},
			Want: persistence.FoodRecordEntry{
				Id: uuid.Nil,
			},
		},
		{
			Name: "Invalid ID is zeroed",
			Args: args{&domain.FoodRecord{
				Id: "bbbbbbbbb",
			}},
			Want: persistence.FoodRecordEntry{
				Id: uuid.Nil,
			},
		},
		{
			Name: "KJ takes precedence",
			Args: args{&domain.FoodRecord{
				Calories: 1000,
				Kj:       10,
			}},
			Want: persistence.FoodRecordEntry{
				KJ: 10,
			},
		},
		{
			Name: "ML takes precedence",
			Args: args{&domain.FoodRecord{
				FlOz: 100,
				Ml:   10,
			}},
			Want: persistence.FoodRecordEntry{
				ML: 10,
			},
		},
		{
			Name: "Grams takes precedence",
			Args: args{&domain.FoodRecord{
				Oz:    100,
				Grams: 10,
			}},
			Want: persistence.FoodRecordEntry{
				Grams: 10,
			},
		},
		{
			Name: "Calories can be used",
			Args: args{&domain.FoodRecord{
				Calories: 1000,
			}},
			Want: persistence.FoodRecordEntry{
				KJ: 4184,
			},
		},
		{
			Name: "Fluid OZ can be used",
			Args: args{&domain.FoodRecord{
				FlOz: 1000,
			}},
			Want: persistence.FoodRecordEntry{
				ML: 29574,
			},
		},
		{
			Name: "Oz can be used",
			Args: args{&domain.FoodRecord{
				Oz: 1000,
			}},
			Want: persistence.FoodRecordEntry{
				Grams: 28350,
			},
		},
	}
	for _, tt := range SuccessTests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Args.Record.UserId = uuid.Nil.String()
			got, err := MapDomainFoodRecordToPersistenceFoodRecordEntry(tt.Args.Record)
			if err != nil {
				t.Errorf("got unexpected err - %v", err)
			}
			if !reflect.DeepEqual(got, tt.Want) {
				t.Errorf("got %v, want %v", got, tt.Want)
			}
		})
	}

	failTests := []struct {
		Name    string
		Input   *domain.FoodRecord
		Want    persistence.FoodRecordEntry
		WantErr error
	}{
		{
			Name:    "nil input is rejected",
			Input:   nil,
			Want:    persistence.FoodRecordEntry{},
			WantErr: errs.ErrNilNotAllowed,
		},
		{
			Name:    "invalid user id is rejected",
			Input:   &domain.FoodRecord{},
			Want:    persistence.FoodRecordEntry{},
			WantErr: errs.ErrBadUserId,
		},
	}

	for _, tt := range failTests {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := MapDomainFoodRecordToPersistenceFoodRecordEntry(tt.Input)

			if !errors.Is(err, tt.WantErr) {
				t.Errorf("got %q error but wanted %q", err, tt.WantErr)
			}

			if !reflect.DeepEqual(got, tt.Want) {
				t.Errorf("got %v but want %v", got, tt.Want)
			}

		})
	}
}

func TestMapPersistenceFoodRecordEntryToDomainFoodRecord(t *testing.T) {
	type args struct {
		entry persistence.FoodRecordEntry
	}
	tests := []struct {
		name string
		args args
		want *domain.FoodRecord
	}{
		{
			name: "Imperial values are created",
			args: args{
				persistence.FoodRecordEntry{
					KJ:      4184,
					ML:      29574,
					Grams:   28350,
					Created: time.Time{},
				},
			},
			want: &domain.FoodRecord{
				Id:       "00000000-0000-0000-0000-000000000000",
				UserId:   "00000000-0000-0000-0000-000000000000",
				Kj:       4184,
				Ml:       29574,
				Grams:    28350,
				Calories: 1000,
				FlOz:     1000,
				Oz:       1000,
				Time:     timestamppb.New(time.Time{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapPersistenceFoodRecordEntryToDomainFoodRecord(tt.args.entry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
