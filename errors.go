package stock

import (
	"context"
	"errors"
	"fmt"

	"github.com/ceheng.io/stock-go/internal/core"
)

// ErrorCode is a stable SDK error code.
type ErrorCode string

const (
	CodeUnknown         ErrorCode = "UNKNOWN"
	CodeInvalidSymbol   ErrorCode = "INVALID_SYMBOL"
	CodeInvalidArgument ErrorCode = "INVALID_ARGUMENT"
	CodeNetwork         ErrorCode = "NETWORK_ERROR"
	CodeHTTP            ErrorCode = "HTTP_ERROR"
	CodeRateLimited     ErrorCode = "RATE_LIMITED"
	CodeTimeout         ErrorCode = "TIMEOUT"
	CodeAborted         ErrorCode = "ABORTED"
	CodeParse           ErrorCode = "PARSE_ERROR"
	CodeCircuitOpen     ErrorCode = "CIRCUIT_OPEN"
	CodeUpstreamEmpty   ErrorCode = "UPSTREAM_EMPTY"
	CodeUpstream        ErrorCode = "UPSTREAM_ERROR"
	CodeNotFound        ErrorCode = "NOT_FOUND"
)

// ErrCircuitBreakerOpen is returned when provider requests are rejected by an open circuit.
var ErrCircuitBreakerOpen = core.ErrCircuitBreakerOpen

// Error is the public SDK error type.
type Error struct {
	Code       ErrorCode
	Message    string
	Provider   ProviderName
	URL        string
	Status     *int
	StatusText string
	Details    map[string]any
	Cause      error
}

// SdkError preserves the TypeScript SDK error type name at the Go root package.
type SdkError = Error

// HttpError preserves the TypeScript SDK HTTP error type name.
type HttpError = Error

// UpstreamEmptyError preserves the TypeScript SDK upstream-empty error type name.
type UpstreamEmptyError = Error

// NotFoundError preserves the TypeScript SDK not-found error type name.
type NotFoundError = Error

// InvalidArgumentError preserves the TypeScript SDK invalid-argument error type name.
type InvalidArgumentError = Error

// InvalidSymbolError preserves the TypeScript SDK invalid-symbol error type name.
type InvalidSymbolError = Error

// UpstreamError preserves the TypeScript SDK upstream error type name.
type UpstreamError = Error

// AbortedError preserves the TypeScript SDK aborted-request error type name.
type AbortedError = Error

// RequestError is a compatibility alias for SDK request errors.
type RequestError = Error

// SdkErrorCode preserves the TypeScript SDK error-code type name.
type SdkErrorCode = ErrorCode

// ErrorMetadata contains structured SDK error context.
type ErrorMetadata struct {
	Code       ErrorCode
	Message    string
	Provider   ProviderName
	URL        string
	Status     int
	StatusText string
	Details    map[string]any
	Cause      error
}

// NewError creates an SDK error with a stable code and optional cause.
func NewError(code ErrorCode, message string, cause error) *Error {
	return NewErrorWithMetadata(ErrorMetadata{
		Code:    code,
		Message: message,
		Cause:   cause,
	})
}

// NewErrorWithMetadata creates an SDK error with structured request context.
func NewErrorWithMetadata(metadata ErrorMetadata) *Error {
	code := metadata.Code
	if code == "" {
		code = CodeUnknown
	}
	details := cloneErrorDetails(metadata.Details)
	var status *int
	if metadata.Status != 0 {
		value := metadata.Status
		status = &value
	}
	return &Error{
		Code:       code,
		Message:    metadata.Message,
		Provider:   metadata.Provider,
		URL:        metadata.URL,
		Status:     status,
		StatusText: metadata.StatusText,
		Details:    details,
		Cause:      metadata.Cause,
	}
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause for errors.Is and errors.As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// NewHTTPError creates an HTTP status error.
func NewHTTPError(status int, statusText string, requestURL string, provider ProviderName) *Error {
	code := CodeHTTP
	if status == 429 {
		code = CodeRateLimited
	}
	details := map[string]any{"statusText": statusText}
	suffix := ""
	if statusText != "" {
		suffix += " " + statusText
	}
	if requestURL != "" {
		suffix += ", url: " + requestURL
	}
	if provider != "" {
		suffix += ", provider: " + string(provider)
	}
	return NewErrorWithMetadata(ErrorMetadata{
		Code:       code,
		Message:    fmt.Sprintf("HTTP error! status: %d%s", status, suffix),
		Provider:   provider,
		URL:        requestURL,
		Status:     status,
		StatusText: statusText,
		Details:    details,
	})
}

// NewUpstreamEmptyError creates an upstream-empty-data error.
func NewUpstreamEmptyError(message string, provider ProviderName, requestURL string) *Error {
	return NewErrorWithMetadata(ErrorMetadata{
		Code:     CodeUpstreamEmpty,
		Message:  message,
		Provider: provider,
		URL:      requestURL,
	})
}

// NewNotFoundError creates a resource-not-found error.
func NewNotFoundError(message string, provider ProviderName, requestURL string) *Error {
	return NewErrorWithMetadata(ErrorMetadata{
		Code:     CodeNotFound,
		Message:  message,
		Provider: provider,
		URL:      requestURL,
	})
}

// NewInvalidArgumentError creates an invalid-argument error.
func NewInvalidArgumentError(message string, details map[string]any) *Error {
	return NewErrorWithMetadata(ErrorMetadata{
		Code:    CodeInvalidArgument,
		Message: message,
		Details: details,
	})
}

// NewInvalidSymbolError creates an invalid-symbol error.
func NewInvalidSymbolError(symbol string, provider ProviderName) *Error {
	return NewErrorWithMetadata(ErrorMetadata{
		Code:     CodeInvalidSymbol,
		Message:  fmt.Sprintf("Invalid symbol: %s", symbol),
		Provider: provider,
		Details:  map[string]any{"symbol": symbol},
	})
}

// NewUpstreamError creates a structured upstream error.
func NewUpstreamError(message string, provider ProviderName, requestURL string, details map[string]any) *Error {
	return NewErrorWithMetadata(ErrorMetadata{
		Code:     CodeUpstream,
		Message:  message,
		Provider: provider,
		URL:      requestURL,
		Details:  details,
	})
}

// NewAbortedError creates an aborted-request error.
func NewAbortedError(message string, provider ProviderName, requestURL string) *Error {
	if message == "" {
		message = "Request aborted"
	}
	return NewErrorWithMetadata(ErrorMetadata{
		Code:     CodeAborted,
		Message:  message,
		Provider: provider,
		URL:      requestURL,
	})
}

// IsSDKError reports whether err wraps an SDK Error.
func IsSDKError(err error) bool {
	var sdkErr *Error
	return errors.As(err, &sdkErr)
}

// IsSdkError reports whether err wraps an SDK Error.
func IsSdkError(err error) bool {
	return IsSDKError(err)
}

// GetErrorCode returns an SDK error code or CodeUnknown for non-SDK errors.
func GetErrorCode(err error) ErrorCode {
	var sdkErr *Error
	if errors.As(err, &sdkErr) {
		return sdkErr.Code
	}
	var codedErr core.CodedError
	if errors.As(err, &codedErr) {
		return ErrorCode(codedErr.SDKCode())
	}
	if errors.Is(err, ErrCircuitBreakerOpen) {
		return CodeCircuitOpen
	}
	if errors.Is(err, context.Canceled) {
		return CodeAborted
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return CodeTimeout
	}
	return CodeUnknown
}

// GetSdkErrorCode returns an SDK error code using the TypeScript SDK naming.
func GetSdkErrorCode(err error) ErrorCode {
	return GetErrorCode(err)
}

// GetSDKErrorCode returns an SDK error code using Go initialism style.
func GetSDKErrorCode(err error) ErrorCode {
	return GetErrorCode(err)
}

func cloneErrorDetails(details map[string]any) map[string]any {
	if details == nil {
		return nil
	}
	clone := make(map[string]any, len(details))
	for key, value := range details {
		clone[key] = value
	}
	return clone
}
