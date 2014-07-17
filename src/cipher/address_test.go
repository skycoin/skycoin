package cipher

import (
	"testing"

	"github.com/skycoin/skycoin/src/lib/base58"
	"github.com/stretchr/testify/assert"
)

func TestMustDecodeBase58Address(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	assert.Nil(t, a.Verify(p))

	assert.Panics(t, func() { MustDecodeBase58Address("") })
	assert.Panics(t, func() { MustDecodeBase58Address("cascs") })
	b := a.Bytes()
	h := string(base58.Hex2Base58(b[:len(b)/2]))
	assert.Panics(t, func() { MustDecodeBase58Address(h) })
	h = string(base58.Hex2Base58(b))
	assert.NotPanics(t, func() { MustDecodeBase58Address(h) })
	a2 := MustDecodeBase58Address(h)
	assert.Equal(t, a, a2)
}

func TestDecodeBase58Address(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	assert.Nil(t, a.Verify(p))

	a2, err := DecodeBase58Address("")
	assert.NotNil(t, err)
	a2, err = DecodeBase58Address("cascs")
	assert.NotNil(t, err)
	b := a.Bytes()
	h := string(base58.Hex2Base58(b[:len(b)/2]))
	a2, err = DecodeBase58Address(h)
	assert.NotNil(t, err)
	h = string(base58.Hex2Base58(b))
	a2, err = DecodeBase58Address(h)
	assert.Nil(t, err)
	assert.Equal(t, a, a2)
}

func TestAddressFromBytes(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	a2, err := addressFromBytes(a.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, a2, a)
	// Invalid number of bytes
	b := a.Bytes()
	_, err = addressFromBytes(b[:len(b)-2])
	assert.NotNil(t, err)
	// Invalid checksum
	b[len(b)-1] += byte(1)
	_, err = addressFromBytes(b)
	assert.NotNil(t, err)
}

func TestAddressVerify(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	// Valid pubkey+address
	assert.Nil(t, a.Verify(p))
	// Invalid pubkey
	assert.NotNil(t, a.Verify(PubKey{}))
	p2, _ := GenerateKeyPair()
	assert.NotNil(t, a.Verify(p2))
	// Bad version
	a.Version = 0x01
	assert.NotNil(t, a.Verify(p))
}

func TestAddressString(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	s := a.String()
	a2, err := DecodeBase58Address(s)
	assert.Nil(t, err)
	assert.Equal(t, a2, a)
	s2 := a2.String()
	a3, err := DecodeBase58Address(s2)
	assert.Nil(t, err)
	assert.Equal(t, a2, a3)
}
