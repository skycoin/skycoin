package xpubwallet

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/sirupsen/logrus"
)

// WalletType represents the xpub wallet type
const WalletType = "xpub"

var defaultWalletDecoder = &JSONDecoder{}
var logger = logging.MustGetLogger("xpubwallet")

// Wallet holds a single xpub (extended public key) and derives child public keys from it.
// Refer to the bip32 spec to understand xpub keys.
// XPub wallets can generate new addresses and receive coins, but can't spend coins
// because the private keys are not available.
type Wallet struct {
	wallet.Meta
	entries wallet.Entries
	xpub    *bip32.PublicKey
	decoder wallet.Decoder
}

// NewWallet creates a xpub wallet with options
func NewWallet(filename, label, xPub string, options ...wallet.Option) (wallet.Wallet, error) {
	key, err := parseXPub(xPub)
	if err != nil {
		return nil, err
	}

	wlt := &Wallet{
		Meta: wallet.Meta{
			wallet.MetaFilename:  filename,
			wallet.MetaLabel:     label,
			wallet.MetaEncrypted: "false",
			wallet.MetaType:      WalletType,
			wallet.MetaVersion:   wallet.Version,
			wallet.MetaCoin:      string(wallet.CoinTypeSkycoin),
			wallet.MetaXPub:      xPub,
			wallet.MetaTimestamp: strconv.FormatInt(time.Now().Unix(), 10),
		},
		xpub: key,
	}

	advOpts := &wallet.AdvancedOptions{}
	for _, opt := range options {
		opt(wlt)
		opt(advOpts)
	}

	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}

	return wlt, nil
}

// SetDecoder sets the wallet decoder
func (w *Wallet) SetDecoder(d wallet.Decoder) {
	w.decoder = d
}

func validateMeta(m wallet.Meta) error {
	if m[wallet.MetaType] != WalletType {
		return errors.New("invalid wallet type")
	}

	return wallet.ValidateMeta(m)
}

// Serialize encodes the xpub wallet to []byte
func (w Wallet) Serialize() ([]byte, error) {
	if w.decoder == nil {
		w.decoder = defaultWalletDecoder
	}

	return w.decoder.Encode(&w)
}

// Deserialize decodes the []byte to a xpub wallet
func (w *Wallet) Deserialize(b []byte) error {
	if w.decoder == nil {
		w.decoder = defaultWalletDecoder
	}

	toW, err := w.decoder.Decode(b)
	if err != nil {
		return err
	}

	toW2 := toW.(*Wallet)
	toW2.decoder = w.decoder
	*w = *toW2
	return nil
}

// IsEncrypted returns whether the wallet is encrypted
func (w Wallet) IsEncrypted() bool {
	return w.Meta.IsEncrypted()
}

// Lock will do nothing to the xpub wallet
func (w Wallet) Lock(_ []byte) error {
	return nil
}

// Unlock will return the origin xpub wallet
func (w *Wallet) Unlock(_ []byte) (wallet.Wallet, error) {
	return w, nil
}

// Fingerprint returns a unique ID fingerprint for this wallet, using the first
// child address of the xpub key
func (w *Wallet) Fingerprint() string {
	// Note: the xpub key is not used as the fingerprint, because it is
	// partially sensitive data
	addr := ""
	if len(w.entries) == 0 {
		entries, err := w.generateEntries(1, 0)
		if err != nil {
			logger.WithError(err).Panic("Fingerprint failed to generate initial entry for empty wallet")
		}
		addr = entries[0].Address.String()
	} else {
		addr = w.entries[0].Address.String()
	}

	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

func (w *Wallet) generateEntries(num uint64, initialChildIdx uint32) (wallet.Entries, error) {
	if num > math.MaxUint32 {
		return nil, wallet.NewError(errors.New("XPubWallet.generateEntries num too large"))
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

	entries := make(wallet.Entries, len(pubkeys))
	addressFromPubKey := wallet.ResolveAddressDecoder(w.Coin()).AddressFromPubKey
	for i, xp := range pubkeys {
		pk := cipher.MustNewPubKey(xp.Key)
		entries[i] = wallet.Entry{
			Address:     addressFromPubKey(pk),
			Public:      pk,
			ChildNumber: addressIndices[i],
		}
	}

	return entries, nil
}

// Clone returns a copy of the wallet
func (w Wallet) Clone() wallet.Wallet {
	return &Wallet{
		Meta:    w.Meta.Clone(),
		entries: w.entries.Clone(),
		decoder: w.decoder,
	}
}

// CopyFrom copy wallet from specific wallet
func (w *Wallet) CopyFrom(src wallet.Wallet) {
	w.copyFrom(src.(*Wallet))
}

func (w *Wallet) copyFrom(wlt *Wallet) {
	w.Meta = wlt.Meta.Clone()
	w.entries = wlt.entries.Clone()
	w.decoder = wlt.decoder
}

// CopyFromRef copies the src wallet with a pointer dereference
func (w *Wallet) CopyFromRef(src wallet.Wallet) {
	*w = *(src.(*Wallet))
}

// Accounts is not implemented for xpub wallet
func (w *Wallet) Accounts() []wallet.Bip44Account {
	return nil
}

// GetEntries returns a copy of all entries held by the wallet
func (w *Wallet) GetEntries(_ ...wallet.Option) (wallet.Entries, error) {
	return w.entries.Clone(), nil
}

// Erase removes sensitive data
func (w *Wallet) Erase() {
}

// ScanAddresses scans ahead N addresses, truncating up to the highest address with any transaction history.
func (w *Wallet) ScanAddresses(scanN uint64, tf wallet.TransactionsFinder) ([]cipher.Addresser, error) {
	if scanN == 0 {
		return nil, nil
	}

	w2 := w.Clone().(*Wallet)

	nExistingAddrs := uint64(len(w2.entries))

	// Generate the addresses to scan
	addrs, err := w2.GenerateAddresses(scanN)
	if err != nil {
		return nil, err
	}

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

	w2.reset()
	if _, err := w2.GenerateAddresses(nExistingAddrs + keepNum); err != nil {
		return nil, err
	}

	*w = *w2

	return addrs[:keepNum], nil
}

// GetAddresses returns all addresses of the wallet
func (w *Wallet) GetAddresses(_ ...wallet.Option) ([]cipher.Addresser, error) {
	return w.entries.GetAddresses(), nil
}

// GenerateAddresses generates addresses for the external chain, and appends them to the wallet's entries array
func (w *Wallet) GenerateAddresses(num uint64, _ ...wallet.Option) ([]cipher.Addresser, error) {
	if num > math.MaxUint32 {
		return nil, wallet.NewError(errors.New("XPubWallet.GenerateAddresses num too large"))
	}

	var addrs []cipher.Addresser
	initLen := uint32(len(w.entries))
	_, err := mathutil.AddUint32(initLen, uint32(num))
	if err != nil {
		return nil, fmt.Errorf("generate %d more addresses failed: %v", num, err)
	}

	makeAddress := wallet.ResolveAddressDecoder(w.Coin())

	for i := uint32(0); i < uint32(num); i++ {
		index := initLen + i
		pk, err := w.xpub.NewPublicChildKey(index)
		if err != nil {
			return nil, err
		}
		cpk, err := cipher.NewPubKey(pk.Key)
		if err != nil {
			return nil, err
		}

		addr := makeAddress.AddressFromPubKey(cpk)
		e := wallet.Entry{
			Address:     addr,
			Public:      cpk,
			ChildNumber: index,
		}

		w.entries = append(w.entries, e)
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func parseXPub(xp string) (*bip32.PublicKey, error) {
	xPub, err := bip32.DeserializeEncodedPublicKey(xp)
	if err != nil {
		return nil, fmt.Errorf("invalid xpub key: %v", err)
	}

	return xPub, nil
}

// GetEntryAt returns the entry at a given index in the entries array
func (w *Wallet) GetEntryAt(i int, _ ...wallet.Option) (wallet.Entry, error) {
	if i < 0 || i >= len(w.entries) {
		return wallet.Entry{}, fmt.Errorf("entry index %d is out of range", i)
	}
	return w.entries[i], nil
}

// GetEntry returns a entry of given address
func (w *Wallet) GetEntry(addr cipher.Addresser, _ ...wallet.Option) (wallet.Entry, error) {
	e, ok := w.entries.Get(addr)
	if !ok {
		return wallet.Entry{}, wallet.ErrEntryNotFound
	}
	return e, nil
}

// HasEntry returns true if the wallet has an Entry with a given address
func (w *Wallet) HasEntry(addr cipher.Addresser, _ ...wallet.Option) (bool, error) {
	return w.entries.Has(addr), nil
}

// EntriesLen returns the number of entries in the wallet
func (w *Wallet) EntriesLen(_ ...wallet.Option) (int, error) {
	return len(w.entries), nil
}

// reset resets the wallet entries and move the lastSeed to origin
func (w *Wallet) reset() {
	w.entries = wallet.Entries{}
}
