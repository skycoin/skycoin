package daemon

import (
	"time"

	"github.com/skycoin/skycoin/src/daemon/pex"
)

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

type Peers struct {
	Config PeersConfig
	// Peer list
	Peers *pex.Pex
}

func NewPeers(c PeersConfig) *Peers {
	if c.Disabled {
		logger.Info("PEX is disabled")
	}
	return &Peers{
		Config: c,
		Peers:  nil,
	}
}

//do "default_peers file"
//read file, write, if does not exist
var DefaultConnections = []string{}

// Configure the pex.PeerList and load local data
func (self *Peers) Init() {
	peers := pex.NewPex(self.Config.Max)
	err := peers.Load(self.Config.DataDirectory)
	if err != nil {
		logger.Notice("Failed to load peer database")
		logger.Notice("Reason: %v", err)
	}
	logger.Debug("Init peers")
	peers.AllowLocalhost = self.Config.AllowLocalhost

	//Boot strap peers
	for _, addr := range DefaultConnections {
		peers.AddPeer(addr)
	}

	self.Peers = peers
}

// Shutdown the PeerList
func (self *Peers) Shutdown() error {
	if self.Peers == nil {
		return nil
	}

	logger.Debug("Saving Peer List")

	err := self.Peers.Save(self.Config.DataDirectory)
	if err != nil {
		logger.Warning("Failed to save peer database")
		logger.Warning("Reason: %v", err)
		return err
	}
	logger.Debug("Shutdown peers")
	return nil
}

// Removes a peer, if not private
func (self *Peers) RemovePeer(a string) {
	p := self.Peers.Peerlist[a]
	if p != nil && !p.Private {
		delete(self.Peers.Peerlist, a)
	}
}

// Requests peers from our connections
func (self *Peers) requestPeers(pool *Pool) {
	if self.Config.Disabled {
		return
	}
	if self.Peers.Full() {
		return
	}
	m := NewGetPeersMessage()
	pool.Pool.BroadcastMessage(m)
}
