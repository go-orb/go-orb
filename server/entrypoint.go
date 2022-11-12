package server

import (
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/types"
	"go-micro.dev/v5/types/component"
)

// RegistrationFunc is executed to register a handler to a server (entrypoint)
// passed as srv. srv can be of any of the various server types, should be a pointer.
//
// You can write your own custom registration functions to register extra handlers.
//
// Inside the registration function, you need to convert the server type and
// assert that you are working with the server type you are expecting, and
// otherwise no-op. For an example, see the implementation of NewRegistrationFunc.
type RegistrationFunc func(srv any)

// EntrypointOption are functional options for entrypoints.
type EntrypointOption func(v any)

// NewDefault is a factory function type for entrypoint defaults, registerd by the plugins.
type NewDefault func() any

// Entrypoint is a server, and represents an entrypoint into the web.
type Entrypoint interface {
	component.Component

	// Register is used to register handlers.
	//
	// A registration function takes a pointer to the server, which can then
	// be used to register handlers in the server specific way.
	Register(RegistrationFunc)

	// Name returns the entrypoint name.
	Name() string
}

// TODO: check Option type, should probably be something like func(any) and then convert
// TODO: check return type here and in registry of providerfunc, interface?

// ProviderFunc is the function type to create a new entrypoint.
// ProviderFuncs are registered in the plugins container, and can be called
// at runtime depending on the configuration.
type ProviderFunc func(
	name string,
	service types.ServiceName,
	data types.ConfigData,
	logger log.Logger,
	opts ...EntrypointOption,
) (Entrypoint, error)

// EntrypointTemplate is the configuation used to create a single entrypoint.
//
// You will rarely need to manually create a template object, it will be done
// for you through the provided server options.
type EntrypointTemplate struct {
	// Type is the entrypoint type to use. To use a specific server type as entrypiont
	// the provider function needs to be registerd as an entrypoint plugin. This
	// is done by importing the package, typically done with a named import as _.
	Type string

	// Options are the options used to create the config. The default options
	// are used as starting point, to which this list of options will be applied.
	// The result will be your full entrypoint configuration.
	//
	// By default, a random port will be chosen for the entrypoint to listen on,
	// defined as ":0". For all options not specified, default values will be used.
	// TODO: do we really use :0 as address, or the v4 code to identify an interface?
	Options []EntrypointOption
}

// EntrypointTemplates is a collection of entrypoint templates.
// The map index is the entrypoint name.
//
// Each entrypoint needs a unique name, as each entrypoint can be dynamically
// configured by referencing the name. The default name used in an entrypoint
// is the format of "http-<uuid>", used if no custom name is provided.
type EntrypointTemplates map[string]EntrypointTemplate

// NewRegistrationFunc takes a registration function and handler and returns
// a registration func that can be used with one specific server type.
//
// This functions is not exepected to be used a lot, and is merely a utility function.
// This function could be useful if a user wants to register a non-micro project proto handler.
//
// This function is located in the core.
func NewRegistrationFunc[Tsrv any, Thandler any](register func(Tsrv, Thandler), handler Thandler) RegistrationFunc {
	return RegistrationFunc(func(s any) {
		sr, ok := s.(Tsrv)
		if !ok {
			// Maybe we should log here
			return
		}

		register(sr, handler)
	})
}
