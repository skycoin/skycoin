package pex

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
)

// Peers peer list
type Peers []Peer

// ToAddrs returns the address list
func (ps Peers) ToAddrs() []string {
	addrs := make([]string, 0, len(ps))
	for _, p := range ps {
		addrs = append(addrs, p.Addr)
	}
	return addrs
}

// peerlist is a map of addresses to *PeerStates
type peerlist struct {
	peers map[string]*Peer
}

func newPeerlist() peerlist {
	return peerlist{
		peers: make(map[string]*Peer),
	}
}

// Filter peers filter
type Filter func(peer Peer) bool

// loadCachedPeersFile loads peers from the cached peers.json file
func loadCachedPeersFile(path string) (map[string]*Peer, error) {
	peersJSON := make(map[string]PeerJSON)
	err := file.LoadJSON(path, &peersJSON)

	if os.IsNotExist(err) {
		logger.WithField("path", path).Info("File does not exist")
		return nil, nil
	} else if err == io.EOF {
		logger.WithField("path", path).Error("Corrupt or empty file")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	peers := make(map[string]*Peer, len(peersJSON))
	for addr, peerJSON := range peersJSON {
		a, err := validateAddress(addr, true)

		if err != nil {
			logger.Errorf("Invalid address in peers JSON file %s: %v", addr, err)
			continue
		}

		peer, err := newPeerFromJSON(peerJSON)
		if err != nil {
			logger.Errorf("newPeerFromJSON failed: %v", err)
			continue
		}

		if a != peer.Addr {
			logger.Errorf("Address key %s does not match Peer.Addr %s", a, peer.Addr)
			continue
		}

		peers[a] = peer
	}

	return peers, nil
}

func (pl *peerlist) setPeers(peers []Peer) {
	for _, p := range peers {
		np := p
		pl.peers[p.Addr] = &np
	}
}

func (pl *peerlist) hasPeer(addr string) bool {
	p, ok := pl.peers[addr]
	return ok && p != nil
}

func (pl *peerlist) addPeer(addr string) {
	if p, ok := pl.peers[addr]; ok && p != nil {
		p.Seen()
		return
	}

	peer := NewPeer(addr)
	pl.peers[addr] = peer
}

func (pl *peerlist) addPeers(addrs []string) {
	for _, addr := range addrs {
		pl.addPeer(addr)
	}
}

// getCanTryPeers returns all peers that are triable(retried times blew exponential backoff times)
// and are able to pass the filters.
func (pl *peerlist) getCanTryPeers(flts []Filter) Peers {
	ps := make(Peers, 0)
	flts = append([]Filter{canTry}, flts...)
loop:
	for _, p := range pl.peers {
		for i := range flts {
			if !flts[i](*p) {
				continue loop
			}
		}

		ps = append(ps, *p)
	}

	return ps
}

// getPeers returns all peers that can pass the filters.
func (pl *peerlist) getPeers(flts []Filter) Peers {
	ps := make(Peers, 0)
loop:
	for _, p := range pl.peers {
		for i := range flts {
			if !flts[i](*p) {
				continue loop
			}
		}

		ps = append(ps, *p)
	}

	return ps
}

// filters
func isPrivate(p Peer) bool {
	return p.Private
}

func isPublic(p Peer) bool {
	return !p.Private
}

func isTrusted(p Peer) bool {
	return p.Trusted
}

func hasIncomingPort(p Peer) bool {
	return p.HasIncomingPort
}

func canTry(p Peer) bool {
	return p.CanTry()
}

func zeroRetryTimes(p Peer) bool {
	return p.RetryTimes == 0
}

// isExchangeable filters exchangeable peers
var isExchangeable = []Filter{hasIncomingPort, isPublic, zeroRetryTimes}

// removePeer removes peer
func (pl *peerlist) removePeer(addr string) {
	delete(pl.peers, addr)
}

// SetPrivate sets specific peer as private
func (pl *peerlist) setPrivate(addr string, private bool) error {
	if p, ok := pl.peers[addr]; ok {
		p.Private = private
		return nil
	}

	return fmt.Errorf("set peer.Private failed: %v does not exist in peer list", addr)
}

// SetTrusted sets peer as trusted peer
func (pl *peerlist) setTrusted(addr string, trusted bool) error {
	if p, ok := pl.peers[addr]; ok {
		p.Trusted = trusted
		return nil
	}

	return fmt.Errorf("set peer.Trusted failed: %v does not exist in peer list", addr)
}

// setAllUntrusted unsets the trusted field on all peers
func (pl *peerlist) setAllUntrusted() {
	for _, p := range pl.peers {
		p.Trusted = false
	}
}

// setHasIncomingPort updates whether the peer is valid and has public incoming port
func (pl *peerlist) setHasIncomingPort(addr string, hasIncomingPort bool) error {
	if p, ok := pl.peers[addr]; ok {
		p.HasIncomingPort = hasIncomingPort
		p.Seen()
		return nil
	}

	return fmt.Errorf("set peer.HasIncomingPort failed: %v does not exist in peer list", addr)
}

// len returns number of peers
func (pl *peerlist) len() int {
	return len(pl.peers)
}

// getPeerByAddr returns peer of given address
func (pl *peerlist) getPeerByAddr(addr string) (Peer, bool) {
	p, ok := pl.peers[addr]
	if ok {
		return *p, true
	}
	return Peer{}, false
}

// ClearOld removes public peers that haven't been seen in timeAgo seconds
func (pl *peerlist) clearOld(timeAgo time.Duration) {
	t := time.Now().UTC()
	for addr, peer := range pl.peers {
		lastSeen := time.Unix(peer.LastSeen, 0)
		if !peer.Private && t.Sub(lastSeen) > timeAgo {
			delete(pl.peers, addr)
		}
	}
}

// Returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.
func (pl *peerlist) random(count int, flts []Filter) Peers {
	keys := pl.getCanTryPeers(flts).ToAddrs()
	if len(keys) == 0 {
		return Peers{}
	}
	max := count
	if count == 0 || count > len(keys) {
		max = len(keys)
	}

	ps := make(Peers, 0)
	perm := rand.Perm(len(keys))
	for _, i := range perm[:max] {
		ps = append(ps, *pl.peers[keys[i]])
	}
	return ps
}

// save saves known peers to disk as a newline delimited list of addresses to
// <dir><PeerCacheFilename>
func (pl *peerlist) save(fn string) error {
	// filter the peers that has retrytime > MaxPeerRetryTimes
	peers := make(map[string]PeerJSON)
	for k, p := range pl.peers {
		if p.RetryTimes <= MaxPeerRetryTimes {
			peers[k] = newPeerJSON(*p)
		}
	}

	if err := file.SaveJSON(fn, peers, 0600); err != nil {
		return fmt.Errorf("save peer list failed: %s", err)
	}
	return nil
}

// increaseRetryTimes increases retry times
func (pl *peerlist) increaseRetryTimes(addr string) {
	if p, ok := pl.peers[addr]; ok {
		p.IncreaseRetryTimes()
		p.Seen()
	}
}

// resetRetryTimes reset retry times
func (pl *peerlist) resetRetryTimes(addr string) {
	if p, ok := pl.peers[addr]; ok {
		p.ResetRetryTimes()
		p.Seen()
	}
}

// resetAllRetryTimes reset all peers' retry times
func (pl *peerlist) resetAllRetryTimes() {
	logger.Info("Reset all peer's retry times")
	for _, p := range pl.peers {
		p.ResetRetryTimes()
	}
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

// newPeerJSON returns a PeerJSON from a Peer
func newPeerJSON(p Peer) PeerJSON {
	return PeerJSON{
		Addr:            p.Addr,
		LastSeen:        p.LastSeen,
		Private:         p.Private,
		Trusted:         p.Trusted,
		HasIncomingPort: &p.HasIncomingPort,
	}
}

// newPeerFromJSON converts a PeerJSON to a Peer
func newPeerFromJSON(p PeerJSON) (*Peer, error) {
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
			return nil, err
		}
		lastSeen = lastSeenTime.Unix()
	case json.Number:
		lastSeenNum := p.LastSeen.(json.Number)
		var err error
		lastSeen, err = lastSeenNum.Int64()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Invalid type %T for LastSeen field", p.LastSeen)
	}

	addr, err := validateAddress(p.Addr, true)
	if err != nil {
		return nil, err
	}

	return &Peer{
		Addr:            addr,
		LastSeen:        lastSeen,
		Private:         p.Private,
		Trusted:         p.Trusted,
		HasIncomingPort: hasIncomingPort,
	}, nil
}
