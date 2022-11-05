package log

import (
	"golang.org/x/exp/slog"
)

// TODO:  Something like this would be nice
// type LevelT interface {
// 	slog.Level | string | constraints.Integer
// }

// DefaultPlugin is the default plugin to use.
const (
	DefaultLevel  = InfoLevel
	DefaultPlugin = "textstderr"
)

type Option func(*Config)

// Config is the loggers config.
type Config struct {
	Plugin string     `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	Level  slog.Level `json:"level,omitempty" yaml:"level,omitempty"`
}

// NewConfig creates a new config with the defaults.
func NewConfig() Config {
	return Config{Level: DefaultLevel, Plugin: DefaultPlugin}
}

// WithLevel sets the log level to user.
// TODO: would love to take in something like (	slog.Level | string | constraints.Integer) here, but not sure how that would work
func WithLevel(level slog.Level) Option {
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
