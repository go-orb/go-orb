package cli

import (
	"net/url"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/config/source"
	"github.com/go-orb/go-orb/types"
)

// ProvideConfigData provides configData from cli, this requires for example urfave.Provide to be registered first.
func ProvideConfigData(serviceName types.ServiceName, cliParser ParseFunc) (types.ConfigData, error) {
	if err := source.Plugins.Add(New(cliParser)); err != nil {
		return nil, err
	}

	u, err := url.Parse("cli://")
	if err != nil {
		return nil, err
	}

	cfgSections := types.SplitServiceName(serviceName)

	data, err := config.Read([]*url.URL{u}, cfgSections)

	return data, err
}
