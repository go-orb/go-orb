// Package marshaler contains a basic interface for marshalers.
package marshaler

import (
	"errors"
	"fmt"
	"io"
)

var (
	// ErrNoSocket happens when you give wrong options to Marshaler.Init().
	ErrNoSocket = errors.New("no socket given")

	ErrNoFileMarshaler = errors.New("no marshaler for the given file found")
)

// Marshaler is not goroutines save.
type Marshaler interface {
	fmt.Stringer

	// Init sets the sockets of Marshaler.
	Init(r io.Reader, w io.Writer) error

	FileExtension() string

	// EncodeSocket writes msg to the socket.
	EncodeSocket(msg any) error

	// DecodeSocket reads msg from the socket.
	DecodeSocket(msg any) error
}
