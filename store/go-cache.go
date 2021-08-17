package store

import (
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


func (repo *goCacheClient) Set(key string, value interface{}, expire time.Duration) error {
	repo.client.Set(key, value, expire)
	return nil
}

func (repo *goCacheClient) Get(key string, value interface{}) error {
	data, ok := repo.client.Get(key)
	if !ok {
		return KeyNorFound
	}
	return copier.Copy(value, data)
}
