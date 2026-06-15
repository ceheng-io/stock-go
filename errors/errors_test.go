package errors

import (
	"context"
	stderrors "errors"
	"net/http"
	"testing"

	stock "github.com/ceheng.io/stock-go"
)

func TestErrorsSubpackageConstructorsDelegateToRoot(t *testing.T) {
	err := NewHTTPError(http.StatusTooManyRequests, "Too Many Requests", "https://example.test", stock.ProviderTencent)
	if err.Code != CodeRateLimited || GetSdkErrorCode(err) != CodeRateLimited {
		t.Fatalf("HTTP error code = %s", err.Code)
	}
	if err.Status == nil || *err.Status != http.StatusTooManyRequests {
		t.Fatalf("status = %#v", err.Status)
	}

	argErr := NewInvalidArgumentError("bad argument", map[string]any{"field": "code"})
	if !IsSdkError(argErr) {
		t.Fatal("IsSdkError returned false")
	}
	var requestErr *RequestError = argErr
	if requestErr.Details["field"] != "code" {
		t.Fatalf("details = %#v", requestErr.Details)
	}
}

func TestAttachErrorMetadataWrapsPlainError(t *testing.T) {
	cause := stderrors.New("dial tcp failed")

	err := AttachErrorMetadata(cause, ErrorMetadata{
		Code:     CodeNetwork,
		Provider: stock.ProviderTencent,
		URL:      "https://qt.gtimg.cn/?q=sh600519",
		Details:  map[string]any{"retry": false},
	})

	if err.Code != CodeNetwork || err.Provider != stock.ProviderTencent || err.URL == "" {
		t.Fatalf("metadata = %+v", err)
	}
	if !stderrors.Is(err, cause) {
		t.Fatal("AttachErrorMetadata did not preserve the original cause")
	}
	if err.Details["retry"] != false {
		t.Fatalf("details = %#v", err.Details)
	}
}

func TestNormalizeRequestErrorMapsCommonRequestErrors(t *testing.T) {
	timeout := NormalizeRequestError(context.DeadlineExceeded, ErrorContext{
		Provider: stock.ProviderEastmoney,
		URL:      "https://push2.eastmoney.com/api",
		Timeout: 1000,
	})
	if timeout.Code != CodeTimeout || timeout.Provider != stock.ProviderEastmoney || timeout.Details["timeout"] != 1000 {
		t.Fatalf("timeout error = %+v", timeout)
	}

	plain := NormalizeRequestError(stderrors.New("connection reset"), ErrorContext{
		Provider: stock.ProviderSina,
		URL:      "https://stock.finance.sina.com.cn",
	})
	if plain.Code != CodeNetwork || plain.Provider != stock.ProviderSina || plain.URL == "" {
		t.Fatalf("plain error = %+v", plain)
	}
}
