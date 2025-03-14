package client

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/go-orb/go-orb/registry"
)

// NodeMap hold registry nodes grouped by transport.
type NodeMap map[string][]*registry.Node

// SelectorFunc get's executed by client.SelectNode which get it's info's from client.ResolveService.
type SelectorFunc func(
	ctx context.Context,
	service string,
	nodes NodeMap,
	preferredTransports []string,
	anyTransport bool,
) (*registry.Node, error)

// SelectRandomNode selects a random node, it tries' on preferredTransport after another, if anyTransport is true it
// will return transports that are not listet as well.
func SelectRandomNode(
	_ context.Context,
	_ string,
	nodes NodeMap,
	preferredTransports []string,
	anyTransport bool,
) (*registry.Node, error) {
	// try preferredTransports
	for _, pt := range preferredTransports {
		tNodes, ok := nodes[pt]
		if !ok {
			continue
		}

		rInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(tNodes))))
		if err != nil {
			return nil, err
		}

		return tNodes[rInt.Int64()], nil
	}

	// Return random
	if anyTransport {
		aNodes := []*registry.Node{}
		for _, v := range nodes {
			aNodes = append(aNodes, v...)
		}

		rInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(aNodes))))
		if err != nil {
			return nil, err
		}

		return aNodes[rInt.Int64()], nil
	}

	return nil, fmt.Errorf("%w: requested transports was: %s", ErrNoNodeFound, preferredTransports)
}
