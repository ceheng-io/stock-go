package cache

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMemoryCacheSetGetDeleteClear(t *testing.T) {
	store := NewMemory()

	store.Set("quote:600519", "maotai", 0)

	value, ok := store.Get("quote:600519")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if value != "maotai" {
		t.Fatalf("value = %v, want maotai", value)
	}

	if !store.Delete("quote:600519") {
		t.Fatal("expected Delete to report existing key was removed")
	}
	if _, ok := store.Get("quote:600519"); ok {
		t.Fatal("expected miss after delete")
	}
	if store.Delete("quote:600519") {
		t.Fatal("expected Delete to report missing key was not removed")
	}

	store.Set("quote:000001", "pingan", 0)
	store.Clear()
	if _, ok := store.Get("quote:000001"); ok {
		t.Fatal("expected miss after clear")
	}
}

func TestMemoryCacheTTLExpires(t *testing.T) {
	store := NewMemory()

	store.Set("short", 1, 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	if _, ok := store.Get("short"); ok {
		t.Fatal("expected expired value to miss")
	}
}

func TestMemoryCacheHasMatchesTypeScript(t *testing.T) {
	store := NewMemory()

	if store.Has("missing") {
		t.Fatal("expected Has to miss unknown key")
	}

	store.Set("present", "value", 0)
	if !store.Has("present") {
		t.Fatal("expected Has to hit present key")
	}

	store.Delete("present")
	if store.Has("present") {
		t.Fatal("expected Has to miss deleted key")
	}

	store.Set("expired", "value", time.Nanosecond)
	time.Sleep(time.Millisecond)
	if store.Has("expired") {
		t.Fatal("expected Has to miss expired key")
	}
	if size := store.Size(); size != 0 {
		t.Fatalf("Size = %d, want 0 after Has cleans expired key", size)
	}
}

func TestMemoryCacheMaxSizeEvictsLRU(t *testing.T) {
	store := NewMemoryWithOptions(Options{MaxSize: 2})

	store.Set("old", "old-value", 0)
	time.Sleep(time.Millisecond)
	store.Set("keep", "keep-value", 0)
	if _, ok := store.Get("old"); !ok {
		t.Fatal("expected old to exist before LRU eviction")
	}
	time.Sleep(time.Millisecond)
	store.Set("new", "new-value", 0)

	if _, ok := store.Get("keep"); ok {
		t.Fatal("expected least recently used key to be evicted")
	}
	if _, ok := store.Get("old"); !ok {
		t.Fatal("expected recently accessed key to remain")
	}
	if _, ok := store.Get("new"); !ok {
		t.Fatal("expected new key to remain")
	}
	if size := store.Size(); size != 2 {
		t.Fatalf("Size = %d, want 2", size)
	}
}

func TestMemoryCacheDefaultMaxSizeMatchesTypeScript(t *testing.T) {
	store := NewMemory()

	for i := 0; i < 1001; i++ {
		store.Set(CreateKey("key", i), i, 0)
	}

	if size := store.Size(); size != 1000 {
		t.Fatalf("Size = %d, want 1000", size)
	}
	if _, ok := store.Get("key:0"); ok {
		t.Fatal("expected oldest key to be evicted by default max size")
	}
	if value, ok := store.Get("key:1000"); !ok || value != 1000 {
		t.Fatalf("newest value = %v/%v, want 1000/true", value, ok)
	}
}

func TestMemoryCacheDefaultTTLAndCleanup(t *testing.T) {
	store := NewMemoryWithOptions(Options{DefaultTTL: 10 * time.Millisecond})
	store.Set("short", "value", 0)
	store.Set("forever", "value", -1)
	time.Sleep(20 * time.Millisecond)

	cleaned := store.Cleanup()
	if cleaned != 1 {
		t.Fatalf("Cleanup = %d, want 1", cleaned)
	}
	if _, ok := store.Get("short"); ok {
		t.Fatal("expected default TTL value to be cleaned")
	}
	if _, ok := store.Get("forever"); !ok {
		t.Fatal("expected negative TTL value to live forever")
	}
}

func TestMemoryCacheClearDropsInflightLikeTypeScript(t *testing.T) {
	store := NewMemory()
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseAll := func() {
		releaseOnce.Do(func() {
			close(release)
		})
	}
	defer releaseAll()

	var calls int32
	fetcher := func(context.Context) (int32, error) {
		call := atomic.AddInt32(&calls, 1)
		started <- struct{}{}
		<-release
		return call, nil
	}

	firstDone := make(chan error, 1)
	go func() {
		_, err := Through(context.Background(), store, "same-key", fetcher, time.Minute)
		firstDone <- err
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("first fetcher did not start")
	}

	store.Clear()

	secondDone := make(chan error, 1)
	go func() {
		_, err := Through(context.Background(), store, "same-key", fetcher, time.Minute)
		secondDone <- err
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("expected Clear to drop in-flight call so a second fetch can start")
	}

	releaseAll()
	for _, done := range []chan error{firstDone, secondDone} {
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(time.Second):
			t.Fatal("Through did not finish")
		}
	}
	if calls != 2 {
		t.Fatalf("fetcher calls = %d, want 2", calls)
	}
}

func TestSharedCacheRegistry(t *testing.T) {
	ClearSharedCaches()
	first := GetSharedCache("quotes", Options{MaxSize: 2})
	second := GetSharedCache("quotes")
	if first != second {
		t.Fatal("GetSharedCache returned different instances for one namespace")
	}

	first.Set("key", "value", 0)
	if value, ok := second.Get("key"); !ok || value != "value" {
		t.Fatalf("shared value = %v/%v, want value/true", value, ok)
	}
	ClearSharedCaches()
	if _, ok := first.Get("key"); ok {
		t.Fatal("expected ClearSharedCaches to clear existing shared caches")
	}
	third := GetSharedCache("quotes")
	if third != first {
		t.Fatal("expected ClearSharedCaches to preserve shared cache registry like TypeScript")
	}
}

func TestCreateKey(t *testing.T) {
	got := CreateKey("quotes", []string{"sh600519", "sz000001"}, map[string]string{
		"period": "daily",
	})
	want := `quotes:["sh600519","sz000001"]:{"period":"daily"}`
	if got != want {
		t.Fatalf("CreateKey = %q, want %q", got, want)
	}
}

func TestCreateKeySkipsNilPartsLikeTypeScript(t *testing.T) {
	got := CreateKey("quotes", nil, "sh600519", true, 1)
	want := "quotes:sh600519:true:1"
	if got != want {
		t.Fatalf("CreateKey = %q, want %q", got, want)
	}
}

func TestTypeScriptStyleCacheNames(t *testing.T) {
	var _ CacheStore = NewMemoryCache()
	var _ *MemoryCache = NewMemoryCache(CacheOptions{MaxSize: 1})
	var _ *MemoryCacheStore = NewMemoryCache()

	store := NewMemoryCache()
	key := CreateCacheKey("quotes", "sh600519")
	calls := 0
	fetcher := func(context.Context) (string, error) {
		calls++
		return "fresh", nil
	}

	first, err := CacheThrough(context.Background(), store, key, fetcher, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	second, err := CacheThrough(context.Background(), store, key, fetcher, time.Minute)
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

func TestMemoryCacheGetOrFetchMatchesTypeScript(t *testing.T) {
	store := NewMemoryCache()
	calls := 0

	fetcher := func(context.Context) (any, error) {
		calls++
		return "fresh", nil
	}

	first, err := store.GetOrFetch(context.Background(), "quote:sh600519", fetcher, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.GetOrFetch(context.Background(), "quote:sh600519", fetcher, time.Minute)
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

func TestThroughCachesValue(t *testing.T) {
	store := NewMemory()
	var calls int32

	fetcher := func(context.Context) (string, error) {
		atomic.AddInt32(&calls, 1)
		return "fresh", nil
	}

	first, err := Through(context.Background(), store, "key", fetcher, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Through(context.Background(), store, "key", fetcher, time.Minute)
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

func TestThroughSingleFlight(t *testing.T) {
	store := NewMemory()
	start := make(chan struct{})
	var calls int32
	var wg sync.WaitGroup
	results := make(chan string, 8)

	fetcher := func(context.Context) (string, error) {
		atomic.AddInt32(&calls, 1)
		<-start
		return "shared", nil
	}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value, err := Through(context.Background(), store, "same-key", fetcher, time.Minute)
			if err != nil {
				t.Errorf("Through returned error: %v", err)
				return
			}
			results <- value
		}()
	}

	for atomic.LoadInt32(&calls) == 0 {
		time.Sleep(time.Millisecond)
	}
	close(start)
	wg.Wait()
	close(results)

	for value := range results {
		if value != "shared" {
			t.Fatalf("value = %q, want shared", value)
		}
	}
	if calls != 1 {
		t.Fatalf("fetcher calls = %d, want 1", calls)
	}
}
