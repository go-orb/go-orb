package client

import (
	"context"
	"errors"
	"time"

	"github.com/go-orb/go-orb/util/orberrors"
)

// RetryFunc is the type for a retry func.
// note that returning either false or a non-nil error will result in the call not being retried.
type RetryFunc func(ctx context.Context, err error, options *CallOptions) (bool, error)

// RetryAlways always retry on error.
func RetryAlways(_ context.Context, _ error, _ *CallOptions) (bool, error) {
	return true, nil
}

// RetryOnTimeoutError retries a request on a 408 timeout error.
func RetryOnTimeoutError(ctx context.Context, err error, options *CallOptions) (bool, error) {
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
		// Retry on connection error: Service Unavailable
		case 503:
			timeout := time.After(options.DialTimeout)
			select {
			case <-ctx.Done():
				return false, nil
			case <-timeout:
				return true, nil
			}
		default:
			return false, nil
		}
	}

	return false, nil
}

// RetryOnConnectionError retries a request on a 503 connection error.
func RetryOnConnectionError(ctx context.Context, err error, options *CallOptions) (bool, error) {
	if err == nil {
		return false, nil
	}

	var orbe *orberrors.Error

	err = orberrors.From(err)
	if errors.As(err, &orbe) {
		switch orbe.Code {
		// Retry on connection error: Service Unavailable
		case 503:
			timeout := time.After(options.DialTimeout)
			select {
			case <-ctx.Done():
				return false, nil
			case <-timeout:
				return true, nil
			}
		default:
			return false, nil
		}
	}

	return false, nil
}
