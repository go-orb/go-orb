//go:build wireinject
// +build wireinject

package quickorb

import (
	"github.com/google/wire"
	"jochum.dev/orb/orb/cli"
)

func newService(
	options *Options,
	cliConfig cli.Config,
	registryConfig registry.Config,
) (*Service, error) {
	panic(wire.Build(
		ProvideService,
	))
}
