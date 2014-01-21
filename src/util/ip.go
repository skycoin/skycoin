// IP Address Tools
package util

import (
    "errors"
    "net"
)

// Returns the localhost ip address
func LocalIP() (net.IP, error) {
    tt, err := net.Interfaces()
    if err != nil {
        return nil, err
    }
    for _, t := range tt {
        aa, err := t.Addrs()
        if err != nil {
            return nil, err
        }
        for _, a := range aa {
            ipnet, ok := a.(*net.IPNet)
            if !ok {
                continue
            }
            v4 := ipnet.IP.To4()
            if v4 == nil || v4[0] == 127 { // loopback address
                continue
            }
            return v4, nil
        }
    }
    return nil, errors.New("cannot find local IP address")
}

// Returns the localhost ip address as a string
func LocalIPString() (string, error) {
    _ip, err := LocalIP()
    var ip string = ""
    if err == nil {
        ip = _ip.String()
    }
    return ip, err
}
