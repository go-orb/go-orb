package types

import (
	"context"

	"github.com/go-orb/go-orb/util/container"
)

// Priority constants.
const (
	PriorityLogger   = 1000
	PriorityMetrics  = 1100
	PriorityRegistry = 1200
	PriorityEvent    = 1300
	PriorityServer   = 1400
	PriorityClient   = 1500
)

// Component needs to be implemented by every component.
type Component interface {
	// Start the component. E.g. connect to the broker.
	Start() error

	// Stop the component. E.g. disconnect from the broker.
	// The context will contain a timeout, and cancelation should be respected.
	Stop(ctx context.Context) error

	// Type returns the component type, e.g. broker.
	Type() string

	// String returns the component plugin name.
	String() string
}

// Components is the container for client implementations.
//
//nolint:gochecknoglobals
var Components = container.NewPriorityList[Component]()

// RegisterComponent adds a component to the container.
func RegisterComponent(component Component, priority int) error {
	return Components.Add(component, priority)
}
