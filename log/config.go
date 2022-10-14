package log

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"jochum.dev/orb/orb/config/chelp"
)

const (
	configKeyFields          = "fields"
	configKeyLevel           = "level"
	configKeyCallerSkipFrame = "caller_skip_frame"
)

// Config is the basic configuration which every log plugin config should implement.
type Config interface {
	chelp.Plugin

	Fields() map[string]any
	Level() string
	CallerSkipFrame() int

	SetFields(map[string]any)
	SetLevel(string)
	SetCallerSkipFrame(int)
}

// BaseConfig is a basic configuration for loggers.
type BaseConfig struct {
	chelp.Plugin
	fields          map[string]any
	level           string
	callerSkipFrame int
}

// NewConfig creates a new BaseConfig.
func NewConfig() *BaseConfig {
	return &BaseConfig{
		Plugin: chelp.NewPluginConfig(),
	}
}

// DefaultConfig returns the default config for a given Plugin.
func DefaultConfig(name string) (any, error) {
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

	return DefaultConfig(pconf.Plugin())
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

	if err := c.Plugin.Load(m); err != nil {
		result = multierror.Append(err)
	}

	// Optional
	var err error

	c.fields, err = chelp.Get(m, configKeyFields, c.fields)
	if err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	c.level, err = chelp.Get(m, configKeyLevel, c.level)
	if err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	c.callerSkipFrame, err = chelp.Get(m, configKeyCallerSkipFrame, c.callerSkipFrame)
	if err != nil && !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	return result
}

// Store stores this config in a map[string]any.
func (c *BaseConfig) Store(m map[string]any) error {
	if err := c.Plugin.Store(m); err != nil {
		return err
	}

	m[configKeyFields] = c.fields
	m[configKeyLevel] = c.level
	m[configKeyCallerSkipFrame] = c.callerSkipFrame

	return nil
}

func (c *BaseConfig) Fields() map[string]any { return c.fields }
func (c *BaseConfig) Level() string          { return c.level }
func (c *BaseConfig) CallerSkipFrame() int   { return c.callerSkipFrame }

func (c *BaseConfig) SetFields(n map[string]any) { c.fields = n }
func (c *BaseConfig) SetLevel(n string)          { c.level = n }
func (c *BaseConfig) SetCallerSkipFrame(n int)   { c.callerSkipFrame = n }
