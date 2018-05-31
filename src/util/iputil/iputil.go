package iputil

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
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
func IsLocalhost(addr string) bool {
	return net.ParseIP(addr).IsLoopback() || addr == "localhost"
}

// SplitAddr splits an ip:port string to ip, port
func SplitAddr(addr string) (string, uint16, error) {
	pts := strings.Split(addr, ":")
	if len(pts) != 2 {
		return pts[0], 0, fmt.Errorf("Invalid addr %s", addr)
	}
	port64, err := strconv.ParseUint(pts[1], 10, 16)
	if err != nil {
		return pts[0], 0, fmt.Errorf("Invalid port in %s", addr)
	}
	return pts[0], uint16(port64), nil
}
