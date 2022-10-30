package registry

import (
	"github.com/go-orb/orb/util/container"
)

var Plugins = container.NewMap[func() Registry]()
