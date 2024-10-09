// Package metadata is a way of defining message headers
package metadata

import (
	"context"
)

// Service is the key for the RPC Service.
const Service = "service"

// Method is the key for the RPC Method.
const Method = "method"

type incomingKey struct{}
type outgoingKey struct{}

// Incoming retrieves incoming metadata from the context.
func Incoming(ctx context.Context) (map[string]string, bool) {
	md, ok := ctx.Value(incomingKey{}).(map[string]string)
	return md, ok
}

// WithIncoming sets metadata as value to the context and returns context as well as the metadata.
func WithIncoming(ctx context.Context) (context.Context, map[string]string) {
	if md, ok := ctx.Value(incomingKey{}).(map[string]string); ok {
		return ctx, md
	}

	md := make(map[string]string)
	ctx = context.WithValue(ctx, incomingKey{}, md)

	return ctx, md
}

// Outgoing retrieves outgoing metadata from the context.
func Outgoing(ctx context.Context) (map[string]string, bool) {
	md, ok := ctx.Value(outgoingKey{}).(map[string]string)
	return md, ok
}

// WithOutgoing sets metadata as value to the context and returns context as well as the metadata.
func WithOutgoing(ctx context.Context) (context.Context, map[string]string) {
	if md, ok := ctx.Value(outgoingKey{}).(map[string]string); ok {
		return ctx, md
	}

	md := make(map[string]string)
	ctx = context.WithValue(ctx, outgoingKey{}, md)

	return ctx, md
}
