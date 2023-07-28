// Package client provides an interface and helpers for go-orb clients.
package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/orberrors"
	"golang.org/x/exp/slog"
)

// ComponentType is the client component type name.
const ComponentType = "client"

// Client is the interface for clients.
type Client interface {
	types.Component

	Config() *Config

	ResolveService(ctx context.Context, service string) (*registry.Node, error)

	Call(ctx context.Context, req *Request[any, any], opts ...CallOption) (resp *RawResponse, err error)
}

// Type is the client type it is returned when you use ProvideClient
// which selects a client to use based on the plugin configuration.
type Type struct {
	Client
}

// RawResponse is a internal struct to pass the transport's response with header and content-type around.
type RawResponse = Response[[]byte]

// Response will be returned by CallWithResponse.
type Response[T any] struct {
	ContentType string
	Headers     map[string][]string
	Body        T
}

// Request is a request for Client.
type Request[TResp any, TReq any] struct {
	service  string
	endpoint string

	// The unencoded request
	request TReq

	client Client

	node *registry.Node
}

// Service returns the Service from the request.
func (r *Request[TResp, TReq]) Service() string {
	return r.service
}

// Endpoint returns the Endpoint from the request.
func (r *Request[TResp, TReq]) Endpoint() string {
	return r.endpoint
}

// Request returns the Request.
func (r *Request[TResp, TReq]) Request() TReq {
	return r.request
}

// Node returns the Node.
func (r *Request[TResp, TReq]) Node(ctx context.Context, opts *CallOptions) (*registry.Node, error) {
	if r.node != nil {
		return r.node, nil
	}

	if opts.URL != "" {
		myU1rl, err := url.Parse(opts.URL)
		if err != nil {
			return nil, orberrors.ErrBadRequest.Wrap(err)
		}

		node := &registry.Node{
			ID:        "url",
			Address:   myU1rl.Host,
			Transport: myU1rl.Scheme,
		}

		r.node = node

		return node, nil
	}

	// Resolve the service and set the request's node.
	node, err := r.client.ResolveService(ctx, r.service)
	if err != nil {
		return nil, err
	}

	r.node = node

	return r.node, nil
}

// Call forward's the Request to Client.Call() and decodes the result into resp with the type TResp.
func (r *Request[TResp, TReq]) Call(ctx context.Context, client Client, opts ...CallOption) (resp *TResp, err error) {
	r.client = client

	var result = new(TResp)

	// Create a copy of Request to forward it.
	// TODO(jochumdev): see if there's a better way to do this.
	fwReq := &Request[any, any]{
		service:  r.service,
		endpoint: r.endpoint,
		request:  r.request,
		client:   r.client,
		node:     r.node,
	}

	cresp, cerr := r.client.Call(ctx, fwReq, opts...)
	if cerr != nil {
		return result, cerr
	}

	codec, err := codecs.GetMime(cresp.ContentType)
	if err != nil {
		return result, orberrors.ErrBadRequest.Wrap(err)
	}

	err = codec.Unmarshal(cresp.Body, result)
	if err != nil {
		return result, orberrors.ErrBadRequest.Wrap(err)
	}

	return result, nil
}

// NewRequest creates a request for a service+endpoint.
//
// Example (with call):
//
// resp, err := client.NewRequest[FooResponse](
// "service1", "Say.Hello", myRequest,
// ).Call(context.Background(), clientFromWire)
//
// Response will be of type FooResponse.
func NewRequest[TResp any, TReq any](
	service string,
	endpoint string,
	req TReq,
) *Request[TResp, TReq] {
	return &Request[TResp, TReq]{
		service:  service,
		endpoint: endpoint,

		request: req,
	}
}

// Call makes a call with the client, it's a shortcut for NewRequest(...).Call(...)
// Example:
//
// resp , err := client.Call[FooResponse](context.Background(), someClient, "service1", "Say.Hello", fooRequest)
//
// Response will be of type FooResponse.
func Call[TResp any, TReq any](
	ctx context.Context,
	client Client,
	service string,
	endpoint string,
	req TReq,
	opts ...CallOption,
) (*TResp, error) {
	return NewRequest[TResp](service, endpoint, req).Call(ctx, client, opts...)
}

// ProvideClient creates a new client instance with the implementation from cfg.Plugin.
func ProvideClient(
	name types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	reg registry.Type,
	opts ...Option) (Type, error) {
	cfg := NewConfig(opts...)

	sections := append(types.SplitServiceName(name), DefaultConfigSection)
	if err := config.Parse(sections, configs, cfg); err != nil {
		return Type{}, err
	}

	if cfg.Plugin == "" {
		logger.Warn("empty client plugin, using the default", "default", DefaultClientPlugin)
		cfg.Plugin = DefaultClientPlugin
	}

	provider, err := plugins.Get(cfg.Plugin)
	if err != nil {
		return Type{}, fmt.Errorf("client plugin (%s) not found, did you register it: %w", cfg.Plugin, err)
	}

	// Configure the logger.
	cLogger, err := logger.WithConfig(sections, configs)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	return provider(name, configs, cLogger, reg, opts...)
}