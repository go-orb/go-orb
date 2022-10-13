package log

import "github.com/rs/zerolog/log"

// GlobalLogger is the global logger.
var GlobalLogger Logger //nolint:gochecknoglobals

func init() {
	GlobalLogger = &zeroLogger{}

	err := GlobalLogger.Init(NewConfig(), WithInternalParent(log.Logger))
	if err != nil {
		panic(err)
	}
}

// Trace creates a new event with log level trace.
func Trace() Event { return GlobalLogger.Trace() }

// Debug creates a new event with log level debug.
func Debug() Event { return GlobalLogger.Debug() }

// Info creates a new event with log level info.
func Info() Event { return GlobalLogger.Info() }

// Warn creates a new event with log level warn.
func Warn() Event { return GlobalLogger.Warn() }

// Err creates a new event with log level error.
func Err() Event { return GlobalLogger.Err() }

// Fatal creates a new event with log level fatal.
func Fatal() Event { return GlobalLogger.Fatal() }

// Panic creates a new event with log level panic.
func Panic() Event { return GlobalLogger.Panic() }
