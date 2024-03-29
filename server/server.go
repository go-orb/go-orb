// Package server provides the go-micro server. It is responsible for managing
// entrypoints.
//
// # Entrypoints
//
// Entrypoints are the actual servers used that listen for incoming requests.
// Various entrypoint plugins are provided by default, but it is straight forward
// to create your own entrypoint implementation. Entrypoints are configured
// through functional options, and your config file. Entrypionts can be
// dynamically added, modified, or disabled through your config files.
//
// # Handler registrations
//
// Entrypoints can be used in any number of combinations. The handlers are
// registered by providing registration functions to the entrypoint config.
// A handler registration function takes care of registering the handler in
// the server specific way. While internal project handlers are designed such
// that they can be used with any type of server out of the box, the way
// they are registered usually differs per server type. Registration functions
// take care of this by switching on the server type. This also allows you to
// create server specific handlers if necessary.
//
// # Internal handlers
//
// The server has been architected with protobuf service definitions as primary
// handler types. Thus registration of these has been made as easy as possible.
// For proto services defined within your go-micro project, registration
// functions will be automatically generated for you, and you only need to
// provide the handler implementation, everything beyond is taken care of.
//
// # External handlers
//
// You may wish to register either external proto service handlers, or server
// specific handlers such as any existing HTTP handlers. THis is also possible.
// External proto services can be registered with the help of the
// NewRegistrationFunc type, which utliizes the power of generics to allow you
// to convert any gRPC registration into an entrypoint registration function.
// It is also possible to manually define your own registration functions. These
// must take one parameter of type any and convert it into the required server
// type, such as the go-micro HTTP server, or the go-micro gRPC server.
package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log/slog"

	"github.com/hashicorp/go-multierror"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

var _ types.Component = (*Server)(nil)

// ComponentType is the server component type name.
const ComponentType = "server"

// Errors.
var (
	ErrEntrypointNotFound = errors.New("requested entrypoint not found")
)

// Server is repsponsible for managing entrypoints. Entrypoints are the actual
// servers that bind to a port and accept connections. Entrypoints can be dynamically configured.
//
// For more info look at the entrypoint types.
type Server struct {
	Logger   log.Logger
	Config   Config
	Registry registry.Type

	// entrypoints are all created entrypoints. All of the entrypoints in this
	// map will be started upon the call of Start method.
	entrypoints *container.SafeMap[string, Entrypoint]
}

// ProvideServer creates a new server.
func ProvideServer(
	name types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	reg registry.Type,
	opts ...Option,
) (Server, error) {
	cfg := NewConfig(opts...)

	srv := Server{
		Config:      cfg,
		Logger:      logger,
		Registry:    reg,
		entrypoints: container.NewSafeMap[string, Entrypoint](),
	}

	sections := append(types.SplitServiceName(name), DefaultConfigSection)
	if err := config.Parse(sections, configs, &srv.Config); err != nil {
		return srv, err
	}

	if err := srv.createEntrypoints(name); err != nil {
		return srv, err
	}

	return srv, nil
}

// Start will start the HTTP servers on all entrypoints.
func (s *Server) Start() error {
	if s == nil {
		return errors.New("failed to create server can't start")
	}

	// TODO(davincible): catch startup errors better from blocking go-routines
	var gErr error

	s.entrypoints.Range(func(addr string, entrypoint Entrypoint) bool {
		if err := entrypoint.Start(); err != nil {
			// Stop any started entrypoints before returning error to give them a chance
			// to free up resources.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			_ = s.Stop(ctx) //nolint:errcheck

			gErr = fmt.Errorf("start entrypoint (%s): %w", addr, err)
			return false
		}

		return true
	})

	return gErr
}

// Stop will stop the servers on all entrypoints and close the listeners.
func (s *Server) Stop(ctx context.Context) error {
	if s == nil {
		return errors.New("failed to create server can't stop")
	}

	errChan := make(chan error, s.entrypoints.Len())

	// Stop all servers.
	s.entrypoints.Range(func(_ string, e Entrypoint) bool {
		errChan <- e.Stop(ctx)

		return true
	})

	var err error

	for i := 0; i < s.entrypoints.Len(); i++ {
		if nerr := <-errChan; nerr != nil {
			err = multierror.Append(err, fmt.Errorf("stop entrypoint: %w", nerr))
		}
	}

	close(errChan)

	return err
}

// GetEntrypoint returns the requested entrypoint, if present.
func (s *Server) GetEntrypoint(name string) (Entrypoint, error) {
	e, ok := s.entrypoints.Get(name)
	if !ok {
		return nil, ErrEntrypointNotFound
	}

	return e, nil
}

// Type returns the micro component type.
func (s *Server) Type() string {
	return ComponentType
}

// String is no-op.
func (s *Server) String() string {
	return ""
}

func (s *Server) createEntrypoints(service types.ServiceName) error {
	for name, template := range s.Config.Templates {
		// If a plugin or specific entrypoint has been globally disabled in config, skip.
		if enabled, ok := s.Config.Enabled[template.Type]; (ok && !enabled) || !template.Enabled {
			continue
		}

		provider, ok := Plugins.Get(template.Type)
		if !ok {
			return fmt.Errorf("entrypoint provider for %s not found, did you register it by importing the package?", template.Type)
		}

		if template.Config == nil {
			return fmt.Errorf("template config for %s is nil", name)
		}

		pluginLogger := s.Logger.With(slog.String("component", ComponentType), slog.String("plugin", template.Type))

		entrypoint, err := provider(service, pluginLogger, s.Registry, template.Config)
		if err != nil {
			return fmt.Errorf("create entrypoint %s (%s): %w", name, template.Type, err)
		}

		s.entrypoints.Set(name, entrypoint)
	}

	return nil
}
