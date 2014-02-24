package blockchain

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "time"
)

type UnconfirmedTxn struct {
    Txn coin.Transaction
    // Time the txn was last received
    Received time.Time
    // Time the txn was last checked against the blockchain
    Checked time.Time
    // Last time we announced this txn
    Announced time.Time
}

// Returns the coin.Transaction's hash
func (self *UnconfirmedTxn) Hash() coin.SHA256 {
    return self.Txn.Hash()
}

// Manages unconfirmed transactions
type UnconfirmedTxnPool struct {
    Txns map[coin.SHA256]UnconfirmedTxn
}

func NewUnconfirmedTxnPool() *UnconfirmedTxnPool {
    return &UnconfirmedTxnPool{
        Txns:    make(map[coin.SHA256]UnconfirmedTxn),
        //Unspent: coin.NewUnspentPool(),
    }
}

func (self *UnconfirmedTxnPool) SetAnnounced(h coin.SHA256, t time.Time) {
    if tx, ok := self.Txns[h]; ok {
        tx.Announced = t
        self.Txns[h] = tx
    }
}

// Note: let wallet read the unconfirmed and the unspents and make its own decision

// Note: use notion of transaction "degree". A "degree 0" transaction requires zero
// other transactions to execute before a spend can occur.  A "degree 1" transaction
// spends inputs that are created in a pending/unconfirmed transaction.
// A "degree 2" transaction spends inputs that are created in an unconfirmed transaction
// that requires inputs created by an unconfirmed transaction.
// A transaction's degree is one greater than the highest degree of any of the transactions
// creating an input spent by the transaction

// Adds a coin.Transaction to the pool
//func (self *UnconfirmedTxnPool) RecordTxn(bc *coin.Blockchain,
//    t coin.Transaction, addrs map[coin.Address]byte, didAnnounce bool) error {
func (self *UnconfirmedTxnPool) RecordTxn(t coin.Transaction) error {

    now := util.Now()
    announcedAt := util.ZeroTime()
    if didAnnounce {
        announcedAt = now
    }
    ut := UnconfirmedTxn{
        Txn:          t,
        Received:     now,
        Checked:      now,
        Announced:    announcedAt,
    }
    self.Txns[t.Hash()] = ut
    return nil
}

// Returns underlying coin.Transactions
func (self *UnconfirmedTxnPool) RawTxns() coin.Transactions {
    txns := make(coin.Transactions, 0, len(self.Txns))
    for _, t := range self.Txns {
        txns = append(txns, t.Txn)
    }
    return txns
}

// Remove a single txn from the unconfirmed transaction pool
func (self *UnconfirmedTxnPool) removeTxn(h coin.SHA256) {
    _, ok := self.Txns[h]
    if !ok {
        return
    }
    delete(self.Txns, h)
}

// Removes multiple txns at once. Slightly more efficient than a series of
// single RemoveTxns
// Note -- efficiency does not matter. Only doing ~10 transactions/second

func (self *UnconfirmedTxnPool) removeTxns(hashes []coin.SHA256) {
    for _, h := range hashes {
        if _, ok := self.Txns[h]; ok {
            delete(self.Txns, h)
        }
    }
}

// Duplicate of removeTxns
// Removes confirmed txns from the pool
func (self *UnconfirmedTxnPool) RemoveTransactions(txns coin.Transactions) {
    for _, tx := range txns {
        self.removeTxn(tx.Hash())
    }
}

// Checks all unconfirmed txns against the blockchain. maxAge is how long
// we'll hold a txn regardless of whether it has been invalidated.
// checkPeriod is how often we check the txn against the blockchain.
func (self *UnconfirmedTxnPool) Refresh(bc *coin.Blockchain, checkPeriod int) {
    now := util.Now()
    toRemove := make([]coin.SHA256, 0)
    for k, t := range self.Txns {
        if now.Sub(t.Checked) >= checkPeriod {
            if bc.VerifyTransaction(t.Txn) == nil {
                t.Checked = now
                self.Txns[k] = t
            } else {
                toRemove = append(toRemove, k)
            }
        }
    }
    self.removeTxns(toRemove)
}

/*
// Returns txn hashes with known ones removed
func (self *UnconfirmedTxnPool) FilterKnown(txns []coin.SHA256) []coin.SHA256 {
    unknown := make([]coin.SHA256, 0)
    for _, h := range txns {
        if _, known := self.Txns[h]; !known {
            unknown = append(unknown, h)
        }
    }
    return unknown
}

// Returns all known coin.Transactions from the pool, given hashes to select
func (self *UnconfirmedTxnPool) GetKnown(txns []coin.SHA256) coin.Transactions {
    known := make(coin.Transactions, 0)
    for _, h := range txns {
        if txn, have := self.Txns[h]; have {
            known = append(known, txn.Txn)
        }
    }
    return known
}
*/