package log

import (
	"os"

	"github.com/go-orb/config/source/cli"
	"github.com/go-orb/config/util/container"
	"golang.org/x/exp/slog"
)

func init() {
	err := cli.Flags.Add(cli.NewFlag(
		"logger",
		DefaultPlugin,
		cli.CPSlice([]string{"logger", "plugin"}),
		cli.Usage("Default logger to use, jsonstderr, jsonstdout, textstderr, textsdout."),
	))
	if err != nil {
		panic(err)
	}

	if err := Plugins.Add("jsonstdout", JSONStdoutPlugin); err != nil {
		panic(err)
	}

	if err := Plugins.Add("jsonstderr", JSONStderrPlugin); err != nil {
		panic(err)
	}

	if err := Plugins.Add("textstdout", TextStdoutPlugin); err != nil {
		panic(err)
	}

	if err := Plugins.Add("textstderr", TextStderrPlugin); err != nil {
		panic(err)
	}
}

// Plugins is the registry for Logger plugins.
var Plugins = container.NewMap[func(level slog.Leveler) (slog.Handler, error)]() //nolint:gochecknoglobals

// JSONStdoutPlugin writes json to stdout.
func JSONStdoutPlugin(level slog.Leveler) (slog.Handler, error) {
	return slog.HandlerOptions{Level: level}.NewJSONHandler(os.Stdout), nil
}

// JSONStderrPlugin writes json to stderr.
func JSONStderrPlugin(level slog.Leveler) (slog.Handler, error) {
	return slog.HandlerOptions{Level: level}.NewJSONHandler(os.Stderr), nil
}

// TextStdoutPlugin writes text to stdout.
func TextStdoutPlugin(level slog.Leveler) (slog.Handler, error) {
	return slog.HandlerOptions{Level: level}.NewTextHandler(os.Stdout), nil
}

// TextStderrPlugin writes text to stderr.
func TextStderrPlugin(level slog.Leveler) (slog.Handler, error) {
	return slog.HandlerOptions{Level: level}.NewTextHandler(os.Stderr), nil
}
