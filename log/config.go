package log

import (
	"github.com/hashicorp/go-multierror"
	"jochum.dev/orb/orb/config/chelp"
)

const (
	CONFIG_KEY_FIELDS          = "fields"
	CONFIG_KEY_LEVEL           = "level"
	CONFIG_KEY_CALLERSKIPFRAME = "caller_skip_frame"
)

type Config interface {
	chelp.PluginConfig

	Fields() map[string]any
	Level() string
	CallerSkipFrame() int

	SetFields(map[string]any)
	SetLevel(string)
	SetCallerSkipFrame(int)
}

type BaseConfig struct {
	chelp.PluginConfig
	fields          map[string]any
	level           string
	callerSkipFrame int
}

func NewConfig() Config {
	return &BaseConfig{
		PluginConfig: chelp.NewPluginConfig(),
	}
}

func (c *BaseConfig) Load(m map[string]any) error {
	var result error

	if err := c.PluginConfig.Load(m); err != nil {
		result = multierror.Append(err)
	}

	// Optional
	var err error
	if c.fields, err = chelp.Get(m, CONFIG_KEY_FIELDS, map[string]any{}); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.level, err = chelp.Get(m, CONFIG_KEY_LEVEL, "info"); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}
	if c.callerSkipFrame, err = chelp.Get(m, CONFIG_KEY_CALLERSKIPFRAME, 0); err != nil && err != chelp.ErrNotExistant {
		result = multierror.Append(err)
	}

	return result
}

func (c *BaseConfig) Store(m map[string]any) error {
	if err := c.PluginConfig.Store(m); err != nil {
		return err
	}

	m[CONFIG_KEY_FIELDS] = c.fields
	m[CONFIG_KEY_LEVEL] = c.level
	m[CONFIG_KEY_CALLERSKIPFRAME] = c.callerSkipFrame

	return nil
}

func (c *BaseConfig) Fields() map[string]any { return c.fields }
func (c *BaseConfig) Level() string          { return c.level }
func (c *BaseConfig) CallerSkipFrame() int   { return c.callerSkipFrame }

func (c *BaseConfig) SetFields(n map[string]any) { c.fields = n }
func (c *BaseConfig) SetLevel(n string)          { c.level = n }
func (c *BaseConfig) SetCallerSkipFrame(n int)   { c.callerSkipFrame = n }
