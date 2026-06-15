package stock

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootReExportsUserAgentPool(t *testing.T) {
	ResetUserAgents()
	values := AllUserAgents()
	if len(values) == 0 {
		t.Fatal("AllUserAgents returned empty list")
	}
	if got := NextUserAgent(); got != values[0] {
		t.Fatalf("NextUserAgent = %q, want first known value", got)
	}
	known := make(map[string]struct{}, len(values))
	for _, value := range values {
		known[value] = struct{}{}
	}
	if got := RandomUserAgent(); got == "" {
		t.Fatal("RandomUserAgent returned empty value")
	} else if _, ok := known[got]; !ok {
		t.Fatalf("RandomUserAgent = %q, want known value", got)
	}
}

func TestRootReExportsTSUserAgentPoolNames(t *testing.T) {
	ResetUserAgents()
	values := GetAllUserAgents()
	if len(values) == 0 {
		t.Fatal("GetAllUserAgents returned empty list")
	}
	if got := GetNextUserAgent(); got != values[0] {
		t.Fatalf("GetNextUserAgent = %q, want first known value", got)
	}
	known := make(map[string]struct{}, len(values))
	for _, value := range values {
		known[value] = struct{}{}
	}
	if got := GetRandomUserAgent(); got == "" {
		t.Fatal("GetRandomUserAgent returned empty value")
	} else if _, ok := known[got]; !ok {
		t.Fatalf("GetRandomUserAgent = %q, want known value", got)
	}
}

func TestNewCanUseUserAgentPoolOptions(t *testing.T) {
	ResetUserAgents()
	values := AllUserAgents()

	nextClient := New(WithNextUserAgent())
	if nextClient.config.UserAgent != values[0] {
		t.Fatalf("next user agent = %q, want %q", nextClient.config.UserAgent, values[0])
	}

	randomClient := New(WithRandomUserAgent())
	if randomClient.config.UserAgent == "" {
		t.Fatal("random user agent is empty")
	}
}

func TestNewCanEnableUserAgentRotation(t *testing.T) {
	ResetUserAgents()
	values := AllUserAgents()
	seen := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("User-Agent"))
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()
	client := New(WithRotateUserAgent())

	if !client.config.RotateUserAgent {
		t.Fatal("RotateUserAgent = false, want true")
	}
	if _, err := client.core.GetText(context.Background(), server.URL); err != nil {
		t.Fatal(err)
	}
	if len(seen) != 1 || seen[0] != values[0] {
		t.Fatalf("seen user agents = %#v, want %q", seen, values[0])
	}
}
