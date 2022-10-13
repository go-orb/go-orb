package cli

import (
	"errors"
	"fmt"
)

var ErrConfigIsNil = errors.New("config is nil")

type Cli interface {
	fmt.Stringer

	Init(config Config) error

	// Add adds a Flag to CLI
	Add(opts ...FlagOption) error

	// Get returns a flag
	Get(name string) (*Flag, bool)

	// Parse parses flags from args you MUST Add Flags first
	Parse(args []string) error
}
