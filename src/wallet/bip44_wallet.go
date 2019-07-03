package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

// Bip44Wallet manages keys using the original Skycoin deterministic
// keypair generator method.
// With this generator, a single chain of addresses is created, each one dependent
// on the previous.
type Bip44Wallet struct {
	Meta
	ExternalEntries Entries
	ChangeEntries   Entries
}

// newBip44Wallet creates a Bip44Wallet
func newBip44Wallet(meta Meta) *Bip44Wallet {
	return &Bip44Wallet{
		Meta: meta,
	}
}

// PackSecrets copies data from decrypted wallets into the secrets container
func (w *Bip44Wallet) PackSecrets(ss Secrets) {
	ss.set(secretSeed, w.Meta.Seed())

	// Saves entry secret keys in secrets
	for _, e := range w.ExternalEntries {
		ss.set(e.Address.String(), e.Secret.Hex())
	}
	for _, e := range w.ChangeEntries {
		ss.set(e.Address.String(), e.Secret.Hex())
	}
}

// UnpackSecrets copies data from decrypted secrets into the wallet
func (w *Bip44Wallet) UnpackSecrets(ss Secrets) error {
	seed, ok := ss.get(secretSeed)
	if !ok {
		return errors.New("seed doesn't exist in secrets")
	}
	w.Meta.setSeed(seed)

	if err := w.ExternalEntries.unpackSecretKeys(ss); err != nil {
		return err
	}
	return w.ChangeEntries.unpackSecretKeys(ss)
}

// Clone clones the wallet a new wallet object
func (w *Bip44Wallet) Clone() Wallet {
	return &Bip44Wallet{
		Meta:            w.Meta.clone(),
		ExternalEntries: w.ExternalEntries.clone(),
		ChangeEntries:   w.ChangeEntries.clone(),
	}
}

// CopyFrom copies the src wallet to w
func (w *Bip44Wallet) CopyFrom(src Wallet) {
	w.Meta = src.(*Bip44Wallet).Meta.clone()
	w.ExternalEntries = src.(*Bip44Wallet).ExternalEntries.clone()
	w.ChangeEntries = src.(*Bip44Wallet).ChangeEntries.clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *Bip44Wallet) CopyFromRef(src Wallet) {
	*w = *(src.(*Bip44Wallet))
}

// Erase wipes secret fields in wallet
func (w *Bip44Wallet) Erase() {
	w.Meta.eraseSeeds()
	w.ExternalEntries.erase()
	w.ChangeEntries.erase()
}

// ToReadable converts the wallet to its readable (serializable) format
func (w *Bip44Wallet) ToReadable() Readable {
	return NewReadableBip44Wallet(w)
}

// Validate validates the wallet
func (w *Bip44Wallet) Validate() error {
	return w.Meta.validate()
}

// GetAddresses returns all addresses in wallet
func (w *Bip44Wallet) GetAddresses() []cipher.Addresser {
	return append(w.ExternalEntries.getAddresses(), w.ChangeEntries.getAddresses()...)
}

// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
func (w *Bip44Wallet) GetSkycoinAddresses() ([]cipher.Address, error) {
	if w.Meta.Coin() != CoinTypeSkycoin {
		return nil, errors.New("Bip44Wallet coin type is not skycoin")
	}

	return append(w.ExternalEntries.getSkycoinAddresses(), w.ChangeEntries.getSkycoinAddresses()...), nil
}

// GetEntries returns a copy of all entries held by the wallet
func (w *Bip44Wallet) GetEntries() Entries {
	return append(w.ExternalEntries.clone(), w.ChangeEntries.clone()...)
}

// EntriesLen returns the number of entries in the wallet
func (w *Bip44Wallet) EntriesLen() int {
	return len(w.ExternalEntries) + len(w.ChangeEntries)
}

// GetEntryAt returns entry at a given index in the entries array
func (w *Bip44Wallet) GetEntryAt(i int) Entry {
	if i >= len(w.ExternalEntries) {
		return w.ChangeEntries[i-len(w.ExternalEntries)]
	}
	return w.ExternalEntries[i]
}

// GetEntry returns entry of given address
func (w *Bip44Wallet) GetEntry(a cipher.Address) (Entry, bool) {
	if e, ok := w.ExternalEntries.get(a); ok {
		return e, true
	}

	return w.ChangeEntries.get(a)
}

// HasEntry returns true if the wallet has an Entry with a given cipher.Address.
func (w *Bip44Wallet) HasEntry(a cipher.Address) bool {
	return w.ExternalEntries.has(a) || w.ChangeEntries.has(a)
}

// GenerateAddresses generates addresses
func (w *Bip44Wallet) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	if num == 0 {
		return nil, nil
	}

	if w.Meta.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	var seckeys []cipher.SecKey
	var seed []byte
	if w.EntriesLen() == 0 {
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
func (w *Bip44Wallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
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

// ScanAddresses scans ahead N addresses, truncating up to the highest address with a non-zero balance.
// If any address has a nonzero balance, it rescans N more addresses from that point, until a entire
// sequence of N addresses has no balance.
func (w *Bip44Wallet) ScanAddresses(scanN uint64, tf TransactionsFinder) error {
	if w.Meta.IsEncrypted() {
		return ErrWalletEncrypted
	}

	if scanN == 0 {
		return nil
	}

	w2 := w.Clone().(*Bip44Wallet)

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
	// w2.reset()
	if _, err := w2.GenerateSkycoinAddresses(nExistingAddrs + nAddAddrs); err != nil {
		return err
	}

	*w = *w2

	return nil
}

// Fingerprint returns a unique ID fingerprint this wallet, composed of its initial address
// and wallet type
func (w *Bip44Wallet) Fingerprint() string {
	addr := ""
	if len(w.ExternalEntries) == 0 {
		if !w.IsEncrypted() {
			_, pk, _ := cipher.MustDeterministicKeyPairIterator([]byte(w.Meta.Seed()))
			addr = w.Meta.AddressConstructor()(pk).String()
		}
	} else {
		addr = w.ExternalEntries[0].Address.String()
	}
	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

// ReadableBip44Wallet used for [de]serialization of a deterministic wallet
type ReadableBip44Wallet struct {
	Meta            `json:"meta"`
	ExternalEntries `json:"external_entries"`
	ChangeEntries   `json:"change_entries"`
}

// LoadReadableBip44Wallet loads a deterministic wallet from disk
func LoadReadableBip44Wallet(wltFile string) (*ReadableBip44Wallet, error) {
	var rw ReadableBip44Wallet
	if err := file.LoadJSON(wltFile, &rw); err != nil {
		return nil, err
	}
	if rw.Type() != WalletTypeBip44 {
		return nil, ErrInvalidWalletType
	}
	return &rw, nil
}

// NewReadableBip44Wallet creates readable wallet
func NewReadableBip44Wallet(w *Bip44Wallet) *ReadableBip44Wallet {
	return &ReadableBip44Wallet{
		Meta:            w.Meta.clone(),
		ExternalEntries: newReadableEntries(w.ExternalEntries, w.Meta.Coin()),
		ChangeEntries:   newReadableEntries(w.ChangeEntries, w.Meta.Coin()),
	}
}

// ToWallet convert readable wallet to Wallet
func (rw *ReadableBip44Wallet) ToWallet() (Wallet, error) {
	w := &Bip44Wallet{
		Meta: rw.Meta.clone(),
	}

	if err := w.Validate(); err != nil {
		err := fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
		logger.WithError(err).Error("ReadableBip44Wallet.ToWallet Validate failed")
		return nil, err
	}

	ets, err := rw.ExternalEntries.toWalletEntries(w.Meta.Coin(), w.Meta.IsEncrypted())
	if err != nil {
		logger.WithError(err).Error("ReadableBip44Wallet.ToWallet ExternalEntries.toWalletEntries failed")
		return nil, err
	}

	w.ExternalEntries = ets

	ets, err = rw.ChangeEntries.toWalletEntries(w.Meta.Coin(), w.Meta.IsEncrypted())
	if err != nil {
		logger.WithError(err).Error("ReadableBip44Wallet.ToWallet ExternalEntries.toWalletEntries failed")
		return nil, err
	}

	w.ChangeEntries = ets

	return w, nil
}
