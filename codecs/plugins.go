package codecs

import (
	"errors"
	"fmt"

	"github.com/go-orb/go-orb/util/container"
	"github.com/go-orb/go-orb/util/slicemap"
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

// GetCodec returns the first codec by preference list.
func GetCodec(preference []string) (Marshaler, error) {
	var codec Marshaler

	allCodecs := Plugins.All()
	for name, c := range allCodecs {
		if slicemap.In(preference, name) {
			codec = c
		}
	}

	if codec == nil {
		return nil, fmt.Errorf("no matching codec plugin found for %v, please import atleast one of them", preference)
	}

	return codec, nil
}

var mimeMap = container.NewSafeMap[Marshaler]() //nolint:gochecknoglobals

func updateMimeMap() {
	for _, codec := range Plugins.All() {
		// One codec can support multiple mime types, we add all of them to the map.
		for _, mime := range codec.ContentTypes() {
			err := mimeMap.Add(mime, codec)
			if errors.Is(err, container.ErrExists) {
				continue
			}
		}
	}
}

// GetMime returns a codec for a mime type.
func GetMime(mime string) (Marshaler, error) {
	codec, err := mimeMap.Get(mime)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnknownMimeType, mime)
	}

	return codec, nil
}
