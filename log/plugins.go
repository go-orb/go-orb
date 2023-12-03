package log

import (
	"context"

	"log/slog"

	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// Provider can be started/stopped and it returns a slog.Handler on request.
// Providers must be cacheable - there will be always one Provider for a given type AND config.
type Provider interface {
	// Key must return a unique key for the cache,
	// this should be unique for this Provider with its config.
	Key() string

	Start() error
	Stop(ctx context.Context) error

	Handler() (slog.Handler, error)
}

// ProviderType is the struct that wraps the interface, so we don't have to pass interfaces everywhere.
type ProviderType struct {
	Provider
}

// ProviderFunc is the function a Plugin must provide which returns a Provider encapsulated into a ProviderType.
type ProviderFunc func(section []string, configs types.ConfigData, opts ...Option) (ProviderType, error)

// Plugins is the registry for Logger plugins.
var plugins = container.NewPlugins[ProviderFunc]() //nolint:gochecknoglobals

// PluginsCache contains plugin's already loaded and started.
var pluginsCache = container.NewSafeMap[string, ProviderType]() //nolint:gochecknoglobals

// Register makes a plugin available by the provided name.
// If Register is called twice with the same name, it panics.
func Register(name string, pFunc ProviderFunc) bool {
	plugins.Register(name, pFunc)
	return true
}
