// Package pex is a toolkit for implementing a peer exchange system
package pex

import (
	"errors"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"math"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
)

//TODO:
// - keep track of last time the peer was connected to
// - last time peer was connected to is more important than "seen"
// - peer "seen" means something else than use here
// - save last time connected to, use 0 for never
// - only transmit peers that have active or recent connections

var (
	// PeerDatabaseFilename filename for disk-cached peers
	PeerDatabaseFilename = "peers.txt"
	// BlacklistedDatabaseFilename  filename for disk-cached blacklisted peers
	BlacklistedDatabaseFilename = "blacklisted_peers.txt"
	// ErrPeerlistFull returned when the Pex is at a maximum
	ErrPeerlistFull = errors.New("Peer list full")
	// ErrInvalidAddress Returned when an address appears malformed
	ErrInvalidAddress = errors.New("Invalid address")
	// ErrBlacklistedAddress returned when attempting to add a blacklisted peer
	ErrBlacklistedAddress = errors.New("Blacklisted address")
	// RefreshBlacklistRate How often to updated expired entries in the blacklist
	RefreshBlacklistRate = time.Second * 30
	// Logging. See http://godoc.org/github.com/op/go-logging for
	// instructions on how to include this log's output
	logger = logging.MustGetLogger("pex")
	// Default rng
	rnum = rand.New(rand.NewSource(time.Now().Unix()))
	// For removing inadvertent whitespace from addresses
	whitespaceFilter = regexp.MustCompile("\\s")
	maxRetryTimes    = 3
)

// validateAddress returns true if ipPort is a valid ip:host string
func validateAddress(ipPort string, allowLocalhost bool) bool {
	ipPort = whitespaceFilter.ReplaceAllString(ipPort, "")
	pts := strings.Split(ipPort, ":")
	if len(pts) != 2 {
		return false
	}

	ip := net.ParseIP(pts[0])
	if ip == nil {
		return false
	} else if ip.IsLoopback() {
		if !allowLocalhost {
			return false
		}
	} else if !ip.IsGlobalUnicast() {
		return false
	}

	port, err := strconv.ParseUint(pts[1], 10, 16)
	if err != nil || port < 1024 {
		return false
	}

	return true
}

// Peer represents a known peer
type Peer struct {
	Addr       string // An address of the form ip:port
	LastSeen   int64  // Unix timestamp when this peer was last seen
	Private    bool   // Whether it should omitted from public requests
	Trusted    bool   // Whether this peer is trusted
	Valid      bool   // Whether this peer is valid(has reachable public port)
	RetryTimes int    `json:"-"` // records the retry times
}

// NewPeer returns a *Peer initialised by an address string of the form ip:port
func NewPeer(address string) *Peer {
	p := &Peer{
		Addr:    address,
		Private: false,
		Trusted: false,
	}
	p.Seen()
	return p
}

// Seen marks the peer as seen
func (peer *Peer) Seen() {
	peer.LastSeen = utc.UnixNow()
}

// IncreaseRetryTimes adds the retry times
func (peer *Peer) IncreaseRetryTimes() {
	peer.RetryTimes++
	logger.Debug("Increase retry times of %v: %v", peer.Addr, peer.RetryTimes)
}

// ResetRetryTimes resets the retry time
func (peer *Peer) ResetRetryTimes() {
	peer.RetryTimes = 0
	logger.Debug("Reset retry times of %v", peer.Addr)
}

// CanTry returns whether this peer is tryable base on the exponential backoff algorithm
func (peer *Peer) CanTry() bool {
	// Exponential backoff
	mod := (math.Exp2(float64(peer.RetryTimes)) - 1) * 5
	if mod == 0 {
		return true
	}

	// Random time elapsed
	now := utc.UnixNow()
	t := rnum.Int63n(int64(mod))
	return now-peer.LastSeen > t
}

// String returns the peer address
func (peer *Peer) String() string {
	return peer.Addr
}

// Config pex config
type Config struct {
	// Folder where peers database should be saved
	DataDirectory string
	// Maximum number of peers to keep account of in the PeerList
	Max int
	// Cull peers after they havent been seen in this much time
	Expiration time.Duration
	// Cull expired peers on this interval
	CullRate time.Duration
	// clear old peers on this interval
	ClearOldRate time.Duration
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
	// Whether the network is disabled
	NetworkDisabled bool
}

// NewConfig creates default pex config.
func NewConfig() Config {
	return Config{
		DataDirectory:       "./",
		Max:                 1000,
		Expiration:          time.Hour * 24 * 7,
		CullRate:            time.Minute * 10,
		ClearOldRate:        time.Minute * 10,
		UpdateBlacklistRate: time.Minute,
		RequestRate:         time.Minute,
		ReplyCount:          30,
		AllowLocalhost:      false,
		Disabled:            false,
		NetworkDisabled:     false,
	}
}

// Pex manages a set of known peers and controls peer acquisition
type Pex struct {
	// All known peers
	*peerlist
	Config    Config
	blacklist *blacklist
	quit      chan struct{}
}

// New creates pex
func New(cfg Config, defaultConns []string) (*Pex, error) {
	pl := newPeerlist(cfg.Max)
	if err := pl.loadIfExist(cfg.DataDirectory); err != nil {
		return nil, err
	}

	bl := newBlacklist()

	pex := &Pex{
		Config:    cfg,
		peerlist:  pl,
		blacklist: bl,
		quit:      make(chan struct{}),
	}

	if err := pex.loadIfExist(cfg.DataDirectory); err != nil {
		return nil, err
	}

	//Boot strap peers
	for _, addr := range defaultConns {
		// default peers will mark as trusted peers.
		if err := pex.AddPeer(addr); err != nil {
			logger.Critical("add peer failed:%v", err)
			continue
		}
		if err := pex.SetTrust(addr, true); err != nil {
			logger.Critical("pex.SetTrust failed: %v", err)
		}
	}

	// persists peers
	if err := pex.save(cfg.DataDirectory); err != nil {
		return nil, err
	}

	return pex, nil
}

// AddPeer adds a peer to the peer list, given an address. If the peer list is
// full, PeerlistFullError is returned */
func (px *Pex) AddPeer(addr string) error {
	if !validateAddress(addr, px.Config.AllowLocalhost) {
		return ErrInvalidAddress
	}

	return px.addPeer(addr)
}

// SetPrivate updates peer's private value
func (px *Pex) SetPrivate(addr string, private bool) error {
	if !validateAddress(addr, px.Config.AllowLocalhost) {
		return ErrInvalidAddress
	}

	return px.setPrivate(addr, private)
}

// SetTrust updates peer's trusted value
func (px *Pex) SetTrust(addr string, trusted bool) error {
	if !validateAddress(addr, px.Config.AllowLocalhost) {
		return ErrInvalidAddress
	}

	return px.setTrusted(addr, trusted)
}

// SetValid updates peer's Valid value
func (px *Pex) SetValid(addr string, valid bool) error {
	if !validateAddress(addr, px.Config.AllowLocalhost) {
		return ErrInvalidAddress
	}

	return px.setPeerIsValid(addr, valid)
}

// AddPeers add multiple peers at once. Any errors will be logged, but not returned
// Returns the number of peers that were added without error.  Note that
// adding a duplicate peer will not cause an error.
func (px *Pex) AddPeers(addrs []string) int {
	return px.addPeers(addrs, func(addr string) error {
		if !validateAddress(addr, px.Config.AllowLocalhost) {
			return ErrInvalidAddress
		}
		return nil
	})
}

// Run starts the pex service
func (px *Pex) Run() error {
	defer func() {
		// save the peerlist
		logger.Info("Save peerlist")
		px.save(px.Config.DataDirectory)
	}()

	cullTicker := time.NewTicker(px.Config.CullRate)
	clearOldTicker := time.NewTicker(px.Config.ClearOldRate)
	for {
		select {
		case <-cullTicker.C:
			if !px.Config.NetworkDisabled {
				px.cullInvalidPeers()
			}
		case <-clearOldTicker.C:
			// Remove peers we haven't seen in a while
			if !px.Config.Disabled && !px.Config.NetworkDisabled {
				px.clearOld(px.Config.Expiration)
			}
		case <-px.quit:
			return nil
		}
	}
}

// Save persists the peerlist
// func (px *Pex) Save() error {
// 	return px.save(px.cfg.DataDirectory)
// }

// Shutdown notifies the pex service to exist
func (px *Pex) Shutdown() {
	if px.quit != nil {
		close(px.quit)
		px.quit = nil
	}
}
