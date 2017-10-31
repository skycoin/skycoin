// Package pex is a toolkit for implementing a peer exchange system
package pex

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
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
	Addr            string // An address of the form ip:port
	LastSeen        int64  // Unix timestamp when this peer was last seen
	Private         bool   // Whether it should omitted from public requests
	Trusted         bool   // Whether this peer is trusted
	HasIncomingPort bool   // Whether this peer has incoming port
	RetryTimes      int    // Records the retry times
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
	logger.Debug("Increase retry times of %v to %v", peer.Addr, peer.RetryTimes)
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

// Peerlist is a map of addresses to *PeerStates
type Peerlist struct {
	lock  sync.RWMutex
	peers map[string]*Peer
}

// GetPublicTrustPeers returns all trusted public peers
func (pl *Peerlist) GetPublicTrustPeers() []Peer {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	addrs := pl.getTrustedAddresses(false)
	return copyPeers(pl.peers, addrs)
}

// GetPrivateTrustPeers returns all trusted private peers
func (pl *Peerlist) GetPrivateTrustPeers() []Peer {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	addrs := pl.getTrustedAddresses(true)
	return copyPeers(pl.peers, addrs)
}

// GetAllTrustedPeers returns all trusted peers, including private and public peers.
func (pl *Peerlist) GetAllTrustedPeers() []Peer {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	addrs := pl.getAllTrustedAddresses()
	return copyPeers(pl.peers, addrs)
}

// Returns a copy of peers from addrs. Called must hold a read lock.
func copyPeers(peerlist map[string]*Peer, addrs []string) []Peer {
	peers := make([]Peer, len(addrs))
	for i, addr := range addrs {
		peers[i] = *peerlist[addr]
	}
	return peers
}

// Returns trusted peers matching Private flag. Caller must hold a read lock.
func (pl *Peerlist) getTrustedAddresses(private bool) []string {
	addrs := []string{}
	for addr, p := range pl.peers {
		if !p.Trusted {
			continue
		}

		if !p.CanTry() {
			continue
		}

		if private && p.Private {
			addrs = append(addrs, addr)
		} else if !private && !p.Private {
			addrs = append(addrs, addr)
		}
	}
	return addrs
}

// Returns all trusted peers, private or not. Caller must hold a read lock.
func (pl *Peerlist) getAllTrustedAddresses() []string {
	return append(pl.getTrustedAddresses(false), pl.getTrustedAddresses(true)...)
}

// GetPublicAddresses returns the string addresses of all public peers
func (pl *Peerlist) GetPublicAddresses() []string {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	return pl.getAddresses(false)
}

// GetPrivateAddresses returns the string addresses of all private peers
func (pl *Peerlist) GetPrivateAddresses() []string {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	return pl.getAddresses(true)
}

// RemovePeer removes peer
func (pl *Peerlist) RemovePeer(a string) {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	delete(pl.peers, a)
}

// GetAllAddresses returns the string addresses of all peers, public or private
func (pl *Peerlist) GetAllAddresses() []string {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	return append(pl.getAddresses(false), pl.getAddresses(true)...)
}

// GetPeerByAddr returns peer of given address
func (pl *Peerlist) GetPeerByAddr(a string) (Peer, bool) {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	if p, ok := pl.peers[a]; ok {
		return *p, true
	}

	return Peer{}, false
}

// ClearOld removes public peers that haven't been seen in timeAgo seconds
func (pl *Peerlist) ClearOld(timeAgo time.Duration) {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	t := utc.Now()

	for addr, peer := range pl.peers {
		lastSeen := time.Unix(peer.LastSeen, 0)
		if !peer.Private && t.Sub(lastSeen) > timeAgo {
			delete(pl.peers, addr)
		}
	}
}

// Returns the string addresses of all public peers.  Called must hold a read lock.
func (pl *Peerlist) getAddresses(private bool) []string {
	keys := make([]string, 0, len(pl.peers))
	for key, p := range pl.peers {
		if !p.CanTry() {
			continue
		}

		if private && p.Private {
			keys = append(keys, key)
		} else if !private && !p.Private {
			keys = append(keys, key)
		}
	}

	return keys
}

// Returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.
// Caller must hold a read lock.
func (pl *Peerlist) random(count int, includePrivate bool) []Peer {
	keys := pl.getAddresses(false)
	if includePrivate {
		keys = append(keys, pl.getAddresses(true)...)
	}

	if len(keys) == 0 {
		return []Peer{}
	}

	max := count
	if count == 0 || count > len(keys) {
		max = len(keys)
	}

	peers := make([]Peer, 0, max)
	perm := rand.Perm(len(keys))
	for _, i := range perm[:max] {
		peers = append(peers, *pl.peers[keys[i]])
	}

	return peers
}

func (pl *Peerlist) getExchgAddr(private bool) []string {
	keys := []string{}
	for a, p := range pl.peers {
		if p.HasIncomingPort && p.Private == private {
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
func (pl *Peerlist) randomExchg(count int, includePrivate bool) []Peer {
	keys := []string{}

	if includePrivate {
		keys = pl.getAllExchgAddr()
	} else {
		keys = pl.getExchgAddr(false)
	}

	if len(keys) == 0 {
		return make([]Peer, 0)
	}

	max := count
	if count == 0 || count > len(keys) {
		max = len(keys)
	}

	peers := make([]Peer, 0, max)
	perm := rand.Perm(len(keys))
	for _, i := range perm[:max] {
		peers = append(peers, *pl.peers[keys[i]])
	}

	return peers
}

// RandomExchgPublic returns n random exchangeable public peers
// return all exchangeable public peers if count is 0.
func (pl *Peerlist) RandomExchgPublic(count int) []Peer {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	return pl.randomExchg(count, false)
}

// RandomExchgAll returns n random exchangeable peers, including private peers.
// return all exchangeable peers if count is 0.
func (pl *Peerlist) RandomExchgAll(count int) []Peer {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	return pl.randomExchg(count, true)
}

// RandomPublic returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.  Will not include
// private peers.
func (pl *Peerlist) RandomPublic(count int) []Peer {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	return pl.random(count, false)
}

// RandomAll returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.  Includes private
// peers.
func (pl *Peerlist) RandomAll(count int) []Peer {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	return pl.random(count, true)
}

// PeerJSON is for saving and loading peers to disk. Some fields are strange,
// to be backwards compatible due to variable name changes
type PeerJSON struct {
	Addr string // An address of the form ip:port
	// Unix timestamp when this peer was last seen.
	// This could be a time.Time string or an int64 timestamp
	LastSeen        interface{}
	Private         bool  // Whether it should omitted from public requests
	Trusted         bool  // Whether this peer is trusted
	HasIncomePort   *bool `json:"HasIncomePort,omitempty"` // Whether this peer has incoming port [DEPRECATED]
	HasIncomingPort *bool // Whether this peer has incoming port
}

// NewPeerJSON returns a PeerJSON from a Peer
func NewPeerJSON(p Peer) PeerJSON {
	hasIncomingPort := p.HasIncomingPort

	return PeerJSON{
		Addr:            p.Addr,
		LastSeen:        p.LastSeen,
		Private:         p.Private,
		Trusted:         p.Trusted,
		HasIncomingPort: &hasIncomingPort,
	}
}

// NewPeerFromJSON converts a PeerJSON to a Peer
func NewPeerFromJSON(p PeerJSON) (Peer, error) {
	hasIncomingPort := false
	if p.HasIncomingPort != nil {
		hasIncomingPort = *p.HasIncomingPort
	} else if p.HasIncomePort != nil {
		hasIncomingPort = *p.HasIncomePort
	}

	// LastSeen could be a RFC3339Nano timestamp or an int64 unix timestamp
	var lastSeen int64
	switch p.LastSeen.(type) {
	case string:
		lastSeenTime, err := time.Parse(time.RFC3339Nano, p.LastSeen.(string))
		if err != nil {
			return Peer{}, err
		}
		lastSeen = lastSeenTime.Unix()
	case json.Number:
		lastSeenNum := p.LastSeen.(json.Number)
		var err error
		lastSeen, err = lastSeenNum.Int64()
		if err != nil {
			return Peer{}, err
		}
	default:
		return Peer{}, fmt.Errorf("Invalid type %T for LastSeen field", p.LastSeen)
	}

	return Peer{
		Addr:            p.Addr,
		LastSeen:        lastSeen,
		Private:         p.Private,
		Trusted:         p.Trusted,
		HasIncomingPort: hasIncomingPort,
	}, nil
}

// Save saves known peers to disk as a newline delimited list of addresses to
// <dir><PeerDatabaseFilename>
func (pl *Peerlist) Save(dir string) error {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	filename := PeerDatabaseFilename
	fn := filepath.Join(dir, filename)

	// filter the peers that has retrytime > 10
	peers := make(map[string]PeerJSON)
	for k, p := range pl.peers {
		if p.RetryTimes <= 10 {
			peers[k] = NewPeerJSON(*p)
		}
	}

	err := file.SaveJSON(fn, peers, 0600)
	if err != nil {
		logger.Notice("SavePeerList Failed: %v", err)
	}

	return err
}

// IncreaseRetryTimes increases retry times
func (pl *Peerlist) IncreaseRetryTimes(addr string) {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	if _, ok := pl.peers[addr]; ok {
		pl.peers[addr].IncreaseRetryTimes()
		pl.peers[addr].Seen()
	}
}

// ResetRetryTimes reset retry times
func (pl *Peerlist) ResetRetryTimes(addr string) {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	if _, ok := pl.peers[addr]; ok {
		pl.peers[addr].ResetRetryTimes()
		pl.peers[addr].Seen()
	}
}

// ResetAllRetryTimes reset all peers' retry times
func (pl *Peerlist) ResetAllRetryTimes() {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	logger.Info("Reset all peer's retry times")
	for _, p := range pl.peers {
		p.ResetRetryTimes()
	}
}

// LoadPeerlist loads a newline delimited list of addresses from
// "<dir>/<PeerDatabaseFilename>"
func LoadPeerlist(dir string) (*Peerlist, error) {
	peersJSON := make(map[string]PeerJSON)

	fn := filepath.Join(dir, PeerDatabaseFilename)
	if err := file.LoadJSON(fn, &peersJSON); err != nil {
		return nil, err
	}

	peers := make(map[string]*Peer, len(peersJSON))
	for addr, peerJSON := range peersJSON {
		peer, err := NewPeerFromJSON(peerJSON)
		if err != nil {
			return nil, err
		}
		peers[addr] = &peer
	}

	return &Peerlist{
		peers: peers,
	}, nil
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
		Peerlist: &Peerlist{
			peers: make(map[string]*Peer, maxPeers),
		},
		maxPeers:       maxPeers,
		AllowLocalhost: false,
	}
}

// AddPeer adds a peer to the peer list, given an address. If the peer list is
// full, PeerlistFullError is returned
func (px *Pex) AddPeer(ip string) (Peer, error) {
	if !ValidateAddress(ip, px.AllowLocalhost) {
		return Peer{}, ErrInvalidAddress
	}

	px.lock.Lock()
	defer px.lock.Unlock()

	return px.addPeer(ip)
}

// addPeer adds a peer to the peer list, given an address. If the peer list is
// full, PeerlistFullError is returned.
// Caller must hold a write lock.
func (px *Pex) addPeer(ip string) (Peer, error) {
	peer, ok := px.peers[ip]

	if peer == nil && ok {
		logger.Warning("PeerList contains a nil peer with address %s", ip)
	}

	if peer != nil {
		return *peer, nil
	} else if px.full() {
		return Peer{}, ErrPeerlistFull
	}

	peer = NewPeer(ip)
	px.peers[ip] = peer
	return *peer, nil
}

// SetPrivate updates the private value of given ip in peerlist
func (px *Pex) SetPrivate(ip string, private bool) error {
	if !ValidateAddress(ip, px.AllowLocalhost) {
		return ErrInvalidAddress
	}

	px.lock.Lock()
	defer px.lock.Unlock()

	if p, ok := px.peers[ip]; ok {
		p.Private = private
		return nil
	}

	return fmt.Errorf("Set peer.Private failed: %v does not exist in peerlist", ip)
}

// SetTrustState updates the peer's Trusted statue
func (px *Pex) SetTrustState(addr string, trusted bool) error {
	if !ValidateAddress(addr, px.AllowLocalhost) {
		return ErrInvalidAddress
	}

	px.lock.Lock()
	defer px.lock.Unlock()

	if p, ok := px.peers[addr]; ok {
		p.Trusted = trusted
		return nil
	}

	return fmt.Errorf("%s does not exist in peel list", addr)

}

// SetPeerHasIncomingPort update whether the peer has incoming port.
func (px *Pex) SetPeerHasIncomingPort(addr string, v bool) error {
	if !ValidateAddress(addr, px.AllowLocalhost) {
		return ErrInvalidAddress
	}

	px.lock.Lock()
	defer px.lock.Unlock()

	if p, ok := px.peers[addr]; ok {
		p.HasIncomingPort = v
		p.Seen()
		return nil
	}

	return fmt.Errorf("peer %s is not in exchange peer list", addr)
}

// Full returns true if no more peers can be added
func (px *Pex) Full() bool {
	px.lock.RLock()
	defer px.lock.RUnlock()

	return px.full()
}

// Returns true if the peer list is full. Called must hold a read lock on the peer list.
func (px *Pex) full() bool {
	return px.maxPeers > 0 && len(px.peers) >= px.maxPeers
}

// AddPeers add multiple peers at once.
// Any errors will be logged, but not returned.
func (px *Pex) AddPeers(peers []string) {
	px.lock.Lock()
	defer px.lock.Unlock()

	for _, p := range peers {
		if _, err := px.addPeer(p); err != nil {
			logger.Warning("Failed to add peer %s, Reason: %v", p, err)
		}
	}
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
