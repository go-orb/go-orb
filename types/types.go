// Package types provides marker's, these are here against dependency cycles.
package types

import (
	"strings"

	"go-micro.dev/v5/config/source"
)

// ServiceName is the name of the Service.
type ServiceName string

// SplitServiceName splits the serviceName into a string slice.
func SplitServiceName(serviceName ServiceName) []string {
	return strings.Split(string(serviceName), ".")
}

// ConfigData holds a single config file marshaled to map[string]any,
// this needs to be done to marshal data back into a components config struct.
//
// After a config source (e.g. a yaml file, or remote resource) has been parsed,
// it will be passed around inside this data type. Each component then gets a
// list of data sources, which layer by layer get applied to eventually construct
// your final component config.
type ConfigData []source.Data
