// Package source is a base for all config sources.
package source

import (
	"fmt"
	"net/url"

	"go-micro.dev/v5/codecs"
)

// Data holds a single config file marshaled to map[string]any.
type Data struct {
	// Source URL.
	URL *url.URL
	// Data holder.
	Data map[string]any
	// The Marshaler used to create Data -> map[string]any
	Marshaler codecs.Marshaler
	// If there was an error while processing the URL.
	Error error

	// AdditionalConfigs is a list of configs that we also have to read.
	// or that have been injected by config.Read().
	AdditionalConfigs []*url.URL
}

// Source is a config source.
type Source interface {
	fmt.Stringer

	// Schemes is a slice of schemes this reader supports.
	Schemes() []string

	// PrependSections indicates whether config.Read() has to prepend the result
	// with sections.
	PrependSections() bool

	// Read reads the url in u and returns it as map[string]any.
	Read(u *url.URL) Data
}
