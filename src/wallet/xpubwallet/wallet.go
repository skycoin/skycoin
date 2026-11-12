// Package xpubwallet src/wallet/xpubwallet/wallet.go c1-coin-wallet
package xpubwallet

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/0magnet/logrus"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip32"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/mathutil"
	"github.com/skycoin/skycoin/src/wallet"
)

// WalletType represents the xpub wallet type
const WalletType = "xpub"

var defaultWalletDecoder = &JSONDecoder{}
var logger = logging.MustGetLogger("xpubwallet")

func init() {
	if err := wallet.RegisterCreator(WalletType, &Creator{}); err != nil {
		panic(err)
	}

	if err := wallet.RegisterLoader(WalletType, &Loader{}); err != nil {
		panic(err)
	}
}

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
func NewWallet(filename, label, xPub string, options ...wallet.Option) (*Wallet, error) {
	if xPub == "" {
		return nil, wallet.ErrMissingXPub
	}

	key, err := parseXPub(xPub)
	if err != nil {
		return nil, wallet.NewError(err)
	}

	meta := wallet.Meta{
		wallet.MetaFilename:  filename,
		wallet.MetaLabel:     label,
		wallet.MetaType:      WalletType,
		wallet.MetaVersion:   wallet.Version,
		wallet.MetaCoin:      string(wallet.CoinTypeSkycoin),
		wallet.MetaXPub:      xPub,
		wallet.MetaTimestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	// A bip44 account key gives its receive addresses on the external chain,
	// m/44'/coin'/account'/0/i. The chain is recorded so a reload derives the same keys.
	if key.Depth == bip44.AccountDepth {
		meta[wallet.MetaXPubChain] = strconv.FormatUint(uint64(bip44.ExternalChainIndex), 10)
	}

	key, err = addressKey(key, meta)
	if err != nil {
		return nil, wallet.NewError(err)
	}

	wlt := &Wallet{
		Meta:    meta,
		decoder: defaultWalletDecoder,
		xpub:    key,
	}

	advOpts := &wallet.AdvancedOptions{}
	for _, opt := range options {
		opt(wlt)
		opt(advOpts)
	}

	if err := validateMeta(wlt.Meta); err != nil {
		return nil, err
	}

	generateN := advOpts.GenerateN
	if generateN > 0 {
		_, err := wlt.GenerateAddresses(wallet.OptionGenerateN(generateN))
		if err != nil {
			return nil, err
		}
	}

	scanN := advOpts.ScanN
	if scanN > 0 {
		if advOpts.TF == nil {
			return nil, errors.New("missing transaction finder for scanning addresses")
		}

		if scanN > generateN {
			scanN = scanN - generateN
		}

		if _, err := wlt.ScanAddresses(scanN, advOpts.TF); err != nil {
			return nil, err
		}
	}

	return wlt, nil
}

// SetDecoder sets the wallet decoder
func (w *Wallet) SetDecoder(d wallet.Decoder) {
	w.decoder = d
}

func validateMeta(m wallet.Meta) error {
	if m[wallet.MetaType] != WalletType {
		return wallet.ErrInvalidWalletType
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
	return wallet.NewError(errors.New("xpub wallet does not support encryption"))
}

// Unlock will return the origin xpub wallet
func (w *Wallet) Unlock(_ []byte) (wallet.Wallet, error) {
	return nil, wallet.NewError(errors.New("xpub wallet does not support encryption"))
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
	for i := uint32(0); i < uint32(num); i++ { //nolint:gosec
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
			} else { //nolint:revive
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
	xpub := w.xpub.Clone()
	return &Wallet{
		Meta:    w.Meta.Clone(),
		entries: w.entries.Clone(),
		xpub:    &xpub,
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

// ScanAddresses scans past the last address for addresses with transactions,
// until scanN in a row have none, the BIP44 address gap limit.
func (w *Wallet) ScanAddresses(scanN uint64, tf wallet.TransactionsFinder) ([]cipher.Addresser, error) {
	if scanN == 0 {
		return nil, nil
	}

	w2 := w.Clone().(*Wallet)

	nExistingAddrs := uint64(len(w2.entries))

	generate := func(n uint64) ([]cipher.Addresser, error) {
		return w2.GenerateAddresses(wallet.OptionGenerateN(n))
	}
	addrs, err := wallet.ScanGap(scanN, generate, tf)
	if err != nil {
		return nil, err
	}

	w2.reset()
	if _, err := w2.GenerateAddresses(wallet.OptionGenerateN(nExistingAddrs + uint64(len(addrs)))); err != nil {
		return nil, err
	}

	*w = *w2

	return addrs, nil
}

// GetAddresses returns all addresses of the wallet
func (w *Wallet) GetAddresses(_ ...wallet.Option) ([]cipher.Addresser, error) {
	return w.entries.GetAddresses(), nil
}

// GenerateAddresses generates addresses for the external chain, and appends them to the wallet's entries array
func (w *Wallet) GenerateAddresses(options ...wallet.Option) ([]cipher.Addresser, error) {
	num := wallet.GetGenerateNFromOptions(options...)
	if num > math.MaxUint32 {
		return nil, wallet.NewError(errors.New("XPubWallet.GenerateAddresses num too large"))
	}

	var addrs []cipher.Addresser
	initLen := uint32(len(w.entries)) //nolint:gosec
	_, err := mathutil.AddUint32(initLen, uint32(num))
	if err != nil {
		return nil, fmt.Errorf("generate %d more addresses failed: %v", num, err)
	}

	makeAddress := wallet.ResolveAddressDecoder(w.Coin())

	for i := uint32(0); i < uint32(num); i++ { //nolint:gosec
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

// addressKey returns the key whose direct children are the wallet's addresses.
// Wallets without the xpubChain meta field derive from the xpub itself, as they always have.
func addressKey(xpub *bip32.PublicKey, m wallet.Meta) (*bip32.PublicKey, error) {
	chainStr, ok := m[wallet.MetaXPubChain]
	if !ok {
		return xpub, nil
	}

	chain, err := strconv.ParseUint(chainStr, 10, 32)
	if err != nil || uint32(chain) != bip44.ExternalChainIndex {
		return nil, fmt.Errorf("invalid %s value %q", wallet.MetaXPubChain, chainStr)
	}

	if xpub.Depth != bip44.AccountDepth {
		return nil, fmt.Errorf("%s is set but the xpub is not a bip44 account key (depth %d)", wallet.MetaXPubChain, xpub.Depth)
	}

	return xpub.NewPublicChildKey(uint32(chain))
}

// checkEntries confirms the stored entries were derived from key, so a wallet
// never mixes addresses from two chains.
func checkEntries(key *bip32.PublicKey, entries wallet.Entries) error {
	if len(entries) == 0 {
		return nil
	}

	e := entries[0]
	child, err := key.NewPublicChildKey(e.ChildNumber)
	if err != nil {
		return fmt.Errorf("derive entry %d: %v", e.ChildNumber, err)
	}

	if cipher.MustNewPubKey(child.Key) != e.Public {
		return errors.New("stored entries were not derived from the wallet's xpub chain")
	}

	return nil
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

// Loader implements the wallet.Loader interface
type Loader struct{}

// Load loads the xpub wallet from byte slice
func (l Loader) Load(data []byte) (wallet.Wallet, error) {
	w := &Wallet{}
	if err := w.Deserialize(data); err != nil {
		return nil, err
	}

	return w, nil
}

// Creator implements the wallet.Creator interface
type Creator struct{}

// Create creates a xpub wallet
func (c Creator) Create(filename, label, _ string, options wallet.Options) (wallet.Wallet, error) {
	if err := validateOptions(options); err != nil {
		return nil, err
	}

	return NewWallet(
		filename,
		label,
		options.XPub,
		convertOptions(options)...)
}

func validateOptions(options wallet.Options) error {
	if options.Encrypt {
		return wallet.NewError(errors.New("xpub wallet does not support encryption"))
	}

	return nil
}

func convertOptions(options wallet.Options) []wallet.Option {
	var opts []wallet.Option

	if options.Coin != "" {
		opts = append(opts, wallet.OptionCoinType(options.Coin))
	}

	if options.Decoder != nil {
		opts = append(opts, wallet.OptionDecoder(options.Decoder))
	}

	if options.GenerateN > 0 {
		opts = append(opts, wallet.OptionGenerateN(options.GenerateN))
	}

	if options.ScanN > 0 {
		opts = append(opts, wallet.OptionScanN(options.ScanN))
		opts = append(opts, wallet.OptionTransactionsFinder(options.TF))
	}

	if options.Temp {
		opts = append(opts, wallet.OptionTemp(true))
	}

	return opts
}
