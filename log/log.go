// Package log contains a golang.org/x/exp/slog compatible logger.
package log

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/config/source"

	"github.com/go-orb/go-orb/types"
)

// This is here to make sure Logger implements the component interface.
var _ types.Component = (*Logger)(nil)

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

	// pluginProvider is the provider of the current log handler.
	pluginProvider ProviderType

	// fields are all parameters passed to Logger.With. We keep track of them
	// in case a sublogger needs to be created with a different plugin, then we
	// manually need to add the fields to the handler plugin.
	fields []any
}

// New creates a n4ew Logger from a Config.
func New(opts ...Option) (Logger, error) {
	return NewConfigDatas([]string{}, nil, opts...)
}

// NewConfigDatas will create a new logger with the given configs,
// as well as this logger's fields.
func NewConfigDatas(sections []string, configs types.ConfigData, opts ...Option) (Logger, error) {
	var cfg Config
	if configs == nil {
		cfg = NewConfig(opts...)

		data, err := config.ParseStruct(append(sections, DefaultConfigSection), &cfg)
		if err != nil {
			return Logger{}, fmt.Errorf("while creating a new config: %w", err)
		}

		configs = []source.Data{data}
	} else {
		cfg = NewConfig(opts...)
		if err := config.Parse(append(sections, DefaultConfigSection), configs, &cfg); err != nil {
			return Logger{}, fmt.Errorf("while creating a new config: %w", err)
		}
	}

	pf, err := plugins.Get(cfg.Plugin)
	if err != nil {
		slog.Error("getting a logger plugin", "plugin", cfg.Plugin, "error", err)
		return Logger{}, fmt.Errorf("while getting the log plugin '%s': %w", cfg.Plugin, err)
	}

	provider, err := pf(sections, configs, opts...)
	if err != nil {
		return Logger{}, err
	}

	cachedProvider, err := pluginsCache.Get(provider.Key())
	if err != nil {
		if err := provider.Start(); err != nil {
			return Logger{}, err
		}

		pluginsCache.Set(cfg.Plugin, provider)
		cachedProvider = provider
	}

	handler, err := cachedProvider.Handler()
	if err != nil {
		return Logger{}, err
	}

	lvlHandler, err := NewLevelHandler(cfg.Level, handler)
	if err != nil {
		return Logger{}, err
	}

	r := Logger{
		Logger:         slog.New(lvlHandler),
		pluginProvider: cachedProvider,
		config:         cfg,
		fields:         []any{},
	}

	return r, nil
}

// ProvideLogger provides a new logger.
// It will set the slog.Logger as package wide default logger.
func ProvideLogger(
	serviceName types.ServiceName,
	configs types.ConfigData,
	opts ...Option,
) (Logger, error) {
	sections := types.SplitServiceName(serviceName)

	logger, err := NewConfigDatas(sections, configs, opts...)
	if err != nil {
		return Logger{}, err
	}

	slog.SetDefault(logger.Logger)

	return logger, nil
}

// WithLevel creates a copy of the logger with a new level.
// It will inherit all the fields and the context from the parent logger.
func (l Logger) WithLevel(level slog.Leveler) Logger {
	if level != nil {
		l.config.Level = level.Level()

		handler, err := l.pluginProvider.Handler()
		if err != nil {
			return l
		}

		l.Logger = slog.New(&LevelHandler{level.Level(), handler})
	}

	return l
}

// WithConfig returns a new logger if there's a config for it in configs else the current one.
// It adds the fields from the current logger.
func (l Logger) WithConfig(sections []string, configs types.ConfigData, opts ...Option) (Logger, error) {
	if !config.HasKey(append(sections, DefaultConfigSection), "plugin", configs) {
		return l, nil
	}

	newLogger, err := NewConfigDatas(sections, configs, opts...)
	if err != nil {
		return Logger{}, err
	}

	return newLogger.With(l.fields...), nil
}

// WithOpts returns a new logger with the given opt's.
// It adds the fields from the current logger.
func (l Logger) WithOpts(opts ...Option) (Logger, error) {
	nl, err := New(opts...)
	if err != nil {
		return Logger{}, err
	}

	return nl.With(l.fields...), nil
}

// With returns a new Logger that includes the given arguments, converted to
// Attrs as in [Logger.Log]. The Attrs will be added to each output from the
// Logger.
//
// The new Logger's handler is the result of calling WithAttrs on the receiver's
// handler.
func (l Logger) With(args ...any) Logger {
	l.fields = append(l.fields, args...)
	l.Logger = l.Logger.With(args...)

	return l
}

// Start no-op.
func (l Logger) Start() error {
	return nil
}

// Stop stops all cached plugins if this is the default logger.
func (l Logger) Stop(ctx context.Context) error {
	if !l.config.SetDefault {
		return nil
	}

	for p, pp := range pluginsCache.All() {
		if err := pp.Stop(ctx); err != nil {
			slog.Error("stopping a logger plugin", "plugin", p, "error", err)
		}
	}

	return nil
}

// String returns current handler plugin used.
func (l Logger) String() string {
	return l.config.Plugin
}

// Type returns the component type.
func (l Logger) Type() string {
	return ComponentType
}

// Level returns the level as int.
func (l Logger) Level() slog.Level {
	return l.config.Level
}

// Trace logs at TraceLevel.
func (l *Logger) Trace(msg string, args ...any) {
	l.Log(context.Background(), LevelTrace, msg, args...)
}
