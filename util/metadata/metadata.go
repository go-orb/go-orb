// Package metadata is a way of defining message headers
package metadata

import (
	"context"
	"strings"
)

type metadataKey struct{}

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

// Ensure returns a context with a Metadata as value,
// it won't overwrite, if metadata exists in the given context.
func Ensure(ctx context.Context) context.Context {
	if _, ok := ctx.Value(metadataKey{}).(Metadata); ok {
		return ctx
	}

	return Metadata{}.To(ctx)
}

// Get returns the value of key.
func (md Metadata) Get(key string) (string, bool) {
	// attempt to get lower case
	val, ok := md[strings.ToLower(key)]

	return val, ok
}

// Set set's the key's value.
func (md Metadata) Set(key, val string) {
	md[strings.ToLower(key)] = val
}

// Delete deletes a key.
func (md Metadata) Delete(key string) {
	delete(md, strings.ToLower(key))
}

// To returns a copy of the given context with metadata as value.
func (md Metadata) To(ctx context.Context) context.Context {
	return context.WithValue(ctx, metadataKey{}, md)
}

// Copy makes a copy of the metadata.
func Copy(md Metadata) Metadata {
	cmd := make(Metadata, len(md))
	for k, v := range md {
		cmd[k] = v
	}

	return cmd
}

// Delete key from metadata.
func Delete(ctx context.Context, k string) context.Context {
	return Set(ctx, k, "")
}

// Set add key with val to metadata.
func Set(ctx context.Context, k, v string) context.Context {
	md, ok := From(ctx)
	if !ok {
		md = make(Metadata)
	}

	if v == "" {
		delete(md, k)
	} else {
		md[k] = v
	}

	return context.WithValue(ctx, metadataKey{}, md)
}

// Get returns a single value from metadata in the context.
func Get(ctx context.Context, key string) (string, bool) {
	md, ok := From(ctx)
	if !ok {
		return "", ok
	}
	// attempt to get as is
	val, ok := md[key]
	if ok {
		return val, ok
	}

	// attempt to get lower case
	val, ok = md[strings.ToLower(key)]

	return val, ok
}

// From returns metadata from the given context.
func From(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return nil, ok
	}

	// lower all values
	newMD := make(Metadata, len(md))
	for k, v := range md {
		newMD[strings.ToLower(k)] = v
	}

	return newMD, ok
}

// Merge merges metadata to existing metadata, overwriting if specified.
func Merge(ctx context.Context, patchMd Metadata, overwrite bool) context.Context {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return patchMd.To(ctx)
	}

	cmd := make(Metadata, len(patchMd)+len(md))
	for k, v := range md {
		cmd[k] = v
	}

	for k, v := range patchMd {
		if _, ok := cmd[k]; ok && !overwrite {
			continue
		}

		if v != "" {
			cmd[k] = v
		} else {
			delete(cmd, k)
		}
	}

	return cmd.To(ctx)
}
