package config

import (
	"errors"
)

// ConfigPluginImpl is the base for every plugin config.
type ConfigPluginImpl struct {
	Plugin  string `json:"plugin" yaml:"plugin"`
	Enabled *bool  `json:"enabled" yaml:"enabled"`
}

// NewPluginConfig creates a new BasePluginConfig, it's used to parse the initial plugin.
func NewPluginConfig() *ConfigPluginImpl {
	return &ConfigPluginImpl{}
}

// MergePrevious merges the previous config into this one.
func (c *ConfigPluginImpl) MergePrevious(aPreviousConfig any) error {
	previousConfig, ok := aPreviousConfig.(ConfigPlugin)
	if !ok {
		return ErrUnknownConfig
	}

	defConfig := NewPluginConfig()

	if c.Plugin == defConfig.Plugin {
		c.Plugin = previousConfig.GetPlugin()
	}

	if c.Enabled == defConfig.Enabled {
		c.Enabled = previousConfig.GetEnabled()
	}

	return nil
}

// Plugin returns the plugin.
func (c *ConfigPluginImpl) GetPlugin() string { return c.Plugin }

// Enabled returns if this component has been enabled.
func (c *ConfigPluginImpl) GetEnabled() *bool { return c.Enabled }

// GetPluginFromConfigData returns the plugin from []Data.
func GetPluginFromConfigData(baseSection string, pluginSection string, configDatas []Data) (string, error) {
	result := ""
	for _, configData := range configDatas {
		// Go own section deeper.
		var err error
		data := configData.Data
		if baseSection != "" {
			if data, err = Get(data, baseSection, map[string]any{}); err != nil {
				// Ignore unknown configSection in config.
				if errors.Is(err, ErrNotExistent) {
					continue
				}
				return result, err
			}
		}

		// Now fetch my own section.
		if data, err = Get(data, pluginSection, map[string]any{}); err != nil {
			// Ignore unknown configSection in config.
			if errors.Is(err, ErrNotExistent) {
				continue
			}
			return result, err
		}

		// Next fetch the plugin name.
		if result, err := Get(data, "plugin", result); !errors.Is(err, ErrNotExistent) {
			return result, err
		}
	}

	return result, nil
}
