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

	for i := 0; i < len(sections); i++ {
		section := sections[i]

		if i+1 < len(sections) && isAlphaNumeric(sections[i+1]) {
			snum, err := strconv.ParseInt(sections[i+1], 10, 64)
			if err != nil {
				return data, fmt.Errorf("while parsing the section number: %w", err)
			}

			sliceData, err := SingleGet(data, section, []any{})
			if err != nil {
				return data, err
			}

			if int64(len(sliceData)) <= snum {
				return data, ErrNotExistent
			}

			tmpData, ok := sliceData[snum].(map[string]any)
			if !ok {
				return data, ErrNotExistent
			}

			data = tmpData
			i++

			continue
		}

		var err error
		if data, err = SingleGet(data, section, map[string]any{}); err != nil {
			return data, err
		}
	}

	return data, nil
}

// Read reads urls into []Data where Data is map[string]any.
//
// By default it will error out if any of these config URLs fail, but you can
// ignore errors for a single url by adding "?ignore_error=true".
func Read(urls []*url.URL) (types.ConfigData, error) {
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
func HasKey[T any](sections []string, key string, configs types.ConfigData) bool {
	var tmp T

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

		if _, err := SingleGet(data, key, tmp); err != nil {
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

// Dump is a helper function to dump configDatas to the console.
func Dump(configs types.ConfigData) error {
	codec, err := codecs.GetMime(codecs.MimeJSON)
	if err != nil {
		return err
	}

	for _, config := range configs {
		jsonb, err := codec.Marshal(config.Data)
		if err == nil {
			fmt.Println(string(jsonb)) //nolint:forbidigo
		}
	}

	return nil
}

// ParseStruct is a helper to make any struct with `json` tags a source.Data (map[string]any{} with some more fields) with sections.
func ParseStruct[TParse any](sections []string, toParse TParse) (source.Data, error) {
	result := source.Data{Data: make(map[string]any)}

	codec, err := codecs.GetMime(codecs.MimeJSON)
	if err != nil {
		result.Error = err
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
