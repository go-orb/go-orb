package registry

import (
	"jochum.dev/orb/orb/util/container"
)

const Name = "registry"

var Plugins = container.NewPlugins(
	func() Registry { return nil }, // Plugin factory
	func() any { return nil },      // Config factory
)
