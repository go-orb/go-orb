// Package net provides net utilities.
package net

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
)

// HostPort format addr and port suitable for dial.
func HostPort(addr string, port interface{}) string {
	host := addr
	if strings.Count(addr, ":") > 0 {
		host = fmt.Sprintf("[%s]", addr)
	}
	// when port is blank or 0, host is a queue name
	if v, ok := port.(string); ok && v == "" {
		return host
	} else if v, ok := port.(int); ok && v == 0 && net.ParseIP(host) == nil {
		return host
	}

	return fmt.Sprintf("%s:%v", host, port)
}

// Listen will opan a net listener on the specified network and address.
func Listen(network, addr string, tlsConf *tls.Config) (net.Listener, error) {
	if tlsConf != nil {
		return tls.Listen(network, addr, tlsConf)
	}

	return net.Listen(network, addr)
}
