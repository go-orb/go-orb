package types

import "fmt"

// Component needs to be implemented by every component.
type Component interface {
	fmt.Stringer

	// Start the component. E.g. connect to the broker
	Start() error
	// Stop the component. E.g. disconnect from the broker.
	Stop() error
	// Type returns the component type, e.g. broker
	Type() string
}
