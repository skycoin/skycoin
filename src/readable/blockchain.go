package readable

import (
	"github.com/SkycoinProject/skycoin/src/daemon"
	"github.com/SkycoinProject/skycoin/src/visor"
)

// BlockchainMetadata encapsulates useful information from the coin.Blockchain
type BlockchainMetadata struct {
	// Most recent block's header
	Head BlockHeader `json:"head"`
	// Number of unspent outputs in the coin.Blockchain
	Unspents uint64 `json:"unspents"`
	// Number of known unconfirmed txns
	Unconfirmed uint64 `json:"unconfirmed"`
}

// NewBlockchainMetadata creates blockchain metadata
func NewBlockchainMetadata(bm visor.BlockchainMetadata) BlockchainMetadata {
	return BlockchainMetadata{
		Head:        NewBlockHeader(bm.HeadBlock.Head),
		Unspents:    bm.Unspents,
		Unconfirmed: bm.Unconfirmed,
	}
}

// BlockchainProgress is the current blockchain syncing status
type BlockchainProgress struct {
	// Our current blockchain length
	Current uint64 `json:"current"`
	// Our best guess at true blockchain length
	Highest uint64 `json:"highest"`
	// Individual blockchain length reports from peers
	Peers []PeerBlockchainHeight `json:"peers"`
}

// PeerBlockchainHeight is a peer's IP address with their reported blockchain height
type PeerBlockchainHeight struct {
	Address string `json:"address"`
	Height  uint64 `json:"height"`
}

// NewBlockchainProgress copies daemon.BlockchainProgress to a struct with json tags
func NewBlockchainProgress(bp *daemon.BlockchainProgress) BlockchainProgress {
	peers := make([]PeerBlockchainHeight, len(bp.Peers))
	for i, p := range bp.Peers {
		peers[i] = PeerBlockchainHeight{
			Address: p.Address,
			Height:  p.Height,
		}
	}

	return BlockchainProgress{
		Current: bp.Current,
		Highest: bp.Highest,
		Peers:   peers,
	}
}
