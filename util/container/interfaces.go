package container

// Hashable is used as the key for maps.
type Hashable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

// Priorizeable has to be implemented by elements that use container.Priorize*.
type Priorizeable interface {
	Priority() int
}
