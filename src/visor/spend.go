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

func calculateBurnAndChange(total, spending, fee,
    factor uint64) (uint64, uint64, error) {
    if total < fee {
        return 0, 0, errors.New("Insufficient total")
    }
    burn := uint64(0)
    if factor > 0 {
        burn = (total - fee) / (factor + 1)
    }
    others := fee + spending + burn
    if total < others {
        return 0, 0, errors.New("Insufficient total")
    }
    change := total - others
    return burn, change, nil
}

func createSpends(headTime uint64, uxa coin.UxArray,
    amt Balance, fee, burnFactor uint64) (coin.UxArray, error) {
    if amt.Coins == 0 {
        return nil, errors.New("Zero spend amount")
    }
    if amt.Coins%1e6 != 0 {
        return nil, errors.New("Coins must be multiple of 1e6")
    }

    uxs := OldestUxOut(uxa)
    sort.Sort(uxs)

    have := Balance{0, 0}
    spending := make(coin.UxArray, 0)
    for i, _ := range uxs {
        burn, _, err := calculateBurnAndChange(have.Hours, amt.Hours,
            fee, burnFactor)
        if err == nil {
            trueHours := amt.Hours + fee + burn
            // Adjust hours as a moving target as outputs change
            if have.Coins > amt.Coins && have.Hours >= trueHours {
                break
            }
            // If we have the exact amount of both, we don't need any extra coins
            // for change
            if have.Coins == amt.Coins && have.Hours == trueHours {
                break
            }
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
    if amt.Hours+fee > have.Hours {
        return nil, errors.New("Not enough hours")
    }
    if _, _, err := calculateBurnAndChange(have.Hours, amt.Hours, fee,
        burnFactor); err != nil {
        return nil, errors.New("Not enough hours to burn")
    }
    return spending, nil
}

// Creates a Transaction spending coins and hours from our coins
func CreateSpendingTransaction(wallet Wallet, unconfirmed *UnconfirmedTxnPool,
    unspent *coin.UnspentPool, headTime uint64, amt Balance,
    fee, burnFactor uint64, dest coin.Address) (coin.Transaction, error) {
    txn := coin.Transaction{}
    auxs := unspent.AllForAddresses(wallet.GetAddresses())
    // Subtract pending spends from available
    puxs := unconfirmed.SpendsForAddresses(unspent, wallet.GetAddressSet())
    auxs = auxs.Sub(puxs)

    // Determine which unspents to spend
    spends, err := createSpends(headTime, auxs.Flatten(), amt, fee, burnFactor)
    if err != nil {
        return txn, err
    }

    // Add these unspents as tx inputs
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

    // Determine how much change we get back, if any
    _, changeHours, err := calculateBurnAndChange(spending.Hours,
        amt.Hours, fee, burnFactor)
    if err != nil {
        // This should not occur, else createSpends is broken
        return txn, err
    }
    change := NewBalance(spending.Coins-amt.Coins, changeHours)
    // TODO -- send change to a new address
    changeAddr := spends[0].Body.Address
    if change.Coins == 0 {
        if change.Hours > 0 {
            msg := ("Have enough coins, but not enough to send coin hours " +
                "change back. Would spend %d more hours than requested.")
            return txn, fmt.Errorf(msg, change.Hours)
        }
    } else {
        txn.PushOutput(changeAddr, change.Coins, change.Hours)
    }

    // Finalize the the transaction
    txn.PushOutput(dest, amt.Coins, amt.Hours)
    txn.SignInputs(toSign)
    txn.UpdateHeader()
    return txn, nil
}
