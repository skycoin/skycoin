package webrpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Tests are setup as subtests, to retain a single *WebRPC instance for scaffolding
// https://blog.golang.org/subtests
func TestClient(t *testing.T) {
	s := setupWebRPC(t)
	errC := make(chan error, 1)
	go func() {
		errC <- s.Run()
	}()
	defer func() {
		time.Sleep(time.Millisecond * 50) // give rpc.Run() enough time to start
		err := s.Shutdown()
		require.NoError(t, err)
		require.NoError(t, <-errC)
	}()

	c := &Client{
		Addr: s.Addr,
	}

	testFuncs := []struct {
		n string
		f func(t *testing.T, c *Client, gw *fakeGateway)
	}{
		{"get unspent outputs", testClientGetUnspentOutputs},
		{"get status", testClientGetStatus},
		{"inject transaction", testClientInjectTransaction},
	}

	for _, f := range testFuncs {
		t.Run(f.n, func(t *testing.T) {
			f.f(t, c, s.Gateway.(*fakeGateway))
		})
	}
}

func testClientGetUnspentOutputs(t *testing.T, c *Client, gw *fakeGateway) {
	// address is copied from outputsStr
	addrs := []string{"fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B", "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW"}

	outputs, err := c.GetUnspentOutputs(addrs)
	require.NoError(t, err)
	require.Len(t, outputs.Outputs.HeadOutputs, 4)
	require.Len(t, outputs.Outputs.IncomingOutputs, 0)
	require.Len(t, outputs.Outputs.OutgoingOutputs, 0)

	require.Equal(t, outputs.Outputs, decodeOutputStr(outputStr))

	// Invalid address
	_, err = c.GetUnspentOutputs([]string{"invalid-address-foo"})
	require.Error(t, err)
	require.Equal(t, "invalid address: invalid-address-foo [code: -32602]", err.Error())
}

func testClientInjectTransaction(t *testing.T, c *Client, gw *fakeGateway) {
	// TODO -- how to make a raw tx?
	// Should we move some of the helper methods from CLI into webrpc?
}

func testClientGetStatus(t *testing.T, c *Client, gw *fakeGateway) {
	status, err := c.GetStatus()
	require.NoError(t, err)
	// values derived from hardcoded `blockString`
	require.Equal(t, status, &StatusResult{
		Running:            true,
		BlockNum:           455,
		LastBlockHash:      "",
		TimeSinceLastBlock: "18446744072232256374s",
	})
}
