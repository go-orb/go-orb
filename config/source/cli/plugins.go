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
	// Flags is the global flag container where you have to register flags with.
	Flags = container.NewList[*Flag]()
)

func init() {
	// Logger can't import config/source/cli thats why this is here.
	Flags.Set(NewFlag(
		"logger",
		log.DefaultPlugin,
		ConfigPathSlice([]string{"logger", "plugin"}),
		Usage("Logger to use (e.g. slog)."),
	))

	Flags.Set(NewFlag(
		"config",
		[]string{},
		ConfigPathSlice([]string{"config"}),
		Usage("Config file"),
	))
}
