package server

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"jochum.dev/orb/orb/config/chelp"
	"jochum.dev/orb/orb/log"
)

const (
	configKeyLogger   = "logger"
	configKeyID       = "id"
	configKeyName     = "name"
	configKeyVersion  = "version"
	configKeyMetadata = "metadata"
)

type Config interface {
	chelp.PluginConfig

	// Required
	Name() string
	Version() string

	// Optional
	Logger() any
	ID() string
	Metadata() map[string]string
	Address() string
	RegisterTTL() int
	RegisterInterval() int

	SetName(n string)
	SetVersion(n string)
	SetLogger(n any)
	SetID(n string)
	SetMetadata(n map[string]string)
	SetAddress(n string)
	SetRegisterTTL(n int)
	SetRegisterInterval(n int)
}

type BaseConfig struct {
	*chelp.BasePluginConfig

	name    string
	version string

	logger           any
	id               string
	metadata         map[string]string
	address          string
	registerTTL      int
	registerInterval int
}

func NewConfig() *BaseConfig {
	return &BaseConfig{
		BasePluginConfig: chelp.NewPluginConfig(),
	}
}

func (c *BaseConfig) Load(m map[string]any) error {
	var result error

	// Required
	if err := c.BasePluginConfig.Load(m); err != nil {
		result = multierror.Append(err)
	}

	var err error
	if c.name, err = chelp.Get(m, configKeyName, ""); err != nil {
		result = multierror.Append(err)
	}

	if c.version, err = chelp.Get(m, configKeyVersion, ""); err != nil {
		result = multierror.Append(err)
	}

	// Optional
	c.logger, err = log.LoadConfig(m, configKeyLogger)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	if c.id, err = chelp.Get(m, configKeyID, ""); !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	c.metadata, err = chelp.Get(m, configKeyMetadata, map[string]string{})
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	return result
}

func (c *BaseConfig) Store(m map[string]any) error {
	var result error

	if err := c.BasePluginConfig.Store(m); err != nil {
		result = multierror.Append(err)
	}

	m[configKeyName] = c.name
	m[configKeyVersion] = c.version

	var err error

	m[configKeyLogger], err = log.StoreConfig(c.logger)
	if !errors.Is(err, chelp.ErrNotExistant) {
		result = multierror.Append(err)
	}

	m[configKeyID] = c.id
	m[configKeyMetadata] = c.metadata

	return result
}

func (c *BaseConfig) Name() string                { return c.name }
func (c *BaseConfig) Version() string             { return c.version }
func (c *BaseConfig) Logger() any                 { return c.logger }
func (c *BaseConfig) ID() string                  { return c.id }
func (c *BaseConfig) Metadata() map[string]string { return c.metadata }
func (c *BaseConfig) Address() string             { return c.address }
func (c *BaseConfig) RegisterTTL() int            { return c.registerTTL }
func (c *BaseConfig) RegisterInterval() int       { return c.registerInterval }

func (c *BaseConfig) SetName(n string)                { c.name = n }
func (c *BaseConfig) SetVersion(n string)             { c.version = n }
func (c *BaseConfig) SetLogger(n any)                 { c.logger = n }
func (c *BaseConfig) SetID(n string)                  { c.id = n }
func (c *BaseConfig) SetMetadata(n map[string]string) { c.metadata = n }
func (c *BaseConfig) SetAddress(n string)             { c.address = n }
func (c *BaseConfig) SetRegisterTTL(n int)            { c.registerTTL = n }
func (c *BaseConfig) SetRegisterInterval(n int)       { c.registerInterval = n }
