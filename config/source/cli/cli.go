// Package cli provides the Cli component of orb.
package cli

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"go-micro.dev/v5/codecs"
	"go-micro.dev/v5/config/source"
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/util/container"
)

var (
	// DefaultCLIPlugin holds the default CLI plugin.
	DefaultCLIPlugin = "urfave" //nolint:gochecknoglobals
)

// ParseFunc is the subplugin of source/cli.
type ParseFunc func(config *Config, flags []*Flag, args []string) error

func init() {
	if err := source.Plugins.Add(New()); err != nil {
		panic(err)
	}
}

// Source cli reads flags and environment variables into a config struct.
type Source struct{}

// New creates a new cli source.
func New() source.Source {
	return &Source{}
}

// Schemes returns the supported schemes by this plugin.
func (s *Source) Schemes() []string {
	return []string{"cli"}
}

// PrependSections indicates whether this needs sections to be prepended,
// which is true in this case.
func (s *Source) PrependSections() bool {
	return true
}

// String returns the name of this plugin.
func (s *Source) String() string {
	return "cli"
}

// Read creates the subplugin for the given url,
// creates its config after and then executes it.
func (s *Source) Read(u *url.URL) source.Data {
	result := source.Data{
		Data: make(map[string]any),
	}

	pName := u.Host
	if pName == "" {
		pName = DefaultCLIPlugin
	}

	// Add the config flag.
	err := Flags.Add(NewFlag(
		"config",
		[]string{},
		ConfigPathSlice([]string{"config"}),
		Usage("Config file"),
	))
	if err != nil && !errors.Is(err, container.ErrExists) {
		result.Error = err
		return result
	}

	config := NewConfig()
	config.Name = u.Query().Get("name")
	config.Version = u.Query().Get("version")

	parseFunc, err := Plugins.Get(pName)
	if err != nil {
		result.Error = fmt.Errorf(
			"failed to get the plugin '%s'. Did you register the plugin by importing it?, error was: %w",
			pName,
			err,
		)

		return result
	}

	if err = parseFunc(&config, Flags.List(), os.Args); err != nil {
		result.Error = err
		return result
	}

	// Parse all Flags into map[string]any.
	for _, flag := range Flags.List() {
		// Special case the "config" flag which might add additionalConfigs.
		if flag.Name == "config" {
			if tmp, ok := flag.Value.([]string); ok {
				for _, t := range tmp {
					u, err := url.Parse(t)
					if err != nil {
						continue
					}

					result.AdditionalConfigs = append(result.AdditionalConfigs, u)
				}
			}

			continue
		}

		sections := flag.ConfigPath[:len(flag.ConfigPath)-1]

		data := result.Data
		for _, s := range sections {
			if tmp, ok := data[s]; ok {
				switch t2 := tmp.(type) {
				case map[string]any:
					data = t2
				default:
					// Should never happen.
					data = result.Data
				}
			} else {
				tmp := map[string]any{}
				data[s] = tmp
				data = tmp
			}
		}

		data[flag.ConfigPath[len(flag.ConfigPath)-1]] = flag.Value
	}

	mJSON, err := codecs.Plugins.Get("json")
	if err != nil {
		log.Error("JSON codec was not found, did you register it by importing?", err)
	}

	result.Marshaler = mJSON

	return result
}
