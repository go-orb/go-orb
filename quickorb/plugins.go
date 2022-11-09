package quickorb

import (
	_ "go-micro.dev/v5/config-plugins/source/cli/urfave"
	_ "go-micro.dev/v5/config-plugins/source/file"
	_ "go-micro.dev/v5/config-plugins/source/http"
	_ "go-micro.dev/v5/config/source/cli"

	_ "github.com/go-orb/plugins/registry/mdnsregistry"
)
