package codecs

import (
	"bytes"
	"encoding/json"
	"io"
)

var _ Marshaler = (*CodecJSON)(nil)

// init() {}.
var _ = Register("json", &CodecJSON{})

// CodecJSON implements the codecs.Marshal interface, and can be used for marshaling
// CodecJSON config files, and web requests.
type CodecJSON struct{}

// Encode marshals any object into json bytes.
// Param v should be a pointer type.
func (j *CodecJSON) Encode(v any) ([]byte, error) {
	b := []byte{}
	buf := bytes.NewBuffer(b)
	err := j.NewEncoder(buf).Encode(v)

	return buf.Bytes(), err
}

// Decode decodes json bytes into object v.
// Param v should be a pointer type.
func (j *CodecJSON) Decode(b []byte, v any) error {
	buf := bytes.NewBuffer(b)

	return j.NewDecoder(buf).Decode(v)
}

type wrapEncoder struct {
	impl *json.Encoder
}

func (j *wrapEncoder) Encode(v any) error {
	return j.impl.Encode(v)
}

// NewEncoder returns a new JSON encoder.
func (j *CodecJSON) NewEncoder(w io.Writer) Encoder {
	return &wrapEncoder{impl: json.NewEncoder(w)}
}

type wrapDecoder struct {
	impl *json.Decoder
}

func (j *wrapDecoder) Decode(v any) error {
	return j.impl.Decode(v)
}

// NewDecoder returns a new JSON decoder.
func (j *CodecJSON) NewDecoder(r io.Reader) Decoder {
	return &wrapDecoder{impl: json.NewDecoder(r)}
}

// Encodes returns if this codec is able to encode the given type.
func (j *CodecJSON) Encodes(v any) bool {
	switch v.(type) {
	case []string:
		return true
	case []byte:
		return true
	case []any:
		return true
	case map[string]any:
		return true
	case string:
		return true
	default:
		return false
	}
}

// Decodes returns if this codec is able to decode the given type.
func (j *CodecJSON) Decodes(v any) bool {
	return j.Encodes(v)
}

// ContentTypes returns the content types the marshaler can handle.
func (j *CodecJSON) ContentTypes() []string {
	return []string{
		"application/json",
	}
}

// String returns the plugin implementation of the marshaler.
func (j *CodecJSON) String() string {
	return "json"
}

// Exts is a list of file extensions this marshaler supports.
func (j *CodecJSON) Exts() []string {
	return []string{".json"}
}
