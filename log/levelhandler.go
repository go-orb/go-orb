package log

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
)

var _ slog.Handler = (*LevelHandler)(nil)

var (
	// ErrNoHandler happens when a the LevelHandler wrapper gets no handler.
	ErrNoHandler = errors.New("no handler defined")
)

// LevelHandler is wrapper for slog.Handler which does Leveling.
type LevelHandler struct {
	level   slog.Level
	handler slog.Handler
}

// NewLevelHandler implements slog.Handler interface. It is used to wrap a
// handler with a new log level. As log level cannot be modified within a handler
// through the interface of slog, you can use this to wrap a handler with a new
// log level.
func NewLevelHandler(level slog.Level, h slog.Handler) (*LevelHandler, error) {
	if h == nil {
		return nil, ErrNoHandler
	}

	return &LevelHandler{level, h}, nil
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
// Enabled is called early, before any arguments are processed,
// to save effort if the log event should be discarded.
func (h *LevelHandler) Enabled(context context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle handles the Record.
// It will only be called if Enabled returns true.
// Handle methods that produce output should observe the following rules:
//   - If r.Time is the zero time, ignore the time.
//   - If an Attr's key is the empty string, ignore the Attr.
func (h *LevelHandler) Handle(context context.Context, r slog.Record) error {
	if h.handler == nil {
		return ErrNoHandler
	}

	return h.handler.Handle(context, r)
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LevelHandler{h.level, h.handler.WithAttrs(attrs)}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	return &LevelHandler{h.level, h.handler.WithGroup(name)}
}
