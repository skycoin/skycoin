package cipher

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// TestAddressDerivationDebug outputs full address derivation values for firmware debugging
func TestAddressDerivationDebug(t *testing.T) {
	// Use the same test mnemonic as firmware testing
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := []byte(mnemonic)

	fmt.Println("=== Address Derivation Debug ===")
	fmt.Printf("Mnemonic: %s\n", mnemonic)
	fmt.Printf("Mnemonic length: %d\n", len(seed))

	// Generate first 5 addresses
	currentSeed := seed
	for i := 0; i < 5; i++ {
		nextSeed, pk, sk := MustDeterministicKeyPairIterator(currentSeed)

		// Generate address from pubkey
		addr := AddressFromPubKey(pk)

		fmt.Printf("\n--- Address %d ---\n", i)
		fmt.Printf("seckey: %s\n", hex.EncodeToString(sk[:]))
		fmt.Printf("pubkey: %s\n", hex.EncodeToString(pk[:]))
		fmt.Printf("address: %s\n", addr.String())

		// Show intermediate hashes for address generation
		// SHA256(pubkey)
		hash1 := SumSHA256(pk[:])
		fmt.Printf("SHA256(pubkey): %s\n", hex.EncodeToString(hash1[:]))

		// SHA256(SHA256(pubkey))
		hash2 := SumSHA256(hash1[:])
		fmt.Printf("SHA256(SHA256(pubkey)): %s\n", hex.EncodeToString(hash2[:]))

		// RIPEMD160(SHA256(SHA256(pubkey)))
		ripemd := HashRipemd160(hash2[:])
		fmt.Printf("RIPEMD160(hash2): %s\n", hex.EncodeToString(ripemd[:]))

		// Show the 25-byte input to base58 (for Skycoin format)
		// Skycoin: ripemd (20) + version (1) + checksum (4)
		var addrBytes [25]byte
		copy(addrBytes[:20], ripemd[:])
		addrBytes[20] = 0x00 // version

		// Checksum: first 4 bytes of SHA256(ripemd + version)
		checksumHash := SumSHA256(addrBytes[:21])
		copy(addrBytes[21:], checksumHash[:4])

		fmt.Printf("Address bytes (before base58): %s\n", hex.EncodeToString(addrBytes[:]))

		currentSeed = nextSeed
	}
}
