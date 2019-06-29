package wallet

import (
	"errors"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
)

// Meta holds wallet metadata
type Meta map[string]string

func (m Meta) clone() Meta {
	mm := make(Meta, len(m))
	for k, v := range m {
		mm[k] = v
	}
	return mm
}

// erase wipes the seed and last seed
func (m Meta) eraseSeeds() {
	m.setSeed("")
	m.setLastSeed("")
}

// validate validates the wallet
func (m Meta) validate() error {
	if fn := m[metaFilename]; fn == "" {
		return errors.New("filename not set")
	}

	if tm := m[metaTimestamp]; tm != "" {
		_, err := strconv.ParseInt(tm, 10, 64)
		if err != nil {
			return errors.New("invalid timestamp")
		}
	}

	walletType, ok := m[metaType]
	if !ok {
		return errors.New("type field not set")
	}
	if !IsValidWalletType(walletType) {
		return ErrInvalidWalletType
	}

	if coinType := m[metaCoin]; coinType == "" {
		return errors.New("coin field not set")
	}

	var isEncrypted bool
	if encStr, ok := m[metaEncrypted]; ok {
		// validate the encrypted value
		var err error
		isEncrypted, err = strconv.ParseBool(encStr)
		if err != nil {
			return errors.New("encrypted field is not a valid bool")
		}
	}

	if isEncrypted {
		cryptoType, ok := m[metaCryptoType]
		if !ok {
			return errors.New("crypto type field not set")
		}

		if _, err := getCrypto(CryptoType(cryptoType)); err != nil {
			return errors.New("unknown crypto type")
		}

		if s := m[metaSecrets]; s == "" {
			return errors.New("wallet is encrypted, but secrets field not set")
		}

		if s := m[metaSeed]; s != "" {
			return errors.New("seed should not be visible in encrypted wallets")
		}

		if s := m[metaLastSeed]; s != "" {
			return errors.New("lastSeed should not be visible in encrypted wallets")
		}
	}

	switch walletType {
	case WalletTypeCollection:
		if s := m[metaSeed]; s != "" {
			return errors.New("seed should not be in collection wallets")
		}

		if s := m[metaLastSeed]; s != "" {
			return errors.New("lastSeed should not be in collection wallets")
		}
	case WalletTypeDeterministic:
		// checks if the secrets field is empty
		if !isEncrypted {
			if s := m[metaSeed]; s == "" {
				return errors.New("seed missing in unencrypted deterministic wallet")
			}

			if s := m[metaLastSeed]; s == "" {
				return errors.New("lastSeed missing in unencrypted deterministic wallet")
			}
		}
	default:
		return ErrInvalidWalletType
	}

	return nil
}

// Find returns a key value from the metadata map
func (m Meta) Find(k string) string {
	return m[k]
}

// Type gets the wallet type
func (m Meta) Type() string {
	return m[metaType]
}

// Version gets the wallet version
func (m Meta) Version() string {
	return m[metaVersion]
}

// SetVersion sets the wallet version
func (m Meta) SetVersion(v string) {
	m[metaVersion] = v
}

// Filename gets the wallet filename
func (m Meta) Filename() string {
	return m[metaFilename]
}

// SetFilename sets the wallet filename
func (m Meta) SetFilename(fn string) {
	m[metaFilename] = fn
}

// Label gets the wallet label
func (m Meta) Label() string {
	return m[metaLabel]
}

// SetLabel sets the wallet label
func (m Meta) SetLabel(label string) {
	m[metaLabel] = label
}

// LastSeed returns the last seed
func (m Meta) LastSeed() string {
	return m[metaLastSeed]
}

func (m Meta) setLastSeed(lseed string) {
	m[metaLastSeed] = lseed
}

// Seed returns the seed
func (m Meta) Seed() string {
	return m[metaSeed]
}

func (m Meta) setSeed(seed string) {
	m[metaSeed] = seed
}

// Coin returns the wallet's coin type
func (m Meta) Coin() CoinType {
	return CoinType(m[metaCoin])
}

// SetCoin sets the wallet's coin type
func (m Meta) SetCoin(ct CoinType) {
	m[metaCoin] = string(ct)
}

func (m Meta) setIsEncrypted(encrypt bool) {
	m[metaEncrypted] = strconv.FormatBool(encrypt)
}

// SetEncrypted sets encryption fields
func (m Meta) SetEncrypted(cryptoType CryptoType, encryptedSecrets string) {
	m.setCryptoType(cryptoType)
	m.setSecrets(encryptedSecrets)
	m.setIsEncrypted(true)
}

// SetDecrypted unsets encryption fields
func (m Meta) SetDecrypted() {
	m.setIsEncrypted(false)
	m.setSecrets("")
	m.setCryptoType("")
}

// IsEncrypted checks whether the wallet is encrypted.
func (m Meta) IsEncrypted() bool {
	encStr, ok := m[metaEncrypted]
	if !ok {
		return false
	}

	b, err := strconv.ParseBool(encStr)
	if err != nil {
		// This can't happen, the meta.encrypted value is either set by
		// setEncrypted() method or converted in ReadableWallet.toWallet().
		// toWallet() method will throw error if the meta.encrypted string is invalid.
		logger.Critical().WithError(err).Error("parse wallet.meta.encrypted string failed")
		return false
	}
	return b
}

func (m Meta) setCryptoType(tp CryptoType) {
	m[metaCryptoType] = string(tp)
}

// CryptoType returns the encryption type
func (m Meta) CryptoType() CryptoType {
	return CryptoType(m[metaCryptoType])
}

// Secrets returns the encrypted wallet secrets
func (m Meta) Secrets() string {
	return m[metaSecrets]
}

func (m Meta) setSecrets(s string) {
	m[metaSecrets] = s
}

// Timestamp returns the timestamp
func (m Meta) Timestamp() int64 {
	// Intentionally ignore the error when parsing the timestamp,
	// if it isn't valid or is missing it will be set to 0.
	// Also, this value is validated by wallet.validate()
	x, _ := strconv.ParseInt(m[metaTimestamp], 10, 64) //nolint:errcheck
	return x
}

// SetTimestamp sets the timestamp
func (m Meta) SetTimestamp(t int64) {
	m[metaTimestamp] = strconv.FormatInt(t, 10)
}

// AddressConstructor returns a function to create a cipher.Addresser from a cipher.PubKey
func (m Meta) AddressConstructor() func(cipher.PubKey) cipher.Addresser {
	switch m.Coin() {
	case CoinTypeSkycoin:
		return func(pk cipher.PubKey) cipher.Addresser {
			return cipher.AddressFromPubKey(pk)
		}
	case CoinTypeBitcoin:
		return func(pk cipher.PubKey) cipher.Addresser {
			return cipher.BitcoinAddressFromPubKey(pk)
		}
	default:
		logger.Panicf("Invalid wallet coin type %q", m.Coin())
		return nil
	}
}
