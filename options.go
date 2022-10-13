package orb

import (
	"jochum.dev/orb/orb/cli"
)

type HookFunc func(service *Service) error

// Options for orb service.
type Options struct {
	ArgPrefix        string
	Address          string
	Name             string
	Description      string
	Version          string
	Usage            string
	NoFlags          bool
	ConfigFile       string
	RegisterTTL      int
	RegisterInterval int
	Metadata         map[string]string
	Flags            []cli.Flag

	// Wrappers
	// WrapSubscriber []server.SubscriberWrapper
	// WrapHandler    []server.HandlerWrapper
	// WrapCall       []client.CallWrapper
	// WrapClient     []client.Wrapper
	// OrigClient     client.Client

	// Before and After funcs
	Actions     []HookFunc
	BeforeStart []HookFunc
	BeforeStop  []HookFunc
	AfterStart  []HookFunc
	AfterStop   []HookFunc
}

type Option func(*Options)

func NewOptions(opts ...Option) *Options {
	options := &Options{}

	for _, o := range opts {
		o(options)
	}

	return options
}
