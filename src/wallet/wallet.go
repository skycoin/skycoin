package wallet

import (
    "github.com/skycoin/skycoin/src/coin"
)

const WalletExt = ".wlt"

type WalletType string

const (
    SimpleWalletType        WalletType = "Simple"
    DeterministicWalletType WalletType = "Deterministic"
)

type WalletID string
type AddressSet map[coin.Address]byte

func (self AddressSet) Update(other AddressSet) AddressSet {
    for k, v := range other {
        self[k] = v
    }
    return self
}

type WalletConstructor func() Wallet

// Wallet interface, to support multiple implementations
type Wallet interface {
    // Returns all entries
    GetEntries() WalletEntries
    // Returns all addresses stored in the wallet as a set
    GetAddressSet() AddressSet
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
    // Saves wallet
    Save(dir string) error
    // Loads wallet
    Load(dir string) error
    // Sets the name of the wallet
    SetName(name string)
    GetName() string
    // Returns the wallet's unique identifier
    GetID() WalletID
    // Sets the wallets filename on disk
    SetFilename(fn string)
    GetFilename() string
    // Returns the type of the wallet (e.g. "Deterministic", "Simple")
    GetType() WalletType
    // Returns extra info to be serialized with the wallet
    GetExtraSerializerData() map[string]interface{}
}
