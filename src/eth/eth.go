// Package eth implements what a wallet needs to hold Ether and ERC-20 tokens:
// BIP44 keys, EIP-55 addresses, and signed EIP-1559 transactions.
package eth

import (
	"golang.org/x/crypto/sha3"
)

// Keccak256 is the hash Ethereum calls SHA3, the pre-standard Keccak.
func Keccak256(data ...[]byte) [32]byte {
	h := sha3.NewLegacyKeccak256()
	for _, d := range data {
		h.Write(d) //nolint:errcheck // a hash never fails to write
	}
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}
