// Package config provides config handling for orb.
package config

import (
	"bytes"
	"errors"
	"net/url"

	"go-micro.dev/v5/config/source"
)

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

// Read reads urls into []Data where Data is basically map[string]any.
//
// By default it will error out if any of these config URL's fail, but you can
// ignore errors for a single url by adding "?ignore_error=true".
//
// prependSections is for url's that don't support sections (cli for example),
// theier result will be prepended, also you can add sections to a single url
// with "?add_section=true".
func Read(urls []*url.URL, prependSections []string) ([]source.Data, error) {
	result := []source.Data{}

	for _, myURL := range urls {
		// Get the config source for the given URL.
		configSource, err := getSourceForURL(myURL)
		if err != nil {
			result = append(result, source.Data{URL: myURL, Error: err})
		}

		// Call the actual read.
		dResult := configSource.Read(myURL)
		if dResult.Error != nil {
			result = append(result, dResult)

			if myURL.Query().Get("ignore_error") == "true" {
				continue
			} else {
				return result, dResult.Error
			}
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

		// Finally append the subresult.
		result = append(result, dResult)
	}

	return result, nil
}

// Parse parses the config from config.Read into the given struct.
func Parse(sections []string, configs []source.Data, target any) error {
	for _, configData := range configs {
		if configData.Error != nil {
			continue
		}

		// Walk the sections
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

		// Create a marshaler for this section.
		buf := bytes.Buffer{}

		// Encode this section into the buf.
		if err := configData.Marshaler.NewEncoder(&buf).Encode(data); err != nil {
			return err
		}

		// Decode this section from the buf.
		if err := configData.Marshaler.NewDecoder(&buf).Decode(target); err != nil {
			return err
		}
	}

	return nil
}
