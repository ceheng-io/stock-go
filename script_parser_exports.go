package stock

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ceheng.io/stock-go/internal/core"
)

// JSONPCallbackMode controls how a callback name is sent to the upstream.
type JSONPCallbackMode = core.JSONPCallbackMode

// JSONPOptions configures a JSONP request.
type JSONPOptions = core.JSONPOptions

// JsonpOptions preserves the TypeScript SDK option type name at the Go root package.
type JsonpOptions = core.JSONPOptions

// FetchJSVarsOptions configures FetchJSVars.
type FetchJSVarsOptions = core.FetchJSVarsOptions

// FetchJsVarsOptions preserves the TypeScript SDK option type name at the Go root package.
type FetchJsVarsOptions = core.FetchJSVarsOptions

const (
	JSONPCallbackModeQuery JSONPCallbackMode = core.JSONPCallbackModeQuery
	JSONPCallbackModePath  JSONPCallbackMode = core.JSONPCallbackModePath

	BROWSER_JSVARS_MUTEX_KEY = "jsVars"
)

type jsonpRequestConfig struct {
	httpClient    *http.Client
	timeout       time.Duration
	callbackParam string
	callbackMode  JSONPCallbackMode
}

type fetchJSVarsConfig struct {
	httpClient *http.Client
	timeout    time.Duration
	headers    map[string]string
}

// JSONPOption configures JSONPRequest.
type JSONPOption func(*jsonpRequestConfig)

// FetchJSVarsOption configures FetchJSVars.
type FetchJSVarsOption func(*fetchJSVarsConfig)

// WithJSONPHTTPClient sets the HTTP client used by JSONPRequest.
func WithJSONPHTTPClient(client *http.Client) JSONPOption {
	return func(config *jsonpRequestConfig) {
		config.httpClient = client
	}
}

// WithJSONPTimeout sets the request timeout used by JSONPRequest.
func WithJSONPTimeout(timeout time.Duration) JSONPOption {
	return func(config *jsonpRequestConfig) {
		config.timeout = timeout
	}
}

// WithJSONPCallbackParam sets the callback query parameter name.
func WithJSONPCallbackParam(param string) JSONPOption {
	return func(config *jsonpRequestConfig) {
		config.callbackParam = param
	}
}

// WithJSONPCallbackMode sets whether the callback is sent by query or path placeholder.
func WithJSONPCallbackMode(mode JSONPCallbackMode) JSONPOption {
	return func(config *jsonpRequestConfig) {
		config.callbackMode = mode
	}
}

// WithFetchJSVarsHTTPClient sets the HTTP client used by FetchJSVars.
func WithFetchJSVarsHTTPClient(client *http.Client) FetchJSVarsOption {
	return func(config *fetchJSVarsConfig) {
		config.httpClient = client
	}
}

// WithFetchJSVarsTimeout sets the request timeout used by FetchJSVars.
func WithFetchJSVarsTimeout(timeout time.Duration) FetchJSVarsOption {
	return func(config *fetchJSVarsConfig) {
		config.timeout = timeout
	}
}

// WithFetchJSVarsHeaders sets request headers used by FetchJSVars.
func WithFetchJSVarsHeaders(headers map[string]string) FetchJSVarsOption {
	return func(config *fetchJSVarsConfig) {
		config.headers = headers
	}
}

// ExtractJSONP extracts and decodes the JSON payload from a JSONP response.
func ExtractJSONP(text string, target any) error {
	return core.ExtractJSONP(text, target)
}

// ExtractJsonFromJsonp extracts and decodes JSONP using the TypeScript SDK naming.
func ExtractJsonFromJsonp(text string, target any) error {
	return ExtractJSONP(text, target)
}

// JSONPRequest fetches a JSONP URL and decodes its JSON payload.
func JSONPRequest(ctx context.Context, requestURL string, target any, options ...JSONPOption) error {
	config := jsonpRequestConfig{}
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}
	err := core.JSONPRequest(ctx, config.httpClient, requestURL, target, core.JSONPOptions{
		Timeout:       config.timeout,
		CallbackParam: config.callbackParam,
		CallbackMode:  core.JSONPCallbackMode(config.callbackMode),
	})
	if err == nil {
		return nil
	}
	var statusErr core.HTTPStatusError
	if errors.As(err, &statusErr) {
		return NewHTTPError(statusErr.StatusCode, statusErr.Status, statusErr.URL, InferProviderFromURL(statusErr.URL, ""))
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return NewErrorWithMetadata(ErrorMetadata{
			Code:    CodeTimeout,
			Message: "JSONP request timed out",
			URL:     requestURL,
			Cause:   err,
		})
	}
	if errors.Is(err, context.Canceled) {
		return NewAbortedError("JSONP request aborted", InferProviderFromURL(requestURL, ""), requestURL)
	}
	return err
}

// JsonpRequest fetches a JSONP URL using the TypeScript SDK naming style.
func JsonpRequest(ctx context.Context, requestURL string, target any, options ...JSONPOption) error {
	return JSONPRequest(ctx, requestURL, target, options...)
}

// ExtractJSVar extracts and decodes one JSON-compatible JavaScript variable declaration.
func ExtractJSVar(text string, name string, target any) error {
	return core.ExtractJSVar(text, name, target)
}

// ParseJSVars extracts JSON-compatible JavaScript variable declarations.
func ParseJSVars(text string, names ...string) map[string]any {
	return core.ParseJSVars(text, names...)
}

// ParseJsVars extracts JavaScript variable declarations using the TypeScript SDK naming style.
func ParseJsVars(text string, names ...string) map[string]any {
	return ParseJSVars(text, names...)
}

// FetchJSVars fetches a JavaScript variable declaration document and extracts variables.
func FetchJSVars(ctx context.Context, requestURL string, names []string, options ...FetchJSVarsOption) (map[string]any, error) {
	config := fetchJSVarsConfig{}
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}
	values, err := core.FetchJSVars(ctx, config.httpClient, requestURL, names, core.FetchJSVarsOptions{
		Timeout: config.timeout,
		Headers: cloneStringMap(config.headers),
	})
	if err == nil {
		return values, nil
	}
	var statusErr core.HTTPStatusError
	if errors.As(err, &statusErr) {
		return nil, NewHTTPError(statusErr.StatusCode, statusErr.Status, statusErr.URL, InferProviderFromURL(statusErr.URL, ""))
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, NewErrorWithMetadata(ErrorMetadata{
			Code:    CodeTimeout,
			Message: "FetchJSVars request timed out",
			URL:     requestURL,
			Cause:   err,
		})
	}
	if errors.Is(err, context.Canceled) {
		return nil, NewAbortedError("FetchJSVars request aborted", InferProviderFromURL(requestURL, ""), requestURL)
	}
	return nil, err
}

// FetchJsVars fetches JavaScript variable declarations using the TypeScript SDK naming style.
func FetchJsVars(ctx context.Context, requestURL string, names []string, options ...FetchJSVarsOption) (map[string]any, error) {
	return FetchJSVars(ctx, requestURL, names, options...)
}
