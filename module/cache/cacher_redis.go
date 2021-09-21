package cache

import (
	"context"
	stderr "errors"
	"time"

	"github.com/arsmn/ontest-server/settings"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type (
	redisDependencies interface {
		settings.Provider
	}
	Redis struct {
		c  *cache.Cache
		dx redisDependencies
	}
)

func NewCacherRedis(dx redisDependencies) *Redis {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
		},
	})

	r := new(Redis)

	r.dx = dx
	r.c = cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(30, time.Minute),
	})

	return r
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	if err := r.c.Get(ctx, key, value); err != nil {
		if stderr.Is(err, cache.ErrCacheMiss) {
			return ErrCacheMissing
		}
		return err
	}
	return nil
}

func (r *Redis) Set(ctx context.Context, item *Item) error {
	return r.c.Set(&cache.Item{
		Ctx:   ctx,
		Key:   item.Key,
		Value: item.Value,
		TTL:   item.TTL,
	})
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.c.Delete(ctx, key)
}
