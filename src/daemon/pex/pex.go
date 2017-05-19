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

	"github.com/skycoin/skycoin/src/util"
)

//TODO:
// - keep track of last time the peer was connected to
// - last time peer was connected to is more important than "seen"
// - peer "seen" means something else than use here
// - save last time connected to, use 0 for never
// - only transmit peers that have active or recent connections

var (
	// Filename for disk-cached peers
	PeerDatabaseFilename = "peers.txt"
	// Filename for disk-cached blacklisted peers
	BlacklistedDatabaseFilename = "blacklisted_peers.txt"
	// Returned when the Pex is at a maximum
	PeerlistFullError = errors.New("Peer list full")
	// Returned when an address appears malformed
	InvalidAddressError = errors.New("Invalid address")
	// Returned when attempting to add a blacklisted peer
	BlacklistedAddressError = errors.New("Blacklisted address")
	// How often to updated expired entries in the blacklist
	RefreshBlacklistRate = time.Second * 30
	// Logging. See http://godoc.org/github.com/op/go-logging for
	// instructions on how to include this log's output
	logger = util.MustGetLogger("pex")
	// Default rng
	rnum = rand.New(rand.NewSource(time.Now().Unix()))
	// For removing inadvertent whitespace from addresses
	whitespaceFilter = regexp.MustCompile("\\s")
)

// Returns true if ipPort is a valid ip:host string
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

// Returns a *Peer initialised by an address string of the form ip:port
func NewPeer(address string) *Peer {
	p := &Peer{Addr: address, Private: false, Trusted: false}
	p.Seen()
	return p
}

// Mark the peer as seen
func (self *Peer) Seen() {
	self.LastSeen = Now()
}

func (self *Peer) IncreaseRetryTimes() {
	self.RetryTimes++
	logger.Info("Increase retry times of %v to %v", self.Addr, self.RetryTimes)
}

func (self *Peer) ResetRetryTimes() {
	self.RetryTimes = 0
	logger.Info("Reset retry times of %v", self.Addr)
}

// CanTry returns whether this peer tryable base on the exponential backoff algorithm
func (self *Peer) CanTry() (rlt bool) {
	// defer func() {
	// 	logger.Info("check if %v can connect, result:%v", self.Addr, rlt)
	// }()
	now := Now()
	mod := (math.Exp2(float64(self.RetryTimes)) - 1) * 5
	if mod == 0 {
		rlt = true
		return
	}

	t := rnum.Int63n(int64(mod))
	timePass := now.Sub(self.LastSeen).Seconds()
	rlt = int64(timePass) > t
	return
}

func (self *Peer) String() string {
	return self.Addr
}

// BlacklistEntry records when an address was blacklisted and how long
// it should be blacklisted for. A duration of 0 is permanent.
type BlacklistEntry struct {
	Start    time.Time
	Duration time.Duration
}

// Returns the time.Time the BlacklistEntry expires
func (b BlacklistEntry) ExpiresAt() time.Time {
	return b.Start.Add(b.Duration)
}

func NewBlacklistEntry(duration time.Duration) BlacklistEntry {
	return BlacklistEntry{Start: Now(), Duration: duration}
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
func (self *Peerlist) GetPublicTrustPeers() []*Peer {
	var peers []*Peer
	self.strand(func() {
		keys := self.getTrustAddresses(false)
		peers = make([]*Peer, len(keys))
		for i, key := range keys {
			peers[i] = self.peers[key]
		}
	}, "GetPublickTrustPeers")
	return peers
}

// GetPrivateTrustPeers returns all trusted private peers
func (self *Peerlist) GetPrivateTrustPeers() []*Peer {
	var peers []*Peer
	self.strand(func() {
		keys := self.getTrustAddresses(true)
		peers = make([]*Peer, len(keys))
		for i, key := range keys {
			peers[i] = self.peers[key]
		}
	}, "GetPrivateTrustPeers")
	return peers
}

// GetAllTrustedPeers returns all trusted peers, including private and public peers.
func (self *Peerlist) GetAllTrustedPeers() []*Peer {
	var peers []*Peer
	self.strand(func() {
		keys := self.getAllTrustPeers()
		peers = make([]*Peer, len(keys))
		for i, key := range keys {
			peers[i] = self.peers[key]
		}
	}, "GetAllTrustedPeers")
	return peers
}

func (self *Peerlist) getTrustAddresses(private bool) []string {
	keys := []string{}
	for key, p := range self.peers {
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

func (self *Peerlist) getAllTrustPeers() []string {
	return append(self.getTrustAddresses(false), self.getTrustAddresses(true)...)
}

// Returns the string addresses of all public peers
func (self *Peerlist) GetPublicAddresses() []string {
	var addrs []string
	self.strand(func() {
		addrs = self.getAddresses(false)
	}, "GetPublicAddresses")
	return addrs
}

// Returns the string addresses of all private peers
func (self *Peerlist) GetPrivateAddresses() []string {
	var addrs []string
	self.strand(func() {
		addrs = self.getAddresses(true)
	}, "GetPrivateAddresses")
	return addrs
}

func (self *Peerlist) RemovePeer(a string) {
	self.strand(func() {
		delete(self.peers, a)
	}, "RemovePeer")
}

// Returns the string addresses of all peers, public or private
func (self *Peerlist) GetAllAddresses() []string {
	var addrs []string
	self.strand(func() {
		addrs = append(self.getAddresses(false), self.getAddresses(true)...)
	}, "GetAllAddreses")
	return addrs
}

func (self *Peerlist) GetPeerByAddr(a string) (Peer, bool) {
	var peer Peer
	var exist bool
	self.strand(func() {
		if p, ok := self.peers[a]; ok {
			peer = *p
			exist = true
			return
		}
	}, "GetPeerByAddr")
	return peer, exist
}

// Removes public peers that haven't been seen in timeAgo seconds
func (self *Peerlist) ClearOld(timeAgo time.Duration) {
	t := Now()
	self.strand(func() {
		for addr, peer := range self.peers {
			if !peer.Private && t.Sub(peer.LastSeen) > timeAgo {
				delete(self.peers, addr)
			}
		}
	}, "ClearOld")
}

// Returns the string addresses of all public peers
func (self *Peerlist) getAddresses(private bool) []string {
	keys := make([]string, 0, len(self.peers))
	for key, p := range self.peers {
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
func (self *Peerlist) random(count int, includePrivate bool) []*Peer {
	keys := []string(nil)
	if includePrivate {
		keys = append(self.getAddresses(true), self.getAddresses(false)...)
	} else {
		keys = self.getAddresses(false)
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
		peers = append(peers, self.peers[keys[i]])
	}
	return peers
}

func (self *Peerlist) getExchgAddr(private bool) []string {
	keys := []string{}
	for a, p := range self.peers {
		if p.HasIncomePort && p.Private == private {
			keys = append(keys, a)
		}
	}
	return keys
}

// returns all exchangeable addresses
func (self *Peerlist) getAllExchgAddr() []string {
	return append(self.getExchgAddr(true), self.getExchgAddr(false)...)
}

// returns n random exchangeable peers, return all if count is 0.
func (self *Peerlist) randomExchg(count int, includePrivate bool) []*Peer {
	keys := []string{}
	if includePrivate {
		keys = self.getAllExchgAddr()
	} else {
		keys = self.getExchgAddr(false)
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
		peers = append(peers, self.peers[keys[i]])
	}
	return peers
}

// RandomExchgPublic returns n random exchangeable public peers
// return all exchangeable public peers if count is 0.
func (self *Peerlist) RandomExchgPublic(count int) []*Peer {
	var peers []*Peer
	self.strand(func() {
		peers = self.randomExchg(count, false)
	}, "RandomExchgPublic")
	return peers
}

// RandomExchgAll returns n random exchangeable peers, including private peers.
// return all exchangeable peers if count is 0.
func (self *Peerlist) RandomExchgAll(count int) []*Peer {
	var peers []*Peer
	self.strand(func() {
		peers = self.randomExchg(count, true)
	}, "RandomExchgAll")
	return peers
}

// Returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.  Will not include
// private peers.
func (self *Peerlist) RandomPublic(count int) []*Peer {
	var peers []*Peer
	self.strand(func() {
		peers = self.random(count, false)
	}, "RandomPublic")
	return peers
}

// Returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.  Includes private
// peers.
func (self *Peerlist) RandomAll(count int) []*Peer {
	var peers []*Peer
	self.strand(func() {
		peers = self.random(count, true)
	}, "RandomAll")
	return peers
}

// Saves known peers to disk as a newline delimited list of addresses to
// <dir><PeerDatabaseFilename>
func (self *Peerlist) Save(dir string) (err error) {
	logger.Debug("PEX: SavingPeerList")
	filename := PeerDatabaseFilename
	fn := filepath.Join(dir, filename)
	self.strand(func() {
		// filter the peers that has retrytime > 10
		peers := make(map[string]*Peer)
		for k, p := range self.peers {
			if p.RetryTimes <= 10 {
				peers[k] = p
			}
		}
		err = util.SaveJSON(fn, peers, 0600)
		if err != nil {
			logger.Notice("SavePeerList Failed: %s", err)
		}
	}, "Save")
	return
}

func (self *Peerlist) IncreaseRetryTimes(addr string) {
	self.strand(func() {
		if _, ok := self.peers[addr]; ok {
			self.peers[addr].IncreaseRetryTimes()
			self.peers[addr].Seen()
		}
	}, "IncreaseRetryTimes")
}

func (self *Peerlist) ResetRetryTimes(addr string) {
	self.strand(func() {
		if _, ok := self.peers[addr]; ok {
			self.peers[addr].ResetRetryTimes()
			self.peers[addr].Seen()
		}
	}, "ResetRetryTimes")
}

// ResetAllRetryTimes reset all peers' retry times
func (self *Peerlist) ResetAllRetryTimes() {
	self.strand(func() {
		for _, p := range self.peers {
			p.ResetRetryTimes()
		}
	})
}

// Loads a newline delimited list of addresses from
// "<dir>/<PeerDatabaseFilename>"
func LoadPeerlist(dir string) (*Peerlist, error) {
	peerlist := Peerlist{peers: make(map[string]*Peer)}
	fn := filepath.Join(dir, PeerDatabaseFilename)
	if err := util.LoadJSON(fn, &peerlist.peers); err != nil {
		return nil, err
	}
	// if err != nil {
	// 	logger.Notice("LoadPeerList Failed: %s", err)
	// }
	return &peerlist, nil

}

// Pex manages a set of known peers and controls peer acquisition
type Pex struct {
	// All known peers
	Peerlist
	// Ignored peers
	// Blacklist Blacklist
	// If false, localhost peers will be rejected from the peerlist
	AllowLocalhost bool
	maxPeers       int
}

func NewPex(maxPeers int) *Pex {
	return &Pex{
		Peerlist: Peerlist{peers: make(map[string]*Peer, maxPeers)},
		// Blacklist:      make(Blacklist, 0),
		maxPeers:       maxPeers,
		AllowLocalhost: false,
	}
}

// AddPeer adds a peer to the peer list, given an address. If the peer list is
// full, PeerlistFullError is returned */
func (px *Pex) AddPeer(ip string) (*Peer, error) {
	if !ValidateAddress(ip, px.AllowLocalhost) {
		return nil, InvalidAddressError
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
			err = PeerlistFullError
		} else {
			peer := NewPeer(ip)
			px.peers[ip] = peer
			p = *peer
		}
	}, "AddPeer")
	return &p, err
}

// SetTrustState updates the peer's Trusted statue
func (px *Pex) SetTrustState(addr string, trusted bool) error {
	if !ValidateAddress(addr, px.AllowLocalhost) {
		return InvalidAddressError
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
		return InvalidAddressError
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

// Load loads both the normal peer and blacklisted peer databases
func (px *Pex) Load(dir string) error {
	peerlist, err := LoadPeerlist(dir)
	if err != nil {
		return err
	}
	px.Peerlist = *peerlist
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

func Now() time.Time {
	return time.Now().UTC()
}
