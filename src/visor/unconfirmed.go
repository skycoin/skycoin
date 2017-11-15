package visor

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/utc"
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
func (utp *UnconfirmedTxnPool) SetTxnsAnnounced(tx *bolt.Tx, hashes []cipher.SHA256, t time.Time) error {
	logger.Debug("UnconfirmedTxnPool.SetTxnsAnnounced: %d txns", len(hashes))

	var txns []*UnconfirmedTxn
	for _, h := range hashes {
		txn, err := utp.txns.get(tx, h)
		if err != nil {
			return err
		}

		if txn == nil {
			logger.Warning("UnconfirmedTxnPool.SetTxnsAnnounced: UnconfirmedTxn %s not found in DB", h.Hex())
			continue
		}

		txns = append(txns, txn)
	}

	now := t.UnixNano()
	for _, txn := range txns {
		txn.Announced = now
		if err := utp.txns.put(tx, txn); err != nil {
			return err
		}
	}

	return nil
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

	fee, err := bc.TransactionFee(tx, head.Time())(&t)
	if err != nil {
		return false, err
	}

	if err := VerifyTransactionFee(&t, fee); err != nil {
		return false, err
	}

	hash := t.Hash()

	if err := bc.VerifyTransaction(tx, t); err != nil {
		return false, err
	}

	known, err := utp.txns.hasKey(tx, hash)
	if err != nil {
		return false, err
	}

	// Update if we already have this txn
	// TODO -- why update to IsValid if we already have the txn?
	// It looks like the other code assumes IsValid txns are txns that we
	// created due to spending.
	if known {
		if err := utp.txns.update(tx, hash, func(txn *UnconfirmedTxn) error {
			now := utc.Now().UnixNano()
			txn.Received = now
			txn.Checked = now
			txn.IsValid = 1
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
	utxns, err := utp.txns.getAll(tx)
	if err != nil {
		return nil, err
	}

	now := utc.Now().UnixNano()

	var hashes []cipher.SHA256
	for _, txn := range utxns {
		txn.Checked = now
		if txn.IsValid == 0 && bc.VerifyTransaction(tx, txn.Txn) == nil {
			txn.IsValid = 1
			hashes = append(hashes, txn.Hash())
		}
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

// // SpendsOfAddresses returns all unconfirmed coin.UxOut spends of addresses
// // Looks at all inputs for unconfirmed txns, gets their source UxOut from the
// // blockchain's unspent pool, and returns as coin.AddressUxOuts
// func (utp *UnconfirmedTxnPool) SpendsOfAddresses(tx *bolt.Tx, addrs []cipher.Address, uxa coin.UxArray) (coin.AddressUxOuts, error) {
// 	addrm := make(map[cipher.Address]struct{}, len(addrs))
// 	for _, addr := range addrs {
// 		addrm[addr] = struct{}{}
// 	}

// 	auxs := make(coin.AddressUxOuts, len(addrs))

// 	uxs, err := utp.GetSpendingOutputs(tx, unspent)
// 	if err != nil {
// 		return err
// 	}

// 	for _, ux := range uxs {
// 		if _, ok := addrm[ux.Body.Address]; ok {
// 			auxs[ux.Body.Address] = append(auxs[ux.Body.Address], ux)
// 		}
// 	}

// 	return auxs, nil
// }

// // GetSpendingOutputs returns all spending outputs in unconfirmed tx pool.
// func (utp *UnconfirmedTxnPool) GetSpendingOutputs(tx *bolt.Tx, unspent blockdb.UnspentPool) (coin.UxArray, error) {
// 	var inputs []cipher.SHA256
// 	txns, err := utps.txns.getAll(tx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, txn := range txns {
// 		inputs = append(inputs, txn.Txn.In...)
// 	}

// 	return unspent.GetArray(tx, inputs)
// }

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
