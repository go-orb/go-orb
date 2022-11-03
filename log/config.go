package log

// DefaultPlugin is the default plugin to use.
const (
	DefaultLevel  = "INFO"
	DefaultPlugin = "textstderr"
)

type Option func(*Config)

// Config is the loggers config.
type Config struct {
	Plugin string `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	Level  string `json:"level,omitempty" yaml:"level,omitempty"`
}

// NewConfig creates a new config with the defaults.
func NewConfig() Config {
	return Config{Level: DefaultLevel, Plugin: DefaultPlugin}
}

// WithLevel sets the log level to user.
func WithLevel(level string) Option {
	return func(c *Config) {
		c.Level = level
	}
}

// WithPlugin sets the logger plugin to be used.
// A logger plugin is the underlying handler the logger will use to process
// log events. To add your custom handler, register it as a plugin.
// See log/plugin.go for more details on how to do so.
func WithPlugin(plugin string) Option {
	return func(c *Config) {
		c.Plugin = plugin
	}
}
