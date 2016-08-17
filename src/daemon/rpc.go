package daemon

import (
	"github.com/skycoin/skycoin/src/cipher"
	//"github.com/skycoin/skycoin/src/visor"
	//"github.com/skycoin/skycoin/src/wallet"
)

// A connection's state within the daemon
type Connection struct {
	Id           int    `json:"id"`
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
}

// An array of connections
// Arrays must be wrapped in structs to avoid certain javascript exploits
type Connections struct {
	Connections []*Connection `json:"connections"`
}

type BlockchainProgress struct {
	// Our current blockchain length
	Current uint64 `json:"current"`
	// Our best guess at true blockchain length
	Highest uint64 `json:"highest"`
}

type ResendResult struct {
}

type RPC struct{}

func (self RPC) GetConnection(d *Daemon, addr string) *Connection {
	if d.Pool.Pool == nil {
		return nil
	}
	c := d.Pool.Pool.Addresses[addr]
	if c == nil {
		return nil
	}
	_, expecting := d.ExpectingIntroductions[addr]
	return &Connection{
		Id:           c.Id,
		Addr:         addr,
		LastSent:     c.LastSent.Unix(),
		LastReceived: c.LastReceived.Unix(),
		Outgoing:     (d.OutgoingConnections[addr] == nil),
		Introduced:   !expecting,
		Mirror:       d.ConnectionMirrors[addr],
		ListenPort:   d.GetListenPort(addr),
	}
}

func (self RPC) GetConnections(d *Daemon) *Connections {
	if d.Pool.Pool == nil {
		return nil
	}
	conns := make([]*Connection, len(d.Pool.Pool.Pool))
	for i, c := range d.Pool.Pool.GetConnections() {
		conns[i] = self.GetConnection(d, c.Addr())
	}
	return &Connections{Connections: conns}
}

func (self RPC) GetDefaultConnections(d *Daemon) []string {
	return d.DefaultConnections
}

func (self RPC) GetBlockchainProgress(v *Visor) *BlockchainProgress {
	if v.Visor == nil {
		return nil
	}
	return &BlockchainProgress{
		Current: v.Visor.HeadBkSeq(),
		Highest: v.EstimateBlockchainLength(),
	}
}

func (self RPC) ResendTransaction(v *Visor, p *Pool,
	txHash cipher.SHA256) *ResendResult {
	if v.Visor == nil {
		return nil
	}
	v.ResendTransaction(txHash, p)
	return &ResendResult{}
}
