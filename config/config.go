// Package config provides config handling for orb.
package config

import (
	"errors"
	"net/url"

	"jochum.dev/orb/orb/config/configsource"
	"jochum.dev/orb/orb/util/marshaler"
)

var (
	ErrUnknownScheme = errors.New("unknown config source scheme")
)

type Data struct {
	URL       url.URL
	Data      map[string]any
	Marshaler marshaler.Marshaler
	Error     error
}

func Read(urls []url.URL) []Data {
	result := make([]Data, len(urls))

	configsources := []configsource.Source{}
	for _, csFunc := range configsource.Plugins.All() {
		configsources = append(configsources, csFunc())
	}

	for idx, u := range urls {
		found := false
		for _, cs := range configsources {
			if u.Scheme == cs.String() {
				d, m, err := cs.Read(u)
				result[idx] = Data{URL: u, Data: d, Marshaler: m, Error: err}
				found = true
				break
			}
		}

		if !found {
			result[idx] = Data{URL: u, Error: ErrUnknownScheme}
		}
	}

	return result
}
