package log

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slog"
)

// Names for common Levels.
const (
	LevelTrace slog.Level = slog.LevelDebug - 1
	LevelDebug slog.Level = slog.LevelDebug
	LevelInfo  slog.Level = slog.LevelInfo
	LevelWarn  slog.Level = slog.LevelWarn
	LevelError slog.Level = slog.LevelError
)

// ParseLevel parses a string level to an Level.
func ParseLevel(l string) (slog.Level, error) {
	switch strings.ToUpper(l) {
	case "TRACE":
		return LevelTrace, nil
	case "DEBUG":
		return LevelDebug, nil
	case "INFO":
		return LevelInfo, nil
	case "WARN":
		return LevelWarn, nil
	case "ERROR":
		return LevelError, nil
	default:
		return LevelInfo, fmt.Errorf("parselevel: unknown level %s", l)
	}
}
