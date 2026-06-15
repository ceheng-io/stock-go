package stock

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRootRateLimiterExportsTokenBucket(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterOptions{
		RequestsPerSecond: 1000,
		MaxBurst:          1,
	})

	if !limiter.TryAcquire() {
		t.Fatal("first token was not available")
	}
	if limiter.TryAcquire() {
		t.Fatal("second token should not be available before refill")
	}
	if wait := limiter.WaitTime(); wait <= 0 {
		t.Fatalf("WaitTime = %s, want positive duration", wait)
	}
	if tokens := limiter.AvailableTokens(); tokens >= 1 {
		t.Fatalf("AvailableTokens = %.2f, want below 1", tokens)
	}
}

func TestRootRateLimiterAcquireHonorsContext(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterOptions{
		RequestsPerSecond: 0.01,
		MaxBurst:          1,
	})

	if err := limiter.Acquire(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	err := limiter.Acquire(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Acquire error = %v, want context deadline exceeded", err)
	}
}
