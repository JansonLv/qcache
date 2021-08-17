package store

import (
	"errors"
	"time"
)


var KeyNorFound = errors.New("key are not exist")

type CacheRepository interface {
	Set(key string, value interface{}, expire time.Duration) error
	Get(key string, value interface{}) error
}


