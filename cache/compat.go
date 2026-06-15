package cache

import (
	"context"
	"time"
)

// CacheOptions keeps the TypeScript SDK cache option name available.
type CacheOptions = Options

// MemoryCache keeps the TypeScript SDK memory cache name available.
type MemoryCache = Memory

// MemoryCacheStore keeps the TypeScript SDK memory store name available.
type MemoryCacheStore = Memory

// CacheStore keeps the TypeScript SDK cache store name available.
type CacheStore = Store

// NewMemoryCache creates an in-memory cache store.
func NewMemoryCache(options ...CacheOptions) *MemoryCache {
	if len(options) > 0 {
		return NewMemoryWithOptions(options[0])
	}
	return NewMemory()
}

// CreateCacheKey creates a stable JSON-based cache key.
func CreateCacheKey(parts ...any) string {
	return CreateKey(parts...)
}

// CacheThrough returns a cached value or obtains and stores it.
func CacheThrough[T any](
	ctx context.Context,
	store CacheStore,
	key string,
	fetcher func(context.Context) (T, error),
	ttl time.Duration,
) (T, error) {
	return Through(ctx, store, key, fetcher, ttl)
}
