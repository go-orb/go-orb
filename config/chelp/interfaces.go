package chelp

// ConfigLoadStore needs to be implemented by ALL config handlers.
type ConfigLoadStore interface {
	Load(map[string]any) error
	Store(map[string]any) error
	Merge(aConfig any) error
}

// PluginConfig is the basic interface for the most plugin related config handlers.
type PluginConfig interface {
	ConfigLoadStore

	Plugin() string
	Enabled() *bool

	SetPlugin(n string)
	SetEnabled(n *bool)
}
