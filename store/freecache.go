package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coocood/freecache"
)

type freeCacheClient struct {
	client *freecache.Cache
}

func NewFreeCache(client *freecache.Cache) CacheRepository {
	return &freeCacheClient{client: client}
}

func (repo *freeCacheClient) Get(_ context.Context, key string, value interface{}) error {
	data, err := repo.client.Get([]byte(key))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil
}

func (repo *freeCacheClient) Set(_ context.Context, key string, value interface{}, expire time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return repo.client.Set([]byte(key), v, int(expire/time.Second))
}

func (repo *freeCacheClient) Del(_ context.Context, key string) error {
	repo.client.Del([]byte(key))
	return nil
}

func (repo *freeCacheClient) Clear(_ context.Context) error {
	repo.client.Clear()
	return nil
}
