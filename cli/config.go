package cli

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"jochum.dev/orb/orb/config/chelp"
)

const (
	configKeyName        = "name"
	configKeyVersion     = "version"
	configKeyDescription = "description"
	configKeyUsage       = "usage"
	configKeyNoFlags     = "no_flags"
	configKeyArgPrefix   = "arg_prefix"
	configKeyConfig      = "config"
)

type Config interface {
	chelp.PluginConfig

	// Required
	Name() string
	Version() string

	// Optional
	Description() string
	Usage() string
	NoFlags() *bool
	ArgPrefix() string
	Config() string

	// Setters for hardcoded settings
	SetName(n string)
	SetVersion(n string)
	SetDescription(n string)
	SetUsage(n string)
	SetNoFlags(n *bool)
	SetArgPrefix(n string)
	SetConfig(n string)

	// Internal to transfer Options
	Flags() []Flag
	SetFlags(n []Flag)
}

type BaseConfig struct {
	chelp.PluginConfig

	// Required
	name    string
	version string

	// Optional
	description string
	usage       string
	noFlags     *bool
	argPrefix   string
	config      string

	// Internal
	flags []Flag
}

func NewConfig() *BaseConfig {
	return &BaseConfig{
		PluginConfig: chelp.NewPluginConfig(),
	}
}

func (c *BaseConfig) Load(m map[string]any) error {
	var result error

	// Required
	if err := c.PluginConfig.Load(m); err != nil {
		result = multierror.Append(err)
	}

	var err error
	if c.name, err = chelp.Get(m, configKeyName, c.name); err != nil {
		result = multierror.Append(err)
	}

	if c.version, err = chelp.Get(m, configKeyVersion, c.version); err != nil {
		result = multierror.Append(err)
	}

	// Optional
	c.description, err = chelp.Get(m, configKeyDescription, c.description)
	if err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	if c.usage, err = chelp.Get(m, configKeyUsage, c.usage); err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	if c.noFlags, err = chelp.Get(m, configKeyNoFlags, c.noFlags); err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	c.argPrefix, err = chelp.Get(m, configKeyArgPrefix, c.argPrefix)
	if err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	if c.config, err = chelp.Get(m, configKeyConfig, c.config); err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	return result
}

func (c *BaseConfig) Store(m map[string]any) error {
	var result error

	if err := c.PluginConfig.Store(m); err != nil {
		result = multierror.Append(err)
	}

	m[configKeyName] = c.name
	m[configKeyVersion] = c.version

	m[configKeyDescription] = c.description
	m[configKeyUsage] = c.usage
	m[configKeyNoFlags] = c.noFlags
	m[configKeyArgPrefix] = c.argPrefix
	m[configKeyConfig] = c.config

	return result
}

func (c *BaseConfig) Name() string        { return c.name }
func (c *BaseConfig) Version() string     { return c.version }
func (c *BaseConfig) Description() string { return c.description }
func (c *BaseConfig) Usage() string       { return c.usage }
func (c *BaseConfig) NoFlags() *bool      { return c.noFlags }
func (c *BaseConfig) ArgPrefix() string   { return c.argPrefix }
func (c *BaseConfig) Config() string      { return c.config }

func (c *BaseConfig) SetName(n string)        { c.name = n }
func (c *BaseConfig) SetVersion(n string)     { c.version = n }
func (c *BaseConfig) SetDescription(n string) { c.description = n }
func (c *BaseConfig) SetUsage(n string)       { c.usage = n }
func (c *BaseConfig) SetNoFlags(n *bool)      { c.noFlags = n }
func (c *BaseConfig) SetArgPrefix(n string)   { c.argPrefix = n }
func (c *BaseConfig) SetConfig(n string)      { c.config = n }

func (c *BaseConfig) Flags() []Flag     { return c.flags }
func (c *BaseConfig) SetFlags(n []Flag) { c.flags = n }
