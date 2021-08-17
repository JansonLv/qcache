package cache

import (
	"errors"
	"github.com/JansonLv/go-cache/store"
	"github.com/jinzhu/copier"
	"golang.org/x/sync/singleflight"
	"reflect"
)

type QCache interface {
	GetKey(key string) *session
	GetCacheWithOptions(key string, value interface{}, getDataFunc func() (interface{}, error), opts ...ConfigOption) error
}

type cacheRepository struct {
	cache store.CacheRepository
	sf *singleflight.Group
}

func NewCacheRepository(cache store.CacheRepository) QCache {
	return &cacheRepository{cache: cache, sf: &singleflight.Group{}}
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

	data, err, _ := repo.sf.Do(key ,func() (interface{}, error) {
		data, err := getDataFunc()
		if err != nil {
			return nil, err
		}
		// 类型判断
		if reflect.TypeOf(data) != reflect.TypeOf(value) {
			return nil, errors.New("get data is not expected as value")
		}
		// 防止重复设置缓存
		if err := setFunc(data); err != nil{
			return nil, err
		}
		return data, nil
	})
	if err != nil {
		return err
	}

	return copier.Copy(value, data)
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
