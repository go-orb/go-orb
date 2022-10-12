package log

import (
	"jochum.dev/orb/orb/util/container"
)

var Plugins = container.NewPlugins(
	func() Logger { return nil }, // Plugin factory
	func() Config { return nil }, // Config factory
)
