package daemon

import (
	"github.com/skycoin/skycoin/src/cipher"
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
}

// Connections an array of connections
// Arrays must be wrapped in structs to avoid certain javascript exploits
type Connections struct {
	Connections []*Connection `json:"connections"`
}

// BlockchainProgress current sync blockchain status
type BlockchainProgress struct {
	// Our current blockchain length
	Current uint64 `json:"current"`
	// Our best guess at true blockchain length
	Highest uint64 `json:"highest"`
	Peers   []struct {
		Address string `json:"address"`
		Height  uint64 `json:"height"`
	} `json:"peers"`
}

// ResendResult rebroadcast tx result
type ResendResult struct {
	Txids []string `json:"txids"` // transaction id
}

// RPC rpc
type RPC struct{}

// GetConnection gets connection of given address
func (rpc RPC) GetConnection(d *Daemon, addr string) *Connection {
	if d.Pool.Pool == nil {
		return nil
	}

	c, err := d.Pool.Pool.GetConnection(addr)
	if err != nil {
		logger.Error("%v", err)
		return nil
	}

	if c == nil {
		return nil
	}

	mirror, exist := d.connectionMirrors.Get(addr)
	if !exist {
		return nil
	}

	return &Connection{
		ID:           c.ID,
		Addr:         addr,
		LastSent:     c.LastSent.Unix(),
		LastReceived: c.LastReceived.Unix(),
		Outgoing:     !d.outgoingConnections.Get(addr),
		Introduced:   !d.needsIntro(addr),
		Mirror:       mirror,
		ListenPort:   d.GetListenPort(addr),
	}
}

// GetConnections gets all connections
func (rpc RPC) GetConnections(d *Daemon) *Connections {
	if d.Pool.Pool == nil {
		return nil
	}

	l, err := d.Pool.Pool.Size()
	if err != nil {
		logger.Error("%v", err)
		return nil
	}

	conns := make([]*Connection, 0, l)
	cs, err := d.Pool.Pool.GetConnections()
	if err != nil {
		logger.Error("%v", err)
		return nil
	}

	for _, c := range cs {
		if c.Solicited {
			conn := rpc.GetConnection(d, c.Addr())
			if conn != nil {
				conns = append(conns, conn)
			}
		}
	}
	return &Connections{Connections: conns}
}

// GetDefaultConnections gets default connections
func (rpc RPC) GetDefaultConnections(d *Daemon) []string {
	return d.DefaultConnections
}

// GetTrustConnections get all trusted transaction
func (rpc RPC) GetTrustConnections(d *Daemon) []string {
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

// GetBlockchainProgress gets the blockchain progress
func (rpc RPC) GetBlockchainProgress(v *Visor) *BlockchainProgress {
	if v.v == nil {
		return nil
	}

	bp := &BlockchainProgress{
		Current: v.HeadBkSeq(),
		Highest: v.EstimateBlockchainHeight(),
	}

	peerHeights := v.GetPeerBlockchainHeights()

	for _, ph := range peerHeights {
		bp.Peers = append(bp.Peers, struct {
			Address string `json:"address"`
			Height  uint64 `json:"height"`
		}{
			Address: ph.Address,
			Height:  ph.Height,
		})
	}

	return bp
}

// ResendTransaction rebroadcast transaction
func (rpc RPC) ResendTransaction(v *Visor, p *Pool, txHash cipher.SHA256) *ResendResult {
	if v.v == nil {
		return nil
	}
	v.ResendTransaction(txHash, p)
	return &ResendResult{}
}

// ResendUnconfirmedTxns rebroadcast unconfirmed transactions
func (rpc RPC) ResendUnconfirmedTxns(v *Visor, p *Pool) *ResendResult {
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
