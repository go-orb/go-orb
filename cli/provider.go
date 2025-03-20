package cli

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/types"
)

// ProvideSingleServiceContext provides a single service context for the application.
func ProvideSingleServiceContext(appContext *AppContext) (*ServiceContext, error) {
	return NewServiceContext(appContext, appContext.Name(), appContext.Version()), nil
}

// ProvideParsedFlagsFromArgs provides parsed flags from the app context.
func ProvideParsedFlagsFromArgs(appContext *AppContext, parser ParserFunc, args []string) ([]*Flag, error) {
	return parser(appContext, args)
}

func flagToMap(globalSections []string, multiServiceConfig bool, flag *Flag, cliResult map[string]any) {
	for _, cp := range flag.ConfigPaths {
		sections := cp.Path[:len(cp.Path)-1]

		if !cp.IsGlobal && !multiServiceConfig {
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

// AppConfigData is the config data type.
type AppConfigData map[string]any

// ServiceContextHasConfigData is a marker type.
type ServiceContextHasConfigData struct{}

// ProvideAppConfigData provides config data from appContext and flags.
func ProvideAppConfigData(appContext *AppContext) (AppConfigData, error) {
	cfg := map[string]any{}

	// Load configs from memory (App().HardcodedConfigs).
	if err := loadHardcodedConfigs(appContext, cfg); err != nil {
		return AppConfigData{}, err
	}

	// Load configs from URLs (App().HardcodedConfigURLs).
	if err := loadHardcodedConfigURLs(appContext, cfg); err != nil {
		return AppConfigData{}, err
	}

	return AppConfigData(cfg), nil
}

// ProvideServiceConfigData provides config data to serviceContext from flags.
func ProvideServiceConfigData(
	serviceContext *ServiceContext,
	appConfigData AppConfigData,
	flags []*Flag,
) (ServiceContextHasConfigData, error) {
	result := map[string]any(appConfigData)

	// Process command-line flags.
	cfg, err := processFlags(serviceContext, flags)
	if err != nil {
		return ServiceContextHasConfigData{}, err
	}

	if err := config.Merge(&result, cfg); err != nil {
		return ServiceContextHasConfigData{}, err
	}

	// Finally, set the config on the service context.
	serviceContext.Config = result

	return ServiceContextHasConfigData{}, nil
}

// loadHardcodedConfigs loads configs from memory strings (serviceContext.App().HardcodedConfigs).
func loadHardcodedConfigs(appContext *AppContext, into map[string]any) error {
	app := appContext.App()
	if app.HardcodedConfigs == nil {
		return nil
	}

	for i, configData := range app.HardcodedConfigs {
		b64 := base64.URLEncoding.EncodeToString([]byte(configData.Data))
		urlString := fmt.Sprintf("file:///memory%d.%s?base64=%s", i, configData.Format, b64)

		cfg, err := loadConfigFromURL(urlString)
		if err != nil {
			return err
		}

		if err := config.Merge(&into, cfg); err != nil {
			return err
		}
	}

	return nil
}

// loadHardcodedConfigURLs loads configs from URL strings (serviceContext.App().HardcodedConfigURLs).
func loadHardcodedConfigURLs(appContext *AppContext, into map[string]any) error {
	app := appContext.App()
	if app.HardcodedConfigURLs == nil {
		return nil
	}

	for _, urlString := range app.HardcodedConfigURLs {
		cfg, err := loadConfigFromURL(urlString)
		if err != nil {
			return err
		}

		if err := config.Merge(&into, cfg); err != nil {
			return err
		}
	}

	return nil
}

// loadConfigFromURL loads config from a single URL string.
func loadConfigFromURL(urlString string) (map[string]any, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %s: %w", urlString, err)
	}

	cfg, err := config.Read(u)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from %s: %w", urlString, err)
	}

	return cfg, nil
}

// processFlags processes CLI flags and loads any config files specified.
func processFlags(serviceContext *ServiceContext, flags []*Flag) (map[string]any, error) {
	cliData := map[string]any{}

	for _, flag := range flags {
		// Special handling for --config flag which loads additional config files.
		if flag.Name == "config" {
			urls, ok := flag.Value.([]string)
			if !ok {
				continue // Skip if not a []string (user redefined config flag).
			}

			for _, urlString := range urls {
				cfg, err := loadConfigFromURL(urlString)
				if err != nil {
					return nil, err
				}

				if err := config.Merge(&cliData, cfg); err != nil {
					return nil, err
				}
			}

			continue
		}

		// Add regular flags to the CLI config data.
		flagToMap(types.SplitServiceName(serviceContext.Name()), serviceContext.App().MultiServiceConfig, flag, cliData)
	}

	return cliData, nil
}
