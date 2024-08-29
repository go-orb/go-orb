package event

import (
	"net/http"

	"github.com/go-orb/go-orb/util/orberrors"
)

var (
	// ErrMissingTopic happens whenever the user doesnt give a topic.
	ErrMissingTopic = orberrors.New(http.StatusInternalServerError, "missing topic")
)
