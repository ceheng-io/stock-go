package stock

import "github.com/ceheng-io/stock-go/internal/core"

// RateLimiterOptions configures token-bucket request throttling.
type RateLimiterOptions = core.RateLimiterOptions

// RateLimiter is a token-bucket limiter that can be reused outside Client.
type RateLimiter = core.RateLimiter

// NewRateLimiter creates a standalone token-bucket limiter.
func NewRateLimiter(options RateLimiterOptions) *RateLimiter {
	return core.NewRateLimiter(options)
}
