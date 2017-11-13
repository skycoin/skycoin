package visor

import (
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// BurnFactor half of coinhours must be burnt
	BurnFactor uint64 = 2

	unconfirmedTxnsBkt     = []byte("unconfirmed_txns")
	unconfirmedUnspentsBkt = []byte("unconfirmed_unspents")

	errUpdateObjectDoesNotExist = errors.New("object does not exist in bucket")
)

// VerifyTransactionFee performs additional transaction verification at the unconfirmed pool level.
// This checks tunable parameters that should prevent the transaction from
// entering the blockchain, but cannot be done at the blockchain level because
// they may be changed.
func VerifyTransactionFee(t *coin.Transaction, fee uint64) error {
	// Calculate total number of coinhours
	var total = t.OutputHours() + fee
	// Make sure at least half (BurnFactor=2) the coin hours are destroyed
	if fee < total/BurnFactor {
		return errors.New("Transaction coinhour fee minimum not met")
	}
	return nil
}

// TransactionFee calculates the current transaction fee in coinhours of a Transaction
func TransactionFee(t *coin.Transaction, headTime uint64, inUxs coin.UxArray) (uint64, error) {
	// Compute input hours
	inHours := uint64(0)
	for _, ux := range inUxs {
		inHours += ux.CoinHours(headTime)
	}

	// Compute output hours
	outHours := uint64(0)
	for i := range t.Out {
		outHours += t.Out[i].Hours
	}

	if inHours < outHours {
		return 0, errors.New("Insufficient coinhours for transaction outputs")
	}

	return inHours - outHours, nil
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
type unconfirmedTxns struct{}

func (utb *unconfirmedTxns) get(tx *bolt.Tx, hash cipher.SHA256) (*UnconfirmedTxn, error) {
	var txn UnconfirmedTxn

	if ok, err := dbutil.GetBucketObjectDecoded(tx, unconfirmedTxnsBkt, []byte(hash.Hex()), &txn); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	if ok, err := dbutil.GetBucketObjectDecoded(tx, unconfirmedTxnsBkt, []byte(hash.Hex()), &txn); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return &txn, nil
}

func (utb *unconfirmedTxns) put(tx *bolt.Tx, v *UnconfirmedTxn) error {
	return dbutil.PutBucketValue(tx, unconfirmedTxnsBkt, []byte(v.Hash().Hex()), encoder.Serialize(v))
}

func (utb *unconfirmedTxns) update(tx *bolt.Tx, hash cipher.SHA256, f func(v *UnconfirmedTxn) error) error {
	txn, err := utb.get(tx, hash)
	if err != nil {
		return err
	}

	if txn == nil {
		return errUpdateObjectDoesNotExist
	}

	if err := f(txn); err != nil {
		return err
	}

	return utb.put(tx, txn)
}

func (utb *unconfirmedTxns) delete(tx *bolt.Tx, hash cipher.SHA256) error {
	return dbutil.Delete(tx, unconfirmedTxnsBkt, []byte(hash.Hex()))
}

func (utb *unconfirmedTxns) getAll(tx *bolt.Tx) ([]UnconfirmedTxn, error) {
	var txns []UnconfirmedTxn

	if err := dbutil.ForEach(tx, unconfirmedTxnsBkt, func(_, v []byte) error {
		var txn UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &txn); err != nil {
			return err
		}

		txns = append(txns, txn)
		return nil
	}); err != nil {
		return nil, err
	}

	return txns, nil
}

func (utb *unconfirmedTxns) rangeUpdate(tx *bolt.Tx, f func(UnconfirmedTxn) (UnconfirmedTxn, error)) error {
	// The UnconfirmedTxns must be pulled from the DB first, it is not
	// safe to iterate the db with ForEach while modifying the DB
	txns, err := utb.getAll(tx)
	if err != nil {
		return err
	}

	for _, txn := range txns {
		modifiedTxn, err := f(txn)
		if err != nil {
			return err
		}

		if err := utb.put(tx, &modifiedTxn); err != nil {
			return err
		}
	}

	return nil
}

func (utb *unconfirmedTxns) hasKey(tx *bolt.Tx, hash cipher.SHA256) (bool, error) {
	return dbutil.BucketHasKey(tx, unconfirmedTxnsBkt, []byte(hash.Hex()))
}

func (utb *unconfirmedTxns) forEach(tx *bolt.Tx, f func(hash cipher.SHA256, tx UnconfirmedTxn) error) error {
	return dbutil.ForEach(tx, unconfirmedTxnsBkt, func(k, v []byte) error {
		hash, err := cipher.SHA256FromHex(string(k))
		if err != nil {
			return err
		}

		var txn UnconfirmedTxn
		if err := encoder.DeserializeRaw(v, &txn); err != nil {
			return err
		}

		return f(hash, txn)
	})
}

func (utb *unconfirmedTxns) length(tx *bolt.Tx) (uint64, error) {
	return dbutil.Len(tx, unconfirmedTxnsBkt)
}

type txUnspents struct{}

func (txus *txUnspents) put(tx *bolt.Tx, hash cipher.SHA256, uxs coin.UxArray) error {
	return dbutil.PutBucketValue(tx, unconfirmedUnspentsBkt, []byte(hash.Hex()), encoder.Serialize(uxs))
}

func (txus *txUnspents) get(tx *bolt.Tx, hash cipher.SHA256) (coin.UxArray, error) {
	var uxs coin.UxArray

	if ok, err := dbutil.GetBucketObjectDecoded(tx, unconfirmedUnspentsBkt, []byte(hash.Hex()), &uxs); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return uxs, nil
}

func (txus *txUnspents) length(tx *bolt.Tx) (uint64, error) {
	return dbutil.Len(tx, unconfirmedUnspentsBkt)
}

func (txus *txUnspents) delete(tx *bolt.Tx, hash cipher.SHA256) error {
	return dbutil.Delete(tx, unconfirmedUnspentsBkt, []byte(hash.Hex()))
}

func (txus *txUnspents) getByAddr(tx *bolt.Tx, a cipher.Address) (coin.UxArray, error) {
	var uxo coin.UxArray

	if err := dbutil.ForEach(tx, unconfirmedUnspentsBkt, func(_, v []byte) error {
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
	}); err != nil {
		return nil, err
	}

	return uxo, nil
}

func (txus *txUnspents) forEach(tx *bolt.Tx, f func(cipher.SHA256, coin.UxArray) error) error {
	return dbutil.ForEach(tx, unconfirmedUnspentsBkt, func(k, v []byte) error {
		hash, err := cipher.SHA256FromHex(string(k))
		if err != nil {
			return err
		}

		var uxa coin.UxArray
		if err := encoder.DeserializeRaw(v, &uxa); err != nil {
			return err
		}

		return f(hash, uxa)
	})
}

// UnconfirmedTxnPool manages unconfirmed transactions
type UnconfirmedTxnPool struct {
	db   *dbutil.DB
	txns *unconfirmedTxns
	// Predicted unspents, assuming txns are valid.  Needed to predict
	// our future balance and avoid double spending our own coins
	// Maps from Transaction.Hash() to UxArray.
	unspent *txUnspents
}

// NewUnconfirmedTxnPool creates an UnconfirmedTxnPool instance
func NewUnconfirmedTxnPool(db *dbutil.DB) (*UnconfirmedTxnPool, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		return dbutil.CreateBuckets(tx, [][]byte{
			unconfirmedTxnsBkt,
			unconfirmedUnspentsBkt,
		})
	}); err != nil {
		return nil, err
	}

	return &UnconfirmedTxnPool{
		db:      db,
		txns:    &unconfirmedTxns{},
		unspent: &txUnspents{},
	}, nil
}

// SetTxnsAnnounced updates announced time of specific tx
func (utp *UnconfirmedTxnPool) SetTxnsAnnounced(tx *bolt.Tx, txns []cipher.SHA256, t time.Time) error {
	for _, h := range txns {
		if err := utp.setAnnounced(tx, h, t); err != nil {
			return err
		}
	}

	return nil
}

func (utp *UnconfirmedTxnPool) setAnnounced(tx *bolt.Tx, h cipher.SHA256, t time.Time) error {
	return utp.txns.update(tx, h, func(tx *UnconfirmedTxn) error {
		tx.Announced = t.UnixNano()
		return nil
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
func (utp *UnconfirmedTxnPool) InjectTransaction(tx *bolt.Tx, bc *Blockchain, t coin.Transaction) (bool, error) {
	head, err := bc.Head(tx)
	if err != nil {
		return false, err
	}

	fee, err := bc.TransactionFee(head.Time())(&t)
	if err != nil {
		return false, err
	}

	if err := VerifyTransactionFee(&t, fee); err != nil {
		return false, err
	}

	hash := t.Hash()

	if err := bc.VerifyTransaction(head, t); err != nil {
		return false, err
	}

	known, err := utp.txns.hasKey(tx, hash)
	if err != nil {
		return false, err
	}

	// Update if we already have this txn
	if known {
		if err := utp.txns.update(tx, hash, func(tx *UnconfirmedTxn) error {
			now := utc.Now().UnixNano()
			tx.Received = now
			tx.Checked = now
			tx.IsValid = 1
			return nil
		}); err != nil {
			return false, err
		}

		return true, nil
	}

	utx := utp.createUnconfirmedTxn(t)
	// add txn to index
	if err := utp.txns.put(tx, &utx); err != nil {
		return false, err
	}

	// update unconfirmed unspent
	head, err = bc.Head(tx)
	if err != nil {
		return false, err
	}

	if err := utp.unspent.put(tx, hash, coin.CreateUnspents(head.Head, t)); err != nil {
		return false, err
	}

	return false, nil
}

// RawTxns returns underlying coin.Transactions
func (utp *UnconfirmedTxnPool) RawTxns(tx *bolt.Tx) (coin.Transactions, error) {
	utxns, err := utp.txns.getAll(tx)
	if err != nil {
		return nil, err
	}

	txns := make(coin.Transactions, len(utxns))
	for i := range utxns {
		txns[i] = utxns[i].Txn
	}
	return txns, nil
}

// Remove a single txn by hash
func (utp *UnconfirmedTxnPool) removeTxn(tx *bolt.Tx, txHash cipher.SHA256) error {
	if err := utp.txns.delete(tx, txHash); err != nil {
		return err
	}

	return utp.unspent.delete(tx, txHash)
}

// RemoveTransactions remove transactions with bolt.Tx
func (utp *UnconfirmedTxnPool) RemoveTransactions(tx *bolt.Tx, txHashes []cipher.SHA256) error {
	for i := range txHashes {
		if err := utp.removeTxn(tx, txHashes[i]); err != nil {
			return err
		}
	}

	return nil
}

// Refresh checks all unconfirmed txns against the blockchain.
// verify the transaction and returns all those txns that turn to valid.
func (utp *UnconfirmedTxnPool) Refresh(tx *bolt.Tx, bc *Blockchain) ([]cipher.SHA256, error) {
	logger.Debug("UnconfirmedTxnPool.RefreshUnconfirmed")
	var hashes []cipher.SHA256
	now := utc.Now().UnixNano()

	head, err := bc.Head(tx)
	if err != nil {
		return nil, err
	}

	if err := utp.txns.rangeUpdate(tx, func(txn UnconfirmedTxn) (UnconfirmedTxn, error) {
		txn.Checked = now
		if txn.IsValid == 0 {
			if bc.VerifyTransaction(head, txn.Txn) == nil {
				txn.IsValid = 1
				hashes = append(hashes, txn.Hash())
			}
		}
		return txn, nil
	}); err != nil {
		return nil, err
	}

	return hashes, nil
}

// GetUnknown returns txn hashes with known ones removed
func (utp *UnconfirmedTxnPool) GetUnknown(tx *bolt.Tx, txns []cipher.SHA256) ([]cipher.SHA256, error) {
	var unknown []cipher.SHA256

	for _, h := range txns {
		if hasKey, err := utp.txns.hasKey(tx, h); err != nil {
			return nil, err
		} else if !hasKey {
			unknown = append(unknown, h)
		}
	}

	return unknown, nil
}

// GetKnown returns all known coin.Transactions from the pool, given hashes to select
func (utp *UnconfirmedTxnPool) GetKnown(tx *bolt.Tx, txns []cipher.SHA256) (coin.Transactions, error) {
	var known coin.Transactions

	for _, h := range txns {
		if tx, err := utp.txns.get(tx, h); err != nil {
			return nil, err
		} else if tx != nil {
			known = append(known, tx.Txn)
		}
	}

	return known, nil
}

// RecvOfAddresses returns unconfirmed receiving uxouts of addresses
func (utp *UnconfirmedTxnPool) RecvOfAddresses(tx *bolt.Tx, bh coin.BlockHeader, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}

	auxs := make(coin.AddressUxOuts, len(addrs))
	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, tx UnconfirmedTxn) error {
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
func (utp *UnconfirmedTxnPool) SpendsOfAddresses(tx *bolt.Tx, addrs []cipher.Address, unspent blockdb.UnspentGetter) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}

	auxs := make(coin.AddressUxOuts, len(addrs))

	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTxn) error {
		for _, h := range txn.Txn.In {
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
		return nil, err
	}

	return auxs, nil
}

// GetSpendingOutputs returns all spending outputs in unconfirmed tx pool.
func (utp *UnconfirmedTxnPool) GetSpendingOutputs(tx *bolt.Tx, bcUnspent blockdb.UnspentPool) (coin.UxArray, error) {
	var outs coin.UxArray

	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTxn) error {
		uxs, err := bcUnspent.GetArray(txn.Txn.In)
		if err != nil {
			return err
		}

		outs = append(outs, uxs...)
		return nil
	}); err != nil {
		return nil, err
	}

	return outs, nil
}

// GetIncomingOutputs returns all predicted incoming outputs.
func (utp *UnconfirmedTxnPool) GetIncomingOutputs(tx *bolt.Tx, bh coin.BlockHeader) (coin.UxArray, error) {
	var outs coin.UxArray

	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTxn) error {
		uxOuts := coin.CreateUnspents(bh, txn.Txn)
		outs = append(outs, uxOuts...)
		return nil
	}); err != nil {
		return nil, err
	}

	return outs, nil
}

// Get returns the unconfirmed transaction of given tx hash.
func (utp *UnconfirmedTxnPool) Get(tx *bolt.Tx, hash cipher.SHA256) (*UnconfirmedTxn, error) {
	return utp.txns.get(tx, hash)
}

// GetTxns returns all transactions that can pass the filter
func (utp *UnconfirmedTxnPool) GetTxns(tx *bolt.Tx, filter func(UnconfirmedTxn) bool) ([]UnconfirmedTxn, error) {
	var txns []UnconfirmedTxn

	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTxn) error {
		if filter(txn) {
			txns = append(txns, txn)
		}
		return nil
	}); err != nil {
		logger.Debug("GetTxns error:%v", err)
		return nil, err
	}

	return txns, nil
}

// GetTxHashes returns transaction hashes that can pass the filter
func (utp *UnconfirmedTxnPool) GetTxHashes(tx *bolt.Tx, filter func(UnconfirmedTxn) bool) ([]cipher.SHA256, error) {
	var hashes []cipher.SHA256

	if err := utp.txns.forEach(tx, func(hash cipher.SHA256, txn UnconfirmedTxn) error {
		if filter(txn) {
			hashes = append(hashes, hash)
		}
		return nil
	}); err != nil {
		logger.Debug("GetTxHashes error:%v", err)
		return nil, err
	}

	return hashes, nil
}

// ForEach iterate the pool with given callback function
func (utp *UnconfirmedTxnPool) ForEach(tx *bolt.Tx, f func(cipher.SHA256, UnconfirmedTxn) error) error {
	return utp.txns.forEach(tx, f)
}

// GetUnspentsOfAddr returns unspent outputs of given address in unspent tx pool
func (utp *UnconfirmedTxnPool) GetUnspentsOfAddr(tx *bolt.Tx, addr cipher.Address) (coin.UxArray, error) {
	return utp.unspent.getByAddr(tx, addr)
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
func (utp *UnconfirmedTxnPool) Len(tx *bolt.Tx) (uint64, error) {
	return utp.txns.length(tx)
}

func nanoToTime(n int64) time.Time {
	zeroTime := time.Time{}
	if n == zeroTime.UnixNano() {
		// maximum time
		return zeroTime
	}
	return time.Unix(n/int64(time.Second), n%int64(time.Second))
}
