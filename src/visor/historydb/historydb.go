// Package historydb is in charge of parsing the consuses blokchain, and providing
// apis for blockchain explorer.
package historydb

import (
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
)

// HistoryDB provides apis for blockchain explorer.
type HistoryDB struct {
	db      *bolt.DB      // bolt db instance.
	txns    *Transactions // transactions bucket instance.
	outputs *Outputs      // outputs bucket instance.
	addrIn  *addressIn
	addrOut *addressOut
}

// Start will open a boltdb named history.db,
// and create corresponding buckets if does not exist.
func (hd *HistoryDB) Start() error {
	dbFile := filepath.Join(util.DataDir, "history.db")
	var err error
	hd.db, err = bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// create the transactions instance.
	hd.txns, err = newTransactions(hd.db)
	if err != nil {
		return err
	}

	// create the toAddressTx instance.
	hd.addrIn, err = newAddressIn(hd.db)
	if err != nil {
		return err
	}

	// create the fromAddressTx instance.
	hd.addrOut, err = newAddressOut(hd.db)
	if err != nil {
		return err
	}

	return nil
}

// Stop the historydb.
func (hd *HistoryDB) Stop() {
	hd.db.Close()
}

// ProcessBlockchain process the blocks in the chain.
func (hd *HistoryDB) ProcessBlockchain(bc *coin.Blockchain) error {
	depth := bc.Head().Seq()
	for i := uint64(0); i <= depth; i++ {
		b := bc.GetBlockInDepth(i)
		if err := hd.ProcessBlock(b); err != nil {
			return err
		}
	}
	return nil
}

// ProcessBlock will index the transaction, outputs,etc.
func (hd *HistoryDB) ProcessBlock(b *coin.Block) error {
	// index the transactions
	for _, t := range b.Body.Transactions {
		tx := Transaction{
			Transaction: t,
			BlockSeq:    b.Seq(),
		}
		if err := hd.txns.Add(&tx); err != nil {
			return err
		}

		// handle the tx in, we don't handle the genesis block has no in transaction.
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
				if err := hd.addrOut.Add(o.Address, o.UxId(o.CreateTxID)); err != nil {
					return err
				}
			}
		}

		// handle the tx out
		for _, o := range t.Out {
			out := Output{
				TransactionOutput: o,
				CreateTxID:        t.Hash(),
				CreatedBlockSeq:   b.Seq(),
			}
			// add output.
			if err := hd.outputs.Set(out); err != nil {
				return err
			}

			// index the output for address in.
			hd.addrIn.Add(o.Address, o.UxId(t.Hash()))
		}
	}
	return nil
}

// GetTransaction get transaction by hash.
func (hd *HistoryDB) GetTransaction(hash cipher.SHA256) (*Transaction, error) {
	return hd.txns.Get(hash)
}
