package kvstore

//nolint:gochecknoglobals
var (
	// DefaultConfigSection is the section in the config to use.
	DefaultConfigSection = "kvstore"

	// DefaultKVStore is the default kvstore to use.
	DefaultKVStore = "natsjs"
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

	// Database allows multiple isolated stores to be kept in one backend, if supported.
	Database string `json:"database,omitempty" yaml:"database,omitempty"`

	// Table is analogous to a table in database backends or a key prefix in KV backends
	Table string `json:"table,omitempty" yaml:"table,omitempty"`
}

// config returns the config.
func (c *Config) config() *Config {
	return c
}

// NewConfig creates a config to use with a registry.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		Plugin: DefaultKVStore,
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// WithPlugin set the client implementation to use.
func WithPlugin(n string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Plugin = n
	}
}

// WithDatabase allows multiple isolated stores to be kept in one backend, if supported.
func WithDatabase(db string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Database = db
	}
}

// WithTable is analogous to a table in database backends or a key prefix in KV backends.
func WithTable(t string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Table = t
	}
}
