package server

import (
	"context"

	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// MiddlewareCallHandler is the Handler for unary RPC Calls.
type MiddlewareCallHandler func(ctx context.Context, req any) (any, error)

// MiddlewareStreamHandler is the handler for streaming RPC calls.
type MiddlewareStreamHandler func(ctx context.Context) error

// Middleware is an interface that must be implemented by server Middlewares.
type Middleware interface {
	// Start/Stop will be called many times, if you require the call once, you have to ensure it on your own.
	types.Component

	Call(next MiddlewareCallHandler) MiddlewareCallHandler
}

// MiddlewareProvider is the provider for a middleware, each Middleware must supply this to register itself.
type MiddlewareProvider func(
	configSection []string,
	configKey string,
	configData map[string]any,
	logger log.Logger,
) (Middleware, error)

var (
	// Middlewares is the map of Middlewares available for servers.
	Middlewares = container.NewSafeMap[string, MiddlewareProvider]() //nolint:gochecknoglobals
)
