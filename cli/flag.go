//go:build go1.18
// +build go1.18

package cli

import (
	"errors"
)

const (
	FlagTypeNone = iota
	FlagTypeString
	FlagTypeInt
	FlagTypeStringSlice
)

type Flag struct {
	Name    string
	EnvVars []string
	Usage   string

	FlagType int

	DefaultString string
	ValueString   string

	DefaultInt int
	ValueInt   int

	DefaultStringSlice []string
	ValueStringSlice   []string
}

type FlagOption func(*Flag)

func (f *Flag) AsOptions() []FlagOption {
	result := []FlagOption{
		Name(f.Name),
		EnvVars(f.EnvVars...),
		Usage(f.Usage),
	}

	switch f.FlagType {
	case FlagTypeString:
		result = append(result, Default(f.DefaultString))
	case FlagTypeInt:
		result = append(result, Default(f.DefaultInt))
	case FlagTypeStringSlice:
		result = append(result, Default(f.DefaultStringSlice))
	}

	return result
}

func Name(n string) FlagOption {
	return func(o *Flag) {
		o.Name = n
	}
}

func EnvVars(n ...string) FlagOption {
	return func(o *Flag) {
		o.EnvVars = n
	}
}

func Usage(n string) FlagOption {
	return func(o *Flag) {
		o.Usage = n
	}
}

func Default[T any](n T) FlagOption {
	return func(o *Flag) {
		switch def := any(n).(type) {
		case string:
			o.DefaultString = def
			o.FlagType = FlagTypeString
		case int:
			o.DefaultInt = def
			o.FlagType = FlagTypeInt
		case []string:
			o.DefaultStringSlice = def
			o.FlagType = FlagTypeStringSlice
		default:
			o.FlagType = FlagTypeNone
		}
	}
}

func UpdateFlagValue[T any](f *Flag, v T) error {
	switch def := any(v).(type) {
	case string:
		f.ValueString = def
	case []string:
		f.ValueStringSlice = def
	case int:
		f.ValueInt = def
	default:
		return errors.New("failed to update flag")
	}

	return nil
}

func FlagValue[T any](f *Flag, v T) T {
	switch any(v).(type) {
	case string:
		return any(f.ValueString).(T)
	case []string:
		return any(f.ValueStringSlice).(T)
	case int:
		return any(f.ValueInt).(T)
	default:
		var result T
		return result
	}
}

func NewFlag(opts ...FlagOption) (*Flag, error) {
	options := &Flag{
		Name:          "",
		EnvVars:       []string{},
		Usage:         "",
		FlagType:      FlagTypeNone,
		DefaultString: "",
		DefaultInt:    0,
	}

	for _, o := range opts {
		o(options)
	}

	return options, nil
}
