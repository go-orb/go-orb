package cli

import (
	"fmt"
	"os"

	"jochum.dev/orb/orb/config"
	"jochum.dev/orb/orb/di"
)

// DiParsed is a marker that indicates that the cli has been parsed.
type DiParsed struct{}

func ProvideCli(
	config Config,
) (Cli, error) {
	p, _, err := Plugins.Get(config.GetPlugin())
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
	for _, f := range config.GetFlags() {
		if err := c.Add(f.AsOptions()...); err != nil {
			return result, err
		}
	}

	if config.GetNoFlags() == nil || !*config.GetNoFlags() {
		if err := c.Add(
			Name(PrefixName(config.GetArgPrefix(), "config")),
			Usage("Config urls"),
			Default(config.GetConfig()),
			EnvVars(PrefixEnv(config.GetArgPrefix(), "config")),
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
	if cfg.GetNoFlags() == nil || *cfg.GetNoFlags() {
		aDefConfig, err := NewConfig(cfg.GetPlugin())
		if err != nil {
			return di.DiConfig{}, err
		}

		defConfig, ok := aDefConfig.(Config)
		if !ok {
			return di.DiConfig{}, config.ErrUnknownConfig
		}

		if f, ok := c.Get("config"); ok {
			defConfig.SetConfig(FlagValue(f, defConfig.GetConfig()))
		}

		if err := defConfig.MergePrevious(cfg); err != nil {
			return di.DiConfig{}, err
		}
	}

	return di.DiConfig(cfg.GetConfig()), nil
}
