package log

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// PluginProvider can be started and stopped.
type PluginProvider interface {
	fmt.Stringer

	Start() error
	Stop(context.Context) error

	Handler() (slog.Handler, error)
}

// PluginProviderType is the struct that wraps the interface, so we don't have to pass interfaces everywhere.
type PluginProviderType struct {
	PluginProvider
}

// PluginProviderFunc is the function a Plugin must provide which returns a PluginProvider.
type PluginProviderFunc func(section []string, data types.ConfigData) (PluginProviderType, error)

// Plugins is the registry for Logger plugins.
var Plugins = container.NewPlugins[PluginProviderFunc]() //nolint:gochecknoglobals

// PluginsCache contains plugin's already loaded and started.
var PluginsCache = container.NewSafeMap[PluginProviderType]() //nolint:gochecknoglobals
