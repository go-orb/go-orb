// Package matcher provides a generic map with regex keys, that can be used to
// match items to paths, e.g. when using middleware.
package matcher

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/sanity-io/litter"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	"go-micro.dev/v5/util/container"
	"go-micro.dev/v5/util/slicemap"
)

type selectorKey struct {
	selector string
	re       *regexp.Regexp
}

type itemContainer[T any] struct {
	Name string
	Item T
}

// Matcher is a map with regular expressions as keys. It is used to add
// middleware on a path based level.
//
// In addition to setting regex selectors for each item, they are also
// de-duplocated with a name key, to make sure each item is only added once.
type Matcher[T any] struct {
	globals   []itemContainer[T]
	plugins   *container.Plugins[T]
	selectors map[selectorKey][]itemContainer[T]
}

// NewMatcher creates a new matcher object.
func NewMatcher[T any](plugins *container.Plugins[T]) Matcher[T] {
	return Matcher[T]{
		plugins:   plugins,
		selectors: make(map[selectorKey][]itemContainer[T]),
	}
}

// Use will use the elements provided on all paths.
func (m *Matcher[T]) Use(name string, item T) {
	if itemPresent(m.globals, name) {
		return
	}

	m.globals = append(m.globals, itemContainer[T]{Name: name, Item: item})
}

// AddPlugin will add plugin item, with a selector.
func (m *Matcher[T]) AddPlugin(selector, plugin string) error {
	item, err := m.plugins.Get(plugin)
	if err != nil {
		return fmt.Errorf("plugin not found '%s'", plugin)
	}

	m.Add(selector, plugin, item)

	return nil
}

// Add will add an item to every item the selector matches. The selector
// is a regexp.
//
// Example selector:
//   - /*
//   - /echo
//   - /echo/*
//   - /echo[1-9]
//   - suffix$
func (m *Matcher[T]) Add(selector, name string, item T) {
	switch selector {
	case "/*":
		fallthrough
	case "*":
		fallthrough
	case ".*":
		fallthrough
	case "^.*":
		m.Use(name, item)
		return
	}

	if strings.HasSuffix(selector, "/*") {
		selector = strings.TrimSuffix(selector, "/*")
		selector += "/.*"
	}

	ic := itemContainer[T]{Name: name, Item: item}

	// Check if we already have a similar selector. Only add if not present
	for key := range m.selectors {
		if key.selector == selector && !itemPresent(m.selectors[key], name) {
			m.selectors[key] = append(m.selectors[key], ic)
			return
		}
	}

	// Create a new selector
	re, err := regexp.Compile(selector)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to compile selector as regexp '%s'", selector), err)
		return
	}

	s := selectorKey{
		selector: selector,
		re:       re,
	}

	m.selectors[s] = []itemContainer[T]{ic}
}

// Match will fetch the list of items that match against a path.
func (m *Matcher[T]) Match(operation string) []T {
	ms := extractItems(m.globals)

	if m.selectors == nil {
		return ms
	}

	for selector, val := range m.selectors {
		if selector.re.MatchString(operation) {
			ms = append(ms, extractItems(val)...)
		}
	}

	return ms
}

// Len returns the total number of items defined.
func (m Matcher[T]) Len() int {
	return len(m.globals) + len(m.selectors)
}

func extractItems[T any](items []itemContainer[T]) []T {
	output := make([]T, 0, len(items))

	for _, item := range items {
		output = append(output, item.Item)
	}

	return output
}

func itemPresent[T any](items []itemContainer[T], query string) bool {
	for _, item := range items {
		if item.Name == query {
			return true
		}
	}

	return false
}

// UnmarshalJSON will unmarshal a JSON file into the matcher.
// JSON can contain either a string or <name, selector> map.
//
// All items provided through the config need to be registered as plugins.
// All config will add to, not replace the exisiting items.
//
// JSON config example:
//
//	{
//	   "middleware":[
//	      "abc",
//	      {
//	         "name":"def",
//	         "selector":"/foo"
//	      }
//	   ]
//	}
func (m *Matcher[T]) UnmarshalJSON(data []byte) error {
	var a any

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	b, ok := a.([]any)
	if !ok {
		return nil
	}

	return m.unmarshal(b)
}

// UnmarshalYAML will unmarshal a yaml file into the matcher.
// Yaml can contain either a string or <name, selector> map.
//
// All items provided through the config need to be registered as plugins.
// All config will add to, not replace the exisiting items.
//
// Yaml config example:
//
//	middleware:
//	 - abc
//	 - name: def
//	   selector: "/foo"
func (m *Matcher[T]) UnmarshalYAML(data *yaml.Node) error {
	var a any
	if err := data.Decode(&a); err != nil {
		litter.Dump(a)
		return err
	}

	b, ok := a.([]any)
	if !ok {
		return nil
	}

	return m.unmarshal(b)
}

func (m *Matcher[T]) unmarshal(yee []any) error {
	for _, i := range yee {
		if item, ok := i.(map[string]any); ok {
			selector, ok := slicemap.Get[string](item, "selector")
			if !ok {
				continue
			}

			pattern, ok := slicemap.Get[string](item, "name")
			if !ok {
				continue
			}

			if err := m.AddPlugin(selector, pattern); err != nil {
				return err
			}

			continue
		}

		if global, ok := i.(string); ok {
			if err := m.AddPlugin("/*", global); err != nil {
				return err
			}
		}
	}

	return nil
}
