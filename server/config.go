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

// WithDisabled disables an entrypoint.
func WithDisabled() Option {
	return func(cfg EntrypointConfigType) {
		c := cfg.config()
		c.Enabled = false
	}
}

// WithMiddlewares appends the given middlewares.
func WithMiddlewares(mws ...Middleware) Option {
	return func(cfg EntrypointConfigType) {
		c := cfg.config()
		c.OptMiddlewares = append(c.OptMiddlewares, mws...)
	}
}

// WithHandlers appends the given handlers.
func WithHandlers(hs ...RegistrationFunc) Option {
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

// ConfigOption allows to set options for MyConfig.
type ConfigOption func(*Config)

// WithEntrypointConfig allows you to create an entrypoint functionally.
func WithEntrypointConfig(config EntrypointConfigType) ConfigOption {
	return func(c *Config) {
		c.functionalEntrypoints = append(c.functionalEntrypoints, config)
	}
}
