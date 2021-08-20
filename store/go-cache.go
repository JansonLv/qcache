package store

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/patrickmn/go-cache"
	"time"
)

type goCacheClient struct {
	client *cache.Cache
}

func NewGoCacheStore(client *cache.Cache) CacheRepository {
	return &goCacheClient{client: client}
}

func (repo *goCacheClient) Set(_ context.Context, key string, value interface{}, expire time.Duration) error {
	repo.client.Set(key, value, expire)
	return nil
}

func (repo *goCacheClient) Get(_ context.Context, key string, value interface{}) error {
	data, ok := repo.client.Get(key)
	if !ok {
		return KeyNorFound
	}
	return copier.Copy(value, data)
}

func (repo *goCacheClient) Del(_ context.Context, key string) error {
	repo.client.Delete(key)
	return nil
}

func (repo *goCacheClient) Clear(_ context.Context) error {
	return nil
}
