package daemon

import (
    "github.com/skycoin/pex"
    "time"
)

var (
    // Maximum number of peers to keep account of in the PeerList
    maxPeers = 1000
    // Peer list
    Peers = pex.NewPex(maxPeers)
    // Cull peers after they havent been seen in this much time
    peerExpiration = time.Hour * 24 * 7
    // Cull expired peers on this interval
    cullPeerRate = time.Minute * 10
    // How often to clear expired blacklist entries
    updateBlacklistRate = time.Minute
    // How often to request peers via PEX
    requestPeersRate = time.Minute
    // How many peers to send back in response to a peers request
    peerReplyCount = 30
)

// Configure the pex.PeerList and load local data
func InitPeers(data_directory string) {
    err := Peers.Load(data_directory)
    if err != nil {
        logger.Notice("Failed to load peer database")
        logger.Notice("Reason: %v", err)
    }
    logger.Debug("Init peers")
}

// Shutdown the PeerList
func ShutdownPeers(data_directory string) {
    err := Peers.Save(data_directory)
    if err != nil {
        logger.Warning("Failed to save peer database")
        logger.Warning("Reason: %v", err)
    }
    Peers = nil
    logger.Debug("Shutdown peers")
}
