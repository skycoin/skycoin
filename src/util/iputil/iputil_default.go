//go:build !tinygo

package iputil

import "net"

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
	return "", ErrNoLocalIP
}
