package server

import (
	"go-micro.dev/v5/util/container"
)

// Plugins is the plugins container for registry.
//
//nolint:gochecknoglobals
var Plugins = container.NewMap[ProviderFunc]()
var NewDefaults = container.NewMap[NewDefault]()
