package btc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeElectrum answers server.version and blockchain.headers.subscribe with
// tip height, over plain tcp://, until closed.
func fakeElectrum(t *testing.T, height int) (url string, stop func()) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close() //nolint:errcheck
				s := bufio.NewScanner(c)
				for s.Scan() {
					var req struct {
						ID     json.RawMessage `json:"id"`
						Method string          `json:"method"`
					}
					if json.Unmarshal(s.Bytes(), &req) != nil {
						return
					}
					result := `["fake 1.0","1.4"]`
					if req.Method == "blockchain.headers.subscribe" {
						result = fmt.Sprintf(`{"height":%d,"hex":""}`, height)
					}
					fmt.Fprintf(c, `{"jsonrpc":"2.0","id":%s,"result":%s}`+"\n", req.ID, result) //nolint:errcheck
				}
			}(c)
		}
	}()
	return "tcp://" + l.Addr().String(), func() { l.Close() } //nolint:errcheck,gosec
}

func deadServer(t *testing.T) string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr().String()
	require.NoError(t, l.Close())
	return "tcp://" + addr
}

// The first server that answers is used, and a server that goes away is
// replaced by the next on the following call.
func TestFailoverElectrum(t *testing.T) {
	a, stopA := fakeElectrum(t, 100)
	b, stopB := fakeElectrum(t, 200)
	defer stopB()

	f, err := NewFailoverElectrum([]string{deadServer(t), a, b}, nil)
	require.NoError(t, err)
	defer f.Close() //nolint:errcheck

	h, err := f.Health()
	require.NoError(t, err)
	require.Equal(t, 100, h, "the dead first server is skipped")

	stopA()
	f.mu.Lock()
	f.cur.Close() //nolint:errcheck,gosec
	f.mu.Unlock()
	h, err = f.Health()
	require.NoError(t, err)
	require.Equal(t, 200, h, "a failed call moves to the next server")
}

func TestFailoverElectrumNoneAnswer(t *testing.T) {
	f, err := NewFailoverElectrum([]string{deadServer(t)}, nil)
	require.NoError(t, err)
	_, err = f.Health()
	require.Error(t, err)
}
