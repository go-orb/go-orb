package registry

import (
	"github.com/hashicorp/go-multierror"
	"jochum.dev/orb/orb/config/chelp"
	"jochum.dev/orb/orb/log"
)

const CONFIG_KEY_LOGGER = "logger"
const CONFIG_KEY_ADDRESSES = "addresses"
const CONFIG_KEY_TIMEOUT = "timeout"

type Config interface {
	chelp.PluginConfig

	// Optional
	Logger() log.Config
	Addresses() []string

	// Timeout in milliseconds.
	Timeout() int

	// Setters
	SetLogger(n log.Config)
	SetAddresses(n []string)
	SetTimeout(n int)
}

type BaseConfig struct {
	*chelp.BasePluginConfig

	logger    log.Config
	addresses []string
	timeout   int
}

func NewConfig() Config {
	return &BaseConfig{
		BasePluginConfig: chelp.NewBasePluginConfig(),
		logger:           log.NewConfig(),
	}
}

func (c *BaseConfig) Load(m map[string]any) error {
	var result error

	// Required
	if err := c.BasePluginConfig.Load(m); err != nil {
		result = multierror.Append(err)
	}
	var err error

	// Optional
	if err := c.logger.Load(m); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.addresses, err = chelp.Get(m, CONFIG_KEY_ADDRESSES, c.addresses); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.timeout, err = chelp.Get(m, CONFIG_KEY_TIMEOUT, c.timeout); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	return result
}

func (c *BaseConfig) Store(m map[string]any) error {
	var result error

	if err := c.BasePluginConfig.Store(m); err != nil {
		result = multierror.Append(err)
	}

	logger := make(map[string]any)
	if err := c.logger.Store(logger); err != nil {
		result = multierror.Append(err)
	}
	m[CONFIG_KEY_LOGGER] = logger

	m[CONFIG_KEY_ADDRESSES] = c.addresses
	m[CONFIG_KEY_TIMEOUT] = c.timeout

	return result
}

func (c *BaseConfig) Logger() log.Config  { return c.logger }
func (c *BaseConfig) Addresses() []string { return c.addresses }
func (c *BaseConfig) Timeout() int        { return c.timeout }

func (c *BaseConfig) SetLogger(n log.Config)  { c.logger = n }
func (c *BaseConfig) SetAddresses(n []string) { c.addresses = n }
func (c *BaseConfig) SetTimeout(n int)        { c.timeout = n }
