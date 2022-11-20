// Package slicemap provides simple slice & map utility functions.
package slicemap

import (
	"errors"

	"golang.org/x/exp/constraints"
)

// In check if query is an element of the list.
func In[T constraints.Ordered](s []T, query T) bool {
	for _, element := range s {
		if element == query {
			return true
		}
	}

	return false
}

// Lookup takes a map and will return a nested key.
func Lookup(data map[string]any, path []string) (any, error) {
	for i, key := range path {
		nd, ok := data[key]
		if !ok {
			return nil, errors.New("key not found")
		}

		if len(path)-1 == i {
			return nd, nil
		}

		data, ok = nd.(map[string]any)
		if !ok {
			return nil, errors.New("key not found")
		}
	}

	return nil, errors.New("key not found")
}

// SetValue takes a map and will assign a value at the specified path.
func SetValue(data map[string]any, path []string, value any) {
	for i, key := range path {
		if len(path)-1 == i {
			data[key] = value
			return
		}

		// If key empty assign empty map.
		temp, ok := data[key]
		if !ok {
			data[key] = map[string]any{}
		}

		// Overwrite value if of a different type.
		if _, ok := temp.(map[string]any); !ok {
			data[key] = map[string]any{}
		}

		data = data[key].(map[string]any) //nolint:errcheck
	}
}

// CopyMap makes a deep copy of a map.
func CopyMap(m map[string]any) map[string]any {
	cp := make(map[string]any, len(m))

	for k, v := range m {
		vm, ok := v.(map[string]any)
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

// Get returns a value from a map in a specified type. If either item is not
// present or of a different type, ok = false.
func Get[T any](m map[string]any, field string) (res T, ok bool) {
	value, ok := m[field]
	if !ok {
		return res, false
	}

	res, ok = value.(T)

	return res, ok
}
