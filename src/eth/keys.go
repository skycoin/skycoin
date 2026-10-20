package eth

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip32"
	"github.com/skycoin/skycoin/src/cipher/bip39"
)

// accountPath is the BIP44 external chain of account 0 for coin type 60.
const accountPath = "m/44'/60'/0'/0"

// Key is the secret key of external address i of a BIP39 mnemonic.
func Key(mnemonic, passphrase string, i uint32) (cipher.SecKey, error) {
	seed, err := bip39.NewSeed(mnemonic, passphrase)
	if err != nil {
		return cipher.SecKey{}, err
	}
	chain, err := bip32.NewPrivateKeyFromPath(seed, accountPath)
	if err != nil {
		return cipher.SecKey{}, err
	}
	child, err := chain.NewPrivateChildKey(i)
	if err != nil {
		return cipher.SecKey{}, fmt.Errorf("address %d: %w", i, err)
	}
	return cipher.NewSecKey(child.Key)
}

// AddressOf is the address of a secret key.
func AddressOf(sec cipher.SecKey) (Address, error) {
	pk, err := cipher.PubKeyFromSecKey(sec)
	if err != nil {
		return Address{}, err
	}
	return AddressFromPubKey(pk)
}
