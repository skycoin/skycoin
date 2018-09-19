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

	b := randBytes(t, 33)
	p, err := NewPubKey(b)
	require.NoError(t, err)
	require.True(t, bytes.Equal(p[:], b))
}

func TestMustNewPubKey(t *testing.T) {
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 31)) })
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 32)) })
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 34)) })
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 0)) })
	require.Panics(t, func() { MustNewPubKey(randBytes(t, 100)) })
	require.NotPanics(t, func() { MustNewPubKey(randBytes(t, 33)) })
	b := randBytes(t, 33)
	p := MustNewPubKey(b)
	require.True(t, bytes.Equal(p[:], b))
}

func TestPubKeyFromHex(t *testing.T) {
	// Invalid hex
	p, err := PubKeyFromHex("")
	require.Equal(t, errors.New("Invalid public key length"), err)

	p, err = PubKeyFromHex("cascs")
	require.Equal(t, errors.New("Invalid public key"), err)

	// Invalid hex length
	p = MustNewPubKey(randBytes(t, 33))
	s := hex.EncodeToString(p[:len(p)/2])
	p, err = PubKeyFromHex(s)
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
	// Invalid hex length
	p := MustNewPubKey(randBytes(t, 33))
	s := hex.EncodeToString(p[:len(p)/2])
	require.Panics(t, func() { MustPubKeyFromHex(s) })
	// Valid
	s = hex.EncodeToString(p[:])
	require.NotPanics(t, func() { MustPubKeyFromHex(s) })
	require.Equal(t, p, MustPubKeyFromHex(s))
}

func TestPubKeyHex(t *testing.T) {
	b := randBytes(t, 33)
	p := MustNewPubKey(b)
	h := p.Hex()
	p2, err := PubKeyFromHex(h)
	require.NoError(t, err)
	require.Equal(t, p2, p)
	require.Equal(t, p2.Hex(), h)
}

func TestPubKeyVerify(t *testing.T) {
	// Random bytes should not be valid, most of the time
	failed := false
	for i := 0; i < 10; i++ {
		b := randBytes(t, 33)
		if MustNewPubKey(b).Verify() != nil {
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

func TestChkSig(t *testing.T) {
	p, s := GenerateKeyPair()
	require.NoError(t, p.Verify())
	require.NoError(t, s.Verify())
	a := AddressFromPubKey(p)
	require.NoError(t, a.Verify(p))
	b := randBytes(t, 256)
	h := SumSHA256(b)
	sig := MustSignHash(h, s)
	require.NoError(t, ChkSig(a, h, sig))
	// Empty sig should be invalid
	require.Error(t, ChkSig(a, h, Sig{}))
	// Random sigs should not pass
	for i := 0; i < 100; i++ {
		require.Error(t, ChkSig(a, h, MustNewSig(randBytes(t, 65))))
	}
	// Sig for one hash does not work for another hash
	h2 := SumSHA256(randBytes(t, 256))
	sig2 := MustSignHash(h2, s)
	require.NoError(t, ChkSig(a, h2, sig2))
	require.Error(t, ChkSig(a, h, sig2))
	require.Error(t, ChkSig(a, h2, sig))

	// Different secret keys should not create same sig
	p2, s2 := GenerateKeyPair()
	a2 := AddressFromPubKey(p2)
	h = SHA256{}
	sig = MustSignHash(h, s)
	sig2 = MustSignHash(h, s2)
	require.NoError(t, ChkSig(a, h, sig))
	require.NoError(t, ChkSig(a2, h, sig2))
	require.NotEqual(t, sig, sig2)
	h = SumSHA256(randBytes(t, 256))
	sig = MustSignHash(h, s)
	sig2 = MustSignHash(h, s2)
	require.NoError(t, ChkSig(a, h, sig))
	require.NoError(t, ChkSig(a2, h, sig2))
	require.NotEqual(t, sig, sig2)

	// Bad address should be invalid
	require.Error(t, ChkSig(a, h, sig2))
	require.Error(t, ChkSig(a2, h, sig))
}

func TestSignHash(t *testing.T) {
	p, s := GenerateKeyPair()
	a := AddressFromPubKey(p)
	h := SumSHA256(randBytes(t, 256))
	sig := MustSignHash(h, s)
	require.NotEqual(t, sig, Sig{})
	require.NoError(t, ChkSig(a, h, sig))
}

func TestPubKeyFromSecKey(t *testing.T) {
	p, s := GenerateKeyPair()
	p2, err := PubKeyFromSecKey(s)
	require.NoError(t, err)
	require.Equal(t, p2, p)

	_, err = PubKeyFromSecKey(SecKey{})
	require.Equal(t, errors.New("PubKeyFromSecKey, attempt to load null seckey, unsafe"), err)
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

func TestVerifySignature(t *testing.T) {
	p, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	h2 := SumSHA256(randBytes(t, 256))
	sig := MustSignHash(h, s)
	require.NoError(t, VerifySignature(p, sig, h))
	require.Error(t, VerifySignature(p, Sig{}, h))
	require.Error(t, VerifySignature(p, sig, h2))
	p2, _ := GenerateKeyPair()
	require.Error(t, VerifySignature(p2, sig, h))
	require.Error(t, VerifySignature(PubKey{}, sig, h))
}

func TestGenerateKeyPair(t *testing.T) {
	for i := 0; i < 10; i++ {
		p, s := GenerateKeyPair()
		require.NoError(t, p.Verify())
		require.NoError(t, s.Verify())
		err := TestSecKey(s)
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
	require.Equal(t, errors.New("seed input is empty"), err)

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
	require.Equal(t, errors.New("seed input is empty"), err)

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
	require.Equal(t, errors.New("seed input is empty"), err)

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
	require.Equal(t, errors.New("seed input is empty"), err)

	require.Panics(t, func() {
		MustDeterministicKeyPairIterator(nil)
	})
}

func TestSecKeyTest(t *testing.T) {
	_, s := GenerateKeyPair()
	require.NoError(t, TestSecKey(s))
	require.Error(t, TestSecKey(SecKey{}))
}

func TestSecKeyHashTest(t *testing.T) {
	_, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	require.NoError(t, TestSecKeyHash(s, h))
	require.Error(t, TestSecKeyHash(SecKey{}, h))
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
	err = TestSecKeyHash(seckey, hash)
	require.NoError(t, err)
}
