package codecs

import "go-micro.dev/v5/util/container"

// Plugins is the registry for codec plugins.
var Plugins = container.NewMap[Marshaler]() //nolint:gochecknoglobals
