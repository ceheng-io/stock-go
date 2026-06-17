package stock

import "github.com/ceheng-io/stock-go/internal/core"

// CircuitBreakerStats exposes breaker counters for diagnostics.
type CircuitBreakerStats = core.CircuitBreakerStats

// CircuitBreaker protects upstream providers from repeated failing requests.
type CircuitBreaker = core.CircuitBreaker

// NewCircuitBreaker creates a standalone request circuit breaker.
func NewCircuitBreaker(options CircuitBreakerOptions) *CircuitBreaker {
	return core.NewCircuitBreaker(core.CircuitBreakerOptions{
		FailureThreshold: options.FailureThreshold,
		ResetTimeout:     options.ResetTimeout,
		HalfOpenRequests: options.HalfOpenRequests,
		OnStateChange:    options.OnStateChange,
	})
}
