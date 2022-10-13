package log

import (
	"os"

	"github.com/rs/zerolog"
	"jochum.dev/orb/orb/config/chelp"
)

func init() {
	if err := Plugins.Add("zerolog", newZero, NewConfig); err != nil {
		panic(err)
	}
}

type zeroLogger struct {
	L zerolog.Logger

	config *BaseConfig
}

func newZero() Logger { return &zeroLogger{} }

func (l *zeroLogger) Init(aConfig any, parent Logger) error {
	switch tConfig := aConfig.(type) {
	case *BaseConfig:
		l.config = tConfig
	default:
		return chelp.ErrUnknownConfig
	}

	level, err := zerolog.ParseLevel(l.config.Level())
	if err != nil {
		return err
	}

	if parent == nil {
		l.L = zerolog.New(os.Stderr).Level(level)
	} else if parent.String() != l.String() {
		return ErrSubLoggerNotPossible
	} else {
		l.L = parent.(*zeroLogger).L.Level(level)
	}

	return nil
}

func (l *zeroLogger) Config() any {
	return l.config
}

func (l *zeroLogger) String() string {
	return "zerolog"
}

func (l *zeroLogger) Level() string {
	return l.config.Level()
}

func (l *zeroLogger) Trace() Event { return newZeroEvent(l.L, zerolog.TraceLevel) }
func (l *zeroLogger) Debug() Event { return newZeroEvent(l.L, zerolog.DebugLevel) }
func (l *zeroLogger) Info() Event  { return newZeroEvent(l.L, zerolog.InfoLevel) }
func (l *zeroLogger) Warn() Event  { return newZeroEvent(l.L, zerolog.WarnLevel) }
func (l *zeroLogger) Err() Event   { return newZeroEvent(l.L, zerolog.ErrorLevel) }
func (l *zeroLogger) Fatal() Event { return newZeroEvent(l.L, zerolog.FatalLevel) }
func (l *zeroLogger) Panic() Event { return newZeroEvent(l.L, zerolog.PanicLevel) }
