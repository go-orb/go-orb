package container

import (
	"github.com/cornelk/hashmap"
)

// NewSafeMap creates a wrapper for "github.com/cornelk/hashmap".
func NewSafeMap[K Hashable, T any]() *SafeMap[K, T] {
	return &SafeMap[K, T]{
		Map: hashmap.New[K, T](),
	}
}

// SafeMap is a wrapper for "github.com/cornelk/hashmap".
type SafeMap[K Hashable, T any] struct {
	*hashmap.Map[K, T]
}

// Add is wrapper for hasmap.Insert.
func (m *SafeMap[K, T]) Add(key K, value T) bool {
	return m.Insert(key, value)
}
