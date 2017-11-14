// Package historydb is in charge of parsing the consuses blokchain, and providing
// apis for blockchain explorer.
package historydb

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var logger = logging.MustGetLogger("historydb")

// HistoryDB provides apis for blockchain explorer.
type HistoryDB struct {
	db           *dbutil.DB    // bolt db instance.
	txns         *transactions // transactions bucket.
	outputs      *UxOuts       // outputs bucket.
	addrUx       *addressUx    // bucket which stores all UxOuts that address recved.
	addrTxns     *addressTxns  //  address related transaction bucket
	*historyMeta               // stores history meta info
}

// New create historydb instance and create corresponding buckets if does not exist.
func New(db *dbutil.DB) (*HistoryDB, error) {
	hd := HistoryDB{
		db: db,
	}

	var err error
	hd.txns, err = newTransactions(db)
	if err != nil {
		return nil, err
	}

	hd.outputs, err = newUxOuts(db)
	if err != nil {
		return nil, err
	}

	hd.addrUx, err = newAddressUx(db)
	if err != nil {
		return nil, err
	}

	hd.historyMeta, err = newHistoryMeta(db)
	if err != nil {
		return nil, err
	}

	hd.addrTxns, err = newAddressTxns(db)
	if err != nil {
		return nil, err
	}

	return &hd, nil
}

// ResetIfNeed checks if need to reset the parsed block history,
// If we have a new added bucket, we need to reset to parse
// blockchain again to get the new bucket filled.
func (hd *HistoryDB) ResetIfNeed(tx *bolt.Tx) error {
	if height, err := hd.historyMeta.ParsedHeight(tx); err != nil {
		return err
	} else if height == 0 {
		return nil
	}

	// if any of the following buckets are empty, need to reset
	addrTxnsEmpty, err := hd.addrTxns.IsEmpty(tx)
	if err != nil {
		return err
	}

	addrUxEmpty, err := hd.addrUx.IsEmpty(tx)
	if err != nil {
		return err
	}

	txnsEmpty, err := hd.txns.IsEmpty(tx)
	if err != nil {
		return err
	}

	outputsEmpty, err := hd.outputs.IsEmpty(tx)
	if err != nil {
		return err
	}

	if addrTxnsEmpty || addrUxEmpty || txnsEmpty || outputsEmpty {
		return hd.reset(tx)
	}

	return nil
}

func (hd *HistoryDB) reset(tx *bolt.Tx) error {
	logger.Info("History db reset")
	if err := hd.addrTxns.Reset(tx); err != nil {
		return err
	}

	if err := hd.addrUx.Reset(tx); err != nil {
		return err
	}

	if err := hd.outputs.Reset(tx); err != nil {
		return err
	}

	if err := hd.historyMeta.Reset(tx); err != nil {
		return err
	}

	return hd.txns.Reset(tx)
}

// GetUxout get UxOut of specific uxID.
func (hd *HistoryDB) GetUxout(tx *bolt.Tx, uxID cipher.SHA256) (*UxOut, error) {
	return hd.outputs.Get(tx, uxID)
}

// ParseBlock will index the transaction, outputs,etc.
func (hd *HistoryDB) ParseBlock(tx *bolt.Tx, b *coin.Block) error {
	if b == nil {
		return errors.New("process nil block")
	}

	// all updates will rollback if return error is not nil
	for _, t := range b.Body.Transactions {
		txn := Transaction{
			Tx:       t,
			BlockSeq: b.Seq(),
		}

		if err := hd.txns.Add(tx, &txn); err != nil {
			return err
		}

		// handle tx in, genesis transaction's vin is empty, so should be ignored.
		if b.Seq() > 0 {
			for _, in := range t.In {
				o, err := hd.outputs.Get(tx, in)
				if err != nil {
					return err
				}

				// update output's spent block seq and txid.
				o.SpentBlockSeq = b.Seq()
				o.SpentTxID = t.Hash()
				if err := hd.outputs.Set(tx, *o); err != nil {
					return err
				}

				// store the IN address with txid
				if err := hd.addrTxns.Add(tx, o.Out.Body.Address, t.Hash()); err != nil {
					return err
				}
			}
		}

		// handle the tx out
		uxArray := coin.CreateUnspents(b.Head, t)
		for _, ux := range uxArray {
			if err := hd.outputs.Set(tx, UxOut{Out: ux}); err != nil {
				return err
			}

			if err := hd.addrUx.Add(tx, ux.Body.Address, ux.Hash()); err != nil {
				return err
			}

			if err := hd.addrTxns.Add(tx, ux.Body.Address, t.Hash()); err != nil {
				return err
			}
		}
	}

	return hd.SetParsedHeight(tx, b.Seq())
}

// GetTransaction get transaction by hash.
func (hd HistoryDB) GetTransaction(tx *bolt.Tx, hash cipher.SHA256) (*Transaction, error) {
	return hd.txns.Get(tx, hash)
}

// GetLastTxs gets the latest N transactions.
func (hd HistoryDB) GetLastTxs(tx *bolt.Tx) ([]*Transaction, error) {
	txHashes := hd.txns.GetLastTxs()
	txns := make([]*Transaction, len(txHashes))

	for i, h := range txHashes {
		txn, err := hd.txns.Get(tx, h)
		if err != nil {
			return nil, err
		}
		txns[i] = txn
	}

	return txns, nil
}

// GetAddrUxOuts get all uxout that the address affected.
func (hd HistoryDB) GetAddrUxOuts(tx *bolt.Tx, address cipher.Address) ([]*UxOut, error) {
	hashes, err := hd.addrUx.Get(tx, address)
	if err != nil {
		return nil, err
	}

	uxOuts := make([]*UxOut, len(hashes))
	for i, hash := range hashes {
		ux, err := hd.outputs.Get(tx, hash)
		if err != nil {
			return nil, err
		}
		uxOuts[i] = ux
	}

	return uxOuts, nil
}

// GetAddrTxns returns all the address related transactions
func (hd HistoryDB) GetAddrTxns(tx *bolt.Tx, address cipher.Address) ([]Transaction, error) {
	hashes, err := hd.addrTxns.Get(tx, address)
	if err != nil {
		return nil, err
	}

	return hd.txns.GetSlice(tx, hashes)
}
