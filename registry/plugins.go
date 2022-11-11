package registry

import (
	"go-micro.dev/v5/config/source"
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/types"
	"go-micro.dev/v5/util/container"
)

// ProviderFunc is provider function type used by plugins to create a new registry.
type ProviderFunc func(
	name types.ServiceName,
	data []source.Data,
	logger log.Logger,
	opts ...Option,
) (*MicroRegistry, error)

// Plugins is the plugins container for registry.
//
//nolint:gochecknoglobals
var Plugins = container.NewMap[ProviderFunc]()
