package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/calamity-m/reaphur/central/internal/conf"
	"github.com/calamity-m/reaphur/pkg/errs"
	"github.com/calamity-m/reaphur/pkg/serr"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sagikazarmark/slog-shim"
)

type RedisFoodStore struct {
	logger *slog.Logger
	conf   *conf.Config
	rdb    *redis.Client
}

type redisRecord struct {
	I           int       `json:"i" redis:"i"`
	Id          string    `json:"id" redis:"id"`
	UserId      string    `json:"user_id" redis:"user_id"`
	Name        string    `json:"name" redis:"name"`
	Description string    `json:"description" redis:"description"`
	KJ          float32   `json:"kj" redis:"kj"`
	Grams       float32   `json:"gram" redis:"gram"`
	ML          float32   `json:"ml" redis:"ml"`
	Created     time.Time `json:"created" redis:"created"`
}

func mapRecord(record FoodRecordEntry) redisRecord {
	return redisRecord{
		I:           record.DbId,
		Id:          record.Id.String(),
		UserId:      record.UserId.String(),
		Name:        record.Name,
		Description: record.Description,
		KJ:          record.KJ,
		Grams:       record.Grams,
		ML:          record.ML,
		Created:     record.Created,
	}
}

func mapRedis(redis redisRecord) (FoodRecordEntry, error) {
	id, err := uuid.Parse(redis.Id)
	if err != nil {
		return FoodRecordEntry{}, err
	}

	user, err := uuid.Parse(redis.UserId)
	if err != nil {
		return FoodRecordEntry{}, err
	}

	return FoodRecordEntry{
		DbId:        redis.I,
		Id:          id,
		UserId:      user,
		Name:        redis.Name,
		Description: redis.Description,
		KJ:          redis.KJ,
		Grams:       redis.Grams,
		ML:          redis.ML,
		Created:     redis.Created,
	}, nil
}

// Create a food record entry
func (r *RedisFoodStore) CreateFood(record FoodRecordEntry) error {
	if record.Id == uuid.Nil {
		record.Id = uuid.Must(uuid.NewRandom())
	}

	rrec := mapRecord(record)

	ctx := context.Background()

	res := r.rdb.Get(ctx, rrec.Id)
	if res == nil {
		r.logger.ErrorContext(ctx, "encountered nil when checking existing", slog.Any("redis_record", rrec))
		return fmt.Errorf("failed to check existing in redis - %w", errs.ErrBadRequest)
	}
	if res.Err() != redis.Nil {
		return fmt.Errorf("record id already exists - %w", errs.ErrBadId)
	}

	set, err := r.rdb.JSONSet(ctx, fmt.Sprintf("food:%s", record.Id.String()), "$", rrec).Result()
	if err != nil {
		return err
	}

	r.logger.Info("redis created", slog.Any("set", set))

	return nil
}

// Retrieve a single food record based on the
// record's uuid.
func (r *RedisFoodStore) GetFood(uuid uuid.UUID) (FoodRecordEntry, error) {
	query := fmt.Sprintf("@id:(%s)", strings.ReplaceAll(uuid.String(), "-", " "))

	res, err := r.rdb.FTSearchWithArgs(
		context.Background(),
		"idx:food",
		query,
		&redis.FTSearchOptions{
			Limit: 1,
		},
	).Result()

	if err != nil {
		r.logger.Error("encountered err", slog.Any("res", res), slog.Any("uuid", uuid), slog.Any("query", query))
		return FoodRecordEntry{}, err
	}

	r.logger.Debug("Found result", slog.Any("result", res))

	if res.Total != 1 {
		r.logger.Error("did not find entry", slog.Any("res", res), slog.Any("uuid", uuid), slog.Any("query", query))
		return FoodRecordEntry{}, errs.ErrNotFound
	}

	scanned, err := serr.DecodeJSONS[redisRecord](res.Docs[0].Fields["$"])
	if err != nil {
		r.logger.Error("failed scanning document from redis", slog.Any("err", err), slog.Any("res", res))
		return FoodRecordEntry{}, err
	}

	rtn, err := mapRedis(scanned)
	if err != nil {
		r.logger.Error("failed mapping redis to food record", slog.Any("err", err), slog.Any("scanned", scanned))
		return FoodRecordEntry{}, err
	}

	return rtn, nil
}

// Provided FoodRecordEntry is treated as a filter, allowing
// the caller to retrieve multiple food records at will.
func (r *RedisFoodStore) GetFoods(filter FoodFilter) ([]FoodRecordEntry, error) {

	escape := func(s string) string {
		esc := strings.ReplaceAll(s, "-", " ")

		return esc
	}

	// Setup the builder
	var queryBuilder strings.Builder

	// Setup user uuid
	queryBuilder.WriteString(fmt.Sprintf("@user_id:(%s) ", escape(filter.UserId.String())))

	// Setup ID
	if filter.Id != uuid.Nil {
		queryBuilder.WriteString(fmt.Sprintf("@id:(%s) ", escape(filter.Id.String())))
	}

	// Setup name
	if filter.Name != "" {
		queryBuilder.WriteString(fmt.Sprintf("@name:(%s) ", escape(filter.Name)))
	}

	// Setup description
	if filter.Description != "" {
		queryBuilder.WriteString(fmt.Sprintf("@description:(%s) ", escape(filter.Description)))
	}

	// Ignore time filters and run them on returned results. :^)
	query := queryBuilder.String()

	ctx := context.Background()
	res, err := r.rdb.FTSearchWithArgs(
		ctx,
		"idx:food",
		query,
		&redis.FTSearchOptions{},
	).Result()

	if err != nil {
		r.logger.Error("encountered err", slog.Any("res", res), slog.Any("filter", filter))
		return nil, err
	}

	if res.Total == 0 {
		r.logger.Debug("found no results", slog.Any("query", query))
		return nil, errs.ErrNotFound
	}

	results := make([]FoodRecordEntry, res.Total)

	for i, doc := range res.Docs {
		scanned, err := serr.DecodeJSONS[redisRecord](doc.Fields["$"])
		if err != nil {
			r.logger.Error("failed scanning document from redis", slog.Any("err", err), slog.Any("res", res))
			return nil, errs.ErrInternal
		}

		rtn, err := mapRedis(scanned)
		if err != nil {
			r.logger.Error("failed mapping redis to food record", slog.Any("err", err), slog.Any("scanned", scanned))
			return nil, errs.ErrInternal
		}

		results[i] = rtn
	}

	return results, nil
}

// Update the record in place
func (r *RedisFoodStore) UpdateFood(record FoodRecordEntry) error {
	return errs.ErrNotImplementedYet
}

// Delete matching record
func (r *RedisFoodStore) DeleteFood(uuid uuid.UUID) error {
	return errs.ErrNotImplementedYet
}

func NewRedisFoodStore(logger *slog.Logger, conf *conf.Config) (*RedisFoodStore, error) {
	if logger == nil || conf == nil {
		return nil, errs.ErrNilNotAllowed
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.FoodRedisAddress,
		Password: conf.FoodRedisPassword,
		DB:       conf.FoodRedisDB,
		Protocol: 2,
	})

	status := client.Ping(context.Background())

	if status.Err() != nil {
		return nil, fmt.Errorf("failed to connect redit client - %w", status.Err())
	}

	client.FTDropIndex(context.Background(), "idx:food")
	_, err := client.FTCreate(
		context.Background(),
		"idx:food",
		// Options:
		&redis.FTCreateOptions{
			OnJSON: true,
			Prefix: []interface{}{"food:"},
		},
		// Index schema fields:
		&redis.FieldSchema{
			FieldName: "$.id",
			As:        "id",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "$.user_id",
			As:        "user_id",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "$.name",
			As:        "name",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "$.description",
			As:        "description",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "$.created",
			As:        "created",
			FieldType: redis.SearchFieldTypeText,
		},
	).Result()

	if err != nil {
		logger.Error("INDEX CREATION FAILED - IGNORING FOR NOW", slog.Any("err", err))
	}

	rdb := &RedisFoodStore{logger: logger, conf: conf, rdb: client}

	return rdb, nil
}
