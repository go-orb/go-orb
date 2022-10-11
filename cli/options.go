package cli

import (
	"os"
	"path/filepath"
)

type Options struct {
	// Name is the name of the app
	Name string

	// Description is the description of the app
	Description string

	// Version is the Version of the app
	Version string

	// Usage is the apps usage string
	Usage string
}

type Option func(*Options)

func CliName(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func CliDescription(n string) Option {
	return func(o *Options) {
		o.Description = n
	}
}

func CliVersion(n string) Option {
	return func(o *Options) {
		o.Version = n
	}
}

func CliUsage(n string) Option {
	return func(o *Options) {
		o.Usage = n
	}
}

func NewCLIOptions(opts ...Option) *Options {
	options := &Options{
		Name:        filepath.Base(os.Args[0]),
		Description: "",
	}

	for _, o := range opts {
		o(options)
	}

	return options
}
