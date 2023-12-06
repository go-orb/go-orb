package config

import (
	"testing"

	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config/source"
	"github.com/go-orb/go-orb/types"
	"github.com/stretchr/testify/require"
)

func TestHasKey(t *testing.T) {
	codec, ok := codecs.Plugins.Get("json")
	require.True(t, ok)

	data := types.ConfigData{
		source.Data{
			Error:     nil,
			Marshaler: codec,
			Data: map[string]any{
				"com": map[string]any{
					"test": map[string]any{
						"registry": map[string]any{
							"plugin": "mdns",
						},
					},
				},
			},
		},
		source.Data{
			Error:     nil,
			Marshaler: codec,
			Data: map[string]any{
				"com": map[string]any{
					"test": map[string]any{
						"client": map[string]any{
							"plugin": "orb",
							"middlewares": []map[string]any{
								{
									"name":   "m1",
									"plugin": "log",
								},
								{
									"name":   "m2",
									"plugin": "trace",
								},
							},
						},
					},
				},
			},
		},
	}

	// require.NoError(t, Dump(data))

	require.True(t, HasKey[[]map[string]any](
		[]string{"com", "test", "client"}, "middlewares", data),
		"Should have key com.test.client.middlewares",
	)
	require.False(t, HasKey[[]map[string]any](
		[]string{"com", "test", "client"}, "middlewares2", data),
		"Should not have com.test.client.middlewares2",
	)
}
