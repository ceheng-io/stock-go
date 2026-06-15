# Proxy Pool Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a thread-safe round-robin proxy pool to the Go SDK so configured proxy URLs rotate on every HTTP attempt, including retries and fallback host attempts.

**Architecture:** Root package options store user-facing proxy pool settings and pass them into `internal/core`. Core owns URL validation, round-robin selection, and HTTP client transport cloning because all provider HTTP requests already flow through `core.Client`. The existing `WithHTTPClient` escape hatch remains intact, with proxy wrapping only when the client transport is nil or a `*http.Transport`.

**Tech Stack:** Go standard library `net/http`, `net/url`, `sync/atomic`, existing Go test suite.

---

### Task 1: Root Proxy Pool Options

**Files:**
- Modify: `options.go`
- Test: `stock_test.go`

- [ ] **Step 1: Write the failing root option tests**

Add these tests to `stock_test.go` before `TestNewAppliesOptions`:

```go
func TestWithProxyPoolStoresURLs(t *testing.T) {
	client := New(WithProxyPool([]string{
		" http://proxy-a.test:8080 ",
		"",
		"http://proxy-b.test:8080",
	}))

	if got, want := client.config.ProxyPool.URLs, []string{"http://proxy-a.test:8080", "http://proxy-b.test:8080"}; !equalStringSlices(got, want) {
		t.Fatalf("ProxyPool.URLs = %#v, want %#v", got, want)
	}
}

func TestWithProxyPoolOptionsCopiesURLs(t *testing.T) {
	urls := []string{"http://proxy-a.test:8080"}
	client := New(WithProxyPoolOptions(ProxyPoolOptions{URLs: urls}))
	urls[0] = "http://mutated.test:8080"

	if got := client.config.ProxyPool.URLs; len(got) != 1 || got[0] != "http://proxy-a.test:8080" {
		t.Fatalf("ProxyPool.URLs = %#v, want immutable copy", got)
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
```

- [ ] **Step 2: Run root tests to verify they fail**

Run:

```bash
go test ./... -run 'TestWithProxyPool'
```

Expected: FAIL because `ProxyPoolOptions`, `WithProxyPool`, `WithProxyPoolOptions`, and `Config.ProxyPool` do not exist.

- [ ] **Step 3: Add public proxy pool options**

In `options.go`, add `strings` to imports:

```go
import (
	"net/http"
	"strings"
	"time"

	"github.com/ceheng.io/stock-go/constants"
	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/useragent"
)
```

Add this field to `Config` near `HTTPClient`:

```go
	ProxyPool                        ProxyPoolOptions
```

Add this type near `RateLimitOptions`:

```go
// ProxyPoolOptions controls round-robin proxy selection for provider requests.
type ProxyPoolOptions struct {
	URLs []string
}
```

Add these option functions after `WithHTTPClient`:

```go
// WithProxyPool configures round-robin proxy URLs for provider requests.
func WithProxyPool(urls []string) Option {
	return WithProxyPoolOptions(ProxyPoolOptions{URLs: urls})
}

// WithProxyPoolOptions configures round-robin proxy selection for provider requests.
func WithProxyPoolOptions(options ProxyPoolOptions) Option {
	return func(config *Config) {
		config.ProxyPool = ProxyPoolOptions{
			URLs: normalizeProxyPoolURLs(options.URLs),
		}
	}
}
```

Add this helper near `cloneStringMap`:

```go
func normalizeProxyPoolURLs(urls []string) []string {
	if len(urls) == 0 {
		return nil
	}
	normalized := make([]string, 0, len(urls))
	for _, rawURL := range urls {
		rawURL = strings.TrimSpace(rawURL)
		if rawURL == "" {
			continue
		}
		normalized = append(normalized, rawURL)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}
```

- [ ] **Step 4: Run root tests to verify they pass**

Run:

```bash
go test ./... -run 'TestWithProxyPool'
```

Expected: PASS.

- [ ] **Step 5: Commit task 1**

Run:

```bash
git add options.go stock_test.go
git commit -m "feat: add proxy pool options"
```

### Task 2: Core Proxy Pool Selector

**Files:**
- Create: `internal/core/proxy_pool.go`
- Test: `internal/core/proxy_pool_test.go`

- [ ] **Step 1: Write failing core selector tests**

Create `internal/core/proxy_pool_test.go`:

```go
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
```

- [ ] **Step 2: Run core selector tests to verify they fail**

Run:

```bash
go test ./internal/core -run 'TestProxyPool|TestNewProxyPool'
```

Expected: FAIL because `NewProxyPool` does not exist.

- [ ] **Step 3: Implement core selector**

Create `internal/core/proxy_pool.go`:

```go
package core

import (
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
)

// ProxyPoolConfig contains round-robin proxy pool settings.
type ProxyPoolConfig struct {
	URLs []string
}

// ProxyPool selects proxy URLs in round-robin order.
type ProxyPool struct {
	urls []*url.URL
	next atomic.Uint64
}

// NewProxyPool creates a round-robin proxy pool from valid proxy URLs.
func NewProxyPool(rawURLs []string) *ProxyPool {
	if len(rawURLs) == 0 {
		return nil
	}
	urls := make([]*url.URL, 0, len(rawURLs))
	for _, rawURL := range rawURLs {
		rawURL = strings.TrimSpace(rawURL)
		if rawURL == "" {
			continue
		}
		proxyURL, err := url.Parse(rawURL)
		if err != nil || proxyURL.Scheme == "" || proxyURL.Host == "" {
			continue
		}
		urls = append(urls, proxyURL)
	}
	if len(urls) == 0 {
		return nil
	}
	return &ProxyPool{urls: urls}
}

// Proxy returns the next proxy URL for an HTTP request.
func (p *ProxyPool) Proxy(_ *http.Request) (*url.URL, error) {
	if p == nil || len(p.urls) == 0 {
		return nil, nil
	}
	index := p.next.Add(1) - 1
	return p.urls[index%uint64(len(p.urls))], nil
}
```

- [ ] **Step 4: Run core selector tests to verify they pass**

Run:

```bash
go test ./internal/core -run 'TestProxyPool|TestNewProxyPool'
```

Expected: PASS.

- [ ] **Step 5: Commit task 2**

Run:

```bash
git add internal/core/proxy_pool.go internal/core/proxy_pool_test.go
git commit -m "feat: add core proxy pool"
```

### Task 3: Wire Proxy Pool Into Core HTTP Client

**Files:**
- Modify: `internal/core/client.go`
- Test: `internal/core/tencent_test.go`

- [ ] **Step 1: Write failing core HTTP client tests**

Add these tests to `internal/core/tencent_test.go` before `TestClientGetTextUsesProviderSpecificRetryPolicy`:

```go
func TestClientProxyPoolRotatesOnRetries(t *testing.T) {
	proxyCalls := []string{}
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyCalls = append(proxyCalls, r.Host)
		if len(proxyCalls) == 1 {
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer proxyServer.Close()

	client := NewClient(Config{
		ProxyPool: ProxyPoolConfig{URLs: []string{
			proxyServer.URL,
			proxyServer.URL,
		}},
		Retry: RetryConfig{MaxRetries: 1, BaseDelay: time.Nanosecond},
	})

	text, err := client.GetText(context.Background(), "http://target.test/path")
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
	if len(proxyCalls) != 2 {
		t.Fatalf("proxy calls = %d, want 2", len(proxyCalls))
	}
}

func TestClientProxyPoolClonesHTTPTransport(t *testing.T) {
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer proxyServer.Close()

	transport := &http.Transport{
		DisableCompression: true,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   2 * time.Second,
	}
	client := NewClient(Config{
		HTTPClient: httpClient,
		ProxyPool:  ProxyPoolConfig{URLs: []string{proxyServer.URL}},
	})

	if client.httpClient == httpClient {
		t.Fatal("client reused original HTTP client, want cloned client")
	}
	if client.httpClient.Timeout != 2*time.Second {
		t.Fatalf("Timeout = %s, want 2s", client.httpClient.Timeout)
	}
	clonedTransport, ok := client.httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport = %T, want *http.Transport", client.httpClient.Transport)
	}
	if clonedTransport == transport {
		t.Fatal("transport was not cloned")
	}
	if !clonedTransport.DisableCompression {
		t.Fatal("DisableCompression was not preserved")
	}
	if clonedTransport.Proxy == nil {
		t.Fatal("Proxy was not configured")
	}
}

func TestClientProxyPoolDoesNotOverrideCustomRoundTripper(t *testing.T) {
	calls := 0
	customTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader("ok")),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})
	httpClient := &http.Client{Transport: customTransport}
	client := NewClient(Config{
		HTTPClient: httpClient,
		ProxyPool:  ProxyPoolConfig{URLs: []string{"http://proxy.test:8080"}},
	})

	text, err := client.GetText(context.Background(), "http://target.test/path")
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
	if calls != 1 {
		t.Fatalf("custom transport calls = %d, want 1", calls)
	}
	if client.httpClient != httpClient {
		t.Fatal("custom round tripper client was replaced")
	}
}
```

- [ ] **Step 2: Run core HTTP client tests to verify they fail**

Run:

```bash
go test ./internal/core -run 'TestClientProxyPool'
```

Expected: FAIL because `Config.ProxyPool` is not wired into `NewClient`.

- [ ] **Step 3: Wire proxy pool into `core.NewClient`**

In `internal/core/client.go`, add this field to `Config` near `HTTPClient`:

```go
	ProxyPool        ProxyPoolConfig
```

In `NewClient`, after the nil `HTTPClient` defaulting block, add:

```go
	config.HTTPClient = httpClientWithProxyPool(config.HTTPClient, NewProxyPool(config.ProxyPool.URLs))
```

Add these helpers near `NewClient`:

```go
func httpClientWithProxyPool(client *http.Client, proxyPool *ProxyPool) *http.Client {
	if proxyPool == nil {
		return client
	}
	if client == nil {
		client = http.DefaultClient
	}
	var transport *http.Transport
	switch current := client.Transport.(type) {
	case nil:
		if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
			transport = defaultTransport.Clone()
		}
	case *http.Transport:
		transport = current.Clone()
	default:
		return client
	}
	if transport == nil {
		return client
	}
	transport.Proxy = proxyPool.Proxy
	cloned := *client
	cloned.Transport = transport
	return &cloned
}
```

- [ ] **Step 4: Run core HTTP client tests to verify they pass**

Run:

```bash
go test ./internal/core -run 'TestClientProxyPool'
```

Expected: PASS.

- [ ] **Step 5: Commit task 3**

Run:

```bash
git add internal/core/client.go internal/core/tencent_test.go
git commit -m "feat: wire proxy pool into core client"
```

### Task 4: Root Integration And Full Verification

**Files:**
- Modify: `stock.go`
- Modify: `stock_test.go`
- Verify: all Go packages

- [ ] **Step 1: Write failing root integration test**

Add this test to `stock_test.go` after `TestWithProxyPoolOptionsCopiesURLs`:

```go
func TestNewWiresProxyPoolToCoreClient(t *testing.T) {
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer proxyServer.Close()

	client := New(WithProxyPool([]string{proxyServer.URL}))
	text, err := client.core.GetText(context.Background(), "http://target.test/path")
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
}
```

Add `net/http/httptest` to the `stock_test.go` imports.

- [ ] **Step 2: Run root integration test to verify it fails if wiring is incomplete**

Run:

```bash
go test . -run 'TestNewWiresProxyPoolToCoreClient'
```

Expected: PASS if Task 1 already passed `ProxyPool` from root into core; otherwise FAIL with a target DNS/network error. If it passes, record that Task 1 wiring already satisfied this integration behavior.

- [ ] **Step 3: Wire root config into core**

In `stock.go`, add this field to the `core.Config` literal near `HTTPClient`:

```go
		ProxyPool: core.ProxyPoolConfig{
			URLs: append([]string(nil), config.ProxyPool.URLs...),
		},
```

- [ ] **Step 4: Run root integration test to verify it passes**

Run:

```bash
go test . -run 'TestNewWiresProxyPoolToCoreClient'
```

Expected: PASS.

- [ ] **Step 5: Run full Go test suite**

Run:

```bash
go test ./...
```

Expected: PASS.

- [ ] **Step 6: Review changed files**

Run:

```bash
git diff --stat
git diff -- options.go stock.go stock_test.go internal/core/client.go internal/core/proxy_pool.go internal/core/proxy_pool_test.go internal/core/tencent_test.go
```

Expected: only proxy pool implementation and tests are changed. Existing unrelated `apps/web/src/pages/Rankings/Rankings.vue` remains unstaged and untouched.

- [ ] **Step 7: Commit task 4**

Run:

```bash
git add stock.go stock_test.go
git commit -m "test: cover root proxy pool integration"
```
