package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringFlag(t *testing.T) {
	flag := NewFlag(
		"string",
		"",
		ConfigPathSlice([]string{"registry", "string"}),
		Usage("demo String flag"),
	)

	// Initial value of flag must be nil
	require.Nil(t, flag.Value)

	flag.Value = int64(0)
	_, err := FlagValue[string](flag)
	require.Error(t, err)

	flag.Value = "somevalue"
	v, err := FlagValue[string](flag)
	require.NoError(t, err)
	assert.Equal(t, "somevalue", v)
}

func TestIntFlag(t *testing.T) {
	flag := NewFlag(
		"int",
		300,
		ConfigPathSlice([]string{"registry", "int"}),
		Usage("demo Int flag"),
	)

	// Initial value of flag must be nil
	require.Nil(t, flag.Value)

	flag.Value = ""
	_, err := FlagValue[int](flag)
	require.Error(t, err)

	flag.Value = 10
	v, err := FlagValue[int](flag)
	require.NoError(t, err)
	assert.Equal(t, 10, v)
}

func TestStringSliceFlag(t *testing.T) {
	flag := NewFlag(
		"stringslice",
		[]string{"1", "2"},
		ConfigPathSlice([]string{"registry", "stringslice"}),
		Usage("demo StringSlice flag"),
	)

	// Initial value of flag must be nil
	require.Nil(t, flag.Value)

	flag.Value = ""
	_, err := FlagValue[[]string](flag)
	require.Error(t, err)

	flag.Value = []string{"a", "b"}
	v, err := FlagValue[[]string](flag)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, v)
}
