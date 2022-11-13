package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"

	"go-micro.dev/v5/log"
	"go-micro.dev/v5/types"
	"go-micro.dev/v5/types/component"
	"go-micro.dev/v5/util/container"
)

// TODO: implement

func init() {
	if err := Plugins.Add("mock", NewEntrypointMock); err != nil {
		panic(err)
	}

	if err := Plugins.Add("mock-two", NewEntrypointMock); err != nil {
		panic(err)
	}

	if err := NewDefaults.Add("mock", NewDefaultMockConfig); err != nil {
		panic(err)
	}

	if err := log.Plugins.Add("textstderr", NewHandlerStderr); err != nil {
		panic(err)
	}
}

func TestMock(t *testing.T) {
	ep1 := "mock-" + uuid.NewString()
	ep2 := "mock-" + uuid.NewString()

	srv, err := setupServer(
		WithMockEntrypoint(
			WithMockName(ep1),
			WithTest(t),
			WithDebugLog(false),
		),
		WithMockEntrypoint(
			WithMockName(ep2),
			WithTest(t),
			WithDebugLog(false),
		),
	)
	require.NoError(t, err, "failed to create server")

	require.Equal(t, len(srv.entrypoints), 2, "expected 2 entrypoints")
	require.NotNil(t, srv.entrypoints[ep1], "entrypoint 1 not found")
	require.NotNil(t, srv.entrypoints[ep2], "entrypoint 2 not found")

	// Check if all entrypoints started
	require.NoError(t, srv.Start(), "failed to start server")
	count, err := startCounter.Get(ep1)
	require.NoError(t, err, "failed to fetch start count, ep has not been started")
	require.Equal(t, count, 1, "sever should have been started")
	count, err = startCounter.Get(ep2)
	require.NoError(t, err, "failed to fetch start count, ep has not been started")
	require.Equal(t, count, 1, "sever should have been started")

	// Check if all entrypoints stopped
	require.NoError(t, srv.Stop(context.Background()), "failed to stop server")
	count, err = stopCounter.Get(ep1)
	require.NoError(t, err, "failed to fetch stop count, ep has not been stopped")
	require.Equal(t, count, 1, "sever should have been stopped")
	count, err = stopCounter.Get(ep2)
	require.NoError(t, err, "startedfailed to fetch stop count, ep has not been stopped")
	require.Equal(t, count, 1, "sever should have been stopped")
}

func TestInvalidEntrypoint(t *testing.T) {
	srv, err := setupServer(
		WithInvalidEntrypoint(),
	)
	t.Logf("expected error: %v", err)
	t.Logf("expected error: %v", srv.Start())
	require.Error(t, err, "invalid entrypoint, should error")

	srv, err = setupServer(
		WithInvalidConfigEntrypoint(),
	)
	t.Logf("expected error: %v", err)
	t.Logf("expected error: %v", srv.Start())
	require.Error(t, err, "invalid entrypoint, should error")
}

func TestStartStopError(t *testing.T) {
	// Startup with error.
	srv, err := setupServer(
		WithMockEntrypoint(
			WithMockName("mock-"+uuid.NewString()),
			WithTest(t),
			WithDebugLog(false),
			WithStartError(),
		),
	)
	require.NoError(t, err, "server setup failed")

	err = srv.Start()
	t.Logf("expected error: %v", err)
	require.Error(t, err, "startup should fail")

	// Shutdown with error.
	srv, err = setupServer(
		WithMockEntrypoint(
			WithMockName("mock-"+uuid.NewString()),
			WithTest(t),
			WithDebugLog(false),
			WithStopError(),
		),
	)
	require.NoError(t, err, "server setup failed")
	require.NoError(t, srv.Start(), "startup should fail")
	err = srv.Stop(context.Background())
	t.Logf("expected error: %v", err)
	require.Error(t, err, "stop should fail")
}

func setupServer(opts ...Option) (MicroServer, error) {
	var service types.ServiceName = "test-service"

	logger, err := log.ProvideLogger(service, nil)
	if err != nil {
		return MicroServer{}, fmt.Errorf("failed to setup logger: %w", err)
	}

	srv, err := ProvideServer(service, nil, logger, opts...)
	if err != nil {
		return MicroServer{}, fmt.Errorf("failed to setup server: %w", err)
	}

	return srv, nil
}

var _ (Entrypoint) = (*EntrypointMock)(nil)

type MockOption func(*ConfigMock)

var startCounter = container.NewSafeMap[int]()
var stopCounter = container.NewSafeMap[int]()

type ConfigMock struct {
	t          *testing.T
	Name       string
	debugLog   bool
	startError bool
	stopError  bool
}

type EntrypointMock struct {
	config  ConfigMock
	started bool
}

func NewDefaultMockConfig(service types.ServiceName, data types.ConfigData) (any, error) {
	return ConfigMock{
		Name:     "mock-" + uuid.NewString(),
		debugLog: false,
	}, nil
}

func NewEntrypointMock(
	name string,
	service types.ServiceName,
	data types.ConfigData,
	logger log.Logger,
	c any,
) (Entrypoint, error) {
	cfg, ok := c.(ConfigMock)
	if !ok {
		return nil, errors.New("invalid config, not of type ConfigMock")
	}

	cfg.Name = name

	if cfg.t == nil {
		return nil, fmt.Errorf("test not set for entrypoint %s", name)
	}

	if cfg.debugLog {
		cfg.t.Logf("creating entrypoint %s", name)
	}

	return &EntrypointMock{
		config:  cfg,
		started: false,
	}, nil
}

// Start the component. E.g. connect to the broker.
func (m *EntrypointMock) Start() error {
	m.started = true

	if m.config.startError {
		return errors.New("oops, some error occurred")
	}

	if m.config.debugLog {
		m.config.t.Logf("starting entrypoint %s", m.config.Name)
	}

	count, err := startCounter.Get(m.config.Name)
	if err != nil {
		count = 0
	}

	startCounter.Upsert(m.config.Name, count+1)

	return nil
}

// Stop the component. E.g. disconnect from the broker.
// The context will contain a timeout, and cancelation should be respected.
func (m *EntrypointMock) Stop(_ context.Context) error {
	m.started = false

	if m.config.stopError {
		return errors.New("oops, some error occurred")
	}

	if m.config.debugLog {
		m.config.t.Logf("stopping entrypoint %s", m.config.Name)
	}

	count, err := stopCounter.Get(m.config.Name)
	if err != nil {
		count = 0
	}

	stopCounter.Upsert(m.config.Name, count+1)

	return nil
}

// Type returns the component type, e.g. broker.
func (m *EntrypointMock) Type() component.Type {
	return ComponentType
}

// String returns the component plugin name.
func (m *EntrypointMock) String() string {
	return "mock"
}

func (m *EntrypointMock) Name() string { return m.config.Name }

func (m *EntrypointMock) Register(r RegistrationFunc) {}

func (c *ConfigMock) ApplyOptions(options ...MockOption) {
	for _, option := range options {
		option(c)
	}
}

func WithMockName(name string) MockOption {
	return func(c *ConfigMock) {
		c.Name = name
	}
}

func WithTest(t *testing.T) MockOption {
	return func(c *ConfigMock) {
		c.t = t
	}
}

func WithDebugLog(debug bool) MockOption {
	return func(c *ConfigMock) {
		c.debugLog = debug
	}
}

func WithStartError() MockOption {
	return func(c *ConfigMock) {
		c.startError = true
	}
}

func WithStopError() MockOption {
	return func(c *ConfigMock) {
		c.stopError = true
	}
}

func WithInvalidEntrypoint() Option {
	return func(c *Config) {
		c.Templates["invalid-"+uuid.NewString()] = EntrypointTemplate{
			Type:   "fake",
			Config: struct{}{},
		}
	}
}

func WithInvalidConfigEntrypoint() Option {
	return func(c *Config) {
		c.Templates["invalid-"+uuid.NewString()] = EntrypointTemplate{
			Type:   "mock-two",
			Config: struct{}{},
		}
	}
}

func WithMockEntrypoint(options ...MockOption) Option {
	return func(c *Config) {
		cfgAny, ok := c.Defaults["mock"]
		if !ok {
			// Should never happen, but just in case.
			panic("no defaults for mock entrypoint found")
		}

		cfg, ok := cfgAny.(ConfigMock)
		if !ok {
			// Should never happen, but just in case.
			panic("default config for mock entrypoint is not of type mock.Config")
		}

		cfg.ApplyOptions(options...)

		c.Templates[cfg.Name] = EntrypointTemplate{
			Type:   "mock",
			Config: cfg,
		}
	}
}

// NewHandlerStderr writes text to stderr.
func NewHandlerStderr(level slog.Leveler) (slog.Handler, error) {
	return slog.HandlerOptions{Level: level}.NewJSONHandler(os.Stderr), nil
}
