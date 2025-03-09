package event

import (
	"errors"
)

var (
	// ErrMissingTopic happens whenever the user doesnt give a topic.
	ErrMissingTopic = errors.New("missing topic")

	// ErrEncodingMessage is returned from publish if there was an error encoding the message option.
	ErrEncodingMessage = errors.New("encoding message")
)
