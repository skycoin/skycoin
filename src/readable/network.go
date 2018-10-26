package readable

import (
	"github.com/skycoin/skycoin/src/daemon"
)

// Connection a connection's state within the daemon
type Connection struct {
	GnetID       int    `json:"id"`
	Addr         string `json:"address"`
	LastSent     int64  `json:"last_sent"`
	LastReceived int64  `json:"last_received"`
	ConnectedAt  int64  `json:"connected_at"`
	Outgoing     bool   `json:"outgoing"`
	State        string `json:"state"`
	Mirror       uint32 `json:"mirror"`
	ListenPort   uint16 `json:"listen_port"`
	Height       uint64 `json:"height"`
}

// NewConnection copies daemon.Connection to a struct with json tags
func NewConnection(c *daemon.Connection) Connection {
	var lastSent int64
	var lastReceived int64
	var connectedAt int64

	if !c.LastSent.IsZero() {
		lastSent = c.LastSent.Unix()
	}
	if !c.LastReceived.IsZero() {
		lastReceived = c.LastReceived.Unix()
	}
	if !c.ConnectedAt.IsZero() {
		connectedAt = c.ConnectedAt.Unix()
	}

	return Connection{
		GnetID:       c.GnetID,
		Addr:         c.Addr,
		LastSent:     lastSent,
		LastReceived: lastReceived,
		ConnectedAt:  connectedAt,
		Outgoing:     c.Outgoing,
		State:        string(c.State),
		Mirror:       c.Mirror,
		ListenPort:   c.ListenPort,
		Height:       c.Height,
	}
}
