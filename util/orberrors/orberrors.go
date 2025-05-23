// Package orberrors provides an error that must be transported on client<->server operations.
package orberrors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Error is the orb error, it contain's a "Code" additional to the Message.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Wrapped error  `json:"wrapped"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Wrapped.Error())
	}

	return e.Message
}

// Toerror converts the "Error" to "error",
// same as doing <variableWithTypeError>.(error).
func (e *Error) Toerror() error {
	if e == nil {
		return nil
	}

	return e
}

// Wrap creates a new error and wraps another error.
func (e *Error) Wrap(err error) *Error {
	if e == nil {
		return nil
	}

	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Wrapped: err,
	}
}

// WrapNew creates a new error and wraps another error with a new message.
func (e *Error) WrapNew(message string) *Error {
	if e == nil {
		return nil
	}

	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Wrapped: errors.New(message),
	}
}

// WrapF creates a new error and wraps another error with a new message and args.
func (e *Error) WrapF(message string, args ...any) *Error {
	if e == nil {
		return nil
	}

	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Wrapped: fmt.Errorf(message, args...),
	}
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Wrapped
}

// Is returns if the given error is an (orb)Error.
func (e *Error) Is(err error) bool {
	orberr, ok := err.(*Error)

	if !ok {
		return false
	}

	return orberr.Code == e.Code && orberr.Message == e.Message
}

// New creates a new orb error with the given parameters.
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// newHTTP creates the HTTP error from the given code.
func newHTTP(code int) *Error {
	return &Error{
		Code:    code,
		Message: strings.ToLower(http.StatusText(code)),
	}
}

// HTTP returns an orb error with the given status code and a static message.
func HTTP(code int) *Error {
	switch code {
	case 503:
		return ErrUnavailable
	case 501:
		return ErrNotImplemented
	case 500:
		return ErrInternalServerError
	case 499:
		return ErrCanceled
	case 401:
		return ErrUnauthorized
	case 408:
		return ErrRequestTimeout
	case 400:
		return ErrBadRequest
	default:
		return newHTTP(code)
	}
}

// From wraps an error into orberrors.Error.
func From(err error) *Error {
	// nil input = nil output.
	if err == nil {
		return nil
	}

	// Already an orberror?
	var orbErr *Error
	if errors.As(err, &orbErr) {
		return orbErr
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return ErrRequestTimeout.Wrap(err)
	}

	// Else make a copy.
	return &Error{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Wrapped: err,
	}
}

// As tries to convert an `error` into an `*orberrors.Error`.
func As(err error) (*Error, bool) {
	var orbe *Error

	if errors.As(err, &orbe) {
		return orbe, true
	}

	return orbe, false
}

// A list of default errors.
var (
	ErrBadRequest          = newHTTP(http.StatusBadRequest)     // 400
	ErrUnauthorized        = newHTTP(http.StatusUnauthorized)   // 401
	ErrNotFound            = newHTTP(http.StatusNotFound)       // 404
	ErrRequestTimeout      = newHTTP(http.StatusRequestTimeout) // 408
	ErrCanceled            = newHTTP(499)
	ErrUnimplemented       = newHTTP(http.StatusInternalServerError) // 500
	ErrNotImplemented      = newHTTP(http.StatusNotImplemented)      // 501
	ErrUnavailable         = newHTTP(http.StatusServiceUnavailable)  // 503
	ErrInternalServerError = newHTTP(http.StatusInternalServerError) // 500
)
