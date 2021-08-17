package cache

import (
	"errors"
	"time"
)

var ConditionErr = errors.New("condition are not met")


type cacheConfig struct {
	isSave bool
	expire time.Duration
}

func NewDefaultCacheConfig() *cacheConfig {
	return &cacheConfig{isSave: true, expire: 3 * time.Second}
}

type ConfigOption func(config *cacheConfig)

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
