// Package quickorb is the quick start entry for the orb framework.
package quickorb

import (
	"errors"
	"go-micro.dev/v5/cli"
)

func NewService(opts ...Option) (*Service, error) {
	options := NewOptions(opts...)

	return nil, errors.New("unimplemented")
}
