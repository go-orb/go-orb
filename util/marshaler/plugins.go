package marshaler

import (
	"jochum.dev/orb/orb/util/container"
)

var Plugins = container.New(func() Marshaler { return nil })
