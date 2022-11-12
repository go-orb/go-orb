// Package component implements the types used by components.
package component

import "context"

// Type is the type of a component.
type Type string

// Component needs to be implemented by every component.
type Component interface {
	// Start the component. E.g. connect to the broker.
	Start() error

	// Stop the component. E.g. disconnect from the broker.
	// The context will contain a timeout, and cancelation should be respected.
	Stop(context.Context) error

	// Type returns the component type, e.g. broker
	Type() Type

	// String returns the component plugin name.
	String() string
}
