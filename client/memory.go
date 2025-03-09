package client

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/util/container"
)

// MemoryServer is the interface that a memory server has to implement to be accepted by the client.
type MemoryServer interface {
	// Request is the same as Request but without encoding.
	Request(ctx context.Context, req *Req[any, any], result any, opts *CallOptions) error
}

//nolint:gochecknoglobals
var memoryServers = container.NewMap[string, MemoryServer]()

// RegisterMemoryServer registers a memory server for a service.
func RegisterMemoryServer(service string, server MemoryServer) {
	memoryServers.Set(service, server)
}

// UnregisterMemoryServer unregisters a memory server for a service.
func UnregisterMemoryServer(service string) {
	memoryServers.Del(service)
}

// ResolveMemoryServer resolves a memory server for a service.
func ResolveMemoryServer(service string) (MemoryServer, error) {
	server, ok := memoryServers.Get(service)
	if !ok {
		return nil, fmt.Errorf("memory server not found for service %s", service)
	}

	return server, nil
}
