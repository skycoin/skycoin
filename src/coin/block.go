/*
Package coin defines the core blockchain datastructures.

This package should not have any dependencies except for go stdlib and cipher.
*/
package coin

import (
	"fmt"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
)

//go:generate skyencoder -struct BlockHeader -unexported
//go:generate skyencoder -struct BlockBody -unexported

// MaxBlockTransactions is the maximum number of transactions in a block (see the maxlen struct tag value applied to BlockBody.Transactions)
const MaxBlockTransactions = 65535

// Block represents the block struct
type Block struct {
	Head BlockHeader
	Body BlockBody
}

// HashPair including current block hash and previous block hash.
type HashPair struct {
	Hash     cipher.SHA256
	PrevHash cipher.SHA256
}

// BlockHeader records the block header
type BlockHeader struct {
	Version uint32

	Time  uint64
	BkSeq uint64 // Increment every block
	Fee   uint64 // Fee in block

	PrevHash cipher.SHA256 // Hash of header of previous block
	BodyHash cipher.SHA256 // Hash of transaction block

	UxHash cipher.SHA256 // XOR of sha256 of elements in unspent output set
}

// BlockBody represents the block body
type BlockBody struct {
	Transactions Transactions `enc:",maxlen=65535"`
}

// SignedBlock signed block
type SignedBlock struct {
	Block
	Sig cipher.Sig
}

// VerifySignature verifies that the block is signed by pubkey
func (b SignedBlock) VerifySignature(pubkey cipher.PubKey) error {
	return cipher.VerifyPubKeySignedHash(pubkey, b.Sig, b.HashHeader())
}

// NewBlock creates new block.
func NewBlock(prev Block, currentTime uint64, uxHash cipher.SHA256, txns Transactions, calc FeeCalculator) (*Block, error) {
	if len(txns) == 0 {
		return nil, fmt.Errorf("Refusing to create block with no transactions")
	}

	fee, err := txns.Fees(calc)
	if err != nil {
		// This should have been caught earlier
		return nil, fmt.Errorf("Invalid transaction fees: %v", err)
	}

	body := BlockBody{txns}
	head := NewBlockHeader(prev.Head, uxHash, currentTime, fee, body)
	return &Block{
		Head: head,
		Body: body,
	}, nil
}

// NewGenesisBlock creates genesis block
func NewGenesisBlock(genesisAddr cipher.Address, genesisCoins, timestamp uint64) (*Block, error) {
	txn := Transaction{}
	if err := txn.PushOutput(genesisAddr, genesisCoins, genesisCoins); err != nil {
		return nil, err
	}
	body := BlockBody{Transactions: Transactions{txn}}
	prevHash := cipher.SHA256{}
	bodyHash := body.Hash()
	head := BlockHeader{
		Time:     timestamp,
		BodyHash: bodyHash,
		PrevHash: prevHash,
		BkSeq:    0,
		Version:  0,
		Fee:      0,
		UxHash:   cipher.SHA256{},
	}
	b := &Block{
		Head: head,
		Body: body,
	}

	return b, nil
}

// HashHeader return hash of block head.
func (b Block) HashHeader() cipher.SHA256 {
	return b.Head.Hash()
}

// Time return the head time of the block.
func (b Block) Time() uint64 {
	return b.Head.Time
}

// Seq return the head seq of the block.
func (b Block) Seq() uint64 {
	return b.Head.BkSeq
}

// Size returns the size of the Block's Transactions, in bytes
func (b Block) Size() (uint32, error) {
	return b.Body.Size()
}

// NewBlockHeader creates block header
func NewBlockHeader(prev BlockHeader, uxHash cipher.SHA256, currentTime, fee uint64, body BlockBody) BlockHeader {
	if currentTime <= prev.Time {
		log.Panic("Time can only move forward")
	}
	bodyHash := body.Hash()
	prevHash := prev.Hash()
	return BlockHeader{
		BodyHash: bodyHash,
		Version:  prev.Version,
		PrevHash: prevHash,
		Time:     currentTime,
		BkSeq:    prev.BkSeq + 1,
		Fee:      fee,
		UxHash:   uxHash,
	}
}

// Hash return hash of block header
func (bh *BlockHeader) Hash() cipher.SHA256 {
	return cipher.SumSHA256(bh.Bytes())
}

// Bytes serialize the blockheader and return the byte value.
func (bh *BlockHeader) Bytes() []byte {
	buf, err := encodeBlockHeader(bh)
	if err != nil {
		log.Panicf("encodeBlockHeader failed: %v", err)
	}
	return buf
}

// Hash returns the merkle hash of contained transactions
func (bb BlockBody) Hash() cipher.SHA256 {
	hashes := make([]cipher.SHA256, len(bb.Transactions))
	for i := range bb.Transactions {
		hashes[i] = bb.Transactions[i].Hash()
	}
	// Merkle hash of transactions
	return cipher.Merkle(hashes)
}

// Size returns the size of Transactions, in bytes
func (bb BlockBody) Size() (uint32, error) {
	// We can't use length of self.Bytes() because it has a length prefix
	// Need only the sum of transaction sizes
	return bb.Transactions.Size()
}

// Bytes serialize block body, and return the byte value.
func (bb *BlockBody) Bytes() []byte {
	buf, err := encodeBlockBody(bb)
	if err != nil {
		log.Panicf("encodeBlockBody failed: %v", err)
	}
	return buf
}

// CreateUnspents creates the expected outputs for a transaction.
func CreateUnspents(bh BlockHeader, txn Transaction) UxArray {
	var h cipher.SHA256
	// The genesis block uses the null hash as the SrcTransaction [FIXME hardfork]
	if bh.BkSeq != 0 {
		h = txn.Hash()
	}
	uxo := make(UxArray, len(txn.Out))
	for i := range txn.Out {
		uxo[i] = UxOut{
			Head: UxHead{
				Time:  bh.Time,
				BkSeq: bh.BkSeq,
			},
			Body: UxBody{
				SrcTransaction: h,
				Address:        txn.Out[i].Address,
				Coins:          txn.Out[i].Coins,
				Hours:          txn.Out[i].Hours,
			},
		}
	}
	return uxo
}

// CreateUnspent creates single unspent output
func CreateUnspent(bh BlockHeader, txn Transaction, outIndex int) (UxOut, error) {
	if outIndex < 0 || outIndex >= len(txn.Out) {
		return UxOut{}, fmt.Errorf("Transaction out index overflows transaction outputs")
	}

	var h cipher.SHA256
	// The genesis block uses the null hash as the SrcTransaction [FIXME hardfork]
	if bh.BkSeq != 0 {
		h = txn.Hash()
	}

	return UxOut{
		Head: UxHead{
			Time:  bh.Time,
			BkSeq: bh.BkSeq,
		},
		Body: UxBody{
			SrcTransaction: h,
			Address:        txn.Out[outIndex].Address,
			Coins:          txn.Out[outIndex].Coins,
			Hours:          txn.Out[outIndex].Hours,
		},
	}, nil
}
