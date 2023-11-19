package client

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/util/container"
	"golang.org/x/exp/slices"
)

// SelectorFunc get's executed by client.SelectNode which get it's info's from client.ResolveService.
type SelectorFunc func(
	ctx context.Context,
	service string,
	nodes *container.Map[[]*registry.Node],
	preferredTransports []string,
	anyTransport bool,
) (*registry.Node, error)

// SelectRandomNode selects a random node, it tries' on preferredTransport after another, if anyTransport is true it
// will return transports that are not listet as well.
func SelectRandomNode(
	_ context.Context,
	_ string,
	nodes *container.Map[[]*registry.Node],
	preferredTransports []string,
	anyTransport bool,
) (*registry.Node, error) {
	foundTransports := nodes.Keys()

	for _, pt := range preferredTransports {
		if slices.Contains(foundTransports, pt) {
			tNodes, err := nodes.Get(pt)
			if err != nil {
				// This should never happen.
				return nil, err
			}

			rInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(tNodes))))
			if err != nil {
				return nil, err
			}

			return tNodes[rInt.Int64()], nil
		}
	}

	// Return random
	if anyTransport {
		aNodes := []*registry.Node{}
		for _, v := range nodes.Values() {
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
