package visor

import (
	"errors"

	"time"

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
}

// Hash returns the coin.Transaction's hash
func (ut *UnconfirmedTxn) Hash() cipher.SHA256 {
	return ut.Txn.Hash()
}

// unconfirmed transactions bucket
type uncfmTxnBkt struct {
	bkt *bucket.Bucket
}

func newUncfmTxBkt(db *bolt.DB) (*uncfmTxnBkt, error) {
	bkt, err := bucket.New([]byte("unconfirmed_txns"), db)
	if err != nil {
		return nil, err
	}
	return &uncfmTxnBkt{bkt: bkt}, nil
}

func (utb *uncfmTxnBkt) get(hash cipher.SHA256) (*UnconfirmedTxn, bool) {
	v := utb.bkt.Get([]byte(hash.Hex()))
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
	return utb.bkt.Put(key, d)
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

	return utb.bkt.Update([]byte(key.Hex()), updateFun)
}

func (utb *uncfmTxnBkt) delete(key cipher.SHA256) error {
	return utb.bkt.Delete([]byte(key.Hex()))
}

func (utb *uncfmTxnBkt) getAll() (map[cipher.SHA256]UnconfirmedTxn, error) {
	vs := utb.bkt.GetAll()
	txns := make(map[cipher.SHA256]UnconfirmedTxn, len(vs))
	for k, v := range vs {
		key, err := cipher.SHA256FromHex(k.(string))
		if err != nil {
			return nil, err
		}

		var tx UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &tx); err != nil {
			return nil, err
		}
		txns[key] = tx
	}
	return txns, nil
}

func (utb *uncfmTxnBkt) rangeUpdate(f func(key cipher.SHA256, tx *UnconfirmedTxn)) error {
	return utb.bkt.RangeUpdate(func(k, v []byte) ([]byte, error) {
		key, err := cipher.SHA256FromHex(string(k))
		if err != nil {
			return nil, err
		}

		var tx UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &tx); err != nil {
			return nil, err
		}

		f(key, &tx)

		// encoder the tx
		d := encoder.Serialize(tx)
		return d, nil
	})
}

func (utb *uncfmTxnBkt) isExist(key cipher.SHA256) bool {
	return utb.bkt.IsExist([]byte(key.Hex()))
}

func (utb *uncfmTxnBkt) forEach(f func(key cipher.SHA256, tx *UnconfirmedTxn) error) error {
	return utb.bkt.ForEach(func(k, v []byte) error {
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
	return utb.bkt.Len()
}

// UnconfirmedTxnPool manages unconfirmed transactions
type UnconfirmedTxnPool struct {
	// Txns map[cipher.SHA256]UnconfirmedTxn
	Txns *uncfmTxnBkt
	// Predicted unspents, assuming txns are valid.  Needed to predict
	// our future balance and avoid double spending our own coins
	// Maps from Transaction.Hash() to UxArray.
	Unspent TxnUnspents
}

// NewUnconfirmedTxnPool creates an UnconfirmedTxnPool instance
func NewUnconfirmedTxnPool(db *bolt.DB) *UnconfirmedTxnPool {
	txnBkt, err := newUncfmTxBkt(db)
	if err != nil {
		panic(err)
	}

	return &UnconfirmedTxnPool{
		Txns:    txnBkt,
		Unspent: make(TxnUnspents),
	}
}

// SetAnnounced updates announced time of specific tx
func (utp *UnconfirmedTxnPool) SetAnnounced(h cipher.SHA256, t time.Time) {
	utp.Txns.update(h, func(tx *UnconfirmedTxn) {
		tx.Announced = t.UnixNano()
	})

	// if tx, ok := utp.Txns[h]; ok {
	// 	tx.Announced = t.UnixNano()
	// 	utp.Txns[h] = tx
	// }
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
func (utp *UnconfirmedTxnPool) InjectTxn(bc *Blockchain,
	t coin.Transaction) (error, bool) {

	if err := t.Verify(); err != nil {
		return err, false
	}

	if err := VerifyTransactionFee(bc, &t); err != nil {
		return err, false
	}
	if err := bc.VerifyTransaction(t); err != nil {
		return err, false
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
	})

	if exist {
		return nil, true
	}
	// ut, ok := utp.Txns[h]
	// if ok {
	// 	now := util.Now()
	// 	ut.Received = now.UnixNano()
	// 	ut.Checked = now.UnixNano()
	// 	utp.Txns[h] = ut
	// 	return nil, true
	// }

	// Add txn to index
	unspent := bc.GetUnspent()
	utx := utp.createUnconfirmedTxn(unspent, t)
	utp.Txns.put(&utx)
	// utp.Txns[h] = utp.createUnconfirmedTxn(unspent, t)
	// Add predicted unspents
	utp.Unspent[h] = coin.CreateUnspents(bc.Head().Head, t)

	return nil, false
}

// RawTxns returns underlying coin.Transactions
func (utp *UnconfirmedTxnPool) RawTxns() coin.Transactions {
	allUtx, err := utp.Txns.getAll()
	if err != nil {
		return coin.Transactions{}
	}

	txns := make(coin.Transactions, len(allUtx))
	i := 0
	for _, t := range allUtx {
		txns[i] = t.Txn
		i++
	}
	return txns
}

// Remove a single txn by hash
func (utp *UnconfirmedTxnPool) removeTxn(bc *Blockchain, txHash cipher.SHA256) {
	// delete(utp.Txns, txHash)
	utp.Txns.delete(txHash)
	delete(utp.Unspent, txHash)
}

// Removes multiple txns at once. Slightly more efficient than a series of
// single RemoveTxns.  Hashes is an array of Transaction hashes.
func (utp *UnconfirmedTxnPool) removeTxns(bc *Blockchain,
	hashes []cipher.SHA256) {
	for i := range hashes {
		// delete(utp.Txns, hashes[i])
		utp.Txns.delete(hashes[i])
		delete(utp.Unspent, hashes[i])
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

// Refresh checks all unconfirmed txns against the blockchain. maxAge is how long
// we'll hold a txn regardless of whether it has been invalidated.
// checkPeriod is how often we check the txn against the blockchain.
func (utp *UnconfirmedTxnPool) Refresh(bc *Blockchain,
	checkPeriod, maxAge time.Duration) {

	now := util.Now()
	var toRemove []cipher.SHA256
	utp.Txns.rangeUpdate(func(key cipher.SHA256, tx *UnconfirmedTxn) {
		if now.Sub(nanoToTime(tx.Received)) >= maxAge {
			toRemove = append(toRemove, key)
		} else if now.Sub(nanoToTime(tx.Checked)) >= checkPeriod {
			if bc.VerifyTransaction(tx.Txn) == nil {
				tx.Checked = now.UnixNano()
			} else {
				toRemove = append(toRemove, key)
			}
		}
	})

	// for k, t := range utp.Txns {
	// 	if now.Sub(nanoToTime(t.Received)) >= maxAge {
	// 		toRemove = append(toRemove, k)
	// 	} else if now.Sub(nanoToTime(t.Checked)) >= checkPeriod {
	// 		if bc.VerifyTransaction(t.Txn) == nil {
	// 			t.Checked = now.UnixNano()
	// 			utp.Txns[k] = t
	// 		} else {
	// 			toRemove = append(toRemove, k)
	// 		}
	// 	}
	// }
	utp.removeTxns(bc, toRemove)
}

// FilterKnown returns txn hashes with known ones removed
func (utp *UnconfirmedTxnPool) FilterKnown(txns []cipher.SHA256) []cipher.SHA256 {
	var unknown []cipher.SHA256
	for _, h := range txns {
		if !utp.Txns.isExist(h) {
			unknown = append(unknown, h)
		}
		// if _, known := utp.Txns[h]; !known {
		// 	unknown = append(unknown, h)
		// }
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
		// if txn, have := utp.Txns[h]; have {
		// 	known = append(known, txn.Txn)
		// }
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
	// for _, utx := range utp.Txns {
	// 	for _, h := range utx.Txn.In {
	// 		if ux, ok := bcUnspent.Get(h); ok {
	// 			if _, ok := a[ux.Body.Address]; ok {
	// 				auxs[ux.Body.Address] = append(auxs[ux.Body.Address], ux)
	// 			}
	// 		}
	// 	}
	// }
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
	// for _, tx := range utp.Txns {
	// 	for _, in := range tx.Txn.In {
	// 		if ux, ok := bcUnspent.Get(in); ok {
	// 			outs = append(outs, NewReadableOutput(ux))
	// 		}
	// 	}
	// }
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
	// for _, tx := range utp.Txns {
	// }
	return outs
}

// Get returns the unconfirmed transaction of given tx hash.
func (utp *UnconfirmedTxnPool) Get(key cipher.SHA256) (*UnconfirmedTxn, bool) {
	return utp.Txns.get(key)
}

// GetAllUnconfirmedTxns returns all unconfirmed transactions array
func (utp *UnconfirmedTxnPool) GetAllUnconfirmedTxns() []UnconfirmedTxn {
	all, err := utp.Txns.getAll()
	if err != nil {
		return []UnconfirmedTxn{}
	}

	txns := make([]UnconfirmedTxn, 0, len(all))
	for _, tx := range all {
		txns = append(txns, tx)
	}
	return txns
}

func nanoToTime(n int64) time.Time {
	return time.Unix(n/int64(time.Second), n%int64(time.Second))
}
