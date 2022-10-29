package registry

import (
	"github.com/orb-org/orb/util/container"
)

var Plugins = container.NewMap[func() Registry]()
