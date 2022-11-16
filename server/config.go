package server

import (
	"fmt"

	"go-micro.dev/v5/types"
)

//nolint:gochecknoglobals
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
	// Defaults is the list of defaults the each server plugin.
	// Provisioned with the factory methods registered by the entrypoint plugins.
	Defaults map[string]any

	// Templates contains a set of entrypoint templates to create, indexed by name.
	//
	// Each entrypoint needs a unique name, as each entrypoint can be dynamically
	// configured by referencing the name. The default name used in an entrypoint
	// is the format of "http-<uuid>", used if no custom name is provided.
	Templates map[string]EntrypointTemplate
}

// NewConfig creates a new server config with default values as starting point,
// after which all the functional options are applied. The config data passed
// in, as parsed from the optional config files and CLI, has the highest priority.
func NewConfig(service types.ServiceName, data types.ConfigData, options ...Option) (Config, error) {
	cfg := Config{
		Defaults:  make(map[string]any),
		Templates: make(map[string]EntrypointTemplate),
	}

	var err error

	// Provision defaults for all entrypoints. Factories are provided by the plugins.
	factories := NewDefaults.All()
	for name, factory := range factories {
		cfg.Defaults[name], err = factory(service, data)
		if err != nil {
			return cfg, fmt.Errorf("create %s default config: %w", name, err)
		}
	}

	cfg.ApplyOptions(options...)

	return cfg, nil
}

// ApplyOptions takes a list of options and applies them to the current config.
func (c *Config) ApplyOptions(options ...Option) {
	for _, option := range options {
		option(c)
	}
}
