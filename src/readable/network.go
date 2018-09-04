package readable

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
)

// Connection a connection's state within the daemon
type Connection struct {
	ID           int    `json:"id"`
	Addr         string `json:"address"`
	LastSent     int64  `json:"last_sent"`
	LastReceived int64  `json:"last_received"`
	// Whether the connection is from us to them (true, outgoing),
	// or from them to us (false, incoming)
	Outgoing bool `json:"outgoing"`
	// Whether the client has identified their version, mirror etc
	Introduced bool   `json:"introduced"`
	Mirror     uint32 `json:"mirror"`
	ListenPort uint16 `json:"listen_port"`
	Height     uint64 `json:"height"`
}

// NewConnection copies daemon.Connection to a struct with json tags
func NewConnection(c *daemon.Connection) Connection {
	return Connection{
		ID:           c.ID,
		Addr:         c.Addr,
		LastSent:     c.LastSent,
		LastReceived: c.LastReceived,
		Outgoing:     c.Outgoing,
		Introduced:   c.Introduced,
		Mirror:       c.Mirror,
		ListenPort:   c.ListenPort,
		Height:       c.Height,
	}
}

// Connections wraps []Connection
type Connections struct {
	Connections []Connection `json:"connections"`
}

// NewConnections copies []daemon.Connection to a struct with json tags
func NewConnections(dconns []daemon.Connection) Connections {
	conns := make([]Connection, len(dconns))
	for i, dc := range dconns {
		conns[i] = NewConnection(&dc)
	}

	return Connections{
		Connections: conns,
	}
}

// ResendResult the result of rebroadcasting transaction
type ResendResult struct {
	Txids []string `json:"txids"`
}

// NewResendResult creates a ResendResult from a list of transaction ID hashes
func NewResendResult(hashes []cipher.SHA256) ResendResult {
	txids := make([]string, len(hashes))
	for i, h := range hashes {
		txids[i] = h.Hex()
	}
	return ResendResult{
		Txids: txids,
	}
}
