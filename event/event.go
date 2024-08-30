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

// Events is the interface for events plugins.
type Events interface {
	types.Component

	// Request runs a REST like call on the given topic.
	Request(ctx context.Context, topic string, req *Call[[]byte, any], opts ...RequestOption) ([]byte, error)

	// HandleRequest subscribes to the given topic and handles the requests.
	HandleRequest(ctx context.Context, topic string) (<-chan Call[[]byte, []byte], error)

	// Publish publishes a Event to the given topic.
	// Publish(ctx context.Context, ev *Event[[]byte, []byte], opts ...PublishOption) error

	// Subscribe lets you subscribe to the given topic.
	// Subscribe(ctx context.Context, topic string, opts ...SubscribeOption) (<-chan CallRequest[[]byte, []byte], error)
}

// Type is the client type it is returned when you use ProvideClient
// which selects a client to use based on the plugin configuration.
type Type struct {
	Events
}

// Call contains all data for a request call.
type Call[TReq any, TResp any] struct {
	ContentType string            `json:"contentType"`
	Metadata    metadata.Metadata `json:"metadata"`

	// The Data of type TReq
	Data TReq `json:"data" yaml:"data"`
	// Err is an error that might happened during encoding.
	Err error

	// ReplyHelper contains the internal helper to answer on exact that topic and request.
	replyFunc func(ctx context.Context, result TResp, err error) error

	client Events
}

// SetReplyFunc sets the internal reply func (for example nats.Msg) for the client.
func (e *Call[TReq, TResp]) SetReplyFunc(h func(ctx context.Context, result TResp, err error) error) {
	e.replyFunc = h
}

// Request runs a REST like call on the events topic.
func (e *Call[TReq, TResp]) Request(ctx context.Context, client Events, topic string, opts ...RequestOption) (*TResp, error) {
	e.client = client

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

	bEv := &Call[[]byte, any]{
		ContentType: e.ContentType,
		Metadata:    e.Metadata,
		Data:        d,
		Err:         orberrors.From(err),
		client:      client,
	}

	result := new(TResp)

	reply, err := client.Request(ctx, topic, bEv, opts...)
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
func NewRequest[TResp, TReq any](req TReq) *Call[TReq, TResp] {
	return &Call[TReq, TResp]{
		Data: req,
	}
}

// Request makes a request with using events, it's a shortcut for NewRequest(...).Request(...)
// Example:
//
// resp , err := events.Request[FooResponse](context.Background(), eventsWire, "user.new", fooRequest)
//
// Response will be of type *FooResponse.
func Request[TResp any, TReq any](
	ctx context.Context,
	eventsWire Events,
	topic string,
	req TReq,
	opts ...RequestOption,
) (*TResp, error) {
	return NewRequest[TResp](req).Request(ctx, eventsWire, topic, opts...)
}

// HandleRequest subscribes to the given topic and handles the requests.
func HandleRequest[TReq any, TResp any](
	ctx context.Context,
	eventsWire Events,
	topic string,
	handler func(ctx context.Context, req *TReq) (*TResp, error),
) (context.CancelFunc, error) {
	ctx, cancelFunc := context.WithCancel(ctx)

	inChan, err := eventsWire.HandleRequest(ctx, topic)
	if err != nil {
		cancelFunc()
		return nil, fmt.Errorf("%w: %w", orberrors.ErrInternalServerError, err)
	}

	// This go routine transforms the encoded request from inChan into a decoded request to outChan.
	go func(ctx context.Context, inChan <-chan Call[[]byte, []byte]) {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-inChan:
				rv := new(TReq)

				// Add metadata to the context.
				md := e.Metadata
				md.Set("Content-Type", e.ContentType)
				hCtx := md.To(ctx)

				codec, err := codecs.GetMime(e.ContentType)
				if err != nil {
					e.replyFunc(hCtx, nil, err)
				}

				err = codec.Decode(e.Data, rv)
				if err != nil {
					e.replyFunc(hCtx, nil, err)
				}

				// Run the handler.
				result, err := handler(hCtx, rv)
				if err != nil {
					e.replyFunc(hCtx, nil, err)
				}

				// Encode the result and send it back to the plugin.
				d, err := codec.Encode(result)

				// Send the result.
				e.replyFunc(hCtx, d, err)
			}
		}
	}(ctx, inChan)

	return cancelFunc, nil
}

// Provide creates a new client instance with the implementation from cfg.Plugin.
func Provide(
	name types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	opts ...Option) (Type, error) {
	cfg := NewConfig(opts...)

	sections := append(types.SplitServiceName(name), DefaultConfigSection)
	if err := config.Parse(sections, configs, &cfg); err != nil {
		return Type{}, err
	}

	if cfg.Plugin == "" {
		logger.Warn("empty event plugin, using the default", "default", DefaultEventPlugin)
		cfg.Plugin = DefaultEventPlugin
	}

	provider, ok := plugins.Get(cfg.Plugin)
	if !ok {
		return Type{}, fmt.Errorf("event plugin (%s) not found, did you register it?", cfg.Plugin)
	}

	// Configure the logger.
	cLogger, err := logger.WithConfig(sections, configs)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	return provider(name, configs, cLogger, opts...)
}
