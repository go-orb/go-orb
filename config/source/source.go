// Package source provides a base for all config sources.
// It provides a source interface which can be used to create config sources,
// and a data type, which gets used to pass around parsed config sources.
package source

import (
	"net/url"

	"go-micro.dev/v5/codecs"
)

// Data holds a single config file marshaled to map[string]any,
// this needs to be done to marshal data back into a components config struct.
//
// After a config source (e.g. a yaml file, or remote resource) has been parsed,
// it will be passed around inside this data type. Each component then gets a
// list of data sources, which layer by layer get applied to eventually construct
// your final component config.
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
	// Schemes is a slice of schemes this reader supports.
	Schemes() []string

	// PrependSections indicates whether config.Read() has to prepend the result
	// with sections.
	PrependSections() bool

	// Read reads the url in u and returns it as map[string]any.
	Read(u *url.URL) Data

	// String returns the name of the source.
	String() string
}
