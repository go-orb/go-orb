//go:build wireinject
// +build wireinject

package quickorb

import (
	"github.com/google/wire"
	"github.com/orb-org/orb/cli"
	"jochum.dev/orb/orb/log"
	"jochum.dev/orb/orb/registry"
)

func newService(
	options *Options,
	cliConfig cli.Config,
	logConfig log.Config,
	registryConfig registry.Config,
) (*Service, error) {
	panic(wire.Build(
		ProvideService,
	))
}
