package cipher

import (
	"testing"

	"github.com/skycoin/skycoin/src/cipher/base58"
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

//encode and decode
func TestAddressRoundtrip(t *testing.T) {
	p, _ := GenerateKeyPair()
	a := AddressFromPubKey(p)
	a2, err := addressFromBytes(a.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, a, a2)
	assert.Equal(t, a.String(), a2.String())
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

func TestBitcoinAddress1(t *testing.T) {
	seckey := MustSecKeyFromHex("1111111111111111111111111111111111111111111111111111111111111111")
	pubkey := PubKeyFromSecKey(seckey)
	pubkeyStr := "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa"
	assert.Equal(t, pubkeyStr, pubkey.Hex())
	bitcoinStr := "1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9"
	bitcoinAddr := BitcoinAddressFromPubkey(pubkey)
	assert.Equal(t, bitcoinStr, bitcoinAddr)
}

func TestBitcoinAddress2(t *testing.T) {
	seckey := MustSecKeyFromHex("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	pubkey := PubKeyFromSecKey(seckey)
	pubkeyStr := "02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd"
	assert.Equal(t, pubkeyStr, pubkey.Hex())
	bitcoinStr := "1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a"
	bitcoinAddr := BitcoinAddressFromPubkey(pubkey)
	assert.Equal(t, bitcoinStr, bitcoinAddr)
}

func TestBitcoinAddress3(t *testing.T) {
	seckey := MustSecKeyFromHex("47f7616ea6f9b923076625b4488115de1ef1187f760e65f89eb6f4f7ff04b012")
	pubkey := PubKeyFromSecKey(seckey)
	pubkeyStr := "032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3"
	assert.Equal(t, pubkeyStr, pubkey.Hex())
	bitcoinStr := "19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV"
	bitcoinAddr := BitcoinAddressFromPubkey(pubkey)
	assert.Equal(t, bitcoinStr, bitcoinAddr)
}

func TestBitcoinWIPRoundTrio(t *testing.T) {

	_, seckey1 := GenerateKeyPair()
	wip1 := BitcoinWalletImportFormatFromSeckey(seckey1)
	seckey2, err := SecKeyFromWalletImportFormat(wip1)
	wip2 := BitcoinWalletImportFormatFromSeckey(seckey2)

	assert.Nil(t, err)
	assert.Equal(t, seckey1, seckey2)
	assert.Equal(t, seckey1.Hex(), seckey2.Hex())
	assert.Equal(t, wip1, wip2)

}

func TestBitcoinWIP(t *testing.T) {
	//wallet input format string
	var wip = []string{
		"KwntMbt59tTsj8xqpqYqRRWufyjGunvhSyeMo3NTYpFYzZbXJ5Hp",
		"L4ezQvyC6QoBhxB4GVs9fAPhUKtbaXYUn8YTqoeXwbevQq4U92vN",
		"KydbzBtk6uc7M6dXwEgTEH2sphZxSPbmDSz6kUUHi4eUpSQuhEbq",
	}
	//the expected pubkey to generate
	var pub = []string{
		"034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa",
		"02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd",
		"032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3",
	}
	//the expected addrss to generate
	var addr = []string{
		"1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9",
		"1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a",
		"19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV",
	}

	for i := range wip {
		seckey, err := SecKeyFromWalletImportFormat(wip[i])
		assert.Equal(t, nil, err)
		_ = MustSecKeyFromWalletImportFormat(wip[i])
		pubkey := PubKeyFromSecKey(seckey)
		assert.Equal(t, pub[i], pubkey.Hex())
		bitcoinAddr := BitcoinAddressFromPubkey(pubkey)
		assert.Equal(t, addr[i], bitcoinAddr)
	}

	/*
		seckey := MustSecKeyFromHex("47f7616ea6f9b923076625b4488115de1ef1187f760e65f89eb6f4f7ff04b012")
		pubkey := PubKeyFromSecKey(seckey)
		pubkey_str := "032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3"
		assert.Equal(t, pubkey_str, pubkey.Hex())
		bitcoin_str := "19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV"
		bitcoin_addr := BitcoinAddressFromPubkey(pubkey)
		assert.Equal(t, bitcoin_str, bitcoin_addr)
	*/
}

func TestAddressBulk(t *testing.T) {

	for i := 0; i < 1024; i++ {
		pub, _ := GenerateDeterministicKeyPair(RandByte(32))

		a := AddressFromPubKey(pub)
		assert.Nil(t, a.Verify(pub))
		s := a.String()
		a2, err := DecodeBase58Address(s)
		assert.Nil(t, err)
		assert.Equal(t, a2, a)

	}
}
