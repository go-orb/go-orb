package container

import "errors"

var (
	// ErrExists is returned when the given element exists in the container.
	ErrExists = errors.New("element exists already")
	// ErrUnknown is returned when the given name doesn't exist in the container.
	ErrUnknown = errors.New("unknown element given")
)
