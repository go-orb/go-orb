// Package configsource is a base for all config sources.
package configsource

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrUnknownScheme = errors.New("unknown config source scheme")
)

type ConfigSource interface {
	fmt.Stringer

	Init() error

	Read(u url.URL) (map[string]any, error)
	Write(u url.URL, data map[string]any) error
}
