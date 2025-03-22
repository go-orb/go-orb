package config

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"
)

// Duration is a time.Duration that can be parsed from a string or float64.
// https://stackoverflow.com/a/54571600
type Duration time.Duration

// MarshalJSON implements the json.Marshaler interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))

		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}

		*d = Duration(tmp)

		return nil
	default:
		return errors.New("invalid duration")
	}
}

// URL represents a URL in JSON.
type URL struct {
	// URL is the underlying URL.
	*url.URL
}

// String returns the string representation of the JSONURL.
func (j *URL) String() string {
	if j == nil || j.URL == nil {
		return ""
	}

	return j.URL.String()
}

// Copy returns a copy of the JURL.
func (j *URL) Copy() (*URL, error) {
	if j == nil || j.URL == nil {
		return nil, errors.New("JURL is nil")
	}

	return NewURL(j.URL.String())
}

// NewURL creates a new `json` URL.
func NewURL(s string) (*URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	return &URL{URL: u}, nil
}

// UnmarshalJSON unmarshals the JURL from JSON.
func (j *URL) UnmarshalJSON(data []byte) error {
	u, err := url.Parse(strings.Trim(string(data), `"`))
	if err != nil {
		return err
	}

	j.URL = u

	return nil
}

// MarshalJSON marshals the JURL to JSON.
func (j *URL) MarshalJSON() ([]byte, error) {
	if j == nil || j.URL == nil {
		return nil, nil
	}

	return []byte(`"` + j.URL.String() + `"`), nil
}
