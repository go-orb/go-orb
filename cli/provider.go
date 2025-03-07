package cli

import (
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

	mJSON, err := codecs.GetMime("application/json")
	if err != nil {
		return nil, err
	}

	cliResult := source.Data{
		Data:      make(map[string]any),
		Marshaler: mJSON,
	}

	results = append(results, cliResult)

	for _, flag := range flags {
		// The config flag is a special case, as you can add additional config files.
		// E.g. `--config cfg-a.yaml --config cfg-b.yaml`, here we keep track of them.
		if flag.Name == "config" {
			var (
				urls []string
				ok   bool
			)

			if urls, ok = flag.Value.([]string); !ok {
				// We ignore this here if the user developed another config variable.
				continue
			}

			for _, t := range urls {
				u, err := url.Parse(t)
				if err != nil {
					return nil, err
				}

				config, err := config.Read([]*url.URL{u})
				if err != nil {
					return nil, err
				}

				results = append(results, config...)
			}

			continue
		}

		flagToMap(types.SplitServiceName(serviceContext.Name()), flag, cliResult.Data)
	}

	return results, nil
}
