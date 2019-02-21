package cipher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/base58"
)

func TestMustDecodeBase58Address(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	require.NoError(t, a.Verify(p))

	require.Panics(t, func() { MustDecodeBase58Address("") })
	require.Panics(t, func() { MustDecodeBase58Address("cascs") })
	b := a.Bytes()
	h := string(base58.Encode(b[:len(b)/2]))
	require.Panics(t, func() { MustDecodeBase58Address(h) })
	h = string(base58.Encode(b))
	require.NotPanics(t, func() { MustDecodeBase58Address(h) })
	a2 := MustDecodeBase58Address(h)
	require.Equal(t, a, a2)

	require.NotPanics(t, func() { MustDecodeBase58Address(a.String()) })
	a2 = MustDecodeBase58Address(a.String())
	require.Equal(t, a, a2)

	// preceding whitespace is invalid
	badAddr := " " + a.String()
	require.Panics(t, func() { MustDecodeBase58Address(badAddr) })

	// preceding zeroes are invalid
	badAddr = "000" + a.String()
	require.Panics(t, func() { MustDecodeBase58Address(badAddr) })

	// trailing whitespace is invalid
	badAddr = a.String() + " "
	require.Panics(t, func() { MustDecodeBase58Address(badAddr) })

	// trailing zeroes are invalid
	badAddr = a.String() + "000"
	require.Panics(t, func() { MustDecodeBase58Address(badAddr) })

	null := "1111111111111111111111111"
	require.Panics(t, func() { MustDecodeBase58Address(null) })
}

func TestDecodeBase58Address(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	require.NoError(t, a.Verify(p))

	_, err := DecodeBase58Address("")
	require.Error(t, err)

	_, err = DecodeBase58Address("cascs")
	require.Error(t, err)

	b := a.Bytes()
	h := string(base58.Encode(b[:len(b)/2]))
	_, err = DecodeBase58Address(h)
	require.Error(t, err)

	h = string(base58.Encode(b))
	a2, err := DecodeBase58Address(h)
	require.NoError(t, err)
	require.Equal(t, a, a2)

	as := a.String()
	a2, err = DecodeBase58Address(as)
	require.NoError(t, err)
	require.Equal(t, a, a2)

	// preceding whitespace is invalid
	as2 := " " + as
	_, err = DecodeBase58Address(as2)
	require.Error(t, err)

	// preceding zeroes are invalid
	as2 = "000" + as
	_, err = DecodeBase58Address(as2)
	require.Error(t, err)

	// trailing whitespace is invalid
	as2 = as + " "
	_, err = DecodeBase58Address(as2)
	require.Error(t, err)

	// trailing zeroes are invalid
	as2 = as + "000"
	_, err = DecodeBase58Address(as2)
	require.Error(t, err)

	// null address is invalid
	null := "1111111111111111111111111"
	_, err = DecodeBase58Address(null)
	require.Error(t, err)
	require.Equal(t, ErrAddressInvalidChecksum, err)
}

func TestAddressFromBytes(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	a2, err := AddressFromBytes(a.Bytes())
	require.NoError(t, err)
	require.Equal(t, a2, a)

	// Invalid number of bytes
	b := a.Bytes()
	_, err = AddressFromBytes(b[:len(b)-2])
	require.EqualError(t, err, "Invalid address length")

	// Invalid checksum
	b[len(b)-1] += byte(1)
	_, err = AddressFromBytes(b)
	require.EqualError(t, err, "Invalid checksum")

	a.Version = 2
	b = a.Bytes()
	_, err = AddressFromBytes(b)
	require.EqualError(t, err, "Address version invalid")
}

func TestMustAddressFromBytes(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	a2 := MustAddressFromBytes(a.Bytes())
	require.Equal(t, a2, a)

	// Invalid number of bytes
	b := a.Bytes()
	require.Panics(t, func() {
		MustAddressFromBytes(b[:len(b)-2])
	})

	// Invalid checksum
	b[len(b)-1] += byte(1)
	require.Panics(t, func() {
		MustAddressFromBytes(b)
	})

	a.Version = 2
	b = a.Bytes()
	require.Panics(t, func() {
		MustAddressFromBytes(b)
	})
}

func TestAddressRoundtrip(t *testing.T) {
	// Tests encode and decode
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	a2, err := AddressFromBytes(a.Bytes())
	require.NoError(t, err)
	require.Equal(t, a, a2)
	require.Equal(t, a.String(), a2.String())
}

func TestAddressVerify(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	// Valid pubkey+address
	require.NoError(t, a.Verify(p))
	// Invalid pubkey
	require.Error(t, a.Verify(PubKey{}))
	p2, _ := GenerateKeyPair()
	require.Error(t, a.Verify(p2))
	// Bad version
	a.Version = 0x01
	require.Error(t, a.Verify(p))
}

func TestAddressString(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	s := a.String()
	a2, err := DecodeBase58Address(s)
	require.NoError(t, err)
	require.Equal(t, a2, a)
	s2 := a2.String()
	a3, err := DecodeBase58Address(s2)
	require.NoError(t, err)
	require.Equal(t, a2, a3)
}

func TestAddressBulk(t *testing.T) {
	for i := 0; i < 1024; i++ {
		pub, _, err := GenerateDeterministicKeyPair(RandByte(32))
		require.NoError(t, err)

		a := AddressFromPubKey(pub)
		require.NoError(t, a.Verify(pub))
		s := a.String()
		a2, err := DecodeBase58Address(s)
		require.NoError(t, err)
		require.Equal(t, a2, a)
	}
}

func TestAddressNull(t *testing.T) {
	var a Address
	require.True(t, a.Null())

	p, _ := GenerateKeyPair()
	a = AddressFromPubKey(p)
	require.False(t, a.Null())
}

func TestAddressFromSecKey(t *testing.T) {
	p, s := GenerateKeyPair()
	a, err := AddressFromSecKey(s)
	require.NoError(t, err)
	// Valid pubkey+address
	require.NoError(t, a.Verify(p))

	_, err = AddressFromSecKey(SecKey{})
	require.Equal(t, errors.New("Attempt to load null seckey, unsafe"), err)
}

func TestMustAddressFromSecKey(t *testing.T) {
	p, s := GenerateKeyPair()
	a := MustAddressFromSecKey(s)
	// Valid pubkey+address
	require.NoError(t, a.Verify(p))

	require.Panics(t, func() {
		MustAddressFromSecKey(SecKey{})
	})
}
