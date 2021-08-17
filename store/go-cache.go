package store

import (
	"github.com/jinzhu/copier"
	"github.com/patrickmn/go-cache"
)

type goCacheClient struct {
	client *cache.Cache
}

func NewGoCacheStore(client *cache.Cache) CacheRepository {
	return &goCacheClient{client: client}
}

func (repo *goCacheClient) GetOrSetDataFunc(key string, value interface{}, opts ...ConfigOption) (func(value interface{}) error, error) {
	config := newDefaultCacheConfig()
	for _, opt := range opts {
		opt(config)
	}
	if !config.IsSave {
		return func(value interface{}) error { return nil }, ConditionErr
	}

	setFunc := func(value interface{}) error {
		repo.client.Set(key, value, config.Expire)
		return nil
	}
	data, ok := repo.client.Get(key)
	if !ok {
		return setFunc, KeyNorFound
	}
	if err := copier.Copy(value, data); err != nil {
		return setFunc, err
	}
	return setFunc, nil
}
