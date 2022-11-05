package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"

	"go-micro.dev/v5/types/component"
)

func TestChangeLevel(t *testing.T) {
	l, err := New(NewConfig())
	require.NoError(t, err)

	dCfg := NewConfig()
	dCfg.Level = DebugLevel
	lDebug, err := New(dCfg)
	require.NoError(t, err)

	l.Info("Default logger Test")
	l.Debug("Not shown")
	lDebug.Debug("Debug: logger test")
	lDebug.Log(TraceLevel, "Debug: Trace test")
}

func TestComponentLogger(t *testing.T) {
	l, err := New(NewConfig())
	require.NoError(t, err)

	l.Info("Message One")

	l2, err := NewComponentLogger(l, component.Type("broker"), "nats", "", "debug")
	require.NoError(t, err)

	l2.Info("Message Two")
	l2.Debug("Debug Two")
}

func init() {
	if err := Plugins.Add("textstderr", NewHandlerStderr); err != nil {
		panic(err)
	}
}

// NewHandlerStderr writes text to stderr.
func NewHandlerStderr(level slog.Leveler) (slog.Handler, error) {
	return slog.HandlerOptions{Level: level}.NewJSONHandler(os.Stderr), nil
}
