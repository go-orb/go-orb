package slice

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

type test[T constraints.Ordered] struct {
	Array    []T
	Query    []T
	Expected []bool
}

var tests = []test[string]{
	{
		Array:    []string{"one", "two", "three", "1", "2", "3"},
		Query:    []string{"one", "five", "1"},
		Expected: []bool{true, false, true},
	},
}

func TestSlice(t *testing.T) {
	for i, test := range tests {
		t.Run("SliceTest"+strconv.Itoa(i), func(t *testing.T) {
			for q, query := range test.Query {
				assert.Equal(t, In(test.Array, query), test.Expected[q])
			}
		})
	}
}
