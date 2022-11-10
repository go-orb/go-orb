package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testJSON = `{
	"string": "value",
	"stringslice": [
		"value1",
		"value2",
		0,
		true
	],
	"mixedslice": [
		"value1",
		0,
		1,
		2
	],
	"stringmap": {
		"key1": "value1",
		"key2": "value2",
		"key3": 1,
		"key4": true
	}
}`

func testData(t *testing.T) map[string]any {
	t.Helper()

	data := make(map[string]any)
	if err := json.Unmarshal([]byte(testJSON), &data); err != nil {
		t.Fatalf("error while reading testJSON: %v", err)
	}

	return data
}

func TestReadString(t *testing.T) {
	data := testData(t)

	// Must return the correct value.
	str, err := Get(data, "string", "x")
	assert.Nil(t, err)
	assert.Equal(t, "value", str)

	// Must return default if type don't match
	i, err := Get(data, "string", 0)
	assert.ErrorIs(t, err, ErrTypesDontMatch)
	assert.Equal(t, i, 0)

	// Must return default
	str, err = Get(data, "string2", "x")
	assert.ErrorIs(t, err, ErrNotExistent)
	assert.Equal(t, str, "x")
}

func TestReadStringSlice(t *testing.T) {
	data := testData(t)

	// Must return the correct value.
	strs, err := Get(data, "stringslice", []string{})
	assert.Nil(t, nil, err)
	assert.Equal(t, []string{"value1", "value2", "0", "true"}, strs)

	// Must return error if not a slice
	_, err = Get(data, "string", []string{})
	assert.ErrorIs(t, err, ErrTypesDontMatch)

	// Must return default
	strs, err = Get(data, "stringslice2", []string{"a", "b"})
	assert.ErrorIs(t, err, ErrNotExistent)
	assert.Equal(t, []string{"a", "b"}, strs)
}

func TestReadMixedSlice(t *testing.T) {
	data := testData(t)

	// Must return the correct value.
	anys, err := Get(data, "mixedslice", []any{})
	assert.Nil(t, err)
	assert.Equal(t, []any{"value1", float64(0), float64(1), float64(2)}, anys)

	// Must return error if not a slice
	_, err = Get(data, "string", []any{})
	assert.ErrorIs(t, err, ErrTypesDontMatch)
}

func TestReadStringMap(t *testing.T) {
	data := testData(t)

	// Must return the correct value.
	maps, err := Get(data, "stringmap", map[string]string{})
	assert.Nil(t, nil, err)
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2", "key3": "1", "key4": "true"}, maps)

	// Must return error if not a map.
	_, err = Get(data, "string", map[string]string{})
	assert.ErrorIs(t, err, ErrTypesDontMatch)

	// Must return default
	maps, err = Get(data, "stringslice2", map[string]string{"a": "a"})
	assert.ErrorIs(t, err, ErrNotExistent)
	assert.Equal(t, map[string]string{"a": "a"}, maps)
}

func TestReadMixedMap(t *testing.T) {
	data := testData(t)

	// Must return the correct value.
	mapa, err := Get(data, "stringmap", map[string]any{})
	assert.Nil(t, err)
	assert.Equal(t, map[string]any{"key1": "value1", "key2": "value2", "key3": float64(1), "key4": true}, mapa)

	// Must return error if not a slice
	_, err = Get(data, "string", map[string]any{})
	assert.ErrorIs(t, err, ErrTypesDontMatch)
}
