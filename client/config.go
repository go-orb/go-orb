package client

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-orb/go-orb/config"
)

//nolint:gochecknoglobals
var (
	// DefaultClientPlugin is the default client implementation to use.
	DefaultClientPlugin = "orb"

	// DefaultConfigSection is the default config section for the client.
	DefaultConfigSection = "client"

	// DefaultContentType is the default Content-Type for calls.
	DefaultContentType = "application/x-protobuf"
	// DefaultPreferredTransports set's in which order a transport will be selected.
	DefaultPreferredTransports = []string{"memory", "grpc", "drpc", "http", "grpcs", "h2c", "http2", "http3", "https"}

	// DefaultPoolHosts set the number of hosts in a pool.
	DefaultPoolHosts = 64
	// DefaultPoolSize sets the connection pool size.
	// The effective pool size will be PoolHosts * PoolSize.
	DefaultPoolSize = 10
	// DefaultPoolTTL sets the connection pool ttl.
	DefaultPoolTTL = 30 * time.Minute

	// DefaultSelector is the default node selector.
	DefaultSelector = SelectRandomNode

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

	// DefaultCallOptionsRetryFunc is nil, so it uses the middlewares default.
	DefaultCallOptionsRetryFunc = (RetryFunc)(nil)

	// DefaultCallOptionsRetries is 0, so it uses the middlewares default.
	DefaultCallOptionsRetries = 0

	// DefaultMaxCallRecvMsgSize is the default maximum size of the call receive message size.
	DefaultMaxCallRecvMsgSize = 10 * 1024 * 1024
	// DefaultMaxCallSendMsgSize is the default maximum size of the call send message size.
	DefaultMaxCallSendMsgSize = 10 * 1024 * 1024
)

// RetryFunc is the type for a retry func.
// note that returning either false or a non-nil error will result in the call not being retried.
type RetryFunc func(ctx context.Context, err error, options *CallOptions) (bool, error)

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

	Middleware []MiddlewareConfig

	// Used to select a codec
	ContentType string `json:"contentType" yaml:"contentType"`

	// PreferredTransports contains a list of transport names in preferred order.
	PreferredTransports []string `json:"preferredTransports" yaml:"preferredTransports"`

	// AnyTransport enables Transports which are not in PreferredTransports.
	AnyTransport bool `json:"anyTransport" yaml:"anyTransport"`

	// Connection Pool
	PoolHosts int             `json:"poolHosts" yaml:"poolHosts"`
	PoolSize  int             `json:"poolSize"  yaml:"poolSize"`
	PoolTTL   config.Duration `json:"poolTtl"   yaml:"poolTtl"`

	// SelectorFunc get's executed by client.SelectNode which get it's info's from client.ResolveService.
	Selector SelectorFunc `json:"-" yaml:"-"`

	// Transport Dial Timeout. Used for initial dial to establish a connection.
	DialTimeout config.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	// ConnectionTimeout of one request to the server.
	// Set this lower than the RequestTimeout to enbale retries on connection timeout.
	ConnectionTimeout config.Duration `json:"connectionTimeout" yaml:"connectionTimeout"`
	// Request/Response timeout of entire srv.Call, for single request timeout set ConnectionTimeout.
	RequestTimeout config.Duration `json:"requestTimeout" yaml:"requestTimeout"`
	// Stream timeout for the stream
	StreamTimeout config.Duration `json:"streamTimeout" yaml:"streamTimeout"`
	// TLS config.
	TLSConfig *tls.Config
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
func WithClientPreferredTransports(n ...string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.PreferredTransports = n
	}
}

// WithClientAnyTransport enables Transports which are not in PreferredTransports.
func WithClientAnyTransport() Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.AnyTransport = true
	}
}

// WithClientPoolHosts overrides the PoolHosts of the client.
func WithClientPoolHosts(n int) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.PoolHosts = n
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
		c.PoolTTL = config.Duration(n)
	}
}

// WithClientSelector overrides the clients selector func.
func WithClientSelector(n SelectorFunc) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Selector = n
	}
}

// WithClientDialTimeout overrides the dial timeout.
func WithClientDialTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.DialTimeout = config.Duration(n)
	}
}

// WithClientConnectionTimeout overrides the connection timeout.
func WithClientConnectionTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.ConnectionTimeout = config.Duration(n)
	}
}

// WithClientRequestTimeout overrides the request timeout.
func WithClientRequestTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.RequestTimeout = config.Duration(n)
	}
}

// WithClientStreamTimeout overrides the stream timeout.
func WithClientStreamTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.StreamTimeout = config.Duration(n)
	}
}

// WithClientTLSConfig set's the clients TLS config.
func WithClientTLSConfig(n *tls.Config) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.TLSConfig = n
	}
}

// WithClientMiddleware appends a middleware to the client.
func WithClientMiddleware(m MiddlewareConfig) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Middleware = append(c.Middleware, m)
	}
}

// NewConfig generates a new config with all the defaults.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		Plugin:              DefaultClientPlugin,
		ContentType:         DefaultContentType,
		PreferredTransports: DefaultPreferredTransports,
		PoolHosts:           DefaultPoolHosts,
		PoolSize:            DefaultPoolSize,
		PoolTTL:             config.Duration(DefaultPoolTTL),
		DialTimeout:         config.Duration(DefaultDialTimeout),
		ConnectionTimeout:   config.Duration(DefaultConnectionTimeout),
		RequestTimeout:      config.Duration(DefaultRequestTimeout),
		StreamTimeout:       config.Duration(DefaultStreamTimeout),
		Selector:            DefaultSelector,
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

	// PreferredTransports contains a list of transport names in preferred order.
	PreferredTransports []string

	AnyTransport bool

	// Selector is the node selector.
	Selector SelectorFunc
	// Check if retriable func
	RetryFunc RetryFunc
	// Number of Call attempts
	Retries int
	// Transport Dial Timeout. Used for initial dial to establish a connection.
	DialTimeout time.Duration
	// ConnectionTimeout of one request to the server.
	// Set this lower than the RequestTimeout to enable retries on connection timeout.
	ConnectionTimeout time.Duration
	// Request/Response timeout of entire srv.Call, for single request timeout set ConnectionTimeout.
	RequestTimeout time.Duration
	// Stream timeout for the stream
	StreamTimeout time.Duration
	// ConnClose sets the Connection: close header.
	ConnClose bool
	// URL bypasses the registry when set. This is mainly for tests.
	// Only <scheme>://<host:port> will be used from it.
	URL string
	// TLS config.
	TLSConfig *tls.Config

	// Metadata to be sent with the request.
	Metadata map[string]string

	// ResponseMetadata will be written into `ResponseMetadata` when given.
	ResponseMetadata map[string]string

	// MaxCallRecvMsgSize is the maximum size of the call receive message size.
	MaxCallRecvMsgSize int

	// MaxCallSendMsgSize is the maximum size of the call send message size.
	MaxCallSendMsgSize int
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

// WithPreferredTransports set's the preffered transports for this request.
func WithPreferredTransports(n ...string) CallOption {
	return func(o *CallOptions) {
		o.PreferredTransports = n
	}
}

// WithAnyTransport enables unconfigured (any) transports.
func WithAnyTransport() CallOption {
	return func(o *CallOptions) {
		o.AnyTransport = true
	}
}

// WithSelector overrides the calls SelectorFunc.
func WithSelector(fn SelectorFunc) CallOption {
	return func(o *CallOptions) {
		o.Selector = fn
	}
}

// WithRetryFunc is a CallOption which overrides the retry function.
func WithRetryFunc(fn RetryFunc) CallOption {
	return func(o *CallOptions) {
		o.RetryFunc = fn
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

// WithURL bypasses the registry when set.
// This is mainly for tests.
// Only <scheme>://<host:port> will be used from it.
func WithURL(n string) CallOption {
	return func(o *CallOptions) {
		o.URL = n
	}
}

// WithTLSConfig set's the clients TLS config.
func WithTLSConfig(n *tls.Config) CallOption {
	return func(o *CallOptions) {
		o.TLSConfig = n
	}
}

// WithMetadata sets the metadata to be sent with the request.
func WithMetadata(n map[string]string) CallOption {
	return func(o *CallOptions) {
		o.Metadata = n
	}
}

// WithRegistryMetadata adds a metadata key-value pair for node selection
func WithRegistryMetadata(key, value string) CallOption {
	return func(o *CallOptions) {
		if o.Metadata == nil {
			o.Metadata = make(map[string]string)
		}
		o.Metadata[key] = value
	}
}

// WithResponseMetadata will write response Metadata into the give map.
func WithResponseMetadata(n map[string]string) CallOption {
	return func(o *CallOptions) {
		o.ResponseMetadata = n
	}
}

// WithMaxCallRecvMsgSize sets the maximum size of the call receive message size.
func WithMaxCallRecvMsgSize(n int) CallOption {
	return func(o *CallOptions) {
		o.MaxCallRecvMsgSize = n
	}
}

// WithMaxCallSendMsgSize sets the maximum size of the call send message size.
func WithMaxCallSendMsgSize(n int) CallOption {
	return func(o *CallOptions) {
		o.MaxCallSendMsgSize = n
	}
}
