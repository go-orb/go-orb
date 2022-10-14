package cli

import (
	"fmt"
	"os"

	"jochum.dev/orb/orb/di"
)

// DiParsed is a marker that indicates that the cli has been parsed.
type DiParsed struct{}

func ProvideCli(
	config Config,
) (Cli, error) {
	p, _, err := Plugins.Get(config.Plugin())
	if err != nil {
		return nil, fmt.Errorf("unknown cli given: %w", err)
	}

	return p(), nil
}

// ProvideParsed parses the cli and provides DiParsed.
func ProvideParsed(
	config Config,
	c Cli,
) (DiParsed, error) {
	result := DiParsed{}

	// User flags
	for _, f := range config.Flags() {
		if err := c.Add(f.AsOptions()...); err != nil {
			return result, err
		}
	}

	if config.NoFlags() == nil || !*config.NoFlags() {
		if err := c.Add(
			Name(PrefixName(config.ArgPrefix(), "config")),
			Usage("Config file"),
			Default(config.Config()),
			EnvVars(PrefixEnv(config.ArgPrefix(), "config")),
		); err != nil {
			return result, err
		}
	}

	// Initialize the CLI / parse flags
	if err := c.Parse(
		os.Args,
	); err != nil {
		return result, err
	}

	return result, nil
}

func ProvideConfig(
	_ di.DiFlags,
	diFlags DiParsed,
	c Cli,
	cfg Config,
) (di.DiConfig, error) {
	if cfg.NoFlags() == nil || *cfg.NoFlags() {
		defConfig := NewConfig()

		if f, ok := c.Get("config"); ok {
			defConfig.SetConfig(FlagValue(f, defConfig.Config()))
		}

		if err := cfg.Merge(defConfig); err != nil {
			return di.DiConfig{}, err
		}
	}

	return di.DiConfig(cfg.Config()), nil
}
