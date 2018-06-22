package ssr

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotFound is used to indicate a cache hit failure
	ErrNotFound = errors.New("not found")
)

type (
	nocache struct{}
)

func NoCache() Option {
	cache := &nocache{}
	return func(ssr *SSR) {
		ssr.Cache = cache
	}
}

func (c *nocache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, ErrNotFound
}

func (c *nocache) Put(ctx context.Context, key string, expiration time.Duration, data []byte) error {
	return nil
}
