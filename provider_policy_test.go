package stock

import (
	"testing"
	"time"
)

func TestMergeProviderPolicyMergesNestedOptions(t *testing.T) {
	baseRotate := true
	overrideRotate := false
	base := ProviderPolicy{
		Timeout:         10 * time.Second,
		UserAgent:       "base-agent",
		RotateUserAgent: &baseRotate,
		Headers:         map[string]string{"Accept": "json", "X-Base": "1"},
		Retry:           &RetryOptions{MaxRetries: 1, BaseDelay: 100 * time.Millisecond},
		RateLimit:       &RateLimitOptions{RequestsPerSecond: 1, MaxBurst: 2},
		CircuitBreaker: &CircuitBreakerOptions{
			FailureThreshold: 2,
			ResetTimeout:     time.Second,
			HalfOpenRequests: 1,
		},
	}
	override := ProviderPolicy{
		RotateUserAgent: &overrideRotate,
		Headers:         map[string]string{"X-Override": "2", "Accept": "text"},
		Retry:           &RetryOptions{MaxRetries: 3},
		RateLimit: &RateLimitOptions{
			MaxBurst: 5,
		},
		CircuitBreaker: &CircuitBreakerOptions{
			HalfOpenRequests: 4,
		},
	}

	merged := MergeProviderPolicy(base, &override)
	if merged.Timeout != 10*time.Second || merged.UserAgent != "base-agent" {
		t.Fatalf("top-level merge = %+v", merged)
	}
	if merged.RotateUserAgent == nil || *merged.RotateUserAgent {
		t.Fatalf("RotateUserAgent = %v, want false", merged.RotateUserAgent)
	}
	if merged.Headers["Accept"] != "text" || merged.Headers["X-Base"] != "1" || merged.Headers["X-Override"] != "2" {
		t.Fatalf("headers = %#v", merged.Headers)
	}
	if merged.Retry == nil || merged.Retry.MaxRetries != 3 || merged.Retry.BaseDelay != 100*time.Millisecond {
		t.Fatalf("retry merge = %+v", merged.Retry)
	}
	if merged.RateLimit == nil || merged.RateLimit.RequestsPerSecond != 1 || merged.RateLimit.MaxBurst != 5 {
		t.Fatalf("rate limit merge = %+v", merged.RateLimit)
	}
	if merged.CircuitBreaker == nil || merged.CircuitBreaker.FailureThreshold != 2 || merged.CircuitBreaker.HalfOpenRequests != 4 {
		t.Fatalf("circuit breaker merge = %+v", merged.CircuitBreaker)
	}

	merged.Headers["Accept"] = "mutated"
	if base.Headers["Accept"] != "json" || override.Headers["Accept"] != "text" {
		t.Fatalf("MergeProviderPolicy reused input header maps")
	}
}

func TestResolveProviderPolicyFillsDefaultsAndNormalizesHeaders(t *testing.T) {
	rotateUserAgent := true
	resolved := ResolveProviderPolicy(ProviderPolicy{
		UserAgent:       "ua",
		RotateUserAgent: &rotateUserAgent,
		Headers:         map[string]string{"Accept": "json"},
		Retry:           &RetryOptions{MaxRetries: 4},
	})

	if resolved.Timeout <= 0 {
		t.Fatalf("Timeout = %s, want positive", resolved.Timeout)
	}
	if resolved.Headers["User-Agent"] != "ua" || resolved.Headers["Accept"] != "json" {
		t.Fatalf("headers = %#v", resolved.Headers)
	}
	if resolved.Retry.MaxRetries != 4 || resolved.Retry.BaseDelay <= 0 {
		t.Fatalf("retry = %+v", resolved.Retry)
	}
	if len(resolved.Retry.RetryableStatusCodes) == 0 || resolved.Retry.RetryableStatusCodes[0] != 408 {
		t.Fatalf("retryable statuses = %#v", resolved.Retry.RetryableStatusCodes)
	}
	if !resolved.RotateUserAgent {
		t.Fatal("RotateUserAgent = false, want true")
	}

	noNetworkRetry := false
	noTimeoutRetry := false
	onRetry := func(int, error, time.Duration) {}
	custom := ResolveProviderPolicy(ProviderPolicy{
		Retry: &RetryOptions{
			MaxRetries:           2,
			BaseDelay:            2 * time.Second,
			MaxDelay:             9 * time.Second,
			BackoffMultiplier:    3,
			RetryableStatusCodes: []int{418, 429},
			RetryOnNetworkError:  &noNetworkRetry,
			RetryOnTimeout:       &noTimeoutRetry,
			OnRetry:              onRetry,
		},
	})
	if custom.Retry.MaxRetries != 2 || custom.Retry.BaseDelay != 2*time.Second || custom.Retry.MaxDelay != 9*time.Second {
		t.Fatalf("custom retry timing = %+v", custom.Retry)
	}
	if custom.Retry.BackoffMultiplier != 3 || custom.Retry.RetryOnNetworkError || custom.Retry.RetryOnTimeout {
		t.Fatalf("custom retry switches = %+v", custom.Retry)
	}
	if custom.Retry.OnRetry == nil {
		t.Fatal("custom retry OnRetry is nil")
	}
	if len(custom.Retry.RetryableStatusCodes) != 2 || custom.Retry.RetryableStatusCodes[0] != 418 || custom.Retry.RetryableStatusCodes[1] != 429 {
		t.Fatalf("custom retry statuses = %#v", custom.Retry.RetryableStatusCodes)
	}

	explicitUA := ResolveProviderPolicy(ProviderPolicy{
		UserAgent: "ignored",
		Headers:   map[string]string{"user-agent": "explicit"},
	})
	if explicitUA.Headers["user-agent"] != "explicit" {
		t.Fatalf("explicit user-agent was not preserved: %#v", explicitUA.Headers)
	}
	if _, ok := explicitUA.Headers["User-Agent"]; ok {
		t.Fatalf("unexpected User-Agent header added: %#v", explicitUA.Headers)
	}
}

func TestInferProviderFromURL(t *testing.T) {
	tests := map[string]ProviderName{
		"https://qt.gtimg.cn/?q=sh600519":                                      ProviderTencent,
		"https://push2his.eastmoney.com/api/qt/stock/kline/get?secid=1.600519": ProviderEastmoney,
		"https://stock.finance.sina.com.cn/futures/api/openapi.php":            ProviderSina,
		"https://assets.linkdiary.cn/shares/zh_a_list.json":                    ProviderLinkdiary,
		"://bad-url": ProviderUnknown,
	}

	for input, want := range tests {
		if got := InferProviderFromURL(input, ""); got != want {
			t.Fatalf("InferProviderFromURL(%q) = %s, want %s", input, got, want)
		}
	}
	if got := InferProviderFromURL("https://example.com", ProviderSina); got != ProviderSina {
		t.Fatalf("explicit provider = %s, want %s", got, ProviderSina)
	}
}

func TestInferProviderFromUrlTSNaming(t *testing.T) {
	if got := InferProviderFromUrl("https://qt.gtimg.cn/?q=sh600519", ""); got != ProviderTencent {
		t.Fatalf("InferProviderFromUrl = %s, want %s", got, ProviderTencent)
	}
}
