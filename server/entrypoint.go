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

// HandlerRegistrations type is a map of regirstration functions for an entrypoint.
// A custom type is used to manually define the json/yaml unmarshal behavior,
// as we don't want to overwrite the list, we only want to add on to it with
// a config file.
type HandlerRegistrations map[string]RegistrationFunc

// NewDefault is a factory function type for entrypoint defaults, registered by
// the plugins.
type NewDefault func() EntrypointConfig

// Entrypoint is a server, and represents an entrypoint into the web.
type Entrypoint interface {
	component.Component

	// Register is used to register handlers.
	//
	// A registration function takes a pointer to the server, which can then
	// be used to register handlers in the server specific way.
	Register(RegistrationFunc)

	// Address returns the address the entrypoint is listening on.
	Address() string
}

// EntrypointConfig provides a primitive way to constrain entrypoint config
// types. It should be implemented by every server plugin.
//
// This interface is a hack around the fact that you cannot create custom type
// constraints on common struct fields, as described in golang/go/issues/48522.
// Once this issue is solved, this interface should be replaced in favor of
// which ever new semantics get introduced.
type EntrypointConfig interface {
	// TODO: as long as https://github.com/golang/go/issues/48522 is open, we need
	// this interface. But should be removed after in favor of some generic
	// constraint to identify entrypoint configs and access common fields.
	GetAddress() string
	Copy() EntrypointConfig
}

// ProviderFunc is the function type to create a new entrypoint.
// ProviderFuncs are registered in the plugins container, and can be called
// at runtime depending on the configuration.
type ProviderFunc func(
	service types.ServiceName,
	logger log.Logger,
	// config is the entrypoint plugin config. Here it is passed as an any to
	// allow any config type to be passed through. The entrypoint provider should
	// convert the any back into its own type, and error on type mismatch.
	config any,
) (Entrypoint, error)

// EntrypointTemplate is the configuation used to create a single entrypoint.
//
// You will rarely need to manually create a template object, it will be done
// for you through the provided server options.
type EntrypointTemplate struct {
	Enabled bool

	// Type is the entrypoint type to use. To use a specific server type as
	// entrypiont the provider function needs to be registered as an entrypoint
	// plugin. This is done by importing the package, typically done with a named
	// import as _.
	Type string

	// Config is the configuration used to create the entrypoint. The default
	// options are used as starting point, to which this list of options will be
	// applied as provided through both the fuctional options and any file config.
	// The result will be your full entrypoint configuration.
	//
	// By default, a random port will be chosen for the entrypoint to listen on,
	// defined as ":0". For all options not specified, default values will be used.
	// TODO: do we really use :0 as address, or the v4 code to identify an interface?
	Config EntrypointConfig
}

// NewRegistrationFunc takes a registration function and handler and returns
// a registration func that can be used with one specific server type.
//
// This function is useful if a user wants to register a non-micro project proto
// handler. For any internal proto services generated in your project you will
// already have a pre-defined registration function which only needs the handler
// implementation.
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
