// Package di provides marker's, these are here against dependency cycles.
package di

// DiFlags is a marker that the config has been loaded from compiled in opts.
type DiFlags struct{}

// DiConfigData is a marker that the config has been loaded from different sources (yaml,json,toml,name it here).
type DiConfigData struct{}

// DiConfig is a list of config URL's.
type DiConfig []string

// DiConfigor is a marker that indicates that the Config loader is available.
type DiConfigor struct{}
