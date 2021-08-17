package store

import (
	"encoding/json"
	"time"

	"github.com/coocood/freecache"
)

type freeCacheClient struct {
	client *freecache.Cache
}

func NewFreeCache(client *freecache.Cache) *freeCacheClient {
	return &freeCacheClient{client: client}
}

func (repo *freeCacheClient) GetOrSetDataFunc(key string, value interface{}, opts ...ConfigOption) (func(value interface{}) error, error) {
	config := newDefaultCacheConfig()
	for _, opt := range opts {
		opt(config)
	}
	if !config.IsSave {
		return func(value interface{}) error { return nil }, ConditionErr
	}
	keyBys := []byte(key)

	setFunc := func(value interface{}) error {
		v, err := json.Marshal(value)
		if err != nil {
			return err
		}
		// 闭包使用key和expire
		return repo.client.Set(keyBys, v, int(config.Expire/time.Second))
	}
	data, err := repo.client.Get(keyBys)
	if err != nil {
		return setFunc, err
	}
	err = json.Unmarshal(data, &value)
	if err != nil {
		return setFunc, err
	}
	return setFunc, nil
}
