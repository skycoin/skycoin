package visor

import (
    "errors"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "log"
    "sort"
)

// Sorts Balances with coins ascending, and hours ascending if coins equal
type CoinsAscending struct {
    Unspents coin.AddressUnspents
    HeadTime uint64
}

func (self CoinsAscending) Len() int {
    return len(self.Unspents)
}

func (self CoinsAscending) Swap(i, j int) {
    self.Unspents[i], self.Unspents[j] = self.Unspents[j], self.Unspents[i]
}

func (self CoinsAscending) Less(i, j int) bool {
    c := self.Unspents[i].Unspent.Body.Coins
    d := self.Unspents[j].Unspent.Body.Coins
    if c == d {
        c = self.Unspents[i].Unspent.CoinHours(self.HeadTime)
        d = self.Unspents[j].Unspent.CoinHours(self.HeadTime)
    }
    return c < d
}

// Sorts AddressUxOuts with hours descending, and coins descending if equal
type HoursDescending struct {
    Unspents coin.AddressUnspents
    HeadTime uint64
}

func (self HoursDescending) Len() int {
    return len(self.Unspents)
}

func (self HoursDescending) Swap(i, j int) {
    self.Unspents[i], self.Unspents[j] = self.Unspents[j], self.Unspents[i]
}

func (self HoursDescending) Less(i, j int) bool {
    c := self.Unspents[i].Unspent.CoinHours(self.HeadTime)
    d := self.Unspents[j].Unspent.CoinHours(self.HeadTime)
    if c == d {
        c = self.Unspents[i].Unspent.Body.Coins
        d = self.Unspents[j].Unspent.Body.Coins
    }
    return c > d
}

// Removes any balances that are 0 or not multiples of 1e6.
// Note: transactions with outputs having those values are rejected, there
// shouldn't be any.
func removePartialCoins(ix coin.AddressUxOuts) coin.AddressUxOuts {
    ox := make(coin.AddressUxOuts, len(ix))
    for a, uxs := range ix {
        // Disallowed unspents should be nonexistent initially, so its fine
        // preallocate everything.
        oxs := make(coin.UxArray, 0, len(uxs))
        for _, ux := range uxs {
            if ux.Body.Coins != 0 && ux.Body.Coins%1e6 == 0 {
                oxs = append(oxs, ux)
            } else {
                logger.Warning("Found unspent with invalid coins: %v", ux)
            }
        }
        ox[a] = oxs
    }
    return ox
}

// Returns a list of coin.AddressUnspents to be used for txn construction.
// Note: amt should include the fee.  auxs should not include unconfirmed
// spends
// Goals:
//   1. Use the least number of unspents
//   2. Preserve coin hours, i.e. always change with at least 1e6 coins if
//      hours need to be returned
func createSpends(headTime uint64, auxs coin.AddressUxOuts,
    amt Balance) (coin.AddressUnspents, error) {
    if amt.IsZero() {
        return nil, errors.New("Zero spend amount")
    }
    if amt.Coins == 0 || amt.Coins%1e6 != 0 {
        return nil, errors.New("Spends must be 1e6 multiple")
    }

    // 1. Remove all balances that are not 1e6 multiples, we can't spend them
    auxs = removePartialCoins(auxs)

    // 2. Sort balances with coins,hours ascending.
    uxs := auxs.Flatten()
    asc := CoinsAscending{
        Unspents: uxs,
        HeadTime: headTime,
    }
    sort.Sort(asc)
    uxs = asc.Unspents

    // 3. For each balance, add coins + hours towards amt.
    //      If hours are not exactly satisfied, we amt to spend from one
    //      more address so that it can receive hours as change, due to the
    //      1e6 restriction
    spending := make(coin.AddressUnspents, 0)
    have := Balance{0, 0}
    for i, _ := range uxs {
        if have.Coins > amt.Coins ||
            (have.Coins == amt.Coins && have.Hours == amt.Hours) {
            break
        }
        have.Coins += uxs[i].Unspent.Body.Coins
        have.Hours += uxs[i].Unspent.CoinHours(headTime)
        spending = append(spending, uxs[i])
    }

    // 4. If coins cannot be met, fail
    if have.Coins < amt.Coins {
        return nil, errors.New("Not enough coins")
    }

    // 5. If hours are not met, sort remaining balance hours descending.
    uxs = uxs.Sub(spending)
    dsc := HoursDescending{
        Unspents: uxs,
        HeadTime: headTime,
    }
    sort.Sort(dsc)
    uxs = asc.Unspents

    // 6. For each balance, add until hours are met.
    for i, _ := range uxs {
        if have.Hours >= amt.Hours {
            break
        }
        have.Coins += uxs[i].Unspent.Body.Coins
        have.Hours += uxs[i].Unspent.CoinHours(headTime)
        spending = append(spending, uxs[i])
    }

    // 7. If hours are not met, fail
    if have.Hours < amt.Hours {
        return nil, errors.New("Not enough hours")
    }

    return spending, nil
}

// Creates a Transaction spending coins and hours from our coins
func CreateSpendingTransaction(wallet Wallet, unconfirmed *UnconfirmedTxnPool,
    unspent *coin.UnspentPool, headTime uint64, amt Balance, fee uint64,
    dest coin.Address) (coin.Transaction, error) {
    txn := coin.Transaction{}
    if amt.Coins == 0 {
        return txn, errors.New("Zero spend amount")
    }
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
        entry, exists := wallet.GetEntry(au.Address)
        if !exists {
            log.Panic("On second thought, the wallet entry does not exist")
        }
        txn.PushInput(au.Unspent.Hash())
        toSign[i] = entry.Secret
        spending.Coins += au.Unspent.Body.Coins
        spending.Hours += au.Unspent.CoinHours(headTime)
    }

    change := spending.Sub(need)
    // TODO -- send change to a new address
    changeAddr := spends[0].Address
    if change.Coins == 0 {
        if change.Hours > fee {
            msg := ("Have enough coins, but not enough to send coin hours change " +
                "back. Would spend %d more hours than requested.")
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
