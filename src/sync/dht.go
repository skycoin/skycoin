package sync

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/nictuku/dht"
	"github.com/skycoin/pex"
	"log"
	"time"
)

type DHTConfig struct {
	// Disable the DHT
	Disabled bool
	// Info to be hashed for identifying peers on the skycoin network
	Info string
	// Number of peers to try to get via DHT
	DesiredPeers int
	// How many local peers, from any source, before we stop requesting DHT peers
	PeerLimit int
	// DHT Bootstrap routers
	BootstrapNodes []string
	// How often to request more peers via DHT
	BootstrapRequestRate time.Duration

	// These should be set by the controlling daemon:
	// Port for DHT traffic (uses UDP)
	port int
}

/*
   Todo:
   - support multiple info hashes
   - callback for each info hashes when peer is found
*/

func NewDHTConfig() DHTConfig {
	return DHTConfig{
		Disabled:     false,
		Info:         "skycoin-skycoin-skycoin-skycoin-skycoin-skycoin-skycoin",
		DesiredPeers: 20,
		PeerLimit:    100,
		BootstrapNodes: []string{
			"1.a.magnets.im:6881",
			"router.bittorrent.com:6881",
			"router.utorrent.com:6881",
			"dht.transmissionbt.com:6881",
		},
		BootstrapRequestRate: time.Second * 10,
		port:                 6677,
	}
}

type DHT struct {
	Config DHTConfig
	// DHT manager
	DHT *dht.DHT
	// Hex encoded sha1 sum of Info
	InfoHash dht.InfoHash
}

func NewDHT(c DHTConfig) *DHT {
	return &DHT{
		Config:   c,
		DHT:      nil,
		InfoHash: "",
	}
}

// Sets up the DHT node for peer bootstrapping
func (self *DHT) Init() error {
	sum := sha1.Sum([]byte(self.Config.Info))
	// Create a hex encoded sha1 sum of a string to be used for DH
	infoHash, err := dht.DecodeInfoHash(hex.EncodeToString(sum[:]))
	if err != nil {
		return err
	}
	cfg := dht.NewConfig()
	cfg.Port = self.Config.port
	cfg.NumTargetPeers = self.Config.DesiredPeers
	d, err := dht.New(cfg)
	if err != nil {
		return err
	}
	self.InfoHash = infoHash
	self.DHT = d

	if self.Config.Disabled {
		// We have to initialize the DHT anyway because daemon loop needs
		// to read from its initialized chans. As long as Start() is prevented,
		// the DHT will not run.
		logger.Info("DHT is disabled")
	} else {
		logger.Info("Init DHT on port %d", self.Config.port)
	}
	return nil
}

// Stops the DHT
func (self *DHT) Shutdown() {
	if self.DHT != nil {
		logger.Debug("Stopping the DHT")
		self.DHT.Stop()
		// The DHT cannot be restarted once stopped, so we clear it
		self.DHT = nil
	}
}

// Starts the DHT
func (self *DHT) Start() {
	if self.Config.Disabled {
		return
	}
	self.DHT.Run()
}

// Requests peers from the DHT
func (self *DHT) RequestPeers() {
	if self.Config.Disabled {
		return
	}
	ih := string(self.InfoHash)
	if ih == "" {
		log.Panic("InfoHash is not initialized")
		return
	}
	logger.Info("Requesting DHT Peers")
	self.DHT.PeersRequest(ih, true)
}

// Called when the DHT finds a peer
func (self *DHT) ReceivePeers(r map[dht.InfoHash][]string, peers *pex.Pex) {
	for _, results := range r {
		for _, p := range results {
			peer := dht.DecodePeerAddress(p)
			logger.Debug("DHT Peer: %s", peer)
			_, err := peers.AddPeer(peer)
			if err != nil {
				logger.Info("Failed to add DHT peer: %v", err)
			}
		}
	}
}
