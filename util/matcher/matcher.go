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

// Matcher is a map with regular expressions as keys. It is used to add
// middleware on a path based level.
type Matcher[T any] struct {
	globals   []T
	plugins   *container.Plugins[T]
	selectors map[selectorKey][]T
}

// NewMatcher creates a new matcher object.
func NewMatcher[T any](plugins *container.Plugins[T]) Matcher[T] {
	return Matcher[T]{
		plugins:   plugins,
		selectors: make(map[selectorKey][]T),
	}
}

// Use will use the elements provided on all paths.
func (m *Matcher[T]) Use(ms ...T) {
	m.globals = append(m.globals, ms...)
}

// AddPlugin will add plugin item, with a selector.
func (m *Matcher[T]) AddPlugin(selector, item string) error {
	i, err := m.plugins.Get(item)
	if err != nil {
		return fmt.Errorf("plugin not found '%s'", item)
	}

	m.Add(selector, i)

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
func (m *Matcher[T]) Add(selector string, items ...T) {
	switch selector {
	case "/*":
		fallthrough
	case "*":
		fallthrough
	case ".*":
		fallthrough
	case "^.*":
		m.Use(items...)
		return
	}

	if strings.HasSuffix(selector, "/*") {
		selector = strings.TrimSuffix(selector, "/*")
		selector += "/.*"
	}

	// Check if we already have a similar selector
	for key := range m.selectors {
		if key.selector == selector {
			m.selectors[key] = append(m.selectors[key], items...)
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

	m.selectors[s] = items
}

// Match will fetch the list of items that match against a path.
func (m *Matcher[T]) Match(operation string) []T {
	ms := make([]T, 0, len(m.globals))

	ms = append(ms, m.globals...)

	if m.selectors == nil {
		return ms
	}

	for selector, val := range m.selectors {
		if selector.re.MatchString(operation) {
			ms = append(ms, val...)
		}
	}

	return ms
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
