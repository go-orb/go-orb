package client

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/go-orb/go-orb/registry"
)

// SelectorFunc get's executed by client.SelectNode which get it's info's from client.ResolveService.
type SelectorFunc func(
	ctx context.Context,
	service string,
	nodes []registry.ServiceNode,
) (registry.ServiceNode, error)

// SelectRandomNode selects a random node, it tries' on preferredTransport after another, if anyTransport is true it
// will return transports that are not listet as well.
func SelectRandomNode(
	_ context.Context,
	_ string,
	nodes []registry.ServiceNode,
) (registry.ServiceNode, error) {
	rInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(nodes))))
	if err != nil {
		return registry.ServiceNode{}, err
	}

	return nodes[rInt.Int64()], nil
}
