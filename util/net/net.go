// Package net provides net utilities.
package net

import (
	"crypto/tls"
	"net"
)

// Listen will opan a net listener on the specified network and address.
func Listen(network, addr string, tlsConf *tls.Config) (net.Listener, error) {
	if tlsConf != nil {
		return tls.Listen(network, addr, tlsConf)
	}

	return net.Listen(network, addr)
}
