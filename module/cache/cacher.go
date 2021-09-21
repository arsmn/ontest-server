package cache

import (
	"context"
	"time"
)

type Item struct {
	Key   string
	Value interface{}
	TTL   time.Duration
}

type Cacher interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, item *Item) error
	Delete(ctx context.Context, key string) error
}

type Provider interface {
	Cacher() Cacher
}
