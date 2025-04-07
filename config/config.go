// Package config provides config handling for go-micro.
package config

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"dario.cat/mergo"
	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config/source"
)

func isNumeric(s string) bool {
	for _, v := range s {
		if v < '0' || v > '9' {
			return false
		}
	}

	return true
}

// WalkMap walks into the sections and returns the map[string]any of that section.
//
// If a section is a slice, the next section should be a number.
//
// Example:
//
//	config := map[string]any{
//	    "foo": map[string]any{
//	        "bar": map[string]any{
//	            "baz": "value",
//	        },
//	    },
//	}
//
// WalkMap([]string{"foo", "bar"}, config) returns map[string]any{"baz": "value"}.
func WalkMap(sections []string, in map[string]any) (map[string]any, error) {
	data := in

	for i := 0; i < len(sections); i++ {
		section := sections[i]

		if i+1 < len(sections) && isNumeric(sections[i+1]) {
			snum, err := strconv.ParseInt(sections[i+1], 10, 64)
			if err != nil {
				return data, fmt.Errorf("while parsing the section number: %w", err)
			}

			sliceData, err := SingleGet(data, section, []any{})
			if err != nil {
				return data, fmt.Errorf("while walking sections '%s': %w", sections, err)
			}

			if int64(len(sliceData)) <= snum {
				return data, fmt.Errorf("while walking sections '%s': %w", sections, ErrNoSuchKey)
			}

			tmpData, ok := sliceData[snum].(map[string]any)
			if !ok {
				return data, fmt.Errorf("while walking sections '%s': %w", sections, ErrNoSuchKey)
			}

			data = tmpData
			i++

			continue
		}

		var err error
		if data, err = SingleGet(data, section, map[string]any{}); err != nil {
			return data, fmt.Errorf("while walking sections '%s': %w", sections, err)
		}
	}

	return data, nil
}

// Read reads url into map[string]any.
func Read(url *url.URL) (map[string]any, error) {
	configSource, err := getSourceForURL(url)
	if err != nil {
		return nil, err
	}

	result, err := configSource.Read(url)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Parse parses the config from config.Read into the given struct.
// Param target should be a pointer to the config to parse into.
func Parse[TMap any](sections []string, key string, config map[string]any, target TMap) error {
	if config == nil {
		return nil
	}

	var (
		data map[string]any
		err  error
	)

	if len(sections) > 0 {
		data, err = WalkMap(sections, config)
		if err != nil {
			return err
		}
	} else {
		data = config
	}

	if key != "" {
		data, err = SingleGet(data, key, map[string]any{})
		if err != nil {
			return err
		}
	}

	codec, err := codecs.GetMime(codecs.MimeJSON)
	if err != nil {
		return err
	}

	b, err := codec.Marshal(data)
	if err != nil {
		return err
	}

	if err := codec.Unmarshal(b, target); err != nil {
		return err
	}

	return nil
}

// ParseSlice parses the config from config.Read into the given slice.
// Param target should be a pointer to the slice to parse into.
func ParseSlice[TSlice any](sections []string, key string, config map[string]any, target TSlice) error {
	if config == nil {
		return nil
	}

	var (
		data map[string]any
		err  error
	)

	if len(sections) > 0 {
		data, err = WalkMap(sections, config)
		if err != nil {
			return err
		}
	} else {
		data = config
	}

	sliceData, err := SingleGet(data, key, []any{})
	if err != nil {
		return err
	}

	if sliceData == nil {
		return nil
	}

	codec, err := codecs.GetMime(codecs.MimeJSON)
	if err != nil {
		return err
	}

	b, err := codec.Marshal(sliceData)
	if err != nil {
		return err
	}

	if err := codec.Unmarshal(b, target); err != nil {
		return err
	}

	return nil
}

// Merge merges the given source into the destination.
func Merge[T any](dst *T, src T) error {
	return mergo.Merge(dst, src, mergo.WithOverride)
}

// Dump is a helper function to dump config to []byte.
func Dump(codecMime string, config map[string]any) ([]byte, error) {
	codec, err := codecs.GetMime(codecMime)
	if err != nil {
		return nil, err
	}

	return codec.Marshal(config)
}

// HasKey returns a boolean which indidcates if the given sections and key exists in the configs.
func HasKey[T any](sections []string, key string, config map[string]any) bool {
	var tmp T

	var err error

	// Walk into the sections.
	data, err := WalkMap(sections, config)
	if err != nil {
		return false
	}

	if _, err := SingleGet(data, key, tmp); err != nil {
		return false
	}

	return true
}

// ParseStruct is a helper to make any struct with `json` tags a map[string]any with sections.
func ParseStruct[TParse any](sections []string, toParse TParse) (map[string]any, error) {
	codec, err := codecs.GetMime(codecs.MimeJSON)
	if err != nil {
		return nil, err
	}

	data := map[string]any{}
	for _, s := range sections {
		if tmp, ok := data[s]; ok {
			switch t2 := tmp.(type) {
			case map[string]any:
				data = t2
			default:
				// Should never happen.
				data = map[string]any{}
			}
		} else {
			tmp := map[string]any{}
			data[s] = tmp
			data = tmp
		}
	}

	buf := bytes.Buffer{}
	if err := codec.NewEncoder(&buf).Encode(toParse); err != nil {
		return nil, fmt.Errorf("encoding: %w", err)
	}

	if err := codec.NewDecoder(&buf).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return data, nil
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
