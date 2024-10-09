package server

import "errors"

// Errors.
var (
	ErrUnknownMiddleware = errors.New("unknown middleware")
	ErrUnknownHandler    = errors.New("unknown handler")
	ErrUnknownEntrypoint = errors.New("unknown entrypoint")
)
