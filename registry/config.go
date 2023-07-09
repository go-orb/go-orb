package registry

import (
	"errors"

	"github.com/go-orb/go-orb/config/source/cli"
	"github.com/go-orb/go-orb/log"
)

//nolint:gochecknoglobals
var (
	DefaultRegistry = "mdns"
	DefaultTimeout  = 100
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

func init() {
	// Register registry CLI flags.
	err := cli.Flags.Add(cli.NewFlag(
		"registry",
		DefaultRegistry,
		cli.ConfigPathSlice([]string{"registry", "plugin"}),
		cli.Usage("Registry for discovery. etcd, mdns"),
		cli.EnvVars("REGISTRY"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}

	err = cli.Flags.Add(cli.NewFlag(
		"registry_timeout",
		DefaultTimeout,
		cli.ConfigPathSlice([]string{"registry", "timeout"}),
		cli.Usage("Registry timeout in milliseconds."),
		cli.EnvVars("REGISTRY_TIMEOUT"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}
}

// TODO: this config misses stuff compared to v4, should that stuff be added here?

// Config is the configuration that can be used in a registry.
type Config struct {
	Plugin  string     `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	Timeout int        `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Logger  log.Logger `json:"logger,omitempty" yaml:"logger,omitempty"`
}

func (c *Config) config() *Config {
	return c
}

// ApplyOptions applies a set of options to the config.
func (c *Config) ApplyOptions(opts ...Option) {
	for _, o := range opts {
		o(c)
	}
}

// WithTimeout sets the default registry timeout used.
func WithTimeout(timeout int) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Timeout = timeout
	}
}

// WithLogger sets a specific logger to use.
func WithLogger(logger log.Logger) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Logger = logger
	}
}

// NewConfig creates a new default config to use with a registry.
func NewConfig() Config {
	return Config{
		Plugin:  DefaultRegistry,
		Timeout: DefaultTimeout,
	}
}
