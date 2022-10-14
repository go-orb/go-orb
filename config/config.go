// Package config provides config handling for orb.
package config

import (
	"errors"
	"net/url"

	"jochum.dev/orb/orb/config/configsource"
)

var (
	ErrUnknownScheme = errors.New("unknown config source scheme")
)

type Data struct {
	URL   url.URL
	Data  map[string]any
	Error error
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
				d, err := cs.Read(u)
				result[idx] = Data{URL: u, Data: d, Error: err}
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
