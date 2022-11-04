package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestChangeLevel(t *testing.T) {
	l, err := New(NewConfig())
	require.NoError(t, err)

	dCfg := NewConfig()
	dCfg.Level = "debug"
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
}

func init() {
	if err := Plugins.Add("textstderr", NewHandlerStderr); err != nil {
		panic(err)
	}
}

// NewHandlerStderr writes text to stderr.
func NewHandlerStderr(level slog.Leveler) (slog.Handler, error) {
	return slog.HandlerOptions{Level: level}.NewTextHandler(os.Stderr), nil
}
