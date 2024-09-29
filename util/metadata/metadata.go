// Package metadata is a way of defining message headers
package metadata

import (
	"context"
)

// Service is the key for the RPC Service.
const Service = "Service"

// Method is the key for the RPC Method.
const Method = "Method"

type incomingKey struct{}
type outgoingKey struct{}

// IncomingFrom returns metadata from the given context.
func IncomingFrom(ctx context.Context) (map[string]string, bool) {
	md, ok := ctx.Value(incomingKey{}).(map[string]string)

	return md, ok
}

// EnsureIncoming returns a context with incoming Metadata as a value,
// it won't overwrite, if metadata exists in the given context.
func EnsureIncoming(ctx context.Context) context.Context {
	if _, ok := ctx.Value(incomingKey{}).(map[string]string); ok {
		return ctx
	}

	return context.WithValue(ctx, incomingKey{}, make(map[string]string))
}

// OutgoingFrom returns metadata from the given context.
func OutgoingFrom(ctx context.Context) (map[string]string, bool) {
	md, ok := ctx.Value(outgoingKey{}).(map[string]string)

	return md, ok
}

// EnsureOutgoing returns a context with outgoing Metadata as a value,
// it won't overwrite, if metadata exists in the given context.
func EnsureOutgoing(ctx context.Context) context.Context {
	if _, ok := ctx.Value(outgoingKey{}).(map[string]string); ok {
		return ctx
	}

	return context.WithValue(ctx, outgoingKey{}, make(map[string]string))
}

// WithOutgoing sets metadata as value to the context and returns context as well as the new context.
func WithOutgoing(ctx context.Context) (context.Context, map[string]string) {
	if md, ok := ctx.Value(outgoingKey{}).(map[string]string); ok {
		return ctx, md
	}

	md := make(map[string]string)
	ctx = context.WithValue(ctx, outgoingKey{}, md)

	return ctx, md
}
