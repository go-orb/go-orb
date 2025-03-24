// Package registry is a component for service discovery
package registry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log/slog"

	"github.com/go-orb/go-orb/cli"
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

// ServiceNode is a service node.
type ServiceNode struct {
	// Name is the name of the service. Should be DNS compatible.
	Name string `json:"name"`
	// Version is the version of the service.
	Version string `json:"version"`

	// Metadata is the metadata of the service.
	Metadata map[string]string `json:"metadata"`

	// Scheme is the scheme of the service.
	Scheme string `json:"scheme"`
	// Address is the address of the service.
	Address string `json:"address"`

	// Namespace is the namespace of the node.
	Namespace string `json:"namespace"`

	// Region is the region of the node.
	Region string `json:"region"`

	// TTL is the time to live for the service.
	// Keep it 0 if you don't want to use TTL.
	TTL time.Duration `json:"ttl"`
}

func (r ServiceNode) String() string {
	return fmt.Sprintf("ServiceNode{%s %s %s %s %s %s}", r.Namespace, r.Region, r.Name, r.Version, r.Address, r.Scheme)
}

// Registry is a component for service discovery.
type Registry interface {
	types.Component

	// Register registers a service within the registry.
	Register(ctx context.Context, srv ServiceNode) error

	// Deregister deregisters a service within the registry.
	Deregister(ctx context.Context, srv ServiceNode) error

	// GetService returns a service from the registry.
	// Leave schemes empty to get all schemes.
	GetService(ctx context.Context, namespace, region, name string, schemes []string) ([]ServiceNode, error)

	// ListServices lists services within the registry.
	// Leave schemes empty to get all schemes.
	ListServices(ctx context.Context, namespace, region string, schemes []string) ([]ServiceNode, error)

	// Watch returns a Watcher which you can watch on.
	Watch(ctx context.Context, opts ...WatchOption) (Watcher, error)
}

// Type is the registry type it is returned when you use ProvideRegistry
// which selects a registry to use based on the plugin configuration.
type Type struct {
	Registry
}

// New creates a new registry without side-effects.
func New(
	configData map[string]any,
	components *types.Components,
	logger log.Logger,
	opts ...Option,
) (Type, error) {
	cfg := NewConfig(opts...)

	if err := config.Parse(nil, DefaultConfigSection, configData, &cfg); err != nil && !errors.Is(err, config.ErrNoSuchKey) {
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
	cLogger, err := logger.WithConfig([]string{DefaultConfigSection}, configData)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	instance, err := provider(configData, components, cLogger, opts...)
	if err != nil {
		return Type{}, err
	}

	return Type{Registry: instance}, nil
}

// Provide is the registry provider for wire.
// It parses the config from "configs", fetches the "Plugin" from the config and
// then forwards all it's arguments to the factory which it get's from "Plugins".
func Provide(
	svcCtx *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
	opts ...Option,
) (Type, error) {
	reg, err := New(svcCtx.Config, components, logger, opts...)
	if err != nil {
		return Type{}, err
	}

	// Register the registry as a component.
	err = components.Add(&reg, types.PriorityRegistry)
	if err != nil {
		logger.Warn("while registering registry as a component", "error", err)
	}

	return reg, nil
}

// ProvideNoOpts is the registry provider for wire without options.
func ProvideNoOpts(
	svcCtx *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
) (Type, error) {
	return Provide(svcCtx, components, logger)
}
