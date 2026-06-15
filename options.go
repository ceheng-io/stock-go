package stock

import (
	"net/http"
	"strings"
	"time"

	"github.com/ceheng.io/stock-go/constants"
	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/useragent"
)

// Option configures a Client.
type Option func(*Config)

// Config contains root SDK configuration.
type Config struct {
	HTTPClient                       *http.Client
	ProxyPool                        ProxyPoolOptions
	BaseURL                          string
	TencentMinuteURL                 string
	SearchBaseURL                    string
	CalendarURL                      string
	AShareListURL                    string
	USListURL                        string
	HKListURL                        string
	FundListURL                      string
	EastmoneyKlineURL                string
	EastmoneyTrendsURL               string
	EastmoneyHKKlineURL              string
	EastmoneyHKTrendsURL             string
	EastmoneyUSKlineURL              string
	EastmoneyUSTrendsURL             string
	EastmoneyIndustryListURL         string
	EastmoneyIndustrySpotURL         string
	EastmoneyIndustryConstituentsURL string
	EastmoneyIndustryKlineURL        string
	EastmoneyIndustryTrendsURL       string
	EastmoneyConceptListURL          string
	EastmoneyConceptSpotURL          string
	EastmoneyConceptConstituentsURL  string
	EastmoneyConceptKlineURL         string
	EastmoneyConceptTrendsURL        string
	EastmoneyFundFlowURL             string
	EastmoneyClistURL                string
	EastmoneyNorthboundMinuteURL     string
	EastmoneyDatacenterURL           string
	EastmoneyTopicURL                string
	EastmoneyFundGZURL               string
	EastmoneyFundPingzhongURL        string
	EastmoneyFundDataIndexURL        string
	EastmoneyFuturesKlineURL         string
	EastmoneyFuturesGlobalSpotURL    string
	EastmoneyFuturesGlobalKlineURL   string
	EastmoneyOptionCFFEXURL          string
	EastmoneyOptionLHBURL            string
	SinaETFOptionListURL             string
	SinaETFOptionExpireURL           string
	SinaETFOptionMinuteURL           string
	SinaETFOptionDailyURL            string
	SinaETFOption5DayURL             string
	SinaIndexOptionSpotURL           string
	SinaIndexOptionKlineURL          string
	SinaCommodityOptionSpotURL       string
	SinaCommodityOptionKlineURL      string
	THSLimitUpPoolURL                string
	Timeout                          time.Duration
	UserAgent                        string
	RotateUserAgent                  bool
	Headers                          map[string]string
	Retry                            RetryOptions
	RateLimit                        RateLimitOptions
	CircuitBreaker                   CircuitBreakerOptions
	ProviderPolicies                 map[ProviderName]ProviderPolicy
	RequestHooks                     RequestHooks
}

// RequestClientOptions preserves the TypeScript SDK configuration type name.
type RequestClientOptions = Config

// RetryOptions controls request retry behavior.
type RetryOptions struct {
	MaxRetries           int
	BaseDelay            time.Duration
	MaxDelay             time.Duration
	BackoffMultiplier    float64
	RetryableStatusCodes []int
	RetryOnNetworkError  *bool
	RetryOnTimeout       *bool
	OnRetry              func(attempt int, err error, delay time.Duration)
}

// RateLimitOptions controls token-bucket request throttling.
type RateLimitOptions struct {
	RequestsPerSecond float64
	MaxBurst          float64
}

// ProxyPoolOptions controls round-robin proxy selection for provider requests.
type ProxyPoolOptions struct {
	URLs []string
}

// ProviderName identifies an upstream data source for provider-level policies.
type ProviderName = core.ProviderName

const (
	ProviderTencent   ProviderName = core.ProviderTencent
	ProviderEastmoney ProviderName = core.ProviderEastmoney
	ProviderSina      ProviderName = core.ProviderSina
	ProviderLinkdiary ProviderName = core.ProviderLinkdiary
	ProviderTHS       ProviderName = core.ProviderTHS
	ProviderUnknown   ProviderName = core.ProviderUnknown
)

// ProviderPolicy overrides request governance for one provider.
type ProviderPolicy struct {
	Timeout         time.Duration
	UserAgent       string
	RotateUserAgent *bool
	Headers         map[string]string
	Retry           *RetryOptions
	RateLimit       *RateLimitOptions
	CircuitBreaker  *CircuitBreakerOptions
}

// ProviderRequestPolicy preserves the TypeScript SDK provider policy type name.
type ProviderRequestPolicy = ProviderPolicy

// RequestTraceEvent describes a request lifecycle trace event.
type RequestTraceEvent = core.RequestTraceEvent

const (
	TraceRequest  RequestTraceEvent = core.TraceRequest
	TraceResponse RequestTraceEvent = core.TraceResponse
	TraceError    RequestTraceEvent = core.TraceError
	TraceRetry    RequestTraceEvent = core.TraceRetry
	TraceFallback RequestTraceEvent = core.TraceFallback
)

// RequestContext describes one request lifecycle event.
type RequestContext = core.RequestContext

// ResponseMeta describes response timing metadata.
type ResponseMeta = core.ResponseMeta

// RequestHooks observes request lifecycle events.
type RequestHooks = core.RequestHooks

// CircuitBreakerOptions controls repeated-failure protection.
type CircuitBreakerOptions struct {
	FailureThreshold int
	ResetTimeout     time.Duration
	HalfOpenRequests int
	OnStateChange    func(from CircuitState, to CircuitState)
}

// CircuitState describes the current circuit breaker state.
type CircuitState = core.CircuitState

const (
	CircuitClosed   CircuitState = core.CircuitClosed
	CircuitOpen     CircuitState = core.CircuitOpen
	CircuitHalfOpen CircuitState = core.CircuitHalfOpen
)

func defaultConfig() Config {
	return Config{
		HTTPClient:                       http.DefaultClient,
		BaseURL:                          constants.TencentBaseURL,
		TencentMinuteURL:                 constants.TencentMinuteURL,
		SearchBaseURL:                    "https://smartbox.gtimg.cn/s3/",
		CalendarURL:                      "https://proxy.finance.qq.com/ifzqgtimg/appstock/app/newfqkline/get?param=calendar",
		AShareListURL:                    constants.AShareListURL,
		USListURL:                        constants.USListURL,
		HKListURL:                        constants.HKListURL,
		FundListURL:                      constants.FundListURL,
		EastmoneyKlineURL:                constants.EMKlineURL,
		EastmoneyTrendsURL:               constants.EMTrendsURL,
		EastmoneyHKKlineURL:              constants.EMHKKlineURL,
		EastmoneyHKTrendsURL:             constants.EMHKTrendsURL,
		EastmoneyUSKlineURL:              constants.EMUSKlineURL,
		EastmoneyUSTrendsURL:             constants.EMUSTrendsURL,
		EastmoneyIndustryListURL:         constants.EMBoardListURL,
		EastmoneyIndustrySpotURL:         constants.EMBoardSpotURL,
		EastmoneyIndustryConstituentsURL: constants.EMBoardConsURL,
		EastmoneyIndustryKlineURL:        constants.EMBoardKlineURL,
		EastmoneyIndustryTrendsURL:       constants.EMBoardTrendsURL,
		EastmoneyConceptListURL:          constants.EMConceptListURL,
		EastmoneyConceptSpotURL:          constants.EMConceptSpotURL,
		EastmoneyConceptConstituentsURL:  constants.EMConceptConsURL,
		EastmoneyConceptKlineURL:         constants.EMConceptKlineURL,
		EastmoneyConceptTrendsURL:        constants.EMConceptTrendsURL,
		EastmoneyFundFlowURL:             constants.EMFFlowURL,
		EastmoneyClistURL:                constants.EMClistURL,
		EastmoneyNorthboundMinuteURL:     constants.EMNorthboundMinuteURL,
		EastmoneyDatacenterURL:           constants.EMDatacenterURL,
		EastmoneyTopicURL:                constants.EMTopicBaseURL,
		EastmoneyFundGZURL:               "https://fundgz.1234567.com.cn/js",
		EastmoneyFundPingzhongURL:        "https://fund.eastmoney.com/pingzhongdata",
		EastmoneyFundDataIndexURL:        "https://fund.eastmoney.com/Data/funddataIndex_Interface.aspx",
		EastmoneyFuturesKlineURL:         constants.EMFuturesKlineURL,
		EastmoneyFuturesGlobalSpotURL:    constants.EMFuturesGlobalSpotURL,
		EastmoneyFuturesGlobalKlineURL:   constants.EMFuturesKlineURL,
		EastmoneyOptionCFFEXURL:          constants.EMOptionCFFEXURL,
		EastmoneyOptionLHBURL:            constants.EMOptionLHBURL,
		SinaETFOptionListURL:             constants.SinaSSEOptionListURL,
		SinaETFOptionExpireURL:           constants.SinaSSEOptionExpireURL,
		SinaETFOptionMinuteURL:           constants.SinaSSEOptionMinuteURL,
		SinaETFOptionDailyURL:            constants.SinaSSEOptionDailyURL,
		SinaETFOption5DayURL:             constants.SinaSSEOption5DayURL,
		SinaIndexOptionSpotURL:           constants.SinaOptionAPIURL,
		SinaIndexOptionKlineURL:          constants.SinaOptionDaylineURL,
		SinaCommodityOptionSpotURL:       constants.SinaOptionAPIURL,
		SinaCommodityOptionKlineURL:      constants.SinaOptionDaylineURL,
		THSLimitUpPoolURL:                constants.THSLimitUpPoolURL,
		Timeout:                          30 * time.Second,
		UserAgent:                        "ceheng-stock-go/0.1",
		Retry: RetryOptions{
			MaxRetries:           constants.DefaultMaxRetries,
			BaseDelay:            time.Duration(constants.DefaultBaseDelayMS) * time.Millisecond,
			MaxDelay:             time.Duration(constants.DefaultMaxDelayMS) * time.Millisecond,
			BackoffMultiplier:    constants.DefaultBackoffMultiplier,
			RetryableStatusCodes: constants.DefaultRetryableStatusCodes(),
			RetryOnNetworkError:  boolPtr(true),
			RetryOnTimeout:       boolPtr(true),
		},
		ProviderPolicies: map[ProviderName]ProviderPolicy{
			ProviderTHS: {
				UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
				Headers: map[string]string{
					"Accept":  "application/json, text/plain, */*",
					"Referer": "https://data.10jqka.com.cn/market/ztStock/",
					"Cookie":  "v=A0aSl97zW6psJw9OiWEn2CdlkTfNp4vvXOm-xTBvMghEJ-jpmDfacSx7DtgD",
				},
			},
		},
	}
}

// WithHTTPClient injects the HTTP client used by the SDK.
func WithHTTPClient(client *http.Client) Option {
	return func(config *Config) {
		if client != nil {
			config.HTTPClient = client
		}
	}
}

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

// WithTimeout sets the default request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(config *Config) {
		if timeout > 0 {
			config.Timeout = timeout
		}
	}
}

// WithUserAgent sets the User-Agent sent by provider requests.
func WithUserAgent(userAgent string) Option {
	return func(config *Config) {
		if userAgent != "" {
			config.UserAgent = userAgent
		}
	}
}

// WithNextUserAgent sets the next User-Agent from the shared rotation pool.
func WithNextUserAgent() Option {
	return func(config *Config) {
		if value := useragent.Next(); value != "" {
			config.UserAgent = value
		}
	}
}

// WithRandomUserAgent sets a random User-Agent from the shared pool.
func WithRandomUserAgent() Option {
	return func(config *Config) {
		if value := useragent.Random(); value != "" {
			config.UserAgent = value
		}
	}
}

// WithRotateUserAgent rotates the User-Agent per provider request.
func WithRotateUserAgent() Option {
	return func(config *Config) {
		config.RotateUserAgent = true
	}
}

// WithHeaders sets default request headers sent by provider requests.
func WithHeaders(headers map[string]string) Option {
	return func(config *Config) {
		if len(headers) == 0 {
			return
		}
		config.Headers = mergeStringMap(config.Headers, headers)
	}
}

// WithRetry sets retry behavior for provider requests.
func WithRetry(retry RetryOptions) Option {
	return func(config *Config) {
		if retry.MaxRetries >= 0 {
			config.Retry.MaxRetries = retry.MaxRetries
		}
		if retry.BaseDelay > 0 {
			config.Retry.BaseDelay = retry.BaseDelay
		}
		if retry.MaxDelay > 0 {
			config.Retry.MaxDelay = retry.MaxDelay
		}
		if retry.BackoffMultiplier > 0 {
			config.Retry.BackoffMultiplier = retry.BackoffMultiplier
		}
		if retry.RetryableStatusCodes != nil {
			config.Retry.RetryableStatusCodes = append([]int(nil), retry.RetryableStatusCodes...)
		}
		if retry.RetryOnNetworkError != nil {
			config.Retry.RetryOnNetworkError = boolPtr(*retry.RetryOnNetworkError)
		}
		if retry.RetryOnTimeout != nil {
			config.Retry.RetryOnTimeout = boolPtr(*retry.RetryOnTimeout)
		}
		if retry.OnRetry != nil {
			config.Retry.OnRetry = retry.OnRetry
		}
	}
}

// WithRateLimit enables token-bucket request throttling.
func WithRateLimit(rateLimit RateLimitOptions) Option {
	return func(config *Config) {
		if rateLimit.RequestsPerSecond > 0 {
			config.RateLimit.RequestsPerSecond = rateLimit.RequestsPerSecond
		}
		if rateLimit.MaxBurst > 0 {
			config.RateLimit.MaxBurst = rateLimit.MaxBurst
		}
	}
}

// WithCircuitBreaker enables provider request circuit breaking.
func WithCircuitBreaker(circuitBreaker CircuitBreakerOptions) Option {
	return func(config *Config) {
		if circuitBreaker.FailureThreshold > 0 {
			config.CircuitBreaker.FailureThreshold = circuitBreaker.FailureThreshold
		}
		if circuitBreaker.ResetTimeout > 0 {
			config.CircuitBreaker.ResetTimeout = circuitBreaker.ResetTimeout
		}
		if circuitBreaker.HalfOpenRequests > 0 {
			config.CircuitBreaker.HalfOpenRequests = circuitBreaker.HalfOpenRequests
		}
		if circuitBreaker.OnStateChange != nil {
			config.CircuitBreaker.OnStateChange = circuitBreaker.OnStateChange
		}
	}
}

// WithProviderPolicy overrides request governance for a specific provider.
func WithProviderPolicy(provider ProviderName, policy ProviderPolicy) Option {
	return func(config *Config) {
		if provider == "" || provider == ProviderUnknown {
			return
		}
		if config.ProviderPolicies == nil {
			config.ProviderPolicies = make(map[ProviderName]ProviderPolicy)
		}
		existing := config.ProviderPolicies[provider]
		if policy.Timeout > 0 {
			existing.Timeout = policy.Timeout
		}
		if policy.UserAgent != "" {
			existing.UserAgent = policy.UserAgent
		}
		if policy.RotateUserAgent != nil {
			existing.RotateUserAgent = cloneBoolPtr(policy.RotateUserAgent)
		}
		if len(policy.Headers) > 0 {
			existing.Headers = mergeStringMap(existing.Headers, policy.Headers)
		}
		if policy.Retry != nil {
			retry := *policy.Retry
			existing.Retry = &retry
		}
		if policy.RateLimit != nil {
			rateLimit := *policy.RateLimit
			existing.RateLimit = &rateLimit
		}
		if policy.CircuitBreaker != nil {
			circuitBreaker := *policy.CircuitBreaker
			existing.CircuitBreaker = &circuitBreaker
		}
		config.ProviderPolicies[provider] = existing
	}
}

// WithRequestHooks configures request lifecycle observers.
func WithRequestHooks(hooks RequestHooks) Option {
	return func(config *Config) {
		config.RequestHooks = hooks
	}
}

func cloneStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

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

func boolPtr(value bool) *bool {
	return &value
}

func cloneBoolPtr(value *bool) *bool {
	if value == nil {
		return nil
	}
	return boolPtr(*value)
}

func boolValue(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func mergeStringMap(base map[string]string, override map[string]string) map[string]string {
	merged := cloneStringMap(base)
	if len(override) == 0 {
		return merged
	}
	if merged == nil {
		merged = make(map[string]string, len(override))
	}
	for key, value := range override {
		merged[key] = value
	}
	return merged
}
