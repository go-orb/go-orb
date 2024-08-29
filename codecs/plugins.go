package codecs

import (
	"fmt"

	"github.com/go-orb/go-orb/util/container"
)

// Plugins is the registry for codec plugins.
var Plugins = container.NewPlugins[Marshaler]() //nolint:gochecknoglobals

// Register makes a plugin available by the provided name.
// If Register is called twice with the same name, it panics.
func Register(name string, codec Marshaler) bool {
	Plugins.Register(name, codec)
	updateMimeMap()

	return true
}

var mimeMap = container.NewSafeMap[string, []Marshaler]() //nolint:gochecknoglobals

func updateMimeMap() {
	Plugins.Range(func(_ string, encoder Marshaler) bool {
		for _, mime := range encoder.ContentTypes() {
			v, ok := mimeMap.Get(mime)
			if !ok {
				// Mime-type is unknown, just add the encoder.
				mimeMap.Set(mime, []Marshaler{encoder})
			} else {
				// We already know that mime-type, see if we know the encoder.
				found := false

				for _, ve := range v {
					if encoder.String() == ve.String() {
						found = true
						break
					}
				}

				if found {
					// Already know that encoder for the given mime-type.
					continue
				}

				// Write back with the new encoder appended.
				v = append(v, encoder)
				mimeMap.Set(mime, v)
			}
		}

		return true
	})
}

// filterMime trims before a space or a semicolon.
func filterMime(mime string) string {
	for i, char := range mime {
		if char == ' ' || char == ';' {
			return mime[:i]
		}
	}

	return mime
}

// GetMime returns a codec for a mime type.
func GetMime(mime string) (Marshaler, error) {
	codec, ok := mimeMap.Get(filterMime(mime))
	if !ok {
		return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownMimeType, mime)
	}

	return codec[0], nil
}

// GetEncoder returns a encoder codec for a mime and golang type.
func GetEncoder(mime string, v any) (Marshaler, error) {
	codecs, ok := mimeMap.Get(filterMime(mime))
	if !ok {
		return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownMimeType, mime)
	}

	for _, codec := range codecs {
		if codec.Encodes(v) {
			return codec, nil
		}
	}

	return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownMimeType, mime)
}

// GetDecoder returns a decoder codec for a mime and golang type.
func GetDecoder(mime string, v any) (Marshaler, error) {
	codecs, ok := mimeMap.Get(filterMime(mime))
	if !ok {
		return nil, fmt.Errorf("%w for %s", ErrUnknownMimeType, mime)
	}

	for _, codec := range codecs {
		if codec.Decodes(v) {
			return codec, nil
		}
	}

	return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownMimeType, mime)
}
