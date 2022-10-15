package cli

import (
	"golang.org/x/exp/slices"
	"jochum.dev/orb/orb/config"
)

// Config is the interface every plugin must implement.
type Config interface { // /nolint:interfacebloat
	config.ConfigPlugin

	// Required.
	GetName() string
	GetVersion() string

	// Optional.
	GetDescription() string
	GetUsage() string
	GetNoFlags() *bool
	GetConfigSection() string
	GetArgPrefix() string

	GetConfig() []string
	SetConfig(n []string)

	// Internal to transfer custom Flags.
	GetFlags() []Flag
}

// ConfigImpl is the base config for this component.
type ConfigImpl struct {
	*config.ConfigPluginImpl

	// Required
	Name    string `json:"name" yaml:"Name"`
	Version string `json:"version" yaml:"Version"`

	// Optional
	Description   string   `json:"description" yaml:"Description"`
	Usage         string   `json:"usage" yaml:"Usage"`
	NoFlags       *bool    `json:"no_flags" yaml:"NoFlags"`
	ConfigSection string   `json:"config_section" yaml:"ConfigSection"`
	ArgPrefix     string   `json:"arg_prefix" yaml:"ArgPrefix"`
	Config        []string `json:"config" yaml:"Config"`

	// Internal
	Flags []Flag `json:"-" yaml:"-"`
}

// NewComponentConfig returns the base component config.
func NewComponentConfig() *ConfigImpl {
	return &ConfigImpl{
		ConfigPluginImpl: config.NewPluginConfig(),
	}
}

// NewConfig returns the default config for a given Plugin.
func NewConfig(name string) (any, error) {
	confFactory, err := Plugins.Config(name)
	if err != nil {
		return nil, err
	}

	return confFactory(), nil
}

// MergePrevious merges the previous config into this one.
func (c *ConfigImpl) MergePrevious(aPreviousConfig any) error {
	if err := c.ConfigPluginImpl.MergePrevious(aPreviousConfig); err != nil {
		return err
	}

	previous, ok := aPreviousConfig.(Config)
	if !ok {
		return config.ErrUnknownConfig
	}

	defConfig := NewComponentConfig()

	if c.Name == defConfig.Name {
		c.Name = previous.GetName()
	}

	if c.Version == defConfig.Version {
		c.Version = previous.GetVersion()
	}

	if c.Description == defConfig.Description {
		c.Description = previous.GetDescription()
	}

	if c.Usage == defConfig.Usage {
		c.Usage = previous.GetUsage()
	}

	if c.NoFlags == defConfig.NoFlags {
		c.NoFlags = previous.GetNoFlags()
	}

	if c.ArgPrefix == defConfig.ArgPrefix {
		c.ArgPrefix = previous.GetArgPrefix()
	}

	if slices.Compare(c.Config, defConfig.Config) != 0 {
		c.Config = previous.GetConfig()
	}

	return nil
}

// GetName returns the apps name.
func (c *ConfigImpl) GetName() string { return c.Name }

// GetVersion returns the apps version.
func (c *ConfigImpl) GetVersion() string { return c.Version }

// GetDescription returns the apps description.
func (c *ConfigImpl) GetDescription() string { return c.Description }

// GetUsage returns the apps usage string.
func (c *ConfigImpl) GetUsage() string { return c.Usage }

// GetNoFlags indicates if we want to disable Flags.
func (c *ConfigImpl) GetNoFlags() *bool { return c.NoFlags }

// GetConfigSection returns the section for orb internal components.
func (c *ConfigImpl) GetConfigSection() string { return c.ConfigSection }

// GetArgPrefix sets the prefix for orb internal Flags.
func (c *ConfigImpl) GetArgPrefix() string { return c.ArgPrefix }

// GetConfig sets the config urls.
func (c *ConfigImpl) GetConfig() []string  { return c.Config }
func (c *ConfigImpl) SetConfig(n []string) { c.Config = n }

// GetFlags returns the custom flags.
func (c *ConfigImpl) GetFlags() []Flag { return c.Flags }
