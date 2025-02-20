package log

// These functions are global functions copied over from the slog library.

import (
	"context"

	"log/slog"
)

// Trace calls Logger.Trace on the default logger.
func Trace(msg string, args ...any) {
	slog.Default().Log(context.TODO(), LevelTrace, msg, args...)
}

// Debug calls Logger.Debug on the default logger.
func Debug(msg string, args ...any) {
	slog.Default().Log(context.TODO(), LevelDebug, msg, args...)
}

// Info calls Logger.Info on the default logger.
func Info(msg string, args ...any) {
	slog.Default().Log(context.TODO(), LevelInfo, msg, args...)
}

// Warn calls Logger.Warn on the default logger.
func Warn(msg string, args ...any) {
	slog.Default().Log(context.TODO(), LevelWarn, msg, args...)
}

// Error calls Logger.Error on the default logger.
func Error(msg string, args ...any) {
	slog.Default().Log(context.TODO(), LevelError, msg, args...)
}

// Log calls Logger.Log on the default logger.
func Log(level slog.Level, msg string, args ...any) {
	slog.Default().Log(context.TODO(), level, msg, args...)
}

// LogAttrs calls Logger.LogAttrs on the default logger.
func LogAttrs(level slog.Level, msg string, attrs ...slog.Attr) { //nolint:revive
	slog.Default().LogAttrs(context.TODO(), level, msg, attrs...)
}
