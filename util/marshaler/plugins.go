package marshaler

import (
	"jochum.dev/jochumdev/orb/util/container"
)

var Plugins = container.New(func() Marshaler { return nil })
