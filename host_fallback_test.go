package stock

import (
	"errors"
	"net/url"
	"testing"
	"time"
)

func TestRootHostFallbackManagerExportsFallbackGovernance(t *testing.T) {
	manager := NewHostFallbackManager(HostFallbackOptions{
		Cooldown:         time.Minute,
		FailureThreshold: 1,
	})
	requestURL := "https://push2his.eastmoney.com/api/qt/stock/kline/get?secid=1.600519"

	candidates := manager.CandidateURLs(requestURL, ProviderEastmoney)
	if len(candidates) < 2 {
		t.Fatalf("len(candidates) = %d, want fallback candidates", len(candidates))
	}
	if hostFromURL(t, candidates[0]) != "push2his.eastmoney.com" {
		t.Fatalf("first candidate host = %q, want original host", hostFromURL(t, candidates[0]))
	}

	manager.RecordFailure(requestURL, errors.New("network down"))
	candidates = manager.CandidateURLs(requestURL, ProviderEastmoney)
	if hostFromURL(t, candidates[0]) == "push2his.eastmoney.com" {
		t.Fatalf("first candidate after failure = %q, want a healthy fallback host", candidates[0])
	}

	stats := manager.Stats(ProviderEastmoney)
	if len(stats) != 1 {
		t.Fatalf("len(stats) = %d, want 1", len(stats))
	}
	if stats[0].Host != "push2his.eastmoney.com" || stats[0].FailureCount != 1 || stats[0].LastError == "" {
		t.Fatalf("stats after failure = %+v", stats[0])
	}

	manager.RecordSuccess(requestURL)
	stats = manager.Stats(ProviderEastmoney)
	if stats[0].FailureCount != 0 || !stats[0].CooldownUntil.IsZero() || stats[0].SuccessCount != 1 {
		t.Fatalf("stats after success = %+v", stats[0])
	}
}

func hostFromURL(t *testing.T, rawURL string) string {
	t.Helper()
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	return parsedURL.Hostname()
}
