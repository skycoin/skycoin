package visor

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

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

// UnconfirmedTxn unconfirmed transaction
type UnconfirmedTxn struct {
	Txn coin.Transaction
	// Time the txn was last received
	Received int64
	// Time the txn was last checked against the blockchain
	Checked int64
	// Last time we announced this txn
	Announced int64
	// If this txn is valid
	IsValid int8
}

// Hash returns the coin.Transaction's hash
func (ut *UnconfirmedTxn) Hash() cipher.SHA256 {
	return ut.Txn.Hash()
}

// unconfirmed transactions bucket
type uncfmTxnBkt struct {
	txns *bucket.Bucket
}

func newUncfmTxBkt(db *bolt.DB) *uncfmTxnBkt {
	bkt, err := bucket.New([]byte("unconfirmed_txns"), db)
	if err != nil {
		panic(err)
	}

	return &uncfmTxnBkt{txns: bkt}
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

func (utb *uncfmTxnBkt) putWithTx(tx *bolt.Tx, v *UnconfirmedTxn) error {
	key := []byte(v.Hash().Hex())
	d := encoder.Serialize(v)
	return utb.txns.PutWithTx(tx, key, d)
}

func (utb *uncfmTxnBkt) update(key cipher.SHA256, f func(v *UnconfirmedTxn)) error {
	updateFun := func(v []byte) ([]byte, error) {
		if v == nil {
			return nil, fmt.Errorf("%s does not exist in bucket %s", key.Hex(), utb.txns.Name)
		}

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
	return utb.txns.Delete([]byte(key.Hex()))
}

func (utb *uncfmTxnBkt) deleteWithTx(tx *bolt.Tx, key cipher.SHA256) error {
	return utb.txns.DeleteWithTx(tx, []byte(key.Hex()))
}

func (utb *uncfmTxnBkt) getAll() ([]UnconfirmedTxn, error) {
	vs := utb.txns.GetAll()
	txns := make([]UnconfirmedTxn, 0, len(vs))
	for _, u := range vs {
		var tx UnconfirmedTxn
		if err := encoder.DeserializeRaw(u, &tx); err != nil {
			return nil, err
		}
		txns = append(txns, tx)
	}

	return txns, nil
}

func (utb *uncfmTxnBkt) rangeUpdate(f func(key cipher.SHA256, tx *UnconfirmedTxn) error) error {
	return utb.txns.RangeUpdate(func(k, v []byte) ([]byte, error) {
		key, err := cipher.SHA256FromHex(string(k))
		if err != nil {
			return nil, err
		}

		var tx UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &tx); err != nil {
			return nil, err
		}

		if err := f(key, &tx); err != nil {
			return nil, err
		}

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

func (txus *txUnspents) putWithTx(tx *bolt.Tx, key cipher.SHA256, uxs coin.UxArray) error {
	v := encoder.Serialize(uxs)
	return txus.bkt.PutWithTx(tx, []byte(key.Hex()), v)
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

func (txus *txUnspents) deleteWithTx(tx *bolt.Tx, key cipher.SHA256) error {
	return txus.bkt.DeleteWithTx(tx, []byte(key.Hex()))
}

func (txus *txUnspents) getByAddr(a cipher.Address) (uxo coin.UxArray) {
	txus.bkt.ForEach(func(k, v []byte) error {
		var uxa coin.UxArray
		if err := encoder.DeserializeRaw(v, &uxa); err != nil {
			return err
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
	txns *uncfmTxnBkt
	// Predicted unspents, assuming txns are valid.  Needed to predict
	// our future balance and avoid double spending our own coins
	// Maps from Transaction.Hash() to UxArray.
	unspent *txUnspents
}

// NewUnconfirmedTxnPool creates an UnconfirmedTxnPool instance
func NewUnconfirmedTxnPool(db *bolt.DB) *UnconfirmedTxnPool {
	return &UnconfirmedTxnPool{
		txns:    newUncfmTxBkt(db),
		unspent: newTxUnspents(db),
	}
}

// SetAnnounced updates announced time of specific tx
func (utp *UnconfirmedTxnPool) SetAnnounced(h cipher.SHA256, t int64) error {
	return utp.txns.update(h, func(tx *UnconfirmedTxn) {
		tx.Announced = t
	})
}

// Creates an unconfirmed transaction
func (utp *UnconfirmedTxnPool) createUnconfirmedTxn(t coin.Transaction) UnconfirmedTxn {
	now := utc.Now()
	return UnconfirmedTxn{
		Txn:       t,
		Received:  now.UnixNano(),
		Checked:   now.UnixNano(),
		Announced: time.Time{}.UnixNano(),
	}
}

// InjectTransaction adds a coin.Transaction to the pool, or updates an existing one's timestamps
// Returns an error if txn is invalid, and whether the transaction already
// existed in the pool.
// If the transaction violates hard constraints, it is rejected.
// Soft constraints violations mark a txn as invalid, but the txn is inserted. The soft violation is returned.
func (utp *UnconfirmedTxnPool) InjectTransaction(bc Blockchainer, t coin.Transaction, maxSize int) (bool, *ErrTxnViolatesSoftConstraint, error) {
	var isValid int8 = 1
	var softErr *ErrTxnViolatesSoftConstraint
	if err := bc.VerifySingleTxnAllConstraints(t, maxSize); err != nil {
		logger.Warningf("bc.VerifySingleTxnAllConstraints failed for txn %s: %v", t.TxIDHex(), err)
		switch err.(type) {
		case ErrTxnViolatesSoftConstraint:
			e := err.(ErrTxnViolatesSoftConstraint)
			softErr = &e
			isValid = 0
		case ErrTxnViolatesHardConstraint:
			return false, nil, err
		default:
			return false, nil, err
		}
	}

	// Update if we already have this txn
	h := t.Hash()
	known := false
	utp.txns.update(h, func(tx *UnconfirmedTxn) {
		known = true
		now := utc.Now().UnixNano()
		tx.Received = now
		tx.Checked = now
		tx.IsValid = isValid
	})

	if known {
		return true, softErr, nil
	}

	utx := utp.createUnconfirmedTxn(t)
	utx.IsValid = isValid

	if err := bc.UpdateDB(func(tx *bolt.Tx) error {
		// add txn to index
		if err := utp.txns.putWithTx(tx, &utx); err != nil {
			return err
		}

		// update unconfirmed unspent
		head, err := bc.Head()
		if err != nil {
			return err
		}

		return utp.unspent.putWithTx(tx, h, coin.CreateUnspents(head.Head, t))
	}); err != nil {
		return false, nil, err
	}

	return false, softErr, nil
}

// RawTxns returns underlying coin.Transactions
func (utp *UnconfirmedTxnPool) RawTxns() coin.Transactions {
	utxns, err := utp.txns.getAll()
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
	utp.txns.delete(txHash)
	utp.unspent.delete(txHash)
}

// Removes multiple txns at once. Slightly more efficient than a series of
// single RemoveTxns.  Hashes is an array of Transaction hashes.
func (utp *UnconfirmedTxnPool) removeTxns(hashes []cipher.SHA256) error {
	for i := range hashes {
		if err := utp.txns.delete(hashes[i]); err != nil {
			return err
		}
		if err := utp.unspent.delete(hashes[i]); err != nil {
			return err
		}
	}

	return nil
}

func (utp *UnconfirmedTxnPool) removeTxnsWithTx(tx *bolt.Tx, hashes []cipher.SHA256) {
	for i := range hashes {
		utp.txns.deleteWithTx(tx, hashes[i])
		utp.unspent.deleteWithTx(tx, hashes[i])
	}
}

// RemoveTransactions removes confirmed txns from the pool
func (utp *UnconfirmedTxnPool) RemoveTransactions(txns []cipher.SHA256) error {
	return utp.removeTxns(txns)
}

// RemoveTransactionsWithTx remove transactions with bolt.Tx
func (utp *UnconfirmedTxnPool) RemoveTransactionsWithTx(tx *bolt.Tx, txns []cipher.SHA256) {
	utp.removeTxnsWithTx(tx, txns)
}

// Refresh checks all unconfirmed txns against the blockchain.
// If the transaction becomes invalid it is marked invalid.
// If the transaction becomes valid it is marked valid and is returned to the caller.
func (utp *UnconfirmedTxnPool) Refresh(bc Blockchainer, maxBlockSize int) ([]cipher.SHA256, error) {
	now := utc.Now()

	var nowValid []cipher.SHA256

	if err := utp.txns.rangeUpdate(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		tx.Checked = now.UnixNano()

		err := bc.VerifySingleTxnAllConstraints(tx.Txn, maxBlockSize)

		switch err.(type) {
		case ErrTxnViolatesSoftConstraint, ErrTxnViolatesHardConstraint:
			tx.IsValid = 0
		case nil:
			if tx.IsValid == 0 {
				nowValid = append(nowValid, tx.Hash())
			}
			tx.IsValid = 1
		default:
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return nowValid, nil
}

// RemoveInvalid checks all unconfirmed txns against the blockchain.
// If a transaction violates hard constraints it is removed from the pool.
// The transactions that were removed are returned.
func (utp *UnconfirmedTxnPool) RemoveInvalid(bc Blockchainer) ([]cipher.SHA256, error) {
	var removeTxs []cipher.SHA256

	if err := utp.txns.forEach(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		err := bc.VerifySingleTxnHardConstraints(tx.Txn)

		switch err.(type) {
		case ErrTxnViolatesHardConstraint:
			removeTxs = append(removeTxs, tx.Hash())
		default:
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err := utp.RemoveTransactions(removeTxs); err != nil {
		return nil, err
	}

	return removeTxs, nil
}

// FilterKnown returns txn hashes with known ones removed
func (utp *UnconfirmedTxnPool) FilterKnown(txns []cipher.SHA256) []cipher.SHA256 {
	var unknown []cipher.SHA256
	for _, h := range txns {
		if !utp.txns.isExist(h) {
			unknown = append(unknown, h)
		}
	}
	return unknown
}

// GetKnown returns all known coin.Transactions from the pool, given hashes to select
func (utp *UnconfirmedTxnPool) GetKnown(txns []cipher.SHA256) coin.Transactions {
	var known coin.Transactions
	for _, h := range txns {
		if tx, ok := utp.txns.get(h); ok {
			known = append(known, tx.Txn)
		}
	}
	return known
}

// RecvOfAddresses returns unconfirmed receiving uxouts of addresses
func (utp *UnconfirmedTxnPool) RecvOfAddresses(bh coin.BlockHeader,
	addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}
	auxs := make(coin.AddressUxOuts, len(addrs))
	if err := utp.txns.forEach(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		for i, o := range tx.Txn.Out {
			if _, ok := addrm[o.Address]; ok {
				uxout, err := coin.CreateUnspent(bh, tx.Txn, i)
				if err != nil {
					return err
				}

				auxs[o.Address] = append(auxs[o.Address], uxout)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return auxs, nil
}

// UnspentGetFunc callback function for querying unspent output of given hash
type UnspentGetFunc func(hash cipher.SHA256) (coin.UxOut, bool)

// SpendsOfAddresses returns all unconfirmed coin.UxOut spends of addresses
// Looks at all inputs for unconfirmed txns, gets their source UxOut from the
// blockchain's unspent pool, and returns as coin.AddressUxOuts
func (utp *UnconfirmedTxnPool) SpendsOfAddresses(addrs []cipher.Address,
	unspent blockdb.UnspentGetter) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}

	auxs := make(coin.AddressUxOuts, len(addrs))
	if err := utp.txns.forEach(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		for _, h := range tx.Txn.In {
			ux, ok := unspent.Get(h)
			if !ok {
				// unconfirm transaction's IN is not in the unspent pool, this should not happen
				return fmt.Errorf("unconfirmed transaction's IN: %s is not in unspent pool", h.Hex())
			}

			if _, ok := addrm[ux.Body.Address]; ok {
				auxs[ux.Body.Address] = append(auxs[ux.Body.Address], ux)
			}
		}
		return nil
	}); err != nil {
		return coin.AddressUxOuts{}, fmt.Errorf("get unconfirmed spend error:%v", err)
	}
	return auxs, nil
}

// GetSpendingOutputs returns all spending outputs in unconfirmed tx pool.
func (utp *UnconfirmedTxnPool) GetSpendingOutputs(bcUnspent blockdb.UnspentPool) (coin.UxArray, error) {
	outs := coin.UxArray{}
	err := utp.txns.forEach(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		uxs, err := bcUnspent.GetArray(tx.Txn.In)
		if err != nil {
			return err
		}

		outs = append(outs, uxs...)
		return nil
	})

	if err != nil {
		return coin.UxArray{}, fmt.Errorf("get unconfirmed spending outputs failed: %v", err)
	}

	return outs, nil
}

// GetIncomingOutputs returns all predicted incoming outputs.
func (utp *UnconfirmedTxnPool) GetIncomingOutputs(bh coin.BlockHeader) coin.UxArray {
	outs := coin.UxArray{}
	utp.txns.forEach(func(_ cipher.SHA256, tx *UnconfirmedTxn) error {
		uxOuts := coin.CreateUnspents(bh, tx.Txn)
		outs = append(outs, uxOuts...)
		return nil
	})
	return outs
}

// Get returns the unconfirmed transaction of given tx hash.
func (utp *UnconfirmedTxnPool) Get(key cipher.SHA256) (*UnconfirmedTxn, bool) {
	return utp.txns.get(key)
}

// GetTxns returns all transactions that can pass the filter
func (utp *UnconfirmedTxnPool) GetTxns(filter func(tx UnconfirmedTxn) bool) (txns []UnconfirmedTxn) {
	if err := utp.txns.forEach(func(hash cipher.SHA256, tx *UnconfirmedTxn) error {
		if filter(*tx) {
			txns = append(txns, *tx)
		}
		return nil
	}); err != nil {
		logger.Debugf("GetTxns error:%v", err)
	}
	return
}

// GetTxHashes returns transaction hashes that can pass the filter
func (utp *UnconfirmedTxnPool) GetTxHashes(filter func(tx UnconfirmedTxn) bool) (hashes []cipher.SHA256) {
	if err := utp.txns.forEach(func(hash cipher.SHA256, tx *UnconfirmedTxn) error {
		if filter(*tx) {
			hashes = append(hashes, hash)
		}
		return nil
	}); err != nil {
		logger.Debugf("GetTxHashes error:%v", err)
	}
	return
}

// ForEach iterate the pool with given callback function,
func (utp *UnconfirmedTxnPool) ForEach(f func(cipher.SHA256, *UnconfirmedTxn) error) error {
	return utp.txns.forEach(f)
}

// GetUnspentsOfAddr returns unspent outputs of given address in unspent tx pool
func (utp *UnconfirmedTxnPool) GetUnspentsOfAddr(addr cipher.Address) coin.UxArray {
	return utp.unspent.getByAddr(addr)
}

// IsValid can be used as filter function
func IsValid(tx UnconfirmedTxn) bool {
	return tx.IsValid == 1
}

// All use as return all filter
func All(tx UnconfirmedTxn) bool {
	return true
}

// Len returns the number of unconfirmed transactions
func (utp *UnconfirmedTxnPool) Len() int {
	return utp.txns.len()
}

func nanoToTime(n int64) time.Time {
	zeroTime := time.Time{}
	if n == zeroTime.UnixNano() {
		// maximum time
		return zeroTime
	}
	return time.Unix(n/int64(time.Second), n%int64(time.Second))
}
