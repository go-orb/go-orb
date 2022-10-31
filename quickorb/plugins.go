package quickorb

import (
	_ "github.com/go-orb/config-plugins/source/cli"
	_ "github.com/go-orb/config-plugins/source/file"
	_ "github.com/go-orb/config-plugins/source/http"

	_ "github.com/go-orb/plugins/registry/mdnsregistry"
)
