package chelp

import (
	"github.com/hashicorp/go-multierror"
)

const (
	configPlugin  = "plugin"
	configEnabled = "enabled"
)

// BasePluginConfig is the base for every plugin config.
type BasePluginConfig struct {
	plugin  string
	enabled *bool
}

// NewPluginConfig creates a new BasePluginConfig, it's used to parse the initial plugin.
func NewPluginConfig() *BasePluginConfig {
	return &BasePluginConfig{}
}

// Load loads this config from map[string]any.
func (c *BasePluginConfig) Load(m map[string]any) error {
	var (
		result error
		err    error
	)

	if c.plugin, err = Get(m, configPlugin, ""); err != nil {
		result = multierror.Append(result, err)
	}

	if c.enabled, err = Get[*bool](m, configEnabled, nil); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// Store stores this config in a map[string]any.
func (c *BasePluginConfig) Store(m map[string]any) error {
	m[configPlugin] = c.plugin
	m[configEnabled] = c.enabled

	return nil
}

// Merge merges a config into this config.
func (c *BasePluginConfig) Merge(aConfig any) error {
	toMerge, ok := aConfig.(Plugin)
	if !ok {
		return ErrUnknownConfig
	}

	defConfig := NewPluginConfig()

	if toMerge.Plugin() != defConfig.Plugin() {
		c.SetPlugin(toMerge.Plugin())
	}

	if toMerge.Enabled() != defConfig.Enabled() {
		c.SetEnabled(toMerge.Enabled())
	}

	return nil
}

// Plugin returns the plugin.
func (c *BasePluginConfig) Plugin() string { return c.plugin }

// Enabled returns if this component has been enabled.
func (c *BasePluginConfig) Enabled() *bool { return c.enabled }

// SetPlugin updates plugin.
func (c *BasePluginConfig) SetPlugin(n string) { c.plugin = n }

// SetEnabled updates enabled.
func (c *BasePluginConfig) SetEnabled(n *bool) { c.enabled = n }
