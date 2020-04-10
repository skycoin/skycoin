package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/file"
	"github.com/SkycoinProject/skycoin/src/wallet/entry"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/SkycoinProject/skycoin/src/wallet/secrets"
)

// DeterministicWallet manages keys using the original Skycoin deterministic
// keypair generator method.
// With this generator, a single chain of addresses is created, each one dependent
// on the previous.
type DeterministicWallet struct {
	meta.Meta
	Entries entry.Entries
}

// newDeterministicWallet creates a DeterministicWallet
func newDeterministicWallet(meta meta.Meta) (*DeterministicWallet, error) { //nolint:unparam
	return &DeterministicWallet{
		Meta: meta,
	}, nil
}

// PackSecrets copies data from decrypted wallets into the secrets container
func (w *DeterministicWallet) PackSecrets(ss secrets.Secrets) {
	ss.Set(secrets.SecretSeed, w.Seed())
	ss.Set(secrets.SecretLastSeed, w.LastSeed())

	// Saves entry secret keys in secrets
	for _, e := range w.Entries {
		ss.Set(e.Address.String(), e.Secret.Hex())
	}
}

// UnpackSecrets copies data from decrypted secrets into the wallet
func (w *DeterministicWallet) UnpackSecrets(ss secrets.Secrets) error {
	seed, ok := ss.Get(secrets.SecretSeed)
	if !ok {
		return errors.New("seed doesn't exist in secrets")
	}
	w.SetSeed(seed)

	lastSeed, ok := ss.Get(secrets.SecretLastSeed)
	if !ok {
		return errors.New("lastSeed doesn't exist in secrets")
	}
	w.SetLastSeed(lastSeed)

	return w.Entries.UnpackSecretKeys(ss)
}

// Clone clones the wallet a new wallet object
func (w *DeterministicWallet) Clone() Wallet {
	return &DeterministicWallet{
		Meta:    w.Meta.Clone(),
		Entries: w.Entries.Clone(),
	}
}

// CopyFrom copies the src wallet to w
func (w *DeterministicWallet) CopyFrom(src Wallet) {
	w.Meta = src.(*DeterministicWallet).Meta.Clone()
	w.Entries = src.(*DeterministicWallet).Entries.Clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *DeterministicWallet) CopyFromRef(src Wallet) {
	*w = *(src.(*DeterministicWallet))
}

// Erase wipes secret fields in wallet
func (w *DeterministicWallet) Erase() {
	w.Meta.EraseSeeds()
	w.Entries.Erase()
}

// ToReadable converts the wallet to its readable (serializable) format
func (w *DeterministicWallet) ToReadable() Readable {
	return NewReadableDeterministicWallet(w)
}

// Validate validates the wallet
func (w *DeterministicWallet) Validate() error {
	if err := w.Meta.Validate(); err != nil {
		return err
	}

	walletType := w.Meta.Type()
	if !IsValidWalletType(walletType) {
		return ErrInvalidWalletType
	}

	if !w.IsEncrypted() {
		if s := w.Seed(); s == "" {
			return errors.New("seed missing in unencrypted deterministic wallet")
		}

		if s := w.LastSeed(); s == "" {
			return errors.New("lastSeed missing in unencrypted deterministic wallet")
		}
	}
	return nil
}

// GetAddresses returns all addresses in wallet
func (w *DeterministicWallet) GetAddresses() []cipher.Addresser {
	return w.Entries.GetAddresses()
}

// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
func (w *DeterministicWallet) GetSkycoinAddresses() ([]cipher.Address, error) {
	if w.Meta.Coin() != meta.CoinTypeSkycoin {
		return nil, errors.New("DeterministicWallet coin type is not skycoin")
	}

	return w.Entries.GetSkycoinAddresses(), nil
}

// GetEntries returns a copy of all entries held by the wallet
func (w *DeterministicWallet) GetEntries() entry.Entries {
	return w.Entries.Clone()
}

// EntriesLen returns the number of entries in the wallet
func (w *DeterministicWallet) EntriesLen() int {
	return len(w.Entries)
}

// GetEntryAt returns entry at a given index in the entries array
func (w *DeterministicWallet) GetEntryAt(i int) entry.Entry {
	return w.Entries[i]
}

// GetEntry returns entry of given address
func (w *DeterministicWallet) GetEntry(a cipher.Address) (entry.Entry, bool) {
	return w.Entries.Get(a)
}

// HasEntry returns true if the wallet has an Entry with a given cipher.Address.
func (w *DeterministicWallet) HasEntry(a cipher.Address) bool {
	return w.Entries.Has(a)
}

// GenerateAddresses generates addresses
func (w *DeterministicWallet) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	if w.Meta.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	if num == 0 {
		return nil, nil
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

	w.Meta.SetLastSeed(hex.EncodeToString(seed))

	addrs := make([]cipher.Addresser, len(seckeys))
	makeAddress := AddressConstructor(w.Meta)
	for i, s := range seckeys {
		p := cipher.MustPubKeyFromSecKey(s)
		a := makeAddress(p)
		addrs[i] = a
		w.Entries = append(w.Entries, entry.Entry{
			Address: a,
			Secret:  s,
			Public:  p,
		})
	}
	return addrs, nil
}

// GenerateSkycoinAddresses generates Skycoin addresses. If the wallet's coin type is not Skycoin, returns an error
func (w *DeterministicWallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
	if w.Meta.Coin() != meta.CoinTypeSkycoin {
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
	w.Entries = entry.Entries{}
	w.Meta.SetLastSeed(w.Meta.Seed())
}

// ScanAddresses scans ahead N addresses, truncating up to the highest address with any transaction history.
func (w *DeterministicWallet) ScanAddresses(scanN uint64, tf TransactionsFinder) error {
	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	if scanN == 0 {
		return nil
	}

	w2 := w.Clone().(*DeterministicWallet)

	nExistingAddrs := uint64(len(w2.Entries))

	// Generate the addresses to scan
	addrs, err := w2.GenerateSkycoinAddresses(scanN)
	if err != nil {
		return err
	}

	// Find if these addresses had any activity
	active, err := tf.AddressesActivity(addrs)
	if err != nil {
		return err
	}

	// Check activity from the last one until we find the address that has activity
	var keepNum uint64
	for i := len(active) - 1; i >= 0; i-- {
		if active[i] {
			keepNum = uint64(i + 1)
			break
		}
	}

	// Regenerate addresses up to nExistingAddrs + nAddAddrs.
	// This is necessary to keep the lastSeed updated.
	w2.reset()
	if _, err := w2.GenerateSkycoinAddresses(nExistingAddrs + keepNum); err != nil {
		return err
	}

	*w = *w2

	return nil
}

// Fingerprint returns a unique ID fingerprint for this wallet, composed of its initial address
// and wallet type
func (w *DeterministicWallet) Fingerprint() string {
	addr := ""
	if len(w.Entries) == 0 {
		if !w.IsEncrypted() {
			_, pk, _ := cipher.MustDeterministicKeyPairIterator([]byte(w.Meta.Seed()))
			addr = AddressConstructor(w.Meta)(pk).String()
		}
	} else {
		addr = w.Entries[0].Address.String()
	}
	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

// ReadableDeterministicWallet used for [de]serialization of a deterministic wallet
type ReadableDeterministicWallet struct {
	meta.Meta       `json:"meta"`
	ReadableEntries `json:"entries"`
}

// LoadReadableDeterministicWallet loads a deterministic wallet from disk
func LoadReadableDeterministicWallet(wltFile string) (*ReadableDeterministicWallet, error) {
	var rw ReadableDeterministicWallet
	if err := file.LoadJSON(wltFile, &rw); err != nil {
		return nil, err
	}
	if rw.Type() != WalletTypeDeterministic {
		return nil, ErrInvalidWalletType
	}
	return &rw, nil
}

// NewReadableDeterministicWallet creates readable wallet
func NewReadableDeterministicWallet(w *DeterministicWallet) *ReadableDeterministicWallet {
	return &ReadableDeterministicWallet{
		Meta:            w.Meta.Clone(),
		ReadableEntries: newReadableEntries(w.Entries, w.Meta.Coin(), w.Meta.Type()),
	}
}

// ToWallet convert readable wallet to Wallet
func (rw *ReadableDeterministicWallet) ToWallet() (Wallet, error) {
	w := &DeterministicWallet{
		Meta: rw.Meta.Clone(),
	}

	if err := w.Validate(); err != nil {
		err := fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
		logger.WithError(err).Error("ReadableDeterministicWallet.ToWallet Validate failed")
		return nil, err
	}

	ets, err := rw.ReadableEntries.toWalletEntries(w.Meta.Coin(), w.Meta.Type(), w.Meta.IsEncrypted())
	if err != nil {
		logger.WithError(err).Error("ReadableDeterministicWallet.ToWallet toWalletEntries failed")
		return nil, err
	}

	w.Entries = ets

	return w, nil
}
