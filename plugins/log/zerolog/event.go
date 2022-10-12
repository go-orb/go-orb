package zerolog

import (
	"fmt"

	"github.com/rs/zerolog"
	"jochum.dev/orb/orb/log"
)

func newEvent(l zerolog.Logger, level zerolog.Level) log.Event {
	switch level {
	case zerolog.TraceLevel:
		return &Event{z: l.Trace(), level: level}
	case zerolog.DebugLevel:
		return &Event{z: l.Debug(), level: level}
	case zerolog.InfoLevel:
		return &Event{z: l.Info(), level: level}
	case zerolog.WarnLevel:
		return &Event{z: l.Warn(), level: level}
	case zerolog.ErrorLevel:
		return &Event{z: l.Error(), level: level}
	case zerolog.FatalLevel:
		return &Event{z: l.Fatal(), level: level}
	case zerolog.PanicLevel:
		return &Event{z: l.Panic(), level: level}
	case zerolog.NoLevel:
		return &Event{z: l.Log(), level: level}
	case zerolog.Disabled:
		return &Event{z: nil, level: level}
	default:
		return &Event{z: nil, level: level}
	}
}

type Event struct {
	level zerolog.Level
	z     *zerolog.Event
}

func (e *Event) Enabled() bool {
	return e.z != nil && e.level != zerolog.Disabled
}

func (e *Event) Discard() log.Event {
	e.z.Discard()
	return e
}

func (e *Event) Msg(msg string) {
	e.z.Msg(msg)
}

func (e *Event) Send() {
	e.z.Send()
}

func (e *Event) Msgf(msg string, v ...interface{}) log.Event {
	e.z.Msgf(msg, v...)
	return e
}
func (e *Event) Fields(fields interface{}) log.Event {
	e.z.Fields(fields)
	return e
}
func (e *Event) Strs(key string, vals []string) log.Event {
	e.z.Strs(key, vals)
	return e
}

func (e *Event) Stringer(key string, val fmt.Stringer) log.Event {
	e.z.Stringer(key, val)
	return e
}
