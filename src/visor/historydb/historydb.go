// Package historydb is in charge of parsing the consuses blokchain, and providing
// apis for blockchain explorer.
package historydb

import (
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
)

// Blockchainer interface for isolating the detail of blockchain.
type Blockchainer interface {
	Head() *coin.Block
	GetUnspent() *coin.UnspentPool
	GetBlockInDepth(dep uint64) *coin.Block
	ExecuteBlock(b *coin.Block) (coin.UxArray, error)
	CreateGenesisBlock(genAddress cipher.Address, genCoins, timestamp uint64) coin.Block
	VerifyTransaction(tx coin.Transaction) error
	GetBlock(hash cipher.SHA256) *coin.Block
}

// NewDB create the history bolt db if does not exsit.
func NewDB() (*bolt.DB, error) {
	dbFile := filepath.Join(util.DataDir, "history.db")
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// HistoryDB provides apis for blockchain explorer.
type HistoryDB struct {
	db      *bolt.DB      // bolt db instance.
	blocks  *blocks       // blocks bucket.
	txns    *transactions // transactions bucket.
	outputs *UxOuts       // outputs bucket.
	addrIn  *addressUx    // bucket which stores all UxOuts that address recved.
	addrOut *addressUx    // bucket which stores all UxOuts that address spent.
}

// New create historydb instance and create corresponding buckets if does not exist.
func New(db *bolt.DB) (*HistoryDB, error) {
	hd := HistoryDB{db: db}
	var err error
	hd.blocks, err = newBlockBkt(db)
	if err != nil {
		return nil, err
	}

	hd.txns, err = newTransactionsBkt(db)
	if err != nil {
		return nil, err
	}

	// create the output instance
	hd.outputs, err = newOutputsBkt(db)
	if err != nil {
		return nil, err
	}

	// create the toAddressTx instance.
	hd.addrIn, err = newAddressInBkt(db)
	if err != nil {
		return nil, err
	}

	// create the fromAddressTx instance.
	hd.addrOut, err = newAddressOutBkt(db)
	if err != nil {
		return nil, err
	}

	return &hd, nil
}

// ProcessBlockchain process the blocks in the chain.
func (hd *HistoryDB) ProcessBlockchain(bc Blockchainer) error {
	depth := bc.Head().Seq()
	for i := uint64(0); i <= depth; i++ {
		b := bc.GetBlockInDepth(i)
		if err := hd.ProcessBlock(b); err != nil {
			return err
		}
	}
	return nil
}

// GetUxout get UxOut of specific uxID.
func (hd *HistoryDB) GetUxout(uxID cipher.SHA256) (*UxOut, error) {
	return hd.outputs.Get(uxID)
}

// ProcessBlock will index the transaction, outputs,etc.
func (hd *HistoryDB) ProcessBlock(b *coin.Block) error {
	// store the block
	if err := hd.blocks.Add(b); err != nil {
		return err
	}

	// index the transactions
	for _, t := range b.Body.Transactions {
		tx := Transaction{
			Tx:       t,
			BlockSeq: b.Seq(),
		}
		if err := hd.txns.Add(&tx); err != nil {
			return err
		}

		// handle tx in, genesis transaction's vin is empty, so should be ignored.
		if b.Seq() > 0 {
			for _, in := range t.In {
				o, err := hd.outputs.Get(in)
				if err != nil {
					return err
				}
				// update the spent block seq of the output.
				o.SpentBlockSeq = b.Seq()
				o.SpentTxID = t.Hash()
				if err := hd.outputs.Set(*o); err != nil {
					return err
				}

				// index the output for address out
				if err := hd.addrOut.Add(o.Out.Body.Address, o.Hash()); err != nil {
					return err
				}
			}
		}

		// handle the tx out
		uxArray := coin.CreateUnspents(b.Head, t)
		for _, ux := range uxArray {
			uxOut := UxOut{
				Out: ux,
			}
			if err := hd.outputs.Set(uxOut); err != nil {
				return err
			}

			if err := hd.addrIn.Add(ux.Body.Address, ux.Hash()); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetTransaction get transaction by hash.
func (hd HistoryDB) GetTransaction(hash cipher.SHA256) (*Transaction, error) {
	return hd.txns.Get(hash)
}

// GetTxsInBlock get all transactions in specifc block.
func (hd HistoryDB) GetTxsInBlock(blockHash cipher.SHA256) ([]*Transaction, error) {
	// get block
	b, err := hd.blocks.Get(blockHash)
	if err != nil {
		return nil, err
	}
	// get txs in the block
	txs := b.Body.Transactions
	atxs := make([]*Transaction, len(txs))
	for i, t := range txs {
		atxs[i] = &Transaction{
			Tx:       t,
			BlockSeq: b.Seq(),
		}
	}
	return atxs, nil
}

// GetSpentUxOutOfAddr get all spent uxout of specifc address.
func (hd HistoryDB) GetSpentUxOutOfAddr(address cipher.Address) ([]*UxOut, error) {
	hashes, err := hd.addrOut.Get(address)
	if err != nil {
		return []*UxOut{}, err
	}
	uxOuts := make([]*UxOut, len(hashes))
	for i, hash := range hashes {
		ux, err := hd.outputs.Get(hash)
		if err != nil {
			return []*UxOut{}, err
		}
		uxOuts[i] = ux
	}
	return uxOuts, nil
}

// GetRecvUxOutOfAddr get all uxout that the address received.
func (hd HistoryDB) GetRecvUxOutOfAddr(address cipher.Address) ([]*UxOut, error) {
	hashes, err := hd.addrIn.Get(address)
	if err != nil {
		return []*UxOut{}, err
	}
	uxOuts := make([]*UxOut, len(hashes))
	for i, hash := range hashes {
		ux, err := hd.outputs.Get(hash)
		if err != nil {
			return []*UxOut{}, err
		}
		uxOuts[i] = ux
	}
	return uxOuts, nil
}
