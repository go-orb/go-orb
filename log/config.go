package log

import "jochum.dev/orb/orb/config"

// Config is the basic configuration which every log plugin config should implement.
type Config interface {
	config.ConfigPlugin

	GetFields() map[string]any
	GetLevel() string
	GetCallerSkipFrame() int
}

// ConfigImpl is a basic configuration for loggers.
type ConfigImpl struct {
	*config.ConfigPluginImpl

	Fields          map[string]any `json:"fields" yaml:"Fields"`
	Level           string         `json:"level" yaml:"Level"`
	CallerSkipFrame int            `json:"caller_skip_frame" yaml:"CallerSkipFrame"`
}

// NewComponentConfig creates a new BaseConfig.
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

	if c.Fields == nil {
		c.Fields = previous.GetFields()
	}

	if c.Level == defConfig.Level {
		c.Level = previous.GetLevel()
	}

	if c.CallerSkipFrame != defConfig.CallerSkipFrame {
		c.CallerSkipFrame = previous.GetCallerSkipFrame()
	}

	return nil
}

func (c *ConfigImpl) GetFields() map[string]any { return c.Fields }
func (c *ConfigImpl) GetLevel() string          { return c.Level }
func (c *ConfigImpl) GetCallerSkipFrame() int   { return c.CallerSkipFrame }
