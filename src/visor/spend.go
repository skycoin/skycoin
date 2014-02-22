package visor

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "log"
    "sort"
)

/*

Sort unspents oldest to newest

Keep adding until either exact amount (coins+hours), or if hours exceeded,
also exceed coins by at least 1e6


*/

// Sorts a UxArray oldest to newest.
type OldestUxOut coin.UxArray

func (self OldestUxOut) Len() int      { return len(self) }
func (self OldestUxOut) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
func (self OldestUxOut) Less(i, j int) bool {
    a := self[i].Head.BkSeq
    b := self[j].Head.BkSeq
    // Use hash to break ties
    if a == b {
        ih := self[i].Hash()
        jh := self[j].Hash()
        cmp := bytes.Compare(ih[:], jh[:])
        if cmp == 0 {
            log.Panic("Duplicate UxOut when sorting")
        }
        return cmp < 0
    }
    return a < b
}

func createSpends(headTime uint64, auxs coin.AddressUxOuts,
    amt Balance) (coin.UxArray, error) {
    if amt.Coins == 0 {
        return nil, errors.New("Zero spend amount")
    }
    if amt.Coins%1e6 != 0 {
        return nil, errors.New("Coins must be multiple of 1e6")
    }

    uxs := OldestUxOut(auxs.Flatten())
    sort.Sort(uxs)

    have := Balance{0, 0}
    spending := make(coin.UxArray, 0)
    for i, _ := range uxs {
        if have.Coins > amt.Coins && have.Hours >= amt.Hours {
            break
        }
        // If we have the exact amount of both, we don't need any extra coins
        // for change
        if have.Coins == amt.Coins && have.Hours == amt.Hours {
            break
        }
        b := NewBalanceFromUxOut(headTime, &uxs[i])
        if b.Coins == 0 || b.Coins%1e6 != 0 {
            logger.Error("UxOut coins are 0 or 1e6, can't spend")
            continue
        }
        have = have.Add(b)
        spending = append(spending, uxs[i])
    }
    if amt.Coins > have.Coins {
        return nil, errors.New("Not enough coins")
    }
    if amt.Hours > have.Hours {
        return nil, errors.New("Not enough hours")
    }
    return spending, nil
}

// Creates a Transaction spending coins and hours from our coins
func CreateSpendingTransaction(wallet Wallet, unconfirmed *UnconfirmedTxnPool,
    unspent *coin.UnspentPool, headTime uint64, amt Balance, fee uint64,
    dest coin.Address) (coin.Transaction, error) {
    txn := coin.Transaction{}
    need := amt
    need.Hours += fee
    addrs := wallet.GetAddresses()
    auxs := unspent.AllForAddresses(addrs)
    // Subtract pending spends from available
    puxs := unconfirmed.SpendsForAddresses(unspent, addrs)
    auxs = auxs.Sub(puxs)

    spends, err := createSpends(headTime, auxs, need)
    if err != nil {
        return txn, err
    }
    toSign := make([]coin.SecKey, len(spends))
    spending := Balance{0, 0}
    for i, au := range spends {
        entry, exists := wallet.GetEntry(au.Body.Address)
        if !exists {
            log.Panic("On second thought, the wallet entry does not exist")
        }
        txn.PushInput(au.Hash())
        toSign[i] = entry.Secret
        spending.Coins += au.Body.Coins
        spending.Hours += au.CoinHours(headTime)
    }

    change := spending.Sub(need)
    // TODO -- send change to a new address
    changeAddr := spends[0].Body.Address
    if change.Coins == 0 {
        if change.Hours > fee {
            msg := ("Have enough coins, but not enough to send coin hours " +
                "change back. Would spend %d more hours than requested.")
            return txn, fmt.Errorf(msg, change.Hours-fee)
        }
    } else {
        logger.Info("Sending change to %s: %d, %d", changeAddr.String(),
            change.Coins, change.Hours)
        txn.PushOutput(changeAddr, change.Coins, change.Hours)
    }

    logger.Info("Sending money to %s: %d, %d", dest.String(), amt.Coins,
        amt.Hours)
    txn.PushOutput(dest, amt.Coins, amt.Hours)
    txn.SignInputs(toSign)
    txn.UpdateHeader()
    return txn, nil
}
