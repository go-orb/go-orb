package container

import (
	"iter"
	"sort"
)

// PriorityListElement needs to be implemented by every interface.
type PriorityListElement[T any] struct {
	Item     T
	Priority int
}

// NewPriorityList creates a new priority list of elements.
func NewPriorityList[T any]() *PriorityList[T] {
	return &PriorityList[T]{
		elements: []PriorityListElement[T]{},
	}
}

// PriorityList is a priority list of elements.
type PriorityList[T any] struct {
	elements []PriorityListElement[T]
}

// Add adds a new factory function to this container.
// It returns ErrExists if the plugin already exists.
func (c *PriorityList[T]) Add(element T, priority int) error {
	item := PriorityListElement[T]{
		Item:     element,
		Priority: priority,
	}

	c.elements = append(c.elements, item)

	return nil
}

// Iterate creates a new iterator over the elements of the priority list.
// Requires Go 1.23 or later.
func (c *PriorityList[T]) Iterate(reversed bool) iter.Seq2[int, T] {
	if reversed {
		sort.SliceStable(c.elements, func(p, q int) bool {
			return c.elements[p].Priority > c.elements[q].Priority
		})
	} else {
		sort.SliceStable(c.elements, func(p, q int) bool {
			return c.elements[p].Priority < c.elements[q].Priority
		})
	}

	return func(yield func(int, T) bool) {
		for i := 0; i <= len(c.elements)-1; i++ {
			if !yield(c.elements[i].Priority, c.elements[i].Item) {
				return
			}
		}
	}
}
