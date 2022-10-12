package container

type pluginsNode[T any, U any] struct {
	plugin T
	config U
}

func NewPlugins[T any, U any](plugin T, config U) *Plugins[T, U] {
	return &Plugins[T, U]{
		elements: make(map[string]pluginsNode[T, U]),
	}
}

type Plugins[T any, U any] struct {
	elements map[string]pluginsNode[T, U]
}

func (c *Plugins[T, U]) Add(name string, plugin T, config U) error {
	if _, nok := c.elements[name]; nok {
		return ErrExists
	}

	c.elements[name] = pluginsNode[T, U]{plugin: plugin, config: config}
	return nil
}

func (c *Plugins[T, U]) Get(name string) (T, U, error) {
	p, ok := c.elements[name]
	if !ok {
		var plugin T
		var config U
		return plugin, config, ErrUnknown
	}

	return p.plugin, p.config, nil
}

func (c *Plugins[T, U]) Plugin(name string) (T, error) {
	plugin, _, err := c.Get(name)
	return plugin, err
}

func (c *Plugins[T, U]) Config(name string) (U, error) {
	_, config, err := c.Get(name)
	return config, err
}
