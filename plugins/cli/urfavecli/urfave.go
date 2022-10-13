// Package urfavecli is a cli wrapper for urfave.
package urfavecli

import (
	"errors"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"
	oCli "jochum.dev/orb/orb/cli"
	"jochum.dev/orb/orb/config/chelp"
)

func init() {
	if err := oCli.Plugins.Add(
		"urfave",
		New,
		func() any { return oCli.NewConfig() },
	); err != nil {
		panic(err)
	}
}

type FlagCLI struct {
	flags            map[string]*oCli.Flag
	stringFlags      map[string]*cli.StringFlag
	intFlags         map[string]*cli.IntFlag
	stringSliceFlags map[string]*cli.StringSliceFlag
	stringSliceDests map[string]*cli.StringSlice
	config           oCli.Config
	ctx              *cli.Context
}

func New() oCli.Cli {
	return &FlagCLI{
		flags:            make(map[string]*oCli.Flag),
		stringFlags:      make(map[string]*cli.StringFlag),
		intFlags:         make(map[string]*cli.IntFlag),
		stringSliceFlags: make(map[string]*cli.StringSliceFlag),
		stringSliceDests: make(map[string]*cli.StringSlice),
	}
}

func (c *FlagCLI) Init(aConfig any) error {
	switch config := aConfig.(type) {
	case *oCli.BaseConfig:
		c.config = config
	default:
		return chelp.ErrUnknownConfig
	}

	return nil
}

func (c *FlagCLI) Config() any {
	return c.config
}

func (c *FlagCLI) Add(opts ...oCli.FlagOption) error {
	flag, err := oCli.NewFlag(opts...)
	if err != nil {
		return err
	}

	switch flag.FlagType {
	case oCli.FlagTypeInt:
		f := &cli.IntFlag{
			Name:        flag.Name,
			Usage:       flag.Usage,
			Value:       flag.DefaultInt,
			EnvVars:     flag.EnvVars,
			Destination: &flag.ValueInt,
		}
		c.intFlags[flag.Name] = f
	case oCli.FlagTypeString:
		f := &cli.StringFlag{
			Name:        flag.Name,
			Usage:       flag.Usage,
			Value:       flag.DefaultString,
			EnvVars:     flag.EnvVars,
			Destination: &flag.ValueString,
		}
		c.stringFlags[flag.Name] = f
	case oCli.FlagTypeStringSlice:
		dest := cli.NewStringSlice()
		c.stringSliceDests[flag.Name] = dest
		f := &cli.StringSliceFlag{
			Name:        flag.Name,
			Usage:       flag.Usage,
			Value:       cli.NewStringSlice(flag.DefaultStringSlice...),
			EnvVars:     flag.EnvVars,
			Destination: dest,
		}
		c.stringSliceFlags[flag.Name] = f
	default:
		return errors.New("found a flag without a default option")
	}

	c.flags[flag.Name] = flag

	return nil
}

func (c *FlagCLI) Get(name string) (*oCli.Flag, bool) {
	flag, ok := c.flags[name]
	return flag, ok
}

func (c *FlagCLI) Parse(args []string) error {
	if c.config == nil {
		return oCli.ErrConfigIsNil
	}

	i := 0
	flags := make([]cli.Flag, len(c.stringFlags)+len(c.intFlags)+len(c.stringSliceFlags))

	for _, f := range c.stringFlags {
		flags[i] = f
		i++
	}

	for _, f := range c.intFlags {
		flags[i] = f
		i++
	}

	for _, f := range c.stringSliceFlags {
		flags[i] = f
		i++
	}

	var ctx *cli.Context

	app := &cli.App{
		Version:     c.config.Version(),
		Description: c.config.Description(),
		Usage:       c.config.Usage(),
		Flags:       flags,
		Action: func(fCtx *cli.Context) error {
			// Extract the ctx from the urfave app
			ctx = fCtx

			return nil
		},
	}

	if len(c.config.Version()) < 1 {
		app.HideVersion = true
	}

	if err := app.Run(args); err != nil {
		return err
	}

	c.ctx = ctx

	if c.ctx == nil {
		os.Exit(0)
	}

	var err error
	for n, f := range c.stringSliceDests {
		if err = oCli.UpdateFlagValue(c.flags[n], f.Get()); err != nil {
			err = multierror.Append(err)
		}
	}

	return err
}

func (c *FlagCLI) String() string {
	return "urfave"
}
