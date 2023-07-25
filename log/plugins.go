package log

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// Provider can be started and stopped.
type Provider interface {
	fmt.Stringer

	Start() error
	Stop(context.Context) error

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
var pluginsCache = container.NewSafeMap[ProviderType]() //nolint:gochecknoglobals

// Register makes a plugin available by the provided name.
// If Register is called twice with the same name, it panics.
func Register(name string, pFunc ProviderFunc) bool {
	plugins.Register(name, pFunc)
	return true
}
