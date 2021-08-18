package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type goRedisClient struct {
	client redis.UniversalClient
}

func NewGoRedisCache(client redis.UniversalClient) CacheRepository {
	return &goRedisClient{client: client}
}

func (repo *goRedisClient) Get(ctx context.Context, key string, value interface{}) error {
	data, err := repo.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil
}

func (repo *goRedisClient) Set(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return repo.client.Set(ctx, key, v, expire).Err()
}
