// Package orberrors provides an error that must be transported on client<->server operations.
package orberrors

import (
	"errors"
	"fmt"
	"net/http"
)

// Error is the orb error, it contain's a "Code" additional to the Message.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Wrapped error  `json:"wrapped"`
}

func (e *Error) Error() string {
	if e.Wrapped == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Message, e.Wrapped.Error())
}

// Toerror converts the "Error" to "error",
// same as doing <variableWithTypeError>.(error).
func (e *Error) Toerror() error {
	if e == nil {
		return nil
	}

	return e
}

// Wrap wraps another error into a copy.
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

	return orberr.Code == e.Code && orberr.Message == e.Message && errors.Is(orberr.Wrapped, e.Wrapped)
}

// New creates a new orb error with the given parameters.
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewHTTP creates an orb error with the given status code and a static message.
func NewHTTP(code int) *Error {
	return &Error{
		Code:    code,
		Message: http.StatusText(code),
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
	ErrUnimplemented       = New(http.StatusInternalServerError, "Unimplemented")
	ErrInternalServerError = NewHTTP(http.StatusInternalServerError)
	ErrUnauthorized        = NewHTTP(http.StatusUnauthorized)
	ErrRequestTimeout      = NewHTTP(http.StatusRequestTimeout)
	ErrBadRequest          = NewHTTP(http.StatusBadRequest)
)
