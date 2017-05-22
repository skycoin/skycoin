package dht

import (
	"crypto/sha1"
	"encoding/hex"

	"log"
	"time"

	"github.com/nictuku/dht"
	logging "github.com/op/go-logging"
)

/*
	TODO:
	- fix memory leak in map
	- map does not release memory after lookup
	- will use increasing memory for large number of lookups
	- add timer for last query and push out oldest after queries exceed number

*/
type AddPeerCallback func(infoHash string, peerAddress string)

var (
	logger = logging.MustGetLogger("skycoin.daemon_dht")
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
	Port int

	AddPeerCallback AddPeerCallback
}

/*
   Todo:
   - support multiple info hashes
   - callback for each info hashes when peer is found
*/

func NewDHTConfig() DHTConfig {
	return DHTConfig{
		Disabled: false,
		//Info:         "skycoin-skycoin-skycoin-skycoin-skycoin-skycoin-skycoin",
		DesiredPeers: 20,
		PeerLimit:    100,
		BootstrapNodes: []string{
			"1.a.magnets.im:6881",
			"router.bittorrent.com:6881",
			"router.utorrent.com:6881",
			"dht.transmissionbt.com:6881",
		},
		BootstrapRequestRate: time.Second * 10,
		Port:                 6677,
	}
}

type DHT struct {
	Config DHTConfig
	// DHT manager
	DHT *dht.DHT
	// Hex encoded sha1 sum of Info
	//InfoHash dht.InfoHash
	InfoHashes map[string]string //reverse map
}

func NewDHT(c DHTConfig) *DHT {
	return &DHT{
		Config:     c,
		DHT:        nil,
		InfoHashes: make(map[string]string),
		//InfoHash: "",
	}
}

// Sets up the DHT node for peer bootstrapping
func (self *DHT) Init() error {
	//sum := sha1.Sum([]byte(self.Config.Info))
	// Create a hex encoded sha1 sum of a string to be used for DH
	//infoHash, err := dht.DecodeInfoHash(hex.EncodeToString(sum[:]))

	cfg := dht.NewConfig()
	cfg.Port = self.Config.Port
	cfg.NumTargetPeers = self.Config.DesiredPeers
	d, err := dht.New(cfg)
	if err != nil {
		return err
	}
	//self.InfoHash = infoHash
	self.DHT = d

	//add the bootstrap nodes
	for _, addr := range self.Config.BootstrapNodes {
		self.DHT.AddNode(addr)
	}

	if self.Config.Disabled {
		// We have to initialize the DHT anyway because daemon loop needs
		// to read from its initialized chans. As long as Start() is prevented,
		// the DHT will not run.
		logger.Info("DHT is disabled")
	} else {
		logger.Info("Init DHT on port %d", self.Config.Port)
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

//run in separate goroutine for constantly flushing
func (self *DHT) Listen() {

	if self.Config.Disabled {
		return
	}

	for {
		r := <-self.DHT.PeersRequestResults
		if self.Config.Disabled {
			log.Panic("There should be no DHT peer results")
		}
		self.ReceivePeers(r)
	}

}

//call to flush pending events
func (self *DHT) FlushResults() {

	if self.Config.Disabled {
		return
	}

	select {
	case r := <-self.DHT.PeersRequestResults:
		if self.Config.Disabled {
			log.Panic("There should be no DHT peer results")
		}
		self.ReceivePeers(r)

	default:
	}

}

// Requests peers from the DHT
func (self *DHT) RequestPeers(infoHashString string) {
	if self.Config.Disabled {
		return
	}

	sum := sha1.Sum([]byte(infoHashString))
	// Create a hex encoded sha1 sum of a string to be used for DH
	infoHash, err := dht.DecodeInfoHash(hex.EncodeToString(sum[:]))
	if err != nil {
		log.Panic()
	}

	ih := string(infoHash)
	if ih == "" {
		log.Panic("InfoHash is not initialized")
		return
	}

	self.InfoHashes[ih] = infoHashString

	logger.Info("Requesting DHT Peers: infoHashString= %s", infoHashString)
	self.DHT.PeersRequest(ih, true)
}

//type InfoHash string

// Called when the DHT finds a peer
func (self *DHT) ReceivePeers(r map[dht.InfoHash][]string) {

	for infoHash, results := range r {
		for _, p := range results {
			peerAddress := dht.DecodePeerAddress(p)
			logger.Debug("DHT Peer: %s", peerAddress)

			if self.Config.AddPeerCallback == nil {
				log.Panic("Must set callback for receiving DHT peers")
			}

			str, ok := self.InfoHashes[string(infoHash)]
			if ok != true {
				log.Panic("infohash not requested")
			}

			self.Config.AddPeerCallback(str, peerAddress)

			//_, err := peers.AddPeer(peer)
			//if err != nil {
			//	logger.Info("Failed to add DHT peer: %v", err)
			//}
		}
	}
}
