// Package cli is the Cli component of orb.
package cli

import (
	"errors"
	"net/url"
	"os"

	"go-micro.dev/v5/codecs"
	"go-micro.dev/v5/config/source"
	"go-micro.dev/v5/log"
	"go-micro.dev/v5/util/container"
)

// ParseFunc is the subplugin of source/cli.
type ParseFunc func(config *Config, flags []*Flag, args []string) error

func init() {
	err := source.Plugins.Add(New())
	if err != nil {
		panic(err)
	}
}

// Source is the cli source.
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
		pName = "urfave"
	}

	// Add the config flag.
	err := Flags.Add(NewFlag(
		"config",
		[]string{},
		CPSlice([]string{"config"}),
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
		result.Error = err
		return result
	}

	err = parseFunc(config, Flags.List(), os.Args)
	if err != nil {
		result.Error = err
		return result
	}

	// Parse all Flags into map[string]any.
	for _, flag := range Flags.List() {
		// Special case the "config" flag.
		if flag.Name == "config" {
			switch tmp := flag.Value.(type) {
			case []string:
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

		// All the other flags.
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
		log.Error("no json encoder compiled in, will fail now", err)
	}

	result.Marshaler = mJSON

	return result
}
