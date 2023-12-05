package server

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	uconfig "github.com/go-orb/go-orb/util/config"
	"github.com/go-orb/go-orb/util/slicemap"
)

// DefaultConfigSection is the section key used in config files used to
// configure the server options.
var DefaultConfigSection = "server" //nolint:gochecknoglobals

// Option is a functional HTTP server option.
type Option func(*Config)

// Config is the server config. It contains the list of addresses on which
// entrypoints will be created, and the default config used for each entrypoint.
type Config struct {
	// Enabled keeps track of whether server plugins are globally disabled or not.
	Enabled map[string]bool

	// Defaults is the list of defaults the each server plugin.
	// Provisioned with the factory methods registered by the entrypoint plugins.
	Defaults map[string]EntrypointConfig

	// Templates contains a set of entrypoint templates to create, indexed by name.
	//
	// Each entrypoint needs a unique name, as each entrypoint can be dynamically
	// configured by referencing the name. The default name used in an entrypoint
	// is the format of "http-<uuid>", used if no custom name is provided.
	Templates map[string]EntrypointTemplate
}

// NewConfig creates a new server config with default values as starting point,
// after which all the functional options are applied.
func NewConfig(options ...Option) Config {
	cfg := Config{
		Enabled:   make(map[string]bool, NewDefaults.Len()),
		Defaults:  make(map[string]EntrypointConfig, NewDefaults.Len()),
		Templates: make(map[string]EntrypointTemplate),
	}

	NewDefaults.Range(func(name string, newConfig NewDefault) bool {
		cfg.Enabled[name] = true
		cfg.Defaults[name] = newConfig()

		return true
	})

	// Apply options.
	for _, o := range options {
		o(&cfg)
	}

	return cfg
}

// UnmarshalJSON extracts the entrypoint configuration from a file config.
func (c *Config) UnmarshalJSON(data []byte) error {
	dataMap := map[string]any{}

	if err := json.Unmarshal(data, &dataMap); err != nil {
		return err
	}

	return c.parseEntrypointConfig(dataMap)
}

// UnmarshalYAML extracts the entrypoint configuration from a file config.
func (c *Config) UnmarshalYAML(data *yaml.Node) error {
	dataMap := map[string]any{}

	if err := data.Decode(dataMap); err != nil {
		return err
	}

	return c.parseEntrypointConfig(dataMap)
}

// Errors.
var (
	ErrInvalidEpSecType  = errors.New("config section entrypoints should be a list")
	ErrInvalidEpItemType = errors.New("config section entrypoints item should be a map")
	ErrInvalidServerType = errors.New("config section server plugins should be a map")
)

// parseEntrypointConfig extracts the entrypiont configurations from a file
// config.
//
// This function will look if a config is provided for an entrypoint, and then
// apply that config directly, instead of parsing all at once. We have to do this because the deserializers
// don't allow us to overwrite array elements, they will delete all eelements
// and create new ones. Therefore we have to only apply the speccific yaml
// that directly relates to the entrypoint.
func (c *Config) parseEntrypointConfig(dataMap map[string]any) error { //nolint:gocognit
	for serverType, d := range dataMap {
		if _, ok := Plugins.Get(serverType); !ok {
			return fmt.Errorf("server plugin %s does not exist, did you regiser it?", serverType)
		}

		serverConfig, ok := d.(map[string]any)
		if !ok {
			return ErrInvalidServerType
		}

		// Overlay defaults.
		if err := uconfig.OverlayMap(serverConfig, c.Defaults[serverType]); err != nil {
			return fmt.Errorf("overlay fileconfig for %s defaults: %w", serverType, err)
		}

		// Overlay on static entrypoints
		for _, staticEntrypoint := range c.Templates {
			if err := uconfig.OverlayMap(serverConfig, staticEntrypoint.Config); err != nil {
				return fmt.Errorf("overlay fileconfig for %s defaults: %w", serverType, err)
			}
		}

		// Check if server plugin has been globally disabled.
		if enabled, ok := slicemap.Get[bool](serverConfig, "enabled"); ok {
			c.Enabled[serverType] = enabled
		}

		entrypoints, ok := slicemap.Get[[]any](serverConfig, "entrypoints")
		if !ok {
			return fmt.Errorf("%s: %w", serverType, ErrInvalidEpSecType)
		}

		// Apply file config to each entrypoint individually.
		for _, epData := range entrypoints {
			ep, ok := epData.(map[string]any)
			if !ok {
				return ErrInvalidEpItemType
			}

			name, ok := slicemap.Get[string](ep, "name")
			if !ok {
				continue
			}

			// If entrypoint is only dynamically defined, create new template.
			if _, ok := c.Templates[name]; !ok {
				c.Templates[name] = EntrypointTemplate{
					Enabled: true,
					Type:    serverType,
					Config:  c.Defaults[serverType].Copy(),
				}
			}

			// Check if entrypoint has been disabled.
			if enabled, ok := slicemap.Get[bool](ep, "enabled"); ok {
				epCfg := c.Templates[name]
				epCfg.Enabled = enabled
				c.Templates[name] = epCfg
			}

			if err := overlayEntrypointConfig(c.Templates[name].Config, entrypoints, name); err != nil {
				return err
			}
		}
	}

	return nil
}

// overlayEntrypointConfig looks for the map[string]any that belongs to the
// specified entrypiont, and overlays it on top of the current config.
func overlayEntrypointConfig(cfg EntrypointConfig, entrypointList []any, name string) error {
	if cfg == nil {
		return nil
	}

	for _, epAny := range entrypointList {
		entrypoint, ok := epAny.(map[string]any)
		if !ok {
			return ErrInvalidEpItemType
		}

		if entrypoint["name"] != name {
			continue
		}

		if i, ok := entrypoint["inherit"]; ok {
			inherit, ok := i.(string)
			if !ok {
				return fmt.Errorf("field inherit should be of type string, not %T", i)
			}

			// Make sure we keep static address across inheritance if no manual
			// address is specified in config
			if _, ok := entrypoint["address"]; !ok {
				entrypoint["address"] = cfg.GetAddress()
			}

			// Overlay inherited config.
			if err := overlayEntrypointConfig(cfg, entrypointList, inherit); err != nil {
				return fmt.Errorf("%s is unable to inherit from %s config: %w", name, inherit, err)
			}

			// Only inherit once.
			delete(entrypoint, "inherit")

			// Overlay original file config again.
			// Apply the old name + address, and any entrypoint specific values.
			return overlayEntrypointConfig(cfg, entrypointList, name)
		}

		if err := uconfig.OverlayMap(entrypoint, cfg); err != nil {
			return fmt.Errorf("parse entrypoint config: %w", err)
		}

		return nil
	}

	return fmt.Errorf("unable to overlay config for %s, not found", name)
}
