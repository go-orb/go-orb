// Package config provides config handling for go-micro.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config/source"
	"github.com/go-orb/go-orb/types"
)

func isAlphaNumeric(s string) bool {
	for _, v := range s {
		if v < '0' || v > '9' {
			return false
		}
	}

	return true
}

func walkMap(sections []string, in map[string]any) (map[string]any, error) {
	data := in

	for _, section := range sections {
		if isAlphaNumeric(section) {
			snum, err := strconv.ParseInt(section, 10, 64)
			if err != nil {
				return data, fmt.Errorf("while parsing the section number: %w", err)
			}

			sliceData, err := Get(data, section, []any{})
			if err != nil {
				return data, err
			}

			tmpData, ok := sliceData[snum].(map[string]any)
			if !ok {
				return data, ErrNotExistent
			}

			data = tmpData
		}

		var err error
		if data, err = Get(data, section, map[string]any{}); err != nil {
			return data, err
		}
	}

	return data, nil
}

// Read reads urls into []Data where Data is map[string]any.
//
// By default it will error out if any of these config URLs fail, but you can
// ignore errors for a single url by adding "?ignore_error=true".
//
// prependSections is for url's that don't support sections (cli for example),
// their result will be prepended, also you can add sections to a single url
// with "?add_section=true".
func Read(urls []*url.URL, prependSections []string) (types.ConfigData, error) {
	result := types.ConfigData{}

	for _, myURL := range urls {
		configSource, err := getSourceForURL(myURL)
		if err != nil {
			result = append(result, source.Data{URL: myURL, Error: err})
			return result, err
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
func Parse(sections []string, configs types.ConfigData, target any) error {
	for _, configData := range []source.Data(configs) {
		if configData.Error != nil {
			continue
		}

		var err error

		// Walk into the sections.
		data, err := walkMap(sections, configData.Data)
		if err != nil {
			if errors.Is(err, ErrNotExistent) {
				continue
			}

			return err
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

// HasKey returns a boolean which indidcates if the given sections and key exists in the configs.
func HasKey(sections []string, key string, configs types.ConfigData) bool {
	for _, configData := range []source.Data(configs) {
		if configData.Error != nil {
			continue
		}

		var err error

		// Walk into the sections.
		data, err := walkMap(sections, configData.Data)
		if err != nil {
			if errors.Is(err, ErrNotExistent) {
				continue
			}

			return false
		}

		if _, err := Get(data, key, ""); err != nil {
			// Ignore unknown configSection in config.
			if errors.Is(err, ErrNotExistent) {
				continue
			}

			return false
		}

		return true
	}

	return false
}

// ParseStruct is a helper to make any struct with `json` tags a source.Data (map[string]any{} with some more fields) with sections.
func ParseStruct[TParse any](sections []string, toParse TParse) (source.Data, error) {
	result := source.Data{Data: make(map[string]any)}

	codec, err := codecs.Plugins.Get("json")
	if err != nil {
		result.Error = fmt.Errorf("getting the json codec: %w", err)
		return result, result.Error
	}

	result.Marshaler = codec

	data := result.Data
	for _, s := range sections {
		if tmp, ok := data[s]; ok {
			switch t2 := tmp.(type) {
			case map[string]any:
				data = t2
			default:
				// Should never happen.
				data = result.Data
			}
		} else {
			tmp := map[string]any{}
			data[s] = tmp
			data = tmp
		}
	}

	buf := bytes.Buffer{}
	if err := codec.NewEncoder(&buf).Encode(toParse); err != nil {
		return result, fmt.Errorf("encoding: %w", err)
	}

	if err := codec.NewDecoder(&buf).Decode(&data); err != nil {
		return result, fmt.Errorf("decoding: %w", err)
	}

	return result, nil
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
