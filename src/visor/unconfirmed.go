package visor

import (
	"errors"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// UnconfirmedTxnsBkt holds unconfirmed transactions
	UnconfirmedTxnsBkt = []byte("unconfirmed_txns")
	// UnconfirmedUnspentsBkt holds unconfirmed unspent outputs
	UnconfirmedUnspentsBkt = []byte("unconfirmed_unspents")

	errUpdateObjectDoesNotExist = errors.New("object does not exist in bucket")
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
type unconfirmedTxns struct{}

func (utb *unconfirmedTxns) get(tx *dbutil.Tx, hash cipher.SHA256) (*UnconfirmedTxn, error) {
	var txn UnconfirmedTxn

	if ok, err := dbutil.GetBucketObjectDecoded(tx, UnconfirmedTxnsBkt, []byte(hash.Hex()), &txn); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return &txn, nil
}

func (utb *unconfirmedTxns) put(tx *dbutil.Tx, v *UnconfirmedTxn) error {
	return dbutil.PutBucketValue(tx, UnconfirmedTxnsBkt, []byte(v.Hash().Hex()), encoder.Serialize(v))
}

func (utb *unconfirmedTxns) update(tx *dbutil.Tx, hash cipher.SHA256, f func(v *UnconfirmedTxn) error) error {
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

func (utb *unconfirmedTxns) getAll(tx *dbutil.Tx) ([]UnconfirmedTxn, error) {
	var txns []UnconfirmedTxn

	if err := dbutil.ForEach(tx, UnconfirmedTxnsBkt, func(_, v []byte) error {
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

func (utb *unconfirmedTxns) hasKey(tx *dbutil.Tx, hash cipher.SHA256) (bool, error) {
	return dbutil.BucketHasKey(tx, UnconfirmedTxnsBkt, []byte(hash.Hex()))
}

func (utb *unconfirmedTxns) forEach(tx *dbutil.Tx, f func(hash cipher.SHA256, tx UnconfirmedTxn) error) error {
	return dbutil.ForEach(tx, UnconfirmedTxnsBkt, func(k, v []byte) error {
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

func (utb *unconfirmedTxns) length(tx *dbutil.Tx) (uint64, error) {
	return dbutil.Len(tx, UnconfirmedTxnsBkt)
}

type txUnspents struct{}

func (txus *txUnspents) put(tx *dbutil.Tx, hash cipher.SHA256, uxs coin.UxArray) error {
	return dbutil.PutBucketValue(tx, UnconfirmedUnspentsBkt, []byte(hash.Hex()), encoder.Serialize(uxs))
}

func (txus *txUnspents) get(tx *dbutil.Tx, hash cipher.SHA256) (coin.UxArray, error) {
	var uxs coin.UxArray

	if ok, err := dbutil.GetBucketObjectDecoded(tx, UnconfirmedUnspentsBkt, []byte(hash.Hex()), &uxs); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return uxs, nil
}

func (txus *txUnspents) length(tx *dbutil.Tx) (uint64, error) {
	return dbutil.Len(tx, UnconfirmedUnspentsBkt)
}

func (txus *txUnspents) delete(tx *dbutil.Tx, hash cipher.SHA256) error {
	return dbutil.Delete(tx, UnconfirmedUnspentsBkt, []byte(hash.Hex()))
}

func (txus *txUnspents) getByAddr(tx *dbutil.Tx, a cipher.Address) (coin.UxArray, error) {
	var uxo coin.UxArray

	if err := dbutil.ForEach(tx, UnconfirmedUnspentsBkt, func(_, v []byte) error {
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

func (txus *txUnspents) forEach(tx *dbutil.Tx, f func(cipher.SHA256, coin.UxArray) error) error {
	return dbutil.ForEach(tx, UnconfirmedUnspentsBkt, func(k, v []byte) error {
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

	return &UnconfirmedTxnPool{
		db:      db,
		txns:    &unconfirmedTxns{},
		unspent: &txUnspents{},
	}, nil
}

// SetTxnsAnnounced updates announced time of specific tx
func (utp *UnconfirmedTxnPool) SetTxnsAnnounced(tx *dbutil.Tx, hashes map[cipher.SHA256]int64) error {
	var txns []*UnconfirmedTxn
	for h, t := range hashes {
		txn, err := utp.txns.get(tx, h)
		if err != nil {
			return err
		}

		if txn == nil {
			logger.Warningf("UnconfirmedTxnPool.SetTxnsAnnounced: UnconfirmedTxn %s not found in DB", h.Hex())
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

func createUnconfirmedTxn(txn coin.Transaction) UnconfirmedTxn {
	now := utc.Now()
	return UnconfirmedTxn{
		Txn:       txn,
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
func (utp *UnconfirmedTxnPool) InjectTransaction(tx *dbutil.Tx, bc Blockchainer, txn coin.Transaction, maxSize int) (bool, *ErrTxnViolatesSoftConstraint, error) {
	var isValid int8 = 1
	var softErr *ErrTxnViolatesSoftConstraint
	if err := bc.VerifySingleTxnSoftHardConstraints(tx, txn, maxSize); err != nil {
		logger.Warningf("bc.VerifySingleTxnSoftHardConstraints failed for txn %s: %v", txn.TxIDHex(), err)
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

	hash := txn.Hash()

	known, err := utp.txns.hasKey(tx, hash)
	if err != nil {
		logger.Errorf("InjectTransaction check txn exists failed: %v", err)
		return false, nil, err
	}

	// Update if we already have this txn
	if known {
		if err := utp.txns.update(tx, hash, func(utxn *UnconfirmedTxn) error {
			now := utc.Now().UnixNano()
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

	utx := createUnconfirmedTxn(txn)
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
	if err := utp.unspent.put(tx, hash, coin.CreateUnspents(head.Head, txn)); err != nil {
		logger.Errorf("InjectTransaction put new unspent outputs: %v", err)
		return false, nil, err
	}

	return false, softErr, nil
}

// RawTxns returns underlying coin.Transactions
func (utp *UnconfirmedTxnPool) RawTxns(tx *dbutil.Tx) (coin.Transactions, error) {
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
func (utp *UnconfirmedTxnPool) removeTxn(tx *dbutil.Tx, txHash cipher.SHA256) error {
	if err := utp.txns.delete(tx, txHash); err != nil {
		return err
	}

	return utp.unspent.delete(tx, txHash)
}

// RemoveTransactions remove transactions with dbutil.Tx
func (utp *UnconfirmedTxnPool) RemoveTransactions(tx *dbutil.Tx, txHashes []cipher.SHA256) error {
	for i := range txHashes {
		if err := utp.removeTxn(tx, txHashes[i]); err != nil {
			return err
		}
	}

	return nil
}

// Refresh checks all unconfirmed txns against the blockchain.
// If the transaction becomes invalid it is marked invalid.
// If the transaction becomes valid it is marked valid and is returned to the caller.
func (utp *UnconfirmedTxnPool) Refresh(tx *dbutil.Tx, bc Blockchainer, maxBlockSize int) ([]cipher.SHA256, error) {
	utxns, err := utp.txns.getAll(tx)
	if err != nil {
		return nil, err
	}

	now := utc.Now()
	var nowValid []cipher.SHA256

	for _, utxn := range utxns {
		utxn.Checked = now.UnixNano()

		err := bc.VerifySingleTxnSoftHardConstraints(tx, utxn.Txn, maxBlockSize)

		switch err.(type) {
		case ErrTxnViolatesSoftConstraint, ErrTxnViolatesHardConstraint:
			utxn.IsValid = 0
		case nil:
			if utxn.IsValid == 0 {
				nowValid = append(nowValid, utxn.Hash())
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
func (utp *UnconfirmedTxnPool) RemoveInvalid(tx *dbutil.Tx, bc Blockchainer) ([]cipher.SHA256, error) {
	var removeUtxns []cipher.SHA256

	utxns, err := utp.txns.getAll(tx)
	if err != nil {
		return nil, err
	}

	for _, utxn := range utxns {
		err := bc.VerifySingleTxnHardConstraints(tx, utxn.Txn)
		if err != nil {
			switch err.(type) {
			case ErrTxnViolatesHardConstraint:
				removeUtxns = append(removeUtxns, utxn.Hash())
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

// GetUnknown returns txn hashes with known ones removed
func (utp *UnconfirmedTxnPool) GetUnknown(tx *dbutil.Tx, txns []cipher.SHA256) ([]cipher.SHA256, error) {
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
func (utp *UnconfirmedTxnPool) GetKnown(tx *dbutil.Tx, txns []cipher.SHA256) (coin.Transactions, error) {
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
func (utp *UnconfirmedTxnPool) RecvOfAddresses(tx *dbutil.Tx, bh coin.BlockHeader, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}

	auxs := make(coin.AddressUxOuts, len(addrs))
	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTxn) error {
		for i, o := range txn.Txn.Out {
			if _, ok := addrm[o.Address]; ok {
				uxout, err := coin.CreateUnspent(bh, txn.Txn, i)
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
func (utp *UnconfirmedTxnPool) GetIncomingOutputs(tx *dbutil.Tx, bh coin.BlockHeader) (coin.UxArray, error) {
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
func (utp *UnconfirmedTxnPool) Get(tx *dbutil.Tx, hash cipher.SHA256) (*UnconfirmedTxn, error) {
	return utp.txns.get(tx, hash)
}

// GetTxns returns all transactions that can pass the filter
func (utp *UnconfirmedTxnPool) GetTxns(tx *dbutil.Tx, filter func(UnconfirmedTxn) bool) ([]UnconfirmedTxn, error) {
	var txns []UnconfirmedTxn

	if err := utp.txns.forEach(tx, func(_ cipher.SHA256, txn UnconfirmedTxn) error {
		if filter(txn) {
			txns = append(txns, txn)
		}
		return nil
	}); err != nil {
		logger.Errorf("GetTxns error: %v", err)
		return nil, err
	}

	return txns, nil
}

// GetTxHashes returns transaction hashes that can pass the filter
func (utp *UnconfirmedTxnPool) GetTxHashes(tx *dbutil.Tx, filter func(UnconfirmedTxn) bool) ([]cipher.SHA256, error) {
	var hashes []cipher.SHA256

	if err := utp.txns.forEach(tx, func(hash cipher.SHA256, txn UnconfirmedTxn) error {
		if filter(txn) {
			hashes = append(hashes, hash)
		}
		return nil
	}); err != nil {
		logger.Errorf("GetTxHashes error: %v", err)
		return nil, err
	}

	return hashes, nil
}

// ForEach iterate the pool with given callback function
func (utp *UnconfirmedTxnPool) ForEach(tx *dbutil.Tx, f func(cipher.SHA256, UnconfirmedTxn) error) error {
	return utp.txns.forEach(tx, f)
}

// GetUnspentsOfAddr returns unspent outputs of given address in unspent tx pool
func (utp *UnconfirmedTxnPool) GetUnspentsOfAddr(tx *dbutil.Tx, addr cipher.Address) (coin.UxArray, error) {
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
func (utp *UnconfirmedTxnPool) Len(tx *dbutil.Tx) (uint64, error) {
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
