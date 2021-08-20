package cache

import (
	"context"
	"errors"
	"github.com/JansonLv/qcache/store"
	"github.com/jinzhu/copier"
	"golang.org/x/sync/singleflight"
	"reflect"
	"time"
)

type QCache interface {
	GetKey(key string) *session
	GetCacheWithOptions(ctx context.Context, key string, value interface{}, getDataFunc func() (interface{}, error), opts ...ConfigOption) error
}

type cacheRepository struct {
	sf     *singleflight.Group
	cache  store.CacheRepository
	expire time.Duration
}

type cacheRepoOption func(*cacheRepository)

func WitchCacheExpireOption(expire time.Duration) func(repository *cacheRepository) {
	return func(repository *cacheRepository) {
		repository.expire = expire
	}
}

func NewCacheRepository(cache store.CacheRepository, opts ...cacheRepoOption) QCache {
	repo := &cacheRepository{cache: cache, sf: &singleflight.Group{}, expire: time.Minute}
	for _, opt := range opts {
		opt(repo)
	}
	return repo
}

func (repo *cacheRepository) GetCacheWithOptions(ctx context.Context, key string, value interface{}, getDataFunc func() (interface{}, error), opts ...ConfigOption) error {
	if key == "" {
		return errors.New("key is empty")
	}
	config := newCacheConfig(repo.expire)
	for _, opt := range opts {
		opt(config)
	}

	setFunc, err := repo.getOrSetDataFunc(ctx, key, value, config)
	if err == nil {
		return nil
	}
	if getDataFunc == nil {
		return err
	}

	data, err, _ := repo.sf.Do(key, func() (interface{}, error) {
		data, err := getDataFunc()
		if err != nil {
			return nil, err
		}
		// 类型判断
		if reflect.TypeOf(data) != reflect.TypeOf(value) {
			return nil, errors.New("get data is not expected as value")
		}
		// 防止重复设置缓存
		if err := setFunc(data); err != nil {
			return data, err
		}
		return data, nil
	})
	if data != nil {
		_ = copier.Copy(value, data)
	}
	return err

}

func (repo *cacheRepository) GetKey(key string) *session {
	return newSession(repo, key)
}

func (repo *cacheRepository) getOrSetDataFunc(ctx context.Context, key string, value interface{}, config *cacheConfig) (func(value interface{}) error, error) {
	if !config.isSave {
		return func(value interface{}) error { return nil }, ConditionErr
	}
	err := repo.cache.Get(ctx, key, value)
	if err == nil {
		// 获取到数据返回
		return func(value interface{}) error { return nil }, nil
	}
	// 获取数据失败
	return func(value interface{}) error {
		return repo.cache.Set(ctx, key, value, config.expire)
	}, err
}
