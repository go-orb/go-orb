package log

// DefaultPlugin is the default plugin to use.
const DefaultPlugin = "textstderr"

// Config is the loggers config.
type Config struct {
	Plugin string `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	Level  string `json:"level,omitempty" yaml:"level,omitempty"`
}

// NewConfig creates a new config with the defaults.
func NewConfig() *Config {
	return &Config{Level: "INFO", Plugin: DefaultPlugin}
}
