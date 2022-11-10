// Package container contains generic containers.
package container

import (
	"fmt"
)

// NewList creates a new container that holds element's of type T.
func NewList[T fmt.Stringer]() *List[T] {
	return &List[T]{
		elements: []T{},
	}
}

// List is a list container of T.
type List[T fmt.Stringer] struct {
	elements []T
}

// Add adds a new element to the container.
func (c *List[T]) Add(element T) error {
	for _, e := range c.elements {
		if e.String() == element.String() {
			return ErrExists
		}
	}

	c.elements = append(c.elements, element)

	return nil
}

// List returns the internal list.
func (c *List[T]) List() []T {
	return c.elements
}

// Clear clears the internal list.
func (c *List[T]) Clear() {
	c.elements = []T{}
}
