package visor

import (
    "errors"
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
    // We are a spender
    IsOurSpend bool
    // We are a receiver
    IsOurReceive bool
}

// Returns the coin.Transaction's hash
func (self *UnconfirmedTxn) Hash() coin.SHA256 {
    return self.Txn.Hash()
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
    if tx, ok := self.Txns[h]; ok {
        tx.Announced = t
        self.Txns[h] = tx
    }
}

// Creates an unconfirmed transaction
func (self *UnconfirmedTxnPool) createUnconfirmedTxn(bcUnsp *coin.UnspentPool,
    t coin.Transaction, addrs map[coin.Address]byte) UnconfirmedTxn {
    now := util.Now()
    ut := UnconfirmedTxn{
        Txn:          t,
        Received:     now,
        Checked:      now,
        Announced:    util.ZeroTime(),
        IsOurReceive: false,
        IsOurSpend:   false,
    }

    // Check if this unspent is related to us
    if addrs != nil {
        // Check if this is one of our receiving txns
        for i, _ := range t.Out {
            if _, ok := addrs[t.Out[i].Address]; ok {
                ut.IsOurReceive = true
                break
            }
        }
        // Check if this is one of our spending txns
        for i, _ := range t.In {
            if ux, ok := bcUnsp.Get(t.In[i]); ok {
                if _, ok := addrs[ux.Body.Address]; ok {
                    ut.IsOurSpend = true
                    break
                }
            }
        }
    }

    return ut
}

// Adds a coin.Transaction to the pool, or updates an existing one's timestamps
// Returns an error if txn is invalid, and whether the transaction already
// existed in the pool.
func (self *UnconfirmedTxnPool) RecordTxn(bc *coin.Blockchain,
    t coin.Transaction, addrs map[coin.Address]byte, maxSize int,
    burnFactor uint64) (error, bool) {
    if t.Size() > maxSize {
        return errors.New("Transaction too large"), false
    }
    if fee, err := bc.TransactionFee(&t); err != nil {
        return err, false
    } else if burnFactor != 0 && t.OutputHours()/burnFactor > fee {
        return errors.New("Transaction fee minimum not met"), false
    }
    if err := bc.VerifyTransaction(t); err != nil {
        return err, false
    }

    // Update if we already have this txn
    ut, ok := self.Txns[t.Hash()]
    if ok {
        now := util.Now()
        ut.Received = now
        ut.Checked = now
        self.Txns[ut.Txn.Hash()] = ut
        return nil, true
    }

    // Add txn to index
    self.Txns[t.Hash()] = self.createUnconfirmedTxn(&bc.Unspent, t, addrs)
    // Add predicted unspents
    uxs := coin.CreateExpectedUnspents(t)
    for i, _ := range uxs {
        self.Unspent.Add(uxs[i])
    }

    return nil, false
}

// Returns underlying coin.Transactions
func (self *UnconfirmedTxnPool) RawTxns() coin.Transactions {
    txns := make(coin.Transactions, len(self.Txns))
    i := 0
    for _, t := range self.Txns {
        txns[i] = t.Txn
        i++
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
    outputs := coin.CreateExpectedUnspents(t.Txn)
    hashes := make([]coin.SHA256, len(outputs))
    for i, _ := range outputs {
        hashes[i] = outputs[i].Hash()
    }
    self.Unspent.DelMultiple(hashes)
}

// Removes multiple txns at once. Slightly more efficient than a series of
// single RemoveTxns
func (self *UnconfirmedTxnPool) removeTxns(bc *coin.Blockchain,
    hashes []coin.SHA256) {
    uxo := make([]coin.UxOut, 0)
    for i, _ := range hashes {
        if t, ok := self.Txns[hashes[i]]; ok {
            delete(self.Txns, hashes[i])
            uxo = append(uxo, coin.CreateExpectedUnspents(t.Txn)...)
        }
    }
    uxhashes := make([]coin.SHA256, len(uxo))
    for i, _ := range uxo {
        uxhashes[i] = uxo[i].Hash()
    }
    self.Unspent.DelMultiple(uxhashes)
}

// Removes confirmed txns from the pool
func (self *UnconfirmedTxnPool) RemoveTransactions(bc *coin.Blockchain,
    txns coin.Transactions) {
    toRemove := make([]coin.SHA256, len(txns))
    for i, _ := range txns {
        toRemove[i] = txns[i].Hash()
    }
    self.removeTxns(bc, toRemove)
}

// Checks all unconfirmed txns against the blockchain. maxAge is how long
// we'll hold a txn regardless of whether it has been invalidated.
// checkPeriod is how often we check the txn against the blockchain.
func (self *UnconfirmedTxnPool) Refresh(bc *coin.Blockchain,
    checkPeriod, maxAge time.Duration) {
    now := util.Now()
    toRemove := make([]coin.SHA256, 0)
    for k, t := range self.Txns {
        if now.Sub(t.Received) >= maxAge {
            toRemove = append(toRemove, k)
        } else if now.Sub(t.Checked) >= checkPeriod {
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
func (self *UnconfirmedTxnPool) GetOldOwnedTransactions(ago time.Duration) []UnconfirmedTxn {
    txns := make([]UnconfirmedTxn, 0)
    now := util.Now()
    for _, tx := range self.Txns {
        // TODO -- don't record IsOurSpend/IsOurReceive and do lookup each time?
        // Slower but more correct
        if (tx.IsOurSpend || tx.IsOurReceive) && now.Sub(tx.Announced) > ago {
            txns = append(txns, tx)
        }
    }
    return txns
}

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

// Returns all unconfirmed coin.UxOut spends for addresses
// Looks at all inputs for unconfirmed txns, gets their source UxOut from the
// blockchain's unspent pool, and returns as coin.AddressUxOuts
// TODO -- optimize or cache
func (self *UnconfirmedTxnPool) SpendsForAddresses(bcUnspent *coin.UnspentPool,
    a map[coin.Address]byte) coin.AddressUxOuts {
    auxs := make(coin.AddressUxOuts, len(a))
    for _, utx := range self.Txns {
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
