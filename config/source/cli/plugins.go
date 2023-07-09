// Package cli provides cli.
package cli

import (
	"errors"

	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/util/container"
)

var (
	// ErrFlagExists is returned when the given element exists in the flag container.
	ErrFlagExists = errors.New("element exists already")
)

//nolint:gochecknoglobals
var (
	// Plugins contains source/cli subplugins, for example urfave, pflag, cobra.
	Plugins = container.NewMap[ParseFunc]()
	// Flags is the global flag container where you have to register flags with.
	Flags = container.NewList[*Flag]()
)

func init() {
	flag := NewFlag(
		"logger",
		log.DefaultPlugin,
		ConfigPathSlice([]string{"logger", "plugin"}),
		Usage("Default logger to use (e.g. jsonstderr, jsonstdout, textstderr, textsdout)."),
	)

	if err := Flags.Add(flag); err != nil {
		panic(err)
	}
}
