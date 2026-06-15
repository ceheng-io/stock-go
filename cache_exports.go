package stock

import (
	"context"
	"time"

	"github.com/ceheng.io/stock-go/cache"
)

type CacheOptions = cache.Options
type MemoryCache = cache.Memory
type MemoryCacheStore = cache.Memory
type CacheStore = cache.Store

// NewMemoryCache 创建内存缓存。
func NewMemoryCache(options ...CacheOptions) *MemoryCache {
	if len(options) > 0 {
		return cache.NewMemoryWithOptions(options[0])
	}
	return cache.NewMemory()
}

// GetSharedCache 返回具名共享缓存。
func GetSharedCache(namespace string, options ...CacheOptions) *MemoryCache {
	return cache.GetSharedCache(namespace, options...)
}

// ClearSharedCaches 清空所有具名共享缓存内容，并保留 namespace 实例。
func ClearSharedCaches() {
	cache.ClearSharedCaches()
}

// CreateCacheKey 创建稳定的缓存 key。
func CreateCacheKey(parts ...any) string {
	return cache.CreateKey(parts...)
}

// CacheThrough 读取缓存；未命中时执行 fetcher，并用 single-flight 防止并发击穿。
func CacheThrough[T any](
	ctx context.Context,
	store CacheStore,
	key string,
	fetcher func(context.Context) (T, error),
	ttl time.Duration,
) (T, error) {
	return cache.Through(ctx, store, key, fetcher, ttl)
}
