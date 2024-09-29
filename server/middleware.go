package server

import (
	"context"

	"github.com/go-orb/go-orb/util/container"
)

type MiddlewareHandler func(ctx context.Context, req any) (any, error)

type Middleware interface {
	Call(next MiddlewareHandler) MiddlewareHandler
}

var (
	Middlewares = container.NewSafeMap[string, Middleware]() //nolint:gochecknoglobals
)
