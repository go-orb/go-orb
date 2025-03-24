package registry

// WatchOptions are the options used by the registry watcher.
type WatchOptions struct {
	// Specify a service to watch
	// If blank, the watch is for all services
	Service string
}

// WatchOption is functional option type for the watch config.
type WatchOption func(*WatchOptions)

// WatchService sets a service name to watch.
func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}
