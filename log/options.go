package log

// Options are compile time options for log.
type Options struct {
	// Parent starts a logger with a parent logger.
	Parent Logger

	// InternalParent starts the logger with Parent of its own internal type for example (zerolog.Logger).
	InternalParent any
}

// Option represents a single option.
type Option func(*Options)

// NewOptions merges Option... to Options.
func NewOptions(opts ...Option) Options {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithParent starts a logger with a parent logger.
func WithParent(n Logger) Option {
	return func(o *Options) {
		o.Parent = n
	}
}

// WithInternalParent starts the logger with Parent of its own internal type for example (zerolog.Logger).
func WithInternalParent(n any) Option {
	return func(o *Options) {
		o.InternalParent = n
	}
}
