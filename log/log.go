// Package log contains a golang.org/x/exp/slog compatible logger.
package log

import (
	"github.com/go-orb/config"
	"github.com/go-orb/config/source"
	"github.com/go-orb/orb/types"
	"golang.org/x/exp/slog"
)

// Logger is the logger we use.
type Logger struct {
	*slog.Logger

	plugin string
}

// This is here to make sure Logger implements types.Component.
var _ types.Component = &Logger{}

func (l *Logger) Start() error {
	return nil
}

func (l *Logger) Stop() error {
	return nil
}

func (l *Logger) String() string {
	return l.plugin
}

func (l *Logger) Type() string {
	return "logger"
}

// New creates a new Logger from a Config..
func New(cfg *Config) (*Logger, error) {
	level, err := ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	handlerFunc, err := Plugins.Get(cfg.Plugin)
	if err != nil {
		return nil, err
	}

	h, err := handlerFunc(level)
	if err != nil {
		return nil, err
	}

	return &Logger{
		plugin: cfg.Plugin,
		Logger: slog.New(h),
	}, nil
}

// Provide provides a new logger to wire.
func Provide(
	serviceName types.ServiceName,
	datas []source.Data,
) (*Logger, error) {
	cfg := NewConfig()

	sections := types.SplitServiceName(serviceName)
	if err := config.Parse(append(sections, "logger"), datas, cfg); err != nil {
		return nil, err
	}

	logger, err := New(cfg)
	if err != nil {
		return nil, err
	}

	// Make it the default logger.
	slog.SetDefault(logger.Logger)

	return logger, nil
}
