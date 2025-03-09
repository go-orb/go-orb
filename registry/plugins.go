package registry

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// ProviderFunc is provider function type used by plugins to create a new registry.
type ProviderFunc func(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	data types.ConfigData,
	components *types.Components,
	logger log.Logger,
	opts ...Option,
) (Type, error)

// Plugins is the plugins container for registries.
//
//nolint:gochecknoglobals
var Plugins = container.NewMap[string, ProviderFunc]()
