package codecs

import (
	"encoding/json"
	"io"
)

var _ Marshaler = (*CodecJSON)(nil)

// init() {}.
var _ = Register("json", &CodecJSON{})

// CodecJSON implements the codecs.Marshal interface, and can be used for marshaling
// CodecJSON config files, and web requests.
type CodecJSON struct{}

// Marshal marshals any object into json bytes.
// Param v should be a pointer type.
func (j *CodecJSON) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes json bytes into object v.
// Param v should be a pointer type.
func (j *CodecJSON) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// NewEncoder returns a new JSON encoder.
func (j *CodecJSON) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

// NewDecoder returns a new JSON decoder.
func (j *CodecJSON) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
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
