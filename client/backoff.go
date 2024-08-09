package client

import (
	"context"
	"math"
	"time"
)

const maxAttempts = 13
const backOffSeconds = 120

// BackoffFunc is the type for backoff funcs.
type BackoffFunc func(ctx context.Context, req Request[any, any], attempts int) (time.Duration, error)

// exponentialDo is a function x^e multiplied by a factor of 0.1 second.
// Result is limited to 2 minute.
func exponentialDo(attempts int) time.Duration {
	if attempts > maxAttempts {
		return backOffSeconds * time.Second
	}

	return time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100
}

// BackoffExponential uses expentionalDo to calc the duration to wait.
func BackoffExponential(_ context.Context, _ Request[any, any], attempts int) (time.Duration, error) {
	return exponentialDo(attempts), nil
}
