package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Store is the injectable cache interface used by the SDK.
type Store interface {
	Get(key string) (any, bool)
	Set(key string, value any, ttl time.Duration)
	Delete(key string) bool
	Clear()
}

// Options configures Memory cache behavior.
type Options struct {
	DefaultTTL time.Duration
	MaxSize    int
}

const defaultMaxSize = 1000

type item struct {
	value      any
	expiresAt  time.Time
	lastAccess time.Time
}

// Memory is an in-memory TTL cache.
type Memory struct {
	mu         sync.RWMutex
	items      map[string]item
	defaultTTL time.Duration
	maxSize    int

	inflightMu sync.Mutex
	inflight   map[string]*call
}

type call struct {
	done  chan struct{}
	value any
	err   error
}

// NewMemory creates an in-memory cache store.
func NewMemory() *Memory {
	return NewMemoryWithOptions(Options{})
}

// NewMemoryWithOptions creates an in-memory cache store with options.
func NewMemoryWithOptions(options Options) *Memory {
	maxSize := options.MaxSize
	if maxSize == 0 {
		maxSize = defaultMaxSize
	}
	return &Memory{
		items:      make(map[string]item),
		defaultTTL: options.DefaultTTL,
		maxSize:    maxSize,
		inflight:   make(map[string]*call),
	}
}

// Get returns a value and whether it exists and has not expired.
func (m *Memory) Get(key string) (any, bool) {
	m.mu.RLock()
	entry, ok := m.items[key]
	m.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		m.Delete(key)
		return nil, false
	}
	entry.lastAccess = time.Now()
	m.mu.Lock()
	m.items[key] = entry
	m.mu.Unlock()
	return entry.value, true
}

// Has reports whether key exists and has not expired.
func (m *Memory) Has(key string) bool {
	_, ok := m.Get(key)
	return ok
}

// Set writes a value. ttl <= 0 means no expiration.
func (m *Memory) Set(key string, value any, ttl time.Duration) {
	var expiresAt time.Time
	effectiveTTL := ttl
	if ttl == 0 {
		effectiveTTL = m.defaultTTL
	}
	if effectiveTTL > 0 {
		expiresAt = time.Now().Add(effectiveTTL)
	}

	m.mu.Lock()
	if m.maxSize > 0 && len(m.items) >= m.maxSize {
		if _, exists := m.items[key]; !exists {
			m.evictLRULocked()
		}
	}
	m.items[key] = item{value: value, expiresAt: expiresAt, lastAccess: time.Now()}
	m.mu.Unlock()
}

// Delete removes a value.
func (m *Memory) Delete(key string) bool {
	m.mu.Lock()
	_, existed := m.items[key]
	delete(m.items, key)
	m.mu.Unlock()
	return existed
}

// Clear removes all values.
func (m *Memory) Clear() {
	m.mu.Lock()
	m.items = make(map[string]item)
	m.mu.Unlock()

	m.inflightMu.Lock()
	m.inflight = make(map[string]*call)
	m.inflightMu.Unlock()
}

// Size returns current cache item count, including expired entries not yet cleaned.
func (m *Memory) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.items)
}

// Cleanup removes expired entries and returns the number of removed items.
func (m *Memory) Cleanup() int {
	now := time.Now()
	cleaned := 0
	m.mu.Lock()
	for key, entry := range m.items {
		if !entry.expiresAt.IsZero() && now.After(entry.expiresAt) {
			delete(m.items, key)
			cleaned++
		}
	}
	m.mu.Unlock()
	return cleaned
}

// GetOrFetch returns a cached value or obtains and stores it.
func (m *Memory) GetOrFetch(
	ctx context.Context,
	key string,
	fetcher func(context.Context) (any, error),
	ttl time.Duration,
) (any, error) {
	return Through[any](ctx, m, key, fetcher, ttl)
}

// CreateKey creates a stable JSON-based cache key.
func CreateKey(parts ...any) string {
	encoded := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == nil {
			continue
		}
		if text, ok := part.(string); ok {
			encoded = append(encoded, text)
			continue
		}
		b, err := json.Marshal(part)
		if err != nil {
			encoded = append(encoded, fmt.Sprint(part))
			continue
		}
		encoded = append(encoded, string(b))
	}
	return strings.Join(encoded, ":")
}

var (
	sharedMu     sync.Mutex
	sharedCaches = make(map[string]*Memory)
)

// GetSharedCache returns a named shared memory cache.
func GetSharedCache(namespace string, options ...Options) *Memory {
	sharedMu.Lock()
	defer sharedMu.Unlock()
	if cache := sharedCaches[namespace]; cache != nil {
		return cache
	}
	config := Options{}
	if len(options) > 0 {
		config = options[0]
	}
	cache := NewMemoryWithOptions(config)
	sharedCaches[namespace] = cache
	return cache
}

// ClearSharedCaches clears all named shared caches while preserving namespace instances.
func ClearSharedCaches() {
	sharedMu.Lock()
	for _, cache := range sharedCaches {
		cache.Clear()
	}
	sharedMu.Unlock()
}

// Through returns a cached value or obtains and stores it with single-flight protection.
func Through[T any](
	ctx context.Context,
	store Store,
	key string,
	fetcher func(context.Context) (T, error),
	ttl time.Duration,
) (T, error) {
	if cached, ok := store.Get(key); ok {
		if value, ok := cached.(T); ok {
			return value, nil
		}
	}

	if memory, ok := store.(*Memory); ok {
		return throughMemory(ctx, memory, key, fetcher, ttl)
	}

	value, err := fetcher(ctx)
	if err != nil {
		var zero T
		return zero, err
	}
	store.Set(key, value, ttl)
	return value, nil
}

func throughMemory[T any](
	ctx context.Context,
	memory *Memory,
	key string,
	fetcher func(context.Context) (T, error),
	ttl time.Duration,
) (T, error) {
	memory.inflightMu.Lock()
	if existing := memory.inflight[key]; existing != nil {
		memory.inflightMu.Unlock()
		select {
		case <-existing.done:
			if existing.err != nil {
				var zero T
				return zero, existing.err
			}
			return existing.value.(T), nil
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}
	}

	current := &call{done: make(chan struct{})}
	memory.inflight[key] = current
	memory.inflightMu.Unlock()

	value, err := fetcher(ctx)
	if err == nil {
		memory.Set(key, value, ttl)
	}

	memory.inflightMu.Lock()
	current.value = value
	current.err = err
	close(current.done)
	if memory.inflight[key] == current {
		delete(memory.inflight, key)
	}
	memory.inflightMu.Unlock()

	return value, err
}

func (m *Memory) evictLRULocked() {
	var oldestKey string
	var oldestTime time.Time
	for key, entry := range m.items {
		if oldestKey == "" || entry.lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.lastAccess
		}
	}
	if oldestKey != "" {
		delete(m.items, oldestKey)
	}
}
