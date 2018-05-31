// Package historydb is in charge of parsing the consuses blokchain, and providing
// apis for blockchain explorer.
package historydb

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var logger = logging.MustGetLogger("historydb")

// CreateBuckets creates bolt.DB buckets used by the historydb
func CreateBuckets(tx *dbutil.Tx) error {
	return dbutil.CreateBuckets(tx, [][]byte{
		AddressTxnsBkt,
		AddressUxBkt,
		HistoryMetaBkt,
		UxOutsBkt,
		TransactionsBkt,
	})
}

// HistoryDB provides APIs for blockchain explorer
type HistoryDB struct {
	txns         *transactions // transactions bucket
	outputs      *UxOuts       // outputs bucket
	addrUx       *addressUx    // bucket which stores all UxOuts that address received
	addrTxns     *addressTxns  // address related transaction bucket
	*historyMeta               // stores history meta info
}

// New create HistoryDB instance
func New() *HistoryDB {
	return &HistoryDB{
		outputs:  &UxOuts{},
		txns:     &transactions{},
		addrUx:   &addressUx{},
		addrTxns: &addressTxns{},
	}
}

// NeedsReset checks if need to reset the parsed block history,
// If we have a new added bucket, we need to reset to parse
// blockchain again to get the new bucket filled.
func (hd *HistoryDB) NeedsReset(tx *dbutil.Tx) (bool, error) {
	if height, ok, err := hd.historyMeta.ParsedHeight(tx); err != nil {
		return false, err
	} else if !ok || height == 0 {
		return true, nil
	}

	// if any of the following buckets are empty, need to reset
	addrTxnsEmpty, err := hd.addrTxns.IsEmpty(tx)
	if err != nil {
		return false, err
	}

	addrUxEmpty, err := hd.addrUx.IsEmpty(tx)
	if err != nil {
		return false, err
	}

	txnsEmpty, err := hd.txns.IsEmpty(tx)
	if err != nil {
		return false, err
	}

	outputsEmpty, err := hd.outputs.IsEmpty(tx)
	if err != nil {
		return false, err
	}

	if addrTxnsEmpty || addrUxEmpty || txnsEmpty || outputsEmpty {
		return true, nil
	}

	return false, nil
}

// Erase erases the entire HistoryDB
func (hd *HistoryDB) Erase(tx *dbutil.Tx) error {
	logger.Debug("HistoryDB.reset")
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

// GetUxOuts get UxOut of specific uxIDs.
func (hd *HistoryDB) GetUxOuts(tx *dbutil.Tx, uxIDs []cipher.SHA256) ([]*UxOut, error) {
	return hd.outputs.GetArray(tx, uxIDs)
}

// ParseBlock builds indexes out of the block data
func (hd *HistoryDB) ParseBlock(tx *dbutil.Tx, b coin.Block) error {
	for _, t := range b.Body.Transactions {
		txn := Transaction{
			Tx:       t,
			BlockSeq: b.Seq(),
		}

		if err := hd.txns.Add(tx, &txn); err != nil {
			return err
		}

		for _, in := range t.In {
			o, err := hd.outputs.Get(tx, in)
			if err != nil {
				return err
			}

			if o == nil {
				return errors.New("HistoryDB.ParseBlock: transaction input not found in outputs bucket")
			}

			// update the output's spent block seq and txid
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

		// handle the tx out
		uxArray := coin.CreateUnspents(b.Head, t)
		for _, ux := range uxArray {
			if err := hd.outputs.Set(tx, UxOut{
				Out: ux,
			}); err != nil {
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
func (hd HistoryDB) GetTransaction(tx *dbutil.Tx, hash cipher.SHA256) (*Transaction, error) {
	return hd.txns.Get(tx, hash)
}

// GetAddrUxOuts get all uxout that the address affected.
func (hd HistoryDB) GetAddrUxOuts(tx *dbutil.Tx, address cipher.Address) ([]*UxOut, error) {
	hashes, err := hd.addrUx.Get(tx, address)
	if err != nil {
		return nil, err
	}

	return hd.outputs.GetArray(tx, hashes)
}

// GetAddressTxns returns all the address related transactions
func (hd HistoryDB) GetAddressTxns(tx *dbutil.Tx, address cipher.Address) ([]Transaction, error) {
	hashes, err := hd.addrTxns.Get(tx, address)
	if err != nil {
		return nil, err
	}

	return hd.txns.GetSlice(tx, hashes)
}

// ForEachTxn traverses the transactions bucket
func (hd HistoryDB) ForEachTxn(tx *dbutil.Tx, f func(cipher.SHA256, *Transaction) error) error {
	return hd.txns.ForEach(tx, f)
}
