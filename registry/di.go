// Code generated with jinja2 templates. DO NOT EDIT.

package registry

import (
	"fmt"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"jochum.dev/orb/orb/cli"
	"jochum.dev/orb/orb/config"
	"jochum.dev/orb/orb/config/chelp"
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
	configDatas []config.Data,
) (DiConfig, error) {

	for _, configData := range configDatas {
		// Go own section deeper to the config section.
		var err error
		data := configData.Data
		if cliConfig.ConfigSection() != "" {
			if data, err = chelp.Get(data, cliConfig.ConfigSection(), map[string]any{}); err != nil {
				// Ignore unknown configSection in config.
				if errors.Is(err, chelp.ErrNotExistant) {
					log.Warn().
						Fields(map[string]string{"section": cliConfig.ConfigSection(), "url": configData.URL.String()}).
						Msg("unknown config section in config")
					continue
				}
				return DiConfig{}, err
			}
		}

		// Now fetch my own section.
		if data, err = chelp.Get(data, Name, map[string]any{}); err != nil {
			// Ignore unknown section in config.
			if errors.Is(err, chelp.ErrNotExistant) {
				log.Warn().
					Fields(map[string]string{"section": Name, "url": configData.URL.String()}).
					Msg("unknown config section in config")
				continue
			}
			return DiConfig{}, err
		}

		// Create a new config.
		aNew, err := NewConfig(config.Plugin())
		if err != nil {
			return DiConfig{}, err
		}

		newCfg, ok := aNew.(chelp.ConfigMethods)
		if !ok {
			return DiConfig{}, chelp.ErrUnknownConfig
		}

		// Load the new config.
		if err := newCfg.Load(data); err != nil {
			return DiConfig{}, err
		}

		// Merge it into the previous one.
		if err := config.Merge(newCfg); err != nil {
			return DiConfig{}, err
		}
	}

	if *cliConfig.NoFlags() {
		// Dont parse flags if NoFlags has been given.
		return DiConfig{}, nil
	}

	// Read flags into config.
	newCfg := NewBaseConfig()
	if f, ok := c.Get(cli.PrefixName(cliConfig.ArgPrefix(), cliArgPlugin)); ok {
		newCfg.SetPlugin(cli.FlagValue(f, newCfg.Plugin()))
	}
	if f, ok := c.Get(cli.PrefixName(cliConfig.ArgPrefix(), cliArgAddresses)); ok {
		newCfg.SetAddresses(cli.FlagValue(f, newCfg.Addresses()))
	}
	if err := config.Merge(newCfg); err != nil {
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