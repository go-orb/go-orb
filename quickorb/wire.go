//go:build wireinject
// +build wireinject

package quickorb

import (
	"github.com/go-orb/orb/cli"
	"github.com/go-orb/orb/log"
	"github.com/go-orb/orb/registry"
	"github.com/google/wire"
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
