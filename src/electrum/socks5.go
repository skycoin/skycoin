package electrum

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

// SOCKS5Dialer returns a DialFunc that reaches addr through a SOCKS5 proxy with
// no authentication. The host is sent by name, so the proxy resolves it. forward
// opens the connection to the proxy; nil means net.Dialer. See RFC 1928.
func SOCKS5Dialer(proxyAddr string, timeout time.Duration, forward func(ctx context.Context, network, addr string) (net.Conn, error)) DialFunc {
	if forward == nil {
		forward = (&net.Dialer{}).DialContext
	}
	return func(network, addr string) (net.Conn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		conn, err := forward(ctx, network, proxyAddr)
		if err != nil {
			return nil, fmt.Errorf("socks5 proxy %s: %w", proxyAddr, err)
		}
		if dl, ok := ctx.Deadline(); ok {
			_ = conn.SetDeadline(dl) //nolint:errcheck
		}
		if err := socks5Connect(conn, addr); err != nil {
			conn.Close() //nolint:errcheck,gosec
			return nil, fmt.Errorf("socks5 connect %s via %s: %w", addr, proxyAddr, err)
		}
		_ = conn.SetDeadline(time.Time{}) //nolint:errcheck
		return conn, nil
	}
}

func socks5Connect(conn net.Conn, addr string) error {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("bad port %q", portStr)
	}
	if len(host) > 255 {
		return fmt.Errorf("host name too long")
	}
	if _, err := conn.Write([]byte{5, 1, 0}); err != nil {
		return err
	}
	reply := make([]byte, 2)
	if _, err := io.ReadFull(conn, reply); err != nil {
		return err
	}
	if reply[0] != 5 || reply[1] != 0 {
		return fmt.Errorf("proxy refused no-auth method")
	}
	req := append([]byte{5, 1, 0, 3, byte(len(host))}, host...) //nolint:gosec // at most 255, checked above
	req = binary.BigEndian.AppendUint16(req, uint16(port))      //nolint:gosec // checked above
	if _, err := conn.Write(req); err != nil {
		return err
	}
	head := make([]byte, 4)
	if _, err := io.ReadFull(conn, head); err != nil {
		return err
	}
	if head[1] != 0 {
		return fmt.Errorf("proxy reply code %d", head[1])
	}
	var skip int
	switch head[3] {
	case 1:
		skip = 4
	case 4:
		skip = 16
	case 3:
		n := make([]byte, 1)
		if _, err := io.ReadFull(conn, n); err != nil {
			return err
		}
		skip = int(n[0])
	default:
		return fmt.Errorf("proxy reply address type %d", head[3])
	}
	_, err = io.ReadFull(conn, make([]byte, skip+2))
	return err
}
