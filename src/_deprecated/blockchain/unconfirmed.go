package blockchain

import (
    "github.com/skycoin/skycoin/src/coin"
    //"github.com/skycoin/skycoin/src/util"
    "time"
)

type UnconfirmedTxn struct {
    Txn coin.Transaction
    // Time the txn was last received
    Received int64
    // Time the txn was last checked against the blockchain
    Checked int64
    // Last time we announced this txn
    Announced int64 //unix time
}

// Returns the coin.Transaction's hash
func (self *UnconfirmedTxn) Hash() cipher.SHA256 {
    return self.Txn.Hash()
}

// Manages unconfirmed transactions
type UnconfirmedTxnPool struct {
    Txns map[cipher.SHA256]UnconfirmedTxn
}

func NewUnconfirmedTxnPool() *UnconfirmedTxnPool {
    return &UnconfirmedTxnPool{
        Txns:    make(map[cipher.SHA256]UnconfirmedTxn),
        //Unspent: coin.NewUnspentPool(),
    }
}

func (self *UnconfirmedTxnPool) SetAnnounced(h cipher.SHA256, t int64) {
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
//func (self *UnconfirmedTxnPool) InjectTxn(bc *coin.Blockchain,
//    t coin.Transaction, addrs map[cipher.Address]byte, didAnnounce bool) error {
func (self *UnconfirmedTxnPool) InjectTxn(t coin.Transaction) error {

    now := time.Now().Unix()
    //announcedAt := util.ZeroTime()

    ut := UnconfirmedTxn{
        Txn:          t,
        Received:     now,
        Checked:      now,
        Announced:    0, //set to 0 until announced
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
func (self *UnconfirmedTxnPool) removeTxn(h cipher.SHA256) {
    _, ok := self.Txns[h]
    if !ok {
        return
    }
    delete(self.Txns, h)
}

// Removes multiple txns at once. Slightly more efficient than a series of
// single RemoveTxns
// Note -- efficiency does not matter. Only doing ~10 transactions/second

func (self *UnconfirmedTxnPool) removeTxns(hashes []cipher.SHA256) {
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
    now := time.Now().Unix()
    toRemove := make([]cipher.SHA256, 0)
    for k, t := range self.Txns {
        if now - t.Checked >= int64(checkPeriod) {
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
func (self *UnconfirmedTxnPool) FilterKnown(txns []cipher.SHA256) []cipher.SHA256 {
    unknown := make([]cipher.SHA256, 0)
    for _, h := range txns {
        if _, known := self.Txns[h]; !known {
            unknown = append(unknown, h)
        }
    }
    return unknown
}

// Returns all known coin.Transactions from the pool, given hashes to select
func (self *UnconfirmedTxnPool) GetKnown(txns []cipher.SHA256) coin.Transactions {
    known := make(coin.Transactions, 0)
    for _, h := range txns {
        if txn, have := self.Txns[h]; have {
            known = append(known, txn.Txn)
        }
    }
    return known
}
*/