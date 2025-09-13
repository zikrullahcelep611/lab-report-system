package redisrepository

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{rdb: rdb}
}

func (r *Repository) SetData(ctx context.Context, data any) (string, error) {
	if data == nil {
		log.Error().Str("operation", "SetData").Msg("Attempted to set nil data")
		return "", errors.New("data is nil")
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		log.Error().Str("operation", "SetData").Err(err).Msg("Failed to encode data with gob")
		return "", errors.New("gob encode errors")
	}

	cacheKey := uuid.New().String()
	err = r.rdb.Set(ctx, cacheKey, buf.Bytes(), 5*time.Minute).Err()
	if err != nil {
		log.Error().Str("operation", "SetData").Err(err).Str("cache_key", cacheKey).
			Msg("Failed to set data in Redis")
		return "", errors.New("redis set errors")
	}

	log.Info().Str("operation", "SetData").Str("cache_key", cacheKey).Msg("Data set in successfully")
	return cacheKey, nil
}

func (r *Repository) GetData(ctx context.Context, cacheKey string, target any) error {
	data, err := r.rdb.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Warn().Str("operation", "GetData").Str("cache_key", cacheKey).
				Msg("No data found for the given cache key")
			return errors.New("redis get error: key does not exist")
		}
		log.Error().Str("operation", "GetData").Err(err).Str("cache_key", cacheKey).
			Msg("Failed to get data from redis")
		return errors.New("redis get errors")
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(target)

	if err != nil {
		log.Error().
			Str("operation", "GetData").
			Err(err).
			Str("cache_key", cacheKey).
			Msg("Failed to decode data with gob")
		return errors.New("gob decode errors")
	}

	log.Info().
		Str("operation", "GetData").
		Str("cache_key", cacheKey).
		Msg("Data retrieved from Redis successfully")

	return nil
}

func (r *Repository) DeleteData(ctx context.Context, cacheKey string) error {
	err := r.rdb.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Error().
			Str("operation", "DeleteData").
			Err(err).
			Str("cache_key", cacheKey).
			Msg("Failed to delete data from Redis")
		return errors.New("redis delete errors")
	}
	log.Info().
		Str("operation", "DeleteData").
		Str("cache_key", cacheKey).
		Msg("Data deleted from Redis successfully")

	return nil
}
