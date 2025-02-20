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
	LevelTrace  slog.Level = slog.LevelDebug - 1
	LevelDebug  slog.Level = slog.LevelDebug
	LevelInfo   slog.Level = slog.LevelInfo
	LevelWarn   slog.Level = slog.LevelWarn
	LevelNotice slog.Level = slog.LevelWarn - 2
	LevelError  slog.Level = slog.LevelError
	LevelFatal  slog.Level = slog.LevelError + 4
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
	case "NOTICE":
		return LevelNotice, nil
	case "ERROR":
		return LevelError, nil
	case "FATAL":
		return LevelFatal, nil
	default:
		return LevelInfo, fmt.Errorf("parselevel: unknown level %s", l)
	}
}
