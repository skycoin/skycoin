package pex

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/utc"
)

// peerlist is a map of addresses to *PeerStates
type peerlist struct {
	lock  sync.Mutex
	peers map[string]*Peer
	cap   int // the max number of peers that the peerlist can take
}

func newPeerlist(maxPeers int) *peerlist {
	return &peerlist{
		peers: make(map[string]*Peer, maxPeers),
	}
}

func (pl *peerlist) load(dir string) error {
	fn := filepath.Join(dir, PeerDatabaseFilename)
	return file.LoadJSON(fn, &pl.peers)
}

func (pl *peerlist) strand(f func(), arg ...interface{}) {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	// logger.Critical("%v", arg)
	f()
}

// Full returns true if no more peers can be added
func (pl *peerlist) Full() bool {
	var full bool
	pl.strand(func() {
		full = pl.full()
	}, "Full")
	return full
}

func (pl *peerlist) full() bool {
	return pl.cap > 0 && len(pl.peers) >= pl.cap
}

func (pl *peerlist) add(addr string) error {
	if p, ok := pl.peers[addr]; ok && p != nil {
		p.Seen()
		return nil
	}

	if pl.full() {
		return ErrPeerlistFull
	}

	peer := NewPeer(addr)
	pl.peers[addr] = peer
	return nil
}

func (pl *peerlist) addPeer(addr string) error {
	var err error
	pl.strand(func() {
		err = pl.add(addr)
	}, "AddPeer")

	return err
}

func (pl *peerlist) addPeers(addrs []string, verifyFunc func(string) error) int {
	n := len(addrs)
	pl.strand(func() {
		for _, addr := range addrs {
			if err := verifyFunc(addr); err != nil {
				logger.Warning("Failed to add peer %s, Reason: %v", addr, err)
				n--
				continue
			}

			if err := pl.add(addr); err != nil {
				logger.Warning("Failed to add peer %s, Reason: %v", addr, err)
				n--
			}
		}
	}, "AddPeers")
	return n
}

// GetPublicTrustPeers returns all trusted public peers
func (pl *peerlist) GetPublicTrustPeers() []*Peer {
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
// func (pl *peerlist) GetPrivateTrustPeers() []*Peer {
// 	var peers []*Peer
// 	pl.strand(func() {
// 		keys := pl.getTrustAddresses(true)
// 		peers = make([]*Peer, len(keys))
// 		for i, key := range keys {
// 			peers[i] = pl.peers[key]
// 		}
// 	}, "GetPrivateTrustPeers")
// 	return peers
// }

// GetAllTrustedPeers returns all trusted peers, including private and public peers.
func (pl *peerlist) GetAllTrustedPeers() []*Peer {
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

func (pl *peerlist) getTrustAddresses(private bool) []string {
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

// SetPrivate sets specific peer as private
func (pl *peerlist) setPrivate(addr string, private bool) error {
	var err error
	pl.strand(func() {
		if p, ok := pl.peers[addr]; ok {
			p.Private = private
			return
		}

		err = fmt.Errorf("Set peer.Private failed: %v does not exist in peer list", addr)
	}, "SetPrivate")
	return err
}

// SetTrusted sets peer as trusted peer
func (pl *peerlist) setTrusted(addr string, trusted bool) error {
	var err error
	pl.strand(func() {
		if p, ok := pl.peers[addr]; ok {
			p.Trusted = trusted
			return
		}

		err = fmt.Errorf("Set peer.Trusted failed: %v does not exist in peer list", addr)
	}, "SetTrusted")
	return err
}

// setPeerIsValid updates whether the peer is valid and has public incoming port
func (pl *peerlist) setPeerIsValid(addr string, valid bool) error {
	var err error
	pl.strand(func() {
		if p, ok := pl.peers[addr]; ok {
			p.Valid = valid
			p.Seen()
			return
		}

		err = fmt.Errorf("Set peer.HasIncomePort failed:%v does not exist in peer list", addr)
	}, "setPeerIsValid")
	return err
}

func (pl *peerlist) getAllTrustPeers() []string {
	return append(pl.getTrustAddresses(false), pl.getTrustAddresses(true)...)
}

// GetPublicAddresses returns the string addresses of all public peers
// func (pl *peerlist) GetPublicAddresses() []string {
// 	var addrs []string
// 	pl.strand(func() {
// 		addrs = pl.getAddresses(false)
// 	}, "GetPublicAddresses")
// 	return addrs
// }

// GetPrivateAddresses returns the string addresses of all private peers
func (pl *peerlist) GetPrivateAddresses() []string {
	var addrs []string
	pl.strand(func() {
		addrs = pl.getAddresses(true)
	}, "GetPrivateAddresses")
	return addrs
}

// RemovePeer removes peer
func (pl *peerlist) RemovePeer(addr string) {
	pl.strand(func() {
		delete(pl.peers, addr)
	}, "RemovePeer")
}

// cullInvalidPeers removes those unreachable and untrusted peers
func (pl *peerlist) cullInvalidPeers() []*Peer {
	var blacklistedPeers []*Peer
	pl.strand(func() {
		for _, p := range pl.peers {
			if p.Trusted {
				continue
			}

			if p.RetryTimes >= maxRetryTimes {
				blacklistedPeers = append(blacklistedPeers, p)
				logger.Critical("delete peer:%v", p.Addr)
				delete(pl.peers, p.Addr)
			}
		}
	}, "cullInvalidPeers")

	return blacklistedPeers
}

// Len returns number of peers
func (pl *peerlist) Len() int {
	return len(pl.peers)
}

// GetPeerByAddr returns peer of given address
func (pl *peerlist) GetPeerByAddr(addr string) (Peer, bool) {
	var peer Peer
	var exist bool
	pl.strand(func() {
		if p, ok := pl.peers[addr]; ok {
			peer = *p
			exist = true
			return
		}
	}, "GetPeerByAddr")
	return peer, exist
}

// ClearOld removes public peers that haven't been seen in timeAgo seconds
func (pl *peerlist) clearOld(timeAgo time.Duration) {
	t := utc.Now()
	pl.strand(func() {
		for addr, peer := range pl.peers {
			if !peer.Private && t.Sub(peer.LastSeen) > timeAgo {
				delete(pl.peers, addr)
			}
		}
	}, "ClearOld")
}

// Returns the string addresses of all public peers
func (pl *peerlist) getAddresses(private bool) []string {
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
func (pl *peerlist) random(count int, includePrivate bool) []*Peer {
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

func (pl *peerlist) getExchgAddr(private bool) []string {
	keys := []string{}
	for a, p := range pl.peers {
		if p.Valid && p.Private == private {
			keys = append(keys, a)
		}
	}
	return keys
}

// returns all exchangeable addresses
func (pl *peerlist) getAllExchgAddr() []string {
	return append(pl.getExchgAddr(true), pl.getExchgAddr(false)...)
}

// returns n random exchangeable peers, return all if count is 0.
func (pl *peerlist) randomExchg(count int, includePrivate bool) []*Peer {
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
func (pl *peerlist) RandomExchgPublic(count int) []*Peer {
	var peers []*Peer
	pl.strand(func() {
		peers = pl.randomExchg(count, false)
	}, "RandomExchgPublic")
	return peers
}

// RandomExchgAll returns n random exchangeable peers, including private peers.
// return all exchangeable peers if count is 0.
func (pl *peerlist) RandomExchgAll(count int) []*Peer {
	var peers []*Peer
	pl.strand(func() {
		peers = pl.randomExchg(count, true)
	}, "RandomExchgAll")
	return peers
}

// RandomPublic returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.  Will not include
// private peers.
func (pl *peerlist) RandomPublic(count int) []*Peer {
	var peers []*Peer
	pl.strand(func() {
		peers = pl.random(count, false)
	}, "RandomPublic")
	return peers
}

// RandomAll returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.  Includes private
// peers.
// func (pl *peerlist) RandomAll(count int) []*Peer {
// 	var peers []*Peer
// 	pl.strand(func() {
// 		peers = pl.random(count, true)
// 	}, "RandomAll")
// 	return peers
// }

// save saves known peers to disk as a newline delimited list of addresses to
// <dir><PeerDatabaseFilename>
func (pl *peerlist) save(dir string) (err error) {
	filename := PeerDatabaseFilename
	fn := filepath.Join(dir, filename)
	pl.strand(func() {
		// filter the peers that has retrytime > maxRetryTimes
		peers := make(map[string]*Peer)
		for k, p := range pl.peers {
			if p.RetryTimes <= maxRetryTimes {
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
func (pl *peerlist) IncreaseRetryTimes(addr string) {
	pl.strand(func() {
		if _, ok := pl.peers[addr]; ok {
			pl.peers[addr].IncreaseRetryTimes()
			pl.peers[addr].Seen()
		}
	}, "IncreaseRetryTimes")
}

// ResetRetryTimes reset retry times
func (pl *peerlist) ResetRetryTimes(addr string) {
	pl.strand(func() {
		if _, ok := pl.peers[addr]; ok {
			pl.peers[addr].ResetRetryTimes()
			pl.peers[addr].Seen()
		}
	}, "ResetRetryTimes")
}

// ResetAllRetryTimes reset all peers' retry times
func (pl *peerlist) ResetAllRetryTimes() {
	logger.Info("Reset all peer's retry times")
	pl.strand(func() {
		for _, p := range pl.peers {
			p.ResetRetryTimes()
		}
	}, "ResetAllRetryTimes")
}

// PrintAll print all peers
func (pl *peerlist) PrintAll() {
	pl.strand(func() {
		for _, p := range pl.peers {
			fmt.Println(p.String(), " ", p.RetryTimes, " is valid:", p.Valid)
		}
	})
}
