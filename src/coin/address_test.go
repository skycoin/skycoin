package coin

import (
    "bytes"
    "github.com/skycoin/skycoin/src/lib/base58"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestVersionByName(t *testing.T) {
    v, ok := VersionByName("main")
    assert.True(t, ok)
    assert.Equal(t, v, byte(0x0F))
    v, ok = VersionByName("test")
    assert.True(t, ok)
    assert.Equal(t, v, byte(0x1F))
    v, ok = VersionByName("bad")
    assert.False(t, ok)
    for k, v := range addressVersions {
        w, ok := VersionByName(k)
        assert.True(t, ok)
        assert.Equal(t, w, v)
    }
}

func TestMustVersionByName(t *testing.T) {
    assert.Panics(t, func() { SetAddressVersion("bad") })
    for k, _ := range addressVersions {
        assert.NotPanics(t, func() { SetAddressVersion(k) })
    }
    assert.NotPanics(t, func() { SetAddressVersion("main") })
    assert.NotPanics(t, func() { SetAddressVersion("test") })
}

func TestSetAddressVersion(t *testing.T) {
    // "test" should be the default address version
    assert.Equal(t, addressVersion, byte(0x1F))
    assert.Panics(t, func() { SetAddressVersion("bad") })
    SetAddressVersion("test")
    assert.Equal(t, addressVersion, byte(0x1F))
    SetAddressVersion("main")
    assert.Equal(t, addressVersion, byte(0x0F))
    SetAddressVersion("test")
    assert.Equal(t, addressVersion, byte(0x1F))
    SetAddressVersion("main")
    assert.Equal(t, addressVersion, byte(0x0F))
    SetAddressVersion("test")
    assert.Equal(t, addressVersion, byte(0x1F))
}

func TestAddressFromPubKey(t *testing.T) {
    // Test with addr version "test"
    p, _ := GenerateKeyPair()
    a := AddressFromPubKey(p)
    assert.Nil(t, a.Verify(p))
    assert.NotEqual(t, a, Address{})
    assert.NotEqual(t, a.Checksum, Checksum{})
    v, _ := VersionByName("test")
    assert.Equal(t, a.Version, v)

    // Test with addr version "main"
    SetAddressVersion("main")
    a = AddressFromPubKey(p)
    assert.Nil(t, a.Verify(p))
    assert.NotEqual(t, a, Address{})
    assert.NotEqual(t, a.Checksum, Checksum{})
    v, _ = VersionByName("main")
    assert.Equal(t, a.Version, v)
    SetAddressVersion("test")
}

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
    v, _ := VersionByName("test")
    assert.Equal(t, v, addressVersion)
    p, _ := GenerateKeyPair()
    a := AddressFromPubKey(p)
    // Valid pubkey+address
    assert.Nil(t, a.Verify(p))
    // Invalid pubkey
    assert.NotNil(t, a.Verify(PubKey{}))
    p2, _ := GenerateKeyPair()
    assert.NotNil(t, a.Verify(p2))
    // Invalid checksum
    a.Checksum[0] += byte(1)
    assert.NotNil(t, a.Verify(p))
    // Bad version
    a.Version = 0x01
    assert.NotNil(t, a.Verify(p))
    // Different version than the default
    v, _ = VersionByName("main")
    a.Version = v
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

func TestAddressBytes(t *testing.T) {
    p, _ := GenerateKeyPair()
    a := AddressFromPubKey(p)
    b := a.Bytes()
    assert.True(t, bytes.Equal(b[:20], a.Key[:]))
    assert.Equal(t, b[20], a.Version)
    assert.True(t, bytes.Equal(b[21:], a.Checksum[:]))
    a2, err := addressFromBytes(b)
    assert.Nil(t, err)
    assert.Equal(t, a, a2)
}

func TestAddressChecksum(t *testing.T) {
    p, _ := GenerateKeyPair()
    a := AddressFromPubKey(p)
    c := a.CreateChecksum()
    assert.Equal(t, a.Checksum, c)
    h := SumSHA256(append(a.Key[:], []byte{a.Version}...))
    assert.True(t, bytes.Equal(c[:], h[:4]))
    assert.NotEqual(t, c, Checksum{})
}

func TestAddressHasValidChecksum(t *testing.T) {
    p, _ := GenerateKeyPair()
    a := AddressFromPubKey(p)
    assert.True(t, a.HasValidChecksum())
    a.Checksum[0] += byte(1)
    assert.False(t, a.HasValidChecksum())
    v, _ := VersionByName("test")
    assert.Equal(t, v, addressVersion)
    v, _ = VersionByName("main")
    a.Version = v
    assert.False(t, a.HasValidChecksum())
    a = AddressFromPubKey(p)
    a.Key[0] += byte(1)
    assert.False(t, a.HasValidChecksum())
    a.Key = Ripemd160{}
    assert.False(t, a.HasValidChecksum())
    a = AddressFromPubKey(p)
    a.Checksum = Checksum{}
    assert.False(t, a.HasValidChecksum())
}

func TestAddressSetChecksum(t *testing.T) {
    p, _ := GenerateKeyPair()
    a := AddressFromPubKey(p)
    assert.True(t, a.HasValidChecksum())
    checksum := a.Checksum
    a.Checksum = Checksum{}
    assert.False(t, a.HasValidChecksum())
    a.setChecksum()
    assert.True(t, a.HasValidChecksum())
    assert.Equal(t, a.Checksum, checksum)
}
