// Package wallet src/wallet/meta.go c1-coin-wallet
package wallet

import (
	"errors"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
)

// wallet meta fields
const (
	MetaVersion        = "version"        // wallet version
	MetaFilename       = "filename"       // wallet file name
	MetaLabel          = "label"          // wallet label
	MetaTimestamp      = "tm"             // the timestamp when creating the wallet
	MetaType           = "type"           // wallet type
	MetaCoin           = "coin"           // coin type
	MetaEncrypted      = "encrypted"      // whether the wallet is encrypted
	MetaCryptoType     = "cryptoType"     // encryption/decryption type
	MetaSeed           = "seed"           // wallet seed
	MetaLastSeed       = "lastSeed"       // seed for generating next address [deterministic wallets]
	MetaSecrets        = "secrets"        // secrets which records the encrypted seeds and secrets of address entries
	MetaBip44Coin      = "bip44Coin"      // bip44 coin type
	MetaAccountsHash   = "accountsHash"   // accounts hash
	MetaSeedPassphrase = "seedPassphrase" // seed passphrase [bip44 wallets]
	MetaXPub           = "xpub"           // xpub key [xpub wallets]
	MetaXPubChain      = "xpubChain"      // chain derived from an account xpub before the addresses [xpub wallets]
	MetaTemp           = "temp"           // whether the wallet is a temporary wallet
)

//const (
//	// CoinTypeSkycoin skycoin type
//	CoinTypeSkycoin CoinType = "skycoin"
//	// CoinTypeBitcoin bitcoin type
//	CoinTypeBitcoin CoinType = "bitcoin"
//)

// CoinType represents the wallet coin type, which refers to the pubkey2addr method used
//type CoinType string

// Meta holds wallet metadata
type Meta map[string]string

// Clone make an copy of the Meta
func (m Meta) Clone() Meta {
	mm := make(Meta, len(m))
	for k, v := range m {
		mm[k] = v
	}
	return mm
}

// EraseSeeds wipes the seed and last seed
func (m Meta) EraseSeeds() {
	m.SetSeed("")
}

// Find returns a key value from the metadata map
func (m Meta) Find(k string) string {
	return m[k]
}

// Type gets the wallet type
func (m Meta) Type() string {
	return m[MetaType]
}

// Version gets the wallet version
func (m Meta) Version() string {
	return m[MetaVersion]
}

// SetVersion sets the wallet version
func (m Meta) SetVersion(v string) {
	m[MetaVersion] = v
}

// Filename gets the wallet filename
func (m Meta) Filename() string {
	return m[MetaFilename]
}

// SetFilename sets the wallet filename
func (m Meta) SetFilename(fn string) {
	m[MetaFilename] = fn
}

// Label gets the wallet label
func (m Meta) Label() string {
	return m[MetaLabel]
}

// SetLabel sets the wallet label
func (m Meta) SetLabel(label string) {
	m[MetaLabel] = label
}

// LastSeed returns the last seed
func (m Meta) LastSeed() string {
	return m[MetaLastSeed]
}

// SetLastSeed sets or updates the last seed
func (m Meta) SetLastSeed(lseed string) {
	m[MetaLastSeed] = lseed
}

// Seed returns the seed
func (m Meta) Seed() string {
	return m[MetaSeed]
}

// SetSeed sets the seed
func (m Meta) SetSeed(seed string) {
	m[MetaSeed] = seed
}

// SeedPassphrase returns the seed passphrase
func (m Meta) SeedPassphrase() string {
	return m[MetaSeedPassphrase]
}

// SetSeedPassphrase sets the seed passphrase
func (m Meta) SetSeedPassphrase(p string) {
	m[MetaSeedPassphrase] = p
}

// Coin returns the wallet's coin type
func (m Meta) Coin() CoinType {
	return CoinType(m[MetaCoin])
}

// SetCoin sets the wallet's coin type
func (m Meta) SetCoin(ct CoinType) {
	m[MetaCoin] = string(ct)
}

// Bip44Coin returns the bip44 coin type, or nil if it is not set.
// Validate rejects an unparseable value, so one reads as unset here.
func (m Meta) Bip44Coin() *bip44.CoinType {
	c, ok := m[MetaBip44Coin]
	if !ok {
		return nil
	}
	x, err := strconv.ParseUint(c, 10, 32)
	if err != nil {
		return nil
	}
	t := bip44.CoinType(x)

	return &t
}

// SetBip44Coin sets the bip44 coin type code
func (m Meta) SetBip44Coin(ct bip44.CoinType) {
	m[MetaBip44Coin] = strconv.FormatUint(uint64(ct), 10)
}

func (m Meta) setIsEncrypted(encrypt bool) {
	m[MetaEncrypted] = strconv.FormatBool(encrypt)
}

// SetEncrypted sets encryption fields
func (m Meta) SetEncrypted(cryptoType crypto.CryptoType, encryptedSecrets string) {
	m.setCryptoType(cryptoType)
	m.setSecrets(encryptedSecrets)
	m.setIsEncrypted(true)
}

// SetDecrypted unsets encryption fields
func (m Meta) SetDecrypted() {
	m.setIsEncrypted(false)
	m.setSecrets("")
	delete(m, MetaSecrets)
}

// Encrypted checks whether the wallet is encrypted.
func (m Meta) Encrypted() bool {
	encStr, ok := m[MetaEncrypted]
	if !ok {
		return false
	}

	b, err := strconv.ParseBool(encStr)
	if err != nil {
		// Validate rejects this when a wallet loads. Read it as encrypted so
		// nothing treats the secrets as plain text.
		return true
	}
	return b
}

func (m Meta) setCryptoType(tp crypto.CryptoType) {
	m[MetaCryptoType] = string(tp)
}

// CryptoType returns the encryption type
func (m Meta) CryptoType() crypto.CryptoType {
	return crypto.CryptoType(m[MetaCryptoType])
}

// SetCryptoType sets the encryption type
func (m Meta) SetCryptoType(ct crypto.CryptoType) {
	m[MetaCryptoType] = string(ct)
}

// Secrets returns the encrypted wallet secrets
func (m Meta) Secrets() string {
	return m[MetaSecrets]
}

func (m Meta) setSecrets(s string) {
	m[MetaSecrets] = s
}

// Timestamp returns the timestamp
func (m Meta) Timestamp() int64 {
	// Intentionally ignore the error when parsing the timestamp,
	// if it isn't valid or is missing it will be set to 0.
	// Also, this value is validated by wallet.validate()
	x, _ := strconv.ParseInt(m[MetaTimestamp], 10, 64) //nolint:errcheck
	return x
}

// SetTimestamp sets the timestamp
func (m Meta) SetTimestamp(t int64) {
	m[MetaTimestamp] = strconv.FormatInt(t, 10)
}

// SetXPub sets xpub
func (m Meta) SetXPub(xpub string) {
	m[MetaXPub] = xpub
}

// XPub returns the wallet's configured XPub key
func (m Meta) XPub() string {
	return m[MetaXPub]
}

// Validate validates the meta data
func (m Meta) Validate() error {
	if fn := m[MetaFilename]; fn == "" {
		return errors.New("filename not set")
	}

	if tm := m[MetaTimestamp]; tm != "" {
		_, err := strconv.ParseInt(tm, 10, 64)
		if err != nil {
			return errors.New("invalid timestamp")
		}
	}

	_, ok := m[MetaType]
	if !ok {
		return errors.New("type field not set")
	}

	if coinType := m[MetaCoin]; coinType == "" {
		return errors.New("coin field not set")
	}

	var encrypted bool
	if encStr, ok := m[MetaEncrypted]; ok {
		// validate the encrypted value
		var err error
		encrypted, err = strconv.ParseBool(encStr)
		if err != nil {
			return errors.New("encrypted field is not a valid bool")
		}
	}

	if encrypted {
		cryptoType, ok := m[MetaCryptoType]
		if !ok {
			return errors.New("crypto type field not set")
		}

		if _, err := crypto.GetCrypto(crypto.CryptoType(cryptoType)); err != nil {
			return errors.New("unknown crypto type")
		}

		if s := m[MetaSecrets]; s == "" {
			return errors.New("wallet is encrypted, but secrets field not set")
		}

		if s := m[MetaSeed]; s != "" {
			return errors.New("seed should not be visible in encrypted wallets")
		}

		if s := m[MetaLastSeed]; s != "" {
			return errors.New("lastSeed should not be visible in encrypted wallets")
		}
	} else {
		if s := m[MetaSecrets]; s != "" {
			return errors.New("secrets should not be in unencrypted wallets")
		}
	}

	// Validate the bip44 coin type here, Meta.Bip44Coin reads a malformed value as unset.
	if c, ok := m[MetaBip44Coin]; ok {
		if _, err := strconv.ParseUint(c, 10, 32); err != nil {
			return errors.New("invalid bip44Coin value")
		}
	}
	return nil
}

// ResolveCoinType normalizes a coin type string to a CoinType constant
func ResolveCoinType(s string) (CoinType, error) {
	switch strings.ToLower(s) {
	case "sky", "skycoin":
		return CoinTypeSkycoin, nil
	case "btc", "bitcoin":
		return CoinTypeBitcoin, nil
	case "bitcoin-segwit":
		return CoinTypeBitcoinSegwit, nil
	case "eth", "ethereum":
		return CoinTypeEthereum, nil
	default:
		return CoinType(""), errors.New("invalid coin type")
	}
}

// SetTemp sets temp
func (m Meta) SetTemp(temp bool) {
	if temp {
		m[MetaTemp] = "true"
	}
}

// Temp returns whether the wallet is a temporary wallet
func (m Meta) Temp() bool {
	return m[MetaTemp] == "true"
}
