// Package log contains a golang.org/x/exp/slog compatible logger.
package log

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/config"

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

// New creates a new Logger from a Config.
func New(opts ...Option) (Logger, error) {
	return NewConfigDatas([]string{}, nil, opts...)
}

// NewConfigDatas will create a new logger with the given configs,
// as well as the given fields.
func NewConfigDatas(sections []string, configs map[string]any, opts ...Option) (Logger, error) {
	// Initialize configuration
	cfg := NewConfig(opts...)

	// Parse configuration
	if configs == nil {
		var err error
		data, err := config.ParseStruct(append(sections, DefaultConfigSection), &cfg)

		if err != nil {
			return Logger{}, fmt.Errorf("while creating a new config: %w", err)
		}

		configs = data
	} else if err := config.Parse(sections, DefaultConfigSection, configs, &cfg); err != nil {
		return Logger{}, fmt.Errorf("while creating a new config: %w", err)
	}

	// Get the plugin
	pf, ok := plugins.Get(cfg.Plugin)
	if !ok {
		return Logger{}, fmt.Errorf("while getting the log plugin '%s'", cfg.Plugin)
	}

	// Get or create the provider
	provider, err := pf(sections, configs, opts...)
	if err != nil {
		return Logger{}, err
	}

	// Get provider from cache or initialize it
	cachedProvider, ok := pluginsCache.Get(provider.Key())
	if !ok {
		if err := provider.Start(); err != nil {
			return Logger{}, fmt.Errorf("unknown provider '%s'", provider.Key())
		}

		pluginsCache.Set(cfg.Plugin, provider)
		cachedProvider = provider
	}

	// Get handler and create level handler
	handler, err := cachedProvider.Handler()
	if err != nil {
		return Logger{}, err
	}

	lvlHandler, err := NewLevelHandler(stringToSlogLevel(cfg.Level), handler)
	if err != nil {
		return Logger{}, err
	}

	// Create logger with fields
	fields := make([]any, 0, len(cfg.Fields)*2)
	for k, v := range cfg.Fields {
		fields = append(fields, slog.Any(k, v))
	}

	logger := Logger{
		Logger:         slog.New(lvlHandler),
		pluginProvider: cachedProvider,
		config:         cfg,
		fields:         fields,
	}

	if len(fields) > 0 {
		logger.Logger = logger.Logger.With(fields...)
	}

	return logger, nil
}

// Provide provides a new logger.
// It will set the slog.Logger as package wide default logger.
func Provide(
	_ cli.ServiceContextHasConfigData,
	svcCtx *cli.ServiceContext,
	components *types.Components,
	opts ...Option,
) (Logger, error) {
	logger, err := NewConfigDatas([]string{}, svcCtx.Config, opts...)
	if err != nil {
		return Logger{}, err
	}

	slog.SetDefault(logger.Logger)

	// Register the logger as a component.
	_ = components.Add(logger, types.PriorityLogger) //nolint:errcheck

	return logger, nil
}

// ProvideNoOpts provides a new logger without options.
func ProvideNoOpts(
	hasConfig cli.ServiceContextHasConfigData,
	svcCtx *cli.ServiceContext,
	components *types.Components,
) (Logger, error) {
	return Provide(hasConfig, svcCtx, components)
}

// ProvideWithServiceNameField provides a new logger with the service name field.
func ProvideWithServiceNameField(
	hasConfig cli.ServiceContextHasConfigData,
	svcCtx *cli.ServiceContext,
	components *types.Components,
) (Logger, error) {
	return Provide(hasConfig, svcCtx, components, WithFields(map[string]any{"service": svcCtx.Name()}))
}

// WithLevel creates a copy of the logger with a new level.
// It will inherit all the fields and the context from the parent logger.
func (l Logger) WithLevel(level string) Logger {
	if level != "" {
		l.config.Level = level

		handler, err := l.pluginProvider.Handler()
		if err != nil {
			return l
		}

		l.Logger = slog.New(&LevelHandler{stringToSlogLevel(level), handler})
	}

	return l
}

// WithConfig returns a new logger if there's a config for it in configs else the current one.
// It adds the fields from the current logger.
func (l Logger) WithConfig(sections []string, configs map[string]any, opts ...Option) (Logger, error) {
	if !config.HasKey[string](sections, DefaultConfigSection, configs) {
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
func (l Logger) Start(_ context.Context) error {
	return nil
}

// Stop stops all cached plugins if this is the default logger.
func (l Logger) Stop(ctx context.Context) error {
	if !l.config.SetDefault {
		return nil
	}

	pluginsCache.Range(func(p string, pp ProviderType) bool {
		if err := pp.Stop(ctx); err != nil {
			slog.Error("stopping a logger plugin", "plugin", p, "error", err)
		}

		return true
	})

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
	return stringToSlogLevel(l.config.Level)
}

// Trace logs at TraceLevel.
func (l Logger) Trace(msg string, args ...any) {
	l.Log(context.Background(), LevelTrace, msg, args...)
}

// TraceContext logs with context.Context.
func (l Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelTrace, msg, args...)
}
