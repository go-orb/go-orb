package config

import "errors"

var (
	// ErrUnknownPlugin happens when there's no config factory for the given plugin.
	ErrUnknownPlugin = errors.New("unknown config given. Did you import the config plugin?")

	// ErrNotExistent happens when a config key is not existent.
	ErrNotExistent = errors.New("no such config key")

	// ErrTypesDontMatch happens when types don't match during Get[T]().
	ErrTypesDontMatch = errors.New("config key requested type and actual type don't match")

	// ErrUnknownScheme happens when you didn't import the plugin for the scheme or the scheme is unknown.
	ErrUnknownScheme = errors.New("unknown config source scheme. Did you register the config source plugin for your scheme?")

	// ErrFileNotFound happens when theres no file.
	ErrFileNotFound = errors.New("file not found")

	// ErrCodecNotFound happens when the required codec is not found.
	ErrCodecNotFound = errors.New("marshaler for codec not found. Did you import the codec plugin for your file type?")
)
