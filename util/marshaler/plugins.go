package marshaler

import (
	"github.com/go-orb/orb/util/container"
)

// Plugins is the marshaler plugin container.
var Plugins = container.New(func() Marshaler { return nil }) //nolint:gochecknoglobals

func Marshalers() map[string]Marshaler {
	result := map[string]Marshaler{}
	for _, mFunc := range Plugins.All() {
		m := mFunc()
		result[m.String()] = m
	}

	return result
}
