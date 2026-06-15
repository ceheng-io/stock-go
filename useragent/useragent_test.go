package useragent

import "testing"

func TestAllReturnsCopy(t *testing.T) {
	values := All()
	if len(values) == 0 {
		t.Fatal("All returned empty list")
	}
	first := values[0]
	values[0] = "mutated"

	again := All()
	if again[0] != first {
		t.Fatalf("All did not return a copy: got %q, want %q", again[0], first)
	}
}

func TestNextRotatesThroughUserAgents(t *testing.T) {
	Reset()
	values := All()
	for i, want := range values {
		if got := Next(); got != want {
			t.Fatalf("Next call %d = %q, want %q", i, got, want)
		}
	}
	if got := Next(); got != values[0] {
		t.Fatalf("Next after wrap = %q, want %q", got, values[0])
	}
}

func TestRandomReturnsKnownUserAgent(t *testing.T) {
	values := All()
	known := make(map[string]struct{}, len(values))
	for _, value := range values {
		known[value] = struct{}{}
	}

	for i := 0; i < 20; i++ {
		if got := Random(); got == "" {
			t.Fatal("Random returned empty user agent")
		} else if _, ok := known[got]; !ok {
			t.Fatalf("Random = %q, want a known user agent", got)
		}
	}
}
