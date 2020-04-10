package wallet

import (
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/file"
	"github.com/SkycoinProject/skycoin/src/wallet/entry"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/SkycoinProject/skycoin/src/wallet/secrets"
)

// CollectionWallet manages keys as an arbitrary collection.
// It has no defined keypair generator. The only way to add keys to the
// wallet is to explicitly add them.
// This wallet does not support address scanning or generation.
// This wallet does not use seeds.
type CollectionWallet struct {
	meta.Meta
	Entries entry.Entries
}

// newCollectionWallet creates a CollectionWallet
func newCollectionWallet(meta meta.Meta) (*CollectionWallet, error) { //nolint:unparam
	return &CollectionWallet{
		Meta: meta,
	}, nil
}

// PackSecrets copies data from decrypted wallets into the secrets container
func (w *CollectionWallet) PackSecrets(ss secrets.Secrets) {
	ss.Set(secrets.SecretSeed, w.Meta.Seed())
	ss.Set(secrets.SecretLastSeed, w.Meta.LastSeed())

	// Saves entry secret keys in secrets
	for _, e := range w.Entries {
		ss.Set(e.Address.String(), e.Secret.Hex())
	}
}

// UnpackSecrets copies data from decrypted secrets into the wallet
func (w *CollectionWallet) UnpackSecrets(ss secrets.Secrets) error {
	return w.Entries.UnpackSecretKeys(ss)
}

// Clone clones the wallet a new wallet object
func (w *CollectionWallet) Clone() Wallet {
	return &CollectionWallet{
		Meta:    w.Meta.Clone(),
		Entries: w.Entries.Clone(),
	}
}

// CopyFrom copies the src wallet by reallocating
func (w *CollectionWallet) CopyFrom(src Wallet) {
	w.Meta = src.(*CollectionWallet).Meta.Clone()
	w.Entries = src.(*CollectionWallet).Entries.Clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *CollectionWallet) CopyFromRef(src Wallet) {
	*w = *(src.(*CollectionWallet))
}

// Erase wipes secret fields in wallet
func (w *CollectionWallet) Erase() {
	w.Meta.EraseSeeds()
	w.Entries.Erase()
}

// ToReadable converts the wallet to its readable (serializable) format
func (w *CollectionWallet) ToReadable() Readable {
	return NewReadableCollectionWallet(w)
}

// Validate validates the wallet
func (w *CollectionWallet) Validate() error {
	return metaValidate(w.Meta)
}

// GetEntries returns a copy of all entries held by the wallet
func (w *CollectionWallet) GetEntries() entry.Entries {
	return w.Entries.Clone()
}

// EntriesLen returns the number of entries in the wallet
func (w *CollectionWallet) EntriesLen() int {
	return len(w.Entries)
}

// GetEntryAt returns entry at a given index in the entries array
func (w *CollectionWallet) GetEntryAt(i int) entry.Entry {
	return w.Entries[i]
}

// GetEntry returns entry of given address
func (w *CollectionWallet) GetEntry(a cipher.Address) (entry.Entry, bool) {
	return w.Entries.Get(a)
}

// HasEntry returns true if the wallet has an entry.Entry with a given cipher.Address.
func (w *CollectionWallet) HasEntry(a cipher.Address) bool {
	return w.Entries.Has(a)
}

// GenerateAddresses is a no-op for "collection" wallets
func (w *CollectionWallet) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	return nil, NewError(errors.New("A collection wallet does not implement GenerateAddresses"))
}

// GenerateSkycoinAddresses is a no-op for "collection" wallets
func (w *CollectionWallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
	return nil, NewError(errors.New("A collection wallet does not implement GenerateSkycoinAddresses"))
}

// ScanAddresses is a no-op for "collection" wallets
func (w *CollectionWallet) ScanAddresses(scanN uint64, tf TransactionsFinder) error {
	return NewError(errors.New("A collection wallet does not implement ScanAddresses"))
}

// GetAddresses returns all addresses in wallet
func (w *CollectionWallet) GetAddresses() []cipher.Addresser {
	return w.Entries.GetAddresses()
}

// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
func (w *CollectionWallet) GetSkycoinAddresses() ([]cipher.Address, error) {
	if w.Meta.Coin() != meta.CoinTypeSkycoin {
		return nil, errors.New("CollectionWallet coin type is not skycoin")
	}

	return w.Entries.GetSkycoinAddresses(), nil
}

// Fingerprint returns an empty string; fingerprints are only defined for
// wallets with a seed
func (w *CollectionWallet) Fingerprint() string {
	return ""
}

// AddEntry adds a new entry to the wallet.
func (w *CollectionWallet) AddEntry(e entry.Entry) error {
	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	if err := e.Verify(); err != nil {
		return err
	}

	for _, entry := range w.Entries {
		if e.SkycoinAddress() == entry.SkycoinAddress() {
			return errors.New("wallet already contains entry with this address")
		}
	}

	w.Entries = append(w.Entries, e)
	return nil
}

// ReadableCollectionWallet used for [de]serialization of a collection wallet
type ReadableCollectionWallet struct {
	meta.Meta       `json:"meta"`
	ReadableEntries `json:"entries"`
}

// NewReadableCollectionWallet creates readable wallet
func NewReadableCollectionWallet(w *CollectionWallet) *ReadableCollectionWallet {
	return &ReadableCollectionWallet{
		Meta:            w.Meta.Clone(),
		ReadableEntries: newReadableEntries(w.Entries, w.Meta.Coin(), w.Meta.Type()),
	}
}

// LoadReadableCollectionWallet loads a collection wallet from disk
func LoadReadableCollectionWallet(wltFile string) (*ReadableCollectionWallet, error) {
	logger.WithField("filename", wltFile).Info("LoadReadableCollectionWallet")
	var rw ReadableCollectionWallet
	if err := file.LoadJSON(wltFile, &rw); err != nil {
		return nil, err
	}
	if rw.Type() != WalletTypeCollection {
		return nil, ErrInvalidWalletType
	}
	return &rw, nil
}

// ToWallet convert readable wallet to Wallet
func (rw *ReadableCollectionWallet) ToWallet() (Wallet, error) {
	w := &CollectionWallet{
		Meta: rw.Meta.Clone(),
	}

	if err := w.Validate(); err != nil {
		err := fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
		logger.WithError(err).Error("ReadableCollectionWallet.ToWallet Validate failed")
		return nil, err
	}

	ets, err := rw.ReadableEntries.toWalletEntries(w.Meta.Coin(), w.Meta.Type(), w.Meta.IsEncrypted())
	if err != nil {
		logger.WithError(err).Error("ReadableCollectionWallet.ToWallet entry.toWalletEntries failed")
		return nil, err
	}

	w.Entries = ets

	return w, nil
}
