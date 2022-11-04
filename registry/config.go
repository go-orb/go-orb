package registry

import (
	"github.com/go-orb/config/source/cli"
	"go-micro.dev/v5/log"
)

const ComponentType = "registry"

var DefaultRegistry = "mdns"
var DefaultTimout = 600

func init() {
	err := cli.Flags.Add(cli.NewFlag(
		"registry",
		DefaultRegistry,
		cli.CPSlice([]string{"registry", "plugin"}),
		cli.Usage("Registry for discovery. etcd, mdns"),
	))
	if err != nil {
		panic(err)
	}

	err = cli.Flags.Add(cli.NewFlag(
		"registry_timout",
		DefaultTimout,
		cli.CPSlice([]string{"registry", "timeout"}),
		cli.Usage("Registry timeout."),
	))
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Plugin  string      `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	Timeout int         `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Logger  *log.Config `json:"logger,omitempty" yaml:"logger,omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}
