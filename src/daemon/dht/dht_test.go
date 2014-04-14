package daemon_dht

import (
	"github.com/nictuku/dht"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	port = 6677
)

func testInitShutdownDHT(t *testing.T, disabled bool) {
	c := NewDHTConfig()
	c.Disabled = disabled
	d := NewDHT(c)
	assert.Equal(t, string(d.InfoHash), "")
	e := d.Init()
	assert.Nil(t, e)
	assert.NotEqual(t, string(d.InfoHash), "")
	go d.Start()
	wait()
	d.Shutdown()
	wait()
}

func TestInitShutdownDHT(t *testing.T) {
	testInitShutdownDHT(t, false)
	testInitShutdownDHT(t, true)
}

func TestReceivePeers(t *testing.T) {
	cleanupPeers()
	d := NewDHT(NewDHTConfig())
	peers := NewPeers(NewPeersConfig())
	peers.Init()
	m := make(map[dht.InfoHash][]string)
	ps := make([]string, 0)
	ps = append(ps, string([]byte{013, 026, 041, 054, 013, 013}))
	ps = append(ps, string([]byte{013, 026, 041, 055, 013, 013}))
	m[dht.InfoHash("")] = ps
	d.ReceivePeers(m, peers.Peers)
	assert.Equal(t, len(peers.Peers.Peerlist), 2)
	assert.NotNil(t, peers.Peers.Peerlist["11.22.33.45:2827"])
	assert.NotNil(t, peers.Peers.Peerlist["11.22.33.44:2827"])
	wait()
	d.Shutdown()
	wait()
}

func TestRequestDHTPeers(t *testing.T) {
	d := NewDHT(NewDHTConfig())
	assert.Panics(t, d.RequestPeers)
	e := d.Init()
	assert.Nil(t, e)
	assert.NotPanics(t, d.RequestPeers)
	wait()
	d.Config.Disabled = true
	assert.NotPanics(t, d.RequestPeers)
	wait()
	d.Shutdown()
	wait()
}
