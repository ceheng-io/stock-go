package core

import (
	"net/http"
	"sync"
	"testing"
)

func TestProxyPoolRotatesRoundRobin(t *testing.T) {
	pool := NewProxyPool([]string{
		"http://proxy-a.test:8080",
		"http://proxy-b.test:8080",
	})
	if pool == nil {
		t.Fatal("NewProxyPool returned nil")
	}

	req, err := http.NewRequest(http.MethodGet, "https://target.test/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	first, err := pool.Proxy(req)
	if err != nil {
		t.Fatal(err)
	}
	second, err := pool.Proxy(req)
	if err != nil {
		t.Fatal(err)
	}
	third, err := pool.Proxy(req)
	if err != nil {
		t.Fatal(err)
	}

	if first.String() != "http://proxy-a.test:8080" {
		t.Fatalf("first proxy = %q, want proxy-a", first)
	}
	if second.String() != "http://proxy-b.test:8080" {
		t.Fatalf("second proxy = %q, want proxy-b", second)
	}
	if third.String() != "http://proxy-a.test:8080" {
		t.Fatalf("third proxy = %q, want proxy-a", third)
	}
}

func TestNewProxyPoolIgnoresInvalidURLs(t *testing.T) {
	pool := NewProxyPool([]string{
		"",
		"://bad",
		"missing-host",
		"http://proxy-a.test:8080",
	})
	if pool == nil {
		t.Fatal("NewProxyPool returned nil")
	}

	req, err := http.NewRequest(http.MethodGet, "https://target.test/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	proxyURL, err := pool.Proxy(req)
	if err != nil {
		t.Fatal(err)
	}
	if proxyURL.String() != "http://proxy-a.test:8080" {
		t.Fatalf("proxy = %q, want valid URL", proxyURL)
	}
}

func TestNewProxyPoolReturnsNilWhenNoValidURLs(t *testing.T) {
	if pool := NewProxyPool([]string{"", "://bad", "missing-host"}); pool != nil {
		t.Fatalf("pool = %#v, want nil", pool)
	}
}

func TestProxyPoolIsConcurrentSafe(t *testing.T) {
	pool := NewProxyPool([]string{
		"http://proxy-a.test:8080",
		"http://proxy-b.test:8080",
	})
	if pool == nil {
		t.Fatal("NewProxyPool returned nil")
	}
	req, err := http.NewRequest(http.MethodGet, "https://target.test/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	const calls = 100
	results := make(chan string, calls)
	var wg sync.WaitGroup
	for i := 0; i < calls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			proxyURL, err := pool.Proxy(req)
			if err != nil {
				t.Errorf("Proxy returned error: %v", err)
				return
			}
			results <- proxyURL.String()
		}()
	}
	wg.Wait()
	close(results)

	counts := map[string]int{}
	for value := range results {
		counts[value]++
	}
	if counts["http://proxy-a.test:8080"] != calls/2 || counts["http://proxy-b.test:8080"] != calls/2 {
		t.Fatalf("counts = %#v, want even split", counts)
	}
}
