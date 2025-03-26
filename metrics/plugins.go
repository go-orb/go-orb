package metrics

import (
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/util/container"
)

// ProviderFunc is provider function type is used to create a new metrics plugin.
type ProviderFunc func(
	svcCtx *cli.ServiceContextWithConfig,
	logger log.Logger,
	opts ...Option,
) (Type, error)

// Plugins is the plugins container for registries.
//
//nolint:gochecknoglobals
var Plugins = container.NewMap[string, ProviderFunc]()
