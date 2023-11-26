package codecs

import (
	"errors"
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

var mimeMap = container.NewSafeMap[[]Marshaler]() //nolint:gochecknoglobals

func updateMimeMap() {
	for _, codec := range Plugins.All() {
		// One codec can support multiple mime types, we add all of them to the map.
		for _, mime := range codec.ContentTypes() {
			err := mimeMap.Add(mime, []Marshaler{codec})
			if errors.Is(err, container.ErrExists) {
				existing, err := mimeMap.Get(mime)
				if err != nil {
					continue
				}

				existing = append(existing, codec)
				mimeMap.Set(mime, existing)
			}
		}
	}
}

// filterMime trims before a space or semicolon.
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
	codec, err := mimeMap.Get(filterMime(mime))
	if err != nil {
		return nil, fmt.Errorf("%w for %s", ErrUnknownMimeType, mime)
	}

	return codec[0], nil
}

// GetEncoder returns a encoder codec for a mime and golang type.
func GetEncoder(mime string, v any) (Marshaler, error) {
	codecs, err := mimeMap.Get(filterMime(mime))
	if err != nil {
		return nil, fmt.Errorf("%w for %s", ErrUnknownMimeType, mime)
	}

	for _, codec := range codecs {
		if codec.Encodes(v) {
			return codec, nil
		}
	}

	return nil, fmt.Errorf("%w for '%s'", ErrUnknownMimeType, mime)
}

// GetDecoder returns a decoder codec for a mime and golang type.
func GetDecoder(mime string, v any) (Marshaler, error) {
	codecs, err := mimeMap.Get(filterMime(mime))
	if err != nil {
		return nil, fmt.Errorf("%w for %s", ErrUnknownMimeType, mime)
	}

	for _, codec := range codecs {
		if codec.Decodes(v) {
			return codec, nil
		}
	}

	return nil, fmt.Errorf("%w for '%s'", ErrUnknownMimeType, mime)
}
