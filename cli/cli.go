// Package cli is the Cli component of orb.
package cli

import (
	"errors"
	"fmt"
)

// ErrConfigIsNil indicates that the given config is nil.
var ErrConfigIsNil = errors.New("config is nil")

// Cli is the component interface for every Cli plugin.
type Cli interface {
	fmt.Stringer

	Init(config any) error
	Config() any

	// Add adds a Flag to CLI
	Add(opts ...FlagOption) error

	// Get returns a flag
	Get(name string) (*Flag, bool)

	// Parse parses flags from args you MUST Add Flags first
	Parse(args []string) error
}
