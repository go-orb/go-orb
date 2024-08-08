package registry

import (
	"context"
	"time"
)

// TODO(davincible): investigate all the contexts in the options, are they really needed?
//       maybe there is a better more idomatic way to achieve the same thing.

// RegisterOptions are the options used to register services.
type RegisterOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// WatchOptions are the options used by the registry watcher.
type WatchOptions struct {
	// Specify a service to watch
	// If blank, the watch is for all services
	Service string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// RegisterOption is functional option type for the register config.
type RegisterOption func(*RegisterOptions)

// WatchOption is functional option type for the watch config.
type WatchOption func(*WatchOptions)

// DeregisterOption is functional option type for the deregister config.
type DeregisterOption func(*DeregisterOptions)

// GetOption is functional option type for the get config.
type GetOption func(*GetOptions)

// ListOption is functional option type for the list config.
type ListOption func(*ListOptions)

// DeregisterOptions are the options used to deregister services.
type DeregisterOptions struct {
	Context context.Context
}

// GetOptions are the options used to fetch a service.
type GetOptions struct {
	Context context.Context
}

// ListOptions are the options used to list services.
type ListOptions struct {
	Context context.Context
}

// RegisterTTL sets the TTL for service registration.
func RegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

// RegisterContext sets the context that is used when registering a service.
func RegisterContext(ctx context.Context) RegisterOption {
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

// WatchService sets a service name to watch.
func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}

// WatchContext sets a context that is used to watch.
func WatchContext(ctx context.Context) WatchOption {
	return func(o *WatchOptions) {
		o.Context = ctx
	}
}

// DeregisterContext is the context used to deregister a service.
func DeregisterContext(ctx context.Context) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Context = ctx
	}
}

// GetContext is the context used when fetching a service.
func GetContext(ctx context.Context) GetOption {
	return func(o *GetOptions) {
		o.Context = ctx
	}
}

// ListContext is the context used when listing a service.
func ListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}
