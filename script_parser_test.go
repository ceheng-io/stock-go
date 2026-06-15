package stock

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRootExtractJSONP(t *testing.T) {
	var payload struct {
		Code string `json:"fundcode"`
		Name string `json:"name"`
	}

	err := ExtractJSONP(`/*<script></script>*/ callback({"fundcode":"110011","name":"易方达中小盘"});`, &payload)
	if err != nil {
		t.Fatal(err)
	}
	if payload.Code != "110011" || payload.Name != "易方达中小盘" {
		t.Fatalf("payload = %+v", payload)
	}
}

func TestRootExtractJSONPInvalidResponseReturnsParseError(t *testing.T) {
	var payload map[string]any
	tests := []string{
		`callback`,
		`callback({"code":"110011"}`,
		`callback()`,
		`callback({bad json})`,
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			err := ExtractJsonFromJsonp(input, &payload)
			if err == nil {
				t.Fatal("expected parse error")
			}
			if code := GetErrorCode(err); code != CodeParse {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeParse, err)
			}
		})
	}
}

func TestRootReExportsTSJSONPAndJSVarsNames(t *testing.T) {
	var payload struct {
		Code string `json:"code"`
	}

	if err := ExtractJsonFromJsonp(`callback({"code":"110011"})`, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Code != "110011" {
		t.Fatalf("payload.Code = %q, want 110011", payload.Code)
	}
	if err := ExtractJsonFromJsonp(`callback({"code":"110011"})`, &payload); err != nil {
		t.Fatal(err)
	}

	var jsonpOptions JsonpOptions = JSONPOptions{
		Timeout:       time.Second,
		CallbackParam: "cb",
		CallbackMode:  JSONPCallbackModeQuery,
	}
	if jsonpOptions.CallbackParam != "cb" {
		t.Fatalf("JsonpOptions.CallbackParam = %q, want cb", jsonpOptions.CallbackParam)
	}

	var fetchOptions FetchJsVarsOptions = FetchJSVarsOptions{
		Timeout: time.Second,
		Headers: map[string]string{"X-Test": "yes"},
	}
	if fetchOptions.Headers["X-Test"] != "yes" {
		t.Fatalf("FetchJsVarsOptions.Headers = %#v", fetchOptions.Headers)
	}
}

func TestRootReExportsTSLowercaseJSHelpers(t *testing.T) {
	text := `
var code = "110011";
const rows = [{"name":"易方达中小盘","value":3.5}];
`

	values := ParseJsVars(text, "code", "rows")
	if values["code"] != "110011" {
		t.Fatalf("ParseJsVars code = %#v, want 110011", values["code"])
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(text))
	}))
	defer server.Close()

	fetched, err := FetchJsVars(context.Background(), server.URL, []string{"code"},
		WithFetchJSVarsHTTPClient(server.Client()),
		WithFetchJSVarsTimeout(time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	if fetched["code"] != "110011" {
		t.Fatalf("fetchJsVars code = %#v, want 110011", fetched["code"])
	}
	if BROWSER_JSVARS_MUTEX_KEY == "" {
		t.Fatal("BROWSER_JSVARS_MUTEX_KEY is empty")
	}
}

func TestRootExtractJSVar(t *testing.T) {
	text := `
var fS_code = "110011";
const Data_netWorthTrend = [{"x":1702857600000,"y":3.5,"unitMoney":""}];
`
	var code string
	if err := ExtractJSVar(text, "fS_code", &code); err != nil {
		t.Fatal(err)
	}
	if code != "110011" {
		t.Fatalf("code = %q, want 110011", code)
	}

	var trend []map[string]any
	if err := ExtractJSVar(text, "Data_netWorthTrend", &trend); err != nil {
		t.Fatal(err)
	}
	if len(trend) != 1 || trend[0]["y"] != 3.5 {
		t.Fatalf("trend = %+v", trend)
	}
}

func TestRootParseJSVars(t *testing.T) {
	text := `
var fS_code = "110011";
const Data_netWorthTrend = [{"x":1702857600000,"y":3.5,"unitMoney":""}];
let invalid = {not_json: true};
`

	values := ParseJSVars(text, "fS_code", "Data_netWorthTrend", "missing", "invalid")
	if len(values) != 2 {
		t.Fatalf("len(values) = %d, want 2: %#v", len(values), values)
	}
	if values["fS_code"] != "110011" {
		t.Fatalf("fS_code = %#v, want 110011", values["fS_code"])
	}
	trend, ok := values["Data_netWorthTrend"].([]any)
	if !ok || len(trend) != 1 {
		t.Fatalf("Data_netWorthTrend = %#v, want one row slice", values["Data_netWorthTrend"])
	}
}

func TestRootJSONPRequest(t *testing.T) {
	var callback string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callback = r.URL.Query().Get("cb")
		if callback == "" {
			t.Fatal("missing callback query parameter")
		}
		_, _ = w.Write([]byte(callback + `({"code":"110011"})`))
	}))
	defer server.Close()

	var payload struct {
		Code string `json:"code"`
	}
	err := JSONPRequest(context.Background(), server.URL+"/api", &payload,
		WithJSONPHTTPClient(server.Client()),
		WithJSONPCallbackParam("cb"),
		WithJSONPTimeout(time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	if callback == "" || payload.Code != "110011" {
		t.Fatalf("callback=%q payload=%+v", callback, payload)
	}
}

func TestRootJsonpRequestTSNamingStyle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callback := r.URL.Query().Get("callback")
		if callback == "" {
			t.Fatal("missing callback query parameter")
		}
		_, _ = w.Write([]byte(callback + `({"code":"110011"})`))
	}))
	defer server.Close()

	var payload struct {
		Code string `json:"code"`
	}
	err := JsonpRequest(context.Background(), server.URL+"/api", &payload,
		WithJSONPHTTPClient(server.Client()),
		WithJSONPTimeout(time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	if payload.Code != "110011" {
		t.Fatalf("payload.Code = %q, want 110011", payload.Code)
	}
}

func TestRootJSONPRequestHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer server.Close()

	var payload map[string]any
	err := JSONPRequest(context.Background(), server.URL, &payload,
		WithJSONPHTTPClient(server.Client()),
		WithJSONPTimeout(time.Second),
	)
	if err == nil {
		t.Fatal("expected HTTP error")
	}
	var sdkErr *Error
	if !errors.As(err, &sdkErr) {
		t.Fatalf("err = %#v, want *Error", err)
	}
	if sdkErr.Code != CodeHTTP || sdkErr.Status == nil || *sdkErr.Status != http.StatusBadGateway {
		t.Fatalf("sdk error = %+v", sdkErr)
	}
}

func TestRootFetchJSVars(t *testing.T) {
	var gotHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Test")
		_, _ = w.Write([]byte(`
var code = "110011";
const rows = [{"name":"易方达中小盘","value":3.5}];
`))
	}))
	defer server.Close()

	values, err := FetchJSVars(context.Background(), server.URL, []string{"code", "rows"},
		WithFetchJSVarsHTTPClient(server.Client()),
		WithFetchJSVarsTimeout(time.Second),
		WithFetchJSVarsHeaders(map[string]string{"X-Test": "yes"}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if gotHeader != "yes" {
		t.Fatalf("X-Test header = %q, want yes", gotHeader)
	}
	if values["code"] != "110011" {
		t.Fatalf("code = %#v, want 110011", values["code"])
	}
	if rows, ok := values["rows"].([]any); !ok || len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", values["rows"])
	}
}

func TestRootFetchJSVarsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer server.Close()

	_, err := FetchJSVars(context.Background(), server.URL, []string{"code"},
		WithFetchJSVarsHTTPClient(server.Client()),
		WithFetchJSVarsTimeout(time.Second),
	)
	if err == nil {
		t.Fatal("expected HTTP error")
	}
	var sdkErr *Error
	if !errors.As(err, &sdkErr) {
		t.Fatalf("err = %#v, want *Error", err)
	}
	if sdkErr.Code != CodeHTTP || sdkErr.Status == nil || *sdkErr.Status != http.StatusBadGateway {
		t.Fatalf("sdk error = %+v", sdkErr)
	}
}
