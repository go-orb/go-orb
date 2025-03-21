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
	metadata map[string]string,
) (*registry.Node, error)

// SelectRandomNode selects a random node, it tries' on preferredTransport after another, if anyTransport is true it
// will return transports that are not listet as well.
func SelectRandomNode(
	_ context.Context,
	_ string,
	nodes NodeMap,
	preferredTransports []string,
	anyTransport bool,
	metadata map[string]string,
) (*registry.Node, error) {
	filteredNodes, err := filterNodesByMetadata(nodes, metadata)
	if err != nil {
		return nil, err
	}

	for _, pt := range preferredTransports {
		tNodes, ok := filteredNodes[pt]
		if !ok {
			continue
		}

		return getRandomNodeFromSlice(tNodes)
	}

	if anyTransport {
		var allNodes []*registry.Node

		for _, v := range filteredNodes {
			allNodes = append(allNodes, v...)
		}

		if len(allNodes) > 0 {
			return getRandomNodeFromSlice(allNodes)
		}
	}

	return nil, fmt.Errorf("%w: requested transports was: %s", ErrNoNodeFound, preferredTransports)
}

// filterNodesByMetadata filters nodes based on metadata.
func filterNodesByMetadata(nodes NodeMap, metadata map[string]string) (NodeMap, error) {
	if len(metadata) == 0 {
		return nodes, nil
	}

	filteredNodes := NodeMap{}

	for transport, transportNodes := range nodes {
		var matchingNodes []*registry.Node

		for _, node := range transportNodes {
			matches := true

			for k, v := range metadata {
				if nodeVal, exists := node.Metadata[k]; !exists || nodeVal != v {
					matches = false
					break
				}
			}

			if matches {
				matchingNodes = append(matchingNodes, node)
			}
		}

		if len(matchingNodes) > 0 {
			filteredNodes[transport] = matchingNodes
		}
	}

	if len(filteredNodes) == 0 {
		return nil, fmt.Errorf("%w: no nodes matched the specified metadata", ErrNoNodeFound)
	}

	return filteredNodes, nil
}

// getRandomNodeFromSlice randomly selects a node from a slice of nodes.
func getRandomNodeFromSlice(nodes []*registry.Node) (*registry.Node, error) {
	rInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(nodes))))
	if err != nil {
		return nil, err
	}

	return nodes[rInt.Int64()], nil
}
