package wallet

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/scrypt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/sha256xor"
)

type cryptor interface {
	Encrypt(w *Wallet, password []byte) error
	Decrypt(w *Wallet, password []byte) (*Wallet, error)
}

// CryptoType represents the type of crypto name
type CryptoType string

// StrToCryptoType converts string to CryptoType
func StrToCryptoType(s string) (CryptoType, error) {
	switch CryptoType(s) {
	case CryptoTypeSha256Xor:
		return CryptoTypeSha256Xor, nil
	case CryptoTypeScryptChacha20poly1305:
		return CryptoTypeScryptChacha20poly1305, nil
	default:
		return "", errors.New("unknow crypto type")
	}
}

// Crypto types
const (
	CryptoTypeSha256Xor              = CryptoType("sha256-xor")
	CryptoTypeScryptChacha20poly1305 = CryptoType("scrypt-chacha20poly1305")
)

// Scrypt paraments
const (
	scryptN       = 1 << 15
	scryptR       = 8
	scryptP       = 1
	scryptKeyLen  = 32
	scryptSaltLen = 8

	additionalDataLen = 32
)

// cryptoTable records all supported wallet crypto methods
// If want to support new crypto methods, register here.
var cryptoTable = map[CryptoType]cryptor{
	CryptoTypeSha256Xor:              &sha256xorCrypto{},
	CryptoTypeScryptChacha20poly1305: newScryptChacha20poly1305Crypto(scryptN, scryptR, scryptP, scryptKeyLen, scryptSaltLen),
}

// ErrAuthenticationFailed wraps the error of decryption.
type ErrAuthenticationFailed struct {
	err error
}

func (e ErrAuthenticationFailed) Error() string {
	return e.err.Error()
}

// getCrypto gets crypto of given type
func getCrypto(cryptoType CryptoType) (cryptor, error) {
	c, ok := cryptoTable[cryptoType]
	if !ok {
		return nil, fmt.Errorf("can not find crypto %v in crypto table", cryptoType)
	}

	return c, nil
}

// scryptParams records the scrypt paramenters
type scryptParams struct {
	N       int    `json:"N"`       // scrypt N parament
	R       int    `json:"r"`       // scrypt r parament
	P       int    `json:"p"`       // scrypt p parament
	KeyLen  int    `json:"key_len"` // scrypt keyLen parament
	SaltLen int    `json:"-"`       // scrypt salt length
	Salt    []byte `json:"salt"`    // scrypt salt parament
}

// chacha20poly1305Param records the chacha20poly1305 paramenters
type chacha20poly1305Params struct {
	Nonce          []byte `json:"nonce"`           // chacha20poly1305 nonce parament
	AdditionalData []byte `json:"additional_data"` // chacha20poly1305 additionalData parament
}

// authenticated records the scrypt and chacha20poly1305 paramenters
type authenticated struct {
	scryptParams
	chacha20poly1305Params
}

func (a *authenticated) Serialize() ([]byte, error) {
	d, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (a *authenticated) Deserialize(d []byte) error {
	if a == nil {
		return errors.New("authenticated is nil")
	}

	return json.Unmarshal(d, a)
}

// scryptChacha20poly1305Crypto use scrypt and chacha20poly1305 to encrypt/decrypt wallet
type scryptChacha20poly1305Crypto struct {
	scryptParams // scrypt paramenters
}

func newScryptChacha20poly1305Crypto(N, r, p, keyLen, saltLen int) *scryptChacha20poly1305Crypto {
	return &scryptChacha20poly1305Crypto{
		scryptParams: scryptParams{
			N:       N,
			R:       r,
			P:       p,
			KeyLen:  keyLen,
			SaltLen: saltLen,
		},
	}
}

// Encrypt encrypts wallet using scrypt and chacha20poly1305 with given password
func (s *scryptChacha20poly1305Crypto) Encrypt(w *Wallet, password []byte) error {
	if len(password) == 0 {
		return ErrMissingPassword
	}

	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	wlt := w.clone()

	// Derives password key
	salt := cipher.RandByte(s.SaltLen)
	dk, err := scrypt.Key(password, salt, s.N, s.R, s.P, s.KeyLen)
	if err != nil {
		return err
	}

	aead, err := chacha20poly1305.New(dk)
	if err != nil {
		return err
	}

	nonce := cipher.RandByte(chacha20poly1305.NonceSize)
	// Generates additional data
	ad := cipher.RandByte(additionalDataLen)
	// Encrypt the seed
	ss := aead.Seal(nil, nonce, []byte(wlt.seed()), ad)

	wlt.setEncryptedSeed(base64.StdEncoding.EncodeToString(ss))

	// Encrypt the last seed
	sls := aead.Seal(nil, nonce, []byte(wlt.lastSeed()), ad)

	wlt.setEncryptedLastSeed(base64.StdEncoding.EncodeToString(sls))

	// encrypt private keys in entries
	for i, e := range wlt.Entries {
		se := aead.Seal(nil, nonce, e.Secret[:], ad)

		// Set the encrypted seckey value
		wlt.Entries[i].EncryptedSecret = base64.StdEncoding.EncodeToString(se)
	}

	// Sets wallet as encrypted
	wlt.setEncrypted(true)
	// Sets authenticated data
	wlt.setAuthenticated(authenticated{
		scryptParams: scryptParams{
			N:      s.N,
			R:      s.R,
			P:      s.P,
			KeyLen: s.KeyLen,
			Salt:   salt,
		},
		chacha20poly1305Params: chacha20poly1305Params{
			Nonce:          nonce,
			AdditionalData: ad,
		},
	})
	// Sets crypto type
	wlt.setCryptoType(CryptoTypeScryptChacha20poly1305)
	// Wipes the sercet fields in wlt
	wlt.erase()
	// Wipes the secret fields in w
	w.erase()

	// Replace the unlocked w with locked wlt
	*w = *wlt

	return nil
}

func (s *scryptChacha20poly1305Crypto) Decrypt(w *Wallet, password []byte) (*Wallet, error) {
	if !w.IsEncrypted() {
		return nil, ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return nil, ErrMissingPassword
	}

	if w.cryptoType() != CryptoTypeScryptChacha20poly1305 {
		return nil, ErrWrongCryptoType
	}

	// Gets the authenticated data
	auth, err := w.authenticated()
	if err != nil {
		return nil, err
	}

	if auth == nil {
		return nil, ErrMissingAuthenticated
	}

	dk, err := scrypt.Key(password, auth.Salt, auth.N, auth.R, auth.P, auth.KeyLen)
	if err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.New(dk)
	if err != nil {
		return nil, err
	}

	wlt := w.clone()

	// Base64 decodes the seed string
	ss, err := base64.StdEncoding.DecodeString(wlt.encryptedSeed())
	if err != nil {
		return nil, err
	}

	// Decrypts the seed
	sd, err := aead.Open(nil, auth.Nonce, ss, auth.AdditionalData)
	if err != nil {
		return nil, ErrAuthenticationFailed{err}
	}
	wlt.setSeed(string(sd))

	// Base64 decodes the last seed string
	sls, err := base64.StdEncoding.DecodeString(wlt.encryptedLastSeed())
	if err != nil {
		return nil, err
	}

	// decrypt lastSeed
	ls, err := aead.Open(nil, auth.Nonce, sls, auth.AdditionalData)
	if err != nil {
		return nil, ErrAuthenticationFailed{err}
	}
	wlt.setLastSeed(string(ls))

	// decrypt the entries
	for i := range wlt.Entries {
		ssecret, err := base64.StdEncoding.DecodeString(wlt.Entries[i].EncryptedSecret)
		if err != nil {
			return nil, err
		}

		sk, err := aead.Open(nil, auth.Nonce, ssecret, auth.AdditionalData)
		if err != nil {
			return nil, ErrAuthenticationFailed{err}
		}
		copy(wlt.Entries[i].Secret[:], sk[:])
	}
	wlt.setEncrypted(false)

	return wlt, nil
}

// sha256xorCrypto use sha256xor crypto to encrypt/decrypt wallet
type sha256xorCrypto struct{}

func (s *sha256xorCrypto) Encrypt(w *Wallet, password []byte) error {
	if len(password) == 0 {
		return ErrMissingPassword
	}

	if w.IsEncrypted() {
		return ErrWalletEncrypted
	}

	crypto := sha256xor.New()

	wlt := w.clone()

	// Encrypt the seed
	ss, err := crypto.Encrypt([]byte(wlt.seed()), password)
	if err != nil {
		return err
	}

	wlt.setEncryptedSeed(base64.StdEncoding.EncodeToString(ss))

	// Encrypt the last seed
	sls, err := crypto.Encrypt([]byte(wlt.lastSeed()), password)
	if err != nil {
		return err
	}

	wlt.setEncryptedLastSeed(base64.StdEncoding.EncodeToString(sls))

	// encrypt private keys in entries
	for i, e := range wlt.Entries {
		se, err := crypto.Encrypt(e.Secret[:], password)
		if err != nil {
			return err
		}

		// Set the encrypted seckey value
		wlt.Entries[i].EncryptedSecret = base64.StdEncoding.EncodeToString(se)
	}

	// Sets wallet as encrypted
	wlt.setEncrypted(true)
	// Sets wallet crypto type
	wlt.setCryptoType(CryptoTypeSha256Xor)
	// Wipes the sercet fields in wlt
	wlt.erase()
	// Wipes the secret fields in w
	w.erase()

	// Replace the unlocked w with locked wlt
	*w = *wlt
	return nil
}

func (s *sha256xorCrypto) Decrypt(w *Wallet, password []byte) (*Wallet, error) {
	if !w.IsEncrypted() {
		return nil, ErrWalletNotEncrypted
	}

	if len(password) == 0 {
		return nil, ErrMissingPassword
	}

	ct := w.cryptoType()
	if ct != CryptoTypeSha256Xor {
		return nil, ErrWrongCryptoType
	}

	// Gets the crypto
	crypto := sha256xor.New()

	wlt := w.clone()

	// Base64 decodes the seed string
	ss, err := base64.StdEncoding.DecodeString(wlt.encryptedSeed())
	if err != nil {
		return nil, err
	}

	// Decrypts the seed
	sd, err := crypto.Decrypt(ss, password)
	if err != nil {
		return nil, ErrAuthenticationFailed{err}
	}
	wlt.setSeed(string(sd))

	// Base64 decodes the last seed string
	sls, err := base64.StdEncoding.DecodeString(wlt.encryptedLastSeed())
	if err != nil {
		return nil, err
	}

	// decrypt lastSeed
	ls, err := crypto.Decrypt(sls, password)
	if err != nil {
		return nil, ErrAuthenticationFailed{err}
	}
	wlt.setLastSeed(string(ls))

	// decrypt the entries
	for i := range wlt.Entries {
		ssecret, err := base64.StdEncoding.DecodeString(wlt.Entries[i].EncryptedSecret)
		if err != nil {
			return nil, err
		}

		sk, err := crypto.Decrypt(ssecret, password)
		if err != nil {
			return nil, ErrAuthenticationFailed{err}
		}
		copy(wlt.Entries[i].Secret[:], sk[:])
	}
	wlt.setEncrypted(false)

	return wlt, nil
}
