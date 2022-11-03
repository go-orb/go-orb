package log

import (
	"github.com/go-orb/config/source/cli"
	"github.com/go-orb/config/util/container"
	"golang.org/x/exp/slog"
)

func init() {
	flag := cli.NewFlag(
		"logger",
		DefaultPlugin,
		cli.CPSlice([]string{"logger", "plugin"}),
		cli.Usage("Default logger to use, jsonstderr, jsonstdout, textstderr, textsdout."),
	)

	if err := cli.Flags.Add(flag); err != nil {
		panic(err)
	}
}

// Plugins is the registry for Logger plugins.
var Plugins = container.NewMap[func(level slog.Leveler) (slog.Handler, error)]() //nolint:gochecknoglobals
