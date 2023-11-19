package client

import (
	"context"

	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// MiddlewareConfig is the basic config for every middleware.
type MiddlewareConfig struct {
	Name string `json:"name" yaml:"name"`
}

// MiddlewareCallHandler is the middleware handler for client.Call.
type MiddlewareCallHandler func(ctx context.Context, req *Request[any, any], opts *CallOptions) (*RawResponse, error)

type MiddlewareCallNoCodecHandler func(ctx context.Context, req *Request[any, any], result any, opts *CallOptions) error

// Middleware is the middleware for clients.
type Middleware interface {
	// String returns the name of this middleware.
	String()

	Call(
		next MiddlewareCallHandler,
	) MiddlewareCallHandler

	CallNoCodec(
		next MiddlewareCallNoCodecHandler,
	) MiddlewareCallNoCodecHandler
}

// MiddlewareFactory is used to create a new client Middleware.
type MiddlewareFactory func(configSection []string, configs types.ConfigData, client Type, logger log.Logger) (Middleware, error)

// Middlewares contains a map of all available middlewares.
//
//nolint:gochecknoglobals
var Middlewares = container.NewPlugins[MiddlewareFactory]()
