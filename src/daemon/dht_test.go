package daemon

import (
    "github.com/nictuku/dht"
    "github.com/skycoin/pex"
    "github.com/stretchr/testify/assert"
    "testing"
)

var (
    port = 6677
)

func resetDHT() {
    DHT = nil
    dhtInfoHash = ""
}

func TestInitShutdownDHT(t *testing.T) {
    resetDHT()
    assert.Nil(t, DHT)
    assert.Equal(t, string(dhtInfoHash), "")
    InitDHT(port)
    assert.NotNil(t, DHT)
    assert.NotEqual(t, string(dhtInfoHash), "")
    go DHT.Run()
    wait()
    ShutdownDHT()
    wait()
    resetDHT()
}

func TestReceivedDHTPeers(t *testing.T) {
    Peers = pex.NewPex(maxPeers)
    m := make(map[dht.InfoHash][]string)
    peers := make([]string, 0)
    peers = append(peers, string([]byte{013, 026, 041, 054, 013, 013}))
    peers = append(peers, string([]byte{013, 026, 041, 055, 013, 013}))
    m[dht.InfoHash("")] = peers
    receivedDHTPeers(m)
    assert.Equal(t, len(Peers.Peerlist), 2)
    assert.NotNil(t, Peers.Peerlist["11.22.33.45:2827"])
    assert.NotNil(t, Peers.Peerlist["11.22.33.44:2827"])
    resetPeers()
}

func TestRequestDHTPeers(t *testing.T) {
    resetDHT()
    assert.Panics(t, RequestDHTPeers)
    InitDHT(port)
    assert.NotPanics(t, RequestDHTPeers)
    resetDHT()
}
