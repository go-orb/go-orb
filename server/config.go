package server

// DefaultConfigSection is the section key used in config files used to
// configure the server options.
var DefaultConfigSection = "server" //nolint:gochecknoglobals

// MiddlewareConfig is the base config for all middlewares.
type MiddlewareConfig struct {
	Plugin string `json:"plugin,omitempty" yaml:"plugin,omitempty"`
}

// EntrypointConfigType is here so we can do magic work on options.
type EntrypointConfigType interface {
	config() *EntrypointConfig
}

// Option is a functional entrypoint option.
type Option func(EntrypointConfigType)

var _ (EntrypointConfigType) = (*EntrypointConfig)(nil)

// EntrypointConfig is the base config for all entrypoints.
type EntrypointConfig struct {
	Plugin  string `json:"plugin"         yaml:"plugin"`
	Enabled bool   `json:"enabled"        yaml:"enabled"`
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`

	OptMiddlewares []Middleware       `json:"-" yaml:"-"`
	OptHandlers    []RegistrationFunc `json:"-" yaml:"-"`
}

func (e *EntrypointConfig) config() *EntrypointConfig {
	return e
}

// NewEntrypointConfig creates a new entrypoint config with the given opts.
func NewEntrypointConfig(opts ...Option) *EntrypointConfig {
	cfg := &EntrypointConfig{
		Enabled: true,
	}

	for _, option := range opts {
		option(cfg)
	}

	return cfg
}

// WithEntrypointName sets the name of the entrypoint.
func WithEntrypointName(p string) Option {
	return func(cfg EntrypointConfigType) {
		c := cfg.config()
		c.Name = p
	}
}

// WithEntrypointPlugin sets the plugin of the entrypoint.
func WithEntrypointPlugin(p string) Option {
	return func(cfg EntrypointConfigType) {
		c := cfg.config()
		c.Plugin = p
	}
}

// WithEntrypointDisabled disables an entrypoint.
func WithEntrypointDisabled() Option {
	return func(cfg EntrypointConfigType) {
		c := cfg.config()
		c.Enabled = false
	}
}

// WithEntrypointMiddlewares appends the given middlewares.
func WithEntrypointMiddlewares(mws ...Middleware) Option {
	return func(cfg EntrypointConfigType) {
		c := cfg.config()
		c.OptMiddlewares = append(c.OptMiddlewares, mws...)
	}
}

// WithEntrypointHandlers appends the given handlers.
func WithEntrypointHandlers(hs ...RegistrationFunc) Option {
	return func(cfg EntrypointConfigType) {
		c := cfg.config()
		c.OptHandlers = append(c.OptHandlers, hs...)
	}
}

// Config is the global config for servers.
type Config struct {
	Middlewares []MiddlewareConfig `json:"middlwares,omitempty"  yaml:"middlewares,omitempty"`
	Handlers    []string           `json:"handlers,omitempty"    yaml:"handlers,omitempty"`
	Entrypoints []EntrypointConfig `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`

	functionalEntrypoints []EntrypointConfigType `json:"-" yaml:"-"`
}

// NewConfig creates a new config struct with the given opts.
func NewConfig(opts ...ConfigOption) Config {
	cfg := Config{}

	for _, option := range opts {
		option(&cfg)
	}

	return cfg
}

// ConfigOption allows to set options for MyConfig.
type ConfigOption func(*Config)

// WithEntrypointConfig allows you to create an entrypoint functionally.
func WithEntrypointConfig(config EntrypointConfigType) ConfigOption {
	return func(c *Config) {
		c.functionalEntrypoints = append(c.functionalEntrypoints, config)
	}
}
