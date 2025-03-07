package client

import (
	"context"

	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// MiddlewareComponentType is returned when you call SomeMiddleware.Type().
const MiddlewareComponentType = "middleware"

// MiddlewareConfig is the basic config for every middleware.
type MiddlewareConfig struct {
	Name string `json:"name" yaml:"name"`
}

// MiddlewareRequestHandler is the middleware handler for client.Request.
type MiddlewareRequestHandler func(ctx context.Context, req *Req[any, any], opts *CallOptions) (*RawResponse, error)

// MiddlewareRequestNoCodecHandler is the middleware handler for client.Request without a codec in between.
type MiddlewareRequestNoCodecHandler func(ctx context.Context, req *Req[any, any], result any, opts *CallOptions) error

// Middleware is the middleware for clients.
type Middleware interface {
	types.Component

	Request(
		next MiddlewareRequestHandler,
	) MiddlewareRequestHandler

	RequestNoCodec(
		next MiddlewareRequestNoCodecHandler,
	) MiddlewareRequestNoCodecHandler
}

// MiddlewareFactory is used to create a new client Middleware.
type MiddlewareFactory func(configSection []string, configs types.ConfigData, client Type, logger log.Logger) (Middleware, error)

// Middlewares contains a map of all available middlewares.
//
//nolint:gochecknoglobals
var Middlewares = container.NewMap[string, MiddlewareFactory]()
