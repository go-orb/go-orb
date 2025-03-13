// Package client provides an interface and helpers for go-orb clients.
package client

import (
	"context"
	"errors"
)

// StreamIface is the interface for handling streaming operations.
type StreamIface[TReq any, TResp any] interface {
	// Send sends a message to the stream.
	Send(msg TReq) error

	// Recv receives a message from the stream.
	Recv(msg TResp) error

	// Close closes the stream.
	Close() error

	// CloseSend closes the send direction of the stream but allows receiving responses.
	CloseSend() error

	// Context returns the context for the stream.
	Context() context.Context
}

// Stream is a helper function to create a new stream with the given service and endpoint.
func Stream[TReq any, TResp any](
	ctx context.Context,
	c Client,
	service, endpoint string,
	opts ...CallOption,
) (StreamIface[TReq, TResp], error) {
	// Call the underlying Stream method
	stream, err := c.Stream(ctx, service, endpoint, opts...)
	if err != nil {
		return nil, err
	}

	// Create an adapter to convert between StreamIface[any, any] and StreamIface[TReq, TResp]
	return NewStreamAdapter[TReq, TResp](stream, service, endpoint, c), nil
}

// StreamAdapter is an adapter that maps between StreamIface with generic types (any, any)
// and StreamIface with specific types (TReq, TResp).
type StreamAdapter[TReq any, TResp any] struct {
	service  string
	endpoint string

	client Client

	stream StreamIface[any, any]
}

// NewStreamAdapter creates a new stream adapter.
func NewStreamAdapter[TReq any, TResp any](
	stream StreamIface[any, any],
	service string,
	endpoint string,
	client Client,
) StreamIface[TReq, TResp] {
	return &StreamAdapter[TReq, TResp]{
		stream:   stream,
		service:  service,
		endpoint: endpoint,
		client:   client,
	}
}

// Service returns the Service from the request.
func (s *StreamAdapter[TReq, TResp]) Service() string {
	return s.service
}

// Endpoint returns the Endpoint from the request.
func (s *StreamAdapter[TReq, TResp]) Endpoint() string {
	return s.endpoint
}

// Send sends a message to the stream.
func (s *StreamAdapter[TReq, TResp]) Send(msg TReq) error {
	return s.stream.Send(msg)
}

// Recv receives a message from the stream.
func (s *StreamAdapter[TReq, TResp]) Recv(msg TResp) error {
	return s.stream.Recv(msg)
}

// Close closes the stream.
func (s *StreamAdapter[TReq, TResp]) Close() error {
	return s.stream.Close()
}

// CloseSend closes the send side of the stream.
func (s *StreamAdapter[TReq, TResp]) CloseSend() error {
	return s.stream.CloseSend()
}

// Context returns the context for the stream.
func (s *StreamAdapter[TReq, TResp]) Context() context.Context {
	return s.stream.Context()
}

// ErrStreamNotSupported is returned when the client does not support streaming.
var ErrStreamNotSupported = errors.New("client does not support streaming")
