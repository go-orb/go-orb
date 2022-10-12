package zerolog

import (
	"os"

	"github.com/rs/zerolog"
	"jochum.dev/orb/orb/log"
)

func init() {
	if err := log.Plugins.Add("zerolog", New, log.NewConfig); err != nil {
		panic(err)
	}
}

type Logger struct {
	L zerolog.Logger

	config log.Config
}

func New() log.Logger { return &Logger{} }

func (l *Logger) Init(config log.Config, parent log.Logger) error {
	level, err := zerolog.ParseLevel(config.Level())
	if err != nil {
		return err
	}

	if parent == nil {
		l.L = zerolog.New(os.Stderr).Level(level)
	} else if parent.String() != l.String() {
		return log.ErrSubLoggerNotPossible
	} else {
		l.L = parent.(*Logger).L.Level(level)
	}

	l.config = config

	return nil
}

func (l *Logger) Config() log.Config {
	return l.config
}

func (l *Logger) String() string {
	return "zerolog"
}

func (l *Logger) Level() string {
	return l.config.Level()
}

func (l *Logger) Trace() log.Event { return newEvent(l.L, zerolog.TraceLevel) }
func (l *Logger) Debug() log.Event { return newEvent(l.L, zerolog.DebugLevel) }
func (l *Logger) Info() log.Event  { return newEvent(l.L, zerolog.InfoLevel) }
func (l *Logger) Warn() log.Event  { return newEvent(l.L, zerolog.WarnLevel) }
func (l *Logger) Err() log.Event   { return newEvent(l.L, zerolog.ErrorLevel) }
func (l *Logger) Fatal() log.Event { return newEvent(l.L, zerolog.FatalLevel) }
func (l *Logger) Panic() log.Event { return newEvent(l.L, zerolog.PanicLevel) }
