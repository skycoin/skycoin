package iputil

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

// LocalhostIP returns the address for localhost on the machine
func LocalhostIP() (string, error) {
	tt, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			return "", err
		}
		for _, a := range aa {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("No local IP found")
}

// IsLocalhost returns true if addr is a localhost address
// Works for both ipv4 and ipv6 addresses.
func IsLocalhost(addr string) bool {
	return net.ParseIP(addr).IsLoopback() || addr == "localhost"
}

// SplitAddr splits an ip:port string to ip, port.
// Works for both ipv4 and ipv6 addresses.
// If the IP is not specified, returns an error.
func SplitAddr(addr string) (string, uint16, error) {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}

	if ip == "" {
		return "", 0, fmt.Errorf("IP missing from %s", addr)
	}

	port64, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return "", 0, fmt.Errorf("Invalid port in %s", addr)
	}

	return ip, uint16(port64), nil
}
