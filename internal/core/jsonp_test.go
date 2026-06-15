package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExtractJSONPParsesCallbackPayload(t *testing.T) {
	var payload struct {
		Code string `json:"fundcode"`
		Name string `json:"name"`
	}

	if err := ExtractJSONP(`jsonpgz({"fundcode":"110011","name":"易方达中小盘"});`, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Code != "110011" || payload.Name != "易方达中小盘" {
		t.Fatalf("payload = %+v", payload)
	}
}

func TestExtractJSONPRejectsInvalidPayload(t *testing.T) {
	var payload map[string]any
	if err := ExtractJSONP(`jsonpgz("broken";`, &payload); err == nil {
		t.Fatal("expected invalid JSONP error")
	}
}

func TestJSONPRequestQueryMode(t *testing.T) {
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
	err := JSONPRequest(context.Background(), server.Client(), server.URL+"/api?existing=1", &payload, JSONPOptions{
		CallbackParam: "cb",
		Timeout:       time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if callback == "" || payload.Code != "110011" {
		t.Fatalf("callback=%q payload=%+v", callback, payload)
	}
}

func TestJSONPRequestPathMode(t *testing.T) {
	var requestedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		if requestedPath == "/jsonp/{callback}/data" {
			t.Fatal("callback placeholder was not replaced")
		}
		callback := requestedPath[len("/jsonp/") : len(requestedPath)-len("/data")]
		_, _ = w.Write([]byte(callback + `([{"code":"IF2412"}])`))
	}))
	defer server.Close()

	var payload []struct {
		Code string `json:"code"`
	}
	err := JSONPRequest(context.Background(), server.Client(), server.URL+"/jsonp/{callback}/data", &payload, JSONPOptions{
		CallbackMode: JSONPCallbackModePath,
		Timeout:      time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if requestedPath == "" || len(payload) != 1 || payload[0].Code != "IF2412" {
		t.Fatalf("path=%q payload=%+v", requestedPath, payload)
	}
}

func TestJSONPRequestRejectsHTTPStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer server.Close()

	var payload map[string]any
	err := JSONPRequest(context.Background(), server.Client(), server.URL, &payload, JSONPOptions{
		Timeout: time.Second,
	})
	if err == nil {
		t.Fatal("expected HTTP status error")
	}
	if statusErr, ok := err.(HTTPStatusError); !ok || statusErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("err = %#v, want HTTPStatusError 502", err)
	}
}

func TestJSONPRequestDrainsHTTPErrorResponseBeforeClose(t *testing.T) {
	body := &drainTrackingBody{content: []byte("bad gateway")}
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Status:     http.StatusText(http.StatusBadGateway),
			Body:       body,
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})}

	var payload map[string]any
	err := JSONPRequest(context.Background(), client, "https://jsonp-error.test/api", &payload, JSONPOptions{
		Timeout: time.Second,
	})
	if err == nil {
		t.Fatal("expected HTTP status error")
	}
	if !body.closed {
		t.Fatal("response body was not closed")
	}
	if !body.closedAfterEOF {
		t.Fatal("response body closed before it was drained to EOF")
	}
}
