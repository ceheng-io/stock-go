package stock

import (
	"net/url"
	"strings"
	"time"

	"github.com/ceheng.io/stock-go/constants"
)

// ResolvedRetryOptions contains complete retry settings.
type ResolvedRetryOptions struct {
	MaxRetries           int
	BaseDelay            time.Duration
	MaxDelay             time.Duration
	BackoffMultiplier    float64
	RetryableStatusCodes []int
	RetryOnNetworkError  bool
	RetryOnTimeout       bool
	OnRetry              func(attempt int, err error, delay time.Duration)
}

// ResolvedProviderPolicy contains complete provider request governance settings.
type ResolvedProviderPolicy struct {
	Timeout         time.Duration
	Retry           ResolvedRetryOptions
	Headers         map[string]string
	UserAgent       string
	RotateUserAgent bool
	RateLimit       *RateLimitOptions
	CircuitBreaker  *CircuitBreakerOptions
}

// MergeProviderPolicy merges provider policies with shallow nested option merging.
func MergeProviderPolicy(base ProviderPolicy, override *ProviderPolicy) ProviderPolicy {
	merged := cloneProviderPolicy(base)
	if override == nil {
		return merged
	}
	if override.Timeout > 0 {
		merged.Timeout = override.Timeout
	}
	if override.UserAgent != "" {
		merged.UserAgent = override.UserAgent
	}
	if override.RotateUserAgent != nil {
		merged.RotateUserAgent = boolPtr(*override.RotateUserAgent)
	}
	if len(override.Headers) > 0 {
		merged.Headers = mergeStringMap(merged.Headers, override.Headers)
	}
	if override.Retry != nil {
		retry := RetryOptions{}
		if merged.Retry != nil {
			retry = *merged.Retry
		}
		if override.Retry.MaxRetries >= 0 {
			retry.MaxRetries = override.Retry.MaxRetries
		}
		if override.Retry.BaseDelay > 0 {
			retry.BaseDelay = override.Retry.BaseDelay
		}
		if override.Retry.MaxDelay > 0 {
			retry.MaxDelay = override.Retry.MaxDelay
		}
		if override.Retry.BackoffMultiplier > 0 {
			retry.BackoffMultiplier = override.Retry.BackoffMultiplier
		}
		if override.Retry.RetryableStatusCodes != nil {
			retry.RetryableStatusCodes = append([]int(nil), override.Retry.RetryableStatusCodes...)
		}
		if override.Retry.RetryOnNetworkError != nil {
			retry.RetryOnNetworkError = boolPtr(*override.Retry.RetryOnNetworkError)
		}
		if override.Retry.RetryOnTimeout != nil {
			retry.RetryOnTimeout = boolPtr(*override.Retry.RetryOnTimeout)
		}
		if override.Retry.OnRetry != nil {
			retry.OnRetry = override.Retry.OnRetry
		}
		merged.Retry = &retry
	}
	if override.RateLimit != nil {
		rateLimit := RateLimitOptions{}
		if merged.RateLimit != nil {
			rateLimit = *merged.RateLimit
		}
		if override.RateLimit.RequestsPerSecond > 0 {
			rateLimit.RequestsPerSecond = override.RateLimit.RequestsPerSecond
		}
		if override.RateLimit.MaxBurst > 0 {
			rateLimit.MaxBurst = override.RateLimit.MaxBurst
		}
		merged.RateLimit = &rateLimit
	}
	if override.CircuitBreaker != nil {
		circuitBreaker := CircuitBreakerOptions{}
		if merged.CircuitBreaker != nil {
			circuitBreaker = *merged.CircuitBreaker
		}
		if override.CircuitBreaker.FailureThreshold > 0 {
			circuitBreaker.FailureThreshold = override.CircuitBreaker.FailureThreshold
		}
		if override.CircuitBreaker.ResetTimeout > 0 {
			circuitBreaker.ResetTimeout = override.CircuitBreaker.ResetTimeout
		}
		if override.CircuitBreaker.HalfOpenRequests > 0 {
			circuitBreaker.HalfOpenRequests = override.CircuitBreaker.HalfOpenRequests
		}
		if override.CircuitBreaker.OnStateChange != nil {
			circuitBreaker.OnStateChange = override.CircuitBreaker.OnStateChange
		}
		merged.CircuitBreaker = &circuitBreaker
	}
	return merged
}

// ResolveProviderPolicy fills provider policy defaults and normalizes headers.
func ResolveProviderPolicy(policy ProviderPolicy) ResolvedProviderPolicy {
	timeout := policy.Timeout
	if timeout <= 0 {
		timeout = time.Duration(constants.DefaultTimeoutMS) * time.Millisecond
	}
	headers := normalizeHeaders(policy.Headers, policy.UserAgent)
	return ResolvedProviderPolicy{
		Timeout:         timeout,
		Retry:           resolveRetryOptions(policy.Retry),
		Headers:         headers,
		UserAgent:       policy.UserAgent,
		RotateUserAgent: boolValue(policy.RotateUserAgent, false),
		RateLimit:       cloneRateLimitOptions(policy.RateLimit),
		CircuitBreaker:  cloneCircuitBreakerOptions(policy.CircuitBreaker),
	}
}

// InferProviderFromURL infers a provider from URL host unless explicitProvider is set.
func InferProviderFromURL(rawURL string, explicitProvider ProviderName) ProviderName {
	if explicitProvider != "" {
		return explicitProvider
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ProviderUnknown
	}
	host := parsed.Hostname()
	switch {
	case strings.Contains(host, "eastmoney.com"):
		return ProviderEastmoney
	case strings.Contains(host, "gtimg.cn"):
		return ProviderTencent
	case strings.Contains(host, "sina.com.cn"):
		return ProviderSina
	case strings.Contains(host, "linkdiary.cn"):
		return ProviderLinkdiary
	default:
		return ProviderUnknown
	}
}

// InferProviderFromUrl infers a provider using the TypeScript SDK naming style.
func InferProviderFromUrl(rawURL string, explicitProvider ProviderName) ProviderName {
	return InferProviderFromURL(rawURL, explicitProvider)
}

func resolveRetryOptions(options *RetryOptions) ResolvedRetryOptions {
	resolved := ResolvedRetryOptions{
		MaxRetries:           constants.DefaultMaxRetries,
		BaseDelay:            time.Duration(constants.DefaultBaseDelayMS) * time.Millisecond,
		MaxDelay:             time.Duration(constants.DefaultMaxDelayMS) * time.Millisecond,
		BackoffMultiplier:    constants.DefaultBackoffMultiplier,
		RetryableStatusCodes: constants.DefaultRetryableStatusCodes(),
		RetryOnNetworkError:  true,
		RetryOnTimeout:       true,
	}
	if options == nil {
		return resolved
	}
	if options.MaxRetries >= 0 {
		resolved.MaxRetries = options.MaxRetries
	}
	if options.BaseDelay > 0 {
		resolved.BaseDelay = options.BaseDelay
	}
	if options.MaxDelay > 0 {
		resolved.MaxDelay = options.MaxDelay
	}
	if options.BackoffMultiplier > 0 {
		resolved.BackoffMultiplier = options.BackoffMultiplier
	}
	if options.RetryableStatusCodes != nil {
		resolved.RetryableStatusCodes = append([]int(nil), options.RetryableStatusCodes...)
	}
	if options.RetryOnNetworkError != nil {
		resolved.RetryOnNetworkError = *options.RetryOnNetworkError
	}
	if options.RetryOnTimeout != nil {
		resolved.RetryOnTimeout = *options.RetryOnTimeout
	}
	if options.OnRetry != nil {
		resolved.OnRetry = options.OnRetry
	}
	return resolved
}

func normalizeHeaders(headers map[string]string, userAgent string) map[string]string {
	normalized := cloneStringMap(headers)
	if normalized == nil {
		normalized = map[string]string{}
	}
	if userAgent == "" {
		return normalized
	}
	for key := range normalized {
		if strings.EqualFold(key, "User-Agent") {
			return normalized
		}
	}
	normalized["User-Agent"] = userAgent
	return normalized
}

func cloneProviderPolicy(policy ProviderPolicy) ProviderPolicy {
	clone := policy
	clone.Headers = cloneStringMap(policy.Headers)
	if policy.Retry != nil {
		retry := *policy.Retry
		clone.Retry = &retry
	}
	clone.RotateUserAgent = cloneBoolPtr(policy.RotateUserAgent)
	clone.RateLimit = cloneRateLimitOptions(policy.RateLimit)
	clone.CircuitBreaker = cloneCircuitBreakerOptions(policy.CircuitBreaker)
	return clone
}

func cloneRateLimitOptions(options *RateLimitOptions) *RateLimitOptions {
	if options == nil {
		return nil
	}
	clone := *options
	return &clone
}

func cloneCircuitBreakerOptions(options *CircuitBreakerOptions) *CircuitBreakerOptions {
	if options == nil {
		return nil
	}
	clone := *options
	return &clone
}
