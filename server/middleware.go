package server

import (
	"context"

	"github.com/go-orb/go-orb/util/container"
)

// MiddlewareCallHandler is the Handler for unitary RPC Calls.
type MiddlewareCallHandler func(ctx context.Context, req any) (any, error)

// Middleware is an interface that must be implemented by server Middlewares.
type Middleware interface {
	Call(next MiddlewareCallHandler) MiddlewareCallHandler
}

var (
	// Middlewares is the map of Middlewares available for servers.
	Middlewares = container.NewSafeMap[string, Middleware]() //nolint:gochecknoglobals
)
