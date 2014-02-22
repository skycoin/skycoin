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
    // SimpleWallet entries are shared as a form of identification, the secret key
    // is not required
    // TODO -- fix lib/base58 to not panic on invalid input -- should
    // return error, so we can detect a broken wallet.
    if w.Address == "" {
        log.Panic("ReadableWalletEntry has no Address")
    }
    var s coin.SecKey
    if w.Secret != "" {
        s = coin.MustSecKeyFromHex(w.Secret)
    }
    return WalletEntry{
        Address: coin.MustDecodeBase58Address(w.Address),
        Public:  coin.MustPubKeyFromHex(w.Public),
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

// Loads a WalletEntry from filename but panics is unable to load or contents
// are invalid
func MustLoadWalletEntry(filename string) WalletEntry {
    keys, err := LoadWalletEntry(filename)
    if err != nil {
        log.Panicf("Failed to load wallet entry: %v", err)
    }
    if err := keys.Verify(); err != nil {
        log.Panicf("Invalid wallet entry: %v", err)
    }
    return keys
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

func NewBalanceFromUxOut(headTime uint64, ux *coin.UxOut) Balance {
    return Balance{
        Coins: ux.Body.Coins,
        Hours: ux.CoinHours(headTime),
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
    if other.Coins > self.Coins || other.Hours > self.Hours {
        log.Panic("Cannot subtract balances, second balance is too large")
    }
    return Balance{
        Coins: self.Coins - other.Coins,
        Hours: self.Hours - other.Hours,
    }
}

func (self Balance) Equals(other Balance) bool {
    return self.Coins == other.Coins && self.Hours == other.Hours
}

func (self Balance) IsZero() bool {
    return self.Coins == 0 && self.Hours == 0
}

// Wallet interface, to support multiple implementations
type Wallet interface {
    // Returns all entries
    GetEntries() map[coin.Address]WalletEntry
    // Returns all addresses stored in the wallet as a set
    GetAddressSet() map[coin.Address]byte
    // Returns all addresses stored in the wallet as array
    GetAddresses() []coin.Address
    // Adds an entry that was created externally. Should return error if
    // entry is not valid or already existed.
    AddEntry(entry WalletEntry) error
    // Creates and adds a new entry to the wallet
    CreateEntry() WalletEntry
    // Returns a wallet entry by address and whether it exists
    GetEntry(addr coin.Address) (WalletEntry, bool)
    // Returns the number of entries in the wallet
    NumEntries() int
    // Adds entries to the wallet up to n, if number of entries is less than n
    Populate(max int)
    // Saves wallet to filename
    Save(filename string) error
    // Loads wallet from filename
    Load(filename string) error
    // Converts Wallet to a ReadableWallet, i.e. one that is json serializable
    ToReadable() *ReadableWallet
}

// Simplest wallet implementation
type SimpleWallet struct {
    Entries map[coin.Address]WalletEntry
}

func NewSimpleWallet() *SimpleWallet {
    return &SimpleWallet{
        Entries: make(map[coin.Address]WalletEntry),
    }
}

func LoadSimpleWallet(filename string) (*SimpleWallet, error) {
    w := NewSimpleWallet()
    return w, w.Load(filename)
}

func NewSimpleWalletFromReadable(r *ReadableWallet) *SimpleWallet {
    entries := make(map[coin.Address]WalletEntry, len(r.Entries))
    for _, re := range r.Entries {
        we := WalletEntryFromReadable(&re)
        if err := we.Verify(); err != nil {
            log.Panicf("Invalid wallet entry loaded. Address: %s", re.Address)
        }
        entries[we.Address] = we
    }
    return &SimpleWallet{
        Entries: entries,
    }
}

func (self *SimpleWallet) NumEntries() int {
    return len(self.Entries)
}

func (self *SimpleWallet) GetEntries() map[coin.Address]WalletEntry {
    return self.Entries
}

func (self *SimpleWallet) GetAddressSet() map[coin.Address]byte {
    m := make(map[coin.Address]byte, len(self.Entries))
    for a, _ := range self.Entries {
        m[a] = byte(1)
    }
    return m
}

// Creates a WalletEntry
func (self *SimpleWallet) CreateEntry() WalletEntry {
    e := NewWalletEntry()
    if err := self.AddEntry(e); err != nil {
        log.Panic("Somehow, we managed to create a bad entry: %v", err)
    }
    return e
}

// Creates new WalletEntries to fill the wallet up to n.  No WalletEntries
// are created if the SimpleWallet already contains n or more entries.
func (self *SimpleWallet) Populate(n int) {
    for i := len(self.Entries); i < n; i++ {
        self.CreateEntry()
    }
}

// Returns all coin.Addresses in this SimpleWallet
func (self *SimpleWallet) GetAddresses() []coin.Address {
    addrs := make([]coin.Address, len(self.Entries))
    i := 0
    for a, _ := range self.Entries {
        addrs[i] = a
        i++
    }
    return addrs
}

// Returns the WalletEntry for a coin.Address
func (self *SimpleWallet) GetEntry(a coin.Address) (WalletEntry, bool) {
    we, exists := self.Entries[a]
    return we, exists
}

// Adds a WalletEntry to the wallet. Returns an error if the coin.Address is
// already in the wallet
func (self *SimpleWallet) AddEntry(e WalletEntry) error {
    if err := e.Verify(); err != nil {
        return err
    }
    _, exists := self.Entries[e.Address]
    if exists {
        return fmt.Errorf("SimpleWallet entry already exists for address %s",
            e.Address.String())
    } else {
        self.Entries[e.Address] = e
        return nil
    }
}

// Saves to filename
func (self *SimpleWallet) Save(filename string) error {
    r := NewReadableWallet(self)
    return r.Save(filename)
}

// Loads from filename
func (self *SimpleWallet) Load(filename string) error {
    r := &ReadableWallet{}
    if err := r.Load(filename); err != nil {
        return err
    }
    *self = *(NewSimpleWalletFromReadable(r))
    return nil
}

func (self *SimpleWallet) ToReadable() *ReadableWallet {
    return NewReadableWallet(self)
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
    pubkey := coin.MustPubKeyFromHex(pub)
    addr := coin.AddressFromPubKey(pubkey)
    return ReadableWalletEntry{
        Address: addr.String(),
        Public:  pub,
    }
}

func (self *ReadableWalletEntry) Save(filename string) error {
    return util.SaveJSONSafe(filename, self, 0600)
}

// Used for [de]serialization of the SimpleWallet
type ReadableWallet struct {
    Entries []ReadableWalletEntry `json:"entries"`
}

// Converts a SimpleWallet to a ReadableWallet
func NewReadableWallet(w *SimpleWallet) *ReadableWallet {
    readable := make([]ReadableWalletEntry, len(w.Entries))
    i := 0
    for _, e := range w.Entries {
        readable[i] = NewReadableWalletEntry(&e)
        i++
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
