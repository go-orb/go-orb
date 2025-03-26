// Package client provides an interface and helpers for go-orb clients.
package client

import (
	"context"
	"errors"
	"fmt"

	"log/slog"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
)

// ComponentType is the client component type name.
const ComponentType = "client"

// Client is the interface for clients.
type Client interface {
	types.Component

	// Config returns the internal config, this is for tests.
	Config() Config

	// With closes all transports and configures the client with the given options.
	With(opts ...Option) error

	// SelectService selects a service node.
	SelectService(ctx context.Context, service string, opts ...CallOption) (string, string, error)

	// Request does the actual call.
	Request(ctx context.Context, service string, endpoint string, req any, result any, opts ...CallOption) error

	// Stream creates a streaming client to the specified service endpoint.
	Stream(ctx context.Context, service string, endpoint string, opts ...CallOption) (StreamIface[any, any], error)
}

// Type is the client type it is returned when you use ProvideClient
// which selects a client to use based on the plugin configuration.
type Type struct {
	Client
}

// RequestInfosKey is the key for the request infos in the context.
type RequestInfosKey struct{}

// RequestInfos contains the request infos.
type RequestInfos struct {
	Service   string
	Endpoint  string
	Transport string
	Address   string
}

// RequestInfo returns the request infos from the context.
func RequestInfo(ctx context.Context) (RequestInfos, bool) {
	v, ok := ctx.Value(RequestInfosKey{}).(*RequestInfos)
	return *v, ok
}

// Request is a typesafe shortcut for making a request.
//
// Example:
//
// resp , err := client.Request[FooResponse](context.Background(), clientWire, "service1", "Say.Hello", fooRequest)
//
// Response will be of type *FooResponse.
func Request[TResp any, TReq any](
	ctx context.Context,
	client Client,
	service string,
	endpoint string,
	req TReq,
	opts ...CallOption,
) (*TResp, error) {
	result := new(TResp)

	err := client.Request(ctx, service, endpoint, req, result, opts...)

	return result, err
}

// New creates a new client instance with the implementation from cfg.Plugin.
func New(
	configData map[string]any,
	components *types.Components,
	logger log.Logger,
	registry registry.Type,
	opts ...Option,
) (Type, error) {
	cfg := NewConfig(opts...)

	if err := config.Parse(nil, DefaultConfigSection, configData, &cfg); err != nil && !errors.Is(err, config.ErrNoSuchKey) {
		return Type{}, err
	}

	if cfg.Plugin == "" {
		logger.Warn("empty client plugin, using the default", "default", DefaultClientPlugin)
		cfg.Plugin = DefaultClientPlugin
	}

	provider, ok := plugins.Get(cfg.Plugin)
	if !ok {
		return Type{}, fmt.Errorf("client plugin (%s) not found, did you register it?", cfg.Plugin)
	}

	// Configure the logger.
	cLogger, err := logger.WithConfig([]string{DefaultConfigSection}, configData)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	return provider(configData, components, cLogger, registry, opts...)
}

// Provide creates a new client instance with the implementation from cfg.Plugin.
func Provide(
	svcCtx *cli.ServiceContextWithConfig,
	components *types.Components,
	logger log.Logger,
	reg registry.Type,
	opts ...Option) (Type, error) {
	return New(svcCtx.Config(), components, logger, reg, opts...)
}

// ProvideNoOpts provides a new client without options.
func ProvideNoOpts(
	svcCtx *cli.ServiceContextWithConfig,
	components *types.Components,
	logger log.Logger,
	reg registry.Type,
) (Type, error) {
	return Provide(svcCtx, components, logger, reg)
}
