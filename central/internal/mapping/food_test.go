package mapping

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/calamity-m/reaphur/central/internal/persistence"
	"github.com/calamity-m/reaphur/proto/v1/domain"
	"github.com/google/uuid"
)

// Performs an equals which ignores the created field
func shallowEqualEntry(got persistence.FoodRecordEntry, want persistence.FoodRecordEntry) bool {
	got.Created = time.Time{}
	want.Created = time.Time{}

	return reflect.DeepEqual(got, want)
}

func TestMapRecordToEntryWithoutUuids(t *testing.T) {
	type args struct {
		record *domain.FoodRecord
	}
	tests := []struct {
		name string
		args args
		want persistence.FoodRecordEntry
	}{
		{
			name: "KJ takes precedence",
			args: args{&domain.FoodRecord{
				Calories: 1000,
				Kj:       10,
			}},
			want: persistence.FoodRecordEntry{
				KJ: 10,
			},
		},
		{
			name: "ML takes precedence",
			args: args{&domain.FoodRecord{
				FlOz: 100,
				Ml:   10,
			}},
			want: persistence.FoodRecordEntry{
				ML: 10,
			},
		},
		{
			name: "Grams takes precedence",
			args: args{&domain.FoodRecord{
				Oz:    100,
				Grams: 10,
			}},
			want: persistence.FoodRecordEntry{
				Grams: 10,
			},
		},
		{
			name: "Calories can be used",
			args: args{&domain.FoodRecord{
				Calories: 1000,
			}},
			want: persistence.FoodRecordEntry{
				KJ: 4184,
			},
		},
		{
			name: "Fluid OZ can be used",
			args: args{&domain.FoodRecord{
				FlOz: 1000,
			}},
			want: persistence.FoodRecordEntry{
				ML: 29574,
			},
		},
		{
			name: "Oz can be used",
			args: args{&domain.FoodRecord{
				Oz: 1000,
			}},
			want: persistence.FoodRecordEntry{
				Grams: 28350,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapRecordToEntryWithoutUuids(tt.args.record); !shallowEqualEntry(got, tt.want) {
				t.Errorf("MapRecordToEntryWithoutUuids() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestMapRecordToEntry(t *testing.T) {
	type args struct {
		record *domain.FoodRecord
	}
	tests := []struct {
		name    string
		args    args
		want    persistence.FoodRecordEntry
		wantErr bool
	}{
		{
			name: "id and user id is parsed",
			args: args{
				&domain.FoodRecord{
					Id:     "019458b8-2e00-7663-be30-6cb737d1ab27",
					UserId: "999458b8-2e00-7663-be30-6cb737d1ab27",
				},
			},
			wantErr: false,
			want: persistence.FoodRecordEntry{
				Id:     uuid.MustParse("019458b8-2e00-7663-be30-6cb737d1ab27"),
				UserId: uuid.MustParse("999458b8-2e00-7663-be30-6cb737d1ab27"),
			},
		},
		{
			name: "id is checked",
			args: args{
				&domain.FoodRecord{
					Id: "bob",
				},
			},
			wantErr: true,
		},
		{
			name: "user id is checked",
			args: args{
				&domain.FoodRecord{
					Id:     "019458b8-2e00-7663-be30-6cb737d1ab27",
					UserId: "bob",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapRecordToEntry(tt.args.record)
			if tt.wantErr && err == nil {
				fmt.Println("???")
				t.Fatalf("error = %+v, wantErr %+v", err, tt.wantErr)
			}

			if got.Id != tt.want.Id {
				t.Errorf("got %+v but want %+v", got.Id, tt.want.Id)
			}
			if got.UserId != tt.want.UserId {
				t.Errorf("got %+v but want %+v", got.Id, tt.want.Id)
			}
		})
	}
}

func TestMapEntryToRecord(t *testing.T) {
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
					KJ:    4184,
					ML:    29574,
					Grams: 28350,
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapEntryToRecord(tt.args.entry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapEntryToRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
