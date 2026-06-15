package stock

import (
	"context"
	"testing"
	"time"
)

func TestRootReExportsCacheUtilities(t *testing.T) {
	var _ *MemoryCacheStore = NewMemoryCache()

	ClearSharedCaches()
	store := GetSharedCache("root-cache", CacheOptions{MaxSize: 1})
	store.Set("first", "first", 0)
	store.Set("second", "second", 0)

	if _, ok := store.Get("first"); ok {
		t.Fatal("expected first key to be evicted by max size")
	}
	if value, ok := store.Get("second"); !ok || value != "second" {
		t.Fatalf("second value = %v/%v, want second/true", value, ok)
	}
	if key := CreateCacheKey("quotes", "sh600519", 1); key != "quotes:sh600519:1" {
		t.Fatalf("CreateCacheKey = %q", key)
	}

	ttlStore := NewMemoryCache(CacheOptions{DefaultTTL: time.Millisecond})
	ttlStore.Set("ttl", "value", 0)
	time.Sleep(2 * time.Millisecond)
	if cleaned := ttlStore.Cleanup(); cleaned != 1 {
		t.Fatalf("Cleanup = %d, want 1", cleaned)
	}
}

func TestRootCacheThrough(t *testing.T) {
	store := NewMemoryCache()
	calls := 0

	fetcher := func(context.Context) (string, error) {
		calls++
		return "fresh", nil
	}

	first, err := CacheThrough(context.Background(), store, "quote:sh600519", fetcher, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	second, err := CacheThrough(context.Background(), store, "quote:sh600519", fetcher, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	if first != "fresh" || second != "fresh" {
		t.Fatalf("values = %q/%q, want fresh/fresh", first, second)
	}
	if calls != 1 {
		t.Fatalf("fetcher calls = %d, want 1", calls)
	}
}
