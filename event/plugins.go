package event

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// ProviderFunc is provider function type used by plugins to create a new client.
type ProviderFunc func(
	name types.ServiceName,
	data types.ConfigData,
	logger log.Logger,
	opts ...Option,
) (Type, error)

// plugins is the container for client implementations.
//
//nolint:gochecknoglobals
var plugins = container.NewPlugins[ProviderFunc]()

// Register makes a plugin available by the provided name.
// If Register is called twice with the same name, it panics.
func Register(name string, factory ProviderFunc) bool {
	plugins.Register(name, factory)
	return true
}
