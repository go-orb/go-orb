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
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

// ToError converts the "Error" to "error".
func (e *Error) ToError() error {
	return e
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
)
