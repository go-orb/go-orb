// Package types provides marker's, these are here against dependency cycles.
//
// If this marker's would live in config for example, everything would import from config
// which means config isn't allowed to import a logger from log for example.
package types

import (
	"strings"

	"github.com/go-orb/go-orb/config/source"
)

//nolint:gochecknoglobals
var (
	// DefaultSeperator is used to split a service name into config section keys.
	DefaultSeperator = "."
)

// ServiceName is the name of the Service.
type ServiceName string

// ServiceVersion is the version of the Service.
type ServiceVersion string

// SplitServiceName splits the serviceName into a string slice, separated by
// the global DefaultSeperator. Each item will be used as a key in the config.
//
// Example:
//
//	ServiceName: "com.example.service"
//	Config:
//	```yaml
//	com:
//	  example:
//	    service:
//	      ...
//	```
func SplitServiceName[T ~string](serviceName T) []string {
	return strings.Split(string(serviceName), DefaultSeperator)
}

// JoinServiceName joins a splitted servicename back together
func JoinServiceName(sections []string) ServiceName {
	return ServiceName(strings.Join(sections, DefaultSeperator))
}

// ConfigData holds a single config file marshaled to map[string]any,
// this needs to be done to marshal data back into a components config struct.
//
// After a config source (e.g. a yaml file, or remote resource) has been parsed,
// it will be passed around inside this data type. Each component then gets a
// list of data sources, which layer by layer get applied to eventually construct
// your final component config.
type ConfigData []source.Data
