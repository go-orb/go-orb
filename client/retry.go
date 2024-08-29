package client

import (
	"context"
	"errors"

	"github.com/go-orb/go-orb/util/orberrors"
)

// RetryFunc is the type for a retry func.
// note that returning either false or a non-nil error will result in the call not being retried.
type RetryFunc func(ctx context.Context, req Request[any, any], retryCount int, err error) (bool, error)

// RetryAlways always retry on error.
func RetryAlways(_ context.Context, _ Request[any, any], _ int, _ error) (bool, error) {
	return true, nil
}

// RetryOnTimeoutError retries a request on a 408 timeout error.
func RetryOnTimeoutError(_ context.Context, _ Request[any, any], _ int, err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	var orbe *orberrors.Error

	err = orberrors.From(err)
	if errors.As(err, &orbe) {
		switch orbe.Code {
		// Retry on timeout, not on 500 internal server error, as that is a business
		// logic error that should be handled by the user.
		case 408:
			return true, nil
		default:
			return false, nil
		}
	}

	return false, nil
}
