package cipher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/base58"
)

func TestBitcoinAddress(t *testing.T) {
	cases := []struct {
		seckey string
		pubkey string
		addr   string
	}{
		{
			seckey: "1111111111111111111111111111111111111111111111111111111111111111",
			pubkey: "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa",
			addr:   "1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9",
		},
		{
			seckey: "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
			pubkey: "02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd",
			addr:   "1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a",
		},
		{
			seckey: "47f7616ea6f9b923076625b4488115de1ef1187f760e65f89eb6f4f7ff04b012",
			pubkey: "032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3",
			addr:   "19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV",
		},
	}

	for _, tc := range cases {
		t.Run(tc.addr, func(t *testing.T) {
			seckey := MustSecKeyFromHex(tc.seckey)

			pubkey := MustPubKeyFromSecKey(seckey)
			require.Equal(t, tc.pubkey, pubkey.Hex())

			bitcoinAddr := BitcoinAddressFromPubKey(pubkey)
			require.Equal(t, tc.addr, bitcoinAddr.String())

			secAddr := MustBitcoinAddressFromSecKey(seckey)
			require.Equal(t, tc.addr, secAddr.String())

			secAddr, err := BitcoinAddressFromSecKey(seckey)
			require.NoError(t, err)
			require.Equal(t, tc.addr, secAddr.String())

			pubAddr := BitcoinAddressFromPubKey(pubkey)
			require.Equal(t, tc.addr, pubAddr.String())
		})
	}
}

func TestMustDecodeBase58BitcoinAddress(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := BitcoinAddressFromPubKey(p)
	require.NoError(t, a.Verify(p))

	require.Panics(t, func() { MustDecodeBase58BitcoinAddress("") })
	require.Panics(t, func() { MustDecodeBase58BitcoinAddress("cascs") })
	b := a.Bytes()
	h := string(base58.Encode(b[:len(b)/2]))
	require.Panics(t, func() { MustDecodeBase58BitcoinAddress(h) })
	h = string(base58.Encode(b))
	require.NotPanics(t, func() { MustDecodeBase58BitcoinAddress(h) })
	a2 := MustDecodeBase58BitcoinAddress(h)
	require.Equal(t, a, a2)

	require.NotPanics(t, func() { MustDecodeBase58BitcoinAddress(a.String()) })
	a2 = MustDecodeBase58BitcoinAddress(a.String())
	require.Equal(t, a, a2)

	// preceding whitespace is invalid
	badAddr := " " + a.String()
	require.Panics(t, func() { MustDecodeBase58BitcoinAddress(badAddr) })

	// preceding zeroes are invalid
	badAddr = "000" + a.String()
	require.Panics(t, func() { MustDecodeBase58BitcoinAddress(badAddr) })

	// trailing whitespace is invalid
	badAddr = a.String() + " "
	require.Panics(t, func() { MustDecodeBase58BitcoinAddress(badAddr) })

	// trailing zeroes are invalid
	badAddr = a.String() + "000"
	require.Panics(t, func() { MustDecodeBase58BitcoinAddress(badAddr) })

	null := "1111111111111111111111111"
	require.Panics(t, func() { MustDecodeBase58BitcoinAddress(null) })
}

func TestDecodeBase58BitcoinAddress(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := BitcoinAddressFromPubKey(p)
	require.NoError(t, a.Verify(p))

	_, err := DecodeBase58BitcoinAddress("")
	require.Error(t, err)

	_, err = DecodeBase58BitcoinAddress("cascs")
	require.Error(t, err)

	b := a.Bytes()
	h := string(base58.Encode(b[:len(b)/2]))
	_, err = DecodeBase58BitcoinAddress(h)
	require.Error(t, err)

	h = string(base58.Encode(b))
	a2, err := DecodeBase58BitcoinAddress(h)
	require.NoError(t, err)
	require.Equal(t, a, a2)

	as := a.String()
	a2, err = DecodeBase58BitcoinAddress(as)
	require.NoError(t, err)
	require.Equal(t, a, a2)

	// preceding whitespace is invalid
	as2 := " " + as
	_, err = DecodeBase58BitcoinAddress(as2)
	require.Error(t, err)

	// preceding zeroes are invalid
	as2 = "000" + as
	_, err = DecodeBase58BitcoinAddress(as2)
	require.Error(t, err)

	// trailing whitespace is invalid
	as2 = as + " "
	_, err = DecodeBase58BitcoinAddress(as2)
	require.Error(t, err)

	// trailing zeroes are invalid
	as2 = as + "000"
	_, err = DecodeBase58BitcoinAddress(as2)
	require.Error(t, err)

	// null address is invalid
	null := "1111111111111111111111111"
	_, err = DecodeBase58BitcoinAddress(null)
	require.Error(t, err)
	require.Equal(t, ErrAddressInvalidChecksum, err)
}

func TestBitcoinAddressFromBytes(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := BitcoinAddressFromPubKey(p)
	a2, err := BitcoinAddressFromBytes(a.Bytes())
	require.NoError(t, err)
	require.Equal(t, a2, a)

	// Invalid number of bytes
	b := a.Bytes()
	_, err = BitcoinAddressFromBytes(b[:len(b)-2])
	require.EqualError(t, err, "Invalid address length")

	// Invalid checksum
	b[len(b)-1] += byte(1)
	_, err = BitcoinAddressFromBytes(b)
	require.EqualError(t, err, "Invalid checksum")

	a.Version = 2
	b = a.Bytes()
	_, err = BitcoinAddressFromBytes(b)
	require.EqualError(t, err, "Address version invalid")
}

func TestMustBitcoinAddressFromBytes(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := BitcoinAddressFromPubKey(p)
	a2 := MustBitcoinAddressFromBytes(a.Bytes())
	require.Equal(t, a2, a)

	// Invalid number of bytes
	b := a.Bytes()
	require.Panics(t, func() {
		MustBitcoinAddressFromBytes(b[:len(b)-2])
	})

	// Invalid checksum
	b[len(b)-1] += byte(1)
	require.Panics(t, func() {
		MustBitcoinAddressFromBytes(b)
	})

	a.Version = 2
	b = a.Bytes()
	require.Panics(t, func() {
		MustBitcoinAddressFromBytes(b)
	})
}

func TestBitcoinAddressFromSecKey(t *testing.T) {
	p, s := GenerateKeyPair()
	a, err := BitcoinAddressFromSecKey(s)
	require.NoError(t, err)
	// Valid pubkey+address
	require.NoError(t, a.Verify(p))

	_, err = BitcoinAddressFromSecKey(SecKey{})
	require.Equal(t, errors.New("Attempt to load null seckey, unsafe"), err)
}

func TestMustBitcoinAddressFromSecKey(t *testing.T) {
	p, s := GenerateKeyPair()
	a := MustBitcoinAddressFromSecKey(s)
	// Valid pubkey+address
	require.NoError(t, a.Verify(p))

	require.Panics(t, func() {
		MustBitcoinAddressFromSecKey(SecKey{})
	})
}

func TestBitcoinAddressNull(t *testing.T) {
	var a BitcoinAddress
	require.True(t, a.Null())

	p, _ := GenerateKeyPair()
	a = BitcoinAddressFromPubKey(p)
	require.False(t, a.Null())
}

func TestBitcoinAddressVerify(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := BitcoinAddressFromPubKey(p)
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

func TestBitcoinWIFRoundTrip(t *testing.T) {
	_, seckey1 := GenerateKeyPair()
	wif1 := BitcoinWalletImportFormatFromSeckey(seckey1)
	seckey2, err := SecKeyFromBitcoinWalletImportFormat(wif1)
	wif2 := BitcoinWalletImportFormatFromSeckey(seckey2)

	require.NoError(t, err)
	require.Equal(t, seckey1, seckey2)
	require.Equal(t, seckey1.Hex(), seckey2.Hex())
	require.Equal(t, wif1, wif2)
}

func TestBitcoinWIF(t *testing.T) {
	cases := []struct {
		wif    string
		pubkey string
		addr   string
	}{
		{
			wif:    "KwntMbt59tTsj8xqpqYqRRWufyjGunvhSyeMo3NTYpFYzZbXJ5Hp",
			pubkey: "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa",
			addr:   "1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9",
		},
		{
			wif:    "L4ezQvyC6QoBhxB4GVs9fAPhUKtbaXYUn8YTqoeXwbevQq4U92vN",
			pubkey: "02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd",
			addr:   "1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a",
		},
		{
			wif:    "KydbzBtk6uc7M6dXwEgTEH2sphZxSPbmDSz6kUUHi4eUpSQuhEbq",
			pubkey: "032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3",
			addr:   "19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV",
		},
	}

	for _, tc := range cases {
		t.Run(tc.addr, func(t *testing.T) {
			seckey, err := SecKeyFromBitcoinWalletImportFormat(tc.wif)
			require.NoError(t, err)

			require.NotPanics(t, func() {
				MustSecKeyFromBitcoinWalletImportFormat(tc.wif)
			})

			pubkey := MustPubKeyFromSecKey(seckey)
			require.Equal(t, tc.pubkey, pubkey.Hex())

			bitcoinAddr := BitcoinAddressFromPubKey(pubkey)
			require.Equal(t, tc.addr, bitcoinAddr.String())
		})
	}
}

func TestBitcoinWIFFailures(t *testing.T) {
	a := " asdio"
	_, err := SecKeyFromBitcoinWalletImportFormat(a)
	require.Equal(t, errors.New("Invalid base58 character"), err)

	a = string(base58.Encode(randBytes(t, 37)))
	_, err = SecKeyFromBitcoinWalletImportFormat(a)
	require.Equal(t, errors.New("Invalid length"), err)

	a = string(base58.Encode(randBytes(t, 39)))
	_, err = SecKeyFromBitcoinWalletImportFormat(a)
	require.Equal(t, errors.New("Invalid length"), err)

	b := randBytes(t, 38)
	b[0] = 0x70
	a = string(base58.Encode(b))
	_, err = SecKeyFromBitcoinWalletImportFormat(a)
	require.Equal(t, errors.New("Bitcoin WIF: First byte invalid"), err)

	b = randBytes(t, 38)
	b[0] = 0x80
	b[33] = 0x02
	a = string(base58.Encode(b))
	_, err = SecKeyFromBitcoinWalletImportFormat(a)
	require.Equal(t, errors.New("Bitcoin WIF: Invalid 33rd byte"), err)

	b = randBytes(t, 38)
	b[0] = 0x80
	b[33] = 0x01
	hashed := DoubleSHA256(b[0:34])
	chksum := hashed[0:4]
	chksum[0] = chksum[0] + 1
	copy(b[34:38], chksum[:])
	a = string(base58.Encode(b))
	_, err = SecKeyFromBitcoinWalletImportFormat(a)
	require.Equal(t, errors.New("Bitcoin WIF: Checksum fail"), err)
}

func TestMustBitcoinWIFFailures(t *testing.T) {
	a := " asdio"
	require.Panics(t, func() {
		MustSecKeyFromBitcoinWalletImportFormat(a)
	})

	a = string(base58.Encode(randBytes(t, 37)))
	require.Panics(t, func() {
		MustSecKeyFromBitcoinWalletImportFormat(a)
	})

	a = string(base58.Encode(randBytes(t, 39)))
	require.Panics(t, func() {
		MustSecKeyFromBitcoinWalletImportFormat(a)
	})

	b := randBytes(t, 38)
	b[0] = 0x70
	a = string(base58.Encode(b))
	require.Panics(t, func() {
		MustSecKeyFromBitcoinWalletImportFormat(a)
	})

	b = randBytes(t, 38)
	b[0] = 0x80
	b[33] = 0x02
	a = string(base58.Encode(b))
	require.Panics(t, func() {
		MustSecKeyFromBitcoinWalletImportFormat(a)
	})

	b = randBytes(t, 38)
	b[0] = 0x80
	b[33] = 0x01
	hashed := DoubleSHA256(b[0:34])
	chksum := hashed[0:4]
	chksum[0] = chksum[0] + 1
	copy(b[34:38], chksum[:])
	a = string(base58.Encode(b))
	require.Panics(t, func() {
		MustSecKeyFromBitcoinWalletImportFormat(a)
	})
}
