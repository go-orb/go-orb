package container

// Map is a map container for function factories.
type Map[K Hashable, T any] struct {
	elements map[K]T
}

// NewMap creates a new map of any type.
// This is not concurrent safe.
func NewMap[K Hashable, T any]() *Map[K, T] {
	return &Map[K, T]{
		elements: make(map[K]T),
	}
}

// Add adds a new factory function to this container.
// It returns false if the element already exists.
func (c *Map[K, T]) Add(name K, element T) bool {
	if _, nok := c.elements[name]; nok {
		return false
	}

	c.elements[name] = element

	return true
}

// Set will set a value regardless of whether it already exists in the map.
func (c *Map[K, T]) Set(name K, element T) {
	c.elements[name] = element
}

// All returns the internal map.
func (c *Map[K, T]) All() map[K]T {
	return c.elements
}

// Get returns a single item by its name.
func (c *Map[K, T]) Get(name K) (T, bool) {
	p, ok := c.elements[name]
	return p, ok
}

// Len returns the length of the internal map.
func (c *Map[K, T]) Len() int {
	return len(c.elements)
}

// Del calls delete(me, key) and returns true if the element existed.
func (c *Map[K, T]) Del(key K) bool {
	_, ok := c.Get(key)
	delete(c.elements, key)

	return ok
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (c *Map[K, T]) Range(f func(k K, v T) bool) {
	for k, v := range c.elements {
		if c := f(k, v); !c {
			break
		}
	}
}
