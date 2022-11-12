package server

import (
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/types"
)

var (
	// DefaultConfigSection is the section key used in config files used to
	// configure the server options.
	DefaultConfigSection = "server"
)

// Option is a functional HTTP server option.
type Option func(*Config)

// Config is the server config. It contains the list of addresses on which
// entrypoints will be created, and the default config used for each entrypoint.
type Config struct {
	// Defaults is the list of defaults for a server.
	// Provisioned with the factory methods registerd by the entrypoint plugins.
	Defaults map[string]any

	// Templates contains a set of entrypoint templates to create, indexed by name.
	Templates EntrypointTemplates
}

// NewConfig creates a new server config with default values as starting point,
// after which all the functional options are applied. The config data passed
// in, as parsed from the optional config files and CLI, has the highest priority.
func NewConfig(serviceName types.ServiceName, data types.ConfigData, options ...Option) (Config, error) {
	cfg := Config{
		Defaults:  make(map[string]any),
		Templates: make(EntrypointTemplates),
	}

	// Provision defaults from all entrypoints.
	factories := NewDefaults.All()
	for name, factory := range factories {
		cfg.Defaults[name] = factory()
	}

	cfg.ApplyOptions(options...)

	sections := types.SplitServiceName(serviceName)
	if err := config.Parse(append(sections, DefaultConfigSection), data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// ApplyOptions takes a list of options and applies them to the current config.
func (c *Config) ApplyOptions(options ...Option) {
	for _, option := range options {
		option(c)
	}
}
