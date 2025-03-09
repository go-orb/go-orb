package event

import (
	"time"

	"github.com/go-orb/go-orb/codecs"
	"github.com/google/uuid"
)

//nolint:gochecknoglobals
var (
	// DefaultEventPlugin is the default client implementation to use.
	DefaultEventPlugin = "natsjs"

	// DefaultConfigSection is the default config section for the client.
	DefaultConfigSection = "client"

	// DefaultRequestContentType is the default content type used to transport data around.
	DefaultRequestContentType = codecs.MimeProto

	// DefaultPublishContentType is the default content type used to transport events around.
	DefaultPublishContentType = codecs.MimeJSON

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

// RequestOptions contains options for a call.
type RequestOptions struct {
	// ContentType for transporting the message.
	ContentType string
	// Metadata contains keys which can be used to query the data, for example a customer id
	Metadata map[string]string

	// RequestTimeout defines how long to wait for the server to reply on a request.
	RequestTimeout time.Duration
}

// RequestOption sets attributes on Calloptions.
type RequestOption func(o *RequestOptions)

// WithRequestContentType sets the ContentType field to use for transporting the message.
func WithRequestContentType(ct string) RequestOption {
	return func(o *RequestOptions) {
		o.ContentType = ct
	}
}

// WithRequestResponseMetadata will write response Metadata into the given map.
func WithRequestResponseMetadata(md map[string]string) RequestOption {
	return func(o *RequestOptions) {
		o.Metadata = md
	}
}

// WithRequestTimeout sets the timeout for a request.
func WithRequestTimeout(t time.Duration) RequestOption {
	return func(o *RequestOptions) {
		o.RequestTimeout = t
	}
}

// NewRequestOptions generates new calloptions with defaults.
func NewRequestOptions(opts ...RequestOption) RequestOptions {
	cfg := RequestOptions{
		ContentType:    DefaultRequestContentType,
		Metadata:       make(map[string]string),
		RequestTimeout: DefaultRequestTimeout,
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// PublishOptions contains all the options which can be provided when publishing an event.
type PublishOptions struct {
	// Metadata contains any keys which can be used to query the data, for example a customer id
	Metadata map[string]string
	// Timestamp to set for the event, if the timestamp is a zero value, the current time will be used
	Timestamp time.Time
}

// PublishOption sets attributes on PublishOptions.
type PublishOption func(o *PublishOptions)

// WithPublishMetadata sets the Metadata field on PublishOptions.
func WithPublishMetadata(md map[string]string) PublishOption {
	return func(o *PublishOptions) {
		o.Metadata = md
	}
}

// WithPublishTimestamp sets the timestamp field on PublishOptions.
func WithPublishTimestamp(t time.Time) PublishOption {
	return func(o *PublishOptions) {
		o.Timestamp = t
	}
}

// NewPublishOptions generates new publish options with defaults.
func NewPublishOptions(opts ...PublishOption) PublishOptions {
	cfg := PublishOptions{
		Metadata:  make(map[string]string),
		Timestamp: time.Now(),
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// ConsumeOptions contains all the options which can be provided when subscribing to a topic.
type ConsumeOptions struct {
	// Offset is the time from which the messages should be consumed from. If not provided then
	// the messages will be consumed starting from the moment the Subscription starts.
	Offset time.Time
	// Group is the name of the consumer group, if two consumers have the same group the events
	// are distributed between them
	Group   string
	AckWait time.Duration
	// RetryLimit indicates number of times a message is retried
	RetryLimit int
	// AutoAck if true (default true), automatically acknowledges every message so it will not be redelivered.
	// If false specifies that each message need ts to be manually acknowledged by the subscriber.
	// If processing is successful the message should be ack'ed to remove the message from the stream.
	// If processing is unsuccessful the message should be nack'ed (negative acknowledgement) which will mean it will
	// remain on the stream to be processed again.
	AutoAck bool
	// CustomRetries indicates whether to use RetryLimit
	CustomRetries bool
}

// GetRetryLimit returns the RetryLimit field on ConsumeOptions.
func (s ConsumeOptions) GetRetryLimit() int {
	if !s.CustomRetries {
		return -1
	}

	return s.RetryLimit
}

// ConsumeOption sets attributes on ConsumeOptions.
type ConsumeOption func(o *ConsumeOptions)

// WithGroup sets the consumer group to be part of when consuming events.
func WithGroup(q string) ConsumeOption {
	return func(o *ConsumeOptions) {
		o.Group = q
	}
}

// WithOffset sets the offset time at which to start consuming events.
func WithOffset(t time.Time) ConsumeOption {
	return func(o *ConsumeOptions) {
		o.Offset = t
	}
}

// WithAutoAck sets the AutoAck field on ConsumeOptions and an ackWait duration after which if no ack is received
// the message is requeued in case auto ack is turned off.
func WithAutoAck(ack bool, ackWait time.Duration) ConsumeOption {
	return func(o *ConsumeOptions) {
		o.AutoAck = ack
		o.AckWait = ackWait
	}
}

// WithRetryLimit sets the RetryLimit field on ConsumeOptions.
// Set to -1 for infinite retries (default).
func WithRetryLimit(retries int) ConsumeOption {
	return func(o *ConsumeOptions) {
		o.RetryLimit = retries
		o.CustomRetries = true
	}
}

// NewConsumeOptions generates new subscribe options with defaults.
func NewConsumeOptions(opts ...ConsumeOption) ConsumeOptions {
	cfg := ConsumeOptions{
		Group: uuid.New().String(),
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}
