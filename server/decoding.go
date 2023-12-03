package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// HandlerRegistrations

// MarshalJSON no-op.
func (m HandlerRegistrations) MarshalJSON() ([]byte, error) {
	return nil, nil
}

// MarshalYAML no-op.
func (m HandlerRegistrations) MarshalYAML() ([]byte, error) {
	return nil, nil
}

// UnmarshalText handler registrations list..
func (m HandlerRegistrations) UnmarshalText(data []byte) error {
	if strings.Contains(string(data), "[") && strings.Contains(string(data), "]") {
		return m.UnmarshalJSON(data)
	}

	return m.set([]string{string(data)})
}

// UnmarshalJSON handler registrations list.
func (m HandlerRegistrations) UnmarshalJSON(data []byte) error {
	var handlers []string

	if err := json.Unmarshal(data, &handlers); err != nil {
		return err
	}

	return m.set(handlers)
}

// UnmarshalYAML handler registrations list.
func (m HandlerRegistrations) UnmarshalYAML(data *yaml.Node) error {
	var handlers []string

	if err := data.Decode(&handlers); err != nil {
		return err
	}

	return m.set(handlers)
}

func (m HandlerRegistrations) set(handlers []string) error {
	for _, name := range handlers {
		handler, ok := Handlers.Get(name)
		if !ok {
			return fmt.Errorf("handler '%s' not found, did you register it?", name)
		}

		m[name] = handler
	}

	return nil
}
