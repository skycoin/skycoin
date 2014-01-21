package daemon

import (
    "crypto/sha1"
    "encoding/hex"
    "github.com/nictuku/dht"
    "log"
    "time"
)

var (
    // DHT manager
    DHT *dht.DHT = nil
    // Info to be hashed for identifying peers on the skycoin network
    dhtInfo = "skycoin-skycoin-skycoin-skycoin-skycoin-skycoin-skycoin"
    // Hex encoded sha1 sum of dhtInfo
    dhtInfoHash dht.InfoHash = ""
    // Number of peers to try to get via DHT
    dhtDesiredPeers = 20
    // How many local peers, from any source, before we stop requesting DHT peers
    dhtPeerLimit = 100
    // DHT Bootstrap routers
    dhtBootstrapNodes = []string{
        "1.a.magnets.im:6881",
        "router.bittorrent.com:6881",
        "router.utorrent.com:6881",
        "dht.transmissionbt.com:6881",
    }
    // How often to request more peers via DHT
    dhtBootstrapRequestRate = time.Second * 10
)

// Sets up the DHT node for peer bootstrapping
func InitDHT(port int) {
    var err error
    sum := sha1.Sum([]byte(dhtInfo))
    // Create a hex encoded sha1 sum of a string to be used for DH
    dhtInfoHash, err = dht.DecodeInfoHash(hex.EncodeToString(sum[:]))
    if err != nil {
        log.Panicf("Failed to create InfoHash: %v", err)
        return
    }
    DHT, err = dht.NewDHTNode(port, dhtDesiredPeers, true)
    if err != nil {
        log.Panicf("Failed to init DHT: %v", err)
        return
    }
    logger.Info("Init DHT on port %d", port)
}

// Called when the DHT finds a peer
func receivedDHTPeers(r map[dht.InfoHash][]string) {
    for _, peers := range r {
        for _, p := range peers {
            peer := dht.DecodePeerAddress(p)
            logger.Debug("DHT Peer: %s", peer)
            _, err := Peers.AddPeer(peer)
            if err != nil {
                logger.Info("Failed to add DHT peer: %v", err)
            }
        }
    }
}

// Requests peers from the DHT
func RequestDHTPeers() {
    ih := string(dhtInfoHash)
    if ih == "" {
        log.Panic("dhtInfoHash is not initialized")
        return
    }
    logger.Info("Requesting DHT Peers")
    DHT.PeersRequest(ih, true)
}

// // DHT Event Logger
// type DHTLogger struct{}

// // Logs a GetPeers event
// func (self *DHTLogger) GetPeers(ip *net.UDPAddr, id string,
//     _info dht.InfoHash) {
//     id = hex.EncodeToString([]byte(id))
//     info := hex.EncodeToString([]byte(_info))
//     logger.Debug("DHT GetPeers event occured:\n\tid: %s\n\tinfohash: %s",
//         id, info)
// }
