package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

// JSONPCallbackMode controls how a callback name is sent to the upstream.
type JSONPCallbackMode string

const (
	JSONPCallbackModeQuery JSONPCallbackMode = "query"
	JSONPCallbackModePath  JSONPCallbackMode = "path"
)

// JSONPOptions configures a JSONP request.
type JSONPOptions struct {
	Timeout       time.Duration
	CallbackParam string
	CallbackMode  JSONPCallbackMode
}

// HTTPStatusError reports a non-2xx HTTP response status.
type HTTPStatusError struct {
	StatusCode int
	Status     string
	URL        string
}

func (e HTTPStatusError) Error() string {
	if e.URL == "" {
		return fmt.Sprintf("http status %d", e.StatusCode)
	}
	return fmt.Sprintf("http status %d from %s", e.StatusCode, e.URL)
}

var jsonpCallbackCounter uint64

// ExtractJSONP extracts and decodes the JSON payload from a JSONP response.
func ExtractJSONP(text string, target any) error {
	cleaned := strings.TrimSpace(text)
	if commentEnd := strings.Index(cleaned, "*/"); commentEnd >= 0 {
		cleaned = strings.TrimSpace(cleaned[commentEnd+2:])
	}
	open := strings.IndexByte(cleaned, '(')
	if open < 0 {
		return parseError("Invalid JSONP response: no opening parenthesis found", nil)
	}
	close := strings.LastIndexByte(cleaned, ')')
	if close <= open {
		return parseError("Invalid JSONP response: no closing parenthesis found", nil)
	}
	payload := strings.TrimSpace(cleaned[open+1 : close])
	if payload == "" {
		return parseError("Invalid JSONP response: empty payload", nil)
	}
	if err := json.Unmarshal([]byte(payload), target); err != nil {
		return parseError("Invalid JSONP response: payload is not valid JSON", err)
	}
	return nil
}

// JSONPRequest fetches a JSONP URL and decodes its JSON payload.
func JSONPRequest(ctx context.Context, client *http.Client, requestURL string, target any, options JSONPOptions) error {
	if client == nil {
		client = http.DefaultClient
	}
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	callbackParam := options.CallbackParam
	if callbackParam == "" {
		callbackParam = "callback"
	}
	callbackMode := options.CallbackMode
	if callbackMode == "" {
		callbackMode = JSONPCallbackModeQuery
	}

	callbackName := nextJSONPCallbackName()
	finalURL, err := buildJSONPURL(requestURL, callbackName, callbackParam, callbackMode)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, finalURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer drainAndClose(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return HTTPStatusError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			URL:        finalURL,
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ExtractJSONP(string(body), target)
}

func buildJSONPURL(requestURL string, callbackName string, callbackParam string, callbackMode JSONPCallbackMode) (string, error) {
	if callbackMode == JSONPCallbackModePath {
		return strings.Replace(requestURL, "{callback}", callbackName, 1), nil
	}
	parsed, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set(callbackParam, callbackName)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func nextJSONPCallbackName() string {
	counter := atomic.AddUint64(&jsonpCallbackCounter, 1) - 1
	return fmt.Sprintf("__stock_sdk_jsonp_%d_%d", time.Now().UnixMilli(), counter)
}

func parseError(message string, cause error) error {
	return NewCodedError("PARSE_ERROR", message, cause)
}
