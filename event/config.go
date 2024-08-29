package event

import (
	"time"

	"github.com/go-orb/go-orb/util/metadata"
)

//nolint:gochecknoglobals
var (
	// DefaultClient is the default client implementation to use.
	DefaultEventPlugin = "natsjs"

	// DefaultConfigSection is the default config section for the client.
	DefaultConfigSection = "client"

	// Default Content Type used to transport data around.
	DefaultContentType = "application/protobuf"

	// DefaultRequestTimeout is the default request timeout.
	DefaultRequestTimeout = time.Second * 30
)

var _ (ConfigType) = (*Config)(nil)

// Option is a functional option type for the registry.
type Option func(ConfigType)

// ConfigType is used in the functional options as type to identify a registry
// option. It is used over a static *Config type as this way plugins can also
// easilty set functional options without the complication of contexts, as was
// done in v4. This is possible because plugins will nest the registry.Config
// type, and thus inherit the interface that is used to identify the registry
// config.
type ConfigType interface {
	config() *Config
}

// Config are the Client options.
type Config struct {
	// Plugin selects the client implementation.
	Plugin string `json:"plugin" yaml:"plugin"`
}

func (c *Config) config() *Config {
	return c
}

// WithClientPlugin set the client implementation to use.
func WithClientPlugin(n string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Plugin = n
	}
}

// NewConfig generates a new config with all the defaults.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		Plugin: DefaultEventPlugin,
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// CallOptions contains options for a call.
type CallOptions struct {
	// ContentType for transporting the message.
	ContentType string
	// Metadata contains any keys which can be used to query the data, for example a customer id
	Metadata metadata.Metadata

	// RequestTimeout defines how long to wait for the server to reply on a request.
	RequestTimeout time.Duration
}

// RequestOption sets attributes on Calloptions.
type RequestOption func(o *CallOptions)

// WithCallContentType sets the ContentType field to use for transporting the message.
func WithCallContentType(ct string) RequestOption {
	return func(o *CallOptions) {
		o.ContentType = ct
	}
}

// WithCallMetadata sets the Metadata field on CallOptions.
func WithCallMetadata(md metadata.Metadata) RequestOption {
	return func(o *CallOptions) {
		o.Metadata = md
	}
}

// WithCallRequestTimeout sets the timeout for a request.
func WithCallRequestTimeout(t time.Duration) RequestOption {
	return func(o *CallOptions) {
		o.RequestTimeout = t
	}
}

// NewCallOptions generates new calloptions with defaults.
func NewCallOptions(opts ...RequestOption) CallOptions {
	cfg := CallOptions{
		ContentType:    DefaultContentType,
		Metadata:       metadata.Metadata{},
		RequestTimeout: DefaultRequestTimeout,
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// // SubscribeOptions contains all the options which can be provided when subscribing to a topic.
// type SubscribeOptions struct {
// 	// Offset is the time from which the messages should be consumed from. If not provided then
// 	// the messages will be consumed starting from the moment the Subscription starts.
// 	Offset time.Time
// 	// ConsumerGroup is the name of the consumer group, if two consumers have the same group the events
// 	// are distributed between them
// 	ConsumerGroup string
// 	AckWait       time.Duration
// 	// AutoAck if true (default true), automatically acknowledges every message so it will not be redelivered.
// 	// If false specifies that each message need ts to be manually acknowledged by the subscriber.
// 	// If processing is successful the message should be ack'ed to remove the message from the stream.
// 	// If processing is unsuccessful the message should be nack'ed (negative acknowledgement) which will mean it will
// 	// remain on the stream to be processed again.
// 	AutoAck bool
// }

// // SubscribeOption sets attributes on SubscribeOptions.
// type SubscribeOption func(o *SubscribeOptions)

// // WithSubscribeConsumerGroup sets the consumer group to be part of when consuming events.
// func WithSubscribeConsumerGroup(q string) SubscribeOption {
// 	return func(o *SubscribeOptions) {
// 		o.ConsumerGroup = q
// 	}
// }

// // WithSubscribeOffset sets the offset time at which to start consuming events.
// func WithSubscribeOffset(t time.Time) SubscribeOption {
// 	return func(o *SubscribeOptions) {
// 		o.Offset = t
// 	}
// }

// // WithSubscribeAutoAck sets the AutoAck field on SubscribeOptions and an ackWait duration after which if no ack is received
// // the message is requeued in case auto ack is turned off.
// func WithSubscribeAutoAck(ack bool, ackWait time.Duration) SubscribeOption {
// 	return func(o *SubscribeOptions) {
// 		o.AutoAck = ack
// 		o.AckWait = ackWait
// 	}
// }
