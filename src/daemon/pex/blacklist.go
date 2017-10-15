package pex

import "sync"

type blacklist struct {
	sync.Mutex
	peers map[string]*Peer
}

func newBlacklist() *blacklist {
	return &blacklist{
		peers: make(map[string]*Peer),
	}
}

// Add adds peer to blacklist
func (bl *blacklist) Add(p *Peer) {
	bl.Lock()
	defer bl.Unlock()
	bl.peers[p.Addr] = p
}

// Delete deletes peer of given address from blacklist
func (bl *blacklist) Delete(addr string) {
	bl.Lock()
	bl.Unlock()
	delete(bl.peers, addr)
}
