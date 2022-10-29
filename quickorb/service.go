package quickorb

import (
	"errors"

	"github.com/orb-org/orb/log"
	"github.com/orb-org/orb/registry"
)

// ErrRequiredOption is returned when an required option haven't been given.
var ErrRequiredOption = errors.New("required option not given")

// Service is an interface that wraps the lower level components
// within orb. Its a convenience method for building and initializing services.
type Service struct {
	options *Options

	logger   log.Logger
	registry registry.Registry
}

// ProvideService provides a service with components.
func ProvideService(
	opts *Option,
	logger log.Logger,
	registry registry.Registry,
) (*Service, error) {
	s := &Service{}

	if s.options.Name == "" {
		return nil, ErrRequiredOption
	}

	if s.options.Version == "" {
		return nil, ErrRequiredOption
	}

	if s.logger == nil {
		return nil, ErrRequiredOption
	}

	s.logger = logger
	s.registry = registry

	return s, nil
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.options.Name
}

// Version returns the version of the service.
func (s *Service) Version() string {
	return s.options.Version
}

// Logger returns the services logger.
func (s *Service) Logger() log.Logger {
	return s.logger
}

// Registry returns the services registry, it may returns nil.
func (s *Service) Registry() registry.Registry {
	return s.registry
}
