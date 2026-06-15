package stock

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestSDKErrorWrapsCodeAndCause(t *testing.T) {
	cause := errors.New("network down")
	err := NewError(CodeNetwork, "request failed", cause)

	if err.Code != CodeNetwork {
		t.Fatalf("Code = %s, want %s", err.Code, CodeNetwork)
	}
	if err.Error() != "NETWORK_ERROR: request failed: network down" {
		t.Fatalf("Error() = %q", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Fatal("errors.Is did not match wrapped cause")
	}
	if GetErrorCode(err) != CodeNetwork {
		t.Fatalf("GetErrorCode = %s, want %s", GetErrorCode(err), CodeNetwork)
	}
}

func TestGetErrorCodeUnknownForPlainError(t *testing.T) {
	if got := GetErrorCode(errors.New("plain")); got != CodeUnknown {
		t.Fatalf("GetErrorCode = %s, want %s", got, CodeUnknown)
	}
}

func TestGetErrorCodeCircuitBreakerOpen(t *testing.T) {
	if got := GetErrorCode(ErrCircuitBreakerOpen); got != CodeCircuitOpen {
		t.Fatalf("GetErrorCode = %s, want %s", got, CodeCircuitOpen)
	}
}

func TestGetErrorCodeContextErrors(t *testing.T) {
	if got := GetErrorCode(context.Canceled); got != CodeAborted {
		t.Fatalf("GetErrorCode(context.Canceled) = %s, want %s", got, CodeAborted)
	}
	if got := GetErrorCode(context.DeadlineExceeded); got != CodeTimeout {
		t.Fatalf("GetErrorCode(context.DeadlineExceeded) = %s, want %s", got, CodeTimeout)
	}
}

func TestSDKErrorCarriesRequestMetadata(t *testing.T) {
	cause := errors.New("bad json")
	err := NewErrorWithMetadata(ErrorMetadata{
		Code:     CodeParse,
		Message:  "parse failed",
		Provider: ProviderEastmoney,
		URL:      "https://example.test/data",
		Status:   http.StatusOK,
		Details:  map[string]any{"field": "data"},
		Cause:    cause,
	})

	if err.Code != CodeParse || err.Provider != ProviderEastmoney || err.URL != "https://example.test/data" {
		t.Fatalf("metadata missing: %#v", err)
	}
	if err.Status == nil || *err.Status != http.StatusOK {
		t.Fatalf("Status = %#v, want %d", err.Status, http.StatusOK)
	}
	if err.Details["field"] != "data" {
		t.Fatalf("Details = %#v", err.Details)
	}
	inputDetails := map[string]any{"field": "old"}
	copied := NewInvalidArgumentError("copy", inputDetails)
	inputDetails["field"] = "new"
	if copied.Details["field"] != "old" {
		t.Fatalf("Details reused caller map: %#v", copied.Details)
	}
	if !errors.Is(err, cause) {
		t.Fatal("errors.Is did not match wrapped cause")
	}
	if !IsSDKError(err) {
		t.Fatal("IsSDKError returned false")
	}
}

func TestSpecificSDKErrorConstructors(t *testing.T) {
	httpErr := NewHTTPError(http.StatusTooManyRequests, "Too Many Requests", "https://example.test", ProviderTencent)
	if httpErr.Code != CodeRateLimited {
		t.Fatalf("HTTP 429 code = %s, want %s", httpErr.Code, CodeRateLimited)
	}
	if httpErr.Status == nil || *httpErr.Status != http.StatusTooManyRequests {
		t.Fatalf("HTTP status = %#v", httpErr.Status)
	}
	if httpErr.StatusText != "Too Many Requests" {
		t.Fatalf("StatusText = %q", httpErr.StatusText)
	}

	emptyErr := NewUpstreamEmptyError("empty", ProviderEastmoney, "https://example.test/empty")
	if emptyErr.Code != CodeUpstreamEmpty || emptyErr.Provider != ProviderEastmoney {
		t.Fatalf("UpstreamEmpty error = %#v", emptyErr)
	}

	notFoundErr := NewNotFoundError("missing", ProviderSina, "https://example.test/missing")
	if notFoundErr.Code != CodeNotFound || notFoundErr.Provider != ProviderSina {
		t.Fatalf("NotFound error = %#v", notFoundErr)
	}

	argErr := NewInvalidArgumentError("invalid page", map[string]any{"page": 0})
	if argErr.Code != CodeInvalidArgument || argErr.Details["page"] != 0 {
		t.Fatalf("InvalidArgument error = %#v", argErr)
	}

	symbolErr := NewInvalidSymbolError("bad-symbol", ProviderTencent)
	if symbolErr.Code != CodeInvalidSymbol || symbolErr.Details["symbol"] != "bad-symbol" {
		t.Fatalf("InvalidSymbol error = %#v", symbolErr)
	}

	upstreamErr := NewUpstreamError("bad upstream", ProviderEastmoney, "https://example.test/up", map[string]any{"code": 1})
	if upstreamErr.Code != CodeUpstream || upstreamErr.Details["code"] != 1 {
		t.Fatalf("Upstream error = %#v", upstreamErr)
	}

	abortedErr := NewAbortedError("", ProviderTencent, "https://example.test/abort")
	if abortedErr.Code != CodeAborted || abortedErr.Message != "Request aborted" {
		t.Fatalf("Aborted error = %#v", abortedErr)
	}
}

func TestGetErrorCodeRecognizesStandardWrappedErrors(t *testing.T) {
	err := NewHTTPError(http.StatusInternalServerError, "Internal Server Error", "", ProviderTencent)
	if got := GetErrorCode(err); got != CodeHTTP {
		t.Fatalf("GetErrorCode(HTTP 500) = %s, want %s", got, CodeHTTP)
	}
	wrapped := errors.Join(errors.New("outer"), NewInvalidArgumentError("bad", nil))
	if got := GetErrorCode(wrapped); got != CodeInvalidArgument {
		t.Fatalf("GetErrorCode(joined) = %s, want %s", got, CodeInvalidArgument)
	}
}

func TestRootReExportsTSCoreErrorNames(t *testing.T) {
	err := NewError(CodeNetwork, "request failed", errors.New("network down"))

	var sdkErr *SdkError = err
	var requestErr *RequestError = err
	var httpErr *HttpError = NewHTTPError(http.StatusBadGateway, "Bad Gateway", "https://example.test", ProviderTencent)
	var emptyErr *UpstreamEmptyError = NewUpstreamEmptyError("empty", ProviderEastmoney, "")
	var notFoundErr *NotFoundError = NewNotFoundError("missing", ProviderSina, "")
	var invalidArgErr *InvalidArgumentError = NewInvalidArgumentError("bad", nil)
	var invalidSymbolErr *InvalidSymbolError = NewInvalidSymbolError("bad-symbol", ProviderTencent)
	var upstreamErr *UpstreamError = NewUpstreamError("upstream", ProviderEastmoney, "", nil)
	var abortedErr *AbortedError = NewAbortedError("", ProviderTencent, "")
	var code SdkErrorCode = sdkErr.Code
	if code != CodeNetwork {
		t.Fatalf("SdkErrorCode = %s, want %s", code, CodeNetwork)
	}
	for name, current := range map[string]*Error{
		"HttpError":            httpErr,
		"UpstreamEmptyError":   emptyErr,
		"NotFoundError":        notFoundErr,
		"InvalidArgumentError": invalidArgErr,
		"InvalidSymbolError":   invalidSymbolErr,
		"UpstreamError":        upstreamErr,
		"AbortedError":         abortedErr,
	} {
		if current == nil {
			t.Fatalf("%s is nil", name)
		}
	}

	wrapped := errors.Join(errors.New("outer"), requestErr)
	if !IsSdkError(wrapped) {
		t.Fatal("IsSdkError returned false")
	}
	if got := GetSdkErrorCode(wrapped); got != CodeNetwork {
		t.Fatalf("GetSdkErrorCode = %s, want %s", got, CodeNetwork)
	}
	if got := GetSDKErrorCode(wrapped); got != CodeNetwork {
		t.Fatalf("GetSDKErrorCode = %s, want %s", got, CodeNetwork)
	}
}
