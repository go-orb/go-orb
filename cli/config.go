package cli

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
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

// Config is the interface every plugin must implement.
type Config interface { // /nolint:interfacebloat
	chelp.Plugin

	// Required
	Name() string
	SetName(n string)

	Version() string
	SetVersion(n string)

	// Optional
	Description() string
	SetDescription(n string)

	Usage() string
	SetUsage(n string)

	NoFlags() *bool
	SetNoFlags(n *bool)

	ConfigSection() string
	SetConfigSection(n string)

	ArgPrefix() string
	SetArgPrefix(n string)

	Config() []string
	SetConfig(n []string)

	// Internal to transfer Options
	Flags() []Flag
	SetFlags(n []Flag)
}

// BaseConfig is the base config for this component.
type BaseConfig struct {
	chelp.Plugin

	// Required
	name    string
	version string

	// Optional
	description   string
	usage         string
	noFlags       *bool
	configSection string
	argPrefix     string
	config        []string

	// Internal
	flags []Flag
}

// NewConfig returns the base component config.
func NewConfig() Config {
	return &BaseConfig{
		Plugin: chelp.NewPluginConfig(),
	}
}

// defaultConfig returns the default config for a given Plugin.
func defaultConfig(name string) (any, error) {
	confFactory, err := Plugins.Config(name)
	if err != nil {
		return nil, err
	}

	return confFactory(), nil
}

func getConfig(m map[string]any) (any, error) {
	pconf := chelp.NewPluginConfig()
	if err := pconf.Load(m); err != nil {
		return nil, err
	}

	return defaultConfig(pconf.Plugin())
}

// LoadConfig loads the config from map `m` with the key `key`.
func LoadConfig(m map[string]any, key string) (any, error) {
	// Optional
	myMap, err := chelp.Get(m, key, map[string]any{})
	if err != nil {
		return nil, err
	}

	myConf, err := getConfig(myMap)
	if err != nil {
		return nil, err
	}

	if loader, ok := myConf.(chelp.ConfigMethods); ok {
		if err := loader.Load(myMap); err != nil {
			return nil, err
		}
	} else {
		return nil, chelp.ErrUnknownConfig
	}

	return myConf, nil
}

// StoreConfig stores the config to map[string]any.
func StoreConfig(config any) (map[string]any, error) {
	result := make(map[string]any)
	if config == nil {
		return result, chelp.ErrNotExistant
	}

	if storer, ok := config.(chelp.ConfigMethods); ok {
		if err := storer.Store(result); err != nil {
			return result, err
		}
	} else {
		return result, chelp.ErrUnknownConfig
	}

	return result, nil
}

// Load loads this config from map[string]any.
func (c *BaseConfig) Load(inputMap map[string]any) error {
	var result error

	// Required
	if err := c.Plugin.Load(inputMap); err != nil {
		result = multierror.Append(result, err)
	}

	var err error
	if c.name, err = chelp.Get(inputMap, configKeyName, c.name); err != nil {
		result = multierror.Append(result, err)
	}

	if c.version, err = chelp.Get(inputMap, configKeyVersion, c.version); err != nil {
		result = multierror.Append(result, err)
	}

	// Optional
	c.description, err = chelp.Get(inputMap, configKeyDescription, c.description)
	if err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
	}

	if c.usage, err = chelp.Get(inputMap, configKeyUsage, c.usage); !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
	}

	if c.noFlags, err = chelp.Get(inputMap, configKeyNoFlags, c.noFlags); !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
	}

	c.argPrefix, err = chelp.Get(inputMap, configKeyArgPrefix, c.argPrefix)
	if err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
	}

	if c.config, err = chelp.Get(inputMap, configKeyConfig, c.config); !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
	}

	return result
}

// Store stores this config in a map[string]any.
func (c *BaseConfig) Store(m map[string]any) error {
	var result error

	if err := c.Plugin.Store(m); err != nil {
		result = multierror.Append(result, err)
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

// Merge merges aConfig into this.
func (c *BaseConfig) Merge(aConfig any) error {
	if err := c.Plugin.Merge(aConfig); err != nil {
		return err
	}

	toMerge, ok := aConfig.(Config)
	if !ok {
		return chelp.ErrUnknownConfig
	}

	defConfig := NewConfig()

	if toMerge.Name() != defConfig.Name() {
		c.SetName(toMerge.Name())
	}

	if toMerge.Version() != defConfig.Version() {
		c.SetVersion(toMerge.Version())
	}

	if toMerge.Description() != defConfig.Description() {
		c.SetDescription(toMerge.Description())
	}

	if toMerge.Usage() != defConfig.Usage() {
		c.SetUsage(toMerge.Usage())
	}

	if toMerge.NoFlags() != defConfig.NoFlags() {
		c.SetNoFlags(toMerge.NoFlags())
	}

	if toMerge.ArgPrefix() != defConfig.ArgPrefix() {
		c.SetArgPrefix(toMerge.ArgPrefix())
	}

	if slices.Compare(toMerge.Config(), defConfig.Config()) != 0 {
		c.SetConfig(toMerge.Config())
	}

	return nil
}

// Name returns the apps name.
func (c *BaseConfig) Name() string { return c.name }

// Version returns the apps version.
func (c *BaseConfig) Version() string { return c.version }

// Description returns the apps description.
func (c *BaseConfig) Description() string { return c.description }

// Usage returns the apps usage string.
func (c *BaseConfig) Usage() string { return c.usage }

// NoFlags indicates if we want to disable Flags.
func (c *BaseConfig) NoFlags() *bool { return c.noFlags }

// ConfigSection returns the section for orb internal components.
func (c *BaseConfig) ConfigSection() string { return c.configSection }

// ArgPrefix sets the prefix for orb internal Flags.
func (c *BaseConfig) ArgPrefix() string { return c.argPrefix }

// Config sets the config urls.
func (c *BaseConfig) Config() []string { return c.config }

func (c *BaseConfig) SetName(n string)          { c.name = n }
func (c *BaseConfig) SetVersion(n string)       { c.version = n }
func (c *BaseConfig) SetDescription(n string)   { c.description = n }
func (c *BaseConfig) SetUsage(n string)         { c.usage = n }
func (c *BaseConfig) SetNoFlags(n *bool)        { c.noFlags = n }
func (c *BaseConfig) SetConfigSection(n string) { c.configSection = n }
func (c *BaseConfig) SetArgPrefix(n string)     { c.argPrefix = n }
func (c *BaseConfig) SetConfig(n []string)      { c.config = n }

func (c *BaseConfig) Flags() []Flag     { return c.flags }
func (c *BaseConfig) SetFlags(n []Flag) { c.flags = n }
