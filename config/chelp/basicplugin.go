package chelp

const (
	CONFIG_PLUGIN  = "plugin"
	CONFIG_ENABLED = "enabled"
)

type BasicPlugin struct {
	plugin  string
	enabled *bool
}

func NewBasicPlugin() *BasicPlugin {
	return &BasicPlugin{}
}

func (c *BasicPlugin) Load(m map[string]any) error {
	var err error
	if c.plugin, err = Get(m, CONFIG_PLUGIN, ""); err != nil {
		return err
	}
	if c.enabled, err = Get[*bool](m, CONFIG_ENABLED, nil); err != nil {
		return err
	}

	return nil
}

func (c *BasicPlugin) Store(m map[string]any) error {
	m[CONFIG_PLUGIN] = c.plugin
	m[CONFIG_ENABLED] = c.enabled

	return nil
}

func (c *BasicPlugin) Plugin() string { return c.plugin }
func (c *BasicPlugin) Enabled() *bool { return c.enabled }
