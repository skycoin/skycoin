package daemon

import (
	"sort"
	"strings"
	"sync"
)

// PeerBlockchainHeight is a peer's IP address with their reported blockchain height
type PeerBlockchainHeight struct {
	Address string `json:"address"`
	Height  uint64 `json:"height"`
}

// peerBlockchainHeights tracks reported blockchain heights of peers
type peerBlockchainHeights struct {
	// Peer-reported blockchain height.  Use to estimate download progress
	heights map[string]uint64
	sync.Mutex
}

// newPeerBlockchainHeights creates a peerBlockchainHeights
func newPeerBlockchainHeights() *peerBlockchainHeights {
	return &peerBlockchainHeights{
		heights: make(map[string]uint64),
	}
}

// Remove removes a connection from the records
func (p *peerBlockchainHeights) Remove(addr string) {
	p.Lock()
	defer p.Unlock()

	delete(p.heights, addr)
}

// Record saves a peer-reported blockchain height
func (p *peerBlockchainHeights) Record(addr string, height uint64) {
	p.Lock()
	defer p.Unlock()

	p.heights[addr] = height
}

// Estimate returns the blockchain length estimated from peer reports.
// The highest height reported amongst all peers, and including the node itself,
// is returned.
func (p *peerBlockchainHeights) Estimate(headSeq uint64) uint64 {
	p.Lock()
	defer p.Unlock()

	for _, seq := range p.heights {
		if headSeq < seq {
			headSeq = seq
		}
	}

	return headSeq
}

// All returns recorded peers' blockchain heights as an array.
// The array is sorted by address as strings.
func (p *peerBlockchainHeights) All() []PeerBlockchainHeight {
	p.Lock()
	defer p.Unlock()

	if len(p.heights) == 0 {
		return nil
	}

	peerHeights := make([]PeerBlockchainHeight, 0, len(p.heights))
	for addr, height := range p.heights {
		peerHeights = append(peerHeights, PeerBlockchainHeight{
			Address: addr,
			Height:  height,
		})
	}

	sort.Slice(peerHeights, func(i, j int) bool {
		return strings.Compare(peerHeights[i].Address, peerHeights[j].Address) < 0
	})

	return peerHeights
}
