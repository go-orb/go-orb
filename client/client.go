// Package client provides an interface and helpers for go-orb clients.
package client

import (
	"context"
	"fmt"
	"net/url"

	"log/slog"

	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
	"github.com/go-orb/go-orb/util/orberrors"
)

// ComponentType is the client component type name.
const ComponentType = "client"

// Client is the interface for clients.
type Client interface {
	types.Component

	Config() *Config

	ResolveService(ctx context.Context, service string, preferredTransports ...string) (*container.Map[[]*registry.Node], error)

	NeedsCodec(ctx context.Context, req *Request[any, any], opts ...CallOption) bool

	Call(ctx context.Context, req *Request[any, any], result any, opts ...CallOption) (*RawResponse, error)
	CallNoCodec(ctx context.Context, req *Request[any, any], result any, opts ...CallOption) error
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
	URL         string
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

	// Resolve the service to a list of nodes in a per transport map.
	nodes, err := r.client.ResolveService(ctx, r.service, opts.PreferredTransports...)
	if err != nil {
		return nil, err
	}

	// Run the configured Selector to get a node from the resolved nodes.
	r.node, err = opts.Selector(ctx, r.service, nodes, opts.PreferredTransports, opts.AnyTransport)
	if err != nil {
		r.node = nil
		return nil, err
	}

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

	if r.client.NeedsCodec(ctx, fwReq, opts...) {
		cresp, cerr := r.client.Call(ctx, fwReq, result, opts...)
		if cerr != nil {
			return result, cerr
		}

		codec, err := codecs.GetMime(cresp.ContentType)
		if err != nil {
			return result, orberrors.ErrBadRequest.Wrap(err)
		}

		err = codec.Decode(cresp.Body, result)
		if err != nil {
			return result, orberrors.ErrBadRequest.Wrap(err)
		}

		return result, nil
	}

	cerr := r.client.CallNoCodec(ctx, fwReq, result, opts...)

	return result, cerr
}

// CallResponse is the same as Call with the difference that it returns a Response[*TResp] instead of *TResp.
func (r *Request[TResp, TReq]) CallResponse(ctx context.Context, client Client, opts ...CallOption) (resp Response[*TResp], err error) {
	r.client = client

	var (
		result    = Response[*TResp]{}
		resultVar = new(TResp)
	)

	// Create a copy of Request to forward it.
	// TODO(jochumdev): see if there's a better way to do this.
	fwReq := &Request[any, any]{
		service:  r.service,
		endpoint: r.endpoint,
		request:  r.request,
		client:   r.client,
		node:     r.node,
	}

	if r.client.NeedsCodec(ctx, fwReq, opts...) {
		cresp, cerr := r.client.Call(ctx, fwReq, result, opts...)
		if cerr != nil {
			return result, cerr
		}

		result.ContentType = cresp.ContentType
		result.Headers = cresp.Headers

		codec, err := codecs.GetMime(cresp.ContentType)
		if err != nil {
			return result, orberrors.ErrBadRequest.Wrap(err)
		}

		err = codec.Decode(cresp.Body, result.Body)
		if err != nil {
			return result, orberrors.ErrBadRequest.Wrap(err)
		}

		return result, nil
	}

	cerr := r.client.CallNoCodec(ctx, fwReq, resultVar, opts...)
	result.Body = resultVar

	return result, cerr
}

// NewRequest creates a request for a service+endpoint.
//
// Example (with call):
//
// resp, err := client.NewRequest[FooResponse](
// "service1", "Say.Hello", myRequest,
// ).Call(context.Background(), clientFromWire)
//
// Response will be of type *FooResponse.
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
// Response will be of type *FooResponse.
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

// CallResponse makes a call with the client, it's a shortcut for NewRequest(...).CallResponse(...),
//
// it is the same as Call with the difference that it returns a Response[*TResp] instead of *TResp.
func CallResponse[TResp any, TReq any](
	ctx context.Context,
	client Client,
	service string,
	endpoint string,
	req TReq,
	opts ...CallOption,
) (Response[*TResp], error) {
	return NewRequest[TResp](service, endpoint, req).CallResponse(ctx, client, opts...)
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

	logger.Debug("Client", "plugin", cfg.Plugin)

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
