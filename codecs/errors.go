package codecs

import "errors"

//nolint:gochecknoglobals
var (
	// ErrUknownMimeType happens when you request a codec for a unknown mime type.
	ErrUnknownMimeType = errors.New("unknown mime-type given")
)
