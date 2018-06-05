// Package historydb is in charge of parsing the consuses blokchain, and providing
// apis for blockchain explorer.
package historydb

import (
	"errors"
	"fmt"
	"sync"

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

// IndexesMap is a goroutine safe address indexes map
type IndexesMap struct {
	value map[cipher.Address]AddressIndexes
	lock  sync.RWMutex
}

// NewIndexesMap creates a IndexesMap instance
func NewIndexesMap() *IndexesMap {
	return &IndexesMap{
		value: make(map[cipher.Address]AddressIndexes),
	}
}

// Load returns value of given key
func (im *IndexesMap) Load(address cipher.Address) (AddressIndexes, bool) {
	im.lock.RLock()
	v, ok := im.value[address]
	im.lock.RUnlock()
	return v, ok
}

// Store saves address with indexes
func (im *IndexesMap) Store(address cipher.Address, indexes AddressIndexes) {
	im.lock.Lock()
	im.value[address] = indexes
	im.lock.Unlock()
}

// AddressIndexes represents the address indexes struct
type AddressIndexes struct {
	TxnHashes map[cipher.SHA256]struct{}
	UxHashes  map[cipher.SHA256]struct{}
}

// Verify checks if the historydb is corrupted
func (hd HistoryDB) Verify(tx *dbutil.Tx, b *coin.SignedBlock, indexesMap *IndexesMap) error {
	for _, t := range b.Body.Transactions {
		txnHash := t.Hash()
		txn, err := hd.txns.Get(tx, txnHash)
		if err != nil {
			return err
		}

		if txn == nil {
			err := fmt.Errorf("HistoryDB.Verify: transaction %v does not exist in historydb", txnHash.Hex())
			return ErrHistoryDBCorrupted{err}
		}

		for _, in := range t.In {
			// Checks the existence of transaction input
			o, err := hd.outputs.Get(tx, in)
			if err != nil {
				return err
			}

			if o == nil {
				err := fmt.Errorf("HistoryDB.Verify: transaction input %v does not exist in historydb", in.Hex())
				return ErrHistoryDBCorrupted{err}
			}

			// Checks the output's spend block seq
			if o.SpentBlockSeq != b.Seq() {
				err := fmt.Errorf("HistoryDB.Verify: spend block seq of transaction input %v is wrong, should be: %v, but is %v",
					in.Hex(), b.Seq(), o.SpentBlockSeq)
				return ErrHistoryDBCorrupted{err}
			}

			addr := o.Out.Body.Address
			txnHashesMap := map[cipher.SHA256]struct{}{}
			uxHashesMap := map[cipher.SHA256]struct{}{}

			// Checks if the address indexes already loaded into memory
			indexes, ok := indexesMap.Load(addr)
			if ok {
				txnHashesMap = indexes.TxnHashes
				uxHashesMap = indexes.UxHashes
			} else {
				txnHashes, err := hd.addrTxns.Get(tx, addr)
				if err != nil {
					return err
				}
				for _, hash := range txnHashes {
					txnHashesMap[hash] = struct{}{}
				}

				uxHashes, err := hd.addrUx.Get(tx, addr)
				if err != nil {
					return err
				}
				for _, hash := range uxHashes {
					uxHashesMap[hash] = struct{}{}
				}

				indexesMap.Store(addr, AddressIndexes{
					TxnHashes: txnHashesMap,
					UxHashes:  uxHashesMap,
				})
			}

			if _, ok := txnHashesMap[txnHash]; !ok {
				err := fmt.Errorf("HistoryDB.Verify: index of address transaction [%s:%s] does not exist in historydb",
					addr, txnHash.Hex())
				return ErrHistoryDBCorrupted{err}
			}

			if _, ok := uxHashesMap[in]; !ok {
				err := fmt.Errorf("HistoryDB.Verify: index of address uxout [%s:%s] does not exist in historydb",
					addr, in.Hex())
				return ErrHistoryDBCorrupted{err}
			}
		}

		// Checks the transaction outs
		uxArray := coin.CreateUnspents(b.Head, t)
		for _, ux := range uxArray {
			uxHash := ux.Hash()
			out, err := hd.outputs.Get(tx, uxHash)
			if err != nil {
				return err
			}

			if out == nil {
				err := fmt.Errorf("HistoryDB.Verify: transaction output %s does not exist in historydb", uxHash.Hex())
				return ErrHistoryDBCorrupted{err}
			}

			addr := ux.Body.Address
			txnHashesMap := map[cipher.SHA256]struct{}{}
			uxHashesMap := map[cipher.SHA256]struct{}{}
			indexes, ok := indexesMap.Load(addr)
			if ok {
				txnHashesMap = indexes.TxnHashes
				uxHashesMap = indexes.UxHashes
			} else {
				txnHashes, err := hd.addrTxns.Get(tx, addr)
				if err != nil {
					return err
				}
				for _, hash := range txnHashes {
					txnHashesMap[hash] = struct{}{}
				}

				uxHashes, err := hd.addrUx.Get(tx, addr)
				if err != nil {
					return err
				}

				for _, hash := range uxHashes {
					uxHashesMap[hash] = struct{}{}
				}

				indexesMap.Store(addr, AddressIndexes{
					TxnHashes: txnHashesMap,
					UxHashes:  uxHashesMap,
				})
			}

			if _, ok := txnHashesMap[txnHash]; !ok {
				err := fmt.Errorf("HistoryDB.Verify: index of address transaction [%s:%s] does not exist in historydb",
					addr, txnHash.Hex())
				return ErrHistoryDBCorrupted{err}
			}
		}
	}
	return nil
}

// ErrHistoryDBCorrupted is returned when found the historydb is corrupted
type ErrHistoryDBCorrupted struct {
	error
}

// NewErrHistoryDBCorrupted is for user to be able to create ErrHistoryDBCorrupted instance
// outside of the package
func NewErrHistoryDBCorrupted(err error) ErrHistoryDBCorrupted {
	return ErrHistoryDBCorrupted{err}
}
