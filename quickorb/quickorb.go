// Package quickorb is the quick start entry for the orb framework.
package quickorb

import (
	"errors"

	"jochum.dev/orb/orb/cli"
)

func NewService(opts ...Option) (*Service, error) {
	options := NewOptions(opts...)

	// Setup cli
	cliConfig := cli.NewComponentConfig()
	bTrue := true
	cliConfig.Enabled = &bTrue
	cliConfig.Plugin = "urfave"
	cliConfig.Name = options.Name
	cliConfig.Version = options.Version
	cliConfig.Description = options.Description
	cliConfig.Usage = options.Usage
	cliConfig.ConfigSection = options.ConfigSection
	cliConfig.ArgPrefix = options.ArgPrefix
	cliConfig.NoFlags = &options.NoFlags
	cliConfig.Config = options.ConfigURLs
	cliConfig.Flags = options.Flags

	return nil, errors.New("unimplemented")
}
