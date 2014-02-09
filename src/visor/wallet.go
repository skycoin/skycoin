package visor

import (
    "errors"
    "fmt"
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

// Loads a WalletEntry from filename, where the file contains a
// ReadableWalletEntry
func LoadWalletEntry(filename string) (WalletEntry, error) {
    w, err := LoadReadableWalletEntry(filename)
    if err != nil {
        return WalletEntry{}, err
    } else {
        return WalletEntryFromReadable(&w), nil
    }
}

// Loads a WalletEntry from filename but also panics if the entry is invalid
func MustLoadWalletEntry(filename string) (WalletEntry, error) {
    keys, err := LoadWalletEntry(filename)
    if err != nil {
        return keys, err
    }
    if err := keys.Verify(); err != nil {
        log.Panicf("Invalid wallet entry: %v", err)
    }
    return keys, nil
}

// Checks that the public key is derivable from the secret key,
// and that the public key is associated with the address
func (self *WalletEntry) Verify() error {
    if coin.PubKeyFromSecKey(self.Secret) != self.Public {
        return errors.New("Invalid public key for secret key")
    }
    return self.VerifyPublic()
}

// Checks that the public key is associated with the address
func (self *WalletEntry) VerifyPublic() error {
    if err := self.Public.Verify(); err != nil {
        return err
    } else {
        return self.Address.Verify(self.Public)
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

func NewWalletFromReadable(r *ReadableWallet) *Wallet {
    entries := make([]WalletEntry, 0, len(r.Entries))
    for _, re := range r.Entries {
        we := WalletEntryFromReadable(&re)
        if err := we.Verify(); err != nil {
            log.Panicf("Invalid wallet entry loaded. Address: %s", re.Address)
        }
        entries = append(entries, we)
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

// Creates a WalletEntry
func (self *Wallet) CreateEntry() WalletEntry {
    e := NewWalletEntry()
    if err := e.Verify(); err != nil {
        log.Panic("Creating invalid wallet entry: %v", err)
    }
    if err := self.AddEntry(e); err != nil {
        log.Panic("Somehow, we managed to create a duplicate Address: %v", err)
    }
    return e
}

// Creates new WalletEntries to fill the wallet up to n.  No WalletEntries
// are created if the Wallet already contains n or more entries.
func (self *Wallet) populate(n int) {
    for i := len(self.Entries); i < n; i++ {
        self.CreateEntry()
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

func (self *Wallet) AddEntry(e WalletEntry) error {
    _, exists := self.addressLookup[e.Address]
    if exists {
        return fmt.Errorf("Wallet entry already exists for address %s",
            e.Address.String())
    } else {
        self.Entries = append(self.Entries, e)
        self.addressLookup[e.Address] = len(self.Entries) - 1
        return nil
    }
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

func LoadReadableWalletEntry(filename string) (ReadableWalletEntry, error) {
    w := ReadableWalletEntry{}
    err := util.LoadJSON(filename, &w)
    return w, err
}

// Creates a ReadableWalletEntry given a pubkey hex string.  The Secret field
// is left empty.
func ReadableWalletEntryFromPubkey(pub string) ReadableWalletEntry {
    pubkey := coin.PubKeyFromHex(pub)
    addr := coin.AddressFromPubKey(pubkey)
    return ReadableWalletEntry{
        Address: addr.String(),
        Public:  pub,
    }
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

// Loads a ReadableWallet from disk
func LoadReadableWallet(filename string) (*ReadableWallet, error) {
    w := &ReadableWallet{}
    err := w.Load(filename)
    return w, err
}

// Saves to filename
func (self *ReadableWallet) Save(filename string) error {
    return util.SaveJSON(filename, self, 0600)
}

// Loads from filename
func (self *ReadableWallet) Load(filename string) error {
    return util.LoadJSON(filename, self)
}
