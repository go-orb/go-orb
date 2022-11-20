package container

// Plugins is an alias around the Map type, for nicer plugin method names.
type Plugins[T any] struct {
	Map[T]
}

// NewPlugins creates a new plugins container of any type.
// Not concurency safe.
func NewPlugins[T any]() *Plugins[T] {
	return &Plugins[T]{
		Map: Map[T]{
			elements: make(map[string]T),
		},
	}
}

// Register a plugin.
func (p *Plugins[T]) Register(name string, element T) {
	p.Set(name, element)
}

// Deregister a plugin.
func (p *Plugins[T]) Deregister(name string) {
	delete(p.elements, name)
}
