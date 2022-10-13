package log

import (
	"fmt"

	"github.com/rs/zerolog"
)

func newZeroEvent(l zerolog.Logger, level zerolog.Level) Event {
	switch level {
	case zerolog.TraceLevel:
		return &zeroEvent{z: l.Trace(), level: level}
	case zerolog.DebugLevel:
		return &zeroEvent{z: l.Debug(), level: level}
	case zerolog.InfoLevel:
		return &zeroEvent{z: l.Info(), level: level}
	case zerolog.WarnLevel:
		return &zeroEvent{z: l.Warn(), level: level}
	case zerolog.ErrorLevel:
		return &zeroEvent{z: l.Error().Caller(2), level: level}
	case zerolog.FatalLevel:
		return &zeroEvent{z: l.Fatal().Caller(2), level: level}
	case zerolog.PanicLevel:
		return &zeroEvent{z: l.Panic().Caller(2), level: level}
	case zerolog.NoLevel:
		return &zeroEvent{z: l.Log(), level: level}
	case zerolog.Disabled:
		return &zeroEvent{z: nil, level: level}
	default:
		return &zeroEvent{z: nil, level: level}
	}
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

func (e *zeroEvent) Msgf(msg string, v ...interface{}) Event {
	e.z.Msgf(msg, v...)
	return e
}
func (e *zeroEvent) Fields(fields interface{}) Event {
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
