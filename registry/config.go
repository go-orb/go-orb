package registry

import (
	"golang.org/x/exp/slices"
	"jochum.dev/orb/orb/config"
	"jochum.dev/orb/orb/log"
)

type Config interface {
	config.ConfigPlugin

	// Optional
	GetLogger() *log.ConfigImpl
	GetAddresses() []string

	// Timeout in milliseconds.
	GetTimeout() int
}

type ConfigImpl struct {
	*config.ConfigPluginImpl

	Logger    *log.ConfigImpl
	Addresses []string
	Timeout   int
}

func NewComponentConfig() *ConfigImpl {
	return &ConfigImpl{
		ConfigPluginImpl: config.NewPluginConfig(),
		Logger:           log.NewComponentConfig(),
	}
}

// NewConfig returns the default config for a given Plugin.
func NewConfig(plugin string) (any, error) {
	confFactory, err := Plugins.Config(plugin)
	if err != nil {
		return nil, err
	}

	return confFactory(), nil
}

// MergePrevious merges the previous config into this one.
func (c *ConfigImpl) Merge(aPreviousConfig any) error {
	if err := c.ConfigPluginImpl.MergePrevious(aPreviousConfig); err != nil {
		return err
	}

	previous, ok := aPreviousConfig.(Config)
	if !ok {
		return config.ErrUnknownConfig
	}

	if previous.GetLogger() != nil && c.Logger != nil {
		c.Logger.MergePrevious(previous.GetLogger())
	} else if c.Logger == nil {
		c.Logger = previous.GetLogger()
	}

	defConfig := NewComponentConfig()

	if slices.Compare(c.Addresses, defConfig.Addresses) == 0 {
		c.Addresses = previous.GetAddresses()
	}

	if c.Timeout == defConfig.Timeout {
		c.Timeout = previous.GetTimeout()
	}

	return nil
}

func (c *ConfigImpl) GetLogger() *log.ConfigImpl { return c.Logger }
func (c *ConfigImpl) GetAddresses() []string     { return c.Addresses }
func (c *ConfigImpl) GetTimeout() int            { return c.Timeout }
