package ssr

import (
	"context"
	"time"

	"google.golang.org/appengine/memcache"
)

type (
	memcacheProvider struct {
		prefix string
	}
)

// NewSSR creates a new Server Side Rendering middleware
func Memcache(prefix string) Option {
	cache := &memcacheProvider{prefix}
	return func(ssr *SSR) {
		ssr.Cache = cache
	}
}

func (c *memcacheProvider) Get(ctx context.Context, key string) ([]byte, error) {
	item, err := memcache.Get(ctx, c.prefix+key)
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

func (c *memcacheProvider) Put(ctx context.Context, key string, expiration time.Duration, data []byte) error {
	return memcache.Set(ctx, &memcache.Item{
		Key:        c.prefix + key,
		Value:      data,
		Expiration: expiration,
	})
}
