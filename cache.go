package cache

import (
	"errors"
	"github.com/JansonLv/go-cache/store"
	"github.com/jinzhu/copier"
	"reflect"
	"time"
)

type QCache interface {
	GetCache(key string) *session
	GetCacheWithOptions(key string, value interface{}, getDataFunc func() (interface{}, error), opts ...store.ConfigOption) error
}

type cacheRepository struct {
	cache store.CacheRepository
}

func NewCacheRepository(cache store.CacheRepository) QCache {
	return &cacheRepository{cache: cache}
}

func WithConditionOption(isSave bool) func(config *store.CacheConfig) {
	return func(config *store.CacheConfig) {
		config.IsSave = isSave
	}
}

func WithExpireOption(expire time.Duration) func(config *store.CacheConfig) {
	return func(config *store.CacheConfig) {
		config.Expire = expire
	}
}

func (repo *cacheRepository) GetCacheWithOptions(key string, value interface{}, getDataFunc func() (interface{}, error), opts ...store.ConfigOption) error {
	setFunc, err := repo.cache.GetOrSetDataFunc(key, value, opts...)
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
		return errors.New("getDataFunc is not expected")
	}
	err = copier.Copy(value, data)
	if err != nil {
		return err
	}
	return setFunc(value)
}

// builder模式
type session struct {
	client      store.CacheRepository
	key         string
	isSave      bool
	expire      time.Duration
	getDataFunc func() (interface{}, error)
}

func (repo *cacheRepository) GetCache(key string) *session {
	return newSession(repo.cache, key)
}

func newSession(client store.CacheRepository, key string) *session {
	return &session{
		client:      client,
		key:         key,
		isSave:      true,
		expire:      time.Second * 3,
		getDataFunc: nil,
	}
}

func (s *session) SetIsSave(value bool) *session {
	s.isSave = value
	return s
}

func (s *session) SetExpire(expire time.Duration) *session {
	s.expire = expire
	return s
}

func (s *session) SetGetDataFunc(fn func() (interface{}, error)) *session {
	s.getDataFunc = fn
	return s
}

func (s *session) Find(value interface{}) error {
	setDataFunc, err := s.client.GetOrSetDataFunc(s.key, value, WithConditionOption(s.isSave), WithExpireOption(s.expire))
	if err == nil {
		return nil
	}
	if s.getDataFunc == nil {
		return err
	}
	data, err := s.getDataFunc()
	if err != nil {
		return err
	}
	if reflect.TypeOf(data) != reflect.TypeOf(value) {
		return errors.New("getDataFunc is not expected")
	}
	if err = copier.Copy(value, data); err != nil {
		return err
	}
	return setDataFunc(data)
}
