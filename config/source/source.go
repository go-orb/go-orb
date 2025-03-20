// Package source provides a base for all config sources.
// It provides a source interface which can be used to create config sources,
// and a data type, which gets used to pass around parsed config sources.
package source

import (
	"net/url"
)

// Source is a config source.
type Source interface {
	// Schemes is a slice of schemes this reader supports.
	Schemes() []string

	// Read reads the url in u and returns it as map[string]any.
	Read(u *url.URL) (map[string]any, error)

	// String returns the name of the source.
	String() string
}
