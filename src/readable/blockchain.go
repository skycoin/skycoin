package readable

import "github.com/skycoin/skycoin/src/coin"

// BlockchainMetadata encapsulates useful information from the coin.Blockchain
type BlockchainMetadata struct {
	// Most recent block's header
	Head BlockHeader `json:"head"`
	// Number of unspent outputs in the coin.Blockchain
	Unspents uint64 `json:"unspents"`
	// Number of known unconfirmed txns
	Unconfirmed uint64 `json:"unconfirmed"`
}

// NewBlockchainMetadata creates blockchain meta data
func NewBlockchainMetadata(head *coin.SignedBlock, unconfirmedLen, unspentsLen uint64) (*BlockchainMetadata, error) {
	return &BlockchainMetadata{
		Head:        NewBlockHeader(&head.Head),
		Unspents:    unspentsLen,
		Unconfirmed: unconfirmedLen,
	}, nil
}
