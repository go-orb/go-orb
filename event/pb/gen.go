// Package pb contains the cloudevent proto. All messages MUST use this on wire.
package pb

// Generate proto files
//go:generate protoc -I . --go_out=paths=source_relative:. cloudevent.proto
