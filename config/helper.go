package config

import (
	"fmt"
)

// Get returns either the value of "key" in "data" or the default value "def".
// If types don't match it returns ErrTypesDontMatch.
// If key hasn't been found it returns ErrNotExistent as well as the default value "def".
//
// It supports the following datatypes:
//   - any non-container (string/float64/uvm.)
//   - []string slice
//   - []any slice
//   - []map[string]any
func Get[T any](data map[string]any, key string, def T) (T, error) {
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

			return any(res).(T), nil
		default:
			return def, ErrTypesDontMatch
		}
	case []any:
		switch value.(type) {
		case []any:
			return value.(T), nil
		default:
			return def, ErrTypesDontMatch
		}
	case map[string]any:
		switch value.(type) {
		case map[string]any:
			return value.(T), nil
		default:
			return def, ErrTypesDontMatch
		}
	default:
		switch vt := value.(type) {
		case T:
			return vt, nil
		default:
			return def, ErrTypesDontMatch
		}
	}
}
