// Package metrics provides a Wrapper around hashicorp/go-metrics.
package metrics

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
)

// ComponentType is the components name.
const ComponentType = "metrics"

// Label is Metrics label.
type Label struct {
	Name  string
	Value string
}

// Metrics contains all hashicorp/go-metrics.Metrics methods we have found so far.
//
//nolint:interfacebloat
type Metrics interface {
	types.Component

	SetGauge(key []string, val float32)                           // Set gauge key and value with 32 bit precision
	SetGaugeWithLabels(key []string, val float32, labels []Label) // Set gauge key and value with 32 bit precision
	SetPrecisionGauge(key []string, val float64)
	SetPrecisionGaugeWithLabels(key []string, val float64, labels []Label)
	EmitKey(key []string, val float32)
	IncrCounter(key []string, val float32)
	IncrCounterWithLabels(key []string, val float32, labels []Label)
	AddSample(key []string, val float32)
	AddSampleWithLabels(key []string, val float32, labels []Label)
	MeasureSince(key []string, start time.Time)
	MeasureSinceWithLabels(key []string, start time.Time, labels []Label)
	// UpdateFilter overwrites the existing filter with the given rules.
	UpdateFilter(allow, block []string)
	// UpdateFilterAndLabels overwrites the existing filter with the given rules.
	UpdateFilterAndLabels(allow, block, allowedLabels, blockedLabels []string)
	// Emits various runtime statsitics
	EmitRuntimeStats()
}

// Type is the registry type it is returned when you use metrics.Provide
// which selects a registry to use based on the plugin configuration.
type Type struct {
	Metrics
}

// Provide is the metrics provider for wire.
// It parses the config from "configs", fetches the "Plugin" from the config and
// then forwards all it's arguments to the factory which it get's from "Plugins".
func Provide(
	svcCtx *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
	opts ...Option,
) (Type, error) {
	cfg := NewConfig(opts...)

	if err := config.Parse(nil, DefaultConfigSection, svcCtx.Config, &cfg); err != nil {
		return Type{}, err
	}

	if cfg.Plugin == "" {
		logger.Warn("empty metrics plugin, using the default", "default", DefaultMetricsPlugin)
		cfg.Plugin = DefaultMetricsPlugin
	}

	logger.Debug("Metrics", "plugin", cfg.Plugin)

	provider, ok := Plugins.Get(cfg.Plugin)
	if !ok {
		return Type{}, fmt.Errorf("Metrics plugin '%s' not found, did you import it?", cfg.Plugin)
	}

	// Configure the logger.
	cLogger, err := logger.WithConfig([]string{DefaultConfigSection}, svcCtx.Config)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	instance, err := provider(svcCtx, cLogger, opts...)
	if err != nil {
		return Type{}, err
	}

	// Register metrics as a component.
	err = components.Add(&instance, types.PriorityMetrics)
	if err != nil {
		logger.Warn("while registering metrics as a component", "error", err)
	}

	return instance, nil
}
