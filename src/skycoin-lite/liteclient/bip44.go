package liteclient

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip32"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/cipher/bip44"
)

// ErrBip44Range is returned for a chain other than external (0) or change (1),
// or an address index range outside 0 to 2^31-1.
var ErrBip44Range = errors.New("bip44 chain or address index out of range")

// Bip44Address is an address on a bip44 chain. Secret is empty when the address
// was derived from an xpub.
type Bip44Address struct {
	Index   uint32
	Secret  string
	Public  string
	Address string
}

// bip44Account derives m/44'/8000'/0' the way src/wallet/bip44wallet does,
// with the bip39 seed passphrase going into the PBKDF2 salt.
func bip44Account(mnemonic, passphrase string) (*bip44.Account, error) {
	seed, err := bip39.NewSeed(mnemonic, passphrase)
	if err != nil {
		return nil, err
	}

	c, err := bip44.NewCoin(seed, bip44.CoinTypeSkycoin)
	if err != nil {
		return nil, err
	}

	return c.Account(0)
}

func checkBip44Range(chain, start, num int) error {
	if (chain != int(bip44.ExternalChainIndex) && chain != int(bip44.ChangeChainIndex)) ||
		start < 0 || num < 0 || uint64(start)+uint64(num) > uint64(bip32.FirstHardenedChild) {
		return ErrBip44Range
	}
	return nil
}

// Bip44AccountXPub returns the account xpub (m/44'/8000'/0') of a bip39
// mnemonic, the same string `skycoin-cli walletXPub` reports.
func Bip44AccountXPub(mnemonic, passphrase string) (string, error) {
	a, err := bip44Account(mnemonic, passphrase)
	if err != nil {
		return "", err
	}

	pub, err := a.PublicKey()
	if err != nil {
		return "", err
	}

	return pub.String(), nil
}

// Bip44Addresses derives num addresses with their secret keys from index start
// of a chain of account 0, external (0) or change (1).
func Bip44Addresses(mnemonic, passphrase string, chain, start, num int) ([]Bip44Address, error) {
	if err := checkBip44Range(chain, start, num); err != nil {
		return nil, err
	}

	a, err := bip44Account(mnemonic, passphrase)
	if err != nil {
		return nil, err
	}

	chainKey, err := a.NewPrivateChildKey(uint32(chain)) //nolint:gosec // checked above
	if err != nil {
		return nil, err
	}

	addresses := make([]Bip44Address, num)
	for i := range addresses {
		index := uint32(start + i) //nolint:gosec // checked above
		k, err := chainKey.NewPrivateChildKey(index)
		if err != nil {
			return nil, err
		}

		sk, err := cipher.NewSecKey(k.Key)
		if err != nil {
			return nil, err
		}

		pk, err := cipher.PubKeyFromSecKey(sk)
		if err != nil {
			return nil, err
		}

		addresses[i] = Bip44Address{
			Index:   index,
			Secret:  sk.Hex(),
			Public:  pk.Hex(),
			Address: cipher.AddressFromPubKey(pk).String(),
		}
	}

	return addresses, nil
}

// Bip44XPubAddresses derives num addresses from index start of a chain of the
// account an xpub belongs to, without any secret keys.
func Bip44XPubAddresses(xpub string, chain, start, num int) ([]Bip44Address, error) {
	if err := checkBip44Range(chain, start, num); err != nil {
		return nil, err
	}

	account, err := bip32.DeserializeEncodedPublicKey(xpub)
	if err != nil {
		return nil, err
	}

	chainKey, err := account.NewPublicChildKey(uint32(chain)) //nolint:gosec // checked above
	if err != nil {
		return nil, err
	}

	addresses := make([]Bip44Address, num)
	for i := range addresses {
		index := uint32(start + i) //nolint:gosec // checked above
		k, err := chainKey.NewPublicChildKey(index)
		if err != nil {
			return nil, err
		}

		pk, err := cipher.NewPubKey(k.Key)
		if err != nil {
			return nil, err
		}

		addresses[i] = Bip44Address{
			Index:   index,
			Public:  pk.Hex(),
			Address: cipher.AddressFromPubKey(pk).String(),
		}
	}

	return addresses, nil
}
