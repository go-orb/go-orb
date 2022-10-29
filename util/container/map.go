package container

// NewMap creates a new container that holds factory function's.
func NewMap[T any]() *Map[T] {
	return &Map[T]{
		elements: make(map[string]T),
	}
}

// Map is a map container for function factories.
type Map[T any] struct {
	elements map[string]T
}

// Add adds a new factory function to this container.
// It returns ErrExists if the plugin already exists.
func (c *Map[T]) Add(name string, element T) error {
	if _, nok := c.elements[name]; nok {
		return ErrExists
	}

	c.elements[name] = element

	return nil
}

// All returns the internal map.
func (c *Map[T]) All() map[string]T {
	return c.elements
}

// Get returns a single item by its name.
func (c *Map[T]) Get(name string) (T, error) {
	p, ok := c.elements[name]
	if !ok {
		var result T
		return result, ErrUnknown
	}

	return p, nil
}
