package cipher

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/ripemd160"
)

func TestNewPubKey(t *testing.T) {
	require.Panics(t, func() { NewPubKey(randBytes(t, 31)) })
	require.Panics(t, func() { NewPubKey(randBytes(t, 32)) })
	require.Panics(t, func() { NewPubKey(randBytes(t, 34)) })
	require.Panics(t, func() { NewPubKey(randBytes(t, 0)) })
	require.Panics(t, func() { NewPubKey(randBytes(t, 100)) })
	require.NotPanics(t, func() { NewPubKey(randBytes(t, 33)) })
	b := randBytes(t, 33)
	p := NewPubKey(b)
	require.True(t, bytes.Equal(p[:], b))
}

func TestPubKeyFromHex(t *testing.T) {
	// Invalid hex
	require.Panics(t, func() { MustPubKeyFromHex("") })
	require.Panics(t, func() { MustPubKeyFromHex("cascs") })
	// Invalid hex length
	p := NewPubKey(randBytes(t, 33))
	s := hex.EncodeToString(p[:len(p)/2])
	require.Panics(t, func() { MustPubKeyFromHex(s) })
	// Valid
	s = hex.EncodeToString(p[:])
	require.NotPanics(t, func() { MustPubKeyFromHex(s) })
	require.Equal(t, p, MustPubKeyFromHex(s))
}

func TestPubKeyHex(t *testing.T) {
	b := randBytes(t, 33)
	p := NewPubKey(b)
	h := p.Hex()
	p2 := MustPubKeyFromHex(h)
	require.Equal(t, p2, p)
	require.Equal(t, p2.Hex(), h)
}

func TestPubKeyVerify(t *testing.T) {
	// Random bytes should not be valid, most of the time
	failed := false
	for i := 0; i < 10; i++ {
		b := randBytes(t, 33)
		if NewPubKey(b).Verify() != nil {
			failed = true
			break
		}
	}
	require.True(t, failed)
}

func TestPubKeyVerifyNil(t *testing.T) {
	// Empty public key should not be valid
	p := PubKey{}
	require.NotNil(t, p.Verify())
}

func TestPubKeyVerifyDefault1(t *testing.T) {
	// Generated pub key should be valid
	p, _ := GenerateKeyPair()
	require.Nil(t, p.Verify())
}

func TestPubKeyVerifyDefault2(t *testing.T) {
	for i := 0; i < 1024; i++ {
		p, _ := GenerateKeyPair()
		require.Nil(t, p.Verify())
	}
}

func TestPubKeyToAddressHash(t *testing.T) {
	p, _ := GenerateKeyPair()
	h := p.ToAddressHash()
	// Should be Ripemd160(SHA256(SHA256()))
	x := sha256.Sum256(p[:])
	x = sha256.Sum256(x[:])
	rh := ripemd160.New()
	err := rh.Write(x[:])
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

func TestMustNewSecKey(t *testing.T) {
	require.Panics(t, func() { NewSecKey(randBytes(t, 31)) })
	require.Panics(t, func() { NewSecKey(randBytes(t, 33)) })
	require.Panics(t, func() { NewSecKey(randBytes(t, 34)) })
	require.Panics(t, func() { NewSecKey(randBytes(t, 0)) })
	require.Panics(t, func() { NewSecKey(randBytes(t, 100)) })
	require.NotPanics(t, func() { NewSecKey(randBytes(t, 32)) })
	b := randBytes(t, 32)
	p := NewSecKey(b)
	require.True(t, bytes.Equal(p[:], b))
}

func TestMustSecKeyFromHex(t *testing.T) {
	// Invalid hex
	require.Panics(t, func() { MustSecKeyFromHex("") })
	require.Panics(t, func() { MustSecKeyFromHex("cascs") })
	// Invalid hex length
	p := NewSecKey(randBytes(t, 32))
	s := hex.EncodeToString(p[:len(p)/2])
	require.Panics(t, func() { MustSecKeyFromHex(s) })
	// Valid
	s = hex.EncodeToString(p[:])
	require.NotPanics(t, func() { MustSecKeyFromHex(s) })
	require.Equal(t, p, MustSecKeyFromHex(s))
}

func TestSecKeyHex(t *testing.T) {
	b := randBytes(t, 32)
	p := NewSecKey(b)
	h := p.Hex()
	p2 := MustSecKeyFromHex(h)
	require.Equal(t, p2, p)
	require.Equal(t, p2.Hex(), h)
}

func TestSecKeyVerify(t *testing.T) {
	// Empty secret key should not be valid
	p := SecKey{}
	require.NotNil(t, p.Verify())

	// Generated sec key should be valid
	_, p = GenerateKeyPair()
	require.Nil(t, p.Verify())

	// Random bytes are usually valid
}

func TestECDHonce(t *testing.T) {
	pub1, sec1 := GenerateKeyPair()
	pub2, sec2 := GenerateKeyPair()

	buf1 := ECDH(pub2, sec1)
	buf2 := ECDH(pub1, sec2)

	require.True(t, bytes.Equal(buf1, buf2))
}

func TestECDHloop(t *testing.T) {
	for i := 0; i < 128; i++ {
		pub1, sec1 := GenerateKeyPair()
		pub2, sec2 := GenerateKeyPair()
		buf1 := ECDH(pub2, sec1)
		buf2 := ECDH(pub1, sec2)
		require.True(t, bytes.Equal(buf1, buf2))
	}
}

func TestNewSig(t *testing.T) {
	require.Panics(t, func() { NewSig(randBytes(t, 64)) })
	require.Panics(t, func() { NewSig(randBytes(t, 66)) })
	require.Panics(t, func() { NewSig(randBytes(t, 67)) })
	require.Panics(t, func() { NewSig(randBytes(t, 0)) })
	require.Panics(t, func() { NewSig(randBytes(t, 100)) })
	require.NotPanics(t, func() { NewSig(randBytes(t, 65)) })
	b := randBytes(t, 65)
	p := NewSig(b)
	require.True(t, bytes.Equal(p[:], b))
}

func TestMustSigFromHex(t *testing.T) {
	// Invalid hex
	require.Panics(t, func() { MustSigFromHex("") })
	require.Panics(t, func() { MustSigFromHex("cascs") })
	// Invalid hex length
	p := NewSig(randBytes(t, 65))
	s := hex.EncodeToString(p[:len(p)/2])
	require.Panics(t, func() { MustSigFromHex(s) })
	// Valid
	s = hex.EncodeToString(p[:])
	require.NotPanics(t, func() { MustSigFromHex(s) })
	require.Equal(t, p, MustSigFromHex(s))
}

func TestSigHex(t *testing.T) {
	b := randBytes(t, 65)
	p := NewSig(b)
	h := p.Hex()
	p2 := MustSigFromHex(h)
	require.Equal(t, p2, p)
	require.Equal(t, p2.Hex(), h)
}

func TestChkSig(t *testing.T) {
	p, s := GenerateKeyPair()
	require.Nil(t, p.Verify())
	require.Nil(t, s.Verify())
	a := AddressFromPubKey(p)
	require.Nil(t, a.Verify(p))
	b := randBytes(t, 256)
	h := SumSHA256(b)
	sig := SignHash(h, s)
	require.Nil(t, ChkSig(a, h, sig))
	// Empty sig should be invalid
	require.NotNil(t, ChkSig(a, h, Sig{}))
	// Random sigs should not pass
	for i := 0; i < 100; i++ {
		require.NotNil(t, ChkSig(a, h, NewSig(randBytes(t, 65))))
	}
	// Sig for one hash does not work for another hash
	h2 := SumSHA256(randBytes(t, 256))
	sig2 := SignHash(h2, s)
	require.Nil(t, ChkSig(a, h2, sig2))
	require.NotNil(t, ChkSig(a, h, sig2))
	require.NotNil(t, ChkSig(a, h2, sig))

	// Different secret keys should not create same sig
	p2, s2 := GenerateKeyPair()
	a2 := AddressFromPubKey(p2)
	h = SHA256{}
	sig = SignHash(h, s)
	sig2 = SignHash(h, s2)
	require.Nil(t, ChkSig(a, h, sig))
	require.Nil(t, ChkSig(a2, h, sig2))
	require.NotEqual(t, sig, sig2)
	h = SumSHA256(randBytes(t, 256))
	sig = SignHash(h, s)
	sig2 = SignHash(h, s2)
	require.Nil(t, ChkSig(a, h, sig))
	require.Nil(t, ChkSig(a2, h, sig2))
	require.NotEqual(t, sig, sig2)

	// Bad address should be invalid
	require.NotNil(t, ChkSig(a, h, sig2))
	require.NotNil(t, ChkSig(a2, h, sig))
}

func TestSignHash(t *testing.T) {
	p, s := GenerateKeyPair()
	a := AddressFromPubKey(p)
	h := SumSHA256(randBytes(t, 256))
	sig := SignHash(h, s)
	require.NotEqual(t, sig, Sig{})
	require.Nil(t, ChkSig(a, h, sig))
}

func TestPubKeyFromSecKey(t *testing.T) {
	p, s := GenerateKeyPair()
	require.Equal(t, PubKeyFromSecKey(s), p)
	require.Panics(t, func() { PubKeyFromSecKey(SecKey{}) })
	require.Panics(t, func() { PubKeyFromSecKey(NewSecKey(randBytes(t, 99))) })
	require.Panics(t, func() { PubKeyFromSecKey(NewSecKey(randBytes(t, 31))) })
}

func TestPubKeyFromSig(t *testing.T) {
	p, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	sig := SignHash(h, s)
	p2, err := PubKeyFromSig(sig, h)
	require.Equal(t, p, p2)
	require.NoError(t, err)
	_, err = PubKeyFromSig(Sig{}, h)
	require.NotNil(t, err)
}

func TestVerifySignature(t *testing.T) {
	p, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	h2 := SumSHA256(randBytes(t, 256))
	sig := SignHash(h, s)
	require.Nil(t, VerifySignature(p, sig, h))
	require.NotNil(t, VerifySignature(p, Sig{}, h))
	require.NotNil(t, VerifySignature(p, sig, h2))
	p2, _ := GenerateKeyPair()
	require.NotNil(t, VerifySignature(p2, sig, h))
	require.NotNil(t, VerifySignature(PubKey{}, sig, h))
}

func TestGenerateKeyPair(t *testing.T) {
	p, s := GenerateKeyPair()
	require.Nil(t, p.Verify())
	require.Nil(t, s.Verify())
}

func TestGenerateDeterministicKeyPair(t *testing.T) {
	// TODO -- deterministic key pairs are useless as is because we can't
	// generate pair n+1, only pair 0
	seed := randBytes(t, 32)
	p, s := GenerateDeterministicKeyPair(seed)
	require.Nil(t, p.Verify())
	require.Nil(t, s.Verify())
	p, s = GenerateDeterministicKeyPair(seed)
	require.Nil(t, p.Verify())
	require.Nil(t, s.Verify())
}

func TestSecKeTest(t *testing.T) {
	_, s := GenerateKeyPair()
	require.Nil(t, TestSecKey(s))
	require.NotNil(t, TestSecKey(SecKey{}))
}

func TestSecKeyHashTest(t *testing.T) {
	_, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	require.Nil(t, TestSecKeyHash(s, h))
	require.NotNil(t, TestSecKeyHash(SecKey{}, h))
}

func TestGenerateDeterministicKeyPairsUsesAllBytes(t *testing.T) {
	// Tests that if a seed >128 bits is used, the generator does not ignore bits >128
	seed := "property diet little foster provide disagree witness mountain alley weekend kitten general"
	seckeys := GenerateDeterministicKeyPairs([]byte(seed), 3)
	seckeys2 := GenerateDeterministicKeyPairs([]byte(seed[:16]), 3)
	require.NotEqual(t, seckeys, seckeys2)
}
