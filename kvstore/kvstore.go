// Package kvstore is an interface for distributed key-value data storage.
package kvstore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
)

// ComponentType is the components name.
const ComponentType = "kvstore"

var (
	// ErrNotFound is returned when a key is not found.
	ErrNotFound = errors.New("not found")
	// ErrDatabaseNotFound is returned when a database is not found.
	ErrDatabaseNotFound = errors.New("database not found")
	// ErrTableNotFound is returned when a table is not found.
	ErrTableNotFound = errors.New("table not found")
)

// Deprecated is an interface for deprecated methods in KVStore.
// This is used to provide backwards compatibility for old code, it will be removed on go-orb v1.0.0.
type Deprecated interface {
	// Read takes a single key and optional ReadOptions. It returns matching []*Record or an error.
	// If no Record is found, ErrNotFound is returned.
	// Deprecated: use Get instead.
	Read(key string, opts ...ReadOption) ([]*Record, error)
	// Write takes a single key and value, and optional WriteOptions.
	// Deprecated: use Set instead.
	Write(r *Record, opts ...WriteOption) error
	// Delete purges the record with the corresponding key from the store.
	// Deprecated: use Purge instead.
	Delete(key string, opts ...DeleteOption) error
	// List returns any keys that match, or an empty list with no error if none matched.
	// Deprecated: use Keys instead.
	List(opts ...ListOption) ([]string, error)
}

// KVStore is an interface for distributed key-value data storage.
type KVStore interface {
	types.Component

	Deprecated

	// Get takes a key, database, table and optional GetOptions. It returns the Record or an error.
	// Leave database and/or table empty to use the defaults.
	Get(ctx context.Context, key, database, table string, opts ...GetOption) ([]Record, error)

	// Set takes a key, database, table and data, and optional SetOptions.
	// Leave database and/or table empty to use the defaults.
	Set(ctx context.Context, key, database, table string, data []byte, opts ...SetOption) error

	// Purge takes a key, database and table and purges it.
	// Leave database and/or table empty to use the defaults.
	Purge(ctx context.Context, key, database, table string) error

	// Keys returns any keys that match, or an empty list with no error if none matched.
	// Leave database and/or table empty to use the defaults.
	Keys(ctx context.Context, database, table string, opts ...KeysOption) ([]string, error)

	// DropTable drops the table.
	// Leave database and/or table empty to use the defaults.
	DropTable(ctx context.Context, database, table string) error

	// DropDatabase drops the database.
	// Leave database empty to use the default.
	DropDatabase(ctx context.Context, database string) error
}

// WatchOp represents the type of Watch operation (Update, Delete). It is a
// part of WatchUpdate.
type WatchOp uint8

// Available WatchOp values.
const (
	// WatchOpCreate is a create operation.
	WatchOpCreate WatchOp = iota
	// WatchOpUpdate is an update operation.
	WatchOpUpdate
	// WatchOpDelete is a delete operation.
	WatchOpDelete
)

// WatchEvent is a change to a key-value store.
type WatchEvent struct {
	Record

	// Operation is the type of operation that occurred.
	Operation WatchOp
}

// A Watcher is a component that can watch for changes to a key-value store.
// For example a registry component can watch for changes to a database.
type Watcher interface {
	// Watch starts a watcher for the given database and table.
	// Returns a channel of WatchUpdate and a function to stop the watcher.
	// If an error occurs, it is returned.
	Watch(ctx context.Context, database, table string, opts ...WatchOption) (<-chan WatchEvent, func() error, error)
}

// Type is the kvstore type it is returned when you use Provide
// which selects a kvstore to use based on the plugin configuration.
type Type struct {
	KVStore
}

// Record is an item stored or retrieved from a Store.
type Record struct {
	// The key to store the record.
	Key string `json:"key"`
	// The value within the record.
	Value []byte `json:"value"`
	// Time when the record expires a nil value means no expiry.
	Expiry *time.Time `json:"expiry,omitempty"`
}

// New creates a new kvstore instance with the implementation from cfg.Plugin.
func New(
	configData map[string]any,
	logger log.Logger,
	opts ...Option,
) (Type, error) {
	cfg := NewConfig(opts...)

	if err := config.Parse(nil, DefaultConfigSection, configData, &cfg); err != nil && !errors.Is(err, config.ErrNoSuchKey) {
		return Type{}, err
	}

	if cfg.Plugin == "" {
		logger.Warn("empty kvstore plugin, using the default", "default", DefaultKVStore)
		cfg.Plugin = DefaultKVStore
	}

	logger.Debug("KVStore", "plugin", cfg.Plugin)

	provider, ok := plugins.Get(cfg.Plugin)
	if !ok {
		return Type{}, fmt.Errorf("KVStore plugin '%s' not found, did you import it?", cfg.Plugin)
	}

	// Configure the logger.
	cLogger, err := logger.WithConfig([]string{DefaultConfigSection}, configData)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	instance, err := provider(configData, cLogger, opts...)
	if err != nil {
		return Type{}, err
	}

	return instance, nil
}

// Provide provides a new KVStore.
func Provide(
	svcCtx *cli.ServiceContextWithConfig,
	components *types.Components,
	logger log.Logger,
	opts ...Option,
) (Type, error) {
	instance, err := New(svcCtx.Config(), logger, opts...)
	if err != nil {
		return Type{}, err
	}

	// Register the kvstore as a component.
	err = components.Add(instance, types.PriorityKVStore)
	if err != nil {
		logger.Warn("while registering kvstore as a component", "error", err)
	}

	return instance, nil
}

// ProvideNoOpts provides a new KVStore without options.
func ProvideNoOpts(
	svcCtx *cli.ServiceContextWithConfig,
	components *types.Components,
	logger log.Logger,
) (Type, error) {
	return Provide(svcCtx, components, logger)
}
