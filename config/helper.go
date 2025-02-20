package config

import (
	"fmt"
)

// SingleGet returns either the value of "key" in "data" or the default value "def".
// If types don't match it returns ErrTypesDontMatch.
// If key hasn't been found it returns ErrNotExistent as well as the default value "def".
//
// It supports the following datatypes:
//   - any non-container (string/float64/uvm.)
//   - []string slice
//   - []any slice
//   - []map[string]any
func SingleGet[T any](data map[string]any, key string, def T) (T, error) {
	value, ok := data[key]
	if !ok {
		return def, ErrNotExistent
	}

	switch any(def).(type) {
	case []string:
		switch vt := value.(type) {
		case []any:
			var res = []string{}
			for _, v := range vt {
				res = append(res, fmt.Sprintf("%v", v))
			}

			return any(res).(T), nil //nolint:errcheck
		default:
			return def, fmt.Errorf("%w: []string", ErrTypesDontMatch)
		}
	case []any:
		switch value.(type) {
		case []any:
			return value.(T), nil //nolint:errcheck
		default:
			return def, fmt.Errorf("%w: []any", ErrTypesDontMatch)
		}
	case map[string]any:
		switch value.(type) {
		case map[string]any:
			return value.(T), nil //nolint:errcheck
		default:
			return def, fmt.Errorf("%w: map[string]any", ErrTypesDontMatch)
		}
	default:
		switch vt := value.(type) {
		case T:
			return vt, nil
		default:
			return def, fmt.Errorf("%w: default", ErrTypesDontMatch)
		}
	}
}
