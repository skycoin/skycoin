package daemon

import (
	"time"

	"github.com/skycoin/skycoin/src/aether/daemon/pex"
	gnet "github.com/skycoin/skycoin/src/aether/gnet" //use local gnet
)

// PeersConfig peers config
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
	// Disable exchanging of peers.  Peers are still loaded from disk
	Disabled bool
	//Ephemerial mode does not save or load from disc
	Ephemerial bool
}

// NewPeersConfig creates peer config
func NewPeersConfig() PeersConfig {
	return PeersConfig{
		DataDirectory:       "./",
		Max:                 1000,
		Expiration:          time.Hour * 24 * 7,
		CullRate:            time.Minute * 10,
		UpdateBlacklistRate: time.Minute,
		RequestRate:         time.Minute,
		ReplyCount:          30,
		Disabled:            false,
		Ephemerial:          false, //disable load/save to disc
	}
}

// Peers daemon peers
type Peers struct {
	Config PeersConfig
	// Peer list
	Peers *pex.Pex
}

// NewPeers create peers
func NewPeers(c PeersConfig) *Peers {
	if c.Disabled {
		logger.Info("PEX is disabled")
	}
	return &Peers{
		Config: c,
		Peers:  nil,
	}
}

// Init configure the pex.PeerList and load local data
func (p *Peers) Init() {
	peers := pex.NewPex(p.Config.Max)
	if p.Config.Ephemerial == false {
		err := peers.Load(p.Config.DataDirectory)
		if err != nil {
			logger.Notice("Failed to load peer database")
			logger.Notice("Reason: %v", err)
		}
	}

	logger.Debug("Init peers")
	p.Peers = peers
}

// Shutdown the PeerList
func (p *Peers) Shutdown() error {
	if p.Peers == nil {
		return nil
	}
	if p.Config.Ephemerial == false {
		err := p.Peers.Save(p.Config.DataDirectory)
		if err != nil {
			logger.Warning("Failed to save peer database")
			logger.Warning("Reason: %v", err)
			return err
		}
	}
	logger.Debug("Shutdown peers")
	return nil
}

// RemovePeer removes a peer, if not private
func (p *Peers) RemovePeer(a string) {
	peer := p.Peers.Peerlist[a]
	if peer != nil && !peer.Private {
		delete(p.Peers.Peerlist, a)
	}
}

// Requests peers from our connections
// TODO:
// - needs a callback for handling peers
// - needs DHT and PEX implemented
func (d *Daemon) requestPeers(service *gnet.Service, callback func(string)) {
	//d.DHT.RequestPeers(service.Id)
}
