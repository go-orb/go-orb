//go:build go1.18
// +build go1.18

package cli

import (
	"fmt"
	"strings"
)

// Flag is a Cli Flag and maybe environment variable.
type Flag struct {
	Name    string
	EnvVars []string
	Usage   string

	// The path in map(\[string\])+any
	ConfigPath []string

	Default any
	Value   any
}

func (f *Flag) String() string {
	return f.Name
}

// FlagOption is an option for NewFlag.
type FlagOption func(*Flag)

// ConfigPath sets the ConfigPath for the flag.
func ConfigPath(n string) FlagOption {
	return func(o *Flag) {
		o.ConfigPath = strings.Split(n, ".")
	}
}

// ConfigPathSlice is the same as ConfigPath but it accepts a slice.
func ConfigPathSlice(n []string) FlagOption {
	return func(o *Flag) {
		o.ConfigPath = n
	}
}

// EnvVars set's environment variables for the flag.
func EnvVars(n ...string) FlagOption {
	return func(o *Flag) {
		o.EnvVars = n
	}
}

// Usage set's the usage string for the flag.
func Usage(n string) FlagOption {
	return func(o *Flag) {
		o.Usage = n
	}
}

// Default sets the flags default.
func Default[T any](n T) FlagOption {
	return func(o *Flag) {
		o.Default = n
	}
}

// FlagValue gets a value back from a Flag and enforces types.
func FlagValue[T any](f *Flag) (T, error) {
	switch t := f.Value.(type) {
	case T:
		return t, nil
	default:
		var tmp T
		return tmp, fmt.Errorf("mixed types: %v", f.Value)
	}
}

// NewFlag creates a new flag.
func NewFlag[T any](
	name string,
	defaultValue T,
	opts ...FlagOption,
) *Flag {
	options := Flag{
		Name:    name,
		Default: defaultValue,
	}

	for _, o := range opts {
		o(&options)
	}

	return &options
}
