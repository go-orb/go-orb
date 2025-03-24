package client

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/util/container"
)

// MemoryServer is the interface that a memory server has to implement to be accepted by the client.
type MemoryServer interface {
	// Request does the actual call.
	Request(ctx context.Context, infos RequestInfos, req any, result any, opts *CallOptions) error

	// Stream creates a streaming client to the specified service endpoint.
	Stream(ctx context.Context, infos RequestInfos, opts *CallOptions) (StreamIface[any, any], error)
}

//nolint:gochecknoglobals
var memoryServers = container.NewSafeMap[string, MemoryServer]()

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
