package visor

import (
    "github.com/skycoin/skycoin/src/coin"
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
    // We are a spender
    IsOurSpend bool
    // We are a receiver
    IsOurReceive bool
}

// Manages unconfirmed transactions
type UnconfirmedTxnPool struct {
    Txns map[coin.SHA256]UnconfirmedTxn
    // Predicted unspents, assuming txns are valid.  Needed to predict
    // our future balance and avoid double spending our own coins
    Unspent coin.UnspentPool
}

func NewUnconfirmedTxnPool() *UnconfirmedTxnPool {
    return &UnconfirmedTxnPool{
        Txns:    make(map[coin.SHA256]UnconfirmedTxn),
        Unspent: coin.NewUnspentPool(),
    }
}

func (self *UnconfirmedTxnPool) SetAnnounced(h coin.SHA256, t time.Time) {
    tx, ok := self.Txns[h]
    if ok {
        tx.Announced = t
        self.Txns[h] = tx
    }
}

// Adds a coin.Transaction to the pool
func (self *UnconfirmedTxnPool) RecordTxn(bc *coin.Blockchain,
    t coin.Transaction, addrs map[coin.Address]int, didAnnounce bool) error {
    if err := bc.VerifyTransaction(t); err != nil {
        return err
    }
    now := time.Now().UTC()
    announcedAt := time.Unix(0, 0)
    if didAnnounce {
        announcedAt = now
    }
    ut := UnconfirmedTxn{
        Txn:          t,
        Received:     now,
        Checked:      now,
        Announced:    announcedAt,
        IsOurReceive: false,
        IsOurSpend:   false,
    }
    // Add predicted unspents
    for _, ux := range bc.TxUxOut(t, coin.BlockHeader{}) {
        self.Unspent.Add(ux)
    }
    // Check if this is one of our receiving txns
    for _, to := range t.Out {
        if _, ok := addrs[to.DestinationAddress]; ok {
            ut.IsOurReceive = true
        }
    }
    // Check if this is one of our spending txns
    for _, ti := range t.In {
        if ux, ok := bc.Unspent.Get(ti.UxOut); ok {
            if _, ok := addrs[ux.Body.Address]; ok {
                ut.IsOurSpend = true
            }
        }
    }
    self.Txns[t.Header.Hash] = ut

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

// Remove a single txn
func (self *UnconfirmedTxnPool) removeTxn(bc *coin.Blockchain, h coin.SHA256) {
    t, ok := self.Txns[h]
    if !ok {
        return
    }
    delete(self.Txns, h)
    outputs := bc.TxUxOut(t.Txn, coin.BlockHeader{})
    hashes := make([]coin.SHA256, len(outputs))
    for _, o := range outputs {
        hashes = append(hashes, o.Hash())
    }
    self.Unspent.DelMultiple(hashes)
}

// Removes multiple txns at once. Slightly more efficient than a series of
// single RemoveTxns
func (self *UnconfirmedTxnPool) removeTxns(bc *coin.Blockchain,
    hashes []coin.SHA256) {
    uxo := make([]coin.UxOut, 0, len(hashes))
    for _, h := range hashes {
        t, ok := self.Txns[h]
        if ok {
            delete(self.Txns, h)
            uxo = append(uxo, bc.TxUxOut(t.Txn, coin.BlockHeader{})...)
        }
    }
    uxhashes := make([]coin.SHA256, len(uxo))
    for _, o := range uxo {
        uxhashes = append(uxhashes, o.Hash())
    }
    self.Unspent.DelMultiple(uxhashes)
}

// Checks all unconfirmed txns against the blockchain. maxAge is how long
// we'll hold a txn regardless of whether it has been invalidated.
// checkPeriod is how often we check the txn against the blockchain.
func (self *UnconfirmedTxnPool) Refresh(bc *coin.Blockchain,
    checkPeriod, maxAge time.Duration) {
    now := time.Now().UTC()
    toRemove := make([]coin.SHA256, 0)
    for k, t := range self.Txns {
        if t.Received.Add(maxAge).Before(now) {
            toRemove = append(toRemove, k)
        } else if t.Checked.Add(checkPeriod).Before(now) {
            if bc.VerifyTransaction(t.Txn) == nil {
                t.Checked = now
                self.Txns[k] = t
            } else {
                toRemove = append(toRemove, k)
            }
        }
    }
    self.removeTxns(bc, toRemove)
}

// Returns transactions in which we are a party and have not been announced
// in ago duration
func (self *UnconfirmedTxnPool) GetOwnedTransactionsSince(ago time.Duration) []*UnconfirmedTxn {
    txns := make([]*UnconfirmedTxn, 0)
    now := time.Now().UTC()
    for _, tx := range self.Txns {
        if (tx.IsOurSpend || tx.IsOurReceive) &&
            tx.Announced.Add(ago).Before(now) {
            txns = append(txns, &tx)
        }
    }
    return txns
}

// Removes confirmed txns from the pool
func (self *UnconfirmedTxnPool) RemoveTransactions(bc *coin.Blockchain,
    txns coin.Transactions) {
    toRemove := make([]coin.SHA256, 0, len(txns))
    for _, tx := range txns {
        toRemove = append(toRemove, tx.Header.Hash)
    }
    self.removeTxns(bc, toRemove)
}

// Returns txn hashes with known ones removed
func (self *UnconfirmedTxnPool) FilterKnown(txns []coin.SHA256) []coin.SHA256 {
    unknown := make([]coin.SHA256, 0)
    for _, h := range txns {
        _, known := self.Txns[h]
        if !known {
            unknown = append(unknown, h)
        }
    }
    return unknown
}

// Returns all known coin.Transactions from the pool, given hashes to select
func (self *UnconfirmedTxnPool) GetKnown(txns []coin.SHA256) coin.Transactions {
    known := make(coin.Transactions, 0)
    for _, h := range txns {
        txn, have := self.Txns[h]
        if have {
            known = append(known, txn.Txn)
        }
    }
    return known
}
