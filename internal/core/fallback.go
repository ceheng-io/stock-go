package core

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

// ProviderName identifies the upstream data source for request governance.
type ProviderName string

const (
	ProviderTencent   ProviderName = "tencent"
	ProviderEastmoney ProviderName = "eastmoney"
	ProviderSina      ProviderName = "sina"
	ProviderLinkdiary ProviderName = "linkdiary"
	ProviderUnknown   ProviderName = "unknown"
)

// HostFallbackOptions configures fallback-host health tracking.
type HostFallbackOptions struct {
	Cooldown         time.Duration
	FailureThreshold int
}

// HostHealthStats exposes fallback-host health counters.
type HostHealthStats struct {
	Host          string
	FailureCount  int
	SuccessCount  int
	CooldownUntil time.Time
	LastFailureAt time.Time
	LastError     string
}

// HostFallbackManager orders candidate hosts and tracks temporary host cooldowns.
type HostFallbackManager struct {
	mu               sync.Mutex
	states           map[string]*hostHealthState
	cooldown         time.Duration
	failureThreshold int
	clock            fallbackClock
}

type hostHealthState struct {
	host          string
	failureCount  int
	successCount  int
	cooldownUntil time.Time
	lastFailureAt time.Time
	lastError     string
}

type fallbackClock interface {
	Now() time.Time
}

type realFallbackClock struct{}

func (realFallbackClock) Now() time.Time {
	return time.Now()
}

var eastmoneyPush2HisHosts = []string{
	"push2his.eastmoney.com",
	"7.push2his.eastmoney.com",
	"33.push2his.eastmoney.com",
	"63.push2his.eastmoney.com",
	"91.push2his.eastmoney.com",
}

var eastmoneyPush2Hosts = []string{
	"17.push2.eastmoney.com",
	"29.push2.eastmoney.com",
	"79.push2.eastmoney.com",
	"91.push2.eastmoney.com",
}

// NewHostFallbackManager creates fallback-host governance.
func NewHostFallbackManager(options HostFallbackOptions) *HostFallbackManager {
	return newHostFallbackManagerWithClock(options, realFallbackClock{})
}

func newHostFallbackManagerWithClock(options HostFallbackOptions, clock fallbackClock) *HostFallbackManager {
	cooldown := options.Cooldown
	if cooldown <= 0 {
		cooldown = 30 * time.Second
	}
	failureThreshold := options.FailureThreshold
	if failureThreshold <= 0 {
		failureThreshold = 1
	}
	if clock == nil {
		clock = realFallbackClock{}
	}
	return &HostFallbackManager{
		states:           make(map[string]*hostHealthState),
		cooldown:         cooldown,
		failureThreshold: failureThreshold,
		clock:            clock,
	}
}

// CandidateURLs returns provider-specific fallback candidates for a URL.
func (m *HostFallbackManager) CandidateURLs(requestURL string, provider ProviderName) []string {
	parsedURL, err := url.Parse(requestURL)
	if err != nil || parsedURL.Hostname() == "" {
		return []string{requestURL}
	}
	hostPool := resolveHostPool(parsedURL.Hostname(), provider)
	if len(hostPool) <= 1 {
		return []string{requestURL}
	}

	m.mu.Lock()
	now := m.clock.Now()
	healthyHosts := make([]string, 0, len(hostPool))
	coolingHosts := make([]string, 0, len(hostPool))
	for _, host := range hostPool {
		state := m.states[host]
		if state == nil || !state.cooldownUntil.After(now) {
			healthyHosts = append(healthyHosts, host)
			continue
		}
		coolingHosts = append(coolingHosts, host)
	}
	m.mu.Unlock()

	orderedHosts := uniqueStrings(append(
		orderedWithPreferred(healthyHosts, parsedURL.Hostname()),
		orderedWithPreferred(coolingHosts, parsedURL.Hostname())...,
	))
	candidates := make([]string, 0, len(orderedHosts))
	for _, host := range orderedHosts {
		candidate := *parsedURL
		candidate.Host = replaceURLHost(candidate.Host, host)
		candidates = append(candidates, candidate.String())
	}
	return candidates
}

// RecordSuccess resets the host failure state after a successful request.
func (m *HostFallbackManager) RecordSuccess(requestURL string) {
	host := safeHost(requestURL)
	if host == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	state := m.getStateLocked(host)
	state.failureCount = 0
	state.cooldownUntil = time.Time{}
	state.successCount++
	state.lastError = ""
}

// RecordFailure records a failed host attempt and may place it in cooldown.
func (m *HostFallbackManager) RecordFailure(requestURL string, err error) {
	host := safeHost(requestURL)
	if host == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	state := m.getStateLocked(host)
	state.failureCount++
	state.lastFailureAt = m.clock.Now()
	if err != nil {
		state.lastError = err.Error()
	}
	if state.failureCount >= m.failureThreshold {
		state.cooldownUntil = m.clock.Now().Add(m.cooldown)
	}
}

// ShouldFallback reports whether an error is worth trying on another host.
func (m *HostFallbackManager) ShouldFallback(err error) bool {
	if err == nil {
		return false
	}
	var statusErr httpStatusError
	if errors.As(err, &statusErr) {
		return isRetryableStatus(statusErr.statusCode, defaultRetryableStatusCodes())
	}
	return true
}

// Stats returns current host health snapshots.
func (m *HostFallbackManager) Stats(provider ProviderName) []HostHealthStats {
	m.mu.Lock()
	defer m.mu.Unlock()

	allowedHosts := map[string]struct{}{}
	if provider == ProviderEastmoney {
		for _, host := range eastmoneyPush2HisHosts {
			allowedHosts[host] = struct{}{}
		}
		for _, host := range eastmoneyPush2Hosts {
			allowedHosts[host] = struct{}{}
		}
	}

	stats := make([]HostHealthStats, 0, len(m.states))
	for host, state := range m.states {
		if provider != "" && provider != ProviderUnknown {
			if _, ok := allowedHosts[host]; !ok {
				continue
			}
		}
		stats = append(stats, HostHealthStats{
			Host:          state.host,
			FailureCount:  state.failureCount,
			SuccessCount:  state.successCount,
			CooldownUntil: state.cooldownUntil,
			LastFailureAt: state.lastFailureAt,
			LastError:     state.lastError,
		})
	}
	return stats
}

func (m *HostFallbackManager) getStateLocked(host string) *hostHealthState {
	state := m.states[host]
	if state != nil {
		return state
	}
	state = &hostHealthState{host: host}
	m.states[host] = state
	return state
}

func inferProviderFromURL(requestURL string) ProviderName {
	host := safeHost(requestURL)
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

func resolveHostPool(host string, provider ProviderName) []string {
	if provider != ProviderEastmoney {
		return []string{host}
	}
	if strings.Contains(host, "push2his.eastmoney.com") {
		return eastmoneyPush2HisHosts
	}
	if strings.Contains(host, "push2.eastmoney.com") {
		return eastmoneyPush2Hosts
	}
	return []string{host}
}

func orderedWithPreferred(hosts []string, preferred string) []string {
	ordered := make([]string, 0, len(hosts))
	for _, host := range hosts {
		if host == preferred {
			ordered = append(ordered, host)
			break
		}
	}
	for _, host := range hosts {
		if host != preferred {
			ordered = append(ordered, host)
		}
	}
	return ordered
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	unique := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		unique = append(unique, item)
	}
	return unique
}

func safeHost(requestURL string) string {
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

func replaceURLHost(originalHost string, hostname string) string {
	if _, port, err := net.SplitHostPort(originalHost); err == nil {
		return net.JoinHostPort(hostname, port)
	}
	return hostname
}
