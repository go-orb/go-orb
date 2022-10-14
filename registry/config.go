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
	chelp.Plugin

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

func NewBaseConfig() *BaseConfig {
	return &BaseConfig{
		BasePluginConfig: chelp.NewPluginConfig(),
		logger:           log.NewConfig(),
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

func getConfig(m map[string]any) (any, error) {
	pconf := chelp.NewPluginConfig()
	if err := pconf.Load(m); err != nil {
		return nil, err
	}

	return NewConfig(pconf.Plugin())
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
func (c *BaseConfig) Load(m map[string]any) error {
	var result error

	// Required
	if err := c.BasePluginConfig.Load(m); err != nil {
		result = multierror.Append(result, err)
	}

	var (
		err error
	)

	// Optional
	c.logger, err = log.LoadConfig(m, configKeyLogger)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
	}

	c.addresses, err = chelp.Get(m, configKeyAddresses, c.addresses)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
	}

	c.timeout, err = chelp.Get(m, configKeyTimeout, c.timeout)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
	}

	return result
}

// Store stores this config in a map[string]any.
func (c *BaseConfig) Store(m map[string]any) error {
	var result error

	if err := c.BasePluginConfig.Store(m); err != nil {
		result = multierror.Append(result, err)
	}

	var err error

	m[configKeyLogger], err = log.StoreConfig(c.logger)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(result, err)
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
