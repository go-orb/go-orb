// Package addr provides functions to retrieve local IP addresses from device interfaces.
package addr

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrIPNotFound no IP address found, and explicit IP not provided.
	ErrIPNotFound = errors.New("no IP address found, and explicit IP not provided")
	// ErrNoAddress is returned when no address is provided.
	ErrNoAddress = errors.New("no adddress provided")
	// ErrPortInvalid is returned when the provided port is below 0.
	ErrPortInvalid = errors.New("invalid port provided, must be between 0 and 65535")
	// ErrInvalidIP is returned an invalid IP is provided.
	ErrInvalidIP = errors.New("invalid IP provided")
)

var ipRe = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)

// GetAddress will validate the address if one is provided, otherwise it will
// return an interface and port to listen on.
//
// If you want to listen on all interfaces, you have to explicityly set
// '0.0.0.0:<port>'.
//
// If no IP address is provides, as ':8080' the first private interface will
// be selected to listen on. If none are available, a public IP is selected.
func GetAddress(addr string) (string, error) {
	var host string

	port := "0"

	if len(addr) > 0 {
		if err := ValidateAddress(addr); err != nil {
			return addr, err
		}

		host, port, _ = net.SplitHostPort(addr) //nolint:errcheck
	}

	host, err := Extract(host)
	if err != nil {
		return addr, err
	}

	return host + ":" + port, nil
}

// ValidateAddress will do basic validation on an address string.
//
// Address example:
//   - :8080
//   - 192.168.1.1:8080
//   - [2001:db8::1]:8080
func ValidateAddress(address string) error {
	if len(address) == 0 {
		return ErrNoAddress
	}

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("split host and port from address: %w", err)
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	if p < 0 || p > 65535 {
		return ErrPortInvalid
	}

	// No host is a valid host, to listen on all interfaces.
	if len(host) == 0 {
		return nil
	}

	if ipRe.MatchString(host) && net.ParseIP(host) == nil {
		return ErrInvalidIP
	}

	return nil
}

// ParsePort will take the port from an address and return it as int.
func ParsePort(address string) (int, error) {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return 0, fmt.Errorf("split host and port from address: %w", err)
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}

	return p, nil
}

// Extract returns a valid IP address. If the address provided is a valid
// address, it will be returned directly. Otherwise, the available interfaces
// will be iterated over to find an IP address, preferably private.
func Extract(addr string) (string, error) {
	// if addr is already specified then it's directly returned
	if len(addr) > 0 {
		return addr, nil
	}

	var (
		addrs   []net.Addr
		loAddrs []net.Addr
	)

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("get interfaces: %w", err)
	}

	for _, iface := range ifaces {
		ifaceAddrs, err := iface.Addrs()
		if err != nil {
			// ignore error, interface can disappear from system
			continue
		}

		if iface.Flags&net.FlagLoopback != 0 {
			loAddrs = append(loAddrs, ifaceAddrs...)
			continue
		}

		addrs = append(addrs, ifaceAddrs...)
	}

	// Add loopback addresses to the end of the list
	addrs = append(addrs, loAddrs...)

	// Try to find private IP in list, public IP otherwise
	ip, err := findIP(addrs)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}

// IPs returns all available interface IP addresses.
func IPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ipAddrs []string

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			ipAddrs = append(ipAddrs, ip.String())
		}
	}

	return ipAddrs
}

// findIP will return the first private IP available in the list.
// If no private IP is available it will return the first public IP, if present.
// If no public IP is available, it will return the first loopback IP, if present.
func findIP(addresses []net.Addr) (net.IP, error) {
	var (
		publicIP net.IP
		localIP  net.IP
	)

	for _, rawAddr := range addresses {
		var myIP net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			myIP = addr.IP
		case *net.IPNet:
			myIP = addr.IP
		default:
			continue
		}

		if myIP.IsLoopback() {
			if localIP == nil {
				localIP = myIP
			}

			continue
		}

		if !myIP.IsPrivate() {
			if publicIP == nil {
				publicIP = myIP
			}

			continue
		}

		// Return private IP if available
		return myIP, nil
	}

	// Return public or virtual IP
	if len(publicIP) > 0 {
		return publicIP, nil
	}

	// Return local IP
	if len(localIP) > 0 {
		return localIP, nil
	}

	return nil, ErrIPNotFound
}

// IsLocal checks whether an IP belongs to one of the device's interfaces.
func IsLocal(addr string) bool {
	// Extract the host
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		addr = host
	}

	if addr == "localhost" {
		return true
	}

	// Check against all local ips
	for _, ip := range IPs() {
		if addr == ip {
			return true
		}
	}

	return false
}

// HostPort formats addr and port suitable for dial.
func HostPort(addr string, port any) string {
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
