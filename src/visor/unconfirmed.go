package visor

import (
	"errors"

	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
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
	Received time.Time
	// Time the txn was last checked against the blockchain
	Checked time.Time
	// Last time we announced this txn
	Announced time.Time
}

// Hash returns the coin.Transaction's hash
func (ut *UnconfirmedTxn) Hash() cipher.SHA256 {
	return ut.Txn.Hash()
}

// UnconfirmedTxnPool manages unconfirmed transactions
type UnconfirmedTxnPool struct {
	Txns map[cipher.SHA256]UnconfirmedTxn
	// Predicted unspents, assuming txns are valid.  Needed to predict
	// our future balance and avoid double spending our own coins
	// Maps from Transaction.Hash() to UxArray.
	Unspent TxnUnspents
}

// NewUnconfirmedTxnPool creates an UnconfirmedTxnPool instance
func NewUnconfirmedTxnPool() *UnconfirmedTxnPool {
	return &UnconfirmedTxnPool{
		Txns:    make(map[cipher.SHA256]UnconfirmedTxn),
		Unspent: make(TxnUnspents),
	}
}

// SetAnnounced updates announced time of specific tx
func (utp *UnconfirmedTxnPool) SetAnnounced(h cipher.SHA256, t time.Time) {
	if tx, ok := utp.Txns[h]; ok {
		tx.Announced = t
		utp.Txns[h] = tx
	}
}

// Creates an unconfirmed transaction
func (utp *UnconfirmedTxnPool) createUnconfirmedTxn(bcUnsp *coin.UnspentPool,
	t coin.Transaction) UnconfirmedTxn {
	now := util.Now()
	return UnconfirmedTxn{
		Txn:       t,
		Received:  now,
		Checked:   now,
		Announced: util.ZeroTime(),
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
	ut, ok := utp.Txns[h]
	if ok {
		now := util.Now()
		ut.Received = now
		ut.Checked = now
		utp.Txns[h] = ut
		return nil, true
	}

	// Add txn to index
	unspent := bc.GetUnspent()
	utp.Txns[h] = utp.createUnconfirmedTxn(unspent, t)
	// Add predicted unspents
	utp.Unspent[h] = coin.CreateUnspents(bc.Head().Head, t)

	return nil, false
}

// RawTxns returns underlying coin.Transactions
func (utp *UnconfirmedTxnPool) RawTxns() coin.Transactions {
	txns := make(coin.Transactions, len(utp.Txns))
	i := 0
	for _, t := range utp.Txns {
		txns[i] = t.Txn
		i++
	}
	return txns
}

// Remove a single txn by hash
func (utp *UnconfirmedTxnPool) removeTxn(bc *Blockchain, txHash cipher.SHA256) {
	delete(utp.Txns, txHash)
	delete(utp.Unspent, txHash)
}

// Removes multiple txns at once. Slightly more efficient than a series of
// single RemoveTxns.  Hashes is an array of Transaction hashes.
func (utp *UnconfirmedTxnPool) removeTxns(bc *Blockchain,
	hashes []cipher.SHA256) {
	for i := range hashes {
		delete(utp.Txns, hashes[i])
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
	for k, t := range utp.Txns {
		if now.Sub(t.Received) >= maxAge {
			toRemove = append(toRemove, k)
		} else if now.Sub(t.Checked) >= checkPeriod {
			if bc.VerifyTransaction(t.Txn) == nil {
				t.Checked = now
				utp.Txns[k] = t
			} else {
				toRemove = append(toRemove, k)
			}
		}
	}
	utp.removeTxns(bc, toRemove)
}

// FilterKnown returns txn hashes with known ones removed
func (utp *UnconfirmedTxnPool) FilterKnown(txns []cipher.SHA256) []cipher.SHA256 {
	var unknown []cipher.SHA256
	for _, h := range txns {
		if _, known := utp.Txns[h]; !known {
			unknown = append(unknown, h)
		}
	}
	return unknown
}

// GetKnown returns all known coin.Transactions from the pool, given hashes to select
func (utp *UnconfirmedTxnPool) GetKnown(txns []cipher.SHA256) coin.Transactions {
	var known coin.Transactions
	for _, h := range txns {
		if txn, have := utp.Txns[h]; have {
			known = append(known, txn.Txn)
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
	for _, utx := range utp.Txns {
		for _, h := range utx.Txn.In {
			if ux, ok := bcUnspent.Get(h); ok {
				if _, ok := a[ux.Body.Address]; ok {
					auxs[ux.Body.Address] = append(auxs[ux.Body.Address], ux)
				}
			}
		}
	}
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
	for _, tx := range utp.Txns {
		for _, in := range tx.Txn.In {
			if ux, ok := bcUnspent.Get(in); ok {
				outs = append(outs, NewReadableOutput(ux))
			}
		}
	}
	return outs
}

// AllIncommingOutputs returns all predicted incomming outputs.
func (utp *UnconfirmedTxnPool) AllIncommingOutputs(bh coin.BlockHeader) []ReadableOutput {
	outs := []ReadableOutput{}
	for _, tx := range utp.Txns {
		uxOuts := coin.CreateUnspents(bh, tx.Txn)
		for _, ux := range uxOuts {
			outs = append(outs, NewReadableOutput(ux))
		}
	}
	return outs
}
