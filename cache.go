package cache

import (
	"errors"
	"github.com/JansonLv/go-cache/store"
	"github.com/jinzhu/copier"
	"reflect"
)

type QCache interface {
	GetKey(key string) *session
	GetCacheWithOptions(key string, value interface{}, getDataFunc func() (interface{}, error), opts ...ConfigOption) error
}

type cacheRepository struct {
	cache store.CacheRepository
}

func NewCacheRepository(cache store.CacheRepository) QCache {
	return &cacheRepository{cache: cache}
}

func (repo *cacheRepository) GetCacheWithOptions(key string, value interface{}, getDataFunc func() (interface{}, error), opts ...ConfigOption) error {
	if key == ""{
		return errors.New("key is empty")
	}
	config := NewDefaultCacheConfig()
	for _, opt := range opts {
		opt(config)
	}

	setFunc, err := repo.getOrSetDataFunc(key, value, config)
	if err == nil {
		return nil
	}
	if getDataFunc == nil {
		return err
	}
	data, err := getDataFunc()
	if err != nil {
		return err
	}
	// 类型判断
	if reflect.TypeOf(data) != reflect.TypeOf(value) {
		return errors.New("get data is not expected as value")
	}
	err = copier.Copy(value, data)
	if err != nil {
		return err
	}
	return setFunc(value)
}

func (repo *cacheRepository) GetKey(key string) *session {
	return newSession(repo, key)
}

func (repo *cacheRepository)getOrSetDataFunc(key string, value interface{}, config *cacheConfig) (func(value interface{}) error, error) {
	if !config.isSave {
		return func(value interface{}) error { return nil }, ConditionErr
	}
	err := repo.cache.Get(key, value)
	if err == nil {
		// 获取到数据返回
		return func(value interface{}) error {return nil}, nil
	}
	// 获取数据失败
	return func(value interface{}) error {
		return repo.cache.Set(key, value, config.expire)
	}, err
}
