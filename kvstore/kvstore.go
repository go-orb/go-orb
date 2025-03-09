// Package kvstore is an interface for distributed key-value data storage.
package kvstore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
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
	// Leave database and table empty to use the defaults.
	Get(key, database, table string, opts ...GetOption) ([]Record, error)

	// Set takes a key, database, table and data, and optional SetOptions.
	// Leave database and table empty to use the defaults.
	Set(key, database, table string, data []byte, opts ...SetOption) error

	// Purge takes a key, database and table and purges it.
	// Leave database and table empty to use the defaults.
	Purge(key, database, table string) error

	// Keys returns any keys that match, or an empty list with no error if none matched.
	// Leave database and table empty to use the defaults.
	Keys(database, table string, opts ...KeysOption) ([]string, error)

	// DropTable drops the table.
	// Leave database and table empty to use the defaults.
	DropTable(database, table string) error

	// DropDatabase drops the database.
	// Leave database empty to use the default.
	DropDatabase(database string) error
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

// Provide provides a new KVStore.
func Provide(
	ctx context.Context,
	name types.ServiceName,
	configs types.ConfigData,
	components *types.Components,
	logger log.Logger,
	registry registry.Type,
	client client.Type,
	opts ...Option,
) (Type, error) {
	cfg := NewConfig(opts...)

	sections := append(types.SplitServiceName(name), DefaultConfigSection)
	if err := config.Parse(sections, configs, &cfg); err != nil {
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
	cLogger, err := logger.WithConfig(sections, configs)
	if err != nil {
		return Type{}, err
	}

	cLogger = cLogger.With(slog.String("component", ComponentType), slog.String("plugin", cfg.Plugin))

	instance, err := provider(ctx, name, configs, cLogger, registry, client, opts...)
	if err != nil {
		return Type{}, err
	}

	// Register the registry as a component.
	err = components.Add(instance, types.PriorityKVStore)
	if err != nil {
		logger.Warn("while registering kvstore as a component", "error", err)
	}

	return instance, nil
}

// ProvideNoOpts provides a new KVStore without options.
func ProvideNoOpts(
	ctx context.Context,
	name types.ServiceName,
	configs types.ConfigData,
	components *types.Components,
	logger log.Logger,
	registry registry.Type,
	client client.Type,
) (Type, error) {
	return Provide(ctx, name, configs, components, logger, registry, client)
}
