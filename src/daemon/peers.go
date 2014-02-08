package daemon

import (
    "github.com/skycoin/pex"
    "time"
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

// Configure the pex.PeerList and load local data
func (self *Peers) Init() {
    peers := pex.NewPex(self.Config.Max)
    err := peers.Load(self.Config.DataDirectory)
    if err != nil {
        logger.Notice("Failed to load peer database")
        logger.Notice("Reason: %v", err)
    }
    logger.Debug("Init peers")
    self.Peers = peers
}

// Shutdown the PeerList
func (self *Peers) Shutdown() error {
    if self.Peers == nil {
        return nil
    }
    err := self.Peers.Save(self.Config.DataDirectory)
    if err != nil {
        logger.Warning("Failed to save peer database")
        logger.Warning("Reason: %v", err)
        return err
    }
    logger.Debug("Shutdown peers")
    return nil
}

// Requests peers from our connections
// TODO -- batching all peer requests at once may cause performance issues
func (self *Peers) requestPeers(pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewGetPeersMessage()
    pool.Pool.Dispatcher.BroadcastMessage(m)
}
