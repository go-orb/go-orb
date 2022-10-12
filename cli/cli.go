package cli

import "errors"

var ErrConfigIsNil = errors.New("config is nil")

type Cli interface {
	Init(config Config) error

	// Add adds a Flag to CLI
	Add(opts ...FlagOption) error

	// Get returns a flag
	Get(name string) (*Flag, bool)

	// Parse parses flags from args you MUST Add Flags first
	Parse(args []string) error

	// String returns the name of the current implementation
	String() string
}
