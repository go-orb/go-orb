// Package slice provides slice utility functions.
package slice

import "golang.org/x/exp/constraints"

// In check if query is an element of the list.
func In[T constraints.Ordered](s []T, query T) bool {
	for _, element := range s {
		if element == query {
			return true
		}
	}

	return false
}
