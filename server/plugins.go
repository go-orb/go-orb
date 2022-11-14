package server

import (
	"go-micro.dev/v5/util/container"
)

//nolint:gochecknoglobals
var (
	// Plugins is the plugins container for registry.
	Plugins = container.NewMap[ProviderFunc]()
	// NewDefaults is the factory container for defaults.
	NewDefaults = container.NewMap[NewDefault]()
)
