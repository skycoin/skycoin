package daemon

import (
	"time"

	"os"

	"github.com/skycoin/skycoin/src/daemon/pex"
)

// PeersConfig config for peers
type PeersConfig struct {
	// Folder where peers database should be saved
	DataDirectory string
	// Maximum number of peers to keep account of in the PeerList
	Max int
	// Cull peers after they havent been seen in this much time
	Expiration time.Duration
	// Cull expired peers on this interval
	CullRate time.Duration
	// How often to clear expired blacklist entries
	UpdateBlacklistRate time.Duration
	// How often to request peers via PEX
	RequestRate time.Duration
	// How many peers to send back in response to a peers request
	ReplyCount int
	// Localhost peers are allowed in the peerlist
	AllowLocalhost bool
	// Disable exchanging of peers.  Peers are still loaded from disk
	Disabled bool
}

// NewPeersConfig creates peers config
func NewPeersConfig() PeersConfig {
	return PeersConfig{
		DataDirectory:       "./",
		Max:                 1000,
		Expiration:          time.Hour * 24 * 7,
		CullRate:            time.Minute * 10,
		UpdateBlacklistRate: time.Minute,
		RequestRate:         time.Minute,
		ReplyCount:          30,
		AllowLocalhost:      false,
		Disabled:            false,
	}
}

// Peers maintains the config and peers instance
type Peers struct {
	Config PeersConfig
	// Peer list
	Peers *pex.Pex
}

// NewPeers creates peers
func NewPeers(c PeersConfig) *Peers {
	if c.Disabled {
		logger.Info("PEX is disabled")
	}
	return &Peers{
		Config: c,
		Peers:  nil,
	}
}

// DefaultConnections do "default_peers file"
// read file, write, if does not exist
var DefaultConnections = []string{}

// Init configures the pex.PeerList and load local data
func (ps *Peers) Init() {
	peers := pex.NewPex(ps.Config.Max)
	err := peers.Load(ps.Config.DataDirectory)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Notice("Failed to load peer database")
			logger.Notice("Reason: %v", err)
		}
	}
	logger.Debug("Init peers")
	peers.AllowLocalhost = ps.Config.AllowLocalhost

	//Boot strap peers
	for _, addr := range DefaultConnections {
		// default peers will mark as trusted peers.
		_, err := peers.AddPeer(addr)
		if err != nil {
			logger.Critical("add peer error:%v", err)
		}
		peers.SetTrustState(addr, true)
	}

	ps.Peers = peers
	ps.Peers.Save(ps.Config.DataDirectory)
}

// Shutdown the PeerList
func (ps *Peers) Shutdown() error {
	if ps.Peers == nil {
		return nil
	}

	logger.Debug("Saving Peer List")

	err := ps.Peers.Save(ps.Config.DataDirectory)
	if err != nil {
		logger.Warning("Failed to save peer database")
		logger.Warning("Reason: %v", err)
		return err
	}
	logger.Debug("Shutdown peers")
	return nil
}

// RemovePeer removes a peer, if not private
func (ps *Peers) RemovePeer(a string) {
	ps.Peers.RemovePeer(a)
}

// Requests peers from our connections
func (ps *Peers) requestPeers(pool *Pool) {
	if ps.Config.Disabled {
		return
	}
	if ps.Peers.Full() {
		return
	}
	m := NewGetPeersMessage()
	pool.Pool.BroadcastMessage(m)
}
