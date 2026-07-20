// Package market is the skydex-client's transport-agnostic connection to an
// skydex-market.
//
// It does not open the connection: the caller supplies an established net.Conn
// (over skywire, an appnet dmsg conn dialed on the market's public key by the
// skywire wrapper; in tests, a net.Pipe). This package only speaks the exchange
// protocol on top of it: one length-prefixed frame per JSON Envelope,
// request/response correlated by Envelope.ID.
package market

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/skydex/message"
	"github.com/skycoin/skycoin/src/skydex/protocol"
)

// Conn is a framed, request/response connection to a market. It is safe for
// concurrent use: each Do serializes a single request/response exchange, which
// matches the client's polling model (one outstanding request at a time).
type Conn struct {
	fc  *message.Conn
	mu  sync.Mutex
	rtt time.Duration // per-request read/write deadline; 0 = none
}

// NewConn wraps an established net.Conn with exchange framing. The skywire
// wrapper obtains the conn by dialing the market's public key over appnet and
// passes it here; tests pass a net.Pipe end.
func NewConn(c net.Conn) *Conn {
	return &Conn{fc: message.NewConn(c), rtt: 30 * time.Second}
}

// Close closes the underlying connection.
func (c *Conn) Close() error { return c.fc.Close() }

// Do sends a request of the given type carrying data and returns the market's
// response envelope. A transport error is returned as err; an application-level
// failure is returned as a non-error response with resp.IsError() == true.
func (c *Conn) Do(msgType string, data any) (protocol.Envelope, error) {
	req, err := protocol.NewRequest(msgType, data)
	if err != nil {
		return protocol.Envelope{}, err
	}
	payload, err := req.Marshal()
	if err != nil {
		return protocol.Envelope{}, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rtt > 0 {
		// Best-effort deadline so a dead transport can't pin the caller.
		_ = c.fc.SetDeadline(time.Now().Add(c.rtt)) //nolint:errcheck
		defer c.fc.SetDeadline(time.Time{})         //nolint:errcheck
	}

	if err := c.fc.WriteFrame(payload); err != nil {
		return protocol.Envelope{}, fmt.Errorf("write request: %w", err)
	}
	respPayload, err := c.fc.ReadFrame()
	if err != nil {
		return protocol.Envelope{}, fmt.Errorf("read response: %w", err)
	}
	resp, err := protocol.Unmarshal(respPayload)
	if err != nil {
		return protocol.Envelope{}, fmt.Errorf("decode response: %w", err)
	}
	if resp.ID != req.ID {
		return protocol.Envelope{}, fmt.Errorf("response id mismatch: got %q want %q", resp.ID, req.ID)
	}
	return resp, nil
}
