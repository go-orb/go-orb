package container

import (
	"errors"
)

var ErrExists = errors.New("element exists already")
var ErrUnknown = errors.New("unknown element given")

func New[T any](cType T) *Container[T] {
	return &Container[T]{
		elements: make(map[string]T),
	}
}

type Container[T any] struct {
	elements map[string]T
}

func (c *Container[T]) Add(name string, element T) error {
	if _, nok := c.elements[name]; nok {
		return ErrExists
	}

	c.elements[name] = element
	return nil
}

func (c *Container[T]) Get(name string) (T, error) {
	p, ok := c.elements[name]
	if !ok {
		var result T
		return result, ErrUnknown
	}

	return p, nil
}
