package secp256k1

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// TestKeyDerivationDebug outputs intermediate values for debugging firmware
func TestKeyDerivationDebug(t *testing.T) {
	// Use the same test mnemonic as firmware testing
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := []byte(mnemonic)

	fmt.Println("=== Key Derivation Debug ===")
	fmt.Printf("Mnemonic: %s\n", mnemonic)
	fmt.Printf("Mnemonic bytes: %s\n", hex.EncodeToString(seed))
	fmt.Printf("Mnemonic length: %d\n", len(seed))

	// Step 1: seed1 = Secp256k1Hash(seed)
	fmt.Println("\n--- Step 1: Secp256k1Hash(mnemonic) ---")

	// Inside Secp256k1Hash:
	hash := SumSHA256(seed)
	fmt.Printf("hash = SHA256(mnemonic): %s\n", hex.EncodeToString(hash))

	// _, seckey1 = deterministicKeyPairIteratorStep(hash)
	pubkey1, seckey1 := deterministicKeyPairIteratorStep(hash)
	fmt.Printf("seckey1 from step(hash): %s\n", hex.EncodeToString(seckey1))
	fmt.Printf("pubkey1 from step(hash): %s\n", hex.EncodeToString(pubkey1))

	// pubkeySeed = SumSHA256(hash)
	pubkeySeed := SumSHA256(hash)
	fmt.Printf("pubkeySeed = SHA256(hash): %s\n", hex.EncodeToString(pubkeySeed))

	// pubkey2, _ = deterministicKeyPairIteratorStep(pubkeySeed)
	pubkey2, _ := deterministicKeyPairIteratorStep(pubkeySeed)
	fmt.Printf("pubkey2 from step(pubkeySeed): %s\n", hex.EncodeToString(pubkey2))

	// ecdh = ECDH(pubkey2, seckey1)
	ecdhResult := ECDH(pubkey2, seckey1)
	fmt.Printf("ecdh = ECDH(pubkey2, seckey1): %s\n", hex.EncodeToString(ecdhResult))

	// seed1 = SumSHA256(hash + ecdh)
	combined := append(hash, ecdhResult...)
	fmt.Printf("combined (hash + ecdh), len=%d: %s\n", len(combined), hex.EncodeToString(combined))
	seed1 := SumSHA256(combined)
	fmt.Printf("seed1 = SHA256(combined): %s\n", hex.EncodeToString(seed1))

	// Verify this matches Secp256k1Hash
	seed1Check := Secp256k1Hash(seed)
	fmt.Printf("seed1 check (Secp256k1Hash): %s\n", hex.EncodeToString(seed1Check))

	// Step 2: seed2 = SHA256(seed + seed1)
	fmt.Println("\n--- Step 2: seed2 = SHA256(mnemonic + seed1) ---")
	combined2 := append(seed, seed1...)
	fmt.Printf("combined2 (mnemonic + seed1), len=%d\n", len(combined2))
	seed2 := SumSHA256(combined2)
	fmt.Printf("seed2 = SHA256(combined2): %s\n", hex.EncodeToString(seed2))

	// Step 3: pubkey, seckey = deterministicKeyPairIteratorStep(seed2)
	fmt.Println("\n--- Step 3: keypair from seed2 ---")
	pubkey, seckey := deterministicKeyPairIteratorStep(seed2)
	fmt.Printf("seckey: %s\n", hex.EncodeToString(seckey))
	fmt.Printf("pubkey: %s\n", hex.EncodeToString(pubkey))

	// Verify against DeterministicKeyPairIterator
	fmt.Println("\n--- Verification with DeterministicKeyPairIterator ---")
	nextSeed, pubkeyCheck, seckeyCheck := DeterministicKeyPairIterator(seed)
	fmt.Printf("nextSeed: %s\n", hex.EncodeToString(nextSeed))
	fmt.Printf("seckey check: %s\n", hex.EncodeToString(seckeyCheck))
	fmt.Printf("pubkey check: %s\n", hex.EncodeToString(pubkeyCheck))

	// Generate address intermediate values
	fmt.Println("\n--- Address Generation Intermediate Values ---")
	// Step 1: SHA256(pubkey)
	pubkeyHash1 := SumSHA256(pubkeyCheck)
	fmt.Printf("SHA256(pubkey): %s\n", hex.EncodeToString(pubkeyHash1))

	// Step 2: SHA256(SHA256(pubkey))
	pubkeyHash2 := SumSHA256(pubkeyHash1)
	fmt.Printf("SHA256(SHA256(pubkey)): %s\n", hex.EncodeToString(pubkeyHash2))

	// Step 3: RIPEMD160(SHA256(SHA256(pubkey))) - first 20 bytes
	// We don't have ripemd160 in this package, but the cipher package does

	fmt.Println("\n--- Full Test with Multiple Addresses ---")
	// Generate addresses 0-4
	currentSeed := seed
	for i := 0; i < 5; i++ {
		nextSeed, pk, sk := DeterministicKeyPairIterator(currentSeed)
		fmt.Printf("\nAddress %d:\n", i)
		fmt.Printf("  seckey: %s\n", hex.EncodeToString(sk))
		fmt.Printf("  pubkey: %s\n", hex.EncodeToString(pk))
		currentSeed = nextSeed
	}
}
