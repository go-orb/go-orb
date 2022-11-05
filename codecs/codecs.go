// Package codecs is provides an interface to encode and decode content types
// to and from byte sequences.
package codecs

import "io"

// Marshaler is able to encode/decode a content type to/from a byte sequence.
type Marshaler interface {
	// Marshal marshals "v" into byte sequence.
	Marshal(v any) ([]byte, error)

	// Unmarshal unmarshals "data" into "v".
	// "v" must be a pointer value.
	Unmarshal(data []byte, v any) error

	// NewDecoder returns a Decoder which reads byte sequence from "r".
	NewDecoder(r io.Reader) Decoder

	// NewEncoder returns an Encoder which writes bytes sequence into "w".
	NewEncoder(w io.Writer) Encoder

	// ContentTypes returns the list of content types this marshaller is able to
	// handle.
	ContentTypes() []string

	// String returns the codec name.
	String() string
}

// Decoder decodes a byte sequence.
type Decoder interface {
	Decode(v any) error
}

// Encoder encodes payloads / fields into byte sequence.
type Encoder interface {
	Encode(v any) error
}

// DecoderFunc adapts an decoder function into Decoder.
type DecoderFunc func(v any) error

// Decode delegates invocations to the underlying function itself.
func (f DecoderFunc) Decode(v any) error { return f(v) }

// EncoderFunc adapts an encoder function into Encoder.
type EncoderFunc func(v any) error

// Encode delegates invocations to the underlying function itself.
func (f EncoderFunc) Encode(v any) error { return f(v) }
