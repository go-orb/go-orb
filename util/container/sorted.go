package container

import (
	"fmt"
	"sort"
)

// SortedElement needs to be implemented by every interface.
type SortedElement interface {
	fmt.Stringer

	Priority() int
}

// NewSorted creates a new container that holds elements of type T sorted.
// Adds are costly as it sorts on each add, calls to Sorted() are free.
func NewSorted[T SortedElement]() *Sorted[T] {
	return &Sorted[T]{
		elements: []T{},
	}
}

// Sorted is a sorted container.
type Sorted[T SortedElement] struct {
	elements []T
}

// Add adds a new factory function to this container.
// It returns ErrExists if the plugin already exists.
func (c *Sorted[T]) Add(element T) error {
	for _, e := range c.elements {
		if e.String() == element.String() {
			return ErrExists
		}
	}

	c.elements = append(c.elements, element)

	sort.SliceStable(c.elements, func(p, q int) bool {
		return c.elements[p].Priority() < c.elements[q].Priority()
	})

	return nil
}

// Sorted returns the internal list.
func (c *Sorted[T]) Sorted() []T {
	return c.elements
}
