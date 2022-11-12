package server

import (
	"context"
	"fmt"

	"go-micro.dev/v5/config"
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/types"
	"go-micro.dev/v5/types/component"
)

var _ component.Component = (*MicroServer)(nil)

const ComponentType component.Type = "server"

// MicroServer is repsonsible for managing entrypoints. Entrypoints are the actual
// servers that bind to a port and accept connections. Entrypoints can be dynamically configured.
//
// For more info look at the entrypoint types.
type MicroServer struct {
	service    types.ServiceName
	configData types.ConfigData

	Logger log.Logger

	Config Config

	// entrypoints are all created entrypoints. All of the entrypoints in this
	// map will be started upon the call of Start method.
	entrypoints map[string]Entrypoint

	// fileConfig is used to parse the fileConfig into, it is used to extract
	// the list of entrypoints defined in the fileConfig.
	fileConfig map[string]fileConfigServer
}

// ProvideServer creates a new server.
func ProviderServer(name types.ServiceName, data types.ConfigData, logger log.Logger, opts ...Option) (*MicroServer, error) {
	cfg, err := NewConfig(name, data, opts...)
	if err != nil {
		return nil, fmt.Errorf("create http server config: %w", err)
	}

	s := MicroServer{
		service:     name,
		configData:  data,
		Config:      cfg,
		Logger:      logger,
		entrypoints: make(map[string]Entrypoint),
		fileConfig:  make(map[string]fileConfigServer),
	}

	// Set enabled to true by default for all registered plugins.
	for plugin := range Plugins.All() {
		s.fileConfig[plugin] = fileConfigServer{Enabled: true}
	}

	sections := types.SplitServiceName(name)
	if err := config.Parse(append(sections, DefaultConfigSection), data, &s.fileConfig); err != nil {
		return &s, err
	}

	if err := s.createEntrypoints(); err != nil {
		return nil, err
	}

	return &s, nil
}

// Start will start the HTTP servers on all entrypoints.
func (s *MicroServer) Start() error {
	// TODO: somehow incorporate the yaml entrypoints here; if yaml entry not in map, take "use" key into account, error when not found
	for addr, entrypoint := range s.entrypoints {
		if err := entrypoint.Start(); err != nil {
			return fmt.Errorf("start entrypoint (%s): %w", addr, err)
		}
	}

	return nil
}

// Stop will stop the HTTP servers on all entrypoints and close the listeners.
func (s *MicroServer) Stop(ctx context.Context) error {
	errChan := make(chan error)

	// Stop all servers in parallel to make sure they get equal amount of time
	// to shutdown gracefully.
	for _, e := range s.entrypoints {
		go func(e Entrypoint) {
			errChan <- e.Stop(ctx)
		}(e)
	}

	var err error

	for i := 0; i < len(s.entrypoints); i++ {
		if nerr := <-errChan; nerr != nil {
			err = fmt.Errorf("stop entrypoint: %w", nerr)
		}
	}

	close(errChan)

	return err
}

// Type returns the micro component type.
func (s *MicroServer) Type() component.Type {
	return ComponentType
}

// String is no-op.
func (s *MicroServer) String() string {
	return ""
}

func (s *MicroServer) createEntrypoints() error {
	for name, template := range s.Config.Templates {
		// If a plugin or specific entrypoint has been globally disabled in config, skip.
		c := s.fileConfig[template.Type]
		if !c.Enabled || !c.IsEnabled(name) {
			continue
		}

		newEntrypoint, err := Plugins.Get(template.Type)
		if err != nil {
			return fmt.Errorf("entrypoint provider for %s not found, did you register it by importing the package?", template.Type)
		}

		cfg := s.Config.Defaults[template.Type]

		inherit := c.Inherit(name)
		if len(inherit) > 0 {
			var ok bool

			cfg, ok = s.Config.Templates[inherit]
			if !ok {
				return fmt.Errorf("%s failed to inherit config from %s, entrypoint not found", name, inherit)
			}
		}

		entrypoint, err := newEntrypoint(name, s.service, s.configData, s.Logger, cfg, template.Options...)
		if err != nil {
			return fmt.Errorf("create entrypoint %s (%s): %w", name, template.Type, err)
		}

		s.entrypoints[name] = entrypoint
	}

	return nil
}
