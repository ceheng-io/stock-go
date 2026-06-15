package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExtractJSVarParsesArraysAndObjects(t *testing.T) {
	text := `
var fS_code = "110011";
var Data_netWorthTrend = [{"x":1702857600000,"y":3.5,"equityReturn":"1.2","unitMoney":""}];
var Data_ACWorthTrend = [[1702857600000,5.6]];
`
	var code string
	if err := ExtractJSVar(text, "fS_code", &code); err != nil {
		t.Fatal(err)
	}
	if code != "110011" {
		t.Fatalf("code = %q", code)
	}
	var trend []map[string]any
	if err := ExtractJSVar(text, "Data_netWorthTrend", &trend); err != nil {
		t.Fatal(err)
	}
	if len(trend) != 1 || trend[0]["y"] != 3.5 {
		t.Fatalf("trend = %+v", trend)
	}
}

func TestExtractJSVarReturnsErrorWhenMissing(t *testing.T) {
	var value string
	if err := ExtractJSVar(`var a = 1;`, "missing", &value); err == nil {
		t.Fatal("expected missing variable error")
	}
}

func TestParseJSVarsSkipsMissingAndInvalidValues(t *testing.T) {
	text := `
var code = "110011";
const rows = [{"name":"易方达中小盘","value":3.5}];
let invalid = {not_json: true};
`

	values := ParseJSVars(text, "code", "rows", "missing", "invalid")
	if len(values) != 2 {
		t.Fatalf("len(values) = %d, want 2: %#v", len(values), values)
	}
	if values["code"] != "110011" {
		t.Fatalf("code = %#v, want 110011", values["code"])
	}
	rows, ok := values["rows"].([]any)
	if !ok || len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row slice", values["rows"])
	}
	if _, ok := values["missing"]; ok {
		t.Fatalf("missing variable should be skipped: %#v", values)
	}
	if _, ok := values["invalid"]; ok {
		t.Fatalf("invalid variable should be skipped: %#v", values)
	}
}

func TestFetchJSVarsFetchesTextAndParsesVariables(t *testing.T) {
	var gotHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Test")
		_, _ = w.Write([]byte(`
var code = "110011";
const rows = [{"name":"易方达中小盘","value":3.5}];
let invalid = {not_json: true};
`))
	}))
	defer server.Close()

	values, err := FetchJSVars(context.Background(), server.Client(), server.URL, []string{"code", "rows", "invalid"}, FetchJSVarsOptions{
		Timeout: time.Second,
		Headers: map[string]string{"X-Test": "yes"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotHeader != "yes" {
		t.Fatalf("X-Test header = %q, want yes", gotHeader)
	}
	if values["code"] != "110011" {
		t.Fatalf("code = %#v, want 110011", values["code"])
	}
	rows, ok := values["rows"].([]any)
	if !ok || len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", values["rows"])
	}
	if _, ok := values["invalid"]; ok {
		t.Fatalf("invalid variable should be skipped: %#v", values)
	}
}

func TestFetchJSVarsRejectsHTTPStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer server.Close()

	_, err := FetchJSVars(context.Background(), server.Client(), server.URL, []string{"code"}, FetchJSVarsOptions{
		Timeout: time.Second,
	})
	if err == nil {
		t.Fatal("expected HTTP status error")
	}
	if statusErr, ok := err.(HTTPStatusError); !ok || statusErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("err = %#v, want HTTPStatusError 502", err)
	}
}
