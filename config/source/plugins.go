package source

import (
	"go-micro.dev/v5/util/container"
)

// Plugins is the configsource plugin container.
var Plugins = container.NewList[Source]() //nolint:gochecknoglobals
