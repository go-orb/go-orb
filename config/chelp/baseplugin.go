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

func (c *BasePluginConfig) Load(m map[string]any) error {
	var (
		result error
		err    error
	)

	if c.plugin, err = Get(m, configPlugin, ""); err != nil {
		result = multierror.Append(err)
	}

	if c.enabled, err = Get[*bool](m, configEnabled, nil); err != nil {
		result = multierror.Append(err)
	}

	return result
}

func (c *BasePluginConfig) Store(m map[string]any) error {
	m[configPlugin] = c.plugin
	m[configEnabled] = c.enabled

	return nil
}

func (c *BasePluginConfig) Plugin() string     { return c.plugin }
func (c *BasePluginConfig) Enabled() *bool     { return c.enabled }
func (c *BasePluginConfig) SetPlugin(n string) { c.plugin = n }
func (c *BasePluginConfig) SetEnabled(n *bool) { c.enabled = n }
