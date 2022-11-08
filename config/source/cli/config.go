package cli

// Config is the base config for this component.
type Config struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

// NewConfig returns the cli config.
func NewConfig() *Config {
	return &Config{}
}
