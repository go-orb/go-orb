// Package container provides generic map types.
package container

import "sync"

// SafeMap is a concurrently safe generic map.
type SafeMap[T any] struct {
	mu sync.RWMutex
	Map[T]
}

// NewSafeMap creates a new concurrently map of any types.
func NewSafeMap[T any]() *SafeMap[T] {
	return &SafeMap[T]{
		Map: Map[T]{},
	}
}

// Add adds a new factory function to this container.
// It returns ErrExists if the plugin already exists.
func (c *SafeMap[T]) Add(name string, element T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Map.Add(name, element)
}

// All returns the internal map.
func (c *SafeMap[T]) All() map[string]T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Map.All()
}

// Get returns a single item by its name.
func (c *SafeMap[T]) Get(name string) (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Map.Get(name)
}
