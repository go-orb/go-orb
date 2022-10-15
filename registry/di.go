package registry

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/google/wire"
	"github.com/pkg/errors"
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
	plugin, err := config.GetPluginFromConfigData(cliConfig.GetConfigSection(), Name, configDatas)
	if err != nil {
		return DiConfig{}, err
	}

	previousConfig := any(cfg)
	for _, configData := range configDatas {
		// Go own section deeper.
		var err error
		data := configData.Data
		if cliConfig.GetConfigSection() != "" {
			if data, err = config.Get(data, cliConfig.GetConfigSection(), map[string]any{}); err != nil {
				// Ignore unknown configSection in config.
				if errors.Is(err, config.ErrNotExistent) {
					log.Warn().
						Fields(map[string]string{"section": cliConfig.GetConfigSection(), "url": configData.URL.String()}).
						Msg("unknown config section in config")
					continue
				}
				return DiConfig{}, err
			}
		}

		// Now fetch my own section.
		if data, err = config.Get(data, Name, map[string]any{}); err != nil {
			// Ignore unknown section in config.
			if errors.Is(err, config.ErrNotExistent) {
				log.Warn().
					Fields(map[string]string{"section": Name, "url": configData.URL.String()}).
					Msg("unknown config section in config")
				continue
			}
			return DiConfig{}, err
		}

		// Create a new config.
		aNew, err := NewConfig(plugin)
		if err != nil {
			return DiConfig{}, err
		}

		// Create a marshaler for this section.
		buf := bytes.Buffer{}
		if err := configData.Marshaler.Init(bufio.NewReader(&buf), bufio.NewWriter(&buf)); err != nil {
			return DiConfig{}, err
		}

		// Encode this section into the bufr.
		if err := configData.Marshaler.EncodeSocket(data); err != nil {
			return DiConfig{}, err
		}

		// Decode this section from the bufr.
		if err := configData.Marshaler.DecodeSocket(aNew); err != nil {
			return DiConfig{}, err
		}

		merger, ok := aNew.(config.ConfigMerge)
		if !ok {
			return DiConfig{}, config.ErrUnknownConfig
		}

		// Merge previous config into this one.
		if err := merger.MergePrevious(previousConfig); err != nil {
			return DiConfig{}, err
		}

		// Update previousConfig for the next run
		previousConfig = aNew
	}

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
