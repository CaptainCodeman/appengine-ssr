package ssr

import (
	"context"
	"time"
)

type (
	// CacheProvider allows different cache implementations
	CacheProvider interface {
		Get(ctx context.Context, key string) ([]byte, error)
		Put(ctx context.Context, key string, expiration time.Duration, data []byte) error
	}
)

// Expiration configures the cache expiration
func Expiration(expiration time.Duration) Option {
	return func(ssr *SSR) {
		ssr.Expiration = expiration
	}
}
