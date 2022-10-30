//go:build wireinject
// +build wireinject

package quickorb

import (
	"github.com/go-orb/orb/cli"
	"github.com/google/wire"
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
