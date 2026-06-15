package core

import (
	"errors"
	"sync"
	"time"
)

// CircuitState describes the current circuit breaker state.
type CircuitState string

const (
	CircuitClosed   CircuitState = "CLOSED"
	CircuitOpen     CircuitState = "OPEN"
	CircuitHalfOpen CircuitState = "HALF_OPEN"
)

// ErrCircuitBreakerOpen is returned when an open circuit rejects a request.
var ErrCircuitBreakerOpen = errors.New("circuit breaker is OPEN")

// CircuitBreakerOptions configures failure protection for provider requests.
type CircuitBreakerOptions struct {
	FailureThreshold int
	ResetTimeout     time.Duration
	HalfOpenRequests int
	OnStateChange    func(from CircuitState, to CircuitState)
}

// CircuitBreakerStats exposes breaker counters for diagnostics.
type CircuitBreakerStats struct {
	State                CircuitState
	FailureCount         int
	LastFailureTime      time.Time
	HalfOpenSuccessCount int
}

// CircuitBreaker protects upstream providers from repeated failing requests.
type CircuitBreaker struct {
	mu                   sync.Mutex
	state                CircuitState
	failureCount         int
	lastFailureTime      time.Time
	halfOpenSuccessCount int
	failureThreshold     int
	resetTimeout         time.Duration
	halfOpenRequests     int
	onStateChange        func(from CircuitState, to CircuitState)
	clock                circuitClock
}

type stateTransition struct {
	from     CircuitState
	to       CircuitState
	callback func(from CircuitState, to CircuitState)
}

type circuitClock interface {
	Now() time.Time
}

type realCircuitClock struct{}

func (realCircuitClock) Now() time.Time {
	return time.Now()
}

// NewCircuitBreaker creates a request circuit breaker.
func NewCircuitBreaker(options CircuitBreakerOptions) *CircuitBreaker {
	return newCircuitBreakerWithClock(options, realCircuitClock{})
}

func newCircuitBreakerWithClock(options CircuitBreakerOptions, clock circuitClock) *CircuitBreaker {
	failureThreshold := options.FailureThreshold
	if failureThreshold <= 0 {
		failureThreshold = 5
	}
	resetTimeout := options.ResetTimeout
	if resetTimeout <= 0 {
		resetTimeout = 30 * time.Second
	}
	halfOpenRequests := options.HalfOpenRequests
	if halfOpenRequests <= 0 {
		halfOpenRequests = 1
	}
	if clock == nil {
		clock = realCircuitClock{}
	}
	return &CircuitBreaker{
		state:            CircuitClosed,
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		halfOpenRequests: halfOpenRequests,
		onStateChange:    options.OnStateChange,
		clock:            clock,
	}
}

// State returns the current state after applying time-based transitions.
func (b *CircuitBreaker) State() CircuitState {
	b.mu.Lock()
	transitions := b.checkStateTransitionLocked()
	state := b.state
	b.mu.Unlock()
	notifyStateTransitions(transitions)
	return state
}

// CanRequest reports whether a request can proceed.
func (b *CircuitBreaker) CanRequest() bool {
	b.mu.Lock()
	transitions := b.checkStateTransitionLocked()

	var canRequest bool
	switch b.state {
	case CircuitClosed:
		canRequest = true
	case CircuitOpen:
		canRequest = false
	case CircuitHalfOpen:
		canRequest = b.halfOpenSuccessCount < b.halfOpenRequests
	default:
		canRequest = false
	}
	b.mu.Unlock()
	notifyStateTransitions(transitions)
	return canRequest
}

// RecordSuccess records a successful request and may close a half-open circuit.
func (b *CircuitBreaker) RecordSuccess() {
	b.mu.Lock()
	transitions := b.checkStateTransitionLocked()

	switch b.state {
	case CircuitHalfOpen:
		b.halfOpenSuccessCount++
		if b.halfOpenSuccessCount >= b.halfOpenRequests {
			transitions = append(transitions, b.transitionToLocked(CircuitClosed))
		}
	case CircuitClosed:
		b.failureCount = 0
	}
	b.mu.Unlock()
	notifyStateTransitions(transitions)
}

// RecordFailure records a failed request and may open the circuit.
func (b *CircuitBreaker) RecordFailure() {
	b.mu.Lock()
	b.lastFailureTime = b.clock.Now()

	var transitions []stateTransition
	switch b.state {
	case CircuitHalfOpen:
		transitions = append(transitions, b.transitionToLocked(CircuitOpen))
	case CircuitClosed:
		b.failureCount++
		if b.failureCount >= b.failureThreshold {
			transitions = append(transitions, b.transitionToLocked(CircuitOpen))
		}
	}
	b.mu.Unlock()
	notifyStateTransitions(transitions)
}

// Reset closes the circuit and clears diagnostic counters.
func (b *CircuitBreaker) Reset() {
	b.mu.Lock()
	transition := b.transitionToLocked(CircuitClosed)
	b.failureCount = 0
	b.halfOpenSuccessCount = 0
	b.lastFailureTime = time.Time{}
	b.mu.Unlock()
	notifyStateTransitions([]stateTransition{transition})
}

// Stats returns the current state and counters.
func (b *CircuitBreaker) Stats() CircuitBreakerStats {
	b.mu.Lock()
	transitions := b.checkStateTransitionLocked()
	stats := CircuitBreakerStats{
		State:                b.state,
		FailureCount:         b.failureCount,
		LastFailureTime:      b.lastFailureTime,
		HalfOpenSuccessCount: b.halfOpenSuccessCount,
	}
	b.mu.Unlock()
	notifyStateTransitions(transitions)
	return stats
}

func (b *CircuitBreaker) checkStateTransitionLocked() []stateTransition {
	if b.state != CircuitOpen {
		return nil
	}
	if b.clock.Now().Sub(b.lastFailureTime) >= b.resetTimeout {
		return []stateTransition{b.transitionToLocked(CircuitHalfOpen)}
	}
	return nil
}

func (b *CircuitBreaker) transitionToLocked(newState CircuitState) stateTransition {
	if b.state == newState {
		return stateTransition{}
	}
	oldState := b.state
	b.state = newState

	switch newState {
	case CircuitClosed:
		b.failureCount = 0
		b.halfOpenSuccessCount = 0
	case CircuitHalfOpen:
		b.halfOpenSuccessCount = 0
	}

	return stateTransition{
		from:     oldState,
		to:       newState,
		callback: b.onStateChange,
	}
}

func notifyStateTransitions(transitions []stateTransition) {
	for _, transition := range transitions {
		if transition.callback != nil && transition.from != "" && transition.to != "" {
			transition.callback(transition.from, transition.to)
		}
	}
}
