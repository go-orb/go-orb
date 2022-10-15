package log

import (
	"os"

	"github.com/rs/zerolog"
	"jochum.dev/orb/orb/config"
)

func init() {
	if err := Plugins.Add(
		"zerolog",
		func() Logger { return &zeroLogger{} },
		func() any { return NewComponentConfig() },
	); err != nil {
		panic(err)
	}
}

type zeroLogger struct {
	L zerolog.Logger

	config Config
}

func (l *zeroLogger) Init(aConfig any, opts ...Option) error {
	if cfg, ok := aConfig.(Config); ok {
		l.config = cfg
	} else {
		return config.ErrUnknownConfig
	}

	level, err := zerolog.ParseLevel(l.config.GetLevel())
	if err != nil {
		return err
	}

	// Options handling
	options := NewOptions(opts...)
	switch il := options.InternalParent.(type) {
	case nil:
		switch options.Parent.String() {
		case "":
			l.L = zerolog.New(os.Stderr).Level(level)
		case l.String():
			l.L = options.Parent.(*zeroLogger).L.Level(level)
		default:
			return ErrSubLogger
		}
	case zerolog.Logger:
		l.L = il.Level(level)
	default:
		return ErrSubLogger
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
	return l.config.GetLevel()
}

func (l *zeroLogger) Trace() Event {
	if !l.should(zerolog.TraceLevel) {
		return nil
	}

	return newZeroEvent(l.config, l.L, zerolog.TraceLevel)
}

func (l *zeroLogger) Debug() Event {
	if !l.should(zerolog.DebugLevel) {
		return nil
	}

	return newZeroEvent(l.config, l.L, zerolog.DebugLevel)
}

func (l *zeroLogger) Info() Event {
	if !l.should(zerolog.InfoLevel) {
		return nil
	}

	return newZeroEvent(l.config, l.L, zerolog.InfoLevel)
}

func (l *zeroLogger) Warn() Event {
	if !l.should(zerolog.WarnLevel) {
		return nil
	}

	return newZeroEvent(l.config, l.L, zerolog.WarnLevel)
}

func (l *zeroLogger) Err() Event {
	if !l.should(zerolog.ErrorLevel) {
		return nil
	}

	return newZeroEvent(l.config, l.L, zerolog.ErrorLevel)
}

func (l *zeroLogger) Fatal() Event {
	return newZeroEvent(l.config, l.L, zerolog.FatalLevel)
}

func (l *zeroLogger) Panic() Event {
	return newZeroEvent(l.config, l.L, zerolog.PanicLevel)
}

func (l *zeroLogger) should(lvl zerolog.Level) bool {
	if lvl < l.L.GetLevel() || lvl < zerolog.GlobalLevel() {
		return false
	}

	return true
}
