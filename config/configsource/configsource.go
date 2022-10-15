// Package configsource is a base for all config sources.
package configsource

import (
	"fmt"
	"net/url"

	"jochum.dev/orb/orb/util/marshaler"
)

type Source interface {
	fmt.Stringer

	Init() error

	Read(u url.URL) (map[string]any, marshaler.Marshaler, error)
	Write(u url.URL, data map[string]any) error
}
