// Code generated with jinja2 templates. DO NOT EDIT.

package registry

import (
	"fmt"
	"jochum.dev/orb/orb/log"
	"jochum.dev/orb/orb/cli"
	"jochum.dev/orb/orb/config"
	"jochum.dev/orb/orb/di"
	"github.com/google/wire"
)

type DiFlags struct{}

// DiConfig is marker that DiFlags has been parsed into Config
type DiConfig struct{}

const (
	cliArgPlugin    = "registry"
	cliArgAddresses = "registry_address"
)

func ProvideFlags(
	config Config,
	cliConfig cli.Config,
	c cli.Cli,
) (DiFlags, error) {
	if *cliConfig.NoFlags() {
		// Defined silently ignore that
		return DiFlags{}, nil
	}

	if err := c.Add(
		cli.Name(cli.PrefixName(cliConfig.ArgPrefix(), cliArgPlugin)),
		cli.Usage("Registry for discovery. etcd, mdns"),
		cli.Default(config.Plugin()),
		cli.EnvVars(cli.PrefixEnv(cliConfig.ArgPrefix(), cliArgPlugin)),
	); err != nil {
		return DiFlags{}, err
	}

	if err := c.Add(
		cli.Name(cli.PrefixName(cliConfig.ArgPrefix(), cliArgAddresses)),
		cli.Usage("List of registry addresses"),
		cli.Default(config.Addresses()),
		cli.EnvVars(cli.PrefixEnv(cliConfig.ArgPrefix(), cliArgAddresses)),
	); err != nil {
		return DiFlags{}, err
	}

	return DiFlags{}, nil
}

func ProvideConfig(
	_ di.DiConfig,
	flags DiFlags,
	config Config,
	c cli.Cli,
	cliConfig cli.Config,
	confiData []config.Data,
) (DiConfig, error) {
	defConfig := NewConfig()
	cfg := sourceConfig{Registry: *defConfig}

	if configor != nil {
		if err := configor.Scan(&cfg); err != nil {
			return DiConfig{}, err
		}
	}
	if err := config.Merge(&cfg.Registry); err != nil {
		return DiConfig{}, err
	}

	if *cliConfig.NoFlags() {
		// Dont parse flags if NoFlags has been given
		return DiConfig{}, nil
	}

	defConfig = NewConfig()
	if f, ok := c.Get(cli.PrefixName(cliConfig.ArgPrefix(), cliArgPlugin)); ok {
		defConfig.SetPlugin(cli.FlagValue(f, defConfig.Plugin()))
	}
	if f, ok := c.Get(cli.PrefixName(cliConfig.ArgPrefix(), cliArgAddresses)); ok {
		defConfig.SetAddresses(cli.FlagValue(f, defConfig.Addresses()))
	}
	if err := config.Merge(defConfig); err != nil {
		return DiConfig{}, err
	}

	return DiConfig{}, nil
}

func ProvideConfigNoFlags(
	config Config,
	configor config.Config,
) (DiConfig, error) {
	defConfig := NewConfig()
	c := sourceConfig{Registry: *defConfig}

	if configor != nil {
		if err := configor.Scan(&c); err != nil {
			return DiConfig{}, err
		}
	}
	if err := config.Merge(&c.Registry); err != nil {
		return DiConfig{}, err
	}

	return DiConfig{}, nil
}

func Provide(
	// Marker so cli has been merged into Config
	_ DiConfig,
	parentLogger log.Logger,
	config Config,
) (Registry, error) {
	if !*config.Enabled() {
		// Not enabled silently ignore that
		return nil, nil
	}

	pluginFunc, err := Plugins.Plugin(config.Plugin())
	if err != nil {
		return nil, fmt.Errorf("unknown plugin registry: '%s'", config.Plugin())
	}
	plugin := pluginFunc()

	opts := []Option{}

	logger, err := log.FromConfig(config.Logger(), parentLogger)
	if err != nil {
		return nil, fmt.Errorf("component '%s' is unable to setup its logger", config.Plugin())
	}

	opts = append(opts, WithLogger(logger))

	if err := plugin.Init(config, opts...); err != nil {
		return nil, err
	}
	
	return plugin, nil
}

var DiSet = wire.NewSet(ProvideFlags, ProvideConfig, Provide)
var DiNoCliSet = wire.NewSet(ProvideConfigNoFlags, Provide)
