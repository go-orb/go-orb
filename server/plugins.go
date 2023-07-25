package server

import (
	"github.com/go-orb/go-orb/util/container"
)

//nolint:gochecknoglobals
var (
	// Plugins is the plugins container for registry.
	Plugins = container.NewPlugins[ProviderFunc]()

	// NewDefaults is the factory container for defaults.
	NewDefaults = container.NewPlugins[NewDefault]()

	// Handlers is a container of registration functions that can be used to
	// dynamically configure entrypoints.
	//
	// You need to register your handlers with a registration function to use it.
	// Example:
	//     Handlers.Register("myHandler", RegisterEchoHandler)
	Handlers = container.NewMap[RegistrationFunc]()
)
