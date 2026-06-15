package core

import (
	"testing"
	"time"
)

func TestCircuitBreakerOpensAfterFailureThreshold(t *testing.T) {
	clock := newFakeCircuitClock(time.Unix(0, 0))
	var transitions []string
	breaker := newCircuitBreakerWithClock(CircuitBreakerOptions{
		FailureThreshold: 2,
		ResetTimeout:     30 * time.Second,
		OnStateChange: func(from CircuitState, to CircuitState) {
			transitions = append(transitions, string(from)+"->"+string(to))
		},
	}, clock)

	if breaker.State() != CircuitClosed {
		t.Fatalf("initial state = %s, want CLOSED", breaker.State())
	}
	breaker.RecordFailure()
	if breaker.State() != CircuitClosed {
		t.Fatalf("state after first failure = %s, want CLOSED", breaker.State())
	}
	breaker.RecordFailure()

	if breaker.State() != CircuitOpen {
		t.Fatalf("state after threshold failures = %s, want OPEN", breaker.State())
	}
	if breaker.CanRequest() {
		t.Fatal("open circuit should reject requests")
	}
	if len(transitions) != 1 || transitions[0] != "CLOSED->OPEN" {
		t.Fatalf("transitions = %#v, want CLOSED->OPEN", transitions)
	}
}

func TestCircuitBreakerHalfOpenClosesAfterProbeSuccesses(t *testing.T) {
	clock := newFakeCircuitClock(time.Unix(0, 0))
	breaker := newCircuitBreakerWithClock(CircuitBreakerOptions{
		FailureThreshold: 1,
		ResetTimeout:     30 * time.Second,
		HalfOpenRequests: 2,
	}, clock)

	breaker.RecordFailure()
	clock.Advance(29 * time.Second)
	if breaker.CanRequest() {
		t.Fatal("open circuit should reject before reset timeout")
	}

	clock.Advance(time.Second)
	if breaker.State() != CircuitHalfOpen {
		t.Fatalf("state after reset timeout = %s, want HALF_OPEN", breaker.State())
	}
	if !breaker.CanRequest() {
		t.Fatal("half-open circuit should allow probe request")
	}

	breaker.RecordSuccess()
	if breaker.State() != CircuitHalfOpen {
		t.Fatalf("state after first probe success = %s, want HALF_OPEN", breaker.State())
	}
	breaker.RecordSuccess()

	stats := breaker.Stats()
	if stats.State != CircuitClosed {
		t.Fatalf("state after probe successes = %s, want CLOSED", stats.State)
	}
	if stats.FailureCount != 0 || stats.HalfOpenSuccessCount != 0 {
		t.Fatalf("stats after close = %+v, want counters reset", stats)
	}
}

func TestCircuitBreakerHalfOpenFailureReopens(t *testing.T) {
	clock := newFakeCircuitClock(time.Unix(0, 0))
	breaker := newCircuitBreakerWithClock(CircuitBreakerOptions{
		FailureThreshold: 1,
		ResetTimeout:     time.Second,
	}, clock)

	breaker.RecordFailure()
	clock.Advance(time.Second)
	if !breaker.CanRequest() {
		t.Fatal("half-open circuit should allow a probe")
	}
	breaker.RecordFailure()

	if breaker.State() != CircuitOpen {
		t.Fatalf("state after half-open failure = %s, want OPEN", breaker.State())
	}
	if breaker.CanRequest() {
		t.Fatal("reopened circuit should reject requests")
	}
}

func TestCircuitBreakerResetClosesAndClearsCounters(t *testing.T) {
	clock := newFakeCircuitClock(time.Unix(0, 0))
	breaker := newCircuitBreakerWithClock(CircuitBreakerOptions{
		FailureThreshold: 1,
		ResetTimeout:     time.Second,
	}, clock)

	breaker.RecordFailure()
	breaker.Reset()

	stats := breaker.Stats()
	if stats.State != CircuitClosed {
		t.Fatalf("state after reset = %s, want CLOSED", stats.State)
	}
	if stats.FailureCount != 0 || stats.HalfOpenSuccessCount != 0 || !stats.LastFailureTime.IsZero() {
		t.Fatalf("stats after reset = %+v, want counters and last failure cleared", stats)
	}
	if !breaker.CanRequest() {
		t.Fatal("reset circuit should allow requests")
	}
}

func TestCircuitBreakerStateChangeCallbackCanInspectState(t *testing.T) {
	clock := newFakeCircuitClock(time.Unix(0, 0))
	var breaker *CircuitBreaker
	observed := make(chan CircuitState, 1)
	breaker = newCircuitBreakerWithClock(CircuitBreakerOptions{
		FailureThreshold: 1,
		OnStateChange: func(CircuitState, CircuitState) {
			observed <- breaker.State()
		},
	}, clock)

	go breaker.RecordFailure()

	select {
	case state := <-observed:
		if state != CircuitOpen {
			t.Fatalf("callback observed state = %s, want OPEN", state)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("state change callback deadlocked while inspecting state")
	}
}

type fakeCircuitClock struct {
	now time.Time
}

func newFakeCircuitClock(now time.Time) *fakeCircuitClock {
	return &fakeCircuitClock{now: now}
}

func (c *fakeCircuitClock) Now() time.Time {
	return c.now
}

func (c *fakeCircuitClock) Advance(duration time.Duration) {
	c.now = c.now.Add(duration)
}
