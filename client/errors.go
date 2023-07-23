package client

import "errors"

var (
	// ErrNoNodeFound happens we haven't found a node for the requested service.
	ErrNoNodeFound = errors.New("no node found for the requested service")
	// ErrServiceArgumentEmpty happens when the service argument is empty.
	ErrServiceArgumentEmpty = errors.New("service argument is empty")
	// ErrFailedToCreateTransport happens when we haven't found the transport requested.
	ErrFailedToCreateTransport = errors.New("failed to create a transport")
	// ErrUnknownContentType happens when you request something with a (yet) unknown content-type.
	ErrUnknownContentType = errors.New("unknown content-type has been requested")
)
