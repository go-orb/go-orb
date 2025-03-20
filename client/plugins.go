package client

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// ProviderFunc is provider function type used by plugins to create a new client.
type ProviderFunc func(
	configData map[string]any,
	components *types.Components,
	logger log.Logger,
	registry registry.Type,
	opts ...Option,
) (Type, error)

// plugins is the container for client implementations.
//
//nolint:gochecknoglobals
var plugins = container.NewMap[string, ProviderFunc]()

// Register makes a plugin available by the provided name.
// If Register is called twice with the same name, it panics.
func Register(name string, factory ProviderFunc) bool {
	plugins.Add(name, factory)
	return true
}
