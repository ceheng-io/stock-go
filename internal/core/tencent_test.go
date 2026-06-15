package core

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ceheng.io/stock-go/useragent"
)

func TestParseTencentQuoteResponse(t *testing.T) {
	text := `v_s_sh600519="1~贵州茅台~600519~1700.00~-1.23~-0.07~12345~67890~~25000~GP-A"; v_pv_none_match="1";`

	items := ParseTencentQuoteResponse(text)

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if items[0].Key != "s_sh600519" {
		t.Fatalf("key = %q, want s_sh600519", items[0].Key)
	}
	if items[0].Fields[1] != "贵州茅台" {
		t.Fatalf("name = %q, want 贵州茅台", items[0].Fields[1])
	}
	if items[1].Key != "pv_none_match" {
		t.Fatalf("none match key = %q, want pv_none_match", items[1].Key)
	}
}

func TestClientGetTencentQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "s_sh600519" {
			t.Fatalf("q = %q, want s_sh600519", r.URL.Query().Get("q"))
		}
		_, _ = w.Write([]byte(`v_s_sh600519="1~贵州茅台~600519~1700.00~-1.23~-0.07~12345~67890~~25000~GP-A";`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, HTTPClient: server.Client()})
	items, err := client.GetTencentQuote(context.Background(), "s_sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Key != "s_sh600519" {
		t.Fatalf("items = %+v", items)
	}
}

func TestClientGetTencentQuoteDecodesGBKResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		body := append([]byte(`v_s_sh600519="1~`), []byte{0xb9, 0xf3, 0xd6, 0xdd, 0xc3, 0xa9, 0xcc, 0xa8}...)
		body = append(body, []byte(`~600519~1700.00~-1.23~-0.07~12345";`)...)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, HTTPClient: server.Client()})
	items, err := client.GetTencentQuote(context.Background(), "s_sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if got := items[0].Fields[1]; got != "贵州茅台" {
		t.Fatalf("name = %q, want 贵州茅台", got)
	}
}

func TestClientGetText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" || r.URL.Query().Get("q") != "maotai" {
			t.Fatalf("unexpected request URL %s", r.URL.String())
		}
		_, _ = w.Write([]byte(`v_hint="sh~600519~\u8d35\u5dde\u8305\u53f0~GZMT~GP-A";`))
	}))
	defer server.Close()

	client := NewClient(Config{HTTPClient: server.Client()})
	text, err := client.GetText(context.Background(), server.URL+"/search?q=maotai")
	if err != nil {
		t.Fatal(err)
	}
	if text == "" {
		t.Fatal("GetText returned empty text")
	}
}

func TestClientGetTextRotatesUserAgentWhenEnabled(t *testing.T) {
	useragent.Reset()
	values := useragent.All()
	seen := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("User-Agent"))
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient:      server.Client(),
		UserAgent:       "fixed-agent",
		RotateUserAgent: true,
	})
	for i := 0; i < 2; i++ {
		if _, err := client.GetText(context.Background(), server.URL); err != nil {
			t.Fatal(err)
		}
	}

	if len(seen) != 2 {
		t.Fatalf("seen user agents = %#v, want 2 entries", seen)
	}
	if seen[0] != values[0] || seen[1] != values[1] {
		t.Fatalf("seen user agents = %#v, want first two rotated values", seen)
	}
}

func TestClientProviderPolicyCanDisableUserAgentRotation(t *testing.T) {
	useragent.Reset()
	disableRotation := false
	seen := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("User-Agent"))
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient:      server.Client(),
		UserAgent:       "fixed-agent",
		RotateUserAgent: true,
		ProviderPolicies: map[ProviderName]ProviderPolicy{
			ProviderUnknown: {
				RotateUserAgent: &disableRotation,
			},
		},
	})
	if _, err := client.GetText(context.Background(), server.URL); err != nil {
		t.Fatal(err)
	}

	if len(seen) != 1 || seen[0] != "fixed-agent" {
		t.Fatalf("seen user agents = %#v, want fixed-agent", seen)
	}
}

func TestClientGetJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"list":["sh600000","sz000001"]}`))
	}))
	defer server.Close()

	client := NewClient(Config{HTTPClient: server.Client()})
	var payload struct {
		List []string `json:"list"`
	}
	if err := client.GetJSON(context.Background(), server.URL, &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.List) != 2 || payload.List[0] != "sh600000" {
		t.Fatalf("payload = %+v", payload)
	}
}

func TestClientGetJSONParseErrorReturnsCodedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`<html>not json</html>`))
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient: server.Client(),
		Retry:      RetryConfig{MaxRetries: 0, BaseDelay: time.Nanosecond},
	})
	var payload map[string]any
	err := client.GetJSON(context.Background(), server.URL, &payload)
	if err == nil {
		t.Fatal("expected parse error")
	}
	var coded CodedError
	if !errors.As(err, &coded) {
		t.Fatalf("err = %T %[1]v, want CodedError", err)
	}
	if coded.SDKCode() != "PARSE_ERROR" {
		t.Fatalf("SDKCode = %q, want PARSE_ERROR", coded.SDKCode())
	}
}

func TestClientGetTextRetriesRetryableHTTPStatus(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			http.Error(w, "temporary", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient: server.Client(),
		Retry:      RetryConfig{MaxRetries: 2, BaseDelay: time.Nanosecond},
	})
	text, err := client.GetText(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" || attempts != 3 {
		t.Fatalf("text=%q attempts=%d", text, attempts)
	}
}

func TestClientGetTextDoesNotRetryNonRetryableHTTPStatus(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient: server.Client(),
		Retry:      RetryConfig{MaxRetries: 2, BaseDelay: time.Nanosecond},
	})
	_, err := client.GetText(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestClientGetTextHTTPStatusReturnsCodedError(t *testing.T) {
	for _, tc := range []struct {
		status int
		code   string
	}{
		{status: http.StatusInternalServerError, code: "HTTP_ERROR"},
		{status: http.StatusTooManyRequests, code: "RATE_LIMITED"},
	} {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, http.StatusText(tc.status), tc.status)
			}))
			defer server.Close()

			client := NewClient(Config{
				HTTPClient: server.Client(),
				Retry:      RetryConfig{MaxRetries: 0, BaseDelay: time.Nanosecond},
			})

			_, err := client.GetText(context.Background(), server.URL)
			if err == nil {
				t.Fatal("expected error")
			}
			var coded CodedError
			if !errors.As(err, &coded) {
				t.Fatalf("err = %T %[1]v, want CodedError", err)
			}
			if coded.SDKCode() != tc.code {
				t.Fatalf("SDKCode = %q, want %q", coded.SDKCode(), tc.code)
			}
		})
	}
}

func TestClientGetTextDrainsHTTPErrorResponseBeforeClose(t *testing.T) {
	body := &drainTrackingBody{content: []byte("temporary upstream failure")}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Status:     http.StatusText(http.StatusServiceUnavailable),
				Body:       body,
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})},
		Retry: RetryConfig{MaxRetries: 0, BaseDelay: time.Nanosecond},
	})

	_, err := client.GetText(context.Background(), "https://drain.test/path")
	if err == nil {
		t.Fatal("expected HTTP status error")
	}
	if !body.closed {
		t.Fatal("response body was not closed")
	}
	if !body.closedAfterEOF {
		t.Fatal("response body closed before it was drained to EOF")
	}
}

func TestClientGetTextRetriesNetworkErrors(t *testing.T) {
	transport := &flakyRoundTripper{
		responses: []roundTripResult{
			{err: errors.New("temporary network failure")},
			{body: "ok", status: http.StatusOK},
		},
	}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: transport},
		Retry:      RetryConfig{MaxRetries: 1, BaseDelay: time.Nanosecond},
	})

	text, err := client.GetText(context.Background(), "https://retry.test/path")
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" || transport.calls != 2 {
		t.Fatalf("text=%q calls=%d", text, transport.calls)
	}
}

func TestClientGetTextNetworkErrorReturnsCodedError(t *testing.T) {
	transport := &flakyRoundTripper{
		responses: []roundTripResult{
			{err: errors.New("temporary network failure")},
		},
	}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: transport},
		Retry:      RetryConfig{MaxRetries: 0, BaseDelay: time.Nanosecond},
	})

	_, err := client.GetText(context.Background(), "https://network-error.test/path")
	if err == nil {
		t.Fatal("expected error")
	}
	var coded CodedError
	if !errors.As(err, &coded) {
		t.Fatalf("err = %T %[1]v, want CodedError", err)
	}
	if coded.SDKCode() != "NETWORK_ERROR" {
		t.Fatalf("SDKCode = %q, want NETWORK_ERROR", coded.SDKCode())
	}
}

func TestClientGetTextDoesNotRetryCanceledContext(t *testing.T) {
	transport := &flakyRoundTripper{
		responses: []roundTripResult{
			{err: context.Canceled},
			{body: "unexpected", status: http.StatusOK},
		},
	}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: transport},
		Retry:      RetryConfig{MaxRetries: 1, BaseDelay: time.Nanosecond},
	})

	_, err := client.GetText(context.Background(), "https://cancel.test/path")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if transport.calls != 1 {
		t.Fatalf("calls = %d, want 1", transport.calls)
	}
}

func TestClientGetTextTimeoutReturnsCodedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(20 * time.Millisecond)
		_, _ = w.Write([]byte("late"))
	}))
	defer server.Close()
	client := NewClient(Config{
		HTTPClient: server.Client(),
		Timeout:    time.Nanosecond,
		Retry:      RetryConfig{MaxRetries: 0, BaseDelay: time.Nanosecond},
	})

	_, err := client.GetText(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected error")
	}
	var coded CodedError
	if !errors.As(err, &coded) {
		t.Fatalf("err = %T %[1]v, want CodedError", err)
	}
	if coded.SDKCode() != "TIMEOUT" {
		t.Fatalf("SDKCode = %q, want TIMEOUT", coded.SDKCode())
	}
}

func TestClientGetTextAcquiresRateLimiterBeforeRequests(t *testing.T) {
	limiter := &countingLimiter{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient:  server.Client(),
		RateLimiter: limiter,
	})

	if _, err := client.GetText(context.Background(), server.URL); err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetText(context.Background(), server.URL); err != nil {
		t.Fatal(err)
	}
	if limiter.calls != 2 {
		t.Fatalf("Acquire calls = %d, want 2", limiter.calls)
	}
}

func TestClientGetTextRejectsWhenCircuitBreakerIsOpen(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerOptions{FailureThreshold: 1, ResetTimeout: time.Minute})
	breaker.RecordFailure()
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient:     server.Client(),
		CircuitBreaker: breaker,
	})
	_, err := client.GetText(context.Background(), server.URL)

	if !errors.Is(err, ErrCircuitBreakerOpen) {
		t.Fatalf("error = %v, want ErrCircuitBreakerOpen", err)
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
}

func TestClientGetTextRecordsCircuitBreakerSuccessAfterRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			http.Error(w, "temporary", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerOptions{FailureThreshold: 1, ResetTimeout: time.Minute})
	client := NewClient(Config{
		HTTPClient:     server.Client(),
		Retry:          RetryConfig{MaxRetries: 1, BaseDelay: time.Nanosecond},
		CircuitBreaker: breaker,
	})

	text, err := client.GetText(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
	stats := breaker.Stats()
	if stats.State != CircuitClosed || stats.FailureCount != 0 {
		t.Fatalf("stats after retry success = %+v, want closed with failure count reset", stats)
	}
}

func TestClientGetTextRecordsCircuitBreakerFailureAfterRetryExhausted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "temporary", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerOptions{FailureThreshold: 1, ResetTimeout: time.Minute})
	client := NewClient(Config{
		HTTPClient:     server.Client(),
		Retry:          RetryConfig{MaxRetries: 1, BaseDelay: time.Nanosecond},
		CircuitBreaker: breaker,
	})

	_, err := client.GetText(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected error")
	}
	if breaker.State() != CircuitOpen {
		t.Fatalf("circuit state = %s, want OPEN", breaker.State())
	}
}

func TestClientGetTextFallsBackAcrossEastmoneyHosts(t *testing.T) {
	transport := &hostSwitchRoundTripper{
		statusByHost: map[string]int{
			"push2his.eastmoney.com": http.StatusServiceUnavailable,
		},
		bodyByHost: map[string]string{
			"7.push2his.eastmoney.com": "ok",
		},
	}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: transport},
		Retry:      RetryConfig{MaxRetries: 0, BaseDelay: time.Nanosecond},
	})

	text, err := client.GetText(context.Background(), "https://push2his.eastmoney.com/api/qt/stock/kline/get?secid=1.600519")
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
	if len(transport.hosts) != 2 {
		t.Fatalf("hosts = %#v, want two attempts", transport.hosts)
	}
	if transport.hosts[0] != "push2his.eastmoney.com" || transport.hosts[1] != "7.push2his.eastmoney.com" {
		t.Fatalf("hosts = %#v, want original then fallback", transport.hosts)
	}
}

func TestClientGetTextRetriesOnlyOriginalHostBeforeFallback(t *testing.T) {
	transport := &hostSwitchRoundTripper{
		statusByHost: map[string]int{
			"push2his.eastmoney.com":   http.StatusServiceUnavailable,
			"7.push2his.eastmoney.com": http.StatusServiceUnavailable,
		},
	}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: transport},
		Retry:      RetryConfig{MaxRetries: 1, BaseDelay: time.Nanosecond},
	})

	_, err := client.GetText(context.Background(), "https://push2his.eastmoney.com/api/qt/stock/kline/get")
	if err == nil {
		t.Fatal("expected error")
	}
	counts := map[string]int{}
	for _, host := range transport.hosts {
		counts[host]++
	}
	if counts["push2his.eastmoney.com"] != 2 {
		t.Fatalf("original host attempts = %d, want 2; hosts=%#v", counts["push2his.eastmoney.com"], transport.hosts)
	}
	if counts["7.push2his.eastmoney.com"] != 1 {
		t.Fatalf("fallback host attempts = %d, want 1; hosts=%#v", counts["7.push2his.eastmoney.com"], transport.hosts)
	}
}

func TestClientProxyPoolRotatesOnRetries(t *testing.T) {
	proxyCalls := []string{}
	firstProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyCalls = append(proxyCalls, "first:"+r.Host)
		http.Error(w, "rate limited", http.StatusTooManyRequests)
	}))
	defer firstProxy.Close()
	secondProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyCalls = append(proxyCalls, "second:"+r.Host)
		_, _ = w.Write([]byte("ok"))
	}))
	defer secondProxy.Close()

	client := NewClient(Config{
		ProxyPool: ProxyPoolConfig{URLs: []string{
			firstProxy.URL,
			secondProxy.URL,
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
	want := []string{"first:target.test", "second:target.test"}
	if strings.Join(proxyCalls, "|") != strings.Join(want, "|") {
		t.Fatalf("proxy calls = %#v, want %#v", proxyCalls, want)
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

func TestClientGetTextUsesProviderSpecificRetryPolicy(t *testing.T) {
	transport := &providerPolicyRoundTripper{
		responses: map[string][]roundTripResult{
			"qt.gtimg.cn": {
				{status: http.StatusServiceUnavailable},
			},
			"push2his.eastmoney.com": {
				{status: http.StatusServiceUnavailable},
				{body: "ok", status: http.StatusOK},
			},
		},
	}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: transport},
		Retry:      RetryConfig{MaxRetries: 0, BaseDelay: time.Nanosecond},
		ProviderPolicies: map[ProviderName]ProviderPolicy{
			ProviderEastmoney: {
				Retry: &RetryConfig{MaxRetries: 1, BaseDelay: time.Nanosecond},
			},
		},
		HostFallback: NewHostFallbackManager(HostFallbackOptions{}),
	})

	_, tencentErr := client.GetText(context.Background(), "https://qt.gtimg.cn/?q=s_sh600519")
	if tencentErr == nil {
		t.Fatal("expected Tencent request to fail without retry")
	}
	text, err := client.GetText(context.Background(), "https://push2his.eastmoney.com/api/qt/stock/kline/get")
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
	if got := transport.calls["qt.gtimg.cn"]; got != 1 {
		t.Fatalf("Tencent attempts = %d, want 1", got)
	}
	if got := transport.calls["push2his.eastmoney.com"]; got != 2 {
		t.Fatalf("Eastmoney attempts = %d, want 2", got)
	}
}

func TestClientGetTextUsesCustomRetryableStatusCodes(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			http.Error(w, "teapot", http.StatusTeapot)
			return
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient: server.Client(),
		Retry: RetryConfig{
			MaxRetries:           1,
			BaseDelay:            time.Nanosecond,
			RetryableStatusCodes: []int{http.StatusTeapot},
		},
	})
	text, err := client.GetText(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" || attempts != 2 {
		t.Fatalf("text=%q attempts=%d, want ok after retry", text, attempts)
	}
}

func TestClientGetTextCustomRetryableStatusCodesReplaceDefaults(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient: server.Client(),
		Retry: RetryConfig{
			MaxRetries:           1,
			BaseDelay:            time.Nanosecond,
			RetryableStatusCodes: []int{http.StatusTeapot},
		},
	})
	if _, err := client.GetText(context.Background(), server.URL); err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("attempts=%d, want no retry for 503", attempts)
	}
}

func TestRetryDelayUsesBackoffMultiplierAndMaxDelay(t *testing.T) {
	got := retryDelay(3, RetryConfig{
		BaseDelay:         2 * time.Second,
		MaxDelay:          10 * time.Second,
		BackoffMultiplier: 3,
	})
	if got != 10*time.Second {
		t.Fatalf("retryDelay = %s, want 10s", got)
	}
}

func TestClientGetTextAppliesHeadersAndProviderOverrides(t *testing.T) {
	var gotGlobalHeader string
	var gotProviderHeader string
	var gotUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotGlobalHeader = r.Header.Get("X-Global")
		gotProviderHeader = r.Header.Get("X-Provider")
		gotUserAgent = r.Header.Get("User-Agent")
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(Config{
		HTTPClient: server.Client(),
		UserAgent:  "ceheng-core-test",
		Headers: map[string]string{
			"X-Global":   "global",
			"User-Agent": "header-agent",
		},
		ProviderPolicies: map[ProviderName]ProviderPolicy{
			ProviderUnknown: {
				Headers: map[string]string{"X-Provider": "provider"},
			},
		},
	})

	_, err := client.GetText(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if gotGlobalHeader != "global" {
		t.Fatalf("X-Global = %q, want global", gotGlobalHeader)
	}
	if gotProviderHeader != "provider" {
		t.Fatalf("X-Provider = %q, want provider", gotProviderHeader)
	}
	if gotUserAgent != "ceheng-core-test" {
		t.Fatalf("User-Agent = %q, want ceheng-core-test", gotUserAgent)
	}
}

func TestClientGetTextEmitsRequestHooks(t *testing.T) {
	var events []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()
	client := NewClient(Config{
		HTTPClient: server.Client(),
		Hooks: RequestHooks{
			OnRequest: func(ctx RequestContext) {
				events = append(events, "request:"+string(ctx.Provider)+":"+ctx.URL)
			},
			OnResponse: func(ctx RequestContext, meta ResponseMeta) {
				events = append(events, "response:"+http.StatusText(meta.StatusCode))
			},
			Trace: func(event RequestTraceEvent, ctx RequestContext) {
				events = append(events, "trace:"+string(event))
			},
		},
	})

	text, err := client.GetText(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
	want := []string{"request:unknown:" + server.URL, "trace:request", "response:OK", "trace:response"}
	if strings.Join(events, "|") != strings.Join(want, "|") {
		t.Fatalf("events = %#v, want %#v", events, want)
	}
}

func TestClientGetTextEmitsErrorRetryAndFallbackHooks(t *testing.T) {
	var events []string
	transport := &hostSwitchRoundTripper{
		statusByHost: map[string]int{
			"push2his.eastmoney.com": http.StatusServiceUnavailable,
		},
		bodyByHost: map[string]string{
			"7.push2his.eastmoney.com": "ok",
		},
	}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: transport},
		Retry:      RetryConfig{MaxRetries: 1, BaseDelay: time.Nanosecond},
		Hooks: RequestHooks{
			OnError: func(ctx RequestContext, err error) {
				events = append(events, "error:"+ctx.URL)
			},
			OnRetry: func(ctx RequestContext, err error, delay time.Duration) {
				events = append(events, "retry:"+ctx.URL)
			},
			Trace: func(event RequestTraceEvent, ctx RequestContext) {
				if event == TraceFallback {
					events = append(events, "fallback:"+ctx.URL)
				}
			},
		},
	})

	text, err := client.GetText(context.Background(), "https://push2his.eastmoney.com/api/qt/stock/kline/get")
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
	if !containsString(events, "retry:https://push2his.eastmoney.com/api/qt/stock/kline/get") {
		t.Fatalf("events = %#v, want retry event", events)
	}
	if !containsString(events, "fallback:https://push2his.eastmoney.com/api/qt/stock/kline/get") {
		t.Fatalf("events = %#v, want fallback event", events)
	}
}

func TestClientGetTextCallsRetryOptionCallback(t *testing.T) {
	retryAttempts := []int{}
	retryDelays := []time.Duration{}
	transport := &flakyRoundTripper{
		responses: []roundTripResult{
			{err: errors.New("temporary network failure")},
			{body: "ok", status: http.StatusOK},
		},
	}
	client := NewClient(Config{
		HTTPClient: &http.Client{Transport: transport},
		Retry: RetryConfig{
			MaxRetries: 1,
			BaseDelay:  time.Nanosecond,
			OnRetry: func(attempt int, err error, delay time.Duration) {
				retryAttempts = append(retryAttempts, attempt)
				retryDelays = append(retryDelays, delay)
			},
		},
	})

	text, err := client.GetText(context.Background(), "https://retry-callback.test/path")
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
	if len(retryAttempts) != 1 || retryAttempts[0] != 1 {
		t.Fatalf("retry attempts = %#v, want [1]", retryAttempts)
	}
	if len(retryDelays) != 1 || retryDelays[0] != time.Nanosecond {
		t.Fatalf("retry delays = %#v, want [1ns]", retryDelays)
	}
}

func TestClientGetTextIgnoresHookPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()
	client := NewClient(Config{
		HTTPClient: server.Client(),
		Hooks: RequestHooks{
			OnRequest: func(RequestContext) {
				panic("hook failed")
			},
		},
	})

	text, err := client.GetText(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("text = %q, want ok", text)
	}
}

type roundTripResult struct {
	status int
	body   string
	err    error
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type drainTrackingBody struct {
	content        []byte
	offset         int
	closed         bool
	sawEOF         bool
	closedAfterEOF bool
}

func (b *drainTrackingBody) Read(p []byte) (int, error) {
	if b.offset >= len(b.content) {
		b.sawEOF = true
		return 0, io.EOF
	}
	n := copy(p, b.content[b.offset:])
	b.offset += n
	return n, nil
}

func (b *drainTrackingBody) Close() error {
	b.closed = true
	b.closedAfterEOF = b.sawEOF
	return nil
}

type flakyRoundTripper struct {
	calls     int
	responses []roundTripResult
}

func (f *flakyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	f.calls++
	index := f.calls - 1
	if index >= len(f.responses) {
		return nil, errors.New("unexpected request")
	}
	result := f.responses[index]
	if result.err != nil {
		return nil, result.err
	}
	return &http.Response{
		StatusCode: result.status,
		Status:     http.StatusText(result.status),
		Body:       io.NopCloser(strings.NewReader(result.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

type countingLimiter struct {
	calls int
}

func (l *countingLimiter) Acquire(context.Context) error {
	l.calls++
	return nil
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

type hostSwitchRoundTripper struct {
	hosts        []string
	statusByHost map[string]int
	bodyByHost   map[string]string
}

func (t *hostSwitchRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	host := req.URL.Hostname()
	t.hosts = append(t.hosts, host)
	if status := t.statusByHost[host]; status != 0 {
		return &http.Response{
			StatusCode: status,
			Status:     http.StatusText(status),
			Body:       io.NopCloser(strings.NewReader(http.StatusText(status))),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	}
	body, ok := t.bodyByHost[host]
	if !ok {
		return nil, errors.New("unexpected host " + host)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

type providerPolicyRoundTripper struct {
	calls     map[string]int
	responses map[string][]roundTripResult
}

func (t *providerPolicyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.calls == nil {
		t.calls = map[string]int{}
	}
	host := req.URL.Hostname()
	t.calls[host]++
	results := t.responses[host]
	index := t.calls[host] - 1
	if index >= len(results) {
		return nil, errors.New("unexpected request to " + host)
	}
	result := results[index]
	if result.err != nil {
		return nil, result.err
	}
	return &http.Response{
		StatusCode: result.status,
		Status:     http.StatusText(result.status),
		Body:       io.NopCloser(strings.NewReader(result.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}
