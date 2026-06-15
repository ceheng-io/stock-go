package core

import (
	"context"
	"math"
	"sync"
	"time"
)

// RateLimiterOptions configures token-bucket request throttling.
type RateLimiterOptions struct {
	RequestsPerSecond float64
	MaxBurst          float64
}

// RateLimiter is a token-bucket limiter for provider request governance.
type RateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
	clock      rateClock
}

type rateClock interface {
	Now() time.Time
	Sleep(context.Context, time.Duration) error
}

type realRateClock struct{}

func (realRateClock) Now() time.Time {
	return time.Now()
}

func (realRateClock) Sleep(ctx context.Context, duration time.Duration) error {
	return sleepContext(ctx, duration)
}

// NewRateLimiter creates a token-bucket limiter.
func NewRateLimiter(options RateLimiterOptions) *RateLimiter {
	return newRateLimiterWithClock(options, realRateClock{})
}

func newRateLimiterWithClock(options RateLimiterOptions, clock rateClock) *RateLimiter {
	requestsPerSecond := options.RequestsPerSecond
	if requestsPerSecond <= 0 {
		requestsPerSecond = 5
	}
	maxBurst := options.MaxBurst
	if maxBurst <= 0 {
		maxBurst = requestsPerSecond
	}
	if clock == nil {
		clock = realRateClock{}
	}
	return &RateLimiter{
		tokens:     maxBurst,
		maxTokens:  maxBurst,
		refillRate: requestsPerSecond,
		lastRefill: clock.Now(),
		clock:      clock,
	}
}

// TryAcquire attempts to consume one token without waiting.
func (l *RateLimiter) TryAcquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refillLocked()
	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// WaitTime returns how long the caller needs to wait for one token.
func (l *RateLimiter) WaitTime() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refillLocked()
	return l.waitTimeLocked()
}

// Acquire waits until one token can be consumed or the context is canceled.
func (l *RateLimiter) Acquire(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refillLocked()
	wait := l.waitTimeLocked()
	if wait > 0 {
		if err := l.clock.Sleep(ctx, wait); err != nil {
			return err
		}
		l.refillLocked()
	}
	l.tokens--
	return nil
}

// AvailableTokens returns the current token count for diagnostics.
func (l *RateLimiter) AvailableTokens() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refillLocked()
	return l.tokens
}

func (l *RateLimiter) refillLocked() {
	now := l.clock.Now()
	elapsed := now.Sub(l.lastRefill)
	if elapsed <= 0 {
		return
	}
	tokensToAdd := elapsed.Seconds() * l.refillRate
	l.tokens = minFloat(l.maxTokens, l.tokens+tokensToAdd)
	l.lastRefill = now
}

func (l *RateLimiter) waitTimeLocked() time.Duration {
	if l.tokens >= 1 {
		return 0
	}
	tokensNeeded := 1 - l.tokens
	waitSeconds := tokensNeeded / l.refillRate
	return time.Duration(math.Ceil(waitSeconds * float64(time.Second)))
}

func minFloat(left float64, right float64) float64 {
	if left < right {
		return left
	}
	return right
}
