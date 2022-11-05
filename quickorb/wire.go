//go:build wireinject
// +build wireinject

package quickorb

import (
	"go-micro.dev/v5/cli"
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/registry"
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
