// Package cli provides cli.
package cli

import (
	"errors"

	"go-micro.dev/v5/log"
	"go-micro.dev/v5/util/container"
)

var (
	// ErrFlagExists is returned when the given element exists in the flag container.
	ErrFlagExists = errors.New("element exists already")

	// Plugins contains source/cli subplugins, for example urfave, pflag, cobra.
	Plugins = container.NewMap[ParseFunc]() //nolint:gochecknoglobals
	// Flags is the global flag container where you have to register flags with.
	Flags = container.NewList[*Flag]() //nolint:gochecknoglobals
)

func init() {
	flag := NewFlag(
		"logger",
		log.DefaultPlugin,
		CPSlice([]string{"logger", "plugin"}),
		Usage("Default logger to use, jsonstderr, jsonstdout, textstderr, textsdout."),
	)

	if err := Flags.Add(flag); err != nil {
		panic(err)
	}
}
