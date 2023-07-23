package registry

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// ProviderFunc is provider function type used by plugins to create a new registry.
type ProviderFunc func(
	name types.ServiceName,
	data types.ConfigData,
	logger log.Logger,
	opts ...Option,
) (Wire, error)

// Plugins is the plugins container for registry.
//
//nolint:gochecknoglobals
var Plugins = container.NewPlugins[ProviderFunc]()
