// Package registry is a component for service discovery
package registry

import (
	"errors"
	"fmt"

	"log/slog"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
)

// ComponentType is the components name.
const ComponentType = "registry"

var (
	// ErrNotFound is a not found error when GetService is called.
	ErrNotFound = errors.New("service not found")
	// ErrWatcherStopped is a error when watcher is stopped.
	ErrWatcherStopped = errors.New("watcher stopped")
)

// Registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}.
type Registry interface {
	types.Component

	ServiceName() string
	ServiceVersion() string

	// Register registers a service within the registry.
	Register(srv *Service, opts ...RegisterOption) error

	// Deregister deregisters a service within the registry.
	Deregister(srv *Service, opts ...DeregisterOption) error

	// GetService returns a service from the registry.
	GetService(name string, opts ...GetOption) ([]*Service, error)

	// ListServices lists services within the registry.
	ListServices(opts ...ListOption) ([]*Service, error)

	// Watch returns a Watcher which you can watch on.
	Watch(opts ...WatchOption) (Watcher, error)
}

// Type is the registry type it is returned when you use ProvideRegistry
// which selects a registry to use based on the plugin configuration.
type Type struct {
	Registry
}

// Service represents a service in a registry.
type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

// Node represents a service node in a registry.
// One service can be comprised of multiple nodes.
type Node struct {
	ID string `json:"id"`
	// ip:port
	Address string `json:"address"`
	// grpc/h2c/http/http3 uvm., since go-orb!
	Transport string            `json:"transport"`
	Metadata  map[string]string `json:"metadata"`
}

// Endpoint represents a service endpoint in a registry.
type Endpoint struct {
	Name     string            `json:"name"`
	Request  *Value            `json:"request"`
	Response *Value            `json:"response"`
	Metadata map[string]string `json:"metadata"`
}

// Value is a value container used in the registry.
type Value struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []*Value `json:"values"`
}

// ProvideRegistry is the registry provider for wire.
// It parses the config from "configs", fetches the "Plugin" from the config and
// then forwards all it's arguments to the factory which it get's from "Plugins".
func ProvideRegistry(
	name types.ServiceName,
	version types.ServiceVersion,
	configs types.ConfigData,
	logger log.Logger,
	opts ...Option) (Type, error) {
	cfg := NewConfig(opts...)

	sections := append(types.SplitServiceName(name), DefaultConfigSection)
	if err := config.Parse(sections, configs, &cfg); err != nil {
		return Type{}, err
	}

	if cfg.Plugin == "" {
		logger.Warn("empty registry plugin, using the default", "default", DefaultRegistry)
		cfg.Plugin = DefaultRegistry
	}

	logger.Debug("Registry", "plugin", cfg.Plugin)

	provider, ok := Plugins.Get(cfg.Plugin)
	if !ok {
		return Type{}, fmt.Errorf("Registry plugin '%s' not found, did you import it?", cfg.Plugin)
	}

	// Configure the logger.
	cLogger, err := logger.WithConfig(sections, configs)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	return provider(name, version, configs, cLogger, opts...)
}
