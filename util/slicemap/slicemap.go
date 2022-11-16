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
