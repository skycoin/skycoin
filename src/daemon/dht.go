package daemon

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nictuku/dht"
	"github.com/skycoin/skycoin/src/daemon/pex"
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

	// Port for DHT traffic (uses UDP)
	Port int
}

func NewDHTConfig() DHTConfig {
	return DHTConfig{
		Disabled:     false,
		Info:         "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6", //use genesis address for now
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
	cfg.Port = self.Config.Port
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

func PeerAddressExtract(addr string) string {
	ret := strings.Split(addr, ":")
	if len(ret) != 2 {
		return ""
	}
	//log.Printf("PeerAddr: %s, %s", ret[0], ret[1])
	//extract int
	ix, err := strconv.ParseUint(ret[1], 10, 16)
	if err != nil {
		log.Printf("DHT PeerAddr: Int Parse Error, %s \n", ret[1])
		return ""
	}
	if ix != 5999 {
		return ""
	}
	addr2 := fmt.Sprintf("%s:%d", ret[0], 6000)
	//log.Printf("addr= %s \n", addr2)
	return addr2
}

// Called when the DHT finds a peer
func (self *DHT) ReceivePeers(r map[dht.InfoHash][]string, peers *pex.Pex) {
	for _, results := range r {
		for _, p := range results {
			peer := dht.DecodePeerAddress(p)

			addr := PeerAddressExtract(peer)
			logger.Debug("DHT Peer: %s, Conn: %s", peer, addr)

			if addr != "" {
				_, err := peers.AddPeer(addr)
				if err != nil {
					logger.Info("Failed to add DHT peer: %v", err)
				}
			}

		}
	}
}
