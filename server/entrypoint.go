package server

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
)

// EntrypointType will be returned by each entrypoint on Type().
const EntrypointType = "server.Entrypoint"

// RegistrationFunc is executed to register a handler to a server (entrypoint)
// passed as srv. srv can be of any of the various server types, should be a pointer.
//
// You can write your own custom registration functions to register extra handlers.
//
// Inside the registration function, you need to convert the server type and
// assert that you are working with the server type you are expecting, and
// otherwise no-op. For an example, see the implementation of NewRegistrationFunc.
type RegistrationFunc func(srv any)

// Entrypoint is a server, and represents an entrypoint into the web.
type Entrypoint interface {
	types.Component

	// Name returns the entrypoints name.
	Name() string

	// Enabled returns if this entrypoint has been enabled in config.
	Enabled() bool

	// AddHandler adds a handler for registration during startup.
	AddHandler(fun RegistrationFunc)

	// Register is used to register handlers.
	//
	// A registration function takes a pointer to the server, which can then
	// be used to register handlers in the server specific way.
	Register(fun RegistrationFunc)

	// Transport returns the client transport that is required to talk to this entrypoint.
	Transport() string

	// Address returns the address the entrypoint is listening on.
	Address() string

	// EntrypointID returns the id (uuid) of this entrypoint in the registry.
	EntrypointID() string
}

// EntrypointProvider is the function type to create a new entrypoint.
// It should create a new config, configure it the run EntrypointFromConfig with it.
type EntrypointProvider func(
	sections []string,
	configs types.ConfigData,
	logger log.Logger,
	reg registry.Type,
	opts ...Option,
) (Entrypoint, error)

// EntrypointNew is the function type to create a new entrypoint.
type EntrypointNew func(
	acfg any,
	logger log.Logger,
	reg registry.Type,
) (Entrypoint, error)

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
