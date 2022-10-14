package configsource

import (
	"jochum.dev/orb/orb/util/container"
)

// Plugins is the configsource plugin container.
var Plugins = container.New(func() ConfigSource { return nil }) //nolint:gochecknoglobals
