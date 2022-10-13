package log

import (
	"errors"
	"fmt"
)

var ErrSubLoggerNotPossible = errors.New("making a sublogger of a different parent logger is not possible")

type Event interface {
	Enabled() bool
	Discard() Event
	Msg(msg string)
	Send()
	Msgf(msg string, v ...interface{}) Event
	Fields(fields interface{}) Event
	Strs(key string, vals []string) Event
	Stringer(key string, val fmt.Stringer) Event
	AnErr(key string, err error) Event
	Err(err error) Event
}

type Logger interface {
	fmt.Stringer

	Init(config any, parent Logger) error
	Config() any

	Level() string

	// Trace starts a new message with trace level.
	Trace() Event

	// Debug starts a new message with debug level.
	Debug() Event

	Info() Event

	Warn() Event

	Err() Event

	Fatal() Event

	Panic() Event
}
