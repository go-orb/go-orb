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
	return fmt.Sprintf("%d %s: %s", e.Code, e.Message, e.Wrapped)
}

// Toerror converts the "Error" to "error",
// same as doing <variableWithTypeError>.(error).
func (e *Error) Toerror() error {
	return e
}

// Wrap wraps another error into a copy.
func (e *Error) Wrap(err error) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Wrapped: err,
	}
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Wrapped
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

// From converts an error to orberrors.Error.
func From(err error) *Error {
	var orbe *Error
	if errors.As(err, &orbe) {
		return orbe
	}

	return New(http.StatusInternalServerError, err.Error())
}

// A list of default errors.
var (
	ErrInternalServerError = NewHTTP(http.StatusInternalServerError)
	ErrUnauthorized        = NewHTTP(http.StatusUnauthorized)
	ErrRequestTimeout      = NewHTTP(http.StatusRequestTimeout)
	ErrBadRequest          = NewHTTP(http.StatusBadRequest)
)
