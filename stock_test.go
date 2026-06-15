package stock

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ceheng.io/stock-go/internal/core"
)

func TestNewAppliesOptions(t *testing.T) {
	httpClient := &http.Client{Timeout: 2 * time.Second}

	client := New(
		WithHTTPClient(httpClient),
		WithTimeout(3*time.Second),
		WithUserAgent("ceheng-test-agent"),
		WithHeaders(map[string]string{"X-Global": "global"}),
		WithRetry(RetryOptions{MaxRetries: 2, BaseDelay: 10 * time.Millisecond}),
		WithRateLimit(RateLimitOptions{RequestsPerSecond: 3, MaxBurst: 2}),
		WithCircuitBreaker(CircuitBreakerOptions{
			FailureThreshold: 2,
			ResetTimeout:     time.Minute,
			HalfOpenRequests: 1,
		}),
		WithProviderPolicy(ProviderEastmoney, ProviderPolicy{
			Timeout:   5 * time.Second,
			UserAgent: "ceheng-eastmoney-test",
			Retry:     &RetryOptions{MaxRetries: 4, BaseDelay: time.Millisecond},
			Headers:   map[string]string{"X-Provider": "provider"},
		}),
		WithRequestHooks(RequestHooks{
			OnRequest: func(RequestContext) {},
		}),
	)

	if client == nil {
		t.Fatal("New returned nil")
	}
	if client.config.HTTPClient != httpClient {
		t.Fatal("WithHTTPClient did not set the HTTP client")
	}
	if client.config.Timeout != 3*time.Second {
		t.Fatalf("Timeout = %s, want 3s", client.config.Timeout)
	}
	if client.config.UserAgent != "ceheng-test-agent" {
		t.Fatalf("UserAgent = %q", client.config.UserAgent)
	}
	if client.config.Headers["X-Global"] != "global" {
		t.Fatalf("global header = %q, want global", client.config.Headers["X-Global"])
	}
	if client.config.Retry.MaxRetries != 2 {
		t.Fatalf("MaxRetries = %d, want 2", client.config.Retry.MaxRetries)
	}
	if client.config.RateLimit.RequestsPerSecond != 3 {
		t.Fatalf("RequestsPerSecond = %.0f, want 3", client.config.RateLimit.RequestsPerSecond)
	}
	if client.config.RateLimit.MaxBurst != 2 {
		t.Fatalf("MaxBurst = %.0f, want 2", client.config.RateLimit.MaxBurst)
	}
	if client.config.CircuitBreaker.FailureThreshold != 2 {
		t.Fatalf("FailureThreshold = %d, want 2", client.config.CircuitBreaker.FailureThreshold)
	}
	if client.config.CircuitBreaker.ResetTimeout != time.Minute {
		t.Fatalf("ResetTimeout = %s, want 1m", client.config.CircuitBreaker.ResetTimeout)
	}
	if client.config.CircuitBreaker.HalfOpenRequests != 1 {
		t.Fatalf("HalfOpenRequests = %d, want 1", client.config.CircuitBreaker.HalfOpenRequests)
	}
	policy := client.config.ProviderPolicies[ProviderEastmoney]
	if policy.Timeout != 5*time.Second {
		t.Fatalf("provider Timeout = %s, want 5s", policy.Timeout)
	}
	if policy.UserAgent != "ceheng-eastmoney-test" {
		t.Fatalf("provider UserAgent = %q, want ceheng-eastmoney-test", policy.UserAgent)
	}
	if policy.Headers["X-Provider"] != "provider" {
		t.Fatalf("provider header = %q, want provider", policy.Headers["X-Provider"])
	}
	if policy.Retry == nil || policy.Retry.MaxRetries != 4 || policy.Retry.BaseDelay != time.Millisecond {
		t.Fatalf("provider Retry = %+v, want maxRetries=4 baseDelay=1ms", policy.Retry)
	}
	if client.config.RequestHooks.OnRequest == nil {
		t.Fatal("RequestHooks.OnRequest was not set")
	}
}

func TestNewUsesDefaults(t *testing.T) {
	client := New()

	if client.config.HTTPClient == nil {
		t.Fatal("default HTTPClient is nil")
	}
	if client.config.Timeout <= 0 {
		t.Fatalf("default Timeout = %s, want positive", client.config.Timeout)
	}
	if client.config.Timeout != 30*time.Second {
		t.Fatalf("default Timeout = %s, want 30s", client.config.Timeout)
	}
	if client.config.Retry.MaxRetries != 3 || client.config.Retry.BaseDelay != time.Second {
		t.Fatalf("default Retry = %+v, want maxRetries=3 baseDelay=1s", client.config.Retry)
	}
	if client.config.UserAgent == "" {
		t.Fatal("default UserAgent is empty")
	}
	if client.config.AShareListURL != AShareListURL || client.config.USListURL != USListURL || client.config.HKListURL != HKListURL || client.config.FundListURL != FundListURL {
		t.Fatalf("default code list URLs = (%q, %q, %q, %q), want TS constants", client.config.AShareListURL, client.config.USListURL, client.config.HKListURL, client.config.FundListURL)
	}
}

func TestRootReExportsTSEntryTypeNames(t *testing.T) {
	var sdk *StockSDK = New()
	if sdk == nil {
		t.Fatal("StockSDK alias returned nil")
	}

	var options RequestClientOptions = Config{
		Timeout:   time.Second,
		UserAgent: "ceheng-test-agent",
		Retry:     RetryOptions{MaxRetries: 1},
	}
	var policy ProviderRequestPolicy = ProviderPolicy{
		UserAgent: "ceheng-provider-agent",
		Retry:     &RetryOptions{MaxRetries: 2},
	}
	if options.UserAgent == "" || policy.UserAgent == "" {
		t.Fatalf("compat aliases were not usable: %+v %+v", options, policy)
	}
}

func TestRootReExportsTopLevelServiceTypes(t *testing.T) {
	client := New()

	var quotes *QuoteService = client.Quotes
	var kline *KlineService = client.Kline
	var indicator *IndicatorService = client.Indicator
	var board *BoardService = client.Board
	var calendar *TradingCalendarService = client.Calendar
	var fundFlow *FundFlowService = client.FundFlow
	var northbound *NorthboundService = client.Northbound
	var dragonTiger *DragonTigerService = client.DragonTiger
	var blockTrade *BlockTradeService = client.BlockTrade
	var margin *MarginService = client.Margin
	var dividend *DividendService = client.Dividend
	var data *DataService = client.Data
	var marketEvent *MarketEventService = client.MarketEvent
	var fund *FundService = client.Fund
	var futures *FuturesService = client.Futures
	var options *OptionsService = client.Options

	services := []any{
		quotes, kline, indicator, board, calendar, fundFlow, northbound,
		dragonTiger, blockTrade, margin, dividend, data, marketEvent,
		fund, futures, options,
	}
	for index, service := range services {
		if service == nil {
			t.Fatalf("service alias %d is nil", index)
		}
	}
}

func TestNewWiresRateLimiter(t *testing.T) {
	client := New(WithRateLimit(RateLimitOptions{RequestsPerSecond: 10, MaxBurst: 1}))

	if client.core == nil {
		t.Fatal("core client is nil")
	}
	if client.core.RateLimiter() == nil {
		t.Fatal("core rate limiter is nil")
	}
}

func TestNewWiresCircuitBreaker(t *testing.T) {
	client := New(WithCircuitBreaker(CircuitBreakerOptions{
		FailureThreshold: 1,
		ResetTimeout:     time.Second,
		HalfOpenRequests: 1,
	}))

	if client.core == nil {
		t.Fatal("core client is nil")
	}
	if client.core.CircuitBreaker() == nil {
		t.Fatal("core circuit breaker is nil")
	}
}

func TestNewWiresHostFallback(t *testing.T) {
	client := New()

	if client.core == nil {
		t.Fatal("core client is nil")
	}
	if client.core.HostFallback() == nil {
		t.Fatal("core host fallback manager is nil")
	}
}

func TestNewWiresProviderPolicy(t *testing.T) {
	noTimeoutRetry := false
	client := New(WithProviderPolicy(ProviderEastmoney, ProviderPolicy{
		Retry: &RetryOptions{
			MaxRetries:           3,
			BaseDelay:            time.Nanosecond,
			MaxDelay:             time.Millisecond,
			BackoffMultiplier:    4,
			RetryableStatusCodes: []int{418},
			RetryOnTimeout:       &noTimeoutRetry,
		},
	}))

	if client.core == nil {
		t.Fatal("core client is nil")
	}
	policy, ok := client.core.ProviderPolicy(core.ProviderEastmoney)
	if !ok {
		t.Fatal("core provider policy is missing")
	}
	if policy.Retry.MaxRetries != 3 {
		t.Fatalf("core provider MaxRetries = %d, want 3", policy.Retry.MaxRetries)
	}
	if policy.Retry.MaxDelay != time.Millisecond || policy.Retry.BackoffMultiplier != 4 {
		t.Fatalf("core provider retry timing = %+v", policy.Retry)
	}
	if len(policy.Retry.RetryableStatusCodes) != 1 || policy.Retry.RetryableStatusCodes[0] != 418 {
		t.Fatalf("core provider retry statuses = %#v", policy.Retry.RetryableStatusCodes)
	}
	if policy.Retry.RetryOnTimeout == nil || *policy.Retry.RetryOnTimeout {
		t.Fatalf("core provider RetryOnTimeout = %v, want false", policy.Retry.RetryOnTimeout)
	}
}

func TestRetryOptionsCanDisableNetworkAndTimeoutRetries(t *testing.T) {
	noNetworkRetry := false
	noTimeoutRetry := false
	attempts := 0
	client := New(
		WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				attempts++
				return nil, errors.New("network down")
			}),
		}),
		WithRetry(RetryOptions{
			MaxRetries:          1,
			BaseDelay:           time.Nanosecond,
			RetryOnNetworkError: &noNetworkRetry,
			RetryOnTimeout:      &noTimeoutRetry,
		}),
	)

	_, err := client.core.GetText(context.Background(), "https://retry-disabled.test/path")
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestRetryOptionsOnRetryCallbackIsWired(t *testing.T) {
	retryAttempts := []int{}
	transportCalls := 0
	client := New(
		WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				transportCalls++
				if transportCalls == 1 {
					return nil, errors.New("network down")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body:       http.NoBody,
					Header:     make(http.Header),
					Request:    req,
				}, nil
			}),
		}),
		WithRetry(RetryOptions{
			MaxRetries: 1,
			BaseDelay:  time.Nanosecond,
			OnRetry: func(attempt int, err error, delay time.Duration) {
				retryAttempts = append(retryAttempts, attempt)
			},
		}),
	)

	text, err := client.core.GetText(context.Background(), "https://retry-callback.test/path")
	if err != nil {
		t.Fatal(err)
	}
	if text != "" {
		t.Fatalf("text = %q, want empty response body", text)
	}
	if len(retryAttempts) != 1 || retryAttempts[0] != 1 {
		t.Fatalf("retry attempts = %#v, want [1]", retryAttempts)
	}
}

func TestNewWiresRequestHooks(t *testing.T) {
	client := New(WithRequestHooks(RequestHooks{
		OnRequest: func(RequestContext) {},
	}))

	if client.core == nil {
		t.Fatal("core client is nil")
	}
	if client.core.Hooks().OnRequest == nil {
		t.Fatal("core request hooks are missing")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
