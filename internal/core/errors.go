package core

import "fmt"

// CodedError carries a stable SDK error code without depending on the root package.
type CodedError interface {
	error
	SDKCode() string
}

type codedError struct {
	code    string
	message string
	cause   error
}

// NewCodedError creates an internal error with a stable SDK error code.
func NewCodedError(code string, message string, cause error) error {
	return codedError{code: code, message: message, cause: cause}
}

func (e codedError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e codedError) Unwrap() error {
	return e.cause
}

func (e codedError) SDKCode() string {
	return e.code
}
