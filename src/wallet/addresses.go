package wallet

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
)

// CreateAddresses genCount addresses deterministically from seed.  coinType is either CoinTypeBitcoin or CoinTypeSkycoin.
// hideSecretKey will hide the secret key from the output.
func CreateAddresses(coinType CoinType, seed string, genCount int, hideSecretKey bool) (*ReadableWallet, error) {
	if genCount < 1 {
		return nil, errors.New("genCount must be > 0")
	}

	if seed == "" {
		return nil, errors.New("seed must not be the empty string")
	}

	wallet := &ReadableWallet{
		Meta: map[string]string{
			"coin": string(coinType),
			"seed": seed,
		},
	}

	seckeys := cipher.GenerateDeterministicKeyPairs([]byte(seed), genCount)

	for _, sec := range seckeys {
		pub := cipher.PubKeyFromSecKey(sec)

		var entry ReadableEntry
		switch coinType {
		case CoinTypeBitcoin:
			entry = GetBitcoinWalletEntry(pub, sec)
		case CoinTypeSkycoin:
			entry = GetSkycoinWalletEntry(pub, sec)
		default:
			return nil, fmt.Errorf(`unknown coinType "%s"`, coinType)
		}

		if hideSecretKey {
			entry.Secret = ""
		}

		wallet.Entries = append(wallet.Entries, entry)
	}

	return wallet, nil
}

// GetSkycoinWalletEntry returns a ReadableEntry in Skycoin format
func GetSkycoinWalletEntry(pub cipher.PubKey, sec cipher.SecKey) ReadableEntry {
	return ReadableEntry{
		Address: cipher.AddressFromPubKey(pub).String(),
		Public:  pub.Hex(),
		Secret:  sec.Hex(),
	}
}

// GetBitcoinWalletEntry returns a ReadableEntry in Bitcoin format
func GetBitcoinWalletEntry(pub cipher.PubKey, sec cipher.SecKey) ReadableEntry {
	return ReadableEntry{
		Address: cipher.BitcoinAddressFromPubkey(pub),
		Public:  pub.Hex(),
		Secret:  cipher.BitcoinWalletImportFormatFromSeckey(sec),
	}
}
