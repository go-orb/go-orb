package server

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"

	_ "github.com/go-micro/plugins/codecs/yaml"
	"github.com/go-micro/plugins/config/source/file"

	"go-micro.dev/v5/config"
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/types"
	"go-micro.dev/v5/types/component"
	"go-micro.dev/v5/util/container"
)

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

var configFile = `---
com:
  example:
    test-service:
      server:
        mock:
          fieldOne: "abc-field-one"
          fieldThree: true
          entrypoints:
            - name: mock-ep-1
              fieldTwo: 9
              fieldThree: false
            - name: mock-ep-2
              fieldOne: "def-field-one"
            - name: mock-ep-3
              enabled: false
            - name: mock-ep-4
              inherit: mock-ep-1
    another-test-two:
      server:
        mock:
          enabled: false
          entrypoints:
            - name: mock-ep-1
    another-test-three:
      server:
        mock:
          entrypoints:
            - name: mock-ep
              inherit: mock-ep-fake
    another-test-four:
      server:
        mock:
          entrypoints: fake
    another-test-five:
      server:
        mock:
          entrypoints:
            - 5
`

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

	// Validate entrypoints.
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

	e, err := srv.GetEntrypoint(ep1)
	require.NoError(t, err, "failed to fetch entrypoint 1")
	require.Equal(t, e.Name(), ep1)

	e, err = srv.GetEntrypoint(ep2)
	require.NoError(t, err, "failed to fetch entrypoint 2")
	require.Equal(t, e.Name(), ep2)

	_, err = srv.GetEntrypoint("fake")
	require.Error(t, err, "fetching invalid entrypoint should return error")
}

func TestMockConfigFile(t *testing.T) {
	data, err := config.Read([]*url.URL{file.TempFile([]byte(configFile), "yaml")}, nil)
	require.NoError(t, err, "failed to read config data")

	var service types.ServiceName = "com.example.test-service"

	logger, err := log.ProvideLogger(service, nil)
	require.NoError(t, err, "failed to setup logger")

	srv, err := ProvideServer(service, data, logger, WithMockDefaults(WithTest(t)))
	require.NoError(t, err, "failed to setup server")
	require.NoError(t, srv.Start(), "failed to start server")

	// Validate entrypoints.
	ep, err := srv.GetEntrypoint("mock-ep-1")
	require.NoError(t, err, "failed to retrieve mock-ep-1 entrypoint")
	epCfg := ep.(*EntrypointMock).config
	require.Equal(t, "abc-field-one", epCfg.FieldOne)
	require.Equal(t, 9, epCfg.FieldTwo)
	require.Equal(t, false, epCfg.FieldThree)
	require.Equal(t, 5, epCfg.FieldFour)

	ep, err = srv.GetEntrypoint("mock-ep-2")
	require.NoError(t, err, "failed to retrieve mock-ep-2 entrypoint")
	epCfg = ep.(*EntrypointMock).config
	require.Equal(t, "def-field-one", epCfg.FieldOne)
	require.Equal(t, 0, epCfg.FieldTwo)
	require.Equal(t, true, epCfg.FieldThree)
	require.Equal(t, 5, epCfg.FieldFour)

	_, err = srv.GetEntrypoint("mock-ep-3")
	require.Error(t, err, "should not be able to retrieve mock-ep-3 entrypoint")
	require.NoError(t, srv.Stop(context.Background()), "failed to start server")

	ep, err = srv.GetEntrypoint("mock-ep-4")
	require.NoError(t, err, "failed to retrieve mock-ep-4 entrypoint")
	epCfg = ep.(*EntrypointMock).config
	require.Equal(t, "abc-field-one", epCfg.FieldOne)
	require.Equal(t, 9, epCfg.FieldTwo)
	require.Equal(t, false, epCfg.FieldThree)
	require.Equal(t, 5, epCfg.FieldFour)

	_, err = srv.GetEntrypoint("mock-ep-5")
	require.Error(t, err, "should fail to retrieve mock-ep-5 entrypoint")

	// Test Service Two, all entrypoints disabled.
	service = "com.example.another-test-two"

	logger, err = log.ProvideLogger(service, nil)
	require.NoError(t, err, "failed to setup logger")

	srv, err = ProvideServer(service, data, logger, WithMockDefaults(WithTest(t)))
	require.NoError(t, err, "failed to setup server")
	require.NoError(t, srv.Start(), "failed to start server")

	_, err = srv.GetEntrypoint("mock-ep-1")
	require.Error(t, err, "should not be able to retrieve mock-ep-1 entrypoint")
	require.NoError(t, srv.Stop(context.Background()), "failed to start server")

	// Test Services containing errors.
	shouldError := []types.ServiceName{
		"com.example.another-test-three",
		"com.example.another-test-four",
		"com.example.another-test-five",
	}

	for _, service := range shouldError {
		logger, err = log.ProvideLogger(service, nil)
		require.NoError(t, err, "failed to setup logger")

		srv, err = ProvideServer(service, data, logger, WithMockDefaults(WithTest(t)))
		t.Logf("expected error: %v", err)
		require.Error(t, err, "should fail to setup server")
	}
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
	FieldOne   string `json:"fieldOne,omitempty" yaml:"fieldOne,omitempty"`
	FieldTwo   int    `json:"fieldTwo,omitempty" yaml:"fieldTwo,omitempty"`
	FieldThree bool   `json:"fieldThree,omitempty" yaml:"fieldThree,omitempty"`
	FieldFour  int    `json:"fieldFour,omitempty" yaml:"fieldFour,omitempty"`
}

type EntrypointMock struct {
	config  ConfigMock
	started bool
}

func NewDefaultMockConfig(service types.ServiceName, data types.ConfigData) (any, error) {
	cfg := ConfigMock{
		Name:      "mock-" + uuid.NewString(),
		debugLog:  false,
		FieldFour: 5,
	}

	sections := types.SplitServiceName(service)
	if err := config.Parse(append(sections, DefaultConfigSection, "mock"), data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c *ConfigMock) GetName() string {
	return c.Name
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
		return nil, fmt.Errorf("create mock entrypoint: invalid config, not of type ConfigMock, but '%T'", c)
	}

	cfg.Name = name

	if err := ParseEntrypointConfig(service, data, "mock", &cfg); err != nil {
		return nil, err
	}

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

func WithTest(t *testing.T) MockOption { //nolint:thelper
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

func WithMockDefaults(opts ...MockOption) Option {
	return func(c *Config) {
		cfg, ok := c.Defaults["mock"].(ConfigMock)
		if !ok {
			// Should never happen.
			panic(fmt.Errorf("mock.WithDefaults received invalid type, not ConfigMock, but '%T'", cfg))
		}

		cfg.ApplyOptions(opts...)
		c.Defaults["mock"] = cfg
	}
}

// NewHandlerStderr writes text to stderr.
func NewHandlerStderr(level slog.Leveler) (slog.Handler, error) {
	return slog.HandlerOptions{Level: level}.NewJSONHandler(os.Stderr), nil
}
