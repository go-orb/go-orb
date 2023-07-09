package codecs

import (
	"fmt"

	"github.com/go-orb/go-orb/util/container"
	"github.com/go-orb/go-orb/util/slicemap"
)

// Plugins is the registry for codec plugins.
var Plugins = container.NewMap[Marshaler]() //nolint:gochecknoglobals

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
