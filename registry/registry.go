// Package registry is a component for service discovery
package registry

import (
	"errors"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
)

// ComponentType is the registry component type name.
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

	// Register registers a service within the registry.
	Register(*Service, ...RegisterOption) error

	// Deregister deregisters a service within the registry.
	Deregister(*Service, ...DeregisterOption) error

	// GetService returns a service from the registry.
	GetService(string, ...GetOption) ([]*Service, error)

	// ListServices lists services within the registry.
	ListServices(...ListOption) ([]*Service, error)

	// Watch returns a Watcher which you can watch on.
	Watch(...WatchOption) (Watcher, error)
}

// Instance is the registry type it is returned when you use ProvideRegistry
// which selects a registry to use based on the plugin configuration.
type Instance struct {
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
	// frpc/grpc/http uvm., since v5!
	Scheme   string            `json:"scheme"`
	Metadata map[string]string `json:"metadata"`
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
// It parses the config from "data", fetches the "Plugin" from the config and
// then forwards all it's arguments to the provider which it get's from "Plugins".
func ProvideRegistry(
	name types.ServiceName,
	data types.ConfigData,
	logger log.Logger,
	opts ...Option) (Instance, error) {
	cfg := NewConfig(opts...)

	sections := types.SplitServiceName(name)
	if err := config.Parse(append(sections, DefaultConfigSection), data, cfg); err != nil {
		return Instance{}, err
	}

	if cfg.Plugin == "" {
		logger.Warn("empty registry plugin, using the default", "default", DefaultRegistry)
		cfg.Plugin = DefaultRegistry
	}

	provider, err := Plugins.Get(cfg.Plugin)
	if err != nil {
		return Instance{}, err
	}

	return provider(name, data, logger, opts...)
}
