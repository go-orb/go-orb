package registry

import (
	"fmt"

	"github.com/google/wire"
	"jochum.dev/orb/orb/cli"
	"jochum.dev/orb/orb/config"
	"jochum.dev/orb/orb/di"
	"jochum.dev/orb/orb/log"
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
	if *cliConfig.GetNoFlags() {
		// Defined silently ignore that
		return DiFlags{}, nil
	}

	if err := c.Add(
		cli.Name(cli.PrefixName(cliConfig.GetArgPrefix(), cliArgPlugin)),
		cli.Usage("Registry for discovery. etcd, mdns"),
		cli.Default(config.GetPlugin()),
		cli.EnvVars(cli.PrefixEnv(cliConfig.GetArgPrefix(), cliArgPlugin)),
	); err != nil {
		return DiFlags{}, err
	}

	if err := c.Add(
		cli.Name(cli.PrefixName(cliConfig.GetArgPrefix(), cliArgAddresses)),
		cli.Usage("List of registry addresses"),
		cli.Default(config.GetAddresses()),
		cli.EnvVars(cli.PrefixEnv(cliConfig.GetArgPrefix(), cliArgAddresses)),
	); err != nil {
		return DiFlags{}, err
	}

	return DiFlags{}, nil
}

func ProvideConfig(
	_ di.DiConfig,
	flags DiFlags,
	cfg Config,
	c cli.Cli,
	cliConfig cli.Config,
	configDatas []config.Data,
) (DiConfig, error) {

	if cliConfig.GetNoFlags() != nil && *cliConfig.GetNoFlags() {
		// Dont parse flags if NoFlags has been given.
		return DiConfig{}, nil
	}

	// Read flags into config.
	newCfg := NewComponentConfig()
	if f, ok := c.Get(cli.PrefixName(cliConfig.GetArgPrefix(), cliArgPlugin)); ok {
		newCfg.Plugin = cli.FlagValue(f, newCfg.GetPlugin())
	}
	if f, ok := c.Get(cli.PrefixName(cliConfig.GetArgPrefix(), cliArgAddresses)); ok {
		newCfg.Addresses = cli.FlagValue(f, newCfg.GetAddresses())
	}
	if err := cfg.Merge(newCfg); err != nil {
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
	if !*config.GetEnabled() {
		// Not enabled silently ignore that
		return nil, nil
	}

	pluginFunc, err := Plugins.Plugin(config.GetPlugin())
	if err != nil {
		return nil, fmt.Errorf("unknown plugin registry: '%s'", config.GetPlugin())
	}
	plugin := pluginFunc()

	opts := []Option{}

	logger, err := log.FromConfig(config.GetLogger(), parentLogger)
	if err != nil {
		return nil, fmt.Errorf("component '%s' is unable to setup its logger", config.GetPlugin())
	}

	opts = append(opts, WithLogger(logger))

	if err := plugin.Init(config, opts...); err != nil {
		return nil, err
	}

	return plugin, nil
}

var DiSet = wire.NewSet(ProvideFlags, ProvideConfig, Provide)
