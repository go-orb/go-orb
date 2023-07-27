package codecs

import "errors"

var (
	// ErrUnknownMimeType happens when you request a codec for a unknown mime type.
	ErrUnknownMimeType = errors.New("unknown mime-type given")
)
