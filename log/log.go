// Package log contains a golang.org/x/exp/slog compatible logger.
package log

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/go-orb/go-orb/config"

	"github.com/go-orb/go-orb/types"
)

// This is here to make sure Logger implements the component interface.
var _ types.Component = (*Logger)(nil)

// DefaultConfigSection is the section key used in config files used to
// configure the logger options.
var DefaultConfigSection = "logger" //nolint:gochecknoglobals

// ComponentType is the name of the component type logger.
const ComponentType = "logger"

// Logger is a go-micro logger, it is the slog.Logger, with some added methods
// to implement the component interface.
type Logger struct {
	*slog.Logger

	// config is the config used to create the current logger.
	// It is not exported as it also acts as a state, and should not be modified
	// externally. It keeps track of the current level set, and plugin used.
	config Config

	// fields are all parameters passed to Logger.With. We keep track of them
	// in case a sublogger needs to be created with a different plugin, then we
	// manually need to add the fields to the handler plugin.
	fields []any
}

// New creates a new Logger from a Config.
func New(cfg Config) (Logger, error) {
	return Logger{
		config: cfg,
	}.Plugin(cfg.Plugin, cfg.Level)
}

// ProvideLogger provides a new logger.
// It will set the slog.Logger as package wide default logger.
func ProvideLogger(
	serviceName types.ServiceName,
	data types.ConfigData,
	opts ...Option,
) (Logger, error) {
	cfg := NewConfig()

	for _, o := range opts {
		o(&cfg)
	}

	sections := types.SplitServiceName(serviceName)
	if err := config.Parse(append(sections, DefaultConfigSection), data, cfg); err != nil {
		return Logger{}, err
	}

	logger, err := New(cfg)
	if err != nil {
		return Logger{}, err
	}

	if cfg.SetDefault {
		slog.SetDefault(logger.Logger)
	}

	return logger, nil
}

// Plugin will return the plugin handler with set to TRACE level. To enable
// a custom level wrap it with a LevelHandler.
func (l Logger) Plugin(plugin string, level ...slog.Leveler) (Logger, error) {
	lvl := l.config.Level
	if len(level) > 0 && level[0] != nil {
		lvl = level[0].Level()
		l.config.Level = lvl
	}

	h, err := plugins.Get(plugin)
	handler := h.handler

	// Check if already have an instance of the plugin. If not create a new one.
	if err != nil {
		p, err := Plugins.Get(plugin)
		if err != nil {
			return l, fmt.Errorf("logger plugin '%s' does not exist, please register your plugin", plugin)
		}

		handler, err = p(lvl)
		if err != nil {
			return l, fmt.Errorf("create new plugin handler: %w", err)
		}

		plugins.Set(plugin, pluginHandler{
			handler: handler,
			level:   lvl,
		})
	} else if h.level != lvl {
		handler = &LevelHandler{level: lvl, handler: handler}
	}

	l.config.Plugin = plugin
	l.Logger = slog.New(handler)

	if len(l.fields) > 0 {
		l.Logger = l.Logger.With(l.fields...)
	}

	return l, nil
}

// WithLevel creates a copy of the logger with a new level.
// It will inherit all the fields and the context from the parent logger.
func (l Logger) WithLevel(level slog.Leveler) Logger {
	if level != nil {
		l.config.Level = level.Level()
		l.Logger = slog.New(&LevelHandler{level.Level(), l.Handler()}).WithContext(l.Context())
	}

	return l
}

// With returns a new Logger that includes the given arguments, converted to
// Attrs as in [Logger.Log]. The Attrs will be added to each output from the
// Logger.
//
// The new Logger's handler is the result of calling WithAttrs on the receiver's
// handler.
func (l *Logger) With(args ...any) Logger {
	l.fields = append(l.fields, args...)
	l.Logger = l.Logger.With(args...)

	return *l
}

// WithContext returns a new Logger with the same handler as the receiver and the given context.
func (l Logger) WithContext(ctx context.Context) Logger {
	l.Logger = l.Logger.WithContext(ctx)
	return l
}

// WithComponent will create a new logger for a component inheriting all
// parent logger fields, and optionally set a new level and handler.
//
// If you want to use the parent handler and log level, pass an empty values.
//
// It will add two fields to the sub logger, the component (e.g. broker) as component
// and the component plugin implementation (e.g. NATS) as plugin.
func (l Logger) WithComponent(component, name, plugin string, level slog.Leveler) (Logger, error) {
	l = l.With(
		slog.String("component", component),
		slog.String("plugin", name),
	)

	var err error
	if len(plugin) > 0 {
		l, err = l.Plugin(plugin, level)
		if err != nil {
			return Logger{}, err
		}
	}

	if level.Level() != l.config.Level {
		l = l.WithLevel(level)
	}

	return l, nil
}

// Start no-op.
func (l Logger) Start() error {
	return nil
}

// Stop no-op.
func (l Logger) Stop(_ context.Context) error {
	return nil
}

// String returns current handler plugin used.
func (l Logger) String() string {
	if i, ok := l.Handler().(fmt.Stringer); ok {
		return i.String()
	}

	return l.config.Plugin
}

// Type returns the component type.
func (l Logger) Type() string {
	return ComponentType
}

// Trace logs at TraceLevel.
func (l *Logger) Trace(msg string, args ...any) {
	l.LogDepth(0, TraceLevel, msg, args...)
}
