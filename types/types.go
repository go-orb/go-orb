// Package types provides marker's, these are here against dependency cycles.
package types

import "strings"

// ServiceName is the name of the Service.
type ServiceName string

// SplitServiceName splits the serviceName into a string slice.
func SplitServiceName(serviceName ServiceName) []string {
	return strings.Split(string(serviceName), ".")
}
