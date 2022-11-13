package server

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"go-micro.dev/v5/log"
	"go-micro.dev/v5/types"
	"go-micro.dev/v5/types/component"
)

// TODO: implement

func init() {
	Plugins.Add("mock", NewEntrypointMock)
	NewDefaults.Add("mock", NewDefaultMockConfig)
}

func TestMock(t *testing.T) {
	var service types.ServiceName = "test-service"

	logger, err := log.ProvideLogger(service, nil)
	if err != nil {
		t.Fatalf("failed to setup logger: %v", err)
	}

	ProvideServer(service, nil, logger)
}

var _ (Entrypoint) = (*EntrypointMock)(nil)

type ConfigMock struct {
	Name string
}

type EntrypointMock struct {
	config  ConfigMock
	started bool
}

func NewDefaultMockConfig(service types.ServiceName, data types.ConfigData) (any, error) {
	return ConfigMock{
		Name: "mock-" + uuid.NewString(),
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
		return nil, errors.New("invalid config")
	}

	return &EntrypointMock{
		config:  cfg,
		started: false,
	}, nil
}

// Start the component. E.g. connect to the broker.
func (m *EntrypointMock) Start() error {
	m.started = true
	return nil
}

// Stop the component. E.g. disconnect from the broker.
// The context will contain a timeout, and cancelation should be respected.
func (m *EntrypointMock) Stop(_ context.Context) error {
	m.started = false
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
