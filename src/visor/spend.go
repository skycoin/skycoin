package visor

// Returns the UxArray as a hash to byte map to be used as a set.  The byte's
// value should be ignored, although it will be 1.  Should only be used for
// membership detection.
func (self UxArray) Set() map[SHA256]byte {
    m := make(map[SHA256]byte, len(self))
    for i, _ := range self {
        m[self[i].Hash()] = byte(1)
    }
    return m
}

// Returns a new UxArray with elements in other removed from self
func (self UxArray) Sub(other UxArray) UxArray {
    // Assume everything will be removed for initial allocation
    uxa := make(UxArray, 0, len(self)-len(other))
    selfm := self.Set()
    for i, _ := range other {
        if ux, ok := selfm[other[i].Hash()]; !ok {
            uxa = append(uxa, ux)
        }
    }
    return uxa
}

// Returns a new set of unspents, with unspents found in other removed.
// No address's unspent set will be empty
func (self AddressUxOuts) Sub(other coin.AddressUxOuts) AddressUxOuts {
    ox := make(coin.AddressUxOuts, len(ix))
    for a, uxs := range other {
        if suxs, ok := self[a]; ok {
            ouxs := suxs.Sub(uxs)
            if len(ouxs) > 0 {
                ox[a] = ouxs
            }
        }
    }
    return ox
}

// Converts an AddressUxOuts map to an array
func (self AddressUxOuts) Flatten() AddressUnspents {
    oxs := make([]AddressUxOuts, 0, len(self))
    for a, uxs := range self {
        for i, _ := range uxs {
            oxs = append(oxs, AddressUnspent{
                Address: a,
                Unspent: uxs[i],
            })
        }
    }
    return oxs
}

// Sorts Balances with coins ascending, and hours ascending if coins equal
type CoinsAscending struct {
    Unspents AddressUnspents
    HeadTime uint64
}

func (self CoinsAscending) Len() int {
    return len(self.Unspents)
}

func (self CoinsAscending) Swap(i, j int) {
    self.Unspents[i], self.Unspents[j] = self.Unspents[j], self.Unspents[i]
}

func (self CoinsAscending) Less(i, j int) bool {
    c := self.Unspent[i].Body.Coins
    d := self.Unspent[j].Body.Coins
    if c == d {
        c = self.Unspent[i].CoinHours(self.HeadTime)
        d = self.Unspent[j].CoinHours(self.HeadTime)
    }
    return c < d
}

// Sorts AddressUxOuts with hours descending, and coins descending if hours equal
type HoursDescending struct {
    Unspents AddressUnspents
    HeadTime uint64
}

func (self HoursDescending) Len() int {
    return len(self.Unspents)
}

func (self HoursDescending) Swap(i, j int) {
    self.Unspents[i], self.Unspents[j] = self.Unspents[j], self.Unspents[i]
}

func (self HoursDescending) Less(i, j int) bool {
    c := self.Unspent[i].CoinHours(self.HeadTime)
    d := self.Unspent[j].CoinHours(self.HeadTime)
    if c == d {
        c = self.Unspent[i].Body.Coins
        d = self.Unspent[j].Body.Coins
    }
    return c > d
}

type AddressUnspent struct {
    Address Address
    Unspent UxOut
}

type AddressUnspents []AddressUnspent

func (self AddressUnspents) Set() map[SHA256]byte {
    m := make(map[SHA256]byte, len(self))
    for i, _ := range self {
        m[self[i].Unspent.Hash()] = self[i]
    }
    return m
}

func (self AddressUnspents) Sub(other AddressUnspents) AddressUnspents {
    m := other.Set()
    o := make(AddressUnspents, 0, len(self)-len(other))
    for i, _ := range self {
        if _, ok := m[self[i].Unspent.Hash()]; !ok {
            o = append(o, self[i])
        }
    }
    return o
}

// Removes any balances that are 0 or not multiples of 1e6.
// Note: transactions with outputs having those values are rejected, there
// shouldn't be any.
func removePartialCoins(ix coin.AddressUxOuts) coin.AddressUxOuts {
    ox := make(coin.AddressUxOuts, len(ix))
    for a, uxs := range ix {
        // Disallowed unspents should be nonexistent initially, so its fine
        // preallocate everything.
        oxs := make(UxArray, 0, len(uxs))
        for i, ux := range uxs {
            if ux.Coins == 0 || ux.Coins%1e6 != 0 {
                logger.Warning("Found unspent with invalid coins: %v", ux)
            } else {
                oxs = append(oxs, ux)
            }
        }
        ox[a] = oxs
    }
    return ox
}

// Returns a list of AddressUnspents to be used for txn construction.
// Note: amt should include the fee.  auxs should not include unconfirmed
// spends
func createSpends(headTime uint64, auxs, AddressUxOuts,
    amt Balance) (AddressUnspents, error) {
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
    uxs = asc.uxs

    // 3. For each balance, add coins + hours towards amt.
    //      If hours are not exactly satisfied, we amt to spend from one
    //      more address so that it can receive hours as change, due to the
    //      1e6 restriction
    spending := make(AddressUnspents, 0)
    have := Balance{0, 0}
    for i, _ := range uxs {
        have.Coins += uxs[i].Unspent.Body.Coins
        have.Hours += uxs[i].Unspent.CoinHours(headTime)
        spending = append(spending, uxs[i])
        if have.Coins > amt.Coins ||
            (have.Coins == amt.Coins && have.Hours == amt.Hours) {
            break
        }
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
    uxs = asc.uxs

    // 6. For each balance, add until hours are met.
    for i, _ := range uxs {
        have.Coins += uxs[i].Unspent.Coins
        have.Hours += uxs[i].Unspent.CoinHours(headTime)
        spending = append(spending, uxs[i])
        if have.Hours >= amt.Hours {
            break
        }
    }

    // 7. If hours are not met, fail
    if have.Hours < amt.Hours {
        return nil, errors.New("Not enough hours")
    }

    return spending, nil
}

// Creates a Transaction spending coins and hours from our coins
func (self *Visor) Spend(amt Balance, fee uint64,
    dest coin.Address) (coin.Transaction, error) {
    var txn coin.Transaction
    if !self.Config.CanSpend {
        return txn, errors.New("Spending disabled")
    }
    if amt.IsZero() {
        return txn, errors.New("Zero spend amount")
    }
    need := amt
    need.Hours += fee
    // TODO -- re-enable once prediction is fixed
    // We amt to keep track of only what we spent that is unconfirmed
    // And subtract those from auxs' balances
    // auxs := self.getAvailableBalances()
    addrs := self.Wallet.GetAddresses()
    auxs := self.blockchain.Unspent.AllForAddresses(addrs)
    puxs := self.PendingUnspents.AllForAddresses(addrs)
    auxs = auxs.Sub(puxs)

    headTime := self.blockchain.HeadTime()
    spends, err := createSpends(headTime, auxs, need)
    if err != nil {
        return txn, err
    }
    toSign := make([]coin.SecKey, len(spends))
    spending := Balance{0, 0}
    for i, au := range spends {
        entry, exists := self.Wallet.GetEntry(a)
        if !exists {
            log.Panic("On second thought, the wallet entry does not exist")
        }
        txn.PushInput(au.Unspent.Hash())
        toSign[i] = entry.Secret
        spending.Coins += au.Unspent.Coins
        spending.Hours += au.Unspent.CoinHours(headTime)
    }

    change := spending.Sub(need)
    change = change.Sub(Balance{Coins: 0, Hours: fee})
    // TODO -- send change to a new address
    changeAddr := spends[0].Address

    txn.PushOutput(changeAddr, change.Coins, change.Hours)
    txn.PushOutput(dest, amt.Coins, amt.Hours)
    txn.SignInputs(toSign)
    txn.UpdateHeader()
    return txn
}
