// Package pex is a toolkit for implementing a peer exchange system
package pex

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"math"

	"sync"

	"github.com/skycoin/skycoin/src/util/file"
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
)

// ValidateAddress returns true if ipPort is a valid ip:host string
func ValidateAddress(ipPort string, allowLocalhost bool) bool {
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
	Addr          string    // An address of the form ip:port
	LastSeen      time.Time // Unix timestamp when this peer was last seen
	Private       bool      // Whether it should omitted from public requests
	Trusted       bool      // Whether this peer is trusted
	HasIncomePort bool      // Whether this peer has incomming port
	RetryTimes    int       `json:"-"` // records the retry times
}

// NewPeer returns a *Peer initialised by an address string of the form ip:port
func NewPeer(address string) *Peer {
	p := &Peer{Addr: address, Private: false, Trusted: false}
	p.Seen()
	return p
}

// Seen marks the peer as seen
func (peer *Peer) Seen() {
	peer.LastSeen = Now()
}

// IncreaseRetryTimes adds the retry times
func (peer *Peer) IncreaseRetryTimes() {
	peer.RetryTimes++
	logger.Debug("Increase retry times of %v to %v", peer.Addr, peer.RetryTimes)
}

// ResetRetryTimes resets the retry time
func (peer *Peer) ResetRetryTimes() {
	peer.RetryTimes = 0
	logger.Debug("Reset retry times of %v", peer.Addr)
}

// CanTry returns whether this peer is tryable base on the exponential backoff algorithm
func (peer *Peer) CanTry() (rlt bool) {
	now := Now()
	mod := (math.Exp2(float64(peer.RetryTimes)) - 1) * 5
	if mod == 0 {
		rlt = true
		return
	}

	t := rnum.Int63n(int64(mod))
	timePass := now.Sub(peer.LastSeen).Seconds()
	rlt = int64(timePass) > t
	return
}

// String returns the peer address
func (peer *Peer) String() string {
	return peer.Addr
}

// Peerlist is a map of addresses to *PeerStates
type Peerlist struct {
	lock  sync.Mutex
	peers map[string]*Peer
}

func (pl *Peerlist) strand(f func(), arg ...interface{}) {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	// logger.Critical("%v", arg)
	f()
}

// GetPublicTrustPeers returns all trusted public peers
func (pl *Peerlist) GetPublicTrustPeers() []*Peer {
	var peers []*Peer
	pl.strand(func() {
		keys := pl.getTrustAddresses(false)
		peers = make([]*Peer, len(keys))
		for i, key := range keys {
			peers[i] = pl.peers[key]
		}
	}, "GetPublickTrustPeers")
	return peers
}

// GetPrivateTrustPeers returns all trusted private peers
func (pl *Peerlist) GetPrivateTrustPeers() []*Peer {
	var peers []*Peer
	pl.strand(func() {
		keys := pl.getTrustAddresses(true)
		peers = make([]*Peer, len(keys))
		for i, key := range keys {
			peers[i] = pl.peers[key]
		}
	}, "GetPrivateTrustPeers")
	return peers
}

// GetAllTrustedPeers returns all trusted peers, including private and public peers.
func (pl *Peerlist) GetAllTrustedPeers() []*Peer {
	var peers []*Peer
	pl.strand(func() {
		keys := pl.getAllTrustPeers()
		peers = make([]*Peer, len(keys))
		for i, key := range keys {
			peers[i] = pl.peers[key]
		}
	}, "GetAllTrustedPeers")
	return peers
}

func (pl *Peerlist) getTrustAddresses(private bool) []string {
	keys := []string{}
	for key, p := range pl.peers {
		if p.Trusted {
			if p.CanTry() {
				if private && p.Private {
					keys = append(keys, key)
				} else if !private && !p.Private {
					keys = append(keys, key)
				}
			}
		}
	}
	return keys
}

func (pl *Peerlist) getAllTrustPeers() []string {
	return append(pl.getTrustAddresses(false), pl.getTrustAddresses(true)...)
}

// GetPublicAddresses returns the string addresses of all public peers
func (pl *Peerlist) GetPublicAddresses() []string {
	var addrs []string
	pl.strand(func() {
		addrs = pl.getAddresses(false)
	}, "GetPublicAddresses")
	return addrs
}

// GetPrivateAddresses returns the string addresses of all private peers
func (pl *Peerlist) GetPrivateAddresses() []string {
	var addrs []string
	pl.strand(func() {
		addrs = pl.getAddresses(true)
	}, "GetPrivateAddresses")
	return addrs
}

// RemovePeer removes peer
func (pl *Peerlist) RemovePeer(a string) {
	pl.strand(func() {
		delete(pl.peers, a)
	}, "RemovePeer")
}

// GetAllAddresses returns the string addresses of all peers, public or private
func (pl *Peerlist) GetAllAddresses() []string {
	var addrs []string
	pl.strand(func() {
		addrs = append(pl.getAddresses(false), pl.getAddresses(true)...)
	}, "GetAllAddreses")
	return addrs
}

// GetPeerByAddr returns peer of given address
func (pl *Peerlist) GetPeerByAddr(a string) (Peer, bool) {
	var peer Peer
	var exist bool
	pl.strand(func() {
		if p, ok := pl.peers[a]; ok {
			peer = *p
			exist = true
			return
		}
	}, "GetPeerByAddr")
	return peer, exist
}

// ClearOld removes public peers that haven't been seen in timeAgo seconds
func (pl *Peerlist) ClearOld(timeAgo time.Duration) {
	t := Now()
	pl.strand(func() {
		for addr, peer := range pl.peers {
			if !peer.Private && t.Sub(peer.LastSeen) > timeAgo {
				delete(pl.peers, addr)
			}
		}
	}, "ClearOld")
}

// Returns the string addresses of all public peers
func (pl *Peerlist) getAddresses(private bool) []string {
	keys := make([]string, 0, len(pl.peers))
	for key, p := range pl.peers {
		if p.CanTry() {
			if private && p.Private {
				keys = append(keys, key)
			} else if !private && !p.Private {
				keys = append(keys, key)
			}
		}
	}

	return keys
}

// Returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.
func (pl *Peerlist) random(count int, includePrivate bool) []*Peer {
	keys := []string{}
	if includePrivate {
		keys = append(pl.getAddresses(true), pl.getAddresses(false)...)
	} else {
		keys = pl.getAddresses(false)
	}
	if len(keys) == 0 {
		return make([]*Peer, 0)
	}
	max := count
	if count == 0 || count > len(keys) {
		max = len(keys)
	}
	peers := make([]*Peer, 0, max)
	perm := rand.Perm(len(keys))
	for _, i := range perm[:max] {
		peers = append(peers, pl.peers[keys[i]])
	}
	return peers
}

func (pl *Peerlist) getExchgAddr(private bool) []string {
	keys := []string{}
	for a, p := range pl.peers {
		if p.HasIncomePort && p.Private == private {
			keys = append(keys, a)
		}
	}
	return keys
}

// returns all exchangeable addresses
func (pl *Peerlist) getAllExchgAddr() []string {
	return append(pl.getExchgAddr(true), pl.getExchgAddr(false)...)
}

// returns n random exchangeable peers, return all if count is 0.
func (pl *Peerlist) randomExchg(count int, includePrivate bool) []*Peer {
	keys := []string{}
	if includePrivate {
		keys = pl.getAllExchgAddr()
	} else {
		keys = pl.getExchgAddr(false)
	}

	if len(keys) == 0 {
		return make([]*Peer, 0)
	}

	max := count
	if count == 0 || count > len(keys) {
		max = len(keys)
	}
	peers := make([]*Peer, 0, max)
	perm := rand.Perm(len(keys))
	for _, i := range perm[:max] {
		peers = append(peers, pl.peers[keys[i]])
	}
	return peers
}

// RandomExchgPublic returns n random exchangeable public peers
// return all exchangeable public peers if count is 0.
func (pl *Peerlist) RandomExchgPublic(count int) []*Peer {
	var peers []*Peer
	pl.strand(func() {
		peers = pl.randomExchg(count, false)
	}, "RandomExchgPublic")
	return peers
}

// RandomExchgAll returns n random exchangeable peers, including private peers.
// return all exchangeable peers if count is 0.
func (pl *Peerlist) RandomExchgAll(count int) []*Peer {
	var peers []*Peer
	pl.strand(func() {
		peers = pl.randomExchg(count, true)
	}, "RandomExchgAll")
	return peers
}

// RandomPublic returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.  Will not include
// private peers.
func (pl *Peerlist) RandomPublic(count int) []*Peer {
	var peers []*Peer
	pl.strand(func() {
		peers = pl.random(count, false)
	}, "RandomPublic")
	return peers
}

// RandomAll returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.  Includes private
// peers.
func (pl *Peerlist) RandomAll(count int) []*Peer {
	var peers []*Peer
	pl.strand(func() {
		peers = pl.random(count, true)
	}, "RandomAll")
	return peers
}

// Save saves known peers to disk as a newline delimited list of addresses to
// <dir><PeerDatabaseFilename>
func (pl *Peerlist) Save(dir string) (err error) {
	filename := PeerDatabaseFilename
	fn := filepath.Join(dir, filename)
	pl.strand(func() {
		// filter the peers that has retrytime > 10
		peers := make(map[string]*Peer)
		for k, p := range pl.peers {
			if p.RetryTimes <= 10 {
				peers[k] = p
			}
		}
		err = file.SaveJSON(fn, peers, 0600)
		if err != nil {
			logger.Notice("SavePeerList Failed: %s", err)
		}
	}, "Save")
	return
}

// IncreaseRetryTimes increases retry times
func (pl *Peerlist) IncreaseRetryTimes(addr string) {
	pl.strand(func() {
		if _, ok := pl.peers[addr]; ok {
			pl.peers[addr].IncreaseRetryTimes()
			pl.peers[addr].Seen()
		}
	}, "IncreaseRetryTimes")
}

// ResetRetryTimes reset retry times
func (pl *Peerlist) ResetRetryTimes(addr string) {
	pl.strand(func() {
		if _, ok := pl.peers[addr]; ok {
			pl.peers[addr].ResetRetryTimes()
			pl.peers[addr].Seen()
		}
	}, "ResetRetryTimes")
}

// ResetAllRetryTimes reset all peers' retry times
func (pl *Peerlist) ResetAllRetryTimes() {
	logger.Info("Reset all peer's retry times")
	pl.strand(func() {
		for _, p := range pl.peers {
			p.ResetRetryTimes()
		}
	})
}

// LoadPeerlist loads a newline delimited list of addresses from
// "<dir>/<PeerDatabaseFilename>"
func LoadPeerlist(dir string) (*Peerlist, error) {
	peerlist := Peerlist{peers: make(map[string]*Peer)}
	fn := filepath.Join(dir, PeerDatabaseFilename)
	if err := file.LoadJSON(fn, &peerlist.peers); err != nil {
		return nil, err
	}
	return &peerlist, nil

}

// Pex manages a set of known peers and controls peer acquisition
type Pex struct {
	// All known peers
	*Peerlist
	// If false, localhost peers will be rejected from the peerlist
	AllowLocalhost bool
	maxPeers       int
}

// NewPex creates pex
func NewPex(maxPeers int) *Pex {
	return &Pex{
		Peerlist:       &Peerlist{peers: make(map[string]*Peer, maxPeers)},
		maxPeers:       maxPeers,
		AllowLocalhost: false,
	}
}

// AddPeer adds a peer to the peer list, given an address. If the peer list is
// full, PeerlistFullError is returned */
func (px *Pex) AddPeer(ip string) (*Peer, error) {
	if !ValidateAddress(ip, px.AllowLocalhost) {
		return nil, ErrInvalidAddress
	}
	var p Peer
	var err error
	px.Peerlist.strand(func() {
		peer := px.peers[ip]
		if peer != nil {
			peer.Seen()
			p = *peer
			return
		} else if px.full() {
			err = ErrPeerlistFull
		} else {
			peer := NewPeer(ip)
			px.peers[ip] = peer
			p = *peer
		}
	}, "AddPeer")
	return &p, err
}

// SetPrivate updates the private value of given ip in peerlist
func (px *Pex) SetPrivate(ip string, private bool) error {
	var err error
	px.Peerlist.strand(func() {
		if p, ok := px.peers[ip]; ok {
			p.Private = private
			return
		}

		err = fmt.Errorf("Set peer.Private failed: %v does not exist in peerlist", ip)
	})
	return err
}

// SetTrustState updates the peer's Trusted statue
func (px *Pex) SetTrustState(addr string, trusted bool) error {
	if !ValidateAddress(addr, px.AllowLocalhost) {
		return ErrInvalidAddress
	}

	var err error
	px.strand(func() {
		if p, ok := px.peers[addr]; ok {
			p.Trusted = trusted
		} else {
			err = fmt.Errorf("%s does not exist in peel list", addr)
		}

	}, "SetTrustState")

	return err
}

// SetPeerHasInPort update whether the peer has incomming port.
func (px *Pex) SetPeerHasInPort(addr string, v bool) error {
	if !ValidateAddress(addr, px.AllowLocalhost) {
		return ErrInvalidAddress
	}

	var err error
	px.strand(func() {
		if p, ok := px.peers[addr]; ok {
			p.HasIncomePort = v
			p.Seen()
		} else {
			err = fmt.Errorf("peer %s is not in exchange peer list", addr)
		}

	}, "SetPeerHasInPort")

	return err
}

// Full returns true if no more peers can be added
func (px *Pex) Full() bool {
	var full bool
	px.strand(func() {
		full = px.full()
	}, "Full")
	return full
}

func (px *Pex) full() bool {
	return px.maxPeers > 0 && len(px.peers) >= px.maxPeers
}

// AddPeers add multiple peers at once. Any errors will be logged, but not returned
// Returns the number of peers that were added without error.  Note that
// adding a duplicate peer will not cause an error.
func (px *Pex) AddPeers(peers []string) int {
	n := len(peers)
	for _, p := range peers {
		_, err := px.AddPeer(p)
		if err != nil {
			logger.Warning("Failed to add peer %s, Reason: %v", p, err)
			n--
		}
	}
	return n
}

// Load loads peers
func (px *Pex) Load(dir string) error {
	pl, err := LoadPeerlist(dir)
	if err != nil {
		return err
	}

	px.Peerlist = pl
	return nil
}

/* Common utilities */

// Reads a file located at dir/filename and splits it on newlines
func readLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	data := make([]byte, info.Size())
	_, err = f.Read(data)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

// Now returns UTC time
func Now() time.Time {
	return utc.Now()
}
