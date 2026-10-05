package electrum

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// fakeSOCKS5 accepts one no-auth CONNECT, records the requested host and port,
// replies success and echoes what follows.
func fakeSOCKS5(t *testing.T) (addr string, got chan string) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { l.Close() }) //nolint:errcheck,gosec
	got = make(chan string, 1)
	go func() {
		c, err := l.Accept()
		if err != nil {
			return
		}
		defer c.Close() //nolint:errcheck
		b := make([]byte, 3)
		if _, err := io.ReadFull(c, b); err != nil {
			return
		}
		c.Write([]byte{5, 0}) //nolint:errcheck,gosec
		head := make([]byte, 5)
		if _, err := io.ReadFull(c, head); err != nil {
			return
		}
		rest := make([]byte, int(head[4])+2)
		if _, err := io.ReadFull(c, rest); err != nil {
			return
		}
		port := binary.BigEndian.Uint16(rest[len(rest)-2:])
		got <- net.JoinHostPort(string(rest[:len(rest)-2]), itoa(port))
		c.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0}) //nolint:errcheck,gosec
		io.Copy(c, c)                                 //nolint:errcheck,gosec
	}()
	return l.Addr().String(), got
}

func itoa(p uint16) string { return strconv.Itoa(int(p)) }

func TestSOCKS5DialerSendsHostByName(t *testing.T) {
	proxy, got := fakeSOCKS5(t)
	conn, err := SOCKS5Dialer(proxy, 5*time.Second, nil)("tcp", "electrum.example.org:50002")
	require.NoError(t, err)
	defer conn.Close() //nolint:errcheck
	require.Equal(t, "electrum.example.org:50002", <-got)

	_, err = conn.Write([]byte("ping"))
	require.NoError(t, err)
	b := make([]byte, 4)
	_, err = io.ReadFull(conn, b)
	require.NoError(t, err)
	require.Equal(t, "ping", string(b))
}
