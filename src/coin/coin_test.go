package coin

import (
    //"crypto/sha256"
    //"hash"
    "encoding/hex"
    "log"
    "testing"
)


func TestAddress1(t *testing.T) {
	a := "02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8"
	b, err := hex.DecodeString(a)
	if err != nil {
		log.Panic(err)
	}
	addr := AddressFromRawPubKey(b)
	_ = addr

	///func SignHash(hash SHA256, sec SecKey) (Sig, error) {

}

func TestAddress2(t *testing.T) {
	a := "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d"
	b, err := hex.DecodeString(a)
	if err != nil {
		log.Panic(err)
	}

    seckey := 
	addr := AddressFromRawPubKey(b)
	_ = addr

	///func SignHash(hash SHA256, sec SecKey) (Sig, error) {

}

/*
func TestGetListenPort(t *testing.T) {
    // No connectionMirror found
    assert.Equal(t, getListenPort(addr), uint16(0))
    // No mirrorConnection map exists
    connectionMirrors[addr] = uint32(4)
    assert.Panics(t, func() { getListenPort(addr) })
    // Everything is good
    m := make(map[string]uint16)
    mirrorConnections[uint32(4)] = m
    m[addrIP] = uint16(6667)
    assert.Equal(t, getListenPort(addr), uint16(6667))

    // cleanup
    delete(mirrorConnections, uint32(4))
    delete(connectionMirrors, addr)
}
*/