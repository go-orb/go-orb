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

var (
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
