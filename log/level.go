package log

import (
	"fmt"
	"strings"

	"log/slog"
)

// Names for common Levels.
// TODO(jochumdev):  Something like this would be nice
//
//	type LevelT interface {
//		slog.Level | string | constraints.Integer
//	}
const (
	// LevelTrace must be added, because [slog] package does not have one by default.
	// Generate it by subtracting 4 levels from [slog.Debug] following the example of
	// [slog.LevelWarn] and [slog.LevelError] which are set to 4 and 8.
	LevelTrace  slog.Level = slog.LevelDebug - 4
	LevelDebug  slog.Level = slog.LevelDebug
	LevelInfo   slog.Level = slog.LevelInfo
	LevelWarn   slog.Level = slog.LevelWarn
	LevelNotice slog.Level = slog.LevelWarn - 2
	LevelError  slog.Level = slog.LevelError
	LevelFatal  slog.Level = slog.LevelError + 4
)

// slogLevelToString converts a slog.Level to a string.
func slogLevelToString(l slog.Level) string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelNotice:
		return "NOTICE"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// stringToSlogLevel parses a string level to an Level.
func stringToSlogLevel(l string) slog.Level {
	switch strings.ToUpper(l) {
	case "TRACE":
		return LevelTrace
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "NOTICE":
		return LevelNotice
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		return stringToSlogLevel(DefaultLevel)
	}
}
