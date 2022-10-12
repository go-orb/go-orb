package cli

import (
	"github.com/hashicorp/go-multierror"
	"jochum.dev/orb/orb/config/chelp"
)

const (
	CONFIG_KEY_NAME        = "name"
	CONFIG_KEY_VERSION     = "version"
	CONFIG_KEY_DESCRIPTION = "description"
	CONFIG_KEY_USAGE       = "usage"
	CONFIG_KEY_NO_FLAGS    = "no_flags"
	CONFIG_KEY_ARG_PREFIX  = "arg_prefix"
	CONFIG_KEY_CONFIG      = "config"
)

type Config interface {
	chelp.PluginConfig

	Name() string
	Version() string

	Description() string
	Usage() string
	NoFlags() bool
	ArgPrefix() string
	Config() string
}

type BaseConfig struct {
	*chelp.BasicPlugin

	// Required
	name    string
	version string

	// Optional
	description string
	usage       string
	noFlags     *bool
	argPrefix   string
	config      string

	// No config option but here as helper
	Flags []Flag
}

func NewConfig() Config {
	return &BaseConfig{
		BasicPlugin: chelp.NewBasicPlugin(),
	}
}

func (c *BaseConfig) Load(m map[string]any) error {
	var result error

	// Required
	if err := c.BasicPlugin.Load(m); err != nil {
		result = multierror.Append(err)
	}
	var err error
	if c.name, err = chelp.Get(m, CONFIG_KEY_NAME, c.name); err != nil {
		result = multierror.Append(err)
	}
	if c.version, err = chelp.Get(m, CONFIG_KEY_VERSION, c.version); err != nil {
		result = multierror.Append(err)
	}

	// Optional
	if c.description, err = chelp.Get(m, CONFIG_KEY_DESCRIPTION, c.description); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.usage, err = chelp.Get(m, CONFIG_KEY_USAGE, c.usage); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.noFlags, err = chelp.Get(m, CONFIG_KEY_NO_FLAGS, c.noFlags); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.argPrefix, err = chelp.Get(m, CONFIG_KEY_ARG_PREFIX, c.argPrefix); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.config, err = chelp.Get(m, CONFIG_KEY_CONFIG, c.config); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	return result
}

func (c *BaseConfig) Store(m map[string]any) error {
	var result error

	if err := c.BasicPlugin.Store(m); err != nil {
		result = multierror.Append(err)
	}

	m[CONFIG_KEY_NAME] = c.name
	m[CONFIG_KEY_VERSION] = c.version

	m[CONFIG_KEY_DESCRIPTION] = c.description
	m[CONFIG_KEY_USAGE] = c.usage
	m[CONFIG_KEY_NO_FLAGS] = c.noFlags
	m[CONFIG_KEY_ARG_PREFIX] = c.argPrefix
	m[CONFIG_KEY_CONFIG] = c.config

	return result
}

func (c *BaseConfig) Name() string        { return c.name }
func (c *BaseConfig) Version() string     { return c.version }
func (c *BaseConfig) Description() string { return c.description }
func (c *BaseConfig) Usage() string       { return c.usage }
func (c *BaseConfig) NoFlags() bool       { return *c.noFlags }
func (c *BaseConfig) ArgPrefix() string   { return c.argPrefix }
func (c *BaseConfig) Config() string      { return c.config }
