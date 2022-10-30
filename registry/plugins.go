package registry

import (
	"github.com/go-orb/config/source"
	"github.com/go-orb/orb/log"
	"github.com/go-orb/orb/types"
	"github.com/go-orb/orb/util/container"
)

// Plugins is the plugins container for registry.
var Plugins = container.NewMap[func( //nolint:gochecknoglobals
	serviceName types.ServiceName,
	datas []source.Data,
	logger *log.Logger,
) (*OrbRegistry, error)]()
