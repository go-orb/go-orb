// Package event contains an interface as well as helpers for go-orb events.
package event

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/metadata"
	"github.com/go-orb/go-orb/util/orberrors"
)

// ComponentType is the client component type name.
const ComponentType = "event"

// Handler is the interface for events plugins.
type Handler interface {
	types.Component

	// Request runs a REST like call on the given topic.
	// This is an internal function, clients MUST use event.Request().
	Request(ctx context.Context, req *Req[[]byte, any], opts ...RequestOption) ([]byte, error)

	// HandleRequest subscribes to the given topic and handles the requests.
	// This is a blocking operation.
	// This is an internal function, clients MUST use event.HandleRequest().
	HandleRequest(ctx context.Context, topic string, cb func(context.Context, *Req[[]byte, []byte]))

	// Publish publishes a Event to the given topic.
	// Publish(ctx context.Context, event any) error

	// Subscribe lets you subscribe to the given topic.
	// Subscribe(ctx context.Context, topic string, opts ...SubscribeOption) (<-chan CallRequest[[]byte, []byte], error)
}

// Req contains all data for a request call.
type Req[TReq any, TResp any] struct {
	Topic       string            `json:"topic"`
	ContentType string            `json:"contentType"`
	Metadata    map[string]string `json:"metadata"`

	// The Data of type TReq
	Data TReq `json:"data" yaml:"data"`
	// Err is an error that might happened during encoding.
	Err error

	// ReplyHelper contains the internal helper to answer on exact that topic and request.
	replyFunc func(ctx context.Context, result TResp, err error)

	handler Handler
}

// SetReplyFunc sets the internal reply func (for example nats.Msg) for the client.
func (e *Req[TReq, TResp]) SetReplyFunc(h func(ctx context.Context, result TResp, err error)) {
	e.replyFunc = h
}

// Request runs a REST like call on the events topic.
func (e *Req[TReq, TResp]) Request(ctx context.Context, handler Handler, topic string, opts ...RequestOption) (*TResp, error) {
	e.handler = handler

	options := NewCallOptions(opts...)
	e.ContentType = options.ContentType
	e.Metadata = options.Metadata

	d := []byte{}
	// The err here will be copied into the result.
	codec, err := codecs.GetMime(e.ContentType)
	if err == nil {
		// The err here will be copied into the result.
		d, err = codec.Encode(e.Data)
	}

	bEv := &Req[[]byte, any]{
		Topic:       topic,
		ContentType: e.ContentType,
		Metadata:    e.Metadata,
		Data:        d,
		Err:         err,
		handler:     handler,
	}

	result := new(TResp)

	reply, err := handler.Request(ctx, bEv, opts...)
	if err != nil {
		return result, orberrors.From(err)
	}

	err = codec.Decode(reply, result)
	if err != nil {
		return result, orberrors.From(err)
	}

	return result, nil
}

// NewRequest creates a event for the given topic.
func NewRequest[TResp, TReq any](req TReq) *Req[TReq, TResp] {
	return &Req[TReq, TResp]{
		Data: req,
	}
}

// Request makes a request with using events, it's a shortcut for NewRequest(...).Request(...)
// Example:
//
// resp , err := events.Request[FooResponse](context.Background(), eventsHandler, "user.new", fooRequest)
//
// Response will be of type *FooResponse.
func Request[TResp any, TReq any](
	ctx context.Context,
	handler Handler,
	topic string,
	req TReq,
	opts ...RequestOption,
) (*TResp, error) {
	return NewRequest[TResp](req).Request(ctx, handler, topic, opts...)
}

// HandleRequest subscribes to the given topic and handles the requests.
func HandleRequest[TReq any, TResp any](
	ctx context.Context,
	handler Handler,
	topic string,
	cb func(ctx context.Context, req *TReq) (*TResp, error),
) {
	myCb := func(ctx context.Context, event *Req[[]byte, []byte]) {
		rv := new(TReq)

		// Add metadata to the context.
		myCtx, md := metadata.WithOutgoing(ctx)

		md["Content-Type"] = event.ContentType

		codec, err := codecs.GetMime(event.ContentType)
		if err != nil {
			event.replyFunc(myCtx, nil, err)
			return
		}

		err = codec.Decode(event.Data, rv)
		if err != nil {
			event.replyFunc(myCtx, nil, err)
			return
		}

		// Run the handler.
		result, err := cb(myCtx, rv)
		if err != nil {
			event.replyFunc(myCtx, nil, err)
			return
		}

		// Encode the result and send it back to the plugin.
		d, err := codec.Encode(result)

		// Send the result.
		event.replyFunc(myCtx, d, err)
	}

	go handler.HandleRequest(ctx, topic, myCb)
}

// Provide creates a new client instance with the implementation from cfg.Plugin.
func Provide(
	name types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	opts ...Option) (Handler, error) {
	cfg := NewConfig(opts...)

	sections := append(types.SplitServiceName(name), DefaultConfigSection)
	if err := config.Parse(sections, configs, &cfg); err != nil {
		return nil, err
	}

	if cfg.Plugin == "" {
		logger.Warn("empty event plugin, using the default", "default", DefaultEventPlugin)
		cfg.Plugin = DefaultEventPlugin
	}

	provider, ok := plugins.Get(cfg.Plugin)
	if !ok {
		return nil, fmt.Errorf("event plugin (%s) not found, did you register it?", cfg.Plugin)
	}

	// Configure the logger.
	cLogger, err := logger.WithConfig(sections, configs)
	if err != nil {
		return nil, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	return provider(name, configs, cLogger, opts...)
}
