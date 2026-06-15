package stock

import (
	"testing"
	"time"
)

func TestRootCircuitBreakerExportsStateMachine(t *testing.T) {
	var transitions []string
	breaker := NewCircuitBreaker(CircuitBreakerOptions{
		FailureThreshold: 1,
		ResetTimeout:     time.Hour,
		OnStateChange: func(from CircuitState, to CircuitState) {
			transitions = append(transitions, string(from)+"->"+string(to))
		},
	})

	if breaker.State() != CircuitClosed {
		t.Fatalf("initial state = %s, want CLOSED", breaker.State())
	}
	if !breaker.CanRequest() {
		t.Fatal("closed circuit should allow requests")
	}

	breaker.RecordFailure()
	if breaker.State() != CircuitOpen {
		t.Fatalf("state after failure = %s, want OPEN", breaker.State())
	}
	if breaker.CanRequest() {
		t.Fatal("open circuit should reject requests")
	}
	if len(transitions) != 1 || transitions[0] != "CLOSED->OPEN" {
		t.Fatalf("transitions = %#v, want CLOSED->OPEN", transitions)
	}

	breaker.Reset()
	stats := breaker.Stats()
	if stats.State != CircuitClosed {
		t.Fatalf("state after reset = %s, want CLOSED", stats.State)
	}
	if stats.FailureCount != 0 || stats.HalfOpenSuccessCount != 0 || !stats.LastFailureTime.IsZero() {
		t.Fatalf("stats after reset = %+v, want cleared counters", stats)
	}
}
