package cache

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/copier"
	"reflect"
	"time"

	"github.com/coocood/freecache"
)

var ConditionErr = errors.New("condition are not met")

type CacheRepository interface {
	SetData(key string, value interface{}) error
	GetData(key string, info interface{}) error
	GetOrSetData(key string, value interface{}) (func(key string, value interface{}, expire int) error, error)
	GetOrSetDataWithCondition(key string, value interface{}, isSave bool) (func(key string, value interface{}, expire int) error, error)
	GetOrSetDataOptions(key string, value interface{}, opts ...configOption) (func(key string, value interface{}, expire int) error, error)
	GetOrSetDataFunc(key string, value interface{}, opts ...configOption) (func(value interface{}) error, error)
	GetWithMarshal(key string, value interface{}, getData func() (interface{}, error), opts ...configOption) error
	GetWithCopier(key string, value interface{}, getData func() (interface{}, error), opts ...configOption) error
}

type cacheRepository struct {
	client *freecache.Cache
}

func NewCacheRepository() CacheRepository {
	client := freecache.NewCache(1024 * 1024)
	return &cacheRepository{client: client}
}

func (repo *cacheRepository) SetData(key string, value interface{}) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return repo.client.Set([]byte(key), v, 10)
}

func (repo *cacheRepository) GetData(key string, info interface{}) error {
	data, err := repo.client.Get([]byte(key))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	return nil
}

func (repo *cacheRepository) GetOrSetData(key string, value interface{}) (func(key string, value interface{}, expire int) error, error) {
	setFunc := func(key string, value interface{}, expire int) error {
		v, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return repo.client.Set([]byte(key), v, expire)
	}
	data, err := repo.client.Get([]byte(key))
	if err != nil {
		return setFunc, err
	}
	err = json.Unmarshal(data, &value)
	if err != nil {
		return setFunc, err
	}
	return setFunc, nil
}

func (repo *cacheRepository) GetOrSetDataWithCondition(key string, value interface{}, isSave bool) (func(key string, value interface{}, expire int) error, error) {
	if !isSave {
		return func(key string, value interface{}, expire int) error { return nil }, errors.New("conditions are not met")
	}
	setFunc := func(key string, value interface{}, expire int) error {
		v, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return repo.client.Set([]byte(key), v, expire)
	}
	data, err := repo.client.Get([]byte(key))
	if err != nil {
		return setFunc, err
	}
	err = json.Unmarshal(data, &value)
	if err != nil {
		return setFunc, err
	}
	return setFunc, nil
}

type cacheConfig struct {
	isSave bool
	expire time.Duration
}

type configOption func(config *cacheConfig)

func WithConditionOption(isSave bool) func(config *cacheConfig) {
	return func(config *cacheConfig) {
		config.isSave = isSave
	}
}

func WithExpireOption(expire time.Duration) func(config *cacheConfig) {
	return func(config *cacheConfig) {
		config.expire = expire
	}
}

func (repo *cacheRepository) GetOrSetDataOptions(key string, value interface{}, opts ...configOption) (func(key string, value interface{}, expire int) error, error) {
	config := &cacheConfig{isSave: true}
	for _, opt := range opts {
		opt(config)
	}
	if !config.isSave {
		return func(key string, value interface{}, expire int) error { return nil }, errors.New("conditions are not met")
	}
	setFunc := func(key string, value interface{}, expire int) error {
		v, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return repo.client.Set([]byte(key), v, expire)
	}
	data, err := repo.client.Get([]byte(key))
	if err != nil {
		return setFunc, err
	}
	err = json.Unmarshal(data, &value)
	if err != nil {
		return setFunc, err
	}
	return setFunc, nil
}

// newDefaultCacheConfig 默认cacheConfig
func newDefaultCacheConfig() *cacheConfig {
	return &cacheConfig{isSave: true, expire: 3 * time.Second}
}

func (repo *cacheRepository) GetOrSetDataFunc(key string, value interface{}, opts ...configOption) (func(value interface{}) error, error) {
	config := newDefaultCacheConfig()
	for _, opt := range opts {
		opt(config)
	}
	if !config.isSave {
		return func(value interface{}) error { return nil }, ConditionErr
	}
	keyBys := []byte(key)

	setFunc := func(value interface{}) error {
		v, err := json.Marshal(value)
		if err != nil {
			return err
		}
		// 闭包使用key和expire
		return repo.client.Set(keyBys, v, int(config.expire/time.Second))
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

func (repo *cacheRepository) GetWithMarshal(key string, value interface{}, getDataFunc func() (interface{}, error), opts ...configOption) error {
	config := newDefaultCacheConfig()
	for _, opt := range opts {
		opt(config)
	}
	getData := func() ([]byte, error) {
		data, err := getDataFunc()
		if err != nil {
			return nil, err
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return bytes, json.Unmarshal(bytes, &value)
	}

	if !config.isSave {
		_, err := getData()
		return err
	}

	keyBys := []byte(key)

	data, err := repo.client.Get(keyBys)
	if err == nil {
		return json.Unmarshal(data, &value)
	}

	v, err := getData()
	if err != nil {
		return err
	}

	return repo.client.Set(keyBys, v, int(config.expire/time.Second))
}

func (repo *cacheRepository) GetWithCopier(key string, value interface{}, getDataFunc func() (interface{}, error), opts ...configOption) error {
	setFunc, err := repo.GetOrSetDataFunc(key, value, opts...)
	if err == nil {
		return nil
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
