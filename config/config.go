// Package config provides config handling for orb.
package config

import (
	"github.com/hashicorp/go-multierror"
	"jochum.dev/orb/orb/config/configsource"
)

type Data struct {
	URL  string
	Data map[string]any
}

func Read(urls []string) ([]Data, error) {
	var (
		result    []Data
		resultErr error
	)

	for _, url := range urls {
		d, err := readURL(url)
		result = append(result, d)
		resultErr = multierror.Append(resultErr, err)
	}

	return result, resultErr
}

func readURL(url string) (Data, error) {
	return Data{}, configsource.ErrUnknown
}
