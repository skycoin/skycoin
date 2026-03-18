package cipher

import (
	"fmt"
	"log"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/bech32"
)

// BitcoinSegwitAddress is a native segwit (P2WPKH) bitcoin address (bc1q...).
type BitcoinSegwitAddress struct {
	Version byte      // Witness version (0 for P2WPKH)
	Key     Ripemd160 // 20-byte witness program (hash160 of compressed pubkey)
}

// BitcoinSegwitAddressFromPubKey creates a P2WPKH address from a compressed public key.
// The witness program is RIPEMD160(SHA256(pubkey)), same as P2PKH but encoded as bech32.
func BitcoinSegwitAddressFromPubKey(pubKey PubKey) BitcoinSegwitAddress {
	return BitcoinSegwitAddress{
		Version: 0,
		Key:     BitcoinPubKeyRipemd160(pubKey),
	}
}

// BitcoinSegwitAddressFromSecKey generates a BitcoinSegwitAddress from SecKey.
func BitcoinSegwitAddressFromSecKey(secKey SecKey) (BitcoinSegwitAddress, error) {
	p, err := PubKeyFromSecKey(secKey)
	if err != nil {
		return BitcoinSegwitAddress{}, err
	}
	return BitcoinSegwitAddressFromPubKey(p), nil
}

// MustBitcoinSegwitAddressFromSecKey generates a BitcoinSegwitAddress from SecKey, panics on error.
func MustBitcoinSegwitAddressFromSecKey(secKey SecKey) BitcoinSegwitAddress {
	return BitcoinSegwitAddressFromPubKey(MustPubKeyFromSecKey(secKey))
}

// DecodeBech32BitcoinAddress decodes a bech32 segwit address string into a BitcoinSegwitAddress.
func DecodeBech32BitcoinAddress(addr string) (BitcoinSegwitAddress, error) {
	hrp := "bc"
	if strings.HasPrefix(strings.ToLower(addr), "tb1") {
		hrp = "tb"
	}

	version, program, err := bech32.SegwitDecode(hrp, addr)
	if err != nil {
		return BitcoinSegwitAddress{}, fmt.Errorf("invalid bech32 address: %w", err)
	}

	if version != 0 {
		return BitcoinSegwitAddress{}, fmt.Errorf("unsupported witness version: %d", version)
	}
	if len(program) != 20 {
		return BitcoinSegwitAddress{}, fmt.Errorf("invalid P2WPKH program length: %d (expected 20)", len(program))
	}

	a := BitcoinSegwitAddress{Version: version}
	copy(a.Key[:], program)
	return a, nil
}

// MustDecodeBech32BitcoinAddress decodes a bech32 segwit address, panics on error.
func MustDecodeBech32BitcoinAddress(addr string) BitcoinSegwitAddress {
	a, err := DecodeBech32BitcoinAddress(addr)
	if err != nil {
		log.Panicf("Invalid bech32 bitcoin address %s: %v", addr, err)
	}
	return a
}

// Null returns true if the address is null.
func (addr BitcoinSegwitAddress) Null() bool {
	return addr == BitcoinSegwitAddress{}
}

// Bytes returns the raw witness program with version prefix.
func (addr BitcoinSegwitAddress) Bytes() []byte {
	b := make([]byte, 1+20)
	b[0] = addr.Version
	copy(b[1:], addr.Key[:])
	return b
}

// Verify checks that the address is valid for the given public key.
func (addr BitcoinSegwitAddress) Verify(key PubKey) error {
	if addr.Version != 0 {
		return ErrAddressInvalidVersion
	}
	if addr.Key != BitcoinPubKeyRipemd160(key) {
		return ErrAddressInvalidPubKey
	}
	return nil
}

// String returns the bech32-encoded address (bc1q...).
func (addr BitcoinSegwitAddress) String() string {
	s, err := bech32.SegwitEncode("bc", addr.Version, addr.Key[:])
	if err != nil {
		log.Panicf("BitcoinSegwitAddress.String failed: %v", err)
	}
	return s
}

// Checksum returns the first 4 bytes of SHA256(version+key).
// This satisfies the Addresser interface. For bech32 addresses, the real
// checksum is embedded in the bech32 encoding, but this provides compatibility.
func (addr BitcoinSegwitAddress) Checksum() Checksum {
	r := append([]byte{addr.Version}, addr.Key[:]...)
	h := SumSHA256(r)
	c := Checksum{}
	copy(c[:], h[:4])
	return c
}
