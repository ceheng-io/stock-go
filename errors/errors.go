// Package errors exposes SDK error types through the TypeScript-compatible subpath.
package errors

import (
	"context"
	stderrors "errors"

	stock "github.com/ceheng-io/stock-go"
)

type ErrorCode = stock.ErrorCode

const (
	CodeUnknown         ErrorCode = stock.CodeUnknown
	CodeInvalidSymbol   ErrorCode = stock.CodeInvalidSymbol
	CodeInvalidArgument ErrorCode = stock.CodeInvalidArgument
	CodeNetwork         ErrorCode = stock.CodeNetwork
	CodeHTTP            ErrorCode = stock.CodeHTTP
	CodeRateLimited     ErrorCode = stock.CodeRateLimited
	CodeTimeout         ErrorCode = stock.CodeTimeout
	CodeAborted         ErrorCode = stock.CodeAborted
	CodeParse           ErrorCode = stock.CodeParse
	CodeCircuitOpen     ErrorCode = stock.CodeCircuitOpen
	CodeUpstreamEmpty   ErrorCode = stock.CodeUpstreamEmpty
	CodeUpstream        ErrorCode = stock.CodeUpstream
	CodeNotFound        ErrorCode = stock.CodeNotFound
)

var ErrCircuitBreakerOpen = stock.ErrCircuitBreakerOpen

type Error = stock.Error
type SdkError = stock.SdkError
type HttpError = stock.HttpError
type UpstreamEmptyError = stock.UpstreamEmptyError
type NotFoundError = stock.NotFoundError
type InvalidArgumentError = stock.InvalidArgumentError
type InvalidSymbolError = stock.InvalidSymbolError
type UpstreamError = stock.UpstreamError
type AbortedError = stock.AbortedError
type RequestError = stock.RequestError
type SdkErrorCode = stock.SdkErrorCode
type ErrorMetadata = stock.ErrorMetadata

// ErrorContext carries request context for NormalizeRequestError.
type ErrorContext struct {
	Provider stock.ProviderName
	URL      string
	Timeout  int
}

func NewError(code ErrorCode, message string, cause error) *Error {
	return stock.NewError(code, message, cause)
}

func NewErrorWithMetadata(metadata ErrorMetadata) *Error {
	return stock.NewErrorWithMetadata(metadata)
}

func NewHTTPError(status int, statusText string, requestURL string, provider stock.ProviderName) *Error {
	return stock.NewHTTPError(status, statusText, requestURL, provider)
}

func NewUpstreamEmptyError(message string, provider stock.ProviderName, requestURL string) *Error {
	return stock.NewUpstreamEmptyError(message, provider, requestURL)
}

func NewNotFoundError(message string, provider stock.ProviderName, requestURL string) *Error {
	return stock.NewNotFoundError(message, provider, requestURL)
}

func NewInvalidArgumentError(message string, details map[string]any) *Error {
	return stock.NewInvalidArgumentError(message, details)
}

func NewInvalidSymbolError(symbol string, provider stock.ProviderName) *Error {
	return stock.NewInvalidSymbolError(symbol, provider)
}

func NewUpstreamError(message string, provider stock.ProviderName, requestURL string, details map[string]any) *Error {
	return stock.NewUpstreamError(message, provider, requestURL, details)
}

func NewAbortedError(message string, provider stock.ProviderName, requestURL string) *Error {
	return stock.NewAbortedError(message, provider, requestURL)
}

func GetErrorCode(err error) ErrorCode {
	return stock.GetErrorCode(err)
}

func GetSDKErrorCode(err error) ErrorCode {
	return stock.GetSDKErrorCode(err)
}

func GetSdkErrorCode(err error) ErrorCode {
	return stock.GetSdkErrorCode(err)
}

func IsSDKError(err error) bool {
	return stock.IsSDKError(err)
}

func IsSdkError(err error) bool {
	return stock.IsSdkError(err)
}

// AttachErrorMetadata wraps an existing error with SDK request metadata.
func AttachErrorMetadata(err error, metadata ErrorMetadata) *Error {
	if err == nil {
		return stock.NewErrorWithMetadata(metadata)
	}
	var sdkErr *stock.Error
	if stderrors.As(err, &sdkErr) {
		return mergeErrorMetadata(sdkErr, metadata)
	}
	if metadata.Cause == nil {
		metadata.Cause = err
	}
	if metadata.Message == "" {
		metadata.Message = err.Error()
	}
	return stock.NewErrorWithMetadata(metadata)
}

// NormalizeRequestError maps unknown request errors to SDK errors.
func NormalizeRequestError(err error, requestContext ErrorContext) *Error {
	if err == nil {
		return nil
	}
	var sdkErr *stock.Error
	if stderrors.As(err, &sdkErr) {
		return mergeErrorMetadata(sdkErr, ErrorMetadata{
			Provider: requestContext.Provider,
			URL:      requestContext.URL,
		})
	}
	metadata := ErrorMetadata{
		Code:     CodeNetwork,
		Message:  err.Error(),
		Provider: requestContext.Provider,
		URL:      requestContext.URL,
		Cause:    err,
	}
	if stderrors.Is(err, context.DeadlineExceeded) {
		metadata.Code = CodeTimeout
		if requestContext.Timeout > 0 {
			metadata.Details = map[string]any{"timeout": requestContext.Timeout}
		}
	} else if stderrors.Is(err, context.Canceled) {
		metadata.Code = CodeAborted
	}
	return stock.NewErrorWithMetadata(metadata)
}

func mergeErrorMetadata(err *stock.Error, metadata ErrorMetadata) *Error {
	if err == nil {
		return stock.NewErrorWithMetadata(metadata)
	}
	if metadata.Provider != "" {
		err.Provider = metadata.Provider
	}
	if metadata.URL != "" {
		err.URL = metadata.URL
	}
	if err.Status == nil && metadata.Status != 0 {
		status := metadata.Status
		err.Status = &status
	}
	if err.StatusText == "" {
		err.StatusText = metadata.StatusText
	}
	if err.Details == nil && metadata.Details != nil {
		err.Details = cloneDetails(metadata.Details)
	}
	return err
}

func cloneDetails(details map[string]any) map[string]any {
	if details == nil {
		return nil
	}
	clone := make(map[string]any, len(details))
	for key, value := range details {
		clone[key] = value
	}
	return clone
}
