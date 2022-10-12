package registry

import (
	"jochum.dev/orb/orb/util/container"
)

var Plugins = container.NewPlugins(
	func() Registry { return nil }, // Plugin factory
	func() Config { return nil },   // Config factory
)
