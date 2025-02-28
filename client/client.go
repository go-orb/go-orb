// Package client provides an interface and helpers for go-orb clients.
package client

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"log/slog"

	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/orberrors"
)

// ComponentType is the client component type name.
const ComponentType = "client"

// NodeMap is the type for a string map with list of registry nodes.
type NodeMap map[string][]*registry.Node

// Client is the interface for clients.
type Client interface {
	types.Component

	// Config returns the internal config, this is for tests.
	Config() Config

	// With closes all transports and configures the client with the given options.
	With(opts ...Option) error

	ResolveService(ctx context.Context, service string, preferredTransports ...string) (NodeMap, error)

	// NeedsCodec has to do node resolving and then selects the right transport for that node,
	// it then has to return whatever the selected transport needs a codec or if it does encoding internaly.
	NeedsCodec(ctx context.Context, req *Request[any, any], opts ...CallOption) bool

	// Call with encoding on client side.
	Call(ctx context.Context, req *Request[any, any], result any, opts ...CallOption) (*RawResponse, error)

	// CallNoCodec is the same as Call but without encoding.
	CallNoCodec(ctx context.Context, req *Request[any, any], result any, opts ...CallOption) error
}

// Type is the client type it is returned when you use ProvideClient
// which selects a client to use based on the plugin configuration.
type Type struct {
	Client
}

// RawResponse is a internal struct to pass the transport's response with metadata and content-type around.
type RawResponse = Response[io.Reader]

// Response will be returned by CallWithResponse.
type Response[T any] struct {
	ContentType string
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

	// Create a [any, any] copy of Request to forward it.
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

		err = codec.NewDecoder(cresp.Body).Decode(result)
		if err != nil {
			return result, orberrors.ErrBadRequest.Wrap(err)
		}

		return result, nil
	}

	cerr := r.client.CallNoCodec(ctx, fwReq, result, opts...)

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
// resp , err := client.Call[FooResponse](context.Background(), clientWire, "service1", "Say.Hello", fooRequest)
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

// Provide creates a new client instance with the implementation from cfg.Plugin.
func Provide(
	name types.ServiceName,
	configs types.ConfigData,
	components *types.Components,
	logger log.Logger,
	reg registry.Type,
	opts ...Option) (Type, error) {
	cfg := NewConfig(opts...)

	sections := append(types.SplitServiceName(name), DefaultConfigSection)
	if err := config.Parse(sections, configs, &cfg); err != nil {
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
	cLogger, err := logger.WithConfig(sections, configs)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	return provider(name, configs, components, cLogger, reg, opts...)
}
