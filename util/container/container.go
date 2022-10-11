package container

import (
	"fmt"
)

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
		return fmt.Errorf("element '%s' does already exists", name)
	}

	c.elements[name] = element
	return nil
}

func (c *Container[T]) Get(name string) (T, error) {
	p, ok := c.elements[name]
	if !ok {
		var result T
		return result, fmt.Errorf("unknown element '%s' given", name)
	}

	return p, nil
}
