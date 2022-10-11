package cli

type Cli interface {
	// Add adds a Flag to CLI
	Add(opts ...FlagOption) error

	// Get returns a flag
	Get(name string) (*Flag, bool)

	// Parse parses flags from args you MUST Add Flags first
	Parse(args []string, opts ...Option) error

	// String returns the name of the current implementation
	String() string
}
