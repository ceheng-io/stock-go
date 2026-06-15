package core

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestHostFallbackManagerBuildsEastmoneyCandidateURLs(t *testing.T) {
	clock := newFakeFallbackClock(time.Unix(0, 0))
	manager := newHostFallbackManagerWithClock(HostFallbackOptions{}, clock)

	candidates := manager.CandidateURLs("https://push2his.eastmoney.com/api/qt/stock/kline/get?secid=1.600519", ProviderEastmoney)

	wantHosts := []string{
		"push2his.eastmoney.com",
		"7.push2his.eastmoney.com",
		"33.push2his.eastmoney.com",
		"63.push2his.eastmoney.com",
		"91.push2his.eastmoney.com",
	}
	if len(candidates) != len(wantHosts) {
		t.Fatalf("len(candidates) = %d, want %d: %#v", len(candidates), len(wantHosts), candidates)
	}
	for i, candidate := range candidates {
		if !strings.Contains(candidate, wantHosts[i]) {
			t.Fatalf("candidate[%d] = %q, want host %q", i, candidate, wantHosts[i])
		}
		if !strings.Contains(candidate, "secid=1.600519") {
			t.Fatalf("candidate[%d] = %q, query was not preserved", i, candidate)
		}
	}
}

func TestHostFallbackManagerPreservesURLPort(t *testing.T) {
	clock := newFakeFallbackClock(time.Unix(0, 0))
	manager := newHostFallbackManagerWithClock(HostFallbackOptions{}, clock)

	candidates := manager.CandidateURLs("https://push2his.eastmoney.com:8443/api/qt/stock/kline/get", ProviderEastmoney)

	if candidateHost(t, candidates[1]) != "7.push2his.eastmoney.com" {
		t.Fatalf("fallback host = %q, want 7.push2his.eastmoney.com", candidateHost(t, candidates[1]))
	}
	if candidatePort(t, candidates[1]) != "8443" {
		t.Fatalf("fallback port = %q, want 8443", candidatePort(t, candidates[1]))
	}
}

func TestHostFallbackManagerSkipsCoolingHostsUntilCooldownExpires(t *testing.T) {
	clock := newFakeFallbackClock(time.Unix(0, 0))
	manager := newHostFallbackManagerWithClock(HostFallbackOptions{
		FailureThreshold: 1,
		Cooldown:         time.Minute,
	}, clock)
	requestURL := "https://push2his.eastmoney.com/api/qt/stock/kline/get"

	manager.RecordFailure(requestURL, errors.New("network down"))

	candidates := manager.CandidateURLs(requestURL, ProviderEastmoney)
	if candidateHost(t, candidates[0]) == "push2his.eastmoney.com" {
		t.Fatalf("first candidate = %q, want original host after healthy fallbacks", candidates[0])
	}
	if candidateHost(t, candidates[len(candidates)-1]) != "push2his.eastmoney.com" {
		t.Fatalf("last candidate = %q, want cooling original host last", candidates[len(candidates)-1])
	}

	clock.Advance(time.Minute)
	candidates = manager.CandidateURLs(requestURL, ProviderEastmoney)
	if candidateHost(t, candidates[0]) != "push2his.eastmoney.com" {
		t.Fatalf("first candidate after cooldown = %q, want original host restored", candidates[0])
	}
}

func TestHostFallbackManagerShouldFallbackOnlyForRetryableErrors(t *testing.T) {
	manager := NewHostFallbackManager(HostFallbackOptions{})

	if !manager.ShouldFallback(errors.New("network down")) {
		t.Fatal("network error should fallback")
	}
	if !manager.ShouldFallback(httpStatusError{statusCode: http.StatusTooManyRequests}) {
		t.Fatal("HTTP 429 should fallback")
	}
	if !manager.ShouldFallback(httpStatusError{statusCode: http.StatusBadGateway}) {
		t.Fatal("HTTP 502 should fallback")
	}
	if manager.ShouldFallback(httpStatusError{statusCode: http.StatusNotFound}) {
		t.Fatal("HTTP 404 should not fallback")
	}
}

func TestHostFallbackManagerRecordsStats(t *testing.T) {
	clock := newFakeFallbackClock(time.Unix(0, 0))
	manager := newHostFallbackManagerWithClock(HostFallbackOptions{
		FailureThreshold: 1,
		Cooldown:         time.Minute,
	}, clock)
	requestURL := "https://17.push2.eastmoney.com/api/qt/clist/get"

	manager.RecordFailure(requestURL, httpStatusError{statusCode: http.StatusServiceUnavailable})
	stats := manager.Stats(ProviderEastmoney)

	if len(stats) != 1 {
		t.Fatalf("len(stats) = %d, want 1", len(stats))
	}
	if stats[0].Host != "17.push2.eastmoney.com" {
		t.Fatalf("host = %q, want 17.push2.eastmoney.com", stats[0].Host)
	}
	if stats[0].FailureCount != 1 || !stats[0].CooldownUntil.Equal(clock.Now().Add(time.Minute)) {
		t.Fatalf("stats after failure = %+v", stats[0])
	}

	manager.RecordSuccess(requestURL)
	stats = manager.Stats(ProviderEastmoney)
	if stats[0].FailureCount != 0 || !stats[0].CooldownUntil.IsZero() || stats[0].SuccessCount != 1 {
		t.Fatalf("stats after success = %+v", stats[0])
	}
}

type fakeFallbackClock struct {
	now time.Time
}

func newFakeFallbackClock(now time.Time) *fakeFallbackClock {
	return &fakeFallbackClock{now: now}
}

func (c *fakeFallbackClock) Now() time.Time {
	return c.now
}

func (c *fakeFallbackClock) Advance(duration time.Duration) {
	c.now = c.now.Add(duration)
}

func candidateHost(t *testing.T, candidate string) string {
	t.Helper()
	parsedURL, err := url.Parse(candidate)
	if err != nil {
		t.Fatal(err)
	}
	return parsedURL.Hostname()
}

func candidatePort(t *testing.T, candidate string) string {
	t.Helper()
	parsedURL, err := url.Parse(candidate)
	if err != nil {
		t.Fatal(err)
	}
	return parsedURL.Port()
}
