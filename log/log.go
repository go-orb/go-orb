// Package log contains a golang.org/x/exp/slog compatible logger.
package log

import (
	"fmt"
	"log"

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
		return nil, err
	}

	handler, err := p(TraceLevel)
	if err != nil {
		return nil, fmt.Errorf("create new plugin handler: %w", err)
	}

	return handler, nil
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
	l := Logg  er{
		plugin: cfg.Plugin,
		plugins: make(map[string]slog.Handler),
	}

	level, err := ParseLevel(cfg.Level)
	if err != nil {
		return Logger{}, err
	}

	h, err := l.Plugin(cfg.Plugin)
	if err != nil {
		return Logger{}, err
	}

	h = NewLevelHandler(level, h)

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

func NewComponentLogger(l log.Logger, component, name, plugin, level string) ( log.Logger, error ) {
	lvl, err := ParseLevel(level)
	if err != nil {
		l.LogDepth(1, ErrorLevel, "invalid log level provided", err)
	}

	handler, err := l.Plugin(plugin)
	if err != nil {
		l.LogDepth(1, ErrorLevel, "invalid log level provided", err)
	}

	// Optionally avoid wrapping a handler if the level is the same as the parent
	// logger. To check the handler level it needs to implent the Leveler interface,
	// which is not provided by default on slog handlers.
	noWrapper := false
	if level, ok := l.Handler().(slog.Leveler); ok && level == lvl {
		noWrapper = true
	}

	handler := l.Plugin(plugin)

	handler := LevelHandler{level: lvl, handler: handler}

	ctx := l.With(
		slog.String("component", component),
		slog.String("plugin", name),
	).Context()
	l = slog.New(handler).WithContext(ctx)

	return l
}
