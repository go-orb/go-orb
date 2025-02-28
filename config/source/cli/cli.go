// Package cli provides the CLI config component of go-micro.
package cli

import (
	"net/url"

	"github.com/go-orb/go-orb/config/source"
)

// ParseFunc is the subplugin of source/cli.
type ParseFunc func() source.Data

var _ (source.Source) = (*Source)(nil)

// Source cli reads flags and environment variables into a config struct.
type Source struct {
	parser ParseFunc
}

// New creates a new cli source.
func New(parser ParseFunc) source.Source {
	return &Source{parser: parser}
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
func (s *Source) Read(_ *url.URL) source.Data {
	return s.parser()
}

// ParseFlags takes the list of flags and parses them into a map[string]any
// contained inside the result.
func ParseFlags(result *source.Data, flags []*Flag) {
	for _, flag := range flags {
		// The config flag is a special case, as you can add additional config files.
		// E.g. `--config cfg-a.yaml --config cfg-b.yaml`, here we keep track of them.
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
}
