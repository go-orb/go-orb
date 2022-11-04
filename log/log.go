// Package log contains a golang.org/x/exp/slog compatible logger.
package log

import (
	"context"
	"fmt"

	"github.com/go-orb/config"
	"github.com/go-orb/config/source"
	"golang.org/x/exp/slog"

	"go-micro.dev/v5/types/component"

	"go-micro.dev/v5/types"
)

// This is here to make sure Logger implements the component interface.
var _ component.Component = &Logger{}

const (
	ComponentType component.Type = "logger"
)

type keyOne struct{}
type keyTwo struct{}

// Logger is a go-micro logger, it is the slog.Logger, with some added methods
// to implement the component interface.
type Logger struct {
	slog.Logger

	config Config

	// plugins is a cache of lazyloaded plugin handlers.
	// In order to prevent creating multiple handlers, and thus potentially
	// multiple connections, depending on the handler, we cache the handlers, and
	// wrap them with a LevelHandler by default. This way we only create one
	// handler per plugin, for use in any amount of loggers.
	plugins map[string]slog.Handler
}

// Plugin will return the plugin handler with set to TRACE level. To enable
// a custom level wrap it with a LevelHandler.
func (l *Logger) Plugin(plugin string) (slog.Handler, error) {
	if h, ok := l.plugins[plugin]; ok {
		return h, nil
	}

	p, err := Plugins.Get(plugin)
	if err != nil {
		return nil, fmt.Errorf("logger plugin '%s' does not exist, please register your plugin", plugin)
	}

	handler, err := p(TraceLevel)
	if err != nil {
		return nil, fmt.Errorf("create new plugin handler: %w", err)
	}

	return handler, nil
}

// Start no-op.
func (l *Logger) Start() error {
	return nil
}

// Stop no-op.
func (l *Logger) Stop() error {
	return nil
}

// String returns current plugin used.
func (l *Logger) String() string {
	// TODO: maybe this should call smth like handler.String()
	return l.config.Plugin
}

// Type returns the component type.
func (l *Logger) Type() component.Type {
	return ComponentType
}

// New creates a new Logger from a Config.
func New(cfg Config) (Logger, error) {
	l := Logger{
		config:  cfg,
		plugins: make(map[string]slog.Handler),
	}

	h, err := l.Plugin(cfg.Plugin)
	if err != nil {
		return Logger{}, err
	}

	h, err = NewLevelHandler(cfg.Level, h)
	if err != nil {
		return Logger{}, err
	}

	l.Logger = slog.New(h)

	return l, nil
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

// NewComponentLogger will create a sub logger for a component inheriting all
// parrent logger fields, and optionally set a new level and handler.
// If you want to use the parrent handler and log level, pass empty strings.
// It will add two fields to the sub logger, the component (e.g. broker)
// and the component plugin implementation (e.g. NATS).
func NewComponentLogger(logger Logger, component component.Type, name, plugin, level string) (Logger, error) {
	errMsg := "(component: %s, name: %s, plugin: %s) create component logger: %w"

	var err error

	lvl := logger.config.Level
	if len(level) > 0 {
		lvl, err = ParseLevel(level)
		if err != nil {
			return Logger{}, fmt.Errorf(errMsg, component, name, plugin, err)
		}
	}

	// Optionally avoid wrapping a handler if the level is the same as the parent
	// logger, and not different handler is requested. To check the handler level
	// it needs to implent the Leveler interface, which is not provided by default
	// on slog handlers, and needs to be implemented manually on handler plugins.
	noWrapper := false

	// If a new handler is requested, fetch one from cache or creata a new on.
	// If no new handler is requested check if we can skip handler wrapping.
	handler := logger.Handler()
	if len(plugin) > 0 {
		handler, err = logger.Plugin(plugin)
		if err != nil {
			return Logger{}, fmt.Errorf(errMsg, component, name, plugin, err)
		}
	} else {
		if level, ok := logger.Handler().(slog.Leveler); ok && level == lvl {
			noWrapper = true
		}
	}

	if !noWrapper {
		handler, err = NewLevelHandler(lvl, handler)
		if err != nil {
			return Logger{}, err
		}
	}

	// FIX:this doesn't work. No way to extract fields with context
	ctx := logger.With(
		slog.String("component", string(component)),
		slog.String("plugin", name),
	).Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = context.WithValue(ctx, keyOne{}, string(component))
	ctx = context.WithValue(ctx, keyTwo{}, name)

	logger.Logger = slog.New(handler).WithContext(ctx)

	return logger, nil
}
