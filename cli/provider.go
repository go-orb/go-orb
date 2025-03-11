package cli

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/config/source"
	"github.com/go-orb/go-orb/types"
)

// ProvideSingleServiceContext provides a single service context for the application.
func ProvideSingleServiceContext(appContext *AppContext) (*ServiceContext, error) {
	return NewServiceContext(appContext, appContext.Name(), appContext.Version()), nil
}

// ProvideServiceName extracts the service name from the service context.
func ProvideServiceName(serviceContext *ServiceContext) (types.ServiceName, error) {
	return types.ServiceName(serviceContext.Name()), nil
}

// ProvideServiceVersion extracts the service version from the service context.
func ProvideServiceVersion(serviceContext *ServiceContext) (types.ServiceVersion, error) {
	return types.ServiceVersion(serviceContext.Version()), nil
}

// ProvideParsedFlagsFromArgs provides parsed flags from the app context.
func ProvideParsedFlagsFromArgs(appContext *AppContext, parser ParserFunc, args []string) ([]*Flag, error) {
	return parser(appContext, args)
}

func flagToMap(globalSections []string, flag *Flag, cliResult map[string]any) {
	for _, cp := range flag.ConfigPaths {
		sections := cp.Path[:len(cp.Path)-1]

		if !cp.IsGlobal {
			sections = append(globalSections, sections...)
		}

		data := cliResult
		for _, s := range sections {
			if tmp, ok := data[s]; ok {
				switch t2 := tmp.(type) {
				case map[string]any:
					data = t2
				default:
					// Should never happen.
					data = cliResult
				}
			} else {
				tmp := map[string]any{}
				data[s] = tmp
				data = tmp
			}
		}

		data[cp.Path[len(cp.Path)-1]] = flag.Value
	}
}

// ProvideConfigData provides config data from serviceContext and flags.
func ProvideConfigData(serviceContext *ServiceContext, flags []*Flag) (types.ConfigData, error) {
	results := types.ConfigData{}

	// Load configs from memory (App().Configs and App().ConfigsFormat).
	if err := loadInMemoryConfigs(serviceContext, &results); err != nil {
		return nil, err
	}

	// Load configs from URLs (App().ConfigURLs).
	if err := loadConfigURLs(serviceContext, &results); err != nil {
		return nil, err
	}

	// Initialize CLI-based config.
	mJSON, err := codecs.GetMime(codecs.MimeJSON)
	if err != nil {
		return nil, err
	}

	cliResult := source.Data{
		Data:      make(map[string]any),
		Marshaler: mJSON,
	}
	results = append(results, cliResult)

	// Process command-line flags.
	if err := processFlags(serviceContext, flags, &results, cliResult.Data); err != nil {
		return nil, err
	}

	return results, nil
}

// loadInMemoryConfigs loads configs from memory strings (serviceContext.App().Configs).
func loadInMemoryConfigs(serviceContext *ServiceContext, results *types.ConfigData) error {
	app := serviceContext.App()
	if app.Configs == nil || app.ConfigsFormat == nil || len(app.Configs) != len(app.ConfigsFormat) {
		return nil
	}

	for i, configData := range app.Configs {
		b64 := base64.URLEncoding.EncodeToString([]byte(configData))
		urlString := fmt.Sprintf("file:///memory%d.%s?base64=%s", i, app.ConfigsFormat[i], b64)

		config, err := loadConfigFromURL(urlString)
		if err != nil {
			return err
		}

		*results = append(*results, config...)
	}

	return nil
}

// loadConfigURLs loads configs from URL strings (serviceContext.App().ConfigURLs).
func loadConfigURLs(serviceContext *ServiceContext, results *types.ConfigData) error {
	app := serviceContext.App()
	if app.ConfigURLs == nil {
		return nil
	}

	for _, urlString := range app.ConfigURLs {
		config, err := loadConfigFromURL(urlString)
		if err != nil {
			return err
		}

		*results = append(*results, config...)
	}

	return nil
}

// loadConfigFromURL loads config from a single URL string.
func loadConfigFromURL(urlString string) (types.ConfigData, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %s: %w", urlString, err)
	}

	config, err := config.Read([]*url.URL{u})
	if err != nil {
		return nil, fmt.Errorf("failed to read config from %s: %w", urlString, err)
	}

	return config, nil
}

// processFlags processes CLI flags and loads any config files specified.
func processFlags(serviceContext *ServiceContext, flags []*Flag, results *types.ConfigData, cliData map[string]any) error {
	for _, flag := range flags {
		// Special handling for --config flag which loads additional config files.
		if flag.Name == "config" {
			urls, ok := flag.Value.([]string)
			if !ok {
				continue // Skip if not a []string (user redefined config flag).
			}

			for _, urlString := range urls {
				config, err := loadConfigFromURL(urlString)
				if err != nil {
					return err
				}

				*results = append(*results, config...)
			}

			continue
		}

		// Add regular flags to the CLI config data.
		flagToMap(types.SplitServiceName(serviceContext.Name()), flag, cliData)
	}

	return nil
}
