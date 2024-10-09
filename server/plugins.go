package server

import (
	"github.com/go-orb/go-orb/util/container"
)

//nolint:gochecknoglobals
var (
	// Plugins is the plugins container for registry.
	Plugins = container.NewMap[string, EntrypointProvider]()

	// PluginsNew holds the New function of Entrypoints, it's here
	// to create entrypoints from given configs.
	PluginsNew = container.NewMap[string, EntrypointNew]()

	// Handlers is a container of registration functions that can be used to
	// dynamically configure entrypoints.
	//
	// You need to register your handlers with a registration function to use it.
	// Example:
	//     Handlers.Register("myHandler", RegisterEchoHandler)
	Handlers = container.NewMap[string, RegistrationFunc]()
)
