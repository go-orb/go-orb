// Package event contains an interface as well as helpers for go-orb events.
package event

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/orberrors"
)

// ComponentType is the client component type name.
const ComponentType = "event"

// AckFunc is the function to call to acknowledge a message.
type AckFunc func() error

// NackFunc is the function to call to negatively acknowledge a message.
type NackFunc func() error

// Event is the object returned by the broker when you subscribe to a topic.
type Event struct {
	// Handler is a reference to the client.
	Handler Client

	// Timestamp of the event
	Timestamp time.Time
	// Metadata contains the values the event was indexed by
	Metadata map[string]string

	ackFunc  AckFunc
	nackFunc NackFunc
	// ID to uniquely identify the event
	ID string
	// Topic of event, e.g. "registry.service.created"
	Topic string
	// Payload contains the encoded message
	Payload []byte
}

// Unmarshal the events message into an object.
func (e *Event) Unmarshal(v any) error {
	return e.Handler.GetPublishCodec().Unmarshal(e.Payload, v)
}

// Ack acknowledges successful processing of the event in ManualAck mode.
func (e *Event) Ack() error {
	return e.ackFunc()
}

// SetAckFunc sets the AckFunc for the event.
func (e *Event) SetAckFunc(f AckFunc) {
	e.ackFunc = f
}

// Nack negatively acknowledges processing of the event (i.e. failure) in ManualAck mode.
func (e *Event) Nack() error {
	return e.nackFunc()
}

// SetNackFunc sets the NackFunc for the event.
func (e *Event) SetNackFunc(f NackFunc) {
	e.nackFunc = f
}

// Client is the client interface for events plugins.
type Client interface {
	types.Component

	// Request runs a REST like call on the given topic.
	// This is an internal function, clients MUST use event.Request().
	Request(ctx context.Context, req *Req[[]byte, any], opts ...RequestOption) ([]byte, error)

	// HandleRequest subscribes to the given topic and handles the requests.
	// This is a blocking operation.
	// This is an internal function, clients MUST use event.HandleRequest().
	HandleRequest(ctx context.Context, topic string, cb func(context.Context, *Req[[]byte, []byte]))

	// Clone creates a clone of the handler, this is useful for parallel requests.
	Clone() Client

	// GetCodec returns the codec used by the handler for publish and subscribe.
	GetPublishCodec() codecs.Marshaler

	// Publish publishes a Event to the given topic.
	Publish(ctx context.Context, topic string, event any, opts ...PublishOption) error

	// Consume lets you consume events from a given topic.
	Consume(topic string, opts ...ConsumeOption) (<-chan Event, error)
}

// Type is the client implementation for events.
type Type struct {
	Client
}

// Req contains all data for a request call.
type Req[TReq any, TResp any] struct {
	Topic       string `json:"topic"`
	ContentType string `json:"contentType"`

	// The Data of type TReq
	Data TReq `json:"data" yaml:"data"`
	// Err is an error that might happened during encoding.
	Err error

	// ReplyHelper contains the internal helper to answer on exact that topic and request.
	replyFunc func(ctx context.Context, result TResp, err error)

	handler Client
}

// SetReplyFunc sets the internal reply func (for example nats.Msg) for the client.
func (e *Req[TReq, TResp]) SetReplyFunc(h func(ctx context.Context, result TResp, err error)) {
	e.replyFunc = h
}

// Request runs a REST like call on the events topic.
func (e *Req[TReq, TResp]) Request(ctx context.Context, handler Client, topic string, opts ...RequestOption) (*TResp, error) {
	e.handler = handler

	options := NewRequestOptions(opts...)
	e.ContentType = options.ContentType

	d := []byte{}
	// The err here will be copied into the result.
	codec, err := codecs.GetMime(e.ContentType)
	if err == nil {
		// The err here will be copied into the result.
		d, err = codec.Marshal(e.Data)
	}

	bEv := &Req[[]byte, any]{
		Topic:       topic,
		ContentType: e.ContentType,
		Data:        d,
		Err:         err,
		handler:     handler,
	}

	result := new(TResp)

	reply, err := handler.Request(ctx, bEv, opts...)
	if err != nil {
		return result, orberrors.From(err)
	}

	err = codec.Unmarshal(reply, result)
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
	handler Client,
	topic string,
	req TReq,
	opts ...RequestOption,
) (*TResp, error) {
	return NewRequest[TResp](req).Request(ctx, handler, topic, opts...)
}

// HandleRequest subscribes to the given topic and handles the requests.
func HandleRequest[TReq any, TResp any](
	ctx context.Context,
	handler Client,
	topic string,
	callback func(ctx context.Context, req *TReq) (*TResp, error),
) {
	myCb := func(ctx context.Context, event *Req[[]byte, []byte]) {
		rv := new(TReq)

		codec, err := codecs.GetMime(event.ContentType)
		if err != nil {
			event.replyFunc(ctx, nil, err)
			return
		}

		err = codec.Unmarshal(event.Data, rv)
		if err != nil {
			event.replyFunc(ctx, nil, err)
			return
		}

		// Run the handler.
		result, err := callback(ctx, rv)
		if err != nil {
			event.replyFunc(ctx, nil, err)
			return
		}

		// Encode the result and send it back to the plugin.
		d, err := codec.Marshal(result)

		// Send the result.
		event.replyFunc(ctx, d, err)
	}

	go handler.HandleRequest(ctx, topic, myCb)
}

// Provide creates a new client instance with the implementation from cfg.Plugin.
func Provide(
	name types.ServiceName,
	configs types.ConfigData,
	components *types.Components,
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

	instance, err := provider(name, configs, cLogger, opts...)
	if err != nil {
		return Type{}, err
	}

	// Register the event as a component.
	err = components.Add(instance, types.PriorityEvent)
	if err != nil {
		logger.Warn("while registering event as a component", "error", err)
	}

	return instance, nil
}

// ProvideNoOpts creates a new client instance with the implementation from cfg.Plugin.
func ProvideNoOpts(
	name types.ServiceName,
	configs types.ConfigData,
	components *types.Components,
	logger log.Logger) (Type, error) {
	return Provide(name, configs, components, logger)
}
