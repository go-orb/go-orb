package registry

import (
	"errors"

	"go-micro.dev/v5/config/source/cli"
	"go-micro.dev/v5/log"
)

//nolint:gochecknoglobals
var (
	DefaultRegistry = "mdns"
	DefaultTimeout  = 600
)

func init() {
	err := cli.Flags.Add(cli.NewFlag(
		"registry",
		DefaultRegistry,
		cli.ConfigPathSlice([]string{"registry", "plugin"}),
		cli.Usage("Registry for discovery. etcd, mdns"),
		cli.EnvVars("REGISTRY"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}

	err = cli.Flags.Add(cli.NewFlag(
		"registry_timout",
		DefaultTimeout,
		cli.ConfigPathSlice([]string{"registry", "timeout"}),
		cli.Usage("Registry timeout."),
		cli.EnvVars("REGISTRY_TIMEOUT"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}
}

// TODO: this config misses stuff compared to v4, should that stuff be added here?

// Config is the configuration that can be used in a registry.
type Config struct {
	Plugin  string     `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	Timeout int        `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Logger  log.Logger `json:"logger,omitempty" yaml:"logger,omitempty"`
}

// NewConfig creates a new default config to use with a registry.
func NewConfig() Config {
	return Config{
		Plugin:  DefaultRegistry,
		Timeout: DefaultTimeout,
	}
}
