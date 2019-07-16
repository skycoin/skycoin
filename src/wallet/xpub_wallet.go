package wallet

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip32"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/mathutil"
)

// XPubWallet holds a single xpub (extended public key) and derives child public keys from it.
// Refer to the bip32 spec to understand xpub keys.
// XPub wallets can generate new addresses and receive coins, but can't spend coins
// because the private keys are not available.
type XPubWallet struct {
	Meta
	Entries Entries
	xpub    *bip32.PublicKey
}

// newXPubWallet creates a XPubWallet
func newXPubWallet(meta Meta) (*XPubWallet, error) {
	xpub, err := parseXPub(meta.XPub())
	if err != nil {
		return nil, err
	}

	return &XPubWallet{
		Meta: meta,
		xpub: xpub,
	}, nil
}

func parseXPub(xp string) (*bip32.PublicKey, error) {
	xpub, err := bip32.DeserializeEncodedPublicKey(xp)
	if err != nil {
		logger.WithError(err).Error("bip32.DeserializeEncodedPublicKey failed")
		return nil, NewError(fmt.Errorf("invalid xpub key: %v", err))
	}

	return xpub, nil
}

// PackSecrets does nothing because XPubWallet has no secrets
func (w *XPubWallet) PackSecrets(ss Secrets) {
}

// UnpackSecrets does nothing because XPubWallet has no secrets
func (w *XPubWallet) UnpackSecrets(ss Secrets) error {
	return nil
}

// Clone clones the wallet a new wallet object
func (w *XPubWallet) Clone() Wallet {
	xpub, err := parseXPub(w.Meta.XPub())
	if err != nil {
		logger.WithError(err).Panic("Clone parseXPub failed")
	}

	return &XPubWallet{
		Meta:    w.Meta.clone(),
		Entries: w.Entries.clone(),
		xpub:    xpub,
	}
}

// CopyFrom copies the src wallet to w
func (w *XPubWallet) CopyFrom(src Wallet) {
	xpub, err := parseXPub(src.XPub())
	if err != nil {
		logger.WithError(err).Panic("CopyFrom parseXPub failed")
	}
	w.xpub = xpub
	w.Meta = src.(*XPubWallet).Meta.clone()
	w.Entries = src.(*XPubWallet).Entries.clone()
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *XPubWallet) CopyFromRef(src Wallet) {
	xpub, err := parseXPub(src.XPub())
	if err != nil {
		logger.WithError(err).Panic("CopyFromRef parseXPub failed")
	}

	*w = *(src.(*XPubWallet))
	w.xpub = xpub
}

// Erase wipes secret fields in wallet
func (w *XPubWallet) Erase() {
	w.Meta.eraseSeeds()
	w.Entries.erase()
}

// ToReadable converts the wallet to its readable (serializable) format
func (w *XPubWallet) ToReadable() Readable {
	return NewReadableXPubWallet(w)
}

// Validate validates the wallet
func (w *XPubWallet) Validate() error {
	return w.Meta.validate()
}

// GetAddresses returns all addresses in wallet
func (w *XPubWallet) GetAddresses() []cipher.Addresser {
	return w.Entries.getAddresses()
}

// GetSkycoinAddresses returns all Skycoin addresses in wallet. The wallet's coin type must be Skycoin.
func (w *XPubWallet) GetSkycoinAddresses() ([]cipher.Address, error) {
	if w.Meta.Coin() != CoinTypeSkycoin {
		return nil, errors.New("XPubWallet coin type is not skycoin")
	}

	return w.Entries.getSkycoinAddresses(), nil
}

// GetEntries returns a copy of all entries held by the wallet
func (w *XPubWallet) GetEntries() Entries {
	return w.Entries.clone()
}

// EntriesLen returns the number of entries in the wallet
func (w *XPubWallet) EntriesLen() int {
	return len(w.Entries)
}

// GetEntryAt returns entry at a given index in the entries array
func (w *XPubWallet) GetEntryAt(i int) Entry {
	return w.Entries[i]
}

// GetEntry returns entry of given address
func (w *XPubWallet) GetEntry(a cipher.Address) (Entry, bool) {
	return w.Entries.get(a)
}

// HasEntry returns true if the wallet has an Entry with a given cipher.Address.
func (w *XPubWallet) HasEntry(a cipher.Address) bool {
	return w.Entries.has(a)
}

// generateEntries generates up to `num` addresses
func (w *XPubWallet) generateEntries(num uint64, initialChildIdx uint32) (Entries, error) {
	if w.Meta.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	if num > math.MaxUint32 {
		return nil, NewError(errors.New("XPubWallet.generateEntries num too large"))
	}

	// Cap `num` in case it would exceed the maximum child index number
	if math.MaxUint32-initialChildIdx < uint32(num) {
		num = uint64(math.MaxUint32 - initialChildIdx)
	}

	if num == 0 {
		return nil, nil
	}

	// Generate `num` secret keys from the external chain HDNode, skipping any children that
	// are invalid (note that this has probability ~2^-128)
	var pubkeys []*bip32.PublicKey
	var addressIndices []uint32
	j := initialChildIdx
	for i := uint32(0); i < uint32(num); i++ {
		k, err := w.xpub.NewPublicChildKey(j)

		var addErr error
		j, addErr = mathutil.AddUint32(j, 1)
		if addErr != nil {
			logger.Critical().WithError(addErr).WithFields(logrus.Fields{
				"num":             num,
				"initialChildIdx": initialChildIdx,
				"childIdx":        j,
				"i":               i,
			}).Error("childIdx can't be incremented any further")
			return nil, errors.New("childIdx can't be incremented any further")
		}

		if err != nil {
			if bip32.IsImpossibleChildError(err) {
				logger.Critical().WithError(err).WithField("childIdx", j).Error("ImpossibleChild for xpub child element")
				continue
			} else {
				logger.Critical().WithError(err).WithField("childIdx", j).Error("NewPublicChildKey failed unexpectedly")
				return nil, err
			}
		}

		pubkeys = append(pubkeys, k)
		addressIndices = append(addressIndices, j-1)
	}

	entries := make(Entries, len(pubkeys))
	makeAddress := w.Meta.AddressConstructor()
	for i, xp := range pubkeys {
		pk := cipher.MustNewPubKey(xp.Key)
		entries[i] = Entry{
			Address:     makeAddress(pk),
			Public:      pk,
			ChildNumber: addressIndices[i],
		}
	}

	return entries, nil
}

// GenerateAddresses generates addresses for the external chain, and appends them to the wallet's entries array
func (w *XPubWallet) GenerateAddresses(num uint64) ([]cipher.Addresser, error) {
	entries, err := w.generateEntries(num, nextChildIdx(w.Entries))
	if err != nil {
		return nil, err
	}

	w.Entries = append(w.Entries, entries...)

	return entries.getAddresses(), nil
}

// GenerateSkycoinAddresses generates Skycoin addresses for the external chain, and appends them to the wallet's entries array.
// If the wallet's coin type is not Skycoin, returns an error
func (w *XPubWallet) GenerateSkycoinAddresses(num uint64) ([]cipher.Address, error) {
	if w.Meta.Coin() != CoinTypeSkycoin {
		return nil, errors.New("GenerateSkycoinAddresses called for non-skycoin wallet")
	}

	entries, err := w.generateEntries(num, nextChildIdx(w.Entries))
	if err != nil {
		return nil, err
	}

	w.Entries = append(w.Entries, entries...)

	return entries.getSkycoinAddresses(), nil
}

// ScanAddresses scans ahead N addresses,
// truncating up to the highest address with any transaction history.
func (w *XPubWallet) ScanAddresses(scanN uint64, tf TransactionsFinder) error {
	if w.Meta.IsEncrypted() {
		return ErrWalletEncrypted
	}

	if scanN == 0 {
		return nil
	}

	w2 := w.Clone().(*XPubWallet)

	entries, err := scanAddressesBip32(func(num uint64, childIdx uint32) (Entries, error) {
		return w2.generateEntries(num, childIdx)
	}, scanN, tf, nextChildIdx(w2.Entries))
	if err != nil {
		return err
	}

	w2.Entries = append(w2.Entries, entries...)

	*w = *w2

	return nil
}

// Fingerprint returns a unique ID fingerprint for this wallet, using the first
// child address of the xpub key
func (w *XPubWallet) Fingerprint() string {
	// Note: the xpub key is not used as the fingerprint, because it is
	// partially sensitive data
	addr := ""
	if len(w.Entries) == 0 {
		if !w.IsEncrypted() {
			entries, err := w.generateEntries(1, 0)
			if err != nil {
				logger.WithError(err).Panic("Fingerprint failed to generate initial entry for empty wallet")
			}
			addr = entries[0].Address.String()
		}
	} else {
		addr = w.Entries[0].Address.String()
	}

	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

// ReadableXPubWallet used for [de]serialization of an xpub wallet
type ReadableXPubWallet struct {
	Meta            `json:"meta"`
	ReadableEntries `json:"entries"`
}

// LoadReadableXPubWallet loads an xpub wallet from disk
func LoadReadableXPubWallet(wltFile string) (*ReadableXPubWallet, error) {
	var rw ReadableXPubWallet
	if err := file.LoadJSON(wltFile, &rw); err != nil {
		return nil, err
	}
	if rw.Type() != WalletTypeXPub {
		return nil, ErrInvalidWalletType
	}
	return &rw, nil
}

// NewReadableXPubWallet creates readable wallet
func NewReadableXPubWallet(w *XPubWallet) *ReadableXPubWallet {
	return &ReadableXPubWallet{
		Meta:            w.Meta.clone(),
		ReadableEntries: newReadableEntries(w.Entries, w.Meta.Coin(), w.Meta.Type()),
	}
}

// ToWallet convert readable wallet to Wallet
func (rw *ReadableXPubWallet) ToWallet() (Wallet, error) {
	w := &XPubWallet{
		Meta: rw.Meta.clone(),
	}

	if err := w.Validate(); err != nil {
		err := fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
		logger.WithError(err).Error("ReadableXPubWallet.ToWallet Validate failed")
		return nil, err
	}

	ets, err := rw.ReadableEntries.toWalletEntries(w.Meta.Coin(), w.Meta.Type(), w.Meta.IsEncrypted())
	if err != nil {
		logger.WithError(err).Error("ReadableXPubWallet.ToWallet toWalletEntries failed")
		return nil, err
	}

	w.Entries = ets

	// Sort childNumber low to high
	sort.Slice(w.Entries, func(i, j int) bool {
		return w.Entries[i].ChildNumber < w.Entries[j].ChildNumber
	})

	w.Entries = ets

	return w, nil
}
