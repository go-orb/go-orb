package log

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slog"
)

// Names for common Levels.
const (
	TraceLevel slog.Level = slog.DebugLevel - 1
	DebugLevel slog.Level = slog.DebugLevel
	InfoLevel  slog.Level = slog.InfoLevel
	WarnLevel  slog.Level = slog.WarnLevel
	ErrorLevel slog.Level = slog.ErrorLevel
)

// ParseLevel parses a string level to an Level.
func ParseLevel(l string) (slog.Level, error) {
	switch strings.ToUpper(l) {
	case "TRACE":
		return TraceLevel, nil
	case "DEBUG":
		return DebugLevel, nil
	case "INFO":
		return InfoLevel, nil
	case "WARN":
		return WarnLevel, nil
	case "ERROR":
		return ErrorLevel, nil
	default:
		return InfoLevel, fmt.Errorf("unknown level %s", l)
	}
}
