package registry

import (
	"github.com/orb-org/orb/log"
)

type Config struct {
	Plugin  string      `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	Timeout int         `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Logger  *log.Config `json:"logger,omitempty" yaml:"logger,omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}
