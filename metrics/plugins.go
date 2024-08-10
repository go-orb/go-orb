package metrics

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/container"
)

// ProviderFunc is provider function type is used to create a new metrics plugin.
type ProviderFunc func(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	data types.ConfigData,
	logger log.Logger,
	opts ...Option,
) (Type, error)

// Plugins is the plugins container for registries.
//
//nolint:gochecknoglobals
var Plugins = container.NewPlugins[ProviderFunc]()
