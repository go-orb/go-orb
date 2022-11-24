package matcher

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"go-micro.dev/v5/util/container"
	"go-micro.dev/v5/util/slicemap"
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
		m.Add(test.selector, test.add...)
		res := m.Match(test.selector)
		for _, i := range test.expected {
			assert.EqualValues(t, true, slicemap.In(res, i), test.selector)
		}
	}
}

func TestMatcherJson(t *testing.T) {
	plugins := container.NewPlugins[string]()
	plugins.Register("abc", "abc")
	plugins.Register("def", "def")

	a := struct {
		Middlware Matcher[string] `json:"middleware"`
	}{
		Middlware: NewMatcher(plugins),
	}

	if err := json.Unmarshal([]byte(jsonFile), &a); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []string{"abc"}, a.Middlware.globals)
	assert.Equal(t, 1, len(a.Middlware.selectors))
	for _, val := range a.Middlware.selectors {
		assert.Equal(t, []string{"def"}, val)
	}
	assert.Equal(t, []string{"abc", "def"}, a.Middlware.Match("/foo"))
}

func TestMatcherYaml(t *testing.T) {
	plugins := container.NewPlugins[string]()
	plugins.Register("abc", "abc")
	plugins.Register("def", "def")

	a := struct {
		Middlware Matcher[string] `yaml:"middleware"`
	}{
		Middlware: NewMatcher(plugins),
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &a); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []string{"abc"}, a.Middlware.globals)
	assert.Equal(t, 1, len(a.Middlware.selectors))
	for _, val := range a.Middlware.selectors {
		assert.Equal(t, []string{"def"}, val)
	}
	assert.Equal(t, []string{"abc", "def"}, a.Middlware.Match("/foo9"))
}
