package visor

import (
	"errors"
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// UnconfirmedTxnsBkt holds unconfirmed transactions
	UnconfirmedTxnsBkt = []byte("unconfirmed_txns")
	// UnconfirmedUnspentsBkt holds unconfirmed unspent outputs
	UnconfirmedUnspentsBkt = []byte("unconfirmed_unspents")

	errUpdateObjectDoesNotExist = errors.New("object does not exist in bucket")
)

//go:generate skyencoder -unexported -struct UnconfirmedTransaction
//go:generate skyencoder -unexported -struct UxArray

// UxArray wraps coin.UxArray
type UxArray struct {
	UxArray coin.UxArray
}

// unconfirmed transactions bucket
type unconfirmedTxns struct{}

func (utb *unconfirmedTxns) get(tx *dbutil.Tx, hash cipher.SHA256) (*UnconfirmedTransaction, error) {
	var txn UnconfirmedTransaction

	v, err := dbutil.GetBucketValueNoCopy(tx, UnconfirmedTxnsBkt, []byte(hash.Hex()))
	if err != nil {
		return nil, err
	} else if v == nil {
		return nil, nil
	}

	if err := decodeUnconfirmedTransactionExact(v, &txn); err != nil {
		return nil, err
	}

	txnHash := txn.Transaction.Hash()
	if hash != txnHash {
		return nil, fmt.Errorf("DB key %s does not match block hash header %s", hash, txnHash)
	}

	return &txn, nil
}

func (utb *unconfirmedTxns) put(tx *dbutil.Tx, v *UnconfirmedTransaction) error {
	h := v.Transaction.Hash()
	buf, err := encodeUnconfirmedTransaction(v)
	if err != nil {
		return err
	}

	return dbutil.PutBucketValue(tx, UnconfirmedTxnsBkt, []byte(h.Hex()), buf)
}

func (utb *unconfirmedTxns) update(tx *dbutil.Tx, hash cipher.SHA256, f func(v *UnconfirmedTransaction) error) error {
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

func (utb *unconfirmedTxns) delete(tx *dbutil.Tx, hash cipher.SHA256) error {
	return dbutil.Delete(tx, UnconfirmedTxnsBkt, []byte(hash.Hex()))
}

func (utb *unconfirmedTxns) getAll(tx *dbutil.Tx) ([]UnconfirmedTransaction, error) {
	var txns []UnconfirmedTransaction

	if err := dbutil.ForEach(tx, UnconfirmedTxnsBkt, func(_, v []byte) error {
		var txn UnconfirmedTransaction
		if err := decodeUnconfirmedTransactionExact(v, &txn); err != nil {
			return err
		}

		txns = append(txns, txn)
		return nil
	}); err != nil {
		return nil, err
	}

	return txns, nil
}

func (utb *unconfirmedTxns) hasKey(tx *dbutil.Tx, hash cipher.SHA256) (bool, error) {
	return dbutil.BucketHasKey(tx, UnconfirmedTxnsBkt, []byte(hash.Hex()))
}

func (utb *unconfirmedTxns) forEach(tx *dbutil.Tx, f func(hash cipher.SHA256, tx UnconfirmedTransaction) error) error {
	return dbutil.ForEach(tx, UnconfirmedTxnsBkt, func(k, v []byte) error {
		hash, err := cipher.SHA256FromHex(string(k))
		if err != nil {
			return err
		}

		var txn UnconfirmedTransaction
		if err := decodeUnconfirmedTransactionExact(v, &txn); err != nil {
			return err
		}

		return f(hash, txn)
	})
}

func (utb *unconfirmedTxns) len(tx *dbutil.Tx) (uint64, error) {
	return dbutil.Len(tx, UnconfirmedTxnsBkt)
}

type txnUnspents struct{}

func (txus *txnUnspents) put(tx *dbutil.Tx, hash cipher.SHA256, uxs coin.UxArray) error {
	buf, err := encodeUxArray(&UxArray{
		UxArray: uxs,
	})
	if err != nil {
		return err
	}

	return dbutil.PutBucketValue(tx, UnconfirmedUnspentsBkt, []byte(hash.Hex()), buf)
}

func (txus *txnUnspents) delete(tx *dbutil.Tx, hash cipher.SHA256) error {
	return dbutil.Delete(tx, UnconfirmedUnspentsBkt, []byte(hash.Hex()))
}

func (txus *txnUnspents) getByAddr(tx *dbutil.Tx, a cipher.Address) (coin.UxArray, error) {
	var uxo coin.UxArray

	if err := dbutil.ForEach(tx, UnconfirmedUnspentsBkt, func(_, v []byte) error {
		var uxa UxArray
		if err := decodeUxArrayExact(v, &uxa); err != nil {
			return err
		}

		for i := range uxa.UxArray {
			if uxa.UxArray[i].Body.Address == a {
				uxo = append(uxo, uxa.UxArray[i])
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return uxo, nil
}

// UnconfirmedTransactionPool manages unconfirmed transactions
type UnconfirmedTransactionPool struct {
	db   *dbutil.DB
	txns *unconfirmedTxns
	// Predicted unspents, assuming txns are valid.  Needed to predict
	// our future balance and avoid double spending our own coins
	// Maps from Transaction.Hash() to UxArray.
	unspent *txnUnspents
}

// NewUnconfirmedTransactionPool creates an UnconfirmedTransactionPool instance
func NewUnconfirmedTransactionPool(db *dbutil.DB) (*UnconfirmedTransactionPool, error) {
	if err := db.View("Check unconfirmed txn pool size", func(tx *dbutil.Tx) error {
		n, err := dbutil.Len(tx, UnconfirmedTxnsBkt)
		if err != nil {
			return err
		}

		logger.Infof("Unconfirmed transaction pool size: %d", n)
		return nil
	}); err != nil {
		return nil, err
	}

	return &UnconfirmedTransactionPool{
		db:      db,
		txns:    &unconfirmedTxns{},
		unspent: &txnUnspents{},
	}, nil
}

// SetTransactionsAnnounced updates announced time of specific tx
func (utp *UnconfirmedTransactionPool) SetTransactionsAnnounced(tx *dbutil.Tx, hashes map[cipher.SHA256]int64) error {
	var txns []*UnconfirmedTransaction
	for h, t := range hashes {
		txn, err := utp.txns.get(tx, h)
		if err != nil {
			return err
		}

		if txn == nil {
			logger.Warningf("UnconfirmedTransactionPool.SetTransactionsAnnounced: UnconfirmedTransaction %s not found in DB", h.Hex())
			continue
		}

		if t > txn.Announced {
			txn.Announced = t
			txns = append(txns, txn)
		}
	}

	for _, txn := range txns {
		if err := utp.txns.put(tx, txn); err != nil {
			return err
		}
	}

	return nil
}

// InjectTransaction adds a coin.Transaction to the pool, or updates an existing one's timestamps
// Returns an error if txn is invalid, and whether the transaction already
// existed in the pool.
// If the transaction violates hard constraints, it is rejected.
// Soft constraints violations mark a txn as invalid, but the txn is inserted. The soft violation is returned.
func (utp *UnconfirmedTransactionPool) InjectTransaction(tx *dbutil.Tx, bc Blockchainer, txn coin.Transaction, distParams params.Distribution, verifyParams params.VerifyTxn) (bool, *ErrTxnViolatesSoftConstraint, error) {
	var isValid int8 = 1
	var softErr *ErrTxnViolatesSoftConstraint
	if _, _, err := bc.VerifySingleTxnSoftHardConstraints(tx, txn, distParams, verifyParams, TxnSigned); err != nil {
		logger.Warningf("bc.VerifySingleTxnSoftHardConstraints failed for txn %s: %v", txn.Hash().Hex(), err)
		switch e := err.(type) {
		case ErrTxnViolatesSoftConstraint:
			softErr = &e
			isValid = 0
		case ErrTxnViolatesHardConstraint:
			return false, nil, err
		default:
			return false, nil, err
		}
	}

	hash := txn.Hash()
	known, err := utp.txns.hasKey(tx, hash)
	if err != nil {
		logger.Errorf("InjectTransaction check txn exists failed: %v", err)
		return false, nil, err
	}

	// Update if we already have this txn
	if known {
		if err := utp.txns.update(tx, hash, func(utxn *UnconfirmedTransaction) error {
			now := time.Now().UTC().UnixNano()
			utxn.Received = now
			utxn.Checked = now
			utxn.IsValid = isValid
			return nil
		}); err != nil {
			logger.Errorf("InjectTransaction update known txn failed: %v", err)
			return false, nil, err
		}

		return true, softErr, nil
	}

	utx := NewUnconfirmedTransaction(txn)
	utx.IsValid = isValid

	// add txn to index
	if err := utp.txns.put(tx, &utx); err != nil {
		logger.Errorf("InjectTransaction put new unconfirmed txn failed: %v", err)
		return false, nil, err
	}

	head, err := bc.Head(tx)
	if err != nil {
		logger.Errorf("InjectTransaction bc.Head() failed: %v", err)
		return false, nil, err
	}

	// update unconfirmed unspent
	createdUnspents := coin.CreateUnspents(head.Head, txn)
	if err := utp.unspent.put(tx, hash, createdUnspents); err != nil {
		logger.Errorf("InjectTransaction put new unspent outputs: %v", err)
		return false, nil, err
	}

	return false, softErr, nil
}

// AllRawTransactions returns underlying coin.Transactions
func (utp *UnconfirmedTransactionPool) AllRawTransactions(tx *dbutil.Tx) (coin.Transactions, error) {
	utxns, err := utp.txns.getAll(tx)
	if err != nil {
		return nil, err
	}

	txns := make(coin.Transactions, len(utxns))
	for i := range utxns {
		txns[i] = utxns[i].Transaction
	}
	return txns, nil
}

// Remove a single txn by hash
func (utp *UnconfirmedTransactionPool) removeTransaction(tx *dbutil.Tx, txHash cipher.SHA256) error {
	if err := utp.txns.delete(tx, txHash); err != nil {
		return err
	}

	return utp.unspent.delete(tx, txHash)
}

// RemoveTransactions remove transactions with dbutil.Tx
func (utp *UnconfirmedTransactionPool) RemoveTransactions(tx *dbutil.Tx, txHashes []cipher.SHA256) error {
	for i := range txHashes {
		if err := utp.removeTransaction(tx, txHashes[i]); err != nil {
			return err
		}
	}

	return nil
}

// Refresh checks all unconfirmed txns against the blockchain.
// If the transaction becomes invalid it is marked invalid.
// If the transaction becomes valid it is marked valid and is returned to the caller.
func (utp *UnconfirmedTransactionPool) Refresh(tx *dbutil.Tx, bc Blockchainer, distParams params.Distribution, verifyParams params.VerifyTxn) ([]cipher.SHA256, error) {
	utxns, err := utp.txns.getAll(tx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	var nowValid []cipher.SHA256

	for _, utxn := range utxns {
		utxn.Checked = now.UnixNano()

		_, _, err := bc.VerifySingleTxnSoftHardConstraints(tx, utxn.Transaction, distParams, verifyParams, TxnSigned)

		switch err.(type) {
		case ErrTxnViolatesSoftConstraint, ErrTxnViolatesHardConstraint:
			utxn.IsValid = 0
		case nil:
			if utxn.IsValid == 0 {
				nowValid = append(nowValid, utxn.Transaction.Hash())
			}
			utxn.IsValid = 1
		default:
			return nil, err
		}

		if err := utp.txns.put(tx, &utxn); err != nil {
			return nil, err
		}
	}

	return nowValid, nil
}

// RemoveInvalid checks all unconfirmed txns against the blockchain.
// If a transaction violates hard constraints it is removed from the pool.
// The transactions that were removed are returned.
func (utp *UnconfirmedTransactionPool) RemoveInvalid(tx *dbutil.Tx, bc Blockchainer) ([]cipher.SHA256, error) {
	var removeUtxns []cipher.SHA256

	utxns, err := utp.txns.getAll(tx)
	if err != nil {
		return nil, err
	}

	for _, utxn := range utxns {
		err := bc.VerifySingleTxnHardConstraints(tx, utxn.Transaction, TxnSigned)
		if err != nil {
			switch err.(type) {
			case ErrTxnViolatesHardConstraint:
				removeUtxns = append(removeUtxns, utxn.Transaction.Hash())
			default:
				return nil, err
			}
		}
	}

	if err := utp.RemoveTransactions(tx, removeUtxns); err != nil {
		return nil, err
	}

	return removeUtxns, nil
}

// FilterKnown returns txn hashes with known ones removed
func (utp *UnconfirmedTransactionPool) FilterKnown(tx *dbutil.Tx, txns []cipher.SHA256) ([]cipher.SHA256, error) {
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

// GetKnown returns all known transactions from the pool, given hashes to select
func (utp *UnconfirmedTransactionPool) GetKnown(tx *dbutil.Tx, txns []cipher.SHA256) (coin.Transactions, error) {
	var known coin.Transactions

	for _, h := range txns {
		if tx, err := utp.txns.get(tx, h); err != nil {
			return nil, err
		} else if tx != nil {
			known = append(known, tx.Transaction)
		}
	}

	return known, nil
}

// RecvOfAddresses returns unconfirmed receiving uxouts of addresses
func (utp *UnconfirmedTransactionPool) RecvOfAddresses(tx *dbutil.Tx, bh coin.BlockHeader, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}

	auxs := make(coin.AddressUxOuts, len(addrs))
	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTransaction) error {
		for i, o := range txn.Transaction.Out {
			if _, ok := addrm[o.Address]; ok {
				uxout, err := coin.CreateUnspent(bh, txn.Transaction, i)
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

// txnOutputsForAddrs returns unspent outputs assigned to addresses in addrs, created by a set of transactions
func txnOutputsForAddrs(bh coin.BlockHeader, addrs []cipher.Address, txns []coin.Transaction) (coin.AddressUxOuts, error) {
	if len(txns) == 0 || len(addrs) == 0 {
		return nil, nil
	}

	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}

	auxs := make(coin.AddressUxOuts, len(addrs))

	for _, txn := range txns {
		for i, o := range txn.Out {
			if _, ok := addrm[o.Address]; ok {
				uxout, err := coin.CreateUnspent(bh, txn, i)
				if err != nil {
					return nil, err
				}

				auxs[o.Address] = append(auxs[o.Address], uxout)
			}
		}
	}

	return auxs, nil
}

// GetIncomingOutputs returns all predicted incoming outputs.
func (utp *UnconfirmedTransactionPool) GetIncomingOutputs(tx *dbutil.Tx, bh coin.BlockHeader) (coin.UxArray, error) {
	var outs coin.UxArray

	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTransaction) error {
		outs = append(outs, coin.CreateUnspents(bh, txn.Transaction)...)
		return nil
	}); err != nil {
		return nil, err
	}

	return outs, nil
}

// Get returns the unconfirmed transaction of given tx hash.
func (utp *UnconfirmedTransactionPool) Get(tx *dbutil.Tx, hash cipher.SHA256) (*UnconfirmedTransaction, error) {
	return utp.txns.get(tx, hash)
}

// GetFiltered returns all transactions that can pass the filter
func (utp *UnconfirmedTransactionPool) GetFiltered(tx *dbutil.Tx, filter func(UnconfirmedTransaction) bool) ([]UnconfirmedTransaction, error) {
	var txns []UnconfirmedTransaction

	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTransaction) error {
		if filter(txn) {
			txns = append(txns, txn)
		}
		return nil
	}); err != nil {
		logger.Errorf("GetFiltered error: %v", err)
		return nil, err
	}

	return txns, nil
}

// GetHashes returns transaction hashes that can pass the filter
func (utp *UnconfirmedTransactionPool) GetHashes(tx *dbutil.Tx, filter func(UnconfirmedTransaction) bool) ([]cipher.SHA256, error) {
	var hashes []cipher.SHA256

	if err := utp.txns.forEach(tx, func(hash cipher.SHA256, txn UnconfirmedTransaction) error {
		if filter(txn) {
			hashes = append(hashes, hash)
		}
		return nil
	}); err != nil {
		logger.Errorf("GetHashes error: %v", err)
		return nil, err
	}

	return hashes, nil
}

// ForEach iterate the pool with given callback function
func (utp *UnconfirmedTransactionPool) ForEach(tx *dbutil.Tx, f func(cipher.SHA256, UnconfirmedTransaction) error) error {
	return utp.txns.forEach(tx, f)
}

// GetUnspentsOfAddr returns unspent outputs of given address in unspent tx pool
func (utp *UnconfirmedTransactionPool) GetUnspentsOfAddr(tx *dbutil.Tx, addr cipher.Address) (coin.UxArray, error) {
	return utp.unspent.getByAddr(tx, addr)
}

// IsValid can be used as filter function
func IsValid(tx UnconfirmedTransaction) bool {
	return tx.IsValid == 1
}

// All use as return all filter
func All(tx UnconfirmedTransaction) bool {
	return true
}

// Len returns the number of unconfirmed transactions
func (utp *UnconfirmedTransactionPool) Len(tx *dbutil.Tx) (uint64, error) {
	return utp.txns.len(tx)
}
