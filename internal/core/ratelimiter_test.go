package core

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterTryAcquireAndWaitTime(t *testing.T) {
	clock := newFakeRateClock(time.Unix(0, 0))
	limiter := newRateLimiterWithClock(RateLimiterOptions{
		RequestsPerSecond: 2,
		MaxBurst:          2,
	}, clock)

	if !limiter.TryAcquire() {
		t.Fatal("first token was not available")
	}
	if !limiter.TryAcquire() {
		t.Fatal("second token was not available")
	}
	if limiter.TryAcquire() {
		t.Fatal("third token should not be available before refill")
	}
	if got := limiter.WaitTime(); got != 500*time.Millisecond {
		t.Fatalf("WaitTime = %s, want 500ms", got)
	}

	clock.Advance(250 * time.Millisecond)
	if got := limiter.WaitTime(); got != 250*time.Millisecond {
		t.Fatalf("WaitTime after partial refill = %s, want 250ms", got)
	}

	clock.Advance(250 * time.Millisecond)
	if !limiter.TryAcquire() {
		t.Fatal("expected one token after 500ms refill")
	}
}

func TestRateLimiterAcquireSleepsUntilTokenIsAvailable(t *testing.T) {
	clock := newFakeRateClock(time.Unix(0, 0))
	limiter := newRateLimiterWithClock(RateLimiterOptions{
		RequestsPerSecond: 4,
		MaxBurst:          1,
	}, clock)

	if err := limiter.Acquire(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := limiter.Acquire(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(clock.sleeps) != 1 || clock.sleeps[0] != 250*time.Millisecond {
		t.Fatalf("sleeps = %v, want [250ms]", clock.sleeps)
	}
}

func TestRateLimiterAcquireHonorsContextCancellation(t *testing.T) {
	clock := newFakeRateClock(time.Unix(0, 0))
	clock.sleepErr = context.Canceled
	limiter := newRateLimiterWithClock(RateLimiterOptions{
		RequestsPerSecond: 1,
		MaxBurst:          1,
	}, clock)

	if err := limiter.Acquire(context.Background()); err != nil {
		t.Fatal(err)
	}
	err := limiter.Acquire(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Acquire error = %v, want context.Canceled", err)
	}
}

func TestRateLimiterAcquireSerializesConcurrentCallers(t *testing.T) {
	clock := newFakeRateClock(time.Unix(0, 0))
	limiter := newRateLimiterWithClock(RateLimiterOptions{
		RequestsPerSecond: 2,
		MaxBurst:          1,
	}, clock)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := limiter.Acquire(context.Background()); err != nil {
				t.Errorf("Acquire returned error: %v", err)
			}
		}()
	}
	wg.Wait()

	if len(clock.sleeps) != 2 {
		t.Fatalf("sleep count = %d, want 2; sleeps=%v", len(clock.sleeps), clock.sleeps)
	}
	for _, sleep := range clock.sleeps {
		if sleep != 500*time.Millisecond {
			t.Fatalf("sleep = %s, want 500ms; sleeps=%v", sleep, clock.sleeps)
		}
	}
	if tokens := limiter.AvailableTokens(); tokens > 0 {
		t.Fatalf("available tokens = %.2f, want <= 0", tokens)
	}
}

type fakeRateClock struct {
	mu       sync.Mutex
	now      time.Time
	sleeps   []time.Duration
	sleepErr error
}

func newFakeRateClock(now time.Time) *fakeRateClock {
	return &fakeRateClock{now: now}
}

func (c *fakeRateClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeRateClock) Sleep(ctx context.Context, duration time.Duration) error {
	c.mu.Lock()
	c.sleeps = append(c.sleeps, duration)
	c.now = c.now.Add(duration)
	err := c.sleepErr
	c.mu.Unlock()
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (c *fakeRateClock) Advance(duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(duration)
}
