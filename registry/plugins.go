package registry

import (
	"github.com/go-orb/config/source"
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/types"
	"go-micro.dev/v5/util/container"
)

// Plugins is the plugins container for registry.
var Plugins = container.NewMap[func( //nolint:gochecknoglobals
	serviceName types.ServiceName,
	datas []source.Data,
	logger *log.Logger,
) (*OrbRegistry, error)]()
