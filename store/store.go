package store

import (
	"context"
	"errors"
	"time"
)

var KeyNorFound = errors.New("key are not exist")

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
}
