// Package config provides config handling for orb.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"

	"go-micro.dev/v5/config/source"
)

// Read reads urls into []Data where Data is map[string]any.
//
// By default it will error out if any of these config URLs fail, but you can
// ignore errors for a single url by adding "?ignore_error=true".
//
// prependSections is for url's that don't support sections (cli for example),
// their result will be prepended, also you can add sections to a single url
// with "?add_section=true".
func Read(urls []*url.URL, prependSections []string) ([]source.Data, error) {
	result := []source.Data{}

	for _, myURL := range urls {
		configSource, err := getSourceForURL(myURL)
		if err != nil {
			result = append(result, source.Data{URL: myURL, Error: err})
		}

		dResult := configSource.Read(myURL)
		if dResult.Error != nil {
			result = append(result, dResult)

			if myURL.Query().Get("ignore_error") == "true" {
				continue
			}

			return result, dResult.Error
		}

		// Read additional Configs from the config Source if any.
		if len(dResult.AdditionalConfigs) > 0 {
			aDatas, aErr := Read(dResult.AdditionalConfigs, prependSections)
			result = append(result, aDatas...)

			if aErr != nil {
				return result, aErr
			}
		}

		// Prepend the result with sections if required.
		if len(prependSections) > 0 && (configSource.PrependSections() ||
			myURL.Query().Get("add_section") == "true") {
			data := map[string]any{}
			prependResult := data

			for _, s := range prependSections[:len(prependSections)-1] {
				tmp := map[string]any{}
				data[s] = tmp
				data = tmp
			}

			data[prependSections[len(prependSections)-1]] = dResult.Data
			dResult.Data = prependResult
		}

		result = append(result, dResult)
	}

	return result, nil
}

// Parse parses the config from config.Read into the given struct.
// Param target should be a pointer to the config to parse into.
func Parse(sections []string, configs []source.Data, target any) error {
	for _, configData := range configs {
		if configData.Error != nil {
			continue
		}

		var err error

		data := configData.Data

		for _, section := range sections {
			if data, err = Get(data, section, map[string]any{}); err != nil {
				// Ignore unknown configSection in config.
				if errors.Is(err, ErrNotExistent) {
					continue
				}

				return err
			}
		}

		buf := bytes.Buffer{}

		// Here we need to take the data from the configs, in the format of map[string]any
		// and parse it into the struct. Because we cannot do this in one operation,
		// we first marshal the map[string]any into a (usually json) byte slice,
		// We then unmarshal the byte slice into the target struct.
		// If there is a way to do this in one operation, this code should be updated.

		if err := configData.Marshaler.NewEncoder(&buf).Encode(data); err != nil {
			return fmt.Errorf("parse config: encode: %w", err)
		}

		if err := configData.Marshaler.NewDecoder(&buf).Decode(target); err != nil {
			return fmt.Errorf("parse config: decode: %w", err)
		}
	}

	return nil
}

func getSourceForURL(u *url.URL) (source.Source, error) {
	for _, cs := range source.Plugins.List() {
		for _, scheme := range cs.Schemes() {
			if u.Scheme == scheme {
				return cs, nil
			}
		}
	}

	return nil, ErrUnknownScheme
}
