package store

import (
	"errors"
	"time"
)

var ConditionErr = errors.New("condition are not met")
var KeyNorFound = errors.New("key are not exist")

type CacheRepository interface {
	GetOrSetDataFunc(key string, value interface{}, opts ...ConfigOption) (func(value interface{}) error, error) //   只要所有的结构实现这个方法就好了
}

type CacheConfig struct {
	IsSave bool
	Expire time.Duration
}

func newDefaultCacheConfig() *CacheConfig {
	return &CacheConfig{IsSave: true, Expire: 3 * time.Second}
}

type ConfigOption func(config *CacheConfig)
