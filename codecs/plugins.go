package codecs

import (
	"fmt"

	"github.com/go-orb/go-orb/util/container"
)

// Plugins is the registry for codec plugins.
var Plugins = container.NewMap[string, Marshaler]() //nolint:gochecknoglobals

// Register makes a plugin available by the provided name.
func Register(name string, codec Marshaler) bool {
	Plugins.Add(name, codec)
	updateMimeMap()

	return true
}

var mimeMap = container.NewSafeMap[string, Marshaler]() //nolint:gochecknoglobals
var extMap = container.NewSafeMap[string, Marshaler]()  //nolint:gochecknoglobals

func updateMimeMap() {
	Plugins.Range(func(_ string, encoder Marshaler) bool {
		for _, mime := range encoder.ContentTypes() {
			mimeMap.GetOrInsert(mime, encoder)
		}

		for _, ext := range encoder.Exts() {
			extMap.GetOrInsert(ext, encoder)
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

	return codec, nil
}

// GetExt returns a codec for a file extension.
func GetExt(ext string) (Marshaler, error) {
	codec, ok := extMap.Get(ext)
	if !ok {
		return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownExt, ext)
	}

	return codec, nil
}

// GetEncoder returns a encoder codec for a mime and golang type.
func GetEncoder(mime string, v any) (Marshaler, error) {
	codec, ok := mimeMap.Get(filterMime(mime))
	if !ok {
		return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownMimeType, mime)
	}

	if !codec.Marshals(v) {
		return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownValueType, v)
	}

	return codec, nil
}

// GetDecoder returns a decoder codec for a mime and golang type.
func GetDecoder(mime string, v any) (Marshaler, error) {
	codec, ok := mimeMap.Get(filterMime(mime))
	if !ok {
		return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownMimeType, mime)
	}

	if !codec.Unmarshals(v) {
		return nil, fmt.Errorf("%w for %s, did you import the codec plugin?", ErrUnknownValueType, v)
	}

	return codec, nil
}
