// Package cli provides the cli for go-orb.
package cli

// App represents a CLI Application.
type App struct {
	processID string

	Name           string
	Version        string
	Usage          string
	Commands       []*Command
	Flags          []*Flag
	NoAction       bool
	NoGlobalConfig bool

	// Internal
	InternalAction func() error
}

// ProcessID returns the process ID of the application.
func (a *App) ProcessID() string {
	return a.processID
}

// Command is a CLI Command for App.
type Command struct {
	Name        string
	Service     string
	Category    string
	Usage       string
	Flags       []*Flag
	Subcommands []*Command
	NoAction    bool

	// Internal
	InternalAction func() error
}
