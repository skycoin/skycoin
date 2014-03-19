package wallet

import (
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "log"
)

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

type ReadableWalletEntries []ReadableWalletEntry

func (self ReadableWalletEntries) ToWalletEntries() WalletEntries {
    entries := make(WalletEntries, len(self))
    for _, re := range self {
        we := WalletEntryFromReadable(&re)
        if err := we.Verify(); err != nil {
            log.Panicf("Invalid wallet entry loaded. Address: %s", re.Address)
        }
        entries[we.Address] = we
    }
    return entries
}

// Used for [de]serialization of a Wallet
type ReadableWallet struct {
    Type     WalletType             `json:"name"`
    Name     string                 `json:"name"`
    Filename string                 `json:"filename"`
    Entries  ReadableWalletEntries  `json:"entries"`
    Extra    map[string]interface{} `json:"extra"`
}

// Converts a Wallet to a ReadableWallet
func NewReadableWallet(w Wallet) *ReadableWallet {
    entries := w.GetEntries()
    readable := make(ReadableWalletEntries, len(entries))
    i := 0
    for _, e := range entries {
        readable[i] = NewReadableWalletEntry(&e)
        i++
    }
    return &ReadableWallet{
        Type:     w.GetType(),
        Name:     w.GetName(),
        Filename: w.GetFilename(),
        Entries:  readable,
        Extra:    w.GetExtraSerializerData(),
    }
}

// Loads a ReadableWallet from disk
func LoadReadableWallet(filename string) (*ReadableWallet, error) {
    w := &ReadableWallet{}
    err := w.Load(filename)
    return w, err
}

func (self *ReadableWallet) ToWallet() (Wallet, error) {
    switch self.Type {
    case DeterministicWalletType:
        return NewDeterministicWalletFromReadable(self), nil
    case SimpleWalletType:
        return NewSimpleWalletFromReadable(self), nil
    default:
        return nil, fmt.Errorf("Unknown wallet type \"%s\"", self.Type)
    }
}

// Saves to filename
func (self *ReadableWallet) Save(filename string) error {
    return util.SaveJSON(filename, self, 0600)
}

// Saves to filename, but won't overwrite existing
func (self *ReadableWallet) SaveSafe(filename string) error {
    return util.SaveJSONSafe(filename, self, 0600)
}

// Loads from filename
func (self *ReadableWallet) Load(filename string) error {
    return util.LoadJSON(filename, self)
}
