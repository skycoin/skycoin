/*
Package readable defines JSON-tagged struct representations of internal binary data structures,
for use by the API and CLI.
*/
package readable

import (
	"errors"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
)

// BlockHeader represents the readable block header
type BlockHeader struct {
	BkSeq        uint64 `json:"seq"`
	Hash         string `json:"block_hash"`
	PreviousHash string `json:"previous_block_hash"`
	Time         uint64 `json:"timestamp"`
	Fee          uint64 `json:"fee"`
	Version      uint32 `json:"version"`
	BodyHash     string `json:"tx_body_hash"`
	UxHash       string `json:"ux_hash"`
}

// NewBlockHeader creates a readable block header
func NewBlockHeader(b coin.BlockHeader) BlockHeader {
	return BlockHeader{
		BkSeq:        b.BkSeq,
		Hash:         b.Hash().Hex(),
		PreviousHash: b.PrevHash.Hex(),
		Time:         b.Time,
		Fee:          b.Fee,
		Version:      b.Version,
		BodyHash:     b.BodyHash.Hex(),
		UxHash:       b.UxHash.Hex(),
	}
}

// ToCoinBlockHeader converts BlockHeader back to coin.BlockHeader
func (bh BlockHeader) ToCoinBlockHeader() (coin.BlockHeader, error) {
	prevHash, err := cipher.SHA256FromHex(bh.PreviousHash)
	if err != nil {
		return coin.BlockHeader{}, err
	}

	bodyHash, err := cipher.SHA256FromHex(bh.BodyHash)
	if err != nil {
		return coin.BlockHeader{}, err
	}

	uxHash, err := cipher.SHA256FromHex(bh.UxHash)
	if err != nil {
		return coin.BlockHeader{}, err
	}

	headHash, err := cipher.SHA256FromHex(bh.Hash)
	if err != nil {
		return coin.BlockHeader{}, err
	}

	cbh := coin.BlockHeader{
		Version:  bh.Version,
		Time:     bh.Time,
		BkSeq:    bh.BkSeq,
		Fee:      bh.Fee,
		PrevHash: prevHash,
		BodyHash: bodyHash,
		UxHash:   uxHash,
	}

	if cbh.Hash() != headHash {
		return coin.BlockHeader{}, errors.New("readable.BlockHeader.Hash != recovered coin.BlockHeader.Hash()")
	}

	return cbh, nil
}

// BlockBody represents a readable block body
type BlockBody struct {
	Transactions []Transaction `json:"txns"`
}

// NewBlockBody creates a readable block body
func NewBlockBody(b coin.Block) (*BlockBody, error) {
	txns := make([]Transaction, len(b.Body.Transactions))
	isGenesis := b.Head.BkSeq == 0
	for i := range b.Body.Transactions {
		txn, err := NewTransaction(b.Body.Transactions[i], isGenesis)
		if err != nil {
			return nil, err
		}
		txns[i] = *txn
	}

	return &BlockBody{
		Transactions: txns,
	}, nil
}

// Block represents a readable block
type Block struct {
	Head BlockHeader `json:"header"`
	Body BlockBody   `json:"body"`
	Size uint32      `json:"size"`
}

// NewBlock creates a readable block
func NewBlock(b coin.Block) (*Block, error) {
	body, err := NewBlockBody(b)
	if err != nil {
		return nil, err
	}

	size, err := b.Size()
	if err != nil {
		return nil, err
	}

	return &Block{
		Head: NewBlockHeader(b.Head),
		Body: *body,
		Size: size,
	}, nil
}

// Blocks an array of readable blocks.
type Blocks struct {
	Blocks []Block `json:"blocks"`
}

// NewBlocks converts []coin.SignedBlock to Blocks
func NewBlocks(blocks []coin.SignedBlock) (*Blocks, error) {
	rbs := make([]Block, 0, len(blocks))
	for _, b := range blocks {
		rb, err := NewBlock(b.Block)
		if err != nil {
			return nil, err
		}
		rbs = append(rbs, *rb)
	}
	return &Blocks{
		Blocks: rbs,
	}, nil
}
