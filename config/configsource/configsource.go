// Package configsource is a base for all config sources.
package configsource

import (
	"fmt"
	"net/url"
)

type Source interface {
	fmt.Stringer

	Init() error

	Read(u url.URL) (map[string]any, error)
	Write(u url.URL, data map[string]any) error
}
