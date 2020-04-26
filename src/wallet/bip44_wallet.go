package wallet

import (
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/wallet/bip44wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/SkycoinProject/skycoin/src/wallet/entry"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/sirupsen/logrus"
)

const (
	defaultAccount = 0
)

// Bip44Wallet manages keys using the original Skycoin deterministic
// keypair generator method.
// With this generator, a single chain of addresses is created, each one dependent
// on the previous.
type Bip44Wallet struct {
	*bip44wallet.Bip44WalletNew
}

// LoadBip44Wallet loads wallet from data
func LoadBip44Wallet(data []byte) (Wallet, error) {
	w := &bip44wallet.Bip44WalletNew{}
	if err := w.Deserialize(data); err != nil {
		return nil, err
	}
	return &Bip44Wallet{w}, nil
}

// NewBip44Wallet creates a bip44 wallet
// This function implements the walletCreator for bip44 wallet, which will be registered
// to registeredWalletCreators in wallet.go
func NewBip44Wallet(filename string, opts Options, tf TransactionsFinder) (*Bip44Wallet, error) {
	wltType := opts.Type
	if wltType == "" {
		wltType = WalletTypeBip44
	}

	if wltType != WalletTypeBip44 {
		return nil, NewError(fmt.Errorf("Invalid wallet type %q for creating a bip44 wallet", wltType))
	}

	if opts.Seed == "" {
		return nil, ErrMissingSeed
	}

	if opts.ScanN > 0 && tf == nil {
		return nil, ErrNilTransactionsFinder
	}

	coin := opts.Coin
	if coin == "" {
		coin = meta.CoinTypeSkycoin
	}
	coin, err := meta.ResolveCoinType(string(coin))
	if err != nil {
		return nil, err
	}

	wlt, err := bip44wallet.NewBip44WalletNew(bip44wallet.Bip44WalletCreateOptions{
		Filename:       filename,
		Version:        Version,
		Label:          opts.Label,
		Seed:           opts.Seed,
		SeedPassphrase: opts.SeedPassphrase,
		CoinType:       coin,
		CryptoType:     opts.CryptoType,
	})
	if err != nil {
		return nil, err
	}

	// Create the default bip44 account
	_, err = wlt.NewAccount("default")
	if err != nil {
		return nil, err
	}

	generateN := opts.GenerateN
	if generateN == 0 {
		generateN = 1
	}

	logger.WithFields(logrus.Fields{
		"generateN":  generateN,
		"walletType": wltType,
	}).Infof("Generating addresses for wallet")

	w := &Bip44Wallet{wlt}

	if _, err := w.GenerateAddresses(generateN); err != nil {
		return nil, err
	}

	if opts.ScanN != 0 && coin != meta.CoinTypeSkycoin {
		return nil, errors.New("Wallet scanning is only supported for Skycoin address wallets")
	}

	if opts.ScanN > generateN {
		// Scan for addresses with balances
		logger.WithFields(logrus.Fields{
			"scanN":      opts.ScanN,
			"walletType": wltType,
		}).Info("Scanning addresses for wallet")
		if _, err := w.ScanAddresses(opts.ScanN-generateN, tf); err != nil {
			return nil, err
		}
	}

	if !opts.Encrypt {
		if len(opts.Password) != 0 {
			return nil, ErrMissingEncrypt
		}
		return w, nil
	}

	// Check if the password is provided
	if len(opts.Password) == 0 {
		return nil, ErrMissingPassword
	}

	// Lock the wallet
	if err := w.Lock(opts.Password); err != nil {
		return nil, err
	}

	return w, nil
}

// Clone makes a copy the bip44 wallet
func (w *Bip44Wallet) Clone() Wallet {
	cw := w.Bip44WalletNew.Clone()
	return &Bip44Wallet{Bip44WalletNew: &cw}
}

// CryptoType returns the crypto type that is used for encrypting/decrypting wallet
func (w *Bip44Wallet) CryptoType() crypto.CryptoType {
	return w.Bip44WalletNew.CryptoType()
}

// Lock encrypts the wallet
func (w *Bip44Wallet) Lock(password []byte) error {
	return w.Bip44WalletNew.Lock(password)
}

// Unlock decrypts the wallet
func (w *Bip44Wallet) Unlock(password []byte, f func(w Wallet) error) error {
	wlt, err := w.Bip44WalletNew.Unlock(password)
	if err != nil {
		return err
	}
	defer wlt.Erase()

	return f(&Bip44Wallet{wlt})
}

// newBip44Wallet creates a Bip44Wallet
// func newBip44Wallet(meta meta.Meta) (*Bip44Wallet, error) { //nolint:unparam
// 	return bip44wallet.NewBip44WalletNew(bip44wallet.Bip44WalletCreateOptions{})
// }

// // CopyFrom copies the src wallet to w
// func (w *Bip44Wallet) CopyFrom(src Wallet) {
// 	w.Meta = src.(*Bip44Wallet).Meta.Clone()
// 	w.ExternalEntries = src.(*Bip44Wallet).ExternalEntries.Clone()
// 	w.ChangeEntries = src.(*Bip44Wallet).ChangeEntries.Clone()
// }

// CopyFromRef copies the src wallet with a pointer dereference
// func (w *Bip44Wallet) CopyFromRef(src Wallet) {
// 	*w = *(src.(*Bip44Wallet))
// }

// Erase wipes secret fields in wallet
// func (w *Bip44Wallet) Erase() {
// 	w.Meta.EraseSeeds()
// 	w.ExternalEntries.Erase()
// 	w.ChangeEntries.Erase()
// }

// ToReadable converts the wallet to its readable (serializable) format
// func (w *Bip44Wallet) ToReadable() Readable {
// 	return NewReadableBip44Wallet(w)
// }

// Validate validates the wallet
// func (w *Bip44Wallet) Validate() error {
// if err := w.Meta.Validate(); err != nil {
// 	return err
// }

// walletType := w.Meta.Type()
// if !IsValidWalletType(walletType) {
// 	return ErrInvalidWalletType
// }

// if !w.IsEncrypted() {
// 	// bip44 wallet seeds must be a valid bip39 mnemonic
// 	if s := w.Seed(); s == "" {
// 		return errors.New("seed missing in unencrypted bip44 wallet")
// 	} else if err := bip39.ValidateMnemonic(s); err != nil {
// 		return err
// 	}
// }

// if s := w.Meta[meta.MetaBip44Coin]; s == "" {
// 	return errors.New("bip44Coin missing")
// } else if _, err := strconv.ParseUint(s, 10, 32); err != nil {
// 	return fmt.Errorf("bip44Coin invalid: %v", err)
// }

// if s := w.Meta[meta.MetaLastSeed]; s != "" {
// 	return errors.New("lastSeed should not be in bip44 wallets")
// }
// return nil
// }

// GetEntries returns a copy of all entries held by the wallet
func (w *Bip44Wallet) GetEntries() (entry.Entries, error) {
	eEntries, err := w.ExternalEntries(defaultAccount)
	if err != nil {
		return nil, err
	}

	cEntries, err := w.ChangeEntries(defaultAccount)
	if err != nil {
		return nil, err
	}

	return append(eEntries, cEntries...), nil
}

// EntriesLen returns the number of all entries
func (w *Bip44Wallet) EntriesLen() int {
	el, err := w.ExternalEntriesLen(defaultAccount)
	if err != nil {
		logger.WithError(err).Panic("Get external entries length failed")
		return 0
	}

	cl, err := w.ChangeEntriesLen(defaultAccount)
	if err != nil {
		logger.WithError(err).Panic("Get change entries length failed")
		return 0
	}

	return int(el) + int(cl)
}

// GetEntry returns entry of given address
func (w *Bip44Wallet) GetEntry(a cipher.Address) (entry.Entry, bool) {
	e, ok, err := w.Bip44WalletNew.GetEntry(defaultAccount, a)
	if err != nil {
		logger.WithError(err).Panic("Get entry failed")
		return entry.Entry{}, false
	}

	return e, ok
}

// HasEntry returns true if the wallet has an entry.Entry with a given cipher.Address.
func (w *Bip44Wallet) HasEntry(a cipher.Address) bool {
	_, ok, err := w.Bip44WalletNew.GetEntry(defaultAccount, a)
	if err != nil {
		logger.WithError(err).Panic("HasEntry getting entry failed")
		return false
	}

	return ok
}

// PeekChangeAddress returns a change address, do not generate a new address unless
// the last change address already have transactions associated.
func (w *Bip44Wallet) PeekChangeAddress(tf TransactionsFinder) (cipher.Address, error) {
	// Get the length of the change chain
	len, err := w.ChangeEntriesLen(defaultAccount)
	if err != nil {
		return cipher.Address{}, err
	}

	if len > 0 {
		// Get the last entry of the change chain
		e, err := w.ChangeEntryAt(defaultAccount, len-1)
		if err != nil {
			return cipher.Address{}, err
		}

		// Check whehter the entry has transactions associated
		addr := e.SkycoinAddress()
		hasTxs, err := tf.AddressesActivity([]cipher.Address{addr})
		if err != nil {
			return cipher.Address{}, err
		}

		if !hasTxs[0] {
			return addr, nil
		}
	}

	// Generates a new change address
	addrs, err := w.NewChangeAddresses(defaultAccount, 1)
	if err != nil {
		return cipher.Address{}, err
	}

	return addrs[0].(cipher.Address), nil
}

// generateEntries generates addresses for a change chain (should be 0 or 1) starting from an initial child number.
// func (w *Bip44Wallet) generateEntries(num uint64, changeIdx, initialChildIdx uint32) (entry.Entries, error) {
// 	if w.Meta.IsEncrypted() {
// 		return nil, ErrWalletEncrypted
// 	}

// 	if num > math.MaxUint32 {
// 		return nil, NewError(errors.New("Bip44Wallet.generateEntries num too large"))
// 	}

// 	// Cap `num` in case it would exceed the maximum child index number
// 	if math.MaxUint32-initialChildIdx < uint32(num) {
// 		num = uint64(math.MaxUint32 - initialChildIdx)
// 	}

// 	if num == 0 {
// 		return nil, nil
// 	}

// 	c, err := w.CoinHDNode()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Generate the "account" HDNode. Multiple accounts are not supported; use 0.
// 	account, err := c.Account(0)
// 	if err != nil {
// 		logger.Critical().WithError(err).Error("Failed to derive the bip44 account node")
// 		if bip32.IsImpossibleChildError(err) {
// 			logger.Critical().Error("ImpossibleChild: this seed cannot be used for bip44")
// 		}
// 		return nil, err
// 	}

// 	// Generate the chain parent node
// 	var chain *bip32.PrivateKey
// 	switch changeIdx {
// 	case bip44.ExternalChainIndex:
// 		chain, err = account.External()
// 	case bip44.ChangeChainIndex:
// 		chain, err = account.Change()
// 	default:
// 		err = errors.New("invalid chain index")
// 	}
// 	if err != nil {
// 		logger.Critical().WithError(err).Error("Failed to derive the final bip44 chain node")
// 		if bip32.IsImpossibleChildError(err) {
// 			logger.Critical().Error("ImpossibleChild: this seed cannot be used for bip44")
// 		}
// 		return nil, err
// 	}

// 	// Generate `num` secret keys from the external chain HDNode, skipping any children that
// 	// are invalid (note that this has probability ~2^-128)
// 	var seckeys []*bip32.PrivateKey
// 	var addressIndices []uint32
// 	j := initialChildIdx
// 	for i := uint32(0); i < uint32(num); i++ {
// 		k, err := chain.NewPrivateChildKey(j)

// 		var addErr error
// 		j, addErr = mathutil.AddUint32(j, 1)
// 		if addErr != nil {
// 			logger.Critical().WithError(addErr).WithFields(logrus.Fields{
// 				"num":             num,
// 				"initialChildIdx": initialChildIdx,
// 				"accountIdx":      0,
// 				"changeIdx":       changeIdx,
// 				"childIdx":        j,
// 				"i":               i,
// 			}).Error("childIdx can't be incremented any further")
// 			return nil, errors.New("childIdx can't be incremented any further")
// 		}

// 		if err != nil {
// 			if bip32.IsImpossibleChildError(err) {
// 				logger.Critical().WithError(err).WithFields(logrus.Fields{
// 					"accountIdx": 0,
// 					"changeIdx":  changeIdx,
// 					"childIdx":   j,
// 				}).Error("ImpossibleChild for chain node child element")
// 				continue
// 			} else {
// 				logger.Critical().WithError(err).WithFields(logrus.Fields{
// 					"accountIdx": 0,
// 					"changeIdx":  changeIdx,
// 					"childIdx":   j,
// 				}).Error("NewPrivateChildKey failed unexpectedly")
// 				return nil, err
// 			}
// 		}

// 		seckeys = append(seckeys, k)
// 		addressIndices = append(addressIndices, j-1)
// 	}

// 	entries := make(entry.Entries, len(seckeys))
// 	makeAddress := AddressConstructor(w.Meta)
// 	for i, xprv := range seckeys {
// 		sk := cipher.MustNewSecKey(xprv.Key)
// 		pk := cipher.MustPubKeyFromSecKey(sk)
// 		entries[i] = entry.Entry{
// 			Address:     makeAddress(pk),
// 			Secret:      sk,
// 			Public:      pk,
// 			ChildNumber: addressIndices[i],
// 			Change:      changeIdx,
// 		}
// 	}

// 	return entries, nil
// }

// PeekChangeEntry creates and returns an entry for the change chain.
// If used, the caller the append it with GenerateChangeEntry
// func (w *Bip44Wallet) PeekChangeEntry() (entry.Entry, error) {
// entries, err := w.generateEntries(1, bip44.ChangeChainIndex, nextChildIdx(w.ChangeEntries))
// if err != nil {
// 	return entry.Entry{}, err
// }

// if len(entries) == 0 {
// 	return entry.Entry{}, NewError(errors.New("PeekChangeEntry: no more change addresses"))
// }

// return entries[0], nil
// }

// GenerateChangeEntry creates, appends and returns an entry for the change chain
// func (w *Bip44Wallet) GenerateChangeEntry() (entry.Entry, error) {
// e, err := w.PeekChangeEntry()
// if err != nil {
// 	return entry.Entry{}, err
// }

// w.ChangeEntries = append(w.ChangeEntries, entry.Entries{e}...)

// return w.ChangeEntries[len(w.ChangeEntries)-1], nil
// }

// GenerateAddresses generates addresses on external chain
func (w *Bip44Wallet) GenerateAddresses(num uint64) ([]cipher.Address, error) {
	addrs, err := w.NewExternalAddresses(defaultAccount, uint32(num))
	if err != nil {
		return nil, err
	}

	skyAddrs := make([]cipher.Address, len(addrs))
	for i, a := range addrs {
		skyAddrs[i] = a.(cipher.Address)
	}
	return skyAddrs, nil
}

// GetAddresses returns all external addresses in wallet
func (w *Bip44Wallet) GetAddresses() ([]cipher.Address, error) {
	extEntries, err := w.ExternalEntries(defaultAccount)
	if err != nil {
		return nil, err
	}

	chgEntries, err := w.ChangeEntries(defaultAccount)
	if err != nil {
		return nil, err
	}

	addrs := make([]cipher.Address, len(extEntries)+len(chgEntries))
	for i, e := range append(extEntries, chgEntries...) {
		addrs[i] = e.SkycoinAddress()
	}
	return addrs, nil
}

// ScanAddresses scans ahead N addresses, truncating up to the highest address with any transaction history.
// returns the new generated addresses
func (w *Bip44Wallet) ScanAddresses(scanN uint64, tf TransactionsFinder) ([]cipher.Address, error) {
	if scanN == 0 {
		return nil, nil
	}

	w2 := w.Clone().(*Bip44Wallet)
	// TODO: should not use map, the order is random
	newAddrsFuncs := []func(account, num uint32) ([]cipher.Addresser, error){
		w2.NewExternalAddresses,
		w2.NewChangeAddresses,
	}

	dropEntriesFunc := []func(account, n uint32) error{
		w2.DropExternalLastEntriesN,
		w2.DropChangeLastEntriesN,
	}

	var retAddrs []cipher.Address
	for i, newAddrs := range newAddrsFuncs {
		addrs, err := newAddrs(defaultAccount, uint32(scanN))
		if err != nil {
			return nil, err
		}

		// converts to skycoin addresses
		skyAddrs := convertToSkyAddrs(addrs)

		keepN, err := scanAddressesBip32(skyAddrs, tf)
		if err != nil {
			return nil, err
		}

		retAddrs = append(retAddrs, convertToSkyAddrs(addrs[:keepN])...)

		// drops the last N entreis that without transactions associated
		if err := dropEntriesFunc[i](defaultAccount, uint32(scanN)-keepN); err != nil {
			return nil, err
		}
	}

	*w = *w2

	return retAddrs, nil
}

func convertToSkyAddrs(addrs []cipher.Addresser) []cipher.Address {
	skyAddrs := make([]cipher.Address, len(addrs))
	for i, a := range addrs {
		skyAddrs[i] = a.(cipher.Address)
	}
	return skyAddrs
}

// scanAddressesBip32 implements the address scanning algorithm for bip32
// based (e.g. bip44, xpub) wallets
func scanAddressesBip32(addrs []cipher.Address, tf TransactionsFinder) (uint32, error) {
	if len(addrs) == 0 {
		return 0, nil
	}

	// Find if these addresses had any activity
	active, err := tf.AddressesActivity(addrs)
	if err != nil {
		return 0, err
	}

	// Check activity from the last one until we find the address that has activity
	var keepNum uint32
	for i := len(active) - 1; i >= 0; i-- {
		if active[i] {
			keepNum = uint32(i + 1)
			break
		}
	}

	return keepNum, nil
}

// Fingerprint returns a unique ID fingerprint for this wallet, composed of its initial address
// and wallet type
func (w *Bip44Wallet) Fingerprint() string {
	addr := ""
	entries, err := w.ExternalEntries(defaultAccount)
	if err != nil {
		logger.WithError(err).Panic("Fingerprint get external entries failed")
		return ""
	}

	if len(entries) == 0 {
		if !w.IsEncrypted() {
			addrs, err := w.NewExternalAddresses(defaultAccount, 1)
			if err != nil {
				logger.WithError(err).Panic("Fingerprint failed to generate initial entry for empty wallet")
			}
			addr = addrs[0].String()
		}
	} else {
		addr = entries[0].Address.String()
	}
	return fmt.Sprintf("%s-%s", w.Type(), addr)
}

// ReadableBip44Wallet used for [de]serialization of a deterministic wallet
// type ReadableBip44Wallet struct {
// 	meta.Meta       `json:"meta"`
// 	ReadableEntries `json:"entries"`
// }

// // LoadReadableBip44Wallet loads a deterministic wallet from disk
// func LoadReadableBip44Wallet(wltFile string) (*ReadableBip44Wallet, error) {
// 	var rw ReadableBip44Wallet
// 	if err := file.LoadJSON(wltFile, &rw); err != nil {
// 		return nil, err
// 	}
// 	if rw.Type() != WalletTypeBip44 {
// 		return nil, ErrInvalidWalletType
// 	}
// 	return &rw, nil
// }

// // NewReadableBip44Wallet creates readable wallet
// func NewReadableBip44Wallet(w *Bip44Wallet) *ReadableBip44Wallet {
// 	return &ReadableBip44Wallet{
// 		Meta:            w.Meta.Clone(),
// 		ReadableEntries: newReadableEntries(w.GetEntries(), w.Meta.Coin(), w.Meta.Type()),
// 	}
// }

// // ToWallet convert readable wallet to Wallet
// func (rw *ReadableBip44Wallet) ToWallet() (Wallet, error) {
// 	w := &Bip44Wallet{
// 		Meta: rw.Meta.Clone(),
// 	}

// 	if err := w.Validate(); err != nil {
// 		err := fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
// 		logger.WithError(err).Error("ReadableBip44Wallet.ToWallet Validate failed")
// 		return nil, err
// 	}

// 	ets, err := rw.ReadableEntries.toWalletEntries(w.Meta.Coin(), w.Meta.Type(), w.Meta.IsEncrypted())
// 	if err != nil {
// 		logger.WithError(err).Error("ReadableBip44Wallet.ToWallet ReadableEntries.toWalletEntries failed")
// 		return nil, err
// 	}

// 	// Split the single array of entries into separate external and change chains,
// 	// for easier internal management
// 	for _, e := range ets {
// 		switch e.Change {
// 		case bip44.ExternalChainIndex:
// 			w.ExternalEntries = append(w.ExternalEntries, e)
// 		case bip44.ChangeChainIndex:
// 			w.ChangeEntries = append(w.ChangeEntries, e)
// 		default:
// 			logger.Panicf("invalid change value %d", e.Change)
// 		}
// 	}

// 	// Sort childNumber low to high
// 	sort.Slice(w.ExternalEntries, func(i, j int) bool {
// 		return w.ExternalEntries[i].ChildNumber < w.ExternalEntries[j].ChildNumber
// 	})
// 	sort.Slice(w.ChangeEntries, func(i, j int) bool {
// 		return w.ChangeEntries[i].ChildNumber < w.ChangeEntries[j].ChildNumber
// 	})

// 	return w, err
// }
