package metrics

import (
	"os"
	"time"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
)

//
//nolint:gochecknoglobals
var (
	// DefaultMetricsPlugin is the default metrics plugin to use.
	DefaultMetricsPlugin = "memory"
	// DefaultConfigSection is the section key used in config files used to configure the metrics options.
	DefaultConfigSection = "metrics"
)

var _ (ConfigType) = (*Config)(nil)

// Option is a logger WithXXX Option.
type Option func(ConfigType)

// ConfigType is a wrapper for config, so we can pass it back to the this plugin handler.
type ConfigType interface {
	config() *Config
}

// Config contains the metrics config.
type Config struct {
	Plugin string // Plugin sets the metrics sink plugin to use.

	Hostname             string          // Hostname to use. If not provided and EnableHostname, it will be os.Hostname
	EnableHostname       bool            // Hostname to use. If not provided and EnableHostname, it will be os.Hostname
	EnableHostnameLabel  bool            // Enable prefixing gauge values with hostname
	EnableRuntimeMetrics bool            // Enables profiling of runtime metrics (GC, Goroutines, Memory)
	EnableTypePrefix     bool            // Prefixes key with a type ("counter", "gauge", "timer")
	TimerGranularity     config.Duration // Granularity of timers.
	ProfileInterval      config.Duration // Interval to profile runtime metrics

	AllowedPrefixes []string // A list of metric prefixes to allow, with '.' as the separator
	BlockedPrefixes []string // A list of metric prefixes to block, with '.' as the separator
	AllowedLabels   []string // A list of metric labels to allow, with '.' as the separator
	BlockedLabels   []string // A list of metric labels to block, with '.' as the separator
	FilterDefault   bool     // Whether to allow metrics by default
}

func (c *Config) config() *Config {
	return c
}

// NewConfig creates a new config with the defaults and applys opts on top.
func NewConfig(opts ...Option) *Config {
	cfg := &Config{
		Plugin:               DefaultMetricsPlugin,
		Hostname:             "",
		EnableHostname:       true,
		EnableRuntimeMetrics: true,
		EnableTypePrefix:     false,
		TimerGranularity:     config.Duration(time.Millisecond),
		ProfileInterval:      config.Duration(time.Second),
		FilterDefault:        true,
	}

	// Apply options.
	for _, o := range opts {
		o(cfg)
	}

	if cfg.EnableHostname && cfg.Hostname == "" {
		// Try to get the hostname
		name, err := os.Hostname()
		if err != nil {
			log.Error("While getting the hostname", err)
		} else {
			cfg.Hostname = name
		}
	}

	return cfg
}
