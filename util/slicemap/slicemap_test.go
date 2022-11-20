package slicemap

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

type sliceTest[T constraints.Ordered] struct {
	Array    []T
	Query    []T
	Expected []bool
}

type lookupTest struct {
	Map      map[string]any
	Query    [][]string
	Expected []any
}

type setValueTest struct {
	Input    map[string]any
	Path     []string
	Value    any
	Expected map[string]any
}

var (
	setValueTests = []setValueTest{
		{
			Input: map[string]any{},
			Path:  []string{"my", "custom", "value"},
			Value: 69,
			Expected: map[string]any{"my": map[string]any{
				"custom": map[string]any{
					"value": 69,
				},
			},
			},
		},
		{
			Input: map[string]any{
				"com": map[string]any{
					"example": map[string]any{
						"service": map[string]any{
							"address": ":5050",
						},
					},
				},
			},
			Path:  []string{"com", "example", "service", "address"},
			Value: ":8080",
			Expected: map[string]any{
				"com": map[string]any{
					"example": map[string]any{
						"service": map[string]any{
							"address": ":8080",
						},
					},
				},
			},
		},
	}

	sliceTests = []sliceTest[string]{
		{
			Array:    []string{"one", "two", "three", "1", "2", "3"},
			Query:    []string{"one", "five", "1"},
			Expected: []bool{true, false, true},
		},
	}

	mapTests = []lookupTest{
		{
			Map: map[string]any{
				"one": map[string]any{
					"two": map[string]any{
						"three": map[string]any{
							"int":    1,
							"string": "myfield",
						},
						"twoA": map[string]any{
							"field": 5,
						},
					},
				},
			},
			Query: [][]string{
				{"one", "two", "three", "int"},
				{"one", "two", "three", "string"},
				{"one", "two", "twoA"},
				{"one", "two", "twoB"},
			},
			Expected: []any{
				1,
				"myfield",
				map[string]any{
					"field": 5,
				},
				nil,
			},
		},
	}
)

func TestSlice(t *testing.T) {
	for i, test := range sliceTests {
		t.Run("SliceTest"+strconv.Itoa(i), func(t *testing.T) {
			for q, query := range test.Query {
				assert.Equal(t, In(test.Array, query), test.Expected[q])
			}
		})
	}
}

func TestLookup(t *testing.T) {
	for i, test := range mapTests {
		t.Run("MapTest"+strconv.Itoa(i), func(t *testing.T) {
			for q, query := range test.Query {
				v, err := Lookup(test.Map, query)
				if test.Expected[q] == nil {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, test.Expected[q], v)
			}
		})
	}
}

func TestSetValue(t *testing.T) {
	for i, test := range setValueTests {
		t.Run("SetValueTest"+strconv.Itoa(i), func(t *testing.T) {
			SetValue(test.Input, test.Path, test.Value)
			assert.Equal(t, test.Expected, test.Input)
		})
	}
}

func TestGet(t *testing.T) {
	a := map[string]any{
		"one": 5,
		"two": 5,
	}

	val, ok := Get[int](a, "one")
	assert.Equal(t, true, ok)
	assert.Equal(t, 5, val)

	_, ok = Get[string](a, "two")
	assert.Equal(t, false, ok)

	_, ok = Get[string](a, "three")
	assert.Equal(t, false, ok)
}
