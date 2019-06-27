package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

// DeterministicWallet manages keys using the original Skycoin deterministic
// keypair generator method.
// With this generator, a single chain of addresses is created, each one dependent
// on the previous.
type DeterministicWallet struct {
	Meta
	Entries Entries
}

// NewDeterministicWallet creates a DeterministicWallet
func NewDeterministicWallet(meta Meta) *DeterministicWallet {
	return &DeterministicWallet{
		Meta: meta,
	}
}

// PackSecrets copies data from decrypted wallets into the secrets container
func (w *DeterministicWallet) PackSecrets(ss Secrets) {
	ss.set(secretSeed, w.Meta.Seed())
	ss.set(secretLastSeed, w.Meta.LastSeed())

	// Saves entry secret keys in secrets
	for _, e := range w.Entries {
		ss.set(e.Address.String(), e.Secret.Hex())
	}
}

// UnpackSecrets copies data from decrypted secrets into the wallet
func (w *DeterministicWallet) UnpackSecrets(ss Secrets) error {
	seed, ok := ss.get(secretSeed)
	if !ok {
		return errors.New("seed doesn't exist in secrets")
	}
	w.Meta.setSeed(seed)

	lastSeed, ok := ss.get(secretLastSeed)
	if !ok {
		return errors.New("lastSeed doesn't exist in secrets")
	}
	w.Meta.setLastSeed(lastSeed)

	// Gets addresses related secrets
	for i, e := range w.Entries {
		sstr, ok := ss.get(e.Address.String())
		if !ok {
			return fmt.Errorf("secret of address %s doesn't exist in secrets", e.Address)
		}
		s, err := hex.DecodeString(sstr)
		if err != nil {
			return fmt.Errorf("decode secret hex string failed: %v", err)
		}

		copy(w.Entries[i].Secret[:], s[:])
	}

	return nil
}

// Clone clones the wallet a new wallet object
func (w *DeterministicWallet) Clone() Walleter {
	return &DeterministicWallet{
		Meta:    w.Meta.clone(),
		Entries: w.Entries.clone(),
	}
}

// CopyFrom copies the src wallet to w
func (w *DeterministicWallet) CopyFrom(src Walleter) {
	w.Meta = src.(*DeterministicWallet).Meta.clone()
	w.Entries = src.(*DeterministicWallet).Entries.clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *DeterministicWallet) CopyFromRef(src Walleter) {
	*w = *(src.(*DeterministicWallet))
}

// Erase wipes secret fields in wallet
func (w *DeterministicWallet) Erase() {
	w.Meta.eraseSeeds()
	w.Entries.erase()
}

// Save saves the wallet to given dir
func (w *DeterministicWallet) Save(dir string) error {
	rw := NewReadableDeterministicWallet(w)
	return SaveReadable(rw, dir)
}

// Validate validates the wallet
func (w *DeterministicWallet) Validate() error {
	return w.Meta.validate()
}

// GetAddresses returns all addresses in wallet
func (w *DeterministicWallet) GetAddresses() []cipher.Addresser {
	return w.Entries.getAddresses()
}

// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
func (w *DeterministicWallet) GetSkycoinAddresses() ([]cipher.Address, error) {
	if w.Meta.Coin() != CoinTypeSkycoin {
		return nil, errors.New("DeterministicWallet coin type is not skycoin")
	}

	return w.Entries.getSkycoinAddresses(), nil
}

// GetEntries returns all entries held by the wallet
func (w *DeterministicWallet) GetEntries() Entries {
	return w.Entries
}

// EntriesLen returns the number of entries in the wallet
func (w *DeterministicWallet) EntriesLen() int {
	return len(w.Entries)
}

// GetEntryAt returns entry at a given index in the entries array
func (w *DeterministicWallet) GetEntryAt(i int) Entry {
	return w.Entries[i]
}

// GetEntry returns entry of given address
func (w *DeterministicWallet) GetEntry(a cipher.Address) (Entry, bool) {
	return w.Entries.get(a)
}

// HasEntry returns true if the wallet has an Entry with a given cipher.Address.
func (w *DeterministicWallet) HasEntry(a cipher.Address) bool {
	return w.Entries.has(a)
}

// GenerateAddresses generates addresses
func (w *DeterministicWallet) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	if num == 0 {
		return nil, nil
	}

	if w.Meta.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	var seckeys []cipher.SecKey
	var seed []byte
	if len(w.Entries) == 0 {
		seed, seckeys = cipher.MustGenerateDeterministicKeyPairsSeed([]byte(w.Meta.Seed()), int(num))
	} else {
		sd, err := hex.DecodeString(w.Meta.LastSeed())
		if err != nil {
			return nil, fmt.Errorf("decode hex seed failed: %v", err)
		}
		seed, seckeys = cipher.MustGenerateDeterministicKeyPairsSeed(sd, int(num))
	}

	w.Meta.setLastSeed(hex.EncodeToString(seed))

	addrs := make([]cipher.Addresser, len(seckeys))
	makeAddress := w.Meta.AddressConstructor()
	for i, s := range seckeys {
		p := cipher.MustPubKeyFromSecKey(s)
		a := makeAddress(p)
		addrs[i] = a
		w.Entries = append(w.Entries, Entry{
			Address: a,
			Secret:  s,
			Public:  p,
		})
	}
	return addrs, nil
}

// GenerateSkycoinAddresses generates Skycoin addresses. If the wallet's coin type is not Skycoin, returns an error
func (w *DeterministicWallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
	if w.Meta.Coin() != CoinTypeSkycoin {
		return nil, errors.New("GenerateSkycoinAddresses called for non-skycoin wallet")
	}

	addrs, err := w.GenerateAddresses(num)
	if err != nil {
		return nil, err
	}

	skyAddrs := make([]cipher.Address, len(addrs))
	for i, a := range addrs {
		skyAddrs[i] = a.(cipher.Address)
	}

	return skyAddrs, nil
}

// reset resets the wallet entries and move the lastSeed to origin
func (w *DeterministicWallet) reset() {
	w.Entries = Entries{}
	w.Meta.setLastSeed(w.Meta.Seed())
}

// ScanAddresses scans ahead N addresses, truncating up to the highest address with a non-zero balance.
// If any address has a nonzero balance, it rescans N more addresses from that point, until a entire
// sequence of N addresses has no balance.
func (w *DeterministicWallet) ScanAddresses(scanN uint64, tf TransactionsFinder) error {
	if w.Meta.IsEncrypted() {
		return ErrWalletEncrypted
	}

	if scanN == 0 {
		return nil
	}

	w2 := w.Clone().(*DeterministicWallet)

	nExistingAddrs := uint64(len(w2.Entries))
	nAddAddrs := uint64(0)
	n := scanN
	extraScan := uint64(0)

	for {
		// Generate the addresses to scan
		addrs, err := w2.GenerateSkycoinAddresses(n)
		if err != nil {
			return err
		}

		// Find if these addresses had any activity
		active, err := tf.AddressesActivity(addrs)
		if err != nil {
			return err
		}

		// Check balance from the last one until we find the address that has activity
		var keepNum uint64
		for i := len(active) - 1; i >= 0; i-- {
			if active[i] {
				keepNum = uint64(i + 1)
				break
			}
		}

		if keepNum == 0 {
			break
		}

		nAddAddrs += keepNum + extraScan

		// extraScan is the number of addresses with a zero balance beyond the
		// last address with a nonzero balance
		extraScan = n - keepNum

		// n is the number of addresses to scan the next iteration
		n = scanN - extraScan
	}

	// Regenerate addresses up to nExistingAddrs + nAddAddrss.
	// This is necessary to keep the lastSeed updated.
	w2.reset()
	if _, err := w2.GenerateSkycoinAddresses(nExistingAddrs + nAddAddrs); err != nil {
		return err
	}

	*w = *w2

	return nil
}

// ReadableDeterministicWallet used for [de]serialization of a deterministic wallet
type ReadableDeterministicWallet struct {
	Meta    `json:"meta"`
	Entries ReadableEntries `json:"entries"`
}

// LoadReadableDeterministicWallet loads a deterministic wallet from disk
func LoadReadableDeterministicWallet(wltFile string) (*ReadableDeterministicWallet, error) {
	var rw ReadableDeterministicWallet
	if err := file.LoadJSON(wltFile, &rw); err != nil {
		return nil, err
	}
	return &rw, nil
}

// NewReadableDeterministicWallet creates readable wallet
func NewReadableDeterministicWallet(w *DeterministicWallet) *ReadableDeterministicWallet {
	return &ReadableDeterministicWallet{
		Meta:    w.Meta.clone(),
		Entries: newReadableEntries(w.Entries, w.Meta.Coin()),
	}
}

// ToWallet convert readable wallet to Wallet
func (rw *ReadableDeterministicWallet) ToWallet() (Walleter, error) {
	w := &DeterministicWallet{
		Meta: rw.Meta.clone(),
	}

	if err := w.Validate(); err != nil {
		err := fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
		logger.WithError(err).Error("ReadableDeterministicWallet.ToWallet Validate failed")
		return nil, err
	}

	ets, err := rw.Entries.toWalletEntries(w.Meta.Coin(), w.Meta.IsEncrypted())
	if err != nil {
		logger.WithError(err).Error("ReadableDeterministicWallet.ToWallet toWalletEntries failed")
		return nil, err
	}

	w.Entries = ets

	return w, nil
}
