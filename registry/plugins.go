package registry

import (
	"go-micro.dev/v5/config/source"
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/types"
	"go-micro.dev/v5/util/container"
)

// Option is a functional option type for the registry.
type Option func(*Config)

// Provider is provider function type used by plugins to create a new registry.
type Provider func(name types.ServiceName, data []source.Data, logger log.Logger, opts ...Option)

// Plugins is the plugins container for registry.
//
//nolint:gochecknoglobals
var Plugins = container.NewMap[Provider]()
