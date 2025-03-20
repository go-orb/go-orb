package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// MainActionName is the name of the "main" action.
const MainActionName = "__main"

func configureCommands(appContext *AppContext, commands []*Command, previousCommands []string) {
	for _, c := range commands {
		if !c.NoAction {
			c.InternalAction = func() error {
				appContext.SelectedCommand = append(previousCommands, c.Name) //nolint:gocritic
				appContext.SelectedService = c.Service

				return nil
			}
		}

		if c.Subcommands != nil {
			configureCommands(appContext, c.Subcommands, append(previousCommands, c.Name))
		}
	}
}

// AppContext holds app state of the main application for multiple services.
type AppContext struct {
	app        *App
	cancelFunc context.CancelFunc

	SelectedService string
	SelectedCommand []string

	Context       context.Context
	StopWaitGroup *sync.WaitGroup
}

// App returns the application.
func (c *AppContext) App() *App {
	return c.app
}

// ProcessID returns the process ID of the application.
func (c *AppContext) ProcessID() string {
	return c.app.ProcessID()
}

// Name returns the name of the application.
func (c *AppContext) Name() string {
	return c.app.Name
}

// Version returns the version of the application.
func (c *AppContext) Version() string {
	return c.app.Version
}

// ExitGracefully exits the application gracefully.
func (c *AppContext) ExitGracefully(exitCode int) {
	c.cancelFunc()
	c.StopWaitGroup.Wait()
	os.Exit(exitCode)
}

// NewAppContext creates a new Application context.
func NewAppContext(app *App) *AppContext {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	appContext := &AppContext{
		app: app,

		Context:       ctx,
		cancelFunc:    cancel,
		StopWaitGroup: &sync.WaitGroup{},
	}

	if !app.NoAction {
		app.InternalAction = func() error {
			appContext.SelectedCommand = []string{MainActionName}
			appContext.SelectedService = app.Name

			return nil
		}
	}

	if !app.NoGlobalConfig {
		app.Flags = append(app.Flags, &Flag{
			Name:    "config",
			Usage:   "Config file(s)",
			Default: []string{},
			EnvVars: []string{"CONFIG"},
		})
	}

	if app.Commands != nil {
		configureCommands(appContext, app.Commands, []string{})
	}

	return appContext
}

// ServiceContext holds the service state.
type ServiceContext struct {
	appContext *AppContext
	name       string
	version    string

	Config map[string]any
}

// App returns the application.
func (c *ServiceContext) App() *App {
	return c.appContext.App()
}

// AppName returns the name of the application.
func (c *ServiceContext) AppName() string {
	return c.appContext.Name()
}

// AppVersion returns the version of the application.
func (c *ServiceContext) AppVersion() string {
	return c.appContext.Version()
}

// Context returns the context of the application.
func (c *ServiceContext) Context() context.Context {
	return c.appContext.Context
}

// StopWaitGroup returns the stop wait group of the application.
func (c *ServiceContext) StopWaitGroup() *sync.WaitGroup {
	return c.appContext.StopWaitGroup
}

// ProcessID returns the process ID of the application.
func (c *ServiceContext) ProcessID() string {
	return c.appContext.ProcessID()
}

// Name returns the name of the service.
func (c *ServiceContext) Name() string {
	return c.name
}

// Version returns the version of the service.
func (c *ServiceContext) Version() string {
	return c.version
}

// ExitAppGracefully exits the application gracefully.
func (c *ServiceContext) ExitAppGracefully(exitCode int) {
	c.appContext.ExitGracefully(exitCode)
}

// NewServiceContext creates a new Service context for the given service.
func NewServiceContext(appContext *AppContext, name string, version string) *ServiceContext {
	return &ServiceContext{
		appContext: appContext,
		name:       name,
		version:    version,
	}
}
