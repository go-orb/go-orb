package container

import (
	"iter"
	"slices"
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

// Add adds a new element with the given priority to the list.
func (c *PriorityList[T]) Add(element T, priority int) error {
	item := PriorityListElement[T]{
		Item:     element,
		Priority: priority,
	}

	c.elements = append(c.elements, item)

	return nil
}

// Iterate clones the internal list, sorts it and then returns a new iterator
// over the elements of the priority list.
// Requires Go 1.23 or later.
func (c *PriorityList[T]) Iterate(reversed bool) iter.Seq2[int, T] {
	elements := slices.Clone(c.elements)

	sort.SliceStable(elements, func(p, q int) bool {
		return elements[p].Priority < elements[q].Priority
	})

	return func(yield func(int, T) bool) {
		if reversed {
			for i := len(elements) - 1; i >= 0; i-- {
				if !yield(elements[i].Priority, elements[i].Item) {
					return
				}
			}
		} else {
			for i := 0; i <= len(elements)-1; i++ {
				if !yield(elements[i].Priority, elements[i].Item) {
					return
				}
			}
		}
	}
}
