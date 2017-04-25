package visor

import (
	"errors"

	"time"

	"fmt"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// BurnFactor half of coinhours must be burnt
var BurnFactor uint64 = 2

// VerifyTransactionFee performs additional transaction verification at the unconfirmed pool level.
// This checks tunable parameters that should prevent the transaction from
// entering the blockchain, but cannot be done at the blockchain level because
// they may be changed.
func VerifyTransactionFee(bc *Blockchain, t *coin.Transaction) error {
	fee, err := bc.TransactionFee(t)
	if err != nil {
		return err
	}

	//calculate total number of coinhours
	var total = t.OutputHours() + fee
	//make sure at least half the coin hours are destroyed
	if fee < total/BurnFactor {
		return errors.New("Transaction coinhour fee minimum not met")
	}
	return nil
}

// TxnUnspents maps from coin.Transaction hash to its expected unspents.  The unspents'
// Head can be different at execution time, but the Unspent's hash is fixed.
type TxnUnspents map[cipher.SHA256]coin.UxArray

// AllForAddress returns all Unspents for a single address
func (tus TxnUnspents) AllForAddress(a cipher.Address) coin.UxArray {
	uxo := make(coin.UxArray, 0)
	for _, uxa := range tus {
		for i := range uxa {
			if uxa[i].Body.Address == a {
				uxo = append(uxo, uxa[i])
			}
		}
	}
	return uxo
}

type UnconfirmedTxn struct {
	Txn coin.Transaction
	// Time the txn was last received
	Received int64
	// Time the txn was last checked against the blockchain
	Checked int64
	// Last time we announced this txn
	Announced int64
	// If this txn is valid
	IsValid bool
}

// Hash returns the coin.Transaction's hash
func (ut *UnconfirmedTxn) Hash() cipher.SHA256 {
	return ut.Txn.Hash()
}

// unconfirmed transactions bucket
type uncfmTxnBkt struct {
	txns      *bucket.Bucket
	idx       *bucket.Bucket
	indexName []byte
}

func newUncfmTxBkt(db *bolt.DB) *uncfmTxnBkt {
	bkt, err := bucket.New([]byte("unconfirmed_txns"), db)
	if err != nil {
		panic(err)
	}

	idx, err := bucket.New([]byte("unconfirmed_index"), db)
	if err != nil {
		panic(err)
	}

	indexName := []byte("index")
	return &uncfmTxnBkt{txns: bkt, idx: idx, indexName: indexName}
}

func (utb *uncfmTxnBkt) get(hash cipher.SHA256) (*UnconfirmedTxn, bool) {
	v := utb.txns.Get([]byte(hash.Hex()))
	if v == nil {
		return nil, false
	}
	var tx UnconfirmedTxn
	if err := encoder.DeserializeRaw(v, &tx); err != nil {
		return nil, false
	}
	return &tx, true
}

func (utb *uncfmTxnBkt) put(v *UnconfirmedTxn) error {
	key := []byte(v.Hash().Hex())
	d := encoder.Serialize(v)
	// check if exist
	exist := utb.txns.IsExist(key)
	if err := utb.txns.Put(key, d); err != nil {
		return err
	}

	if exist {
		return nil
	}

	indexes := utb.mustGetTxHashes()
	indexes = append(indexes, v.Hash())
	return utb.idx.Put(utb.indexName, encoder.Serialize(indexes))
}

func (utb *uncfmTxnBkt) update(key cipher.SHA256, f func(v *UnconfirmedTxn)) error {
	updateFun := func(v []byte) ([]byte, error) {
		var tx UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &tx); err != nil {
			return nil, err
		}

		f(&tx)
		return encoder.Serialize(tx), nil
	}

	return utb.txns.Update([]byte(key.Hex()), updateFun)
}

func (utb *uncfmTxnBkt) delete(key cipher.SHA256) error {
	if err := utb.txns.Delete([]byte(key.Hex())); err != nil {
		return err
	}

	// remove the index
	return utb.idx.Update(utb.indexName, func(idx []byte) ([]byte, error) {
		var indexes []cipher.SHA256
		if err := encoder.DeserializeRaw(idx, &indexes); err != nil {
			return nil, err
		}

		for i := range indexes {
			if indexes[i] == key {
				indexes = append(indexes[:i], indexes[i+1:]...)
				break
			}
		}
		return encoder.Serialize(indexes), nil
	})
}

func (utb *uncfmTxnBkt) getAll() ([]UnconfirmedTxn, error) {
	vs := utb.txns.GetAll()
	indexes := utb.mustGetTxHashes()

	if len(indexes) != len(vs) {
		return nil, fmt.Errorf("index size not match transaction size, index:%d, txn:%d", len(indexes), len(vs))
	}

	txns := make([]UnconfirmedTxn, 0, len(indexes))
	for _, ix := range indexes {
		v, ok := vs[ix.Hex()]
		if !ok {
			return nil, fmt.Errorf("There's no tx with hash:%s", ix.Hex())
		}

		var tx UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &tx); err != nil {
			return nil, err
		}
		txns = append(txns, tx)
	}

	return txns, nil
}

func (utb *uncfmTxnBkt) mustGetTxHashes() []cipher.SHA256 {
	idx := utb.idx.Get(utb.indexName)
	if idx == nil {
		return []cipher.SHA256{}
	}
	var indexes []cipher.SHA256
	if err := encoder.DeserializeRaw(idx, &indexes); err != nil {
		panic(err)
	}

	return indexes
}

func (utb *uncfmTxnBkt) getSlice(keys []cipher.SHA256) ([]UnconfirmedTxn, error) {
	ks := make([][]byte, len(keys))
	for i := range keys {
		ks[i] = []byte(keys[i].Hex())
	}
	vs := utb.txns.GetSlice(ks)
	txns := make([]UnconfirmedTxn, len(vs))
	for i := range vs {
		var txn UnconfirmedTxn
		if err := encoder.DeserializeRaw(vs[i], &txn); err != nil {
			return []UnconfirmedTxn{}, err
		}
		txns[i] = txn
	}
	return txns, nil
}

func (utb *uncfmTxnBkt) rangeUpdate(f func(key cipher.SHA256, tx *UnconfirmedTxn)) error {
	return utb.txns.RangeUpdate(func(k, v []byte) ([]byte, error) {
		key, err := cipher.SHA256FromHex(string(k))
		if err != nil {
			return nil, err
		}

		var tx UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &tx); err != nil {
			return nil, err
		}
		f(key, &tx)
		// encode the tx
		d := encoder.Serialize(tx)
		return d, nil
	})
}

func (utb *uncfmTxnBkt) isExist(key cipher.SHA256) bool {
	return utb.txns.IsExist([]byte(key.Hex()))
}

func (utb *uncfmTxnBkt) forEach(f func(key cipher.SHA256, tx *UnconfirmedTxn) error) error {
	return utb.txns.ForEach(func(k, v []byte) error {
		key, err := cipher.SHA256FromHex(string(k))
		if err != nil {
			return err
		}
		var tx UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &tx); err != nil {
			return err
		}

		return f(key, &tx)
	})
}

func (utb *uncfmTxnBkt) len() int {
	// exclude the index
	return utb.txns.Len()
}

func (utb *uncfmTxnBkt) indexLen() int {
	idx := utb.mustGetTxHashes()
	return len(idx)
}

// returns the first N transactions in the pool
func (utb *uncfmTxnBkt) mustFirstN(n int) []UnconfirmedTxn {
	// get the indexes
	indexes := utb.mustGetTxHashes()
	idxs := indexes[:n]
	txns := make([]UnconfirmedTxn, len(idxs))
	for i, hash := range idxs {
		tx, ok := utb.get(hash)
		if !ok {
			panic(fmt.Errorf("hash %s not exist in unconfirmed txns bucket, while it's exist in indexes", hash.Hex()))
		}
		txns[i] = *tx
	}
	return txns
}

type txUnspents struct {
	bkt *bucket.Bucket
}

func newTxUnspents(db *bolt.DB) *txUnspents {
	bkt, err := bucket.New([]byte("unconfirmed_unspents"), db)
	if err != nil {
		panic(err)
	}

	return &txUnspents{bkt: bkt}
}

func (txus *txUnspents) put(key cipher.SHA256, uxs coin.UxArray) error {
	v := encoder.Serialize(uxs)
	return txus.bkt.Put([]byte(key.Hex()), v)
}

func (txus *txUnspents) get(key cipher.SHA256) (coin.UxArray, error) {
	v := txus.bkt.Get([]byte(key.Hex()))
	var uxs coin.UxArray
	if err := encoder.DeserializeRaw(v, &uxs); err != nil {
		return coin.UxArray{}, err
	}
	return uxs, nil
}

func (txus *txUnspents) len() int {
	return txus.bkt.Len()
}

func (txus *txUnspents) delete(key cipher.SHA256) error {
	return txus.bkt.Delete([]byte(key.Hex()))
}

func (txus *txUnspents) getAllForAddress(a cipher.Address) (uxo coin.UxArray) {
	txus.bkt.ForEach(func(k, v []byte) error {
		var uxa coin.UxArray
		if err := encoder.DeserializeRaw(v, &uxa); err != nil {
			panic(err)
		}

		for i := range uxa {
			if uxa[i].Body.Address == a {
				uxo = append(uxo, uxa[i])
			}
		}
		return nil
	})
	return
}

func (txus *txUnspents) forEach(f func(cipher.SHA256, coin.UxArray)) error {
	return txus.bkt.ForEach(func(k, v []byte) error {
		hash, err := cipher.SHA256FromHex(string(k))
		if err != nil {
			return err
		}

		var uxa coin.UxArray
		if err := encoder.DeserializeRaw(v, &uxa); err != nil {
			return err
		}

		f(hash, uxa)
		return nil
	})
}

// UnconfirmedTxnPool manages unconfirmed transactions
type UnconfirmedTxnPool struct {
	// Txns map[cipher.SHA256]UnconfirmedTxn
	Txns *uncfmTxnBkt
	// Predicted unspents, assuming txns are valid.  Needed to predict
	// our future balance and avoid double spending our own coins
	// Maps from Transaction.Hash() to UxArray.
	Unspent *txUnspents
}

// NewUnconfirmedTxnPool creates an UnconfirmedTxnPool instance
func NewUnconfirmedTxnPool(db *bolt.DB) *UnconfirmedTxnPool {
	return &UnconfirmedTxnPool{
		Txns:    newUncfmTxBkt(db),
		Unspent: newTxUnspents(db),
	}
}

// SetAnnounced updates announced time of specific tx
func (utp *UnconfirmedTxnPool) SetAnnounced(h cipher.SHA256, t time.Time) {
	utp.Txns.update(h, func(tx *UnconfirmedTxn) {
		tx.Announced = t.UnixNano()
	})
}

// Creates an unconfirmed transaction
func (utp *UnconfirmedTxnPool) createUnconfirmedTxn(bcUnsp *coin.UnspentPool,
	t coin.Transaction) UnconfirmedTxn {
	now := util.Now()
	return UnconfirmedTxn{
		Txn:       t,
		Received:  now.UnixNano(),
		Checked:   now.UnixNano(),
		Announced: util.ZeroTime().UnixNano(),
	}
}

// InjectTxn adds a coin.Transaction to the pool, or updates an existing one's timestamps
// Returns an error if txn is invalid, and whether the transaction already
// existed in the pool.
func (utp *UnconfirmedTxnPool) InjectTxn(bc *Blockchain, t coin.Transaction) (error, bool) {
	var valid bool
	for {
		if err := VerifyTransactionFee(bc, &t); err != nil {
			if err == ErrUnspentNotExist {
				break
			}
			return err, false
		}

		if err := bc.VerifyTransaction(t); err != nil {
			return err, false
		}

		valid = true
		break
	}

	// Update if we already have this txn
	h := t.Hash()
	// update the time if exist
	var exist bool
	utp.Txns.update(h, func(tx *UnconfirmedTxn) {
		exist = true
		now := util.Now()
		tx.Received = now.UnixNano()
		tx.Checked = now.UnixNano()
		tx.IsValid = valid
	})

	if exist {
		return nil, true
	}

	// Add txn to index
	unspent := bc.GetUnspent()
	utx := utp.createUnconfirmedTxn(unspent, t)
	utx.IsValid = valid
	utp.Txns.put(&utx)
	utp.Unspent.put(h, coin.CreateUnspents(bc.Head().Head, t))

	return nil, false
}

// RawTxns returns underlying coin.Transactions
func (utp *UnconfirmedTxnPool) RawTxns() coin.Transactions {
	utxns, err := utp.Txns.getAll()
	if err != nil {
		return coin.Transactions{}
	}

	txns := make(coin.Transactions, len(utxns))
	for i := range utxns {
		txns[i] = utxns[i].Txn
	}
	return txns
}

// Remove a single txn by hash
func (utp *UnconfirmedTxnPool) removeTxn(bc *Blockchain, txHash cipher.SHA256) {
	// delete(utp.Txns, txHash)
	utp.Txns.delete(txHash)
	utp.Unspent.delete(txHash)
}

// Removes multiple txns at once. Slightly more efficient than a series of
// single RemoveTxns.  Hashes is an array of Transaction hashes.
func (utp *UnconfirmedTxnPool) removeTxns(bc *Blockchain,
	hashes []cipher.SHA256) {
	for i := range hashes {
		// delete(utp.Txns, hashes[i])
		utp.Txns.delete(hashes[i])
		utp.Unspent.delete(hashes[i])
	}
}

// RemoveTransactions removes confirmed txns from the pool
func (utp *UnconfirmedTxnPool) RemoveTransactions(bc *Blockchain,
	txns coin.Transactions) {
	toRemove := make([]cipher.SHA256, len(txns))
	for i := range txns {
		toRemove[i] = txns[i].Hash()
	}
	utp.removeTxns(bc, toRemove)
}

// MustFirstN returns the first N unconfirmed transaction in the pool
func (utp *UnconfirmedTxnPool) MustFirstN(n int) []UnconfirmedTxn {
	return utp.Txns.mustFirstN(n)
}

// Refresh checks all unconfirmed txns against the blockchain.
// verify the transaction and returns all those txns that turn to valid.
func (utp *UnconfirmedTxnPool) Refresh(bc *Blockchain) (hashes []cipher.SHA256) {
	now := util.Now()
	utp.Txns.rangeUpdate(func(key cipher.SHA256, tx *UnconfirmedTxn) {
		tx.Checked = now.UnixNano()
		if !tx.IsValid {
			if bc.VerifyTransaction(tx.Txn) == nil {
				tx.IsValid = true
				hashes = append(hashes, tx.Hash())
			}
		}
	})

	return
}

// FilterKnown returns txn hashes with known ones removed
func (utp *UnconfirmedTxnPool) FilterKnown(txns []cipher.SHA256) []cipher.SHA256 {
	var unknown []cipher.SHA256
	for _, h := range txns {
		if !utp.Txns.isExist(h) {
			unknown = append(unknown, h)
		}
	}
	return unknown
}

// GetKnown returns all known coin.Transactions from the pool, given hashes to select
func (utp *UnconfirmedTxnPool) GetKnown(txns []cipher.SHA256) coin.Transactions {
	var known coin.Transactions
	for _, h := range txns {
		if tx, ok := utp.Txns.get(h); ok {
			known = append(known, tx.Txn)
		}
	}
	return known
}

// SpendsForAddresses returns all unconfirmed coin.UxOut spends for addresses
// Looks at all inputs for unconfirmed txns, gets their source UxOut from the
// blockchain's unspent pool, and returns as coin.AddressUxOuts
func (utp *UnconfirmedTxnPool) SpendsForAddresses(bcUnspent *coin.UnspentPool,
	a map[cipher.Address]byte) coin.AddressUxOuts {
	auxs := make(coin.AddressUxOuts, len(a))
	utp.Txns.forEach(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		for _, h := range tx.Txn.In {
			if ux, ok := bcUnspent.Get(h); ok {
				if _, ok := a[ux.Body.Address]; ok {
					auxs[ux.Body.Address] = append(auxs[ux.Body.Address], ux)
				}
			}
		}
		return nil
	})
	return auxs
}

func (utp *UnconfirmedTxnPool) SpendsForAddress(bcUnspent *coin.UnspentPool,
	a cipher.Address) coin.UxArray {
	ma := map[cipher.Address]byte{a: 1}
	auxs := utp.SpendsForAddresses(bcUnspent, ma)
	return auxs[a]
}

// AllSpendsOutputs returns all spending outputs in unconfirmed tx pool.
func (utp *UnconfirmedTxnPool) AllSpendsOutputs(bcUnspent *coin.UnspentPool) []ReadableOutput {
	outs := []ReadableOutput{}
	utp.Txns.forEach(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		for _, in := range tx.Txn.In {
			if ux, ok := bcUnspent.Get(in); ok {
				outs = append(outs, NewReadableOutput(ux))
			}
		}
		return nil
	})
	return outs
}

// AllIncommingOutputs returns all predicted incomming outputs.
func (utp *UnconfirmedTxnPool) AllIncommingOutputs(bh coin.BlockHeader) []ReadableOutput {
	outs := []ReadableOutput{}
	utp.Txns.forEach(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		uxOuts := coin.CreateUnspents(bh, tx.Txn)
		for _, ux := range uxOuts {
			outs = append(outs, NewReadableOutput(ux))
		}
		return nil
	})
	return outs
}

// Get returns the unconfirmed transaction of given tx hash.
func (utp *UnconfirmedTxnPool) Get(key cipher.SHA256) (*UnconfirmedTxn, bool) {
	return utp.Txns.get(key)
}

// GetAllUnconfirmedTxns returns all unconfirmed transactions array
func (utp *UnconfirmedTxnPool) GetAllUnconfirmedTxns() []UnconfirmedTxn {
	txns, err := utp.Txns.getAll()
	if err != nil {
		return []UnconfirmedTxn{}
	}

	return txns
}

// GetAllTxnHashes returns all unconfirmed txns hashes
func (utp *UnconfirmedTxnPool) GetAllTxnHashes() []cipher.SHA256 {
	return utp.Txns.mustGetTxHashes()
}

// Len returns the number of unconfirmed transactions
func (utp *UnconfirmedTxnPool) Len() int {
	return utp.Txns.len()
}

func nanoToTime(n int64) time.Time {
	zeroTime := time.Time{}
	if n == zeroTime.UnixNano() {
		// maximum time
		return zeroTime
	}
	return time.Unix(n/int64(time.Second), n%int64(time.Second))
}
