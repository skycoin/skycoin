package coin

import (
	"fmt"
	"log"

    logging "gopkg.in/op/go-logging.v1"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
)

var logger = logging.MustGetLogger("skycoin.coin")

type Block struct {
	Head BlockHeader
	Body BlockBody
}

// HashPair including current block hash and previous block hash.
type HashPair struct {
	Hash    cipher.SHA256
	PreHash cipher.SHA256
}

type BlockHeader struct {
	Version uint32

	Time  uint64
	BkSeq uint64 //increment every block
	Fee   uint64 //fee in block, used for Proof of Stake

	PrevHash cipher.SHA256 //hash of header of previous block
	BodyHash cipher.SHA256 //hash of transaction block

	UxHash cipher.SHA256 //XOR of sha256 of elements in unspent output set

}

type BlockBody struct {
	Transactions Transactions
}

type SignedBlock struct {
	Block Block
	Sig   cipher.Sig
}

//TODO: merge header/body and cleanup top level interface

/*
Todo: merge header/body

type Block struct {
    Time  uint64
    BkSeq uint64 //increment every block
    Fee   uint64 //fee in block, used for Proof of Stake

    HashPrevBlock cipher.SHA256 //hash of header of previous block
    BodyHash      cipher.SHA256 //hash of transaction block

    Transactions Transactions
}
*/

// NewBlock creates new block.
func NewBlock(prev Block, currentTime uint64, unspent UnspentPool,
	txns Transactions, calc FeeCalculator) Block {
	if len(txns) == 0 {
		log.Panic("Refusing to create block with no transactions")
	}
	fee, err := txns.Fees(calc)
	if err != nil {
		// This should have been caught earlier
		log.Panicf("Invalid transaction fees: %v", err)
	}
	body := BlockBody{txns}
	return Block{
		Head: NewBlockHeader(prev.Head, unspent, currentTime, fee, body),
		Body: body,
	}
}

// HashHeader return hash of block head.
func (b Block) HashHeader() cipher.SHA256 {
	return b.Head.Hash()
}

// PreHashHeader return hash of prevous block.
func (b Block) PreHashHeader() cipher.SHA256 {
	return b.Head.PrevHash
}

// Time return the head time of the block.
func (b Block) Time() uint64 {
	return b.Head.Time
}

// Seq return the head seq of the block.
func (b Block) Seq() uint64 {
	return b.Head.BkSeq
}

// HashBody return hash of block body.
func (b Block) HashBody() cipher.SHA256 {
	return b.Body.Hash()
}

// Size returns the size of the Block's Transactions, in bytes
func (b Block) Size() int {
	return b.Body.Size()
}

// String return readable string of block.
func (b Block) String() string {
	return b.Head.String()
}

// GetTransaction looks up a Transaction by its Head.Hash.
// Returns the Transaction and whether it was found or not
// TODO -- build a private index on the block, or a global blockchain one
// mapping txns to their block + tx index
// TODO: Deprecate? Utility Function
func (b Block) GetTransaction(txHash cipher.SHA256) (Transaction, bool) {
	txns := b.Body.Transactions
	for i := range txns {
		if txns[i].Hash() == txHash {
			return txns[i], true
		}
	}
	return Transaction{}, false
}

func NewBlockHeader(prev BlockHeader, unspent UnspentPool, currentTime,
	fee uint64, body BlockBody) BlockHeader {
	if currentTime <= prev.Time {
		log.Panic("Time can only move forward")
	}
	prevHash := prev.Hash()
	return BlockHeader{
		BodyHash: body.Hash(),
		Version:  prev.Version,
		PrevHash: prevHash,
		Time:     currentTime,
		BkSeq:    prev.BkSeq + 1,
		Fee:      fee,
		UxHash:   unspent.GetUxHash(),
	}
}

// Hash return hash of block header
func (bh BlockHeader) Hash() cipher.SHA256 {
	b1 := encoder.Serialize(bh)
	return cipher.SumSHA256(b1)
}

// Bytes serialize the blockheader and return the byte value.
func (bh BlockHeader) Bytes() []byte {
	return encoder.Serialize(bh)
}

// String return readable string of block header.
func (bh BlockHeader) String() string {
	return fmt.Sprintf("Version: %d\nTime: %d\nBkSeq: %d\nFee: %d\n"+
		"PrevHash: %s\nBodyHash: %s", bh.Version, bh.Time, bh.BkSeq,
		bh.Fee, bh.PrevHash.Hex(), bh.BodyHash.Hex())
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
func (bb BlockBody) Size() int {
	// We can't use length of self.Bytes() because it has a length prefix
	// Need only the sum of transaction sizes
	return bb.Transactions.Size()
}

// Bytes serialize block body, and return the byte value.
func (bb BlockBody) Bytes() []byte {
	return encoder.Serialize(bb)
}

// CreateUnspents creates the expected outputs for a transaction.
func CreateUnspents(bh BlockHeader, tx Transaction) UxArray {
	var h cipher.SHA256
	if bh.BkSeq != 0 {
		// not genesis block
		h = tx.Hash()
	}
	uxo := make(UxArray, len(tx.Out))
	for i := range tx.Out {
		uxo[i] = UxOut{
			Head: UxHead{
				Time:  bh.Time,
				BkSeq: bh.BkSeq,
			},
			Body: UxBody{
				SrcTransaction: h,
				Address:        tx.Out[i].Address,
				Coins:          tx.Out[i].Coins,
				Hours:          tx.Out[i].Hours,
			},
		}
	}
	return uxo
}
