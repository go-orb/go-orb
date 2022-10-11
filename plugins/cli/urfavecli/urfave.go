package urfave

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"
	oCli "jochum.dev/jochumdev/orb/cli"
)

func init() {
	_ = oCli.Plugins.Add("urfave", NewCLI)
}

type FlagCLI struct {
	flags            map[string]*oCli.Flag
	stringFlags      map[string]*cli.StringFlag
	intFlags         map[string]*cli.IntFlag
	stringSliceFlags map[string]*cli.StringSliceFlag
	stringSliceDests map[string]*cli.StringSlice
	options          *oCli.Options
	ctx              *cli.Context
}

func NewCLI(opts ...oCli.Option) oCli.Cli {
	return &FlagCLI{
		flags:            make(map[string]*oCli.Flag),
		stringFlags:      make(map[string]*cli.StringFlag),
		intFlags:         make(map[string]*cli.IntFlag),
		stringSliceFlags: make(map[string]*cli.StringSliceFlag),
		stringSliceDests: make(map[string]*cli.StringSlice),
		options:          oCli.NewCLIOptions(),
	}
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

func (c *FlagCLI) Parse(args []string, opts ...oCli.Option) error {
	for _, o := range opts {
		o(c.options)
	}

	i := 0
	flags := make([]cli.Flag, len(c.stringFlags)+len(c.intFlags)+len(c.stringSliceFlags))
	for _, f := range c.stringFlags {
		flags[i] = f
		i += 1
	}
	for _, f := range c.intFlags {
		flags[i] = f
		i += 1
	}
	for _, f := range c.stringSliceFlags {
		flags[i] = f
		i += 1
	}

	var ctx *cli.Context
	app := &cli.App{
		Version:     c.options.Version,
		Description: c.options.Description,
		Usage:       c.options.Usage,
		Flags:       flags,
		Action: func(fCtx *cli.Context) error {
			// Extract the ctx from the urfave app
			ctx = fCtx

			return nil
		},
	}
	if len(c.options.Version) < 1 {
		app.HideVersion = true
	}

	if err := app.Run(args); err != nil {
		return err
	}
	c.ctx = ctx

	if c.ctx == nil {
		os.Exit(0)
	}

	for n, f := range c.stringSliceDests {
		_ = oCli.UpdateFlagValue(c.flags[n], f.Get())
	}

	return nil
}

func (c *FlagCLI) String() string {
	return "urfave"
}
