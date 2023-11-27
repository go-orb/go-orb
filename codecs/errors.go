package codecs

import "errors"

var (
	// ErrUnknownMimeType happens when you request a codec for a unknown mime type.
	ErrUnknownMimeType = errors.New("unknown mime type given")

	// ErrUnknownValueType happens when you give a golang type to a codec that doesn't understand it.
	ErrUnknownValueType = errors.New("unknown golang type given")
)
