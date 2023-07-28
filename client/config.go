package client

import (
	"time"
)

//nolint:gochecknoglobals
var (
	// DefaultClient is the default client implementation to use.
	DefaultClientPlugin = "orb"

	// DefaultConfigSection is the default config section for the client.
	DefaultConfigSection = "client"

	// DefaultContentType is the default Content-Type for calls.
	DefaultContentType = "application/protobuf"
	// DefaultPreferredTransports set's in which order a transport will be selected.
	DefaultPreferredTransports = []string{"grpc", "h2c", "http2", "http3"}

	// DefaultPoolSize sets the connection pool size.
	DefaultPoolSize = 100
	// DefaultPoolTTL sets the connection pool ttl.
	DefaultPoolTTL = time.Minute

	// DefaultBackoff is the default backoff function for retries.
	DefaultBackoff = BackoffExponential
	// DefaultRetry is the default check-for-retry function for retries.
	DefaultRetry = RetryOnTimeoutError
	// DefaultRetries is the default number of times a request is tried.
	DefaultRetries = 5

	// DefaultDialTimeout is the default dial timeout.
	DefaultDialTimeout = time.Second * 5
	// DefaultRequestTimeout is the default request timeout.
	DefaultRequestTimeout = time.Second * 30
	// DefaultConnectionTimeout is the default connection timeout.
	DefaultConnectionTimeout = time.Second * 5
	// DefaultStreamTimeout is by default a noop.
	DefaultStreamTimeout = time.Duration(0)
	// DefaultConnClose indicates whetever to close the connection after each request.
	DefaultConnClose = false
	// DefaultReturnHeaders indicates if you want to copy resulting headers to the Response.
	DefaultReturnHeaders = false
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
//
// Plugin specific option example:
//
//		 // WithLogger option located in the MDNS registry package.
//			func WithLogger(logger log.Logger) registry.Option {
//			 	return func(c registry.ConfigType) {
//	        // The config type used here is *mdns.Config
//			 	   cfg, ok := c.(*Config)
//			 	   if ok {
//			 	    	cfg.Logger = logger
//			 	   }
//			 	}
//			}
type ConfigType interface {
	config() *Config
}

// Config are the Client options.
type Config struct {
	// Plugin selects the client implementation.
	Plugin string `json:"plugin" yaml:"plugin"`

	// Used to select a codec
	ContentType string `json:"contentType" yaml:"contentType"`

	PreferredTransports []string `json:"preferredTransports" yaml:"preferredTransports"`

	// Connection Pool
	PoolSize int           `json:"poolSize" yaml:"poolSize"`
	PoolTTL  time.Duration `json:"poolTTL"  yaml:"poolTTL"` //nolint:tagliatelle

	// Backoff func
	Backoff BackoffFunc `json:"-" yaml:"-"`
	// Check if retriable func
	Retry RetryFunc `json:"-" yaml:"-"`
	// Number of Call attempts
	Retries int `json:"retries" yaml:"retries"`
	// Transport Dial Timeout. Used for initial dial to establish a connection.
	DialTimeout time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	// ConnectionTimeout of one request to the server.
	// Set this lower than the RequestTimeout to enbale retries on connection timeout.
	ConnectionTimeout time.Duration `json:"connectionTimeout" yaml:"connectionTimeout"`
	// Request/Response timeout of entire srv.Call, for single request timeout set ConnectionTimeout.
	RequestTimeout time.Duration `json:"requestTimeout" yaml:"requestTimeout"`
	// Stream timeout for the stream
	StreamTimeout time.Duration `json:"streamTimeout" yaml:"streamTimeout"`
	// ReturnHeaders set to true will add Headers to the response
	ReturnHeaders bool `json:"returnHeaders" yaml:"returnHeaders"`
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

// WithClientContentType set's the Content-Type other than the default
// for this client.
func WithClientContentType(n string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.ContentType = n
	}
}

// WithClientPreferredTransports set the order of transports.
func WithClientPreferredTransports(n []string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.PreferredTransports = n
	}
}

// WithClientPoolSize overrides the PoolSize of the client.
func WithClientPoolSize(n int) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.PoolSize = n
	}
}

// WithClientPoolTTL overrides the PoolTTL of the client.
func WithClientPoolTTL(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.PoolTTL = n
	}
}

// WithClientBackoff overrides the backoff func.
func WithClientBackoff(n BackoffFunc) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Backoff = n
	}
}

// WithClientRetry overrides the retry function.
func WithClientRetry(n RetryFunc) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Retry = n
	}
}

// WithClientRetries overrides the number of retries to make.
func WithClientRetries(n int) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Retries = n
	}
}

// WithClientDialTimeout overrides the dial timeout.
func WithClientDialTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.DialTimeout = n
	}
}

// WithClientConnectionTimeout overrides the connection timeout.
func WithClientConnectionTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.ConnectionTimeout = n
	}
}

// WithClientRequestTimeout overrides the request timeout.
func WithClientRequestTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.RequestTimeout = n
	}
}

// WithClientStreamTimeout overrides the stream timeout.
func WithClientStreamTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.StreamTimeout = n
	}
}

// NewConfig generates a new config with all the defaults.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		Plugin:              DefaultClientPlugin,
		ContentType:         DefaultContentType,
		PreferredTransports: DefaultPreferredTransports,
		PoolSize:            DefaultPoolSize,
		PoolTTL:             DefaultPoolTTL,
		Retries:             DefaultRetries,
		DialTimeout:         DefaultDialTimeout,
		ConnectionTimeout:   DefaultConnectionTimeout,
		RequestTimeout:      DefaultRequestTimeout,
		StreamTimeout:       DefaultStreamTimeout,
		ReturnHeaders:       DefaultReturnHeaders,
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// CallOptions are options used to make calls to a server.
type CallOptions struct {
	// Used to select a codec
	ContentType string

	// Backoff func
	Backoff BackoffFunc
	// Check if retriable func
	Retry RetryFunc
	// Number of Call attempts
	Retries int
	// Transport Dial Timeout. Used for initial dial to establish a connection.
	DialTimeout time.Duration
	// ConnectionTimeout of one request to the server.
	// Set this lower than the RequestTimeout to enbale retries on connection timeout.
	ConnectionTimeout time.Duration
	// Request/Response timeout of entire srv.Call, for single request timeout set ConnectionTimeout.
	RequestTimeout time.Duration
	// Stream timeout for the stream
	StreamTimeout time.Duration
	// ConnClose sets the Connection: close header.
	ConnClose bool
	// Headers copies all headers into the RawResponse.
	Headers bool
	// URL bypasses the registry when set. This is mainly for tests.
	// Only <scheme>://<host:port> will be used from it.
	URL string
}

// CallOption used by Call or Stream.
type CallOption func(*CallOptions)

// Call Options.

// WithContentType set's the call's Content-Type.
func WithContentType(ct string) CallOption {
	return func(o *CallOptions) {
		o.ContentType = ct
	}
}

// func WithSelectOption(so ...selector.SelectOption) CallOption {
// return func(o *CallOptions) {
// o.SelectOptions = append(o.SelectOptions, so...)
// }
// }

// WithBackoff is a CallOption which overrides that which
// set in Options.CallOptions.
func WithBackoff(fn BackoffFunc) CallOption {
	return func(o *CallOptions) {
		o.Backoff = fn
	}
}

// WithRetry is a CallOption which overrides that which
// set in Options.CallOptions.
func WithRetry(fn RetryFunc) CallOption {
	return func(o *CallOptions) {
		o.Retry = fn
	}
}

// WithRetries sets the number of tries for a call.
// This CallOption overrides Options.CallOptions.
func WithRetries(i int) CallOption {
	return func(o *CallOptions) {
		o.Retries = i
	}
}

// WithRequestTimeout is a CallOption which overrides that which
// set in Options.CallOptions.
func WithRequestTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.RequestTimeout = d
	}
}

// WithConnClose sets the Connection header to close.
func WithConnClose() CallOption {
	return func(o *CallOptions) {
		o.ConnClose = true
	}
}

// WithStreamTimeout sets the stream timeout.
func WithStreamTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.StreamTimeout = d
	}
}

// WithDialTimeout is a CallOption which overrides that which
// set in Options.CallOptions.
func WithDialTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.DialTimeout = d
	}
}

// WithHeaders copies all headers into the metadata of the context.
func WithHeaders() CallOption {
	return func(o *CallOptions) {
		o.Headers = true
	}
}

// WithURL bypasses the registry when set.
// This is mainly for tests.
// Only <scheme>://<host:port> will be used from it.
func WithURL(n string) CallOption {
	return func(o *CallOptions) {
		o.URL = n
	}
}