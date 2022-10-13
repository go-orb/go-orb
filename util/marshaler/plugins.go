package marshaler

import (
	"jochum.dev/orb/orb/util/container"
)

// Plugins is the marshaler plugin container.
var Plugins = container.New(func() Marshaler { return nil }) //nolint:gochecknoglobals
