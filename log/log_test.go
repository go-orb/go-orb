package log

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestChangeLevel(t *testing.T) {
	l, err := New(NewConfig())
	require.NoError(t, err)

	dCfg := NewConfig()
	dCfg.Level = LevelDebug
	lDebug, err := New(dCfg)
	require.NoError(t, err)

	l.Info("Default logger Test")
	l.Debug("Not shown")
	lDebug.Debug("Debug: logger test")
	lDebug.Log(context.TODO(), LevelTrace, "Debug: Trace test")
}

func TestComponentLogger(t *testing.T) {
	l, err := New(NewConfig())
	require.NoError(t, err)

	l.Info("Message One")

	l2, err := l.WithComponent("broker", "nats", "", LevelDebug)
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
	return slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level}), nil
}
