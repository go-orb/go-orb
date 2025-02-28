// Package server provides the go-orb server. It is responsible for managing
// entrypoints.
package server

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

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

// Server is responsible for managing entrypoints. Entrypoints are the actual
// servers that bind to a port and accept connections. Entrypoints can be dynamically configured.
//
// For more info look at the entrypoint types.
type Server struct {
	// entrypoints are all created entrypoints.
	// All entrypoints will be started upon call of the Start method.
	entrypoints *container.Map[string, Entrypoint]
}

// Provide creates a new server.
//
//nolint:funlen,gocyclo
func Provide(
	name types.ServiceName,
	configs types.ConfigData,
	components *types.Components,
	logger log.Logger,
	reg registry.Type,
	opts ...ConfigOption,
) (Server, error) {
	cfg := Config{}
	for _, o := range opts {
		o(&cfg)
	}

	sections := append(types.SplitServiceName(name), DefaultConfigSection)
	if err := config.Parse(sections, configs, &cfg); err != nil {
		return Server{}, err
	}

	// Configure Middlewares.
	mws := []Middleware{}

	for idx, cfgMw := range cfg.Middlewares {
		pFunc, ok := Middlewares.Get(cfgMw.Plugin)
		if !ok {
			return Server{}, fmt.Errorf("%w: '%s', did you register it?", ErrUnknownMiddleware, cfgMw.Plugin)
		}

		mw, err := pFunc(append(sections, "middlewares", strconv.Itoa(idx)), configs, logger)
		if err != nil {
			return Server{}, err
		}

		mws = append(mws, mw)
	}

	// Get handlers.
	handlers := []RegistrationFunc{}

	for _, k := range cfg.Handlers {
		h, ok := Handlers.Get(k)
		if !ok {
			return Server{}, fmt.Errorf("%w: '%s', did you register it?", ErrUnknownHandler, k)
		}

		handlers = append(handlers, h)
	}

	// Configure entrypoints.
	eps := container.NewMap[string, Entrypoint]()

	for _, cfgNewEp := range cfg.functionalEntrypoints {
		newFunc, ok := PluginsNew.Get(cfgNewEp.config().Plugin)
		if !ok {
			return Server{}, fmt.Errorf("%w: '%s', did you register it?", ErrUnknownEntrypoint, cfgNewEp.config().Plugin)
		}

		ep, err := newFunc(cfgNewEp, logger, reg)
		if err != nil {
			return Server{}, err
		}

		if !ep.Enabled() {
			continue
		}

		eps.Set(ep.Name(), ep)
	}

	for idx, cfgEp := range cfg.Entrypoints {
		pFunc, ok := Plugins.Get(cfgEp.Plugin)
		if !ok {
			return Server{}, fmt.Errorf("%w: '%s', did you register it?", ErrUnknownEntrypoint, cfgEp.Plugin)
		}

		mSections := sections
		mSections = append(mSections, "entrypoints", strconv.Itoa(idx))

		ep, err := pFunc(mSections, configs, logger, reg, WithMiddlewares(mws...), WithHandlers(handlers...))
		if err != nil {
			return Server{}, err
		}

		if !ep.Enabled() {
			continue
		}

		eps.Set(ep.Name(), ep)
	}

	srv := Server{
		entrypoints: eps,
	}

	// Register the server as a component.
	err := components.Add(&srv, types.PriorityServer)
	if err != nil {
		logger.Warn("while registering server as a component", "error", err)
	}

	return srv, nil
}

// Start will start the HTTP servers on all entrypoints.
func (s *Server) Start() error {
	if s == nil {
		return errors.New("failed to create server can't start")
	}

	var gErr error

	s.entrypoints.Range(func(addr string, entrypoint Entrypoint) bool {
		if err := entrypoint.Start(); err != nil {
			// Stop any started entrypoints before returning error to give them a chance
			// to free up resources.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			_ = s.Stop(ctx) //nolint:errcheck

			gErr = multierror.Append(err, fmt.Errorf("start entrypoint (%s): %w", addr, err))

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

// GetEntrypoints returns a map of entrypoints.
func (s *Server) GetEntrypoints() *container.Map[string, Entrypoint] {
	return s.entrypoints
}

// GetEntrypoint returns the requested entrypoint, if present.
func (s *Server) GetEntrypoint(name string) (Entrypoint, error) {
	e, ok := s.entrypoints.Get(name)
	if !ok {
		return nil, errors.New("requested entrypoint was not found")
	}

	return e, nil
}

// Type returns the orb component type.
func (s *Server) Type() string {
	return ComponentType
}

// String is no-op.
func (s *Server) String() string {
	return ""
}
