package cli

import (
	"jochum.dev/jochumdev/orb/util/container"
)

var Plugins = container.New(func(opts ...Option) Cli { return nil })
