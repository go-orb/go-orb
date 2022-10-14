// Package quickorb is the quick start entry for the orb framework.
package quickorb

import (
	"errors"

	"jochum.dev/orb/orb/cli"
	"jochum.dev/orb/orb/server"
)

func NewService(opts ...Option) (*Service, error) {
	options := NewOptions(opts...)

	// Setup cli
	cliConfig := cli.NewConfig()
	cliConfig.SetPlugin("urfave")
	cliConfig.SetName(options.Name)
	cliConfig.SetVersion(options.Version)
	cliConfig.SetDescription(options.Description)
	cliConfig.SetUsage(options.Usage)
	cliConfig.SetConfigSection(options.ConfigSection)
	cliConfig.SetArgPrefix(options.ArgPrefix)
	cliConfig.SetNoFlags(&options.NoFlags)
	cliConfig.SetConfig(options.ConfigURLs)
	cliConfig.SetFlags(options.Flags)

	serverConfig := server.NewConfig()
	serverConfig.SetAddress(options.Address)
	serverConfig.SetRegisterTTL(options.RegisterTTL)
	serverConfig.SetRegisterInterval(options.RegisterInterval)
	serverConfig.SetMetadata(options.Metadata)

	return nil, errors.New("unimplemented")
}
