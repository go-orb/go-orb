package mdnsregistry

import (
	"jochum.dev/orb/orb/registry"
)

type Config interface {
	registry.Config

	Domain() string
	SetDomain(n string)
}

type ConfigImpl struct {
	*registry.BaseConfig

	domain string
}

func NewConfig() *ConfigImpl {
	return &ConfigImpl{
		BaseConfig: registry.NewConfig(),
	}
}

func (c *ConfigImpl) Domain() string     { return c.domain }
func (c *ConfigImpl) SetDomain(n string) { c.domain = n }
