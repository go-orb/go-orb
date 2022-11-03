// Package log contains a golang.org/x/exp/slog compatible logger.
package log

import (
	"github.com/go-orb/config"
	"github.com/go-orb/config/source"
	"golang.org/x/exp/slog"

	"go-micro.dev/v5/types/component"

	"github.com/go-orb/orb/types"
)

// This is here to make sure Logger implements the component interface.
var _ component.Component = &Logger{}

const (
	ComponentType component.Type = "logger"
)

// Logger is a go-micro logger, it is the slog.Logger, with some added methods
// to implement the component interface.
type Logger struct {
	slog.Logger

	plugin string
}

func (l *Logger) Start() error {
	return nil
}

func (l *Logger) Stop() error {
	return nil
}

func (l *Logger) String() string {
	return l.plugin
}

func (l *Logger) Type() component.Type {
	return ComponentType
}

// New creates a new Logger from a Config.
func New(cfg Config) (Logger, error) {
	level, err := ParseLevel(cfg.Level)
	if err != nil {
		return Logger{}, err
	}

	handlerFunc, err := Plugins.Get(cfg.Plugin)
	if err != nil {
		return Logger{}, err
	}

	h, err := handlerFunc(level)
	if err != nil {
		return Logger{}, err
	}

	return Logger{
		plugin: cfg.Plugin,
		Logger: slog.New(h),
	}, nil
}

// ProvideLogger provides a new logger to wire.
func ProvideLogger(serviceName types.ServiceName, data []source.Data, opts ...Option) (Logger, error) {
	cfg := NewConfig()

	for _, o := range opts {
		o(&cfg)
	}

	sections := types.SplitServiceName(serviceName)
	if err := config.Parse(append(sections, "logger"), data, cfg); err != nil {
		return Logger{}, err
	}

	logger, err := New(cfg)
	if err != nil {
		return Logger{}, err
	}

	slog.SetDefault(logger.Logger)

	return logger, nil
}
