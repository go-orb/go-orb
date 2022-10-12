package component

import (
	"github.com/hashicorp/go-multierror"
	"jochum.dev/orb/orb/config/chelp"
	"jochum.dev/orb/orb/log"
)

const (
	CONFIG_KEY_LOGGER   = "logger"
	CONFIG_KEY_ID       = "id"
	CONFIG_KEY_NAME     = "name"
	CONFIG_KEY_VERSION  = "version"
	CONFIG_KEY_METADATA = "metadata"
)

type Config interface {
	chelp.PluginConfig

	// Required
	Name() string
	Version() string

	// Optional
	Logger() log.Config
	ID() string
	Metadata() map[string]string
}

type BaseConfig struct {
	*chelp.BasePluginConfig

	name    string
	version string

	logger   log.Config
	id       string
	metadata map[string]string
}

func NewBaseConfig() Config {
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
	if c.name, err = chelp.Get(m, CONFIG_KEY_NAME, ""); err != nil {
		result = multierror.Append(err)
	}
	if c.version, err = chelp.Get(m, CONFIG_KEY_VERSION, ""); err != nil {
		result = multierror.Append(err)
	}

	// Optional
	if err := c.logger.Load(m); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}

	if c.id, err = chelp.Get(m, CONFIG_KEY_ID, ""); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.metadata, err = chelp.Get(m, CONFIG_KEY_METADATA, map[string]string{}); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}

	return result
}

func (c *BaseConfig) Store(m map[string]any) error {
	var result error

	if err := c.BasePluginConfig.Store(m); err != nil {
		result = multierror.Append(err)
	}

	m[CONFIG_KEY_NAME] = c.name
	m[CONFIG_KEY_VERSION] = c.version

	logger := make(map[string]any)
	if err := c.logger.Store(logger); err != nil {
		result = multierror.Append(err)
	}
	m[CONFIG_KEY_LOGGER] = logger

	m[CONFIG_KEY_ID] = c.id
	m[CONFIG_KEY_METADATA] = c.metadata

	return result
}

func (c *BaseConfig) Name() string                { return c.name }
func (c *BaseConfig) Version() string             { return c.version }
func (c *BaseConfig) Logger() log.Config          { return c.logger }
func (c *BaseConfig) ID() string                  { return c.id }
func (c *BaseConfig) Metadata() map[string]string { return c.metadata }
