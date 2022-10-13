package log

import (
	"fmt"

	"github.com/rs/zerolog"
)

func newZeroEvent(config Config, l zerolog.Logger, level zerolog.Level) Event {
	var event *zeroEvent

	switch level {
	case zerolog.TraceLevel:
		event = &zeroEvent{z: l.Trace(), level: level}
	case zerolog.DebugLevel:
		event = &zeroEvent{z: l.Debug(), level: level}
	case zerolog.InfoLevel:
		event = &zeroEvent{z: l.Info(), level: level}
	case zerolog.WarnLevel:
		event = &zeroEvent{z: l.Warn(), level: level}
	case zerolog.ErrorLevel:
		event = &zeroEvent{z: l.Error().Caller(2), level: level}
	case zerolog.FatalLevel:
		event = &zeroEvent{z: l.Fatal().Caller(2), level: level}
	case zerolog.PanicLevel:
		event = &zeroEvent{z: l.Panic().Caller(2), level: level}
	case zerolog.NoLevel:
		event = &zeroEvent{z: l.Log(), level: level}
	case zerolog.Disabled:
		event = &zeroEvent{z: nil, level: level}
	default:
		event = &zeroEvent{z: nil, level: level}
	}

	if config.Fields() != nil {
		event.Fields(config.Fields())
	}

	return event
}

type zeroEvent struct {
	level zerolog.Level
	z     *zerolog.Event
}

func (e *zeroEvent) Enabled() bool {
	return e.z != nil && e.level != zerolog.Disabled
}

func (e *zeroEvent) Discard() Event {
	e.z.Discard()
	return e
}

func (e *zeroEvent) Msg(msg string) {
	e.z.Msg(msg)
}

func (e *zeroEvent) Send() {
	e.z.Send()
}

func (e *zeroEvent) Msgf(msg string, v ...any) Event {
	e.z.Msgf(msg, v...)
	return e
}
func (e *zeroEvent) Fields(fields any) Event {
	e.z.Fields(fields)
	return e
}
func (e *zeroEvent) Strs(key string, vals []string) Event {
	e.z.Strs(key, vals)
	return e
}

func (e *zeroEvent) Stringer(key string, val fmt.Stringer) Event {
	e.z.Stringer(key, val)
	return e
}

func (e *zeroEvent) AnErr(key string, err error) Event {
	e.z.AnErr(key, err)
	return e
}

func (e *zeroEvent) Err(err error) Event {
	e.z.Err(err)
	return e
}
