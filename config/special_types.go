package config

import (
	"encoding/json"
	"errors"
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
