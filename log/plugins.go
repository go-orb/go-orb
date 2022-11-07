package log

import (
	"github.com/go-orb/config/source/cli"
	"golang.org/x/exp/slog"

	"go-micro.dev/v5/util/container"
)

type pluginHandler struct {
	handler slog.Handler
	level   slog.Level
}

// Plugins is the registry for Logger plugins.
var Plugins = container.NewMap[func(level slog.Leveler) (slog.Handler, error)]() //nolint:gochecknoglobals

// plugins is a cache of lazyloaded plugin handlers.
// In order to prevent creating multiple handlers, and thus potentially
// multiple connections, depending on the handler, we cache the handlers, and
// wrap them with a LevelHandler by default. This way we only create one
// handler per plugin, for use in any amount of loggers.
var plugins = container.NewSafeMap[pluginHandler]()

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
