// Package log is the log component of Orb.
package log

import (
	"errors"
	"fmt"

	"jochum.dev/orb/orb/config/chelp"
)

// ErrSubLogger is returned on Init() when it's not possible to make a sublogger of the parent/internalParent.
var ErrSubLogger = errors.New("making a sublogger of a different parent logger is not possible")

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

	Init(config any, opts ...Option) error
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

// FromConfig converts "aConfig" into a Logger using "parent" as parent.
func FromConfig(aConfig any, parent Logger) (Logger, error) {
	if aConfig == nil {
		return parent, nil
	}

	config, ok := aConfig.(Config)
	if !ok {
		return nil, chelp.ErrUnknownConfig
	}

	pFunc, err := Plugins.Plugin(config.Plugin())
	if err != nil {
		return nil, err
	}

	p := pFunc()

	if parent.String() == config.Plugin() {
		if err := p.Init(aConfig, WithParent(parent)); err != nil {
			return nil, err
		}
	} else {
		if err := p.Init(aConfig); err != nil {
			return nil, err
		}
	}

	return p, nil
}
