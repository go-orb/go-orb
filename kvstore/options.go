package kvstore

import (
	"time"
)

// GetOptions configures an individual Get operation.
type GetOptions struct {
	// Prefix returns all records that are prefixed with key
	Prefix bool
	// Suffix returns all records that have the suffix key
	Suffix bool
	// Limit limits the number of returned records
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// NewGetOptions creates a new GetOptions with the provided options applied.
func NewGetOptions(opts ...GetOption) GetOptions {
	var options GetOptions
	for _, o := range opts {
		o(&options)
	}

	return options
}

// GetOption sets values in GetOptions.
type GetOption func(g *GetOptions)

// GetPrefix returns all records that are prefixed with key.
func GetPrefix() GetOption {
	return func(g *GetOptions) {
		g.Prefix = true
	}
}

// GetSuffix returns all records that have the suffix key.
func GetSuffix() GetOption {
	return func(g *GetOptions) {
		g.Suffix = true
	}
}

// GetLimit limits the number of responses to l.
func GetLimit(l uint) GetOption {
	return func(g *GetOptions) {
		g.Limit = l
	}
}

// GetOffset starts returning responses from o. Use in conjunction with Limit for pagination.
func GetOffset(o uint) GetOption {
	return func(g *GetOptions) {
		g.Offset = o
	}
}

// SetOptions configures an individual Set operation
// If Expiry and TTL are set TTL takes precedence.
type SetOptions struct {
	// Expiry is the time the record expires
	Expiry time.Time
	// TTL is the time until the record expires
	TTL time.Duration
}

// NewSetOptions creates a new SetOptions with the provided options applied.
func NewSetOptions(opts ...SetOption) SetOptions {
	var options SetOptions
	for _, o := range opts {
		o(&options)
	}

	return options
}

// SetOption sets values in SetOptions.
type SetOption func(s *SetOptions)

// SetExpiry is the time the record expires.
func SetExpiry(t time.Time) SetOption {
	return func(s *SetOptions) {
		s.Expiry = t
	}
}

// SetTTL is the time until the record expires.
func SetTTL(d time.Duration) SetOption {
	return func(s *SetOptions) {
		s.TTL = d
	}
}

// KeysOptions configures an individual Keys/List operation.
type KeysOptions struct {
	// Prefix returns all keys that are prefixed with key
	Prefix string
	// Suffix returns all keys that end with key
	Suffix string
	// Limit limits the number of returned keys
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// NewKeysOptions creates a new KeysOptions with the provided options applied.
func NewKeysOptions(opts ...KeysOption) KeysOptions {
	var options KeysOptions
	for _, o := range opts {
		o(&options)
	}

	return options
}

// KeysOption sets values in KeysOptions.
type KeysOption func(k *KeysOptions)

// KeysPrefix returns all keys that are prefixed with key.
func KeysPrefix(p string) KeysOption {
	return func(k *KeysOptions) {
		k.Prefix = p
	}
}

// KeysSuffix returns all keys that have the suffix key.
func KeysSuffix(s string) KeysOption {
	return func(k *KeysOptions) {
		k.Suffix = s
	}
}

// KeysLimit limits the number of responses to l.
func KeysLimit(l uint) KeysOption {
	return func(k *KeysOptions) {
		k.Limit = l
	}
}

// KeysOffset starts returning responses from o. Use in conjunction with Limit for pagination.
func KeysOffset(o uint) KeysOption {
	return func(k *KeysOptions) {
		k.Offset = o
	}
}

// ReadOptions configures an individual Read operation.
type ReadOptions struct {
	Database, Table string
	// Prefix returns all keys that are prefixed with key
	Prefix string
	// Suffix returns all keys that end with key
	Suffix string
	// Limit limits the number of returned keys
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// NewReadOptions creates a new ReadOptions with the provided options applied.
func NewReadOptions(opts ...ReadOption) ReadOptions {
	var options ReadOptions
	for _, o := range opts {
		o(&options)
	}

	return options
}

// ReadOption sets values in ReadOptions.
type ReadOption func(r *ReadOptions)

// ReadFrom the database and table.
func ReadFrom(database, table string) ReadOption {
	return func(r *ReadOptions) {
		r.Database = database
		r.Table = table
	}
}

// ReadPrefix returns all records that are prefixed with key.
func ReadPrefix(p string) ReadOption {
	return func(r *ReadOptions) {
		r.Prefix = p
	}
}

// ReadSuffix returns all records that have the suffix key.
func ReadSuffix(s string) ReadOption {
	return func(r *ReadOptions) {
		r.Suffix = s
	}
}

// ReadLimit limits the number of responses to l.
func ReadLimit(l uint) ReadOption {
	return func(r *ReadOptions) {
		r.Limit = l
	}
}

// ReadOffset starts returning responses from o. Use in conjunction with Limit for pagination.
func ReadOffset(o uint) ReadOption {
	return func(r *ReadOptions) {
		r.Offset = o
	}
}

// WriteOptions configures an individual Write operation
// If Expiry and TTL are set TTL takes precedence.
type WriteOptions struct {
	// Expiry is the time the record expires
	Expiry          time.Time
	Database, Table string
	// TTL is the time until the record expires
	TTL time.Duration
}

// NewWriteOptions creates a new WriteOptions with the provided options applied.
func NewWriteOptions(opts ...WriteOption) WriteOptions {
	var options WriteOptions
	for _, o := range opts {
		o(&options)
	}

	return options
}

// WriteOption sets values in WriteOptions.
type WriteOption func(w *WriteOptions)

// WriteTo the database and table.
func WriteTo(database, table string) WriteOption {
	return func(w *WriteOptions) {
		w.Database = database
		w.Table = table
	}
}

// WriteExpiry is the time the record expires.
func WriteExpiry(t time.Time) WriteOption {
	return func(w *WriteOptions) {
		w.Expiry = t
	}
}

// WriteTTL is the time the record expires.
func WriteTTL(d time.Duration) WriteOption {
	return func(w *WriteOptions) {
		w.TTL = d
	}
}

// DeleteOptions configures an individual Delete operation.
type DeleteOptions struct {
	Database, Table string
}

// NewDeleteOptions creates a new DeleteOptions with the provided options applied.
func NewDeleteOptions(opts ...DeleteOption) DeleteOptions {
	var options DeleteOptions
	for _, o := range opts {
		o(&options)
	}

	return options
}

// DeleteOption sets values in DeleteOptions.
type DeleteOption func(d *DeleteOptions)

// DeleteFrom the database and table.
func DeleteFrom(database, table string) DeleteOption {
	return func(d *DeleteOptions) {
		d.Database = database
		d.Table = table
	}
}

// ListOptions configures an individual List operation.
type ListOptions struct {
	// List from the following
	Database, Table string
	// Prefix returns all keys that are prefixed with key
	Prefix string
	// Suffix returns all keys that end with key
	Suffix string
	// Limit limits the number of returned keys
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset uint
}

// NewListOptions creates a new ListOptions with the provided options applied.
func NewListOptions(opts ...ListOption) ListOptions {
	var options ListOptions
	for _, o := range opts {
		o(&options)
	}

	return options
}

// ListOption sets values in ListOptions.
type ListOption func(l *ListOptions)

// ListFrom the database and table.
func ListFrom(database, table string) ListOption {
	return func(l *ListOptions) {
		l.Database = database
		l.Table = table
	}
}

// ListPrefix returns all keys that are prefixed with key.
func ListPrefix(p string) ListOption {
	return func(l *ListOptions) {
		l.Prefix = p
	}
}

// ListSuffix returns all keys that end with key.
func ListSuffix(s string) ListOption {
	return func(l *ListOptions) {
		l.Suffix = s
	}
}

// ListLimit limits the number of returned keys to l.
func ListLimit(l uint) ListOption {
	return func(lo *ListOptions) {
		lo.Limit = l
	}
}

// ListOffset starts returning responses from o. Use in conjunction with Limit for pagination.
func ListOffset(o uint) ListOption {
	return func(l *ListOptions) {
		l.Offset = o
	}
}
