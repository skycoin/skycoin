package blockdb

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// BlockPQSigsBkt holds the ML-DSA-44 signatures of blocks, keyed by header hash.
	// Only blocks from the post-quantum fork height have one.
	BlockPQSigsBkt = []byte("block_pq_sigs")
)

// blockPQSigs stores post-quantum block signatures apart from the secp256k1 ones,
// so the block_sigs encoding and the nodes that read it are unchanged.
type blockPQSigs struct{}

// Get returns the post-quantum signature of a block, if it has one. A database
// from before the post-quantum code, opened read-only, has no bucket yet.
func (bs *blockPQSigs) Get(tx *dbutil.Tx, hash cipher.SHA256) ([]byte, error) {
	if !dbutil.Exists(tx, BlockPQSigsBkt) {
		return nil, nil
	}
	return dbutil.GetBucketValue(tx, BlockPQSigsBkt, hash[:])
}

// Add stores the post-quantum signature of a block.
func (bs *blockPQSigs) Add(tx *dbutil.Tx, hash cipher.SHA256, sig []byte) error {
	return dbutil.PutBucketValue(tx, BlockPQSigsBkt, hash[:], sig)
}

var (
	// TxnPQWitnessesBkt holds the PQ witnesses of transactions in blocks, keyed by
	// transaction hash, because the block encoding leaves them out.
	TxnPQWitnessesBkt = []byte("txn_pq_witnesses")
)

// txnPQWitnesses stores the post-quantum witnesses of transactions in blocks.
type txnPQWitnesses struct{}

func (tw *txnPQWitnesses) add(tx *dbutil.Tx, b *coin.Block) error {
	for _, txn := range b.Body.Transactions {
		if txn.Type != coin.TxnTypePQ {
			continue
		}
		h := txn.Hash()
		if err := dbutil.PutBucketValue(tx, TxnPQWitnessesBkt, h[:], coin.EncodePQWitness(txn.PQ)); err != nil {
			return err
		}
	}
	return nil
}

// load puts the stored witnesses back onto the transactions of b.
func (tw *txnPQWitnesses) load(tx *dbutil.Tx, b *coin.Block) error {
	for i := range b.Body.Transactions {
		txn := &b.Body.Transactions[i]
		if txn.Type != coin.TxnTypePQ {
			continue
		}
		h := txn.Hash()
		v, err := dbutil.GetBucketValue(tx, TxnPQWitnessesBkt, h[:])
		if err != nil {
			return err
		}
		if v == nil {
			return fmt.Errorf("post-quantum witness of transaction %s not found", h.Hex())
		}
		if txn.PQ, _, err = coin.DecodePQWitness(v); err != nil {
			return err
		}
	}
	return nil
}

// GetBlockPQ returns the post-quantum signature of b, nil if it has none, and puts
// the stored PQ witnesses back onto its transactions.
func (bc *Blockchain) GetBlockPQ(tx *dbutil.Tx, b *coin.Block) ([]byte, error) {
	sig, err := bc.pqSigs.Get(tx, b.HashHeader())
	if err != nil {
		return nil, err
	}
	if err := bc.pqWits.load(tx, b); err != nil {
		return nil, err
	}
	return sig, nil
}
