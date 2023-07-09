package source

import (
	"github.com/go-orb/go-orb/util/container"
)

// Plugins is the configsource plugin container.
var Plugins = container.NewList[Source]() //nolint:gochecknoglobals
