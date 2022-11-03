// Package component implements the types used by components.
package component

import "fmt"

// Type is the type of a component.
type Type string

// Component needs to be implemented by every component.
type Component interface {
	fmt.Stringer

	// Start the component. E.g. connect to the broker
	Start() error
	// Stop the component. E.g. disconnect from the broker.
	Stop() error
	// Type returns the component type, e.g. broker
	Type() Type
}
