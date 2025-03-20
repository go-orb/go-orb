// Package cli provides the cli for go-orb.
package cli

// HardcodedConfig represents a hardcoded config with it's format.
// Format can be any of the importet codecs.
type HardcodedConfig struct {
	Format string
	Data   string
}

// App represents a CLI Application.
type App struct {
	processID string

	Name     string
	Version  string
	Usage    string
	Commands []*Command
	Flags    []*Flag

	// MultiServiceConfig defines if the config is used with sections or without.
	// If true, the config is used with sections.
	// For example:
	// ```yaml
	// service1: # service1 section
	//   logger:
	//     level: INFO
	// service2: # service2 section
	//   logger:
	//     level: INFO
	// ```
	//
	// If false, the config is used without sections.
	// For example:
	// ```yaml
	// logger:
	//   level: INFO
	// ```
	MultiServiceConfig bool

	// NoAction defines if there will be no main action.
	NoAction bool

	// NoGlobalConfig defines if the global config flag should be added and parsed.
	NoGlobalConfig bool

	// HardcodedConfigs defines the hardcoded configs.
	HardcodedConfigs []HardcodedConfig
	// HardcodedConfigURLs defines the hardcoded config URLs.
	HardcodedConfigURLs []string

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
