package wallet

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/util/file"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
	"github.com/SkycoinProject/skycoin/src/wallet/entry"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/SkycoinProject/skycoin/src/wallet/secrets"
)

// Bip44Wallet manages keys using the original Skycoin deterministic
// keypair generator method.
// With this generator, a single chain of addresses is created, each one dependent
// on the previous.
type Bip44Wallet struct {
	meta.Meta
	ExternalEntries entry.Entries
	ChangeEntries   entry.Entries
}

// newBip44Wallet creates a Bip44Wallet
func newBip44Wallet(meta meta.Meta) (*Bip44Wallet, error) { //nolint:unparam
	return &Bip44Wallet{
		Meta: meta,
	}, nil
}

// PackSecrets copies data from decrypted wallets into the secrets container
func (w *Bip44Wallet) PackSecrets(ss secrets.Secrets) {
	ss.Set(secrets.SecretSeed, w.Seed())
	ss.Set(secrets.SecretSeedPassphrase, w.SeedPassphrase())

	// Saves entry secret keys in secrets
	for _, e := range w.ExternalEntries {
		ss.Set(e.Address.String(), e.Secret.Hex())
	}
	for _, e := range w.ChangeEntries {
		ss.Set(e.Address.String(), e.Secret.Hex())
	}
}

// UnpackSecrets copies data from decrypted secrets into the wallet
func (w *Bip44Wallet) UnpackSecrets(ss secrets.Secrets) error {
	seed, ok := ss.Get(secrets.SecretSeed)
	if !ok {
		return errors.New("seed doesn't exist in secrets")
	}
	w.SetSeed(seed)

	passphrase, _ := ss.Get(secrets.SecretSeedPassphrase)
	w.SetSeedPassphrase(passphrase)

	if err := w.ExternalEntries.UnpackSecretKeys(ss); err != nil {
		return err
	}
	return w.ChangeEntries.UnpackSecretKeys(ss)
}

// Clone clones the wallet a new wallet object
func (w *Bip44Wallet) Clone() Wallet {
	return &Bip44Wallet{
		Meta:            w.Meta.Clone(),
		ExternalEntries: w.ExternalEntries.Clone(),
		ChangeEntries:   w.ChangeEntries.Clone(),
	}
}

// CopyFrom copies the src wallet to w
func (w *Bip44Wallet) CopyFrom(src Wallet) {
	w.Meta = src.(*Bip44Wallet).Meta.Clone()
	w.ExternalEntries = src.(*Bip44Wallet).ExternalEntries.Clone()
	w.ChangeEntries = src.(*Bip44Wallet).ChangeEntries.Clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *Bip44Wallet) CopyFromRef(src Wallet) {
	*w = *(src.(*Bip44Wallet))
}

// Erase wipes secret fields in wallet
func (w *Bip44Wallet) Erase() {
	w.Meta.EraseSeeds()
	w.ExternalEntries.Erase()
	w.ChangeEntries.Erase()
}

// ToReadable converts the wallet to its readable (serializable) format
func (w *Bip44Wallet) ToReadable() Readable {
	return NewReadableBip44Wallet(w)
}

// Validate validates the wallet
func (w *Bip44Wallet) Validate() error {
	return metaValidate(w.Meta)
}

// GetAddresses returns all addresses in wallet
func (w *Bip44Wallet) GetAddresses() []cipher.Addresser {
	return append(w.ExternalEntries.GetAddresses(), w.ChangeEntries.GetAddresses()...)
}

// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
func (w *Bip44Wallet) GetSkycoinAddresses() ([]cipher.Address, error) {
	if w.Meta.Coin() != meta.CoinTypeSkycoin {
		return nil, errors.New("Bip44Wallet coin type is not skycoin")
	}

	return append(w.ExternalEntries.GetSkycoinAddresses(), w.ChangeEntries.GetSkycoinAddresses()...), nil
}

// GetEntries returns a copy of all entries held by the wallet
func (w *Bip44Wallet) GetEntries() entry.Entries {
	if w.EntriesLen() == 0 {
		return nil
	}
	return append(w.ExternalEntries.Clone(), w.ChangeEntries.Clone()...)
}

// EntriesLen returns the number of entries in the wallet
func (w *Bip44Wallet) EntriesLen() int {
	return len(w.ExternalEntries) + len(w.ChangeEntries)
}

// GetEntryAt returns entry at a given index in the entries array
func (w *Bip44Wallet) GetEntryAt(i int) entry.Entry {
	if i >= len(w.ExternalEntries) {
		return w.ChangeEntries[i-len(w.ExternalEntries)]
	}
	return w.ExternalEntries[i]
}

// GetEntry returns entry of given address
func (w *Bip44Wallet) GetEntry(a cipher.Address) (entry.Entry, bool) {
	if e, ok := w.ExternalEntries.Get(a); ok {
		return e, true
	}

	return w.ChangeEntries.Get(a)
}

// HasEntry returns true if the wallet has an entry.Entry with a given cipher.Address.
func (w *Bip44Wallet) HasEntry(a cipher.Address) bool {
	return w.ExternalEntries.Has(a) || w.ChangeEntries.Has(a)
}

// CoinHDNode return the "coin" level bip44 HDNode
func (w *Bip44Wallet) CoinHDNode() (*bip44.Coin, error) {
	// w.Meta.Seed() must return a valid bip39 mnemonic
	seed, err := bip39.NewSeed(w.Meta.Seed(), w.Meta.SeedPassphrase())
	if err != nil {
		return nil, err
	}
	bip44Coin, ok := w.Bip44Coin()
	if !ok {
		return nil, errors.New("Bip44Coin does not exist")
	}

	c, err := bip44.NewCoin(seed, bip44Coin)
	if err != nil {
		logger.Critical().WithError(err).Error("Failed to derive the bip44 purpose node")
		if bip32.IsImpossibleChildError(err) {
			logger.Critical().Error("ImpossibleChild: this seed cannot be used for bip44")
		}
		return nil, err
	}

	return c, nil
}

// nextChildIdx returns the next child index from a sequence of entries.
// This assumes that entries are sorted by child number ascending.
func nextChildIdx(e entry.Entries) uint32 {
	if len(e) == 0 {
		return 0
	}
	return e[len(e)-1].ChildNumber + 1
}

// generateEntries generates addresses for a change chain (should be 0 or 1) starting from an initial child number.
func (w *Bip44Wallet) generateEntries(num uint64, changeIdx, initialChildIdx uint32) (entry.Entries, error) {
	if w.Meta.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	if num > math.MaxUint32 {
		return nil, NewError(errors.New("Bip44Wallet.generateEntries num too large"))
	}

	// Cap `num` in case it would exceed the maximum child index number
	if math.MaxUint32-initialChildIdx < uint32(num) {
		num = uint64(math.MaxUint32 - initialChildIdx)
	}

	if num == 0 {
		return nil, nil
	}

	c, err := w.CoinHDNode()
	if err != nil {
		return nil, err
	}

	// Generate the "account" HDNode. Multiple accounts are not supported; use 0.
	account, err := c.Account(0)
	if err != nil {
		logger.Critical().WithError(err).Error("Failed to derive the bip44 account node")
		if bip32.IsImpossibleChildError(err) {
			logger.Critical().Error("ImpossibleChild: this seed cannot be used for bip44")
		}
		return nil, err
	}

	// Generate the chain parent node
	var chain *bip32.PrivateKey
	switch changeIdx {
	case bip44.ExternalChainIndex:
		chain, err = account.External()
	case bip44.ChangeChainIndex:
		chain, err = account.Change()
	default:
		err = errors.New("invalid chain index")
	}
	if err != nil {
		logger.Critical().WithError(err).Error("Failed to derive the final bip44 chain node")
		if bip32.IsImpossibleChildError(err) {
			logger.Critical().Error("ImpossibleChild: this seed cannot be used for bip44")
		}
		return nil, err
	}

	// Generate `num` secret keys from the external chain HDNode, skipping any children that
	// are invalid (note that this has probability ~2^-128)
	var seckeys []*bip32.PrivateKey
	var addressIndices []uint32
	j := initialChildIdx
	for i := uint32(0); i < uint32(num); i++ {
		k, err := chain.NewPrivateChildKey(j)

		var addErr error
		j, addErr = mathutil.AddUint32(j, 1)
		if addErr != nil {
			logger.Critical().WithError(addErr).WithFields(logrus.Fields{
				"num":             num,
				"initialChildIdx": initialChildIdx,
				"accountIdx":      0,
				"changeIdx":       changeIdx,
				"childIdx":        j,
				"i":               i,
			}).Error("childIdx can't be incremented any further")
			return nil, errors.New("childIdx can't be incremented any further")
		}

		if err != nil {
			if bip32.IsImpossibleChildError(err) {
				logger.Critical().WithError(err).WithFields(logrus.Fields{
					"accountIdx": 0,
					"changeIdx":  changeIdx,
					"childIdx":   j,
				}).Error("ImpossibleChild for chain node child element")
				continue
			} else {
				logger.Critical().WithError(err).WithFields(logrus.Fields{
					"accountIdx": 0,
					"changeIdx":  changeIdx,
					"childIdx":   j,
				}).Error("NewPrivateChildKey failed unexpectedly")
				return nil, err
			}
		}

		seckeys = append(seckeys, k)
		addressIndices = append(addressIndices, j-1)
	}

	entries := make(entry.Entries, len(seckeys))
	makeAddress := AddressConstructor(w.Meta)
	for i, xprv := range seckeys {
		sk := cipher.MustNewSecKey(xprv.Key)
		pk := cipher.MustPubKeyFromSecKey(sk)
		entries[i] = entry.Entry{
			Address:     makeAddress(pk),
			Secret:      sk,
			Public:      pk,
			ChildNumber: addressIndices[i],
			Change:      changeIdx,
		}
	}

	return entries, nil
}

// PeekChangeEntry creates and returns an entry for the change chain.
// If used, the caller the append it with GenerateChangeEntry
func (w *Bip44Wallet) PeekChangeEntry() (entry.Entry, error) {
	entries, err := w.generateEntries(1, bip44.ChangeChainIndex, nextChildIdx(w.ChangeEntries))
	if err != nil {
		return entry.Entry{}, err
	}

	if len(entries) == 0 {
		return entry.Entry{}, NewError(errors.New("PeekChangeEntry: no more change addresses"))
	}

	return entries[0], nil
}

// GenerateChangeEntry creates, appends and returns an entry for the change chain
func (w *Bip44Wallet) GenerateChangeEntry() (entry.Entry, error) {
	e, err := w.PeekChangeEntry()
	if err != nil {
		return entry.Entry{}, err
	}

	w.ChangeEntries = append(w.ChangeEntries, entry.Entries{e}...)

	return w.ChangeEntries[len(w.ChangeEntries)-1], nil
}

// GenerateAddresses generates addresses for the external chain, and appends them to the wallet's entries array
func (w *Bip44Wallet) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	entries, err := w.generateEntries(num, bip44.ExternalChainIndex, nextChildIdx(w.ExternalEntries))
	if err != nil {
		return nil, err
	}

	w.ExternalEntries = append(w.ExternalEntries, entries...)

	return entries.GetAddresses(), nil
}

// GenerateSkycoinAddresses generates Skycoin addresses for the external chain, and appends them to the wallet's entries array.
// If the wallet's coin type is not Skycoin, returns an error
func (w *Bip44Wallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
	if w.Meta.Coin() != meta.CoinTypeSkycoin {
		return nil, errors.New("GenerateSkycoinAddresses called for non-skycoin wallet")
	}

	entries, err := w.generateEntries(num, bip44.ExternalChainIndex, nextChildIdx(w.ExternalEntries))
	if err != nil {
		return nil, err
	}

	w.ExternalEntries = append(w.ExternalEntries, entries...)

	return entries.GetSkycoinAddresses(), nil
}

// ScanAddresses scans ahead N addresses, truncating up to the highest address with any transaction history.
func (w *Bip44Wallet) ScanAddresses(scanN uint64, tf TransactionsFinder) error {
	if w.Meta.IsEncrypted() {
		return ErrWalletEncrypted
	}

	if scanN == 0 {
		return nil
	}

	w2 := w.Clone().(*Bip44Wallet)

	externalEntries, err := scanAddressesBip32(func(num uint64, childIdx uint32) (entry.Entries, error) {
		return w.generateEntries(num, bip44.ExternalChainIndex, childIdx)
	}, scanN, tf, nextChildIdx(w2.ExternalEntries))
	if err != nil {
		return err
	}

	changeEntries, err := scanAddressesBip32(func(num uint64, childIdx uint32) (entry.Entries, error) {
		return w.generateEntries(num, bip44.ChangeChainIndex, childIdx)
	}, scanN, tf, nextChildIdx(w2.ChangeEntries))
	if err != nil {
		return err
	}

	// Add scanned entries
	w2.ExternalEntries = append(w2.ExternalEntries, externalEntries...)
	w2.ChangeEntries = append(w2.ChangeEntries, changeEntries...)

	*w = *w2

	return nil
}

// scanAddressesBip32 implements the address scanning algorithm for bip32
// based (e.g. bip44, xpub) wallets
func scanAddressesBip32(generateEntries func(num uint64, childIdx uint32) (entry.Entries, error), scanN uint64, tf TransactionsFinder, initialChildIdx uint32) (entry.Entries, error) {
	if scanN == 0 {
		return nil, nil
	}

	nAddAddrs := uint64(0)
	var newEntries entry.Entries

	// Generate the addresses to scan
	entries, err := generateEntries(scanN, initialChildIdx)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, nil
	}

	newEntries = append(newEntries, entries...)

	addrs := entries.GetSkycoinAddresses()

	// Find if these addresses had any activity
	active, err := tf.AddressesActivity(addrs)
	if err != nil {
		return nil, err
	}

	// Check activity from the last one until we find the address that has activity
	var keepNum uint64
	for i := len(active) - 1; i >= 0; i-- {
		if active[i] {
			keepNum = uint64(i + 1)
			break
		}
	}

	nAddAddrs += keepNum

	return newEntries[:nAddAddrs], nil
}

// Fingerprint returns a unique ID fingerprint for this wallet, composed of its initial address
// and wallet type
func (w *Bip44Wallet) Fingerprint() string {
	addr := ""
	if len(w.ExternalEntries) == 0 {
		if !w.IsEncrypted() {
			entries, err := w.generateEntries(1, bip44.ExternalChainIndex, 0)
			if err != nil {
				logger.WithError(err).Panic("Fingerprint failed to generate initial entry for empty wallet")
			}
			addr = entries[0].Address.String()
		}
	} else {
		addr = w.ExternalEntries[0].Address.String()
	}
	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

// ReadableBip44Wallet used for [de]serialization of a deterministic wallet
type ReadableBip44Wallet struct {
	meta.Meta       `json:"meta"`
	ReadableEntries `json:"entries"`
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
		Meta:            w.Meta.Clone(),
		ReadableEntries: newReadableEntries(w.GetEntries(), w.Meta.Coin(), w.Meta.Type()),
	}
}

// ToWallet convert readable wallet to Wallet
func (rw *ReadableBip44Wallet) ToWallet() (Wallet, error) {
	w := &Bip44Wallet{
		Meta: rw.Meta.Clone(),
	}

	if err := w.Validate(); err != nil {
		err := fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
		logger.WithError(err).Error("ReadableBip44Wallet.ToWallet Validate failed")
		return nil, err
	}

	ets, err := rw.ReadableEntries.toWalletEntries(w.Meta.Coin(), w.Meta.Type(), w.Meta.IsEncrypted())
	if err != nil {
		logger.WithError(err).Error("ReadableBip44Wallet.ToWallet ReadableEntries.toWalletEntries failed")
		return nil, err
	}

	// Split the single array of entries into separate external and change chains,
	// for easier internal management
	for _, e := range ets {
		switch e.Change {
		case bip44.ExternalChainIndex:
			w.ExternalEntries = append(w.ExternalEntries, e)
		case bip44.ChangeChainIndex:
			w.ChangeEntries = append(w.ChangeEntries, e)
		default:
			logger.Panicf("invalid change value %d", e.Change)
		}
	}

	// Sort childNumber low to high
	sort.Slice(w.ExternalEntries, func(i, j int) bool {
		return w.ExternalEntries[i].ChildNumber < w.ExternalEntries[j].ChildNumber
	})
	sort.Slice(w.ChangeEntries, func(i, j int) bool {
		return w.ChangeEntries[i].ChildNumber < w.ChangeEntries[j].ChildNumber
	})

	return w, err
}
