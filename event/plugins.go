package event

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/util/container"
)

// ProviderFunc is provider function type used by plugins to create a new client.
type ProviderFunc func(
	configData map[string]any,
	logger log.Logger,
	opts ...Option,
) (Type, error)

// plugins is the container for client implementations.
//
//nolint:gochecknoglobals
var plugins = container.NewMap[string, ProviderFunc]()

// Register makes a plugin available by the provided name.
func Register(name string, factory ProviderFunc) {
	plugins.Add(name, factory)
}
