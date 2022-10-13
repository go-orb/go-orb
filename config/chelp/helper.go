package chelp

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownConfig = errors.New("unknown config given")

	ErrNotExistant    = errors.New("no such config key")
	ErrTypesDontMatch = errors.New("config key requested type and actual type don't match")
)

/**
 * Get returns either the value of "key" in "data" or the default value "def".
 * If types don't match it returns ErrTypesDontMatch.
 * If key hasn't been found it returns ErrNotExistant as well as the default value "def".
 *
 * It supports the following datatypes:
 * - any non-container (string/float64/uvm.)
 * - []string slice
 * - []any slice
 * - map[string]string
 * - map[string]any
 */
func Get[T any](data map[string]any, key string, def T) (T, error) {
	v, ok := data[key]
	if !ok {
		return def, ErrNotExistant
	}

	var tmp T
	switch any(tmp).(type) {
	case []string:
		switch vt := any(v).(type) {
		case []any:
			var res = []string{}
			for _, v := range vt {
				res = append(res, fmt.Sprintf("%v", v))
			}
			return any(res).(T), nil
		default:
			return tmp, ErrTypesDontMatch
		}
	case []any:
		switch any(v).(type) {
		case []any:
			return any(v).(T), nil
		default:
			return tmp, ErrTypesDontMatch
		}
	case map[string]string:
		switch vt := v.(type) {
		case map[string]any:
			var res = map[string]string{}
			for k, v := range vt {
				res[k] = fmt.Sprintf("%v", v)
			}
			return any(res).(T), nil
		default:
			return tmp, ErrTypesDontMatch
		}
	case map[string]any:
		switch any(v).(type) {
		case map[string]any:
			return any(v).(T), nil
		default:
			return tmp, ErrTypesDontMatch
		}
	default:
		switch vt := v.(type) {
		case T:
			return vt, nil
		default:
			return tmp, ErrTypesDontMatch
		}
	}
}
