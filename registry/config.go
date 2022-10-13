package registry

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"jochum.dev/orb/orb/config/chelp"
	"jochum.dev/orb/orb/log"
)

const configKeyLogger = "logger"
const configKeyAddresses = "addresses"
const configKeyTimeout = "timeout"

type Config interface {
	chelp.PluginConfig

	// Optional
	Logger() any
	Addresses() []string

	// Timeout in milliseconds.
	Timeout() int

	// Setters
	SetLogger(n any)
	SetAddresses(n []string)
	SetTimeout(n int)
}

type BaseConfig struct {
	*chelp.BasePluginConfig

	logger    any
	addresses []string
	timeout   int
}

func NewConfig() *BaseConfig {
	return &BaseConfig{
		BasePluginConfig: chelp.NewPluginConfig(),
		logger:           log.NewConfig(),
	}
}

func (c *BaseConfig) Load(m map[string]any) error {
	var result error

	// Required
	if err := c.BasePluginConfig.Load(m); err != nil {
		result = multierror.Append(err)
	}

	var (
		err error
	)

	// Optional
	c.logger, err = log.LoadConfig(m, configKeyLogger)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	c.addresses, err = chelp.Get(m, configKeyAddresses, c.addresses)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	c.timeout, err = chelp.Get(m, configKeyTimeout, c.timeout)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	return result
}

func (c *BaseConfig) Store(m map[string]any) error {
	var result error

	if err := c.BasePluginConfig.Store(m); err != nil {
		result = multierror.Append(err)
	}

	var err error

	m[configKeyLogger], err = log.StoreConfig(c.logger)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	m[configKeyAddresses] = c.addresses
	m[configKeyTimeout] = c.timeout

	return result
}

func (c *BaseConfig) Logger() any         { return c.logger }
func (c *BaseConfig) Addresses() []string { return c.addresses }
func (c *BaseConfig) Timeout() int        { return c.timeout }

func (c *BaseConfig) SetLogger(n any)         { c.logger = n }
func (c *BaseConfig) SetAddresses(n []string) { c.addresses = n }
func (c *BaseConfig) SetTimeout(n int)        { c.timeout = n }
