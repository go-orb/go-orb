package config

// ConfigMerge needs to be implemented by ALL config handlers.
type ConfigMerge interface {
	MergePrevious(aPreviousConfig any) error
}

// ConfigPlugin is the basic interface for the most plugin related config handlers.
type ConfigPlugin interface {
	ConfigMerge

	GetPlugin() string
	GetEnabled() *bool
}
