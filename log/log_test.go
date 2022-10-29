package log

import (
	"testing"

	"github.com/stretchr/testify/require"
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
