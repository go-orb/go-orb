package chelp

// ConfigMethods needs to be implemented by ALL config handlers.
type ConfigMethods interface {
	Load(map[string]any) error
	Store(map[string]any) error
	Merge(aConfig any) error
}

// Plugin is the basic interface for the most plugin related config handlers.
type Plugin interface {
	ConfigMethods

	Plugin() string
	Enabled() *bool

	SetPlugin(n string)
	SetEnabled(n *bool)
}
