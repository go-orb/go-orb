package registry

import "time"

// Watcher is an interface that returns updates
// about services within the registry.
type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop() error
}

// Result is returned by a call to Next on
// the watcher. Actions can be create, update, delete.
type Result struct {
	Action  string   `json:"action"`
	Service *Service `json:"service"`
}

// EventType defines registry event type.
type EventType int

const (
	// Create is emitted when a new service is registered.
	Create EventType = iota
	// Delete is emitted when an existing service is deregsitered.
	Delete
	// Update is emitted when an existing servicec is updated.
	Update
)

// String returns human readable event type.
func (t EventType) String() string {
	switch t {
	case Create:
		return "create"
	case Delete:
		return "delete"
	case Update:
		return "update"
	default:
		return "unknown"
	}
}

// Event is registry event.
type Event struct {
	// ID is registry id
	ID string
	// Type defines type of event
	Type EventType
	// Timestamp is event timestamp
	Timestamp time.Time
	// Service is registry service
	Service *Service
}
