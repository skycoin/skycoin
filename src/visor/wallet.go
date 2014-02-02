package visor

import (
    "errors"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "log"
)

type WalletEntry struct {
    Address coin.Address
    Public  coin.PubKey
    Secret  coin.SecKey
}

func NewWalletEntry() WalletEntry {
    pub, sec := coin.GenerateKeyPair()
    return WalletEntry{
        Address: coin.AddressFromPubKey(pub),
        Public:  pub,
        Secret:  sec,
    }
}

func WalletEntryFromReadable(w *ReadableWalletEntry) WalletEntry {
    // Wallet entries are shared as a form of identification, the secret key
    // is not required
    // TODO -- fix lib/base58 to not panic on invalid input -- should
    // return error, so we can detect a broken wallet.
    if w.Address == "" {
        log.Panic("ReadableWalletEntry has no Address")
    }
    var s coin.SecKey
    if w.Secret != "" {
        s = coin.SecKeyFromHex(w.Secret)
    }
    return WalletEntry{
        Address: coin.MustDecodeBase58Address(w.Address),
        Public:  coin.PubKeyFromHex(w.Public),
        Secret:  s,
    }
}

// Checks that the public key is derivable from the secret key if present,
// and that the public key is associated with the address
func (self *WalletEntry) Verify(isMaster bool) error {
    var emptySecret coin.SecKey
    if self.Secret == emptySecret {
        if isMaster {
            return errors.New("WalletEntry is master, but has no secret key")
        }
    } else {
        if coin.PubKeyFromSecKey(self.Secret) != self.Public {
            return errors.New("Invalid public key for secret key")
        }
    }
    return self.Address.Verify(self.Public)
}

type ReadableWalletEntry struct {
    Address string `json:"address"`
    Public  string `json:"public_key"`
    Secret  string `json:"secret_key"`
}

func NewReadableWalletEntry(w *WalletEntry) ReadableWalletEntry {
    return ReadableWalletEntry{
        Address: w.Address.String(),
        Public:  w.Public.Hex(),
        Secret:  w.Secret.Hex(),
    }
}

// Loads a WalletEntry from filename, where the file contains a
// ReadableWalletEntry
func LoadWalletEntry(filename string) (WalletEntry, error) {
    w := &ReadableWalletEntry{}
    err := util.LoadJSON(filename, w)
    if err != nil {
        return WalletEntry{}, err
    } else {
        return WalletEntryFromReadable(w), nil
    }
}

type Balance struct {
    Coins uint64 `json:"coins"`
    Hours uint64 `json:"hours"`
}

func NewBalance(coins, hours uint64) Balance {
    return Balance{
        Coins: coins,
        Hours: hours,
    }
}

func (self Balance) Add(other Balance) Balance {
    return Balance{
        Coins: self.Coins + other.Coins,
        Hours: self.Hours + other.Hours,
    }
}

// Subtracts other from self and returns the new Balance.  Will panic if
// other is greater than balance, because Coins and Hours are unsigned.
func (self Balance) Sub(other Balance) Balance {
    if other.GreaterThan(self) {
        log.Panic("Cannot subtract balances, second balance is too large")
    }
    return Balance{
        Coins: self.Coins - other.Coins,
        Hours: self.Hours - other.Hours,
    }
}

func (self Balance) GreaterThan(other Balance) bool {
    return self.Coins > other.Coins && self.Hours > other.Hours
}

func (self Balance) GreaterThanOrEqual(other Balance) bool {
    return self.Coins >= other.Coins && self.Hours >= other.Hours
}

func (self Balance) IsZero() bool {
    return self.Coins == 0 && self.Hours == 0
}

// Simplest wallet implementation -- array of addres, keypairs
type Wallet struct {
    Entries []WalletEntry
    // Lookup table pointing from coin.Address to position in Entries
    addressLookup map[coin.Address]int
}

func NewWallet() *Wallet {
    return &Wallet{
        Entries:       make([]WalletEntry, 0),
        addressLookup: make(map[coin.Address]int),
    }
}

// Creates new WalletEntries to fill the wallet up to n.  No WalletEntries
// are created if the Wallet already contains n or more entries.
func (self *Wallet) Populate(n int) {
    for i := len(self.Entries); i < n; i++ {
        e := NewWalletEntry()
        self.Entries = append(self.Entries, e)
        self.addressLookup[e.Address] = len(self.Entries) - 1
    }
}

func NewWalletFromReadable(r *ReadableWallet) *Wallet {
    entries := make([]WalletEntry, 0, len(r.Entries))
    for _, re := range r.Entries {
        entries = append(entries, WalletEntryFromReadable(&re))
    }
    lookup := make(map[coin.Address]int, len(entries))
    for i, e := range entries {
        lookup[e.Address] = i
    }
    return &Wallet{
        Entries:       entries,
        addressLookup: lookup,
    }
}

// Returns all coin.Addresses in this Wallet
func (self *Wallet) GetAddresses() []coin.Address {
    addrs := make([]coin.Address, 0, len(self.Entries))
    for a, _ := range self.addressLookup {
        addrs = append(addrs, a)
    }
    return addrs
}

// Returns the WalletEntry for a coin.Address
func (self *Wallet) GetEntry(a coin.Address) (WalletEntry, bool) {
    i, exists := self.addressLookup[a]
    if !exists {
        return WalletEntry{}, false
    } else {
        return self.Entries[i], true
    }
}

// Returns the Balance for a single Address
func (self *Wallet) Balance(unspent *coin.UnspentPool, prevTime uint64,
    a coin.Address) Balance {
    b := NewBalance(0, 0)
    uxs := unspent.AllForAddress(a)
    for _, ux := range uxs {
        b = b.Add(NewBalance(ux.Body.Coins, ux.CoinHours(prevTime)))
    }
    return b
}

// Returns the sum of all Balances for each Address in the wallet
func (self *Wallet) TotalBalance(unspent *coin.UnspentPool,
    prevTime uint64) Balance {
    b := NewBalance(0, 0)
    addrs := self.GetAddresses()
    auxs := unspent.AllForAddresses(addrs)
    for _, uxs := range auxs {
        for _, ux := range uxs {
            b = b.Add(NewBalance(ux.Body.Coins, ux.CoinHours(prevTime)))
        }
    }
    return b
}

// Saves to filename
func (self *Wallet) Save(filename string) error {
    r := NewReadableWallet(self)
    return r.Save(filename)
}

// Loads from filename
func (self *Wallet) Load(filename string) error {
    r := &ReadableWallet{}
    if err := r.Load(filename); err != nil {
        return err
    }
    *self = *(NewWalletFromReadable(r))
    return nil
}

// Used for [de]serialization of the Wallet
type ReadableWallet struct {
    Entries []ReadableWalletEntry `json:"entries"`
}

// Converts a Wallet to a ReadableWallet
func NewReadableWallet(w *Wallet) *ReadableWallet {
    readable := make([]ReadableWalletEntry, 0, len(w.Entries))
    for _, e := range w.Entries {
        readable = append(readable, NewReadableWalletEntry(&e))
    }
    return &ReadableWallet{
        Entries: readable,
    }
}

// Saves to filename
func (self *ReadableWallet) Save(filename string) error {
    return util.SaveJSON(filename, self, 0600)
}

// Loads from filename
func (self *ReadableWallet) Load(filename string) error {
    return util.LoadJSON(filename, self)
}
