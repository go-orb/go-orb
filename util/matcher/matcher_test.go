package matcher

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/go-orb/go-orb/util/container"
	"github.com/go-orb/go-orb/util/slicemap"
)

var (
	jsonFile = `{"middleware": ["abc", {"name": "def", "selector": "/foo"}]}`
	yamlFile = `
middleware:
 - abc
 - name: def
   selector: "/foo[1-9]"
`
)

type test struct {
	selector string
	add      []int
	expected []int
}

var tests = []test{
	{selector: "/*", add: []int{1, 2}, expected: []int{1, 2}},
	{selector: "/helloworld", add: []int{3}, expected: []int{1, 2, 3}},
	{selector: "/helloworld/echo", add: []int{4}, expected: []int{1, 2, 3, 4}},
	{selector: "/foo", add: []int{}, expected: []int{1, 2}},
	{selector: "/foo/*", add: []int{5}, expected: []int{1, 2, 5}},
	{selector: "/foo/bar", add: []int{}, expected: []int{1, 2, 5}},
}

func TestMatcher(t *testing.T) {
	m := NewMatcher[int](nil)

	for _, test := range tests {
		// Add test cases.
		for _, item := range test.add {
			name := strconv.Itoa(item)
			m.Add(test.selector, name, item)
		}

		// Verify test cases.
		res := m.Match(test.selector)
		for _, i := range test.expected {
			assert.EqualValues(t, true, slicemap.In(res, i), test.selector)
		}
	}
}

func TestMatcherDuplication(t *testing.T) {
	plugins := container.NewMap[string]()
	plugins.Set("one", "itemOne")
	plugins.Set("two", "itemTwo")

	m := NewMatcher(plugins)

	m.Use("customOne", "itemOneC")
	m.Use("customOne", "itemOneC")
	m.Use("customOne", "itemOneC")

	assert.Len(t, m.globals, 1, "customOne")

	require.NoError(t, m.AddPlugin("/helloworld", "one"))
	require.NoError(t, m.AddPlugin("/helloworld", "one"))

	var seen bool
	for _, i := range m.selectors {
		assert.Len(t, i, 1)
		seen = true
	}
	assert.True(t, seen)

	require.NoError(t, m.AddPlugin("/helloworld", "two"))
	require.NoError(t, m.AddPlugin("/helloworld", "two"))

	seen = false
	for _, i := range m.selectors {
		assert.Len(t, i, 2)
		seen = true
	}
	assert.True(t, seen)
}

func TestMatcherJson(t *testing.T) {
	plugins := container.NewMap[string]()
	plugins.Set("abc", "abc")
	plugins.Set("def", "def")

	a := struct {
		Middlware Matcher[string] `json:"middleware"`
	}{
		Middlware: NewMatcher(plugins),
	}

	if err := json.Unmarshal([]byte(jsonFile), &a); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, a.Middlware.globals, 1)
	assert.Len(t, a.Middlware.selectors, 1)
	assert.Equal(t, []string{"abc"}, a.Middlware.Match("/bar"))
	assert.Equal(t, []string{"abc", "def"}, a.Middlware.Match("/foo"))
}

func TestMatcherYaml(t *testing.T) {
	plugins := container.NewMap[string]()
	plugins.Set("abc", "abc")
	plugins.Set("def", "def")

	a := struct {
		Middlware Matcher[string] `yaml:"middleware"`
	}{
		Middlware: NewMatcher(plugins),
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &a); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, a.Middlware.globals, 1)
	assert.Len(t, a.Middlware.selectors, 1)
	assert.Equal(t, []string{"abc"}, a.Middlware.Match("/bar"))
	assert.Equal(t, []string{"abc", "def"}, a.Middlware.Match("/foo9"))
}
