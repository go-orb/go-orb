package log

//
//nolint:gochecknoglobals
var (
	// DefaultLevel is info.
	DefaultLevel = LevelInfo.String()
	// DefaultPlugin is "slog", it support's json and text output to Stdout/Stderr and file.
	DefaultPlugin = "slog"
	// DefaultSetDefault set's the "log" and slog default logger when true.
	DefaultSetDefault = false
	// DefaultConfigSection is the section key used in config files used to
	// configure the logger options.
	DefaultConfigSection = "logger"
)

var _ (ConfigType) = (*Config)(nil)

// Option is a logger WithXXX Option.
type Option func(ConfigType)

// ConfigType is a wrapper for config, so we can pass it back to the this plugin handler.
type ConfigType interface {
	config() *Config
}

// Config is the loggers config.
type Config struct {
	// Plugin sets the log handler plugin to use.
	// Make sure to register the plugin by importing it.
	Plugin string `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	// Level sets the log level to use.
	Level string `json:"level,omitempty" yaml:"level,omitempty"`
	// SetDefault indicates if this logger should be set as default logger.
	SetDefault bool `json:"setDefault,omitempty" yaml:"setDefault,omitempty"`
	// Fields are additional fields to add to the log.
	Fields map[string]any `json:"fields,omitempty" yaml:"fields,omitempty"`
}

func (c *Config) config() *Config {
	return c
}

// NewConfig creates a new config with the defaults and applys opts on top.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		Level:  DefaultLevel,
		Plugin: DefaultPlugin,
	}

	// Apply options.
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// WithLevel sets the log level to user.
// TODO(davincible): would love to take in something like (	slog.Level | string | constraints.Integer) here,
// but not sure how that would work.
func WithLevel(n string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Level = n
	}
}

// WithPlugin sets the logger plugin to be used.
// A logger plugin is the underlying handler the logger will use to process
// log events. To add your custom handler, register it as a plugin.
// See log/plugin.go for more details on how to do so.
func WithPlugin(n string) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Plugin = n
	}
}

// WithSetDefault makes the resulting logger the default logger.
// TODO(jochumdev): Remove this? SetDefault also stops all Plugins.
func WithSetDefault() Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.SetDefault = true
	}
}

// WithFields adds additional fields to the log.
func WithFields(n map[string]any) Option {
	return func(cfg ConfigType) {
		c := cfg.config()
		c.Fields = n
	}
}
