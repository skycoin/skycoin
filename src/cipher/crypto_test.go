package cipher

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/ripemd160"
)

func TestNewPubKey(t *testing.T) {
	_, err := NewPubKey(randBytes(t, 31))
	require.Equal(t, errors.New("Invalid public key length"), err)
	_, err = NewPubKey(randBytes(t, 32))
	require.Equal(t, errors.New("Invalid public key length"), err)
	_, err = NewPubKey(randBytes(t, 34))
	require.Equal(t, errors.New("Invalid public key length"), err)
	_, err = NewPubKey(randBytes(t, 0))
	require.Equal(t, errors.New("Invalid public key length"), err)
	_, err = NewPubKey(randBytes(t, 100))
	require.Equal(t, errors.New("Invalid public key length"), err)

	_, err = NewPubKey(make([]byte, len(PubKey{})))
	require.Equal(t, errors.New("Invalid public key"), err)

	p, _ := GenerateKeyPair()
	p2, err := NewPubKey(p[:])
	require.NoError(t, err)
	require.Equal(t, p, p2)
}

func TestMustNewPubKey(t *testing.T) {
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 31)) })
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 32)) })
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 34)) })
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 0)) })
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 100)) })
	require.Panics(t, func() { MustNewPubKey(make([]byte, len(PubKey{}))) })

	p, _ := GenerateKeyPair()
	p2 := MustNewPubKey(p[:])
	require.Equal(t, p, p2)
}

func TestPubKeyFromHex(t *testing.T) {
	// Invalid hex
	_, err := PubKeyFromHex("")
	require.Equal(t, errors.New("Invalid public key length"), err)

	_, err = PubKeyFromHex("cascs")
	require.Equal(t, errors.New("Invalid public key"), err)

	// Empty key
	empty := PubKey{}
	h := hex.EncodeToString(empty[:])
	_, err = PubKeyFromHex(h)
	require.Equal(t, errors.New("Invalid public key"), err)

	// Invalid hex length
	p, _ := GenerateKeyPair()
	s := hex.EncodeToString(p[:len(p)/2])
	_, err = PubKeyFromHex(s)
	require.Equal(t, errors.New("Invalid public key length"), err)

	// Valid
	s = hex.EncodeToString(p[:])
	p2, err := PubKeyFromHex(s)
	require.NoError(t, err)
	require.Equal(t, p, p2)
}

func TestMustPubKeyFromHex(t *testing.T) {
	// Invalid hex
	require.Panics(t, func() { MustPubKeyFromHex("") })
	require.Panics(t, func() { MustPubKeyFromHex("cascs") })

	// Empty key
	empty := PubKey{}
	h := hex.EncodeToString(empty[:])
	require.Panics(t, func() { MustPubKeyFromHex(h) })

	// Invalid hex length
	p, _ := GenerateKeyPair()
	s := hex.EncodeToString(p[:len(p)/2])
	require.Panics(t, func() { MustPubKeyFromHex(s) })

	// Valid
	s = hex.EncodeToString(p[:])
	require.NotPanics(t, func() { MustPubKeyFromHex(s) })
	require.Equal(t, p, MustPubKeyFromHex(s))
}

func TestPubKeyHex(t *testing.T) {
	p, _ := GenerateKeyPair()
	h := p.Hex()
	p2, err := PubKeyFromHex(h)
	require.NoError(t, err)
	require.Equal(t, p2, p)
	require.Equal(t, p2.Hex(), h)
}

func TestNewPubKeyRandom(t *testing.T) {
	// Random bytes should not be valid, most of the time
	failed := false
	for i := 0; i < 10; i++ {
		b := randBytes(t, 33)
		if _, err := NewPubKey(b); err != nil {
			failed = true
			break
		}
	}
	require.True(t, failed)
}

func TestPubKeyVerify(t *testing.T) {
	// Random bytes should not be valid, most of the time
	failed := false
	for i := 0; i < 10; i++ {
		b := randBytes(t, 33)
		p := PubKey{}
		copy(p[:], b[:])
		if p.Verify() != nil {
			failed = true
			break
		}
	}
	require.True(t, failed)
}

func TestPubKeyNullVerifyFails(t *testing.T) {
	// Empty public key should not be valid
	p := PubKey{}
	require.Error(t, p.Verify())
}

func TestPubKeyVerifyDefault1(t *testing.T) {
	// Generated pub key should be valid
	p, _ := GenerateKeyPair()
	require.NoError(t, p.Verify())
}

func TestPubKeyVerifyDefault2(t *testing.T) {
	for i := 0; i < 1024; i++ {
		p, _ := GenerateKeyPair()
		require.NoError(t, p.Verify())
	}
}

func TestPubKeyRipemd160(t *testing.T) {
	p, _ := GenerateKeyPair()
	h := PubKeyRipemd160(p)
	// Should be Ripemd160(SHA256(SHA256()))
	x := sha256.Sum256(p[:])
	x = sha256.Sum256(x[:])
	rh := ripemd160.New()
	_, err := rh.Write(x[:])
	require.NoError(t, err)
	y := rh.Sum(nil)
	require.True(t, bytes.Equal(h[:], y))
}

func TestPubKeyToAddress(t *testing.T) {
	p, _ := GenerateKeyPair()
	addr := AddressFromPubKey(p)
	//func (self Address) Verify(key PubKey) error {
	err := addr.Verify(p)
	require.NoError(t, err)
	addrStr := addr.String()
	_, err = DecodeBase58Address(addrStr)
	//func DecodeBase58Address(addr string) (Address, error) {
	require.NoError(t, err)
}

func TestPubKeyToAddress2(t *testing.T) {
	for i := 0; i < 1024; i++ {
		p, _ := GenerateKeyPair()
		addr := AddressFromPubKey(p)
		//func (self Address) Verify(key PubKey) error {
		err := addr.Verify(p)
		require.NoError(t, err)
		addrStr := addr.String()
		_, err = DecodeBase58Address(addrStr)
		//func DecodeBase58Address(addr string) (Address, error) {
		require.NoError(t, err)
	}
}

func TestNewSecKey(t *testing.T) {
	_, err := NewSecKey(randBytes(t, 31))
	require.Equal(t, errors.New("Invalid secret key length"), err)
	_, err = NewSecKey(randBytes(t, 33))
	require.Equal(t, errors.New("Invalid secret key length"), err)
	_, err = NewSecKey(randBytes(t, 34))
	require.Equal(t, errors.New("Invalid secret key length"), err)
	_, err = NewSecKey(randBytes(t, 0))
	require.Equal(t, errors.New("Invalid secret key length"), err)
	_, err = NewSecKey(randBytes(t, 100))
	require.Equal(t, errors.New("Invalid secret key length"), err)

	b := randBytes(t, 32)
	p, err := NewSecKey(b)
	require.NoError(t, err)
	require.True(t, bytes.Equal(p[:], b))
}

func TestMustNewSecKey(t *testing.T) {
	require.Panics(t, func() { MustNewSecKey(randBytes(t, 31)) })
	require.Panics(t, func() { MustNewSecKey(randBytes(t, 33)) })
	require.Panics(t, func() { MustNewSecKey(randBytes(t, 34)) })
	require.Panics(t, func() { MustNewSecKey(randBytes(t, 0)) })
	require.Panics(t, func() { MustNewSecKey(randBytes(t, 100)) })
	require.NotPanics(t, func() { MustNewSecKey(randBytes(t, 32)) })
	b := randBytes(t, 32)
	p := MustNewSecKey(b)
	require.True(t, bytes.Equal(p[:], b))
}

func TestSecKeyFromHex(t *testing.T) {
	// Invalid hex
	_, err := SecKeyFromHex("")
	require.Equal(t, errors.New("Invalid secret key length"), err)

	_, err = SecKeyFromHex("cascs")
	require.Equal(t, errors.New("Invalid secret key"), err)

	// Invalid hex length
	p := MustNewSecKey(randBytes(t, 32))
	s := hex.EncodeToString(p[:len(p)/2])
	_, err = SecKeyFromHex(s)
	require.Equal(t, errors.New("Invalid secret key length"), err)

	// Valid
	s = hex.EncodeToString(p[:])
	p2, err := SecKeyFromHex(s)
	require.NoError(t, err)
	require.Equal(t, p2, p)
}

func TestMustSecKeyFromHex(t *testing.T) {
	// Invalid hex
	require.Panics(t, func() { MustSecKeyFromHex("") })
	require.Panics(t, func() { MustSecKeyFromHex("cascs") })
	// Invalid hex length
	p := MustNewSecKey(randBytes(t, 32))
	s := hex.EncodeToString(p[:len(p)/2])
	require.Panics(t, func() { MustSecKeyFromHex(s) })
	// Valid
	s = hex.EncodeToString(p[:])
	require.NotPanics(t, func() { MustSecKeyFromHex(s) })
	require.Equal(t, p, MustSecKeyFromHex(s))
}

func TestSecKeyHex(t *testing.T) {
	b := randBytes(t, 32)
	p := MustNewSecKey(b)
	h := p.Hex()
	p2 := MustSecKeyFromHex(h)
	require.Equal(t, p2, p)
	require.Equal(t, p2.Hex(), h)
}

func TestSecKeyVerify(t *testing.T) {
	// Empty secret key should not be valid
	p := SecKey{}
	require.Error(t, p.Verify())

	// Generated sec key should be valid
	_, p = GenerateKeyPair()
	require.NoError(t, p.Verify())

	// Random bytes are usually valid
}

func TestECDH(t *testing.T) {
	pub1, sec1 := GenerateKeyPair()
	pub2, sec2 := GenerateKeyPair()

	buf1, err := ECDH(pub2, sec1)
	require.NoError(t, err)
	buf2, err := ECDH(pub1, sec2)
	require.NoError(t, err)

	require.True(t, bytes.Equal(buf1, buf2))

	goodPub, goodSec := GenerateKeyPair()
	var badPub PubKey
	var badSec SecKey

	_, err = ECDH(badPub, goodSec)
	require.Equal(t, errors.New("ECDH invalid pubkey input"), err)
	_, err = ECDH(goodPub, badSec)
	require.Equal(t, errors.New("ECDH invalid seckey input"), err)

	for i := 0; i < 128; i++ {
		pub1, sec1 := GenerateKeyPair()
		pub2, sec2 := GenerateKeyPair()
		buf1, err := ECDH(pub2, sec1)
		require.NoError(t, err)
		buf2, err := ECDH(pub1, sec2)
		require.NoError(t, err)
		require.True(t, bytes.Equal(buf1, buf2))
	}
}

func TestMustECDH(t *testing.T) {
	goodPub, goodSec := GenerateKeyPair()
	var badPub PubKey
	var badSec SecKey

	require.Panics(t, func() {
		MustECDH(badPub, goodSec)
	})
	require.Panics(t, func() {
		MustECDH(goodPub, badSec)
	})

	pub1, sec1 := GenerateKeyPair()
	pub2, sec2 := GenerateKeyPair()

	buf1 := MustECDH(pub2, sec1)
	buf2 := MustECDH(pub1, sec2)

	require.True(t, bytes.Equal(buf1, buf2))
}

func TestNewSig(t *testing.T) {
	_, err := NewSig(randBytes(t, 64))
	require.Equal(t, errors.New("Invalid signature length"), err)
	_, err = NewSig(randBytes(t, 66))
	require.Equal(t, errors.New("Invalid signature length"), err)
	_, err = NewSig(randBytes(t, 67))
	require.Equal(t, errors.New("Invalid signature length"), err)
	_, err = NewSig(randBytes(t, 0))
	require.Equal(t, errors.New("Invalid signature length"), err)
	_, err = NewSig(randBytes(t, 100))
	require.Equal(t, errors.New("Invalid signature length"), err)

	b := randBytes(t, 65)
	p, err := NewSig(b)
	require.NoError(t, err)
	require.True(t, bytes.Equal(p[:], b))
}

func TestMustNewSig(t *testing.T) {
	require.Panics(t, func() { MustNewSig(randBytes(t, 64)) })
	require.Panics(t, func() { MustNewSig(randBytes(t, 66)) })
	require.Panics(t, func() { MustNewSig(randBytes(t, 67)) })
	require.Panics(t, func() { MustNewSig(randBytes(t, 0)) })
	require.Panics(t, func() { MustNewSig(randBytes(t, 100)) })

	require.NotPanics(t, func() { MustNewSig(randBytes(t, 65)) })

	b := randBytes(t, 65)
	p := MustNewSig(b)
	require.True(t, bytes.Equal(p[:], b))
}

func TestSigFromHex(t *testing.T) {
	// Invalid hex
	_, err := SigFromHex("")
	require.Equal(t, errors.New("Invalid signature length"), err)

	_, err = SigFromHex("cascs")
	require.Equal(t, errors.New("Invalid signature"), err)

	// Invalid hex length
	p := MustNewSig(randBytes(t, 65))
	s := hex.EncodeToString(p[:len(p)/2])
	_, err = SigFromHex(s)
	require.Equal(t, errors.New("Invalid signature length"), err)

	// Valid
	s = hex.EncodeToString(p[:])
	s2, err := SigFromHex(s)
	require.NoError(t, err)
	require.Equal(t, p, s2)
}

func TestMustSigFromHex(t *testing.T) {
	// Invalid hex
	require.Panics(t, func() { MustSigFromHex("") })
	require.Panics(t, func() { MustSigFromHex("cascs") })
	// Invalid hex length
	p := MustNewSig(randBytes(t, 65))
	s := hex.EncodeToString(p[:len(p)/2])
	require.Panics(t, func() { MustSigFromHex(s) })
	// Valid
	s = hex.EncodeToString(p[:])
	require.NotPanics(t, func() { MustSigFromHex(s) })
	require.Equal(t, p, MustSigFromHex(s))
}

func TestSigHex(t *testing.T) {
	b := randBytes(t, 65)
	p := MustNewSig(b)
	h := p.Hex()
	p2 := MustSigFromHex(h)
	require.Equal(t, p2, p)
	require.Equal(t, p2.Hex(), h)
}

func TestVerifyAddressSignedHash(t *testing.T) {
	p, s := GenerateKeyPair()
	require.NoError(t, p.Verify())
	require.NoError(t, s.Verify())
	a := AddressFromPubKey(p)
	require.NoError(t, a.Verify(p))
	b := randBytes(t, 256)
	h := SumSHA256(b)
	sig := MustSignHash(h, s)
	require.NoError(t, VerifyAddressSignedHash(a, sig, h))
	// Empty sig should be invalid
	require.Error(t, VerifyAddressSignedHash(a, Sig{}, h))
	// Random sigs should not pass
	for i := 0; i < 100; i++ {
		require.Error(t, VerifyAddressSignedHash(a, MustNewSig(randBytes(t, 65)), h))
	}
	// Sig for one hash does not work for another hash
	h2 := SumSHA256(randBytes(t, 256))
	sig2 := MustSignHash(h2, s)
	require.NoError(t, VerifyAddressSignedHash(a, sig2, h2))
	require.Error(t, VerifyAddressSignedHash(a, sig2, h))
	require.Error(t, VerifyAddressSignedHash(a, sig, h2))

	// Different secret keys should not create same sig
	p2, s2 := GenerateKeyPair()
	a2 := AddressFromPubKey(p2)
	h = SumSHA256(randBytes(t, 256))
	sig = MustSignHash(h, s)
	sig2 = MustSignHash(h, s2)
	require.NoError(t, VerifyAddressSignedHash(a, sig, h))
	require.NoError(t, VerifyAddressSignedHash(a2, sig2, h))
	require.NotEqual(t, sig, sig2)
	h = SumSHA256(randBytes(t, 256))
	sig = MustSignHash(h, s)
	sig2 = MustSignHash(h, s2)
	require.NoError(t, VerifyAddressSignedHash(a, sig, h))
	require.NoError(t, VerifyAddressSignedHash(a2, sig2, h))
	require.NotEqual(t, sig, sig2)

	// Bad address should be invalid
	require.Error(t, VerifyAddressSignedHash(a, sig2, h))
	require.Error(t, VerifyAddressSignedHash(a2, sig, h))

	// Empty hash should panic
	require.Panics(t, func() {
		MustSignHash(SHA256{}, s)
	})
}

func TestSignHash(t *testing.T) {
	p, s := GenerateKeyPair()
	a := AddressFromPubKey(p)
	h := SumSHA256(randBytes(t, 256))
	sig, err := SignHash(h, s)
	require.NoError(t, err)
	require.NotEqual(t, sig, Sig{})
	require.NoError(t, VerifyAddressSignedHash(a, sig, h))
	require.NoError(t, VerifyPubKeySignedHash(p, sig, h))

	p2, err := PubKeyFromSig(sig, h)
	require.NoError(t, err)
	require.Equal(t, p, p2)

	_, err = SignHash(h, SecKey{})
	require.Equal(t, ErrInvalidSecKey, err)

	_, err = SignHash(SHA256{}, s)
	require.Equal(t, ErrNullSignHash, err)
}

func TestMustSignHash(t *testing.T) {
	p, s := GenerateKeyPair()
	a := AddressFromPubKey(p)
	h := SumSHA256(randBytes(t, 256))
	sig := MustSignHash(h, s)
	require.NotEqual(t, sig, Sig{})
	require.NoError(t, VerifyAddressSignedHash(a, sig, h))

	require.Panics(t, func() {
		MustSignHash(h, SecKey{})
	})
}

func TestPubKeyFromSecKey(t *testing.T) {
	p, s := GenerateKeyPair()
	p2, err := PubKeyFromSecKey(s)
	require.NoError(t, err)
	require.Equal(t, p2, p)

	_, err = PubKeyFromSecKey(SecKey{})
	require.Equal(t, errors.New("Attempt to load null seckey, unsafe"), err)
}

func TestMustPubKeyFromSecKey(t *testing.T) {
	p, s := GenerateKeyPair()
	require.Equal(t, MustPubKeyFromSecKey(s), p)
	require.Panics(t, func() { MustPubKeyFromSecKey(SecKey{}) })
}

func TestPubKeyFromSig(t *testing.T) {
	p, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	sig := MustSignHash(h, s)
	p2, err := PubKeyFromSig(sig, h)
	require.Equal(t, p, p2)
	require.NoError(t, err)
	_, err = PubKeyFromSig(Sig{}, h)
	require.Error(t, err)
}

func TestMustPubKeyFromSig(t *testing.T) {
	p, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	sig := MustSignHash(h, s)
	p2 := MustPubKeyFromSig(sig, h)
	require.Equal(t, p, p2)

	require.Panics(t, func() {
		_ = MustPubKeyFromSig(Sig{}, h)
	})
}

func TestVerifyPubKeySignedHash(t *testing.T) {
	p, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	h2 := SumSHA256(randBytes(t, 256))
	sig := MustSignHash(h, s)
	require.NoError(t, VerifyPubKeySignedHash(p, sig, h))
	require.Error(t, VerifyPubKeySignedHash(p, Sig{}, h))
	require.Error(t, VerifyPubKeySignedHash(p, sig, h2))
	p2, _ := GenerateKeyPair()
	require.Error(t, VerifyPubKeySignedHash(p2, sig, h))
	require.Error(t, VerifyPubKeySignedHash(PubKey{}, sig, h))
}

func TestGenerateKeyPair(t *testing.T) {
	for i := 0; i < 10; i++ {
		p, s := GenerateKeyPair()
		require.NoError(t, p.Verify())
		require.NoError(t, s.Verify())
		err := CheckSecKey(s)
		require.NoError(t, err)
	}
}

func TestGenerateDeterministicKeyPair(t *testing.T) {
	// TODO -- deterministic key pairs are useless as is because we can't
	// generate pair n+1, only pair 0
	seed := randBytes(t, 32)
	p, s := MustGenerateDeterministicKeyPair(seed)
	require.NoError(t, p.Verify())
	require.NoError(t, s.Verify())
	p, s = MustGenerateDeterministicKeyPair(seed)
	require.NoError(t, p.Verify())
	require.NoError(t, s.Verify())

	_, _, err := GenerateDeterministicKeyPair(nil)
	require.Equal(t, errors.New("Seed input is empty"), err)

	require.Panics(t, func() {
		MustGenerateDeterministicKeyPair(nil)
	})
}

func TestGenerateDeterministicKeyPairs(t *testing.T) {
	seed := randBytes(t, 32)
	keys, err := GenerateDeterministicKeyPairs(seed, 4)
	require.NoError(t, err)
	require.Len(t, keys, 4)
	for _, s := range keys {
		require.NoError(t, s.Verify())
	}

	keys2 := MustGenerateDeterministicKeyPairs(seed, 4)
	require.Equal(t, keys, keys2)

	_, err = GenerateDeterministicKeyPairs(nil, 1)
	require.Equal(t, errors.New("Seed input is empty"), err)

	require.Panics(t, func() {
		MustGenerateDeterministicKeyPairs(nil, 1)
	})
}

func TestGenerateDeterministicKeyPairsSeed(t *testing.T) {
	seed := randBytes(t, 32)
	newSeed, keys, err := GenerateDeterministicKeyPairsSeed(seed, 4)
	require.NoError(t, err)
	require.Len(t, newSeed, 32)
	require.NotEqual(t, seed, newSeed)
	require.Len(t, keys, 4)
	for _, s := range keys {
		require.NoError(t, s.Verify())
	}

	newSeed2, keys2 := MustGenerateDeterministicKeyPairsSeed(seed, 4)
	require.Equal(t, newSeed, newSeed2)
	require.Equal(t, keys, keys2)

	_, _, err = GenerateDeterministicKeyPairsSeed(nil, 4)
	require.Equal(t, errors.New("Seed input is empty"), err)

	require.Panics(t, func() {
		MustGenerateDeterministicKeyPairsSeed(nil, 4)
	})
}

func TestDeterministicKeyPairIterator(t *testing.T) {
	seed := randBytes(t, 32)
	newSeed, p, s, err := DeterministicKeyPairIterator(seed)
	require.NoError(t, err)
	require.NoError(t, p.Verify())
	require.NoError(t, s.Verify())
	require.NotEqual(t, seed, newSeed)
	require.Len(t, newSeed, 32)

	newSeed2, p2, s2 := MustDeterministicKeyPairIterator(seed)
	require.Equal(t, newSeed, newSeed2)
	require.Equal(t, p, p2)
	require.Equal(t, s, s2)

	_, _, _, err = DeterministicKeyPairIterator(nil)
	require.Equal(t, errors.New("Seed input is empty"), err)

	require.Panics(t, func() {
		MustDeterministicKeyPairIterator(nil)
	})
}

func TestCheckSecKey(t *testing.T) {
	_, s := GenerateKeyPair()
	require.NoError(t, CheckSecKey(s))
	require.Error(t, CheckSecKey(SecKey{}))
}

func TestCheckSecKeyHash(t *testing.T) {
	_, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	require.NoError(t, CheckSecKeyHash(s, h))
	require.Error(t, CheckSecKeyHash(SecKey{}, h))
}

func TestGenerateDeterministicKeyPairsUsesAllBytes(t *testing.T) {
	// Tests that if a seed >128 bits is used, the generator does not ignore bits >128
	seed := "property diet little foster provide disagree witness mountain alley weekend kitten general"
	seckeys := MustGenerateDeterministicKeyPairs([]byte(seed), 3)
	seckeys2 := MustGenerateDeterministicKeyPairs([]byte(seed[:16]), 3)
	require.NotEqual(t, seckeys, seckeys2)
}

func TestPubkey1(t *testing.T) {
	// This was migrated from coin/coin_test.go
	a := "02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8"
	b, err := hex.DecodeString(a)
	require.NoError(t, err)

	p, err := NewPubKey(b)
	require.NoError(t, err)
	require.NoError(t, p.Verify())

	addr := AddressFromPubKey(p)
	require.NoError(t, addr.Verify(p))
}

func TestSecKey1(t *testing.T) {
	// This was migrated from coin/coin_test.go
	a := "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d"
	b, err := hex.DecodeString(a)
	require.NoError(t, err)
	require.Len(t, b, 32)

	seckey, err := NewSecKey(b)
	require.NoError(t, err)
	require.NoError(t, seckey.Verify())

	pubkey, err := PubKeyFromSecKey(seckey)
	require.NoError(t, err)
	require.NoError(t, pubkey.Verify())

	addr := AddressFromPubKey(pubkey)
	require.NoError(t, addr.Verify(pubkey))

	test := []byte("test message")
	hash := SumSHA256(test)
	err = CheckSecKeyHash(seckey, hash)
	require.NoError(t, err)
}

func TestSecKeyPubKeyNull(t *testing.T) {
	var pk PubKey
	require.True(t, pk.Null())
	pk[0] = 1
	require.False(t, pk.Null())

	var sk SecKey
	require.True(t, sk.Null())
	sk[0] = 1
	require.False(t, sk.Null())

	sk, err := NewSecKey(randBytes(t, 32))
	require.NoError(t, err)
	pk = MustPubKeyFromSecKey(sk)

	require.False(t, sk.Null())
	require.False(t, pk.Null())
}

func TestVerifySignedHash(t *testing.T) {
	h := MustSHA256FromHex("127e9b0d6b71cecd0363b366413f0f19fcd924ae033513498e7486570ff2a1c8")
	sig := MustSigFromHex("63c035b0c95d0c5744fc1c0bdf38af02cef2d2f65a8f923732ab44e436f8a491216d9ab5ff795e3144f4daee37077b8b9db54d2ba3a3df8d4992f06bb21f724401")

	err := VerifySignatureRecoverPubKey(sig, h)
	require.NoError(t, err)

	// Fails with ErrInvalidHashForSig
	badSigHex := "71f2c01516fe696328e79bcf464eb0db374b63d494f7a307d1e77114f18581d7a81eed5275a9e04a336292dd2fd16977d9bef2a54ea3161d0876603d00c53bc9dd"
	badSig := MustSigFromHex(badSigHex)
	err = VerifySignatureRecoverPubKey(badSig, h)
	require.Equal(t, ErrInvalidHashForSig, err)

	// Fails with ErrInvalidSigPubKeyRecovery
	badSig = MustSigFromHex("63c035b0c95d0c5744fc1c0bdf39af02cef2d2f65a8f923732ab44e436f8a491216d9ab5ff795e3144f4daee37077b8b9db54d2ba3a3df8d4992f06bb21f724401")
	err = VerifySignatureRecoverPubKey(badSig, h)
	require.Equal(t, ErrInvalidSigPubKeyRecovery, err)
}

func TestHighSPointSigInvalid(t *testing.T) {
	// Verify that signatures that were generated with forceLowS=false
	// are not accepted as valid, to avoid a signature malleability case.
	// Refer to secp256k1go's TestSigForceLowS for the reference test inputs

	h := MustSHA256FromHex("DD72CBF2203C1A55A411EEC4404AF2AFB2FE942C434B23EFE46E9F04DA8433CA")

	// This signature has a high S point (the S point is above the half-order of the curve)
	sigHexHighS := "8c20a668be1b5a910205de46095023fe4823a3757f4417114168925f28193bffadf317cc256cec28d90d5b2b7e1ce6a45cd5f3b10880ab5f99c389c66177d39a01"
	s := MustSigFromHex(sigHexHighS)
	err := VerifySignatureRecoverPubKey(s, h)

	require.Error(t, err)
	require.Equal(t, "Signature not valid for hash", err.Error())

	// This signature has a low S point (the S point is below the half-order of the curve).
	// It is equal to forceLowS(sigHighS).
	sigHexLowS := "8c20a668be1b5a910205de46095023fe4823a3757f4417114168925f28193bff520ce833da9313d726f2a4d481e3195a5dd8e935a6c7f4dc260ed4c66ebe6da700"
	s2 := MustSigFromHex(sigHexLowS)
	err = VerifySignatureRecoverPubKey(s2, h)
	require.NoError(t, err)
}
