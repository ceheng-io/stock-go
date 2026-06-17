package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ceheng.io/stock-go/parser"
	"github.com/ceheng.io/stock-go/useragent"
)

const maxDrainBodyBytes int64 = 512 << 10

// Config contains internal request client configuration.
type Config struct {
	BaseURL           string
	TencentMinuteURL  string
	SearchBaseURL     string
	CalendarURL       string
	AShareListURL     string
	USListURL         string
	HKListURL         string
	FundListURL       string
	EastmoneyKlineURL string
	HTTPClient        *http.Client
	ProxyPool         ProxyPoolConfig
	Timeout           time.Duration
	UserAgent         string
	RotateUserAgent   bool
	Headers           map[string]string
	Retry             RetryConfig
	RateLimiter       RequestLimiter
	CircuitBreaker    RequestCircuitBreaker
	HostFallback      *HostFallbackManager
	ProviderPolicies  map[ProviderName]ProviderPolicy
	Hooks             RequestHooks
	EastmoneySession  EastmoneySessionConfig
}

// RetryConfig controls internal request retries.
type RetryConfig struct {
	MaxRetries           int
	BaseDelay            time.Duration
	MaxDelay             time.Duration
	BackoffMultiplier    float64
	RetryableStatusCodes []int
	RetryOnNetworkError  *bool
	RetryOnTimeout       *bool
	OnRetry              func(attempt int, err error, delay time.Duration)
}

// ProviderPolicy overrides request governance for a specific provider.
type ProviderPolicy struct {
	Timeout         time.Duration
	UserAgent       string
	RotateUserAgent *bool
	Headers         map[string]string
	Retry           *RetryConfig
	RateLimiter     RequestLimiter
	CircuitBreaker  RequestCircuitBreaker
}

// EastmoneySessionConfig controls optional browser-session bootstrapping for Eastmoney requests.
type EastmoneySessionConfig struct {
	AutoInit  bool
	InitURL   string
	UserAgent string
	Headers   map[string]string
}

// ResolvedProviderPolicy contains complete governance settings for one provider.
type ResolvedProviderPolicy struct {
	Timeout         time.Duration
	UserAgent       string
	RotateUserAgent bool
	Headers         map[string]string
	Retry           RetryConfig
	RateLimiter     RequestLimiter
	CircuitBreaker  RequestCircuitBreaker
}

// RequestTraceEvent describes a request lifecycle trace event.
type RequestTraceEvent string

const (
	TraceRequest  RequestTraceEvent = "request"
	TraceResponse RequestTraceEvent = "response"
	TraceError    RequestTraceEvent = "error"
	TraceRetry    RequestTraceEvent = "retry"
	TraceFallback RequestTraceEvent = "fallback"
)

// RequestContext describes one request lifecycle event.
type RequestContext struct {
	Provider ProviderName
	URL      string
	Timeout  time.Duration
	Attempt  int
}

// ResponseMeta describes response timing metadata.
type ResponseMeta struct {
	StatusCode int
	Duration   time.Duration
}

// RequestHooks observes request lifecycle events.
type RequestHooks struct {
	OnRequest  func(RequestContext)
	OnResponse func(RequestContext, ResponseMeta)
	OnError    func(RequestContext, error)
	OnRetry    func(RequestContext, error, time.Duration)
	Trace      func(RequestTraceEvent, RequestContext)
}

// RequestLimiter gates outgoing HTTP attempts.
type RequestLimiter interface {
	Acquire(context.Context) error
}

// RequestCircuitBreaker gates requests after repeated provider failures.
type RequestCircuitBreaker interface {
	CanRequest() bool
	RecordSuccess()
	RecordFailure()
}

// Client performs provider HTTP requests.
type Client struct {
	baseURL           string
	tencentMinuteURL  string
	searchBaseURL     string
	calendarURL       string
	aShareListURL     string
	usListURL         string
	hkListURL         string
	fundListURL       string
	eastmoneyKlineURL string
	httpClient        *http.Client
	timeout           time.Duration
	userAgent         string
	rotateUserAgent   bool
	headers           map[string]string
	retry             RetryConfig
	rateLimiter       RequestLimiter
	circuitBreaker    RequestCircuitBreaker
	hostFallback      *HostFallbackManager
	providerPolicies  map[ProviderName]ResolvedProviderPolicy
	hooks             RequestHooks
	eastmoneySession  EastmoneySessionConfig
	eastmoneyMu       sync.Mutex
	eastmoneyCookies  map[string]string
}

// NewClient creates an internal request client.
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = "https://qt.gtimg.cn"
	}
	if config.SearchBaseURL == "" {
		config.SearchBaseURL = "https://smartbox.gtimg.cn/s3/"
	}
	if config.TencentMinuteURL == "" {
		config.TencentMinuteURL = "https://web.ifzq.gtimg.cn/appstock/app/minute/query"
	}
	if config.CalendarURL == "" {
		config.CalendarURL = "https://proxy.finance.qq.com/ifzqgtimg/appstock/app/newfqkline/get?param=calendar"
	}
	if config.AShareListURL == "" {
		config.AShareListURL = "https://assets.linkdiary.cn/shares/zh_a_list.json"
	}
	if config.USListURL == "" {
		config.USListURL = "https://assets.linkdiary.cn/shares/us_list.json"
	}
	if config.HKListURL == "" {
		config.HKListURL = "https://assets.linkdiary.cn/shares/hk_list.json"
	}
	if config.FundListURL == "" {
		config.FundListURL = "https://assets.linkdiary.cn/shares/fund_list"
	}
	if config.EastmoneyKlineURL == "" {
		config.EastmoneyKlineURL = "https://push2his.eastmoney.com/api/qt/stock/kline/get"
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}
	config.HTTPClient = httpClientWithProxyPool(config.HTTPClient, NewProxyPool(config.ProxyPool.URLs))
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}
	if config.Retry.MaxRetries < 0 {
		config.Retry.MaxRetries = 0
	}
	if config.Retry.BaseDelay <= 0 {
		config.Retry.BaseDelay = time.Second
	}
	if config.Retry.MaxDelay <= 0 {
		config.Retry.MaxDelay = 30 * time.Second
	}
	if config.Retry.BackoffMultiplier <= 0 {
		config.Retry.BackoffMultiplier = 2
	}
	if config.Retry.RetryableStatusCodes == nil {
		config.Retry.RetryableStatusCodes = defaultRetryableStatusCodes()
	}
	if config.Retry.RetryOnNetworkError == nil {
		config.Retry.RetryOnNetworkError = boolPtr(true)
	}
	if config.Retry.RetryOnTimeout == nil {
		config.Retry.RetryOnTimeout = boolPtr(true)
	}
	if config.HostFallback == nil {
		config.HostFallback = NewHostFallbackManager(HostFallbackOptions{})
	}
	defaultPolicy := ResolvedProviderPolicy{
		Timeout:         config.Timeout,
		UserAgent:       config.UserAgent,
		RotateUserAgent: config.RotateUserAgent,
		Headers:         cloneStringMap(config.Headers),
		Retry:           config.Retry,
		RateLimiter:     config.RateLimiter,
		CircuitBreaker:  config.CircuitBreaker,
	}
	return &Client{
		baseURL:           strings.TrimRight(config.BaseURL, "/"),
		tencentMinuteURL:  config.TencentMinuteURL,
		searchBaseURL:     config.SearchBaseURL,
		calendarURL:       config.CalendarURL,
		aShareListURL:     config.AShareListURL,
		usListURL:         config.USListURL,
		hkListURL:         config.HKListURL,
		fundListURL:       config.FundListURL,
		eastmoneyKlineURL: config.EastmoneyKlineURL,
		httpClient:        config.HTTPClient,
		timeout:           config.Timeout,
		userAgent:         config.UserAgent,
		rotateUserAgent:   config.RotateUserAgent,
		headers:           cloneStringMap(config.Headers),
		retry:             config.Retry,
		rateLimiter:       config.RateLimiter,
		circuitBreaker:    config.CircuitBreaker,
		hostFallback:      config.HostFallback,
		providerPolicies:  resolveProviderPolicies(defaultPolicy, config.ProviderPolicies),
		hooks:             config.Hooks,
		eastmoneySession:  normalizeEastmoneySessionConfig(config.EastmoneySession),
		eastmoneyCookies:  make(map[string]string),
	}
}

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

// TencentMinuteURL returns the configured Tencent minute URL.
func (c *Client) TencentMinuteURL() string { return c.tencentMinuteURL }

// CalendarURL returns the configured trading calendar URL.
func (c *Client) CalendarURL() string {
	return c.calendarURL
}

// AShareListURL returns the configured A-share code list URL.
func (c *Client) AShareListURL() string { return c.aShareListURL }

// USListURL returns the configured US code list URL.
func (c *Client) USListURL() string { return c.usListURL }

// HKListURL returns the configured HK code list URL.
func (c *Client) HKListURL() string { return c.hkListURL }

// FundListURL returns the configured fund code list URL.
func (c *Client) FundListURL() string { return c.fundListURL }

// EastmoneyKlineURL returns the configured Eastmoney K-line URL.
func (c *Client) EastmoneyKlineURL() string { return c.eastmoneyKlineURL }

// RateLimiter returns the configured request limiter.
func (c *Client) RateLimiter() RequestLimiter { return c.rateLimiter }

// CircuitBreaker returns the configured request circuit breaker.
func (c *Client) CircuitBreaker() RequestCircuitBreaker { return c.circuitBreaker }

// HostFallback returns the configured fallback-host manager.
func (c *Client) HostFallback() *HostFallbackManager { return c.hostFallback }

// ProviderPolicy returns the resolved policy for a provider.
func (c *Client) ProviderPolicy(provider ProviderName) (ResolvedProviderPolicy, bool) {
	policy, ok := c.providerPolicies[provider]
	return policy, ok
}

// Hooks returns the configured request lifecycle hooks.
func (c *Client) Hooks() RequestHooks { return c.hooks }

// EastmoneySession returns the configured Eastmoney browser-session options.
func (c *Client) EastmoneySession() EastmoneySessionConfig {
	session := c.eastmoneySession
	session.Headers = cloneStringMap(session.Headers)
	return session
}

// GetTencentQuote fetches and parses Tencent quote response text.
func (c *Client) GetTencentQuote(ctx context.Context, params string) ([]TencentQuoteItem, error) {
	requestURL := c.baseURL + "/?q=" + url.QueryEscape(params)
	body, err := c.GetBytes(ctx, requestURL)
	if err != nil {
		return nil, err
	}
	text, err := decodeGBK(body)
	if err != nil {
		return nil, err
	}
	return ParseTencentQuoteResponse(text), nil
}

// TencentSearchURL builds the Smartbox search URL for a keyword.
func (c *Client) TencentSearchURL(keyword string) string {
	return c.searchBaseURL + "?v=2&t=all&q=" + url.QueryEscape(keyword)
}

// GetText performs an HTTP GET request and returns response text.
func (c *Client) GetText(ctx context.Context, requestURL string) (string, error) {
	body, err := c.GetBytes(ctx, requestURL)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetBytes performs an HTTP GET request and returns raw response bytes.
func (c *Client) GetBytes(ctx context.Context, requestURL string) ([]byte, error) {
	provider := inferProviderFromURL(requestURL)
	policy := c.policyForProvider(provider)
	ctx, cancel := context.WithTimeout(ctx, policy.Timeout)
	defer cancel()

	return c.getBytes(ctx, requestURL, provider, policy)
}

func (c *Client) getBytes(ctx context.Context, requestURL string, provider ProviderName, policy ResolvedProviderPolicy) ([]byte, error) {
	if policy.CircuitBreaker != nil && !policy.CircuitBreaker.CanRequest() {
		return nil, ErrCircuitBreakerOpen
	}

	candidates := []string{requestURL}
	if c.hostFallback != nil {
		candidates = c.hostFallback.CandidateURLs(requestURL, provider)
	}

	var lastErr error
	for index, candidateURL := range candidates {
		retry := policy.Retry
		if index > 0 {
			retry.MaxRetries = 0
		}
		body, err := c.getBytesWithRetry(ctx, candidateURL, retry, policy, provider)
		if err == nil {
			if policy.CircuitBreaker != nil {
				policy.CircuitBreaker.RecordSuccess()
			}
			if c.hostFallback != nil {
				c.hostFallback.RecordSuccess(candidateURL)
			}
			return body, nil
		}
		if c.shouldInitEastmoneySession(candidateURL, provider, policy, err) && c.initializeEastmoneySession(ctx, policy) {
			sessionRetry := policy.Retry
			sessionRetry.MaxRetries = 0
			body, err = c.getBytesWithRetry(ctx, candidateURL, sessionRetry, policy, provider)
			if err == nil {
				if policy.CircuitBreaker != nil {
					policy.CircuitBreaker.RecordSuccess()
				}
				if c.hostFallback != nil {
					c.hostFallback.RecordSuccess(candidateURL)
				}
				return body, nil
			}
		}
		lastErr = err
		if c.hostFallback != nil {
			c.hostFallback.RecordFailure(candidateURL, err)
		}
		if index < len(candidates)-1 && c.hostFallback != nil && c.hostFallback.ShouldFallback(err) {
			c.trace(TraceFallback, RequestContext{
				Provider: provider,
				URL:      candidateURL,
				Timeout:  policy.Timeout,
			})
			continue
		}
		if policy.CircuitBreaker != nil {
			policy.CircuitBreaker.RecordFailure()
		}
		return nil, err
	}
	if policy.CircuitBreaker != nil {
		policy.CircuitBreaker.RecordFailure()
	}
	return nil, lastErr
}

func (c *Client) getBytesWithRetry(ctx context.Context, requestURL string, retry RetryConfig, policy ResolvedProviderPolicy, provider ProviderName) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= retry.MaxRetries; attempt++ {
		body, err := c.getBytesOnce(ctx, requestURL, policy, provider, attempt)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !c.shouldRetry(err, attempt, retry) {
			return nil, err
		}
		delay := retryDelay(attempt, retry)
		requestContext := RequestContext{
			Provider: provider,
			URL:      requestURL,
			Timeout:  policy.Timeout,
			Attempt:  attempt,
		}
		c.safe(func() {
			if retry.OnRetry != nil {
				retry.OnRetry(attempt+1, err, delay)
			}
		})
		c.safe(func() {
			if c.hooks.OnRetry != nil {
				c.hooks.OnRetry(requestContext, err, delay)
			}
		})
		c.trace(TraceRetry, requestContext)
		if err := sleepContext(ctx, delay); err != nil {
			return nil, err
		}
	}
	return nil, lastErr
}

func (c *Client) getBytesOnce(ctx context.Context, requestURL string, policy ResolvedProviderPolicy, provider ProviderName, attempt int) ([]byte, error) {
	if policy.RateLimiter != nil {
		if err := policy.RateLimiter.Acquire(ctx); err != nil {
			return nil, codedRequestError(err, provider, requestURL, policy.Timeout)
		}
	}
	policy = c.withEastmoneySessionCookies(requestURL, provider, policy)
	requestContext := RequestContext{
		Provider: provider,
		URL:      requestURL,
		Timeout:  policy.Timeout,
		Attempt:  attempt,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		err = codedRequestError(err, provider, requestURL, policy.Timeout)
		c.emitError(requestContext, err)
		return nil, err
	}
	for key, value := range policy.Headers {
		req.Header.Set(key, value)
	}
	if policy.UserAgent != "" {
		req.Header.Set("User-Agent", policy.UserAgent)
	}
	if policy.RotateUserAgent {
		if userAgent := useragent.Next(); userAgent != "" {
			req.Header.Set("User-Agent", userAgent)
		}
	}

	c.safe(func() {
		if c.hooks.OnRequest != nil {
			c.hooks.OnRequest(requestContext)
		}
	})
	c.trace(TraceRequest, requestContext)
	startedAt := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = codedRequestError(err, provider, requestURL, policy.Timeout)
		c.emitError(requestContext, err)
		return nil, err
	}
	defer drainAndClose(resp.Body)

	c.safe(func() {
		if c.hooks.OnResponse != nil {
			c.hooks.OnResponse(requestContext, ResponseMeta{
				StatusCode: resp.StatusCode,
				Duration:   time.Since(startedAt),
			})
		}
	})
	c.trace(TraceResponse, requestContext)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := codedRequestError(httpStatusError{
			statusCode: resp.StatusCode,
			statusText: resp.Status,
			url:        requestURL,
		}, provider, requestURL, policy.Timeout)
		c.emitError(requestContext, err)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = codedRequestError(err, provider, requestURL, policy.Timeout)
		c.emitError(requestContext, err)
		return nil, err
	}
	return body, nil
}

func drainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(body, maxDrainBodyBytes))
	_ = body.Close()
}

type httpStatusError struct {
	statusCode int
	statusText string
	url        string
}

func (e httpStatusError) Error() string {
	return fmt.Sprintf("http status %d from %s", e.statusCode, e.url)
}

func (c *Client) shouldRetry(err error, attempt int, retry RetryConfig) bool {
	if attempt >= retry.MaxRetries {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return boolValue(retry.RetryOnTimeout, true)
	}
	var statusErr httpStatusError
	if errors.As(err, &statusErr) {
		return isRetryableStatus(statusErr.statusCode, retry.RetryableStatusCodes)
	}
	return boolValue(retry.RetryOnNetworkError, true)
}

func isRetryableStatus(statusCode int, retryableStatusCodes []int) bool {
	for _, retryable := range retryableStatusCodes {
		if statusCode == retryable {
			return true
		}
	}
	return false
}

func defaultRetryableStatusCodes() []int {
	return []int{http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout}
}

func retryDelay(attempt int, retry RetryConfig) time.Duration {
	delay := float64(retry.BaseDelay)
	for i := 0; i < attempt; i++ {
		delay *= retry.BackoffMultiplier
	}
	if maxDelay := float64(retry.MaxDelay); maxDelay > 0 && delay > maxDelay {
		delay = maxDelay
	}
	return time.Duration(delay)
}

func sleepContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// GetJSON performs an HTTP GET request and decodes a JSON response.
func (c *Client) GetJSON(ctx context.Context, requestURL string, target any) error {
	text, err := c.GetText(ctx, requestURL)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(text), target); err != nil {
		return parseRequestError(err, inferProviderFromURL(requestURL), requestURL)
	}
	return nil
}

func codedRequestError(err error, provider ProviderName, requestURL string, timeout time.Duration) error {
	if err == nil {
		return nil
	}
	var coded CodedError
	if errors.As(err, &coded) {
		return err
	}
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return NewCodedError("TIMEOUT", fmt.Sprintf("Request timeout after %s: %s", timeout, requestURL), err)
	case errors.Is(err, context.Canceled):
		return NewCodedError("ABORTED", fmt.Sprintf("Request aborted: %s", requestURL), err)
	}
	var statusErr httpStatusError
	if errors.As(err, &statusErr) {
		code := "HTTP_ERROR"
		if statusErr.statusCode == http.StatusTooManyRequests {
			code = "RATE_LIMITED"
		}
		return NewCodedError(code, statusErr.Error(), err)
	}
	return NewCodedError("NETWORK_ERROR", fmt.Sprintf("Network request failed for %s: %s", provider, requestURL), err)
}

func parseRequestError(err error, provider ProviderName, requestURL string) error {
	if err == nil {
		return nil
	}
	return NewCodedError("PARSE_ERROR", fmt.Sprintf("Failed to parse JSON response from %s: %s", provider, requestURL), err)
}

const (
	defaultEastmoneySessionInitURL   = "https://quote.eastmoney.com/center/boardlist.html"
	defaultEastmoneySessionUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36"
)

func normalizeEastmoneySessionConfig(config EastmoneySessionConfig) EastmoneySessionConfig {
	if config.InitURL == "" {
		config.InitURL = defaultEastmoneySessionInitURL
	}
	if config.UserAgent == "" {
		config.UserAgent = defaultEastmoneySessionUserAgent
	}
	config.Headers = cloneStringMap(config.Headers)
	if config.Headers == nil {
		config.Headers = map[string]string{}
	}
	if !hasHeader(config.Headers, "Accept") {
		config.Headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"
	}
	if !hasHeader(config.Headers, "Accept-Language") {
		config.Headers["Accept-Language"] = "zh-CN,zh;q=0.9"
	}
	if !hasHeader(config.Headers, "Cache-Control") {
		config.Headers["Cache-Control"] = "no-cache"
	}
	if !hasHeader(config.Headers, "Pragma") {
		config.Headers["Pragma"] = "no-cache"
	}
	return config
}

func (c *Client) shouldInitEastmoneySession(requestURL string, provider ProviderName, policy ResolvedProviderPolicy, err error) bool {
	if provider != ProviderEastmoney || !c.eastmoneySession.AutoInit || err == nil {
		return false
	}
	if !strings.Contains(safeHost(requestURL), "eastmoney.com") || hasHeader(policy.Headers, "Cookie") {
		return false
	}
	var coded CodedError
	if errors.As(err, &coded) {
		switch coded.SDKCode() {
		case "NETWORK_ERROR", "TIMEOUT", "HTTP_ERROR", "RATE_LIMITED":
			return true
		default:
			return false
		}
	}
	var statusErr httpStatusError
	if errors.As(err, &statusErr) {
		return isRetryableStatus(statusErr.statusCode, defaultRetryableStatusCodes()) || statusErr.statusCode == http.StatusForbidden || statusErr.statusCode == http.StatusUnauthorized
	}
	return true
}

func (c *Client) initializeEastmoneySession(ctx context.Context, policy ResolvedProviderPolicy) bool {
	config := c.eastmoneySession
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.InitURL, nil)
	if err != nil {
		return false
	}
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}
	if !hasHeader(config.Headers, "User-Agent") && config.UserAgent != "" {
		req.Header.Set("User-Agent", config.UserAgent)
	}
	for key, value := range policy.Headers {
		if strings.EqualFold(key, "Cookie") || req.Header.Get(key) != "" {
			continue
		}
		req.Header.Set(key, value)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer drainAndClose(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false
	}
	return c.storeEastmoneyCookies(resp.Cookies())
}

func (c *Client) withEastmoneySessionCookies(requestURL string, provider ProviderName, policy ResolvedProviderPolicy) ResolvedProviderPolicy {
	if provider != ProviderEastmoney || !strings.Contains(safeHost(requestURL), "eastmoney.com") || hasHeader(policy.Headers, "Cookie") {
		return policy
	}
	cookieHeader := c.eastmoneyCookieHeader()
	if cookieHeader == "" {
		return policy
	}
	policy.Headers = mergeStringMap(policy.Headers, map[string]string{"Cookie": cookieHeader})
	return policy
}

func (c *Client) storeEastmoneyCookies(cookies []*http.Cookie) bool {
	if len(cookies) == 0 {
		return false
	}
	c.eastmoneyMu.Lock()
	defer c.eastmoneyMu.Unlock()
	for _, cookie := range cookies {
		if cookie == nil || cookie.Name == "" || cookie.Value == "" {
			continue
		}
		c.eastmoneyCookies[cookie.Name] = cookie.Value
	}
	return len(c.eastmoneyCookies) > 0
}

func (c *Client) eastmoneyCookieHeader() string {
	c.eastmoneyMu.Lock()
	defer c.eastmoneyMu.Unlock()
	if len(c.eastmoneyCookies) == 0 {
		return ""
	}
	names := make([]string, 0, len(c.eastmoneyCookies))
	for name := range c.eastmoneyCookies {
		names = append(names, name)
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, name+"="+c.eastmoneyCookies[name])
	}
	return strings.Join(parts, "; ")
}

func decodeGBK(data []byte) (string, error) {
	return parser.DecodeGBK(data)
}

func (c *Client) policyForProvider(provider ProviderName) ResolvedProviderPolicy {
	if policy, ok := c.providerPolicies[provider]; ok {
		policy.Headers = cloneStringMap(policy.Headers)
		return policy
	}
	return ResolvedProviderPolicy{
		Timeout:         c.timeout,
		UserAgent:       c.userAgent,
		RotateUserAgent: c.rotateUserAgent,
		Headers:         cloneStringMap(c.headers),
		Retry:           c.retry,
		RateLimiter:     c.rateLimiter,
		CircuitBreaker:  c.circuitBreaker,
	}
}

func resolveProviderPolicies(defaultPolicy ResolvedProviderPolicy, policies map[ProviderName]ProviderPolicy) map[ProviderName]ResolvedProviderPolicy {
	if len(policies) == 0 {
		return nil
	}
	resolved := make(map[ProviderName]ResolvedProviderPolicy, len(policies))
	for provider, policy := range policies {
		resolved[provider] = mergeProviderPolicy(defaultPolicy, policy)
	}
	return resolved
}

func mergeProviderPolicy(base ResolvedProviderPolicy, override ProviderPolicy) ResolvedProviderPolicy {
	merged := base
	merged.Headers = cloneStringMap(base.Headers)
	if override.Timeout > 0 {
		merged.Timeout = override.Timeout
	}
	if override.UserAgent != "" {
		merged.UserAgent = override.UserAgent
	}
	if override.RotateUserAgent != nil {
		merged.RotateUserAgent = *override.RotateUserAgent
	}
	if len(override.Headers) > 0 {
		merged.Headers = mergeStringMap(merged.Headers, override.Headers)
	}
	if override.Retry != nil {
		retry := merged.Retry
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
			retry.RetryOnNetworkError = cloneBoolPtr(override.Retry.RetryOnNetworkError)
		}
		if override.Retry.RetryOnTimeout != nil {
			retry.RetryOnTimeout = cloneBoolPtr(override.Retry.RetryOnTimeout)
		}
		if override.Retry.OnRetry != nil {
			retry.OnRetry = override.Retry.OnRetry
		}
		merged.Retry = retry
	}
	if override.RateLimiter != nil {
		merged.RateLimiter = override.RateLimiter
	}
	if override.CircuitBreaker != nil {
		merged.CircuitBreaker = override.CircuitBreaker
	}
	return merged
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

func hasHeader(headers map[string]string, key string) bool {
	for existing := range headers {
		if strings.EqualFold(existing, key) {
			return true
		}
	}
	return false
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

func (c *Client) emitError(ctx RequestContext, err error) {
	c.safe(func() {
		if c.hooks.OnError != nil {
			c.hooks.OnError(ctx, err)
		}
	})
	c.trace(TraceError, ctx)
}

func (c *Client) trace(event RequestTraceEvent, ctx RequestContext) {
	c.safe(func() {
		if c.hooks.Trace != nil {
			c.hooks.Trace(event, ctx)
		}
	})
}

func (c *Client) safe(fn func()) {
	defer func() {
		_ = recover()
	}()
	fn()
}
