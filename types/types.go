// Package types provides marker's, these are here against dependency cycles.
//
// If this marker's would live in config for example, everything would import from config
// which means config isn't allowed to import a logger from log for example.
package types

import (
	"strings"
)

//nolint:gochecknoglobals
var (
	// DefaultSeparator is used to split a service name into config section keys.
	DefaultSeparator = "."
)

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
	return strings.Split(string(serviceName), DefaultSeparator)
}

// JoinServiceName joins a splitted servicename back together.
func JoinServiceName(sections []string) string {
	return strings.Join(sections, DefaultSeparator)
}
