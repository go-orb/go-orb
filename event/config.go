package event

import (
	"time"
)

//nolint:gochecknoglobals
var (
	// DefaultEventPlugin is the default client implementation to use.
	DefaultEventPlugin = "natsjs"

	// DefaultConfigSection is the default config section for the client.
	DefaultConfigSection = "client"

	// DefaultContentType is the default content type used to transport data around.
	DefaultContentType = "application/x-protobuf"

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
		ContentType:    DefaultContentType,
		Metadata:       make(map[string]string),
		RequestTimeout: DefaultRequestTimeout,
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}
