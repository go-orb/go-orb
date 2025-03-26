package registry

import (
	"time"

	"github.com/go-orb/go-orb/config"
)

//nolint:gochecknoglobals
var (
	// DefaultConfigSection is the section in the config to use.
	DefaultConfigSection = "registry"

	// DefaultRegistry is the default registry to use.
	DefaultRegistry = "mdns"

	// DefaultTimeout is the default timeout for the registry.
	DefaultTimeout = config.Duration(500 * time.Millisecond)
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

// TODO(jochumdev): this config misses things compared to v4, should they be added here?

// Config is the configuration that can be used in a registry.
type Config struct {
	Plugin  string          `json:"plugin,omitempty"  yaml:"plugin,omitempty"`
	Timeout config.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (c *Config) config() *Config {
	return c
}

// WithPlugin set the implementation to use.
func WithPlugin(n string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Plugin = n
	}
}

// WithTimeout sets the default registry timeout used.
func WithTimeout(n time.Duration) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Timeout = config.Duration(n)
	}
}

// NewConfig creates a config to use with a registry.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		Plugin:  DefaultRegistry,
		Timeout: DefaultTimeout,
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}
