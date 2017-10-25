package pex

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/utc"
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
	lock  sync.Mutex
	peers map[string]*Peer
	cap   int // the max number of peers that the peerlist can take
}

func newPeerlist(maxPeers int) *peerlist {
	return &peerlist{
		peers: make(map[string]*Peer, maxPeers),
		cap:   maxPeers,
	}
}

// Filter peers filter
type Filter func(peers Peers) Peers

// loadIfExist loads if the peer.txt file does exist
func (pl *peerlist) loadIfExist(dir string) error {
	fn := filepath.Join(dir, PeerDatabaseFilename)
	// check if the file does exist
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return nil
	}
	return file.LoadJSON(fn, &pl.peers)
}

func (pl *peerlist) strand(f func(), arg ...interface{}) {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	// logger.Critical("%v", arg)
	f()
}

// IsFull returns if the peerlist is full
func (pl *peerlist) IsFull() bool {
	var full bool
	pl.strand(func() {
		full = pl.isFull()
	}, "Full")
	return full
}

func (pl *peerlist) isFull() bool {
	return pl.cap > 0 && len(pl.peers) >= pl.cap
}

func (pl *peerlist) add(addr string) error {
	if p, ok := pl.peers[addr]; ok && p != nil {
		p.Seen()
		return nil
	}

	if pl.isFull() {
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
	}, "addPeer")

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
	}, "addPeers")
	return n
}

// GetPeers returns peers that can pass through the filters if any,
// otherwise returns all peers that are allowed to connect to.
func (pl *peerlist) GetPeers(flts ...Filter) Peers {
	var ps Peers
	pl.strand(func() {
		ps = pl.getPeers(flts...)
	}, "GetPeers")

	return ps
}

func (pl *peerlist) getPeers(flts ...Filter) Peers {
	var ps Peers
	for _, p := range pl.peers {
		if p.CanTry() {
			ps = append(ps, *p)
		}
	}

	for _, flt := range flts {
		ps = flt(ps)
	}
	return ps
}

// filters

// isPrivate filters private peers
func isPrivate(peers Peers) Peers {
	var ps Peers
	for _, p := range peers {
		if p.Private {
			ps = append(ps, p)
		}
	}

	return ps
}

// isPublic filters public peers
func isPublic(peers Peers) Peers {
	var ps Peers
	for _, p := range peers {
		if !p.Private {
			ps = append(ps, p)
		}
	}
	return ps
}

// isTrusted filters trusted peers
func isTrusted(peers Peers) Peers {
	var ps Peers
	for _, p := range peers {
		if p.Trusted {
			ps = append(ps, p)
		}
	}
	return ps
}

// hasPublicPort filters peers that have public port
func hasPublicPort(peers Peers) Peers {
	var ps Peers
	for _, p := range peers {
		if p.HasPublicPort {
			ps = append(ps, p)
		}
	}
	return ps
}

// isExchangeable filters exchangeable peers
var isExchangeable = hasPublicPort

// RemovePeer removes peer
func (pl *peerlist) RemovePeer(addr string) {
	pl.strand(func() {
		delete(pl.peers, addr)
	}, "RemovePeer")
}

// GetPublicTrustPeers returns all trusted public peers
// func (pl *peerlist) GetPublicTrustPeers() Peers {
// 	return pl.getPeersSafe(IsPublic, IsTrusted)
// }

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
// func (pl *peerlist) GetAllTrustedPeers() Peers {
// var peers []*Peer
// pl.strand(func() {
// 	keys := pl.getAllTrustPeers()
// 	peers = make([]*Peer, len(keys))
// 	for i, key := range keys {
// 		peers[i] = pl.peers[key]
// 	}
// }, "GetAllTrustedPeers")
// return peers
// return pl.getPeersSafe(trustFilter)
// }

// func (pl *peerlist) getTrustAddresses(private bool) []string {
// 	keys := []string{}
// 	for key, p := range pl.peers {
// 		if p.Trusted {
// 			if p.CanTry() {
// 				if private && p.Private {
// 					keys = append(keys, key)
// 				} else if !private && !p.Private {
// 					keys = append(keys, key)
// 				}
// 			}
// 		}
// 	}
// 	return keys
// }

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

// setPeerHasPublicPort updates whether the peer is valid and has public incoming port
func (pl *peerlist) setPeerHasPublicPort(addr string, hasPublicPort bool) error {
	var err error
	pl.strand(func() {
		if p, ok := pl.peers[addr]; ok {
			p.HasPublicPort = hasPublicPort
			p.Seen()
			return
		}

		err = fmt.Errorf("Set peer.HasIncomePort failed:%v does not exist in peer list", addr)
	}, "setPeerIsValid")
	return err
}

// func (pl *peerlist) getAllTrustPeers() []string {
// 	return append(pl.getTrustAddresses(false), pl.getTrustAddresses(true)...)
// }

// GetPublicAddresses returns the string addresses of all public peers
// func (pl *peerlist) GetPublicAddresses() []string {
// 	var addrs []string
// 	pl.strand(func() {
// 		addrs = pl.getAddresses(false)
// 	}, "GetPublicAddresses")
// 	return addrs
// }

// GetPrivateAddresses returns the string addresses of all private peers
// func (pl *peerlist) GetPrivatePeers() Peers {
// var addrs []string
// pl.strand(func() {
// 	addrs = pl.getAddresses(true)
// }, "GetPrivateAddresses")
// return addrs
// 	return pl.getPeersSafe(privateFilter)
// }

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
	var l int
	pl.strand(func() {
		l = len(pl.peers)
	}, "Len")
	return l
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
	}, "clearOld")
}

// Returns n random peers, or all of the peers, whichever is lower.
// If count is 0, all of the peers are returned, shuffled.
func (pl *peerlist) random(count int, flts ...Filter) Peers {
	keys := pl.getPeers(flts...).ToAddrs()
	if len(keys) == 0 {
		return Peers{}
	}
	max := count
	if count == 0 || count > len(keys) {
		max = len(keys)
	}
	var ps Peers
	perm := rand.Perm(len(keys))
	for _, i := range perm[:max] {
		ps = append(ps, *pl.peers[keys[i]])
	}
	return ps
}

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
			fmt.Println(p.String(), " ", p.RetryTimes, " public:", p.HasPublicPort)
		}
	})
}

// GetTrustPeers returns trusted peers
func (pl *peerlist) Trust() Peers {
	var ps Peers
	pl.strand(func() {
		ps = pl.getPeers(isTrusted)
	}, "Trust")
	return ps
}

// GetPrivate returns private peers
func (pl *peerlist) Private() Peers {
	var ps Peers
	pl.strand(func() {
		ps = pl.getPeers(isPrivate)
	}, "Private")
	return ps
}

// GetTrustPublicPeers returns trusted public peers
func (pl *peerlist) TrustPublic() Peers {
	var ps Peers
	pl.strand(func() {
		ps = pl.getPeers(isPublic, isTrusted)
	}, "TrustPublic")
	return ps
}

// GetRandomPublicPeers returns N random public peers
func (pl *peerlist) RandomPublic(n int) Peers {
	var ps Peers
	pl.strand(func() {
		ps = pl.random(n, isPublic)
	}, "GetRandomPublic")
	return ps
}

// GetRandomValidPublic returns N random valid peers
func (pl *peerlist) RandomValidPublic(n int) Peers {
	var ps Peers
	pl.strand(func() {
		ps = pl.random(n, isPublic)
	}, "RandomValidPublic")
	return ps
}

// RandomValid returns N random valid peers, returns all if N is 0
func (pl *peerlist) RandomValid(n int) Peers {
	var ps Peers
	pl.strand(func() {
		ps = pl.random(n)
	}, "RandomValid")
	return ps
}
