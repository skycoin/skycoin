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
	Txids []string `json:"txids"` // transaction id
}

type RPC struct{}

func (self RPC) GetConnection(d *Daemon, addr string) *Connection {
	if d.Pool.Pool == nil {
		return nil
	}

	c := d.Pool.Pool.GetConnection(addr)
	if c == nil {
		return nil
	}

	mirror, exist := d.connectionMirrors.Get(addr)
	if !exist {
		return nil
	}

	return &Connection{
		Id:           c.Id,
		Addr:         addr,
		LastSent:     c.LastSent.Unix(),
		LastReceived: c.LastReceived.Unix(),
		Outgoing:     !d.outgoingConnections.Get(addr),
		Introduced:   !d.needsIntro(addr),
		Mirror:       mirror,
		ListenPort:   d.GetListenPort(addr),
	}
}

func (self RPC) GetConnections(d *Daemon) *Connections {
	if d.Pool.Pool == nil {
		return nil
	}
	conns := make([]*Connection, 0, d.Pool.Pool.Size())
	for _, c := range d.Pool.Pool.GetConnections() {
		conn := self.GetConnection(d, c.Addr())
		if conn != nil {
			conns = append(conns, conn)
		}
	}
	return &Connections{Connections: conns}
}

func (self RPC) GetDefaultConnections(d *Daemon) []string {
	return d.DefaultConnections
}

func (self RPC) GetTrustConnections(d *Daemon) []string {
	peers := d.Peers.Peers.GetAllTrustedPeers()
	addrs := make([]string, len(peers))
	for i, p := range peers {
		addrs[i] = p.Addr
	}
	return addrs
}

// GetAllExchgConnections return all exchangeable connections
func (rpc RPC) GetAllExchgConnections(d *Daemon) []string {
	peers := d.Peers.Peers.RandomExchgAll(0)
	addrs := make([]string, len(peers))
	for i, p := range peers {
		addrs[i] = p.Addr
	}
	return addrs
}

func (self RPC) GetBlockchainProgress(v *Visor) *BlockchainProgress {
	if v.v == nil {
		return nil
	}
	return &BlockchainProgress{
		Current: v.HeadBkSeq(),
		Highest: v.EstimateBlockchainLength(),
	}
}

func (self RPC) ResendTransaction(v *Visor, p *Pool, txHash cipher.SHA256) *ResendResult {
	if v.v == nil {
		return nil
	}
	v.ResendTransaction(txHash, p)
	return &ResendResult{}
}

func (self RPC) ResendUnconfirmedTxns(v *Visor, p *Pool) *ResendResult {
	if v.v == nil {
		return nil
	}
	txids := v.ResendUnconfirmedTxns(p)
	var rlt ResendResult
	for _, txid := range txids {
		rlt.Txids = append(rlt.Txids, txid.Hex())
	}
	return &rlt
}
