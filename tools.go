//go:build tools
// +build tools

// following https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package orb

import (
	_ "github.com/go-orb/wire/cmd/wire"
)
