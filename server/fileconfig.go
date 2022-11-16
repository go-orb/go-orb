package server

// fileConfigServer is used to parse the config section from config files.
// It checks if the user has any custom dynamic entrypoints defined.
//
// The config itself contains a lot more data, but we only need to know the
// entrypoint plugins, and list of entrypoints with their names.
type fileConfigServer struct {
	// Enabled allows you to easily disable all entrypionts of one plugin type.
	Enabled     *bool `json:"enabled" yaml:"enabled"`
	Entrypoints []*struct {
		Name string `json:"name,omitempty" yaml:"name,omitempty"`

		// Enabled allows you to disable specific entrypoints at runtime.
		// a pointer is used here to distinguish between unset vs default value.
		Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

		// Inherit allows you to Inherit a config from a different entrypoint
		Inherit string `json:"inherit,omitempty" yaml:"inherit,omitempty"`
	} `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
}

// IsEnabled checks whether an entrypoint should be enabled by name.
func (c *fileConfigServer) IsEnabled(name string) bool {
	for _, e := range c.Entrypoints {
		if name == e.Name && e.Enabled != nil {
			return *e.Enabled
		}
	}

	return true
}

// Inherit returns the value of the inherit field for a specific entrypoint,
// if present.
func (c *fileConfigServer) Inherit(name string) string {
	for _, e := range c.Entrypoints {
		if name == e.Name {
			return e.Inherit
		}
	}

	return ""
}
