package chelp

const (
	CONFIG_PLUGIN  = "plugin"
	CONFIG_ENABLED = "enabled"
)

type BasePluginConfig struct {
	plugin  string
	enabled *bool
}

func NewPluginConfig() *BasePluginConfig {
	return &BasePluginConfig{}
}

func (c *BasePluginConfig) Load(m map[string]any) error {
	var err error
	if c.plugin, err = Get(m, CONFIG_PLUGIN, ""); err != nil {
		return err
	}
	if c.enabled, err = Get[*bool](m, CONFIG_ENABLED, nil); err != nil {
		return err
	}

	return nil
}

func (c *BasePluginConfig) Store(m map[string]any) error {
	m[CONFIG_PLUGIN] = c.plugin
	m[CONFIG_ENABLED] = c.enabled

	return nil
}

func (c *BasePluginConfig) Plugin() string { return c.plugin }
func (c *BasePluginConfig) Enabled() *bool { return c.enabled }
