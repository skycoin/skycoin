package readable

import "github.com/skycoin/skycoin/src/coin"

// BlockHeader represents the readable block header
type BlockHeader struct {
	BkSeq             uint64 `json:"seq"`
	BlockHash         string `json:"block_hash"`
	PreviousBlockHash string `json:"previous_block_hash"`
	Time              uint64 `json:"timestamp"`
	Fee               uint64 `json:"fee"`
	Version           uint32 `json:"version"`
	BodyHash          string `json:"tx_body_hash"`
}

// NewBlockHeader creates a readable block header
func NewBlockHeader(b *coin.BlockHeader) BlockHeader {
	return BlockHeader{
		BkSeq:             b.BkSeq,
		BlockHash:         b.Hash().Hex(),
		PreviousBlockHash: b.PrevHash.Hex(),
		Time:              b.Time,
		Fee:               b.Fee,
		Version:           b.Version,
		BodyHash:          b.BodyHash.Hex(),
	}
}

// BlockBody represents a readable block body
type BlockBody struct {
	Transactions []Transaction `json:"txns"`
}

// NewBlockBody creates a readable block body
func NewBlockBody(b *coin.Block) (*BlockBody, error) {
	txns := make([]Transaction, len(b.Body.Transactions))
	for i := range b.Body.Transactions {
		txn, err := NewTransaction(b.Body.Transactions[i], true)
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
	Size int         `json:"size"`
}

// NewBlock creates a readable block
func NewBlock(b *coin.Block) (*Block, error) {
	body, err := NewBlockBody(b)
	if err != nil {
		return nil, err
	}
	return &Block{
		Head: NewBlockHeader(&b.Head),
		Body: *body,
		Size: b.Size(),
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
		rb, err := NewBlock(&b.Block)
		if err != nil {
			return nil, err
		}
		rbs = append(rbs, *rb)
	}
	return &Blocks{
		Blocks: rbs,
	}, nil
}
