package marshaler

import (
	"errors"
	"io"
)

var (
	ErrNoSocket = errors.New("no socket given")
)

// Marshaler is not goroutines save.
type Marshaler interface {
	// Init sets the sockets of Marshaler.
	Init(r io.Reader, w io.Writer) error

	// EncodeSocket writes msg to the socket.
	EncodeSocket(msg any) error

	// DecodeSocket reads msg from the socket.
	DecodeSocket(msg any) error
}
