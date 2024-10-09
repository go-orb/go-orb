package container

// Plugins is an alias around the Map type, for nicer plugin method names.
type Plugins[T any] struct {
	*Map[string, T]
}

// NewPlugins creates a new plugins container of any type.
// This is not concurrent safe.
func NewPlugins[T any]() *Plugins[T] {
	return &Plugins[T]{
		NewMap[string, T](),
	}
}

// Register a plugin.
func (p *Plugins[T]) Register(name string, element T) bool {
	return p.Add(name, element)
}

// Deregister a plugin.
func (p *Plugins[T]) Deregister(name string) bool {
	return p.Del(name)
}
