package daemon

// import (
// 	"os"
// 	"strings"
// 	"testing"

// 	//"github.com/skycoin/skycoin/src/daemon/gnet"
// 	"github.com/skycoin/skycoin/src/daemon/gnet"
// 	"github.com/skycoin/skycoin/src/daemon/pex"
// 	"github.com/stretchr/testify/assert"
// )

// func TestInitPeers(t *testing.T) {
// 	defer cleanupPeers()
// 	c := NewPeersConfig()
// 	peers := NewPeers(c)

// 	// Write dummy peer db
// 	fn := "./" + pex.PeerDatabaseFilename
// 	cleanupPeers()
// 	f, err := os.Create(fn)
// 	assert.Nil(t, err)
// 	if err != nil {
// 		t.Fatalf("Error creating %s", fn)
// 	}
// 	_, err = f.Write([]byte(addr + " 0 0\n"))
// 	assert.Nil(t, err)
// 	f.Close()

// 	peers.Config.DataDirectory = "./"
// 	assert.NotPanics(t, func() { peers.Init() })
// 	assert.Equal(t, len(peers.Peers.Peerlist), 1)
// 	assert.NotNil(t, peers.Peers.Peerlist[addr])
// 	assert.False(t, peers.Peers.AllowLocalhost)

// 	peers.Config.AllowLocalhost = true
// 	assert.NotPanics(t, func() { peers.Init() })
// 	assert.True(t, peers.Peers.AllowLocalhost)
// }

// func TestShutdownPeers(t *testing.T) {
// 	defer cleanupPeers()
// 	peers := setupPeersShutdown(t)
// 	peers.Shutdown()
// 	confirmPeersShutdown(t)
// }

// func TestShutdownPeersNil(t *testing.T) {
// 	p := NewPeers(NewPeersConfig())
// 	assert.Nil(t, p.Peers)
// 	assert.Nil(t, p.Shutdown())
// }

// func TestRequestPeers(t *testing.T) {
// 	defer cleanupPeers()
// 	gnet.EraseMessages()
// 	defer gnet.EraseMessages()

// 	// Disabled
// 	c := NewPeersConfig()
// 	c.Disabled = true
// 	p := NewPeers(c)
// 	assert.NotPanics(t, func() { p.requestPeers(nil) })

// 	// Full
// 	c.Disabled = false
// 	c.Max = 1
// 	p = NewPeers(c)
// 	p.Init()
// 	p.Peers.AddPeer(addr)
// 	assert.NotPanics(t, func() { p.requestPeers(nil) })

// 	// Not full, will send message
// 	c.Max = 10
// 	p = NewPeers(c)
// 	p.Init()
// 	p.Peers.AddPeer(addr)
// 	pool, gc := setupPool()
// 	go pool.Pool.ConnectionWriteLoop(gc)
// 	defer gc.Close()
// 	assert.NotPanics(t, func() { p.requestPeers(pool) })
// 	wait()
// 	assert.Equal(t, len(pool.Pool.SendResults), 1)
// 	if len(pool.Pool.SendResults) == 0 {
// 		t.Fatalf("SendResults empty, would block")
// 	}
// 	sr := <-pool.Pool.SendResults
// 	assert.Equal(t, sr.Connection, gc)
// 	assert.Nil(t, sr.Error)
// 	_, ok := sr.Message.(*GetPeersMessage)
// 	assert.True(t, ok)
// }

// func TestRemovePeer(t *testing.T) {
// 	cleanupPeers()
// 	defer cleanupPeers()
// 	p := NewPeers(NewPeersConfig())
// 	p.Init()
// 	assert.Equal(t, len(p.Peers.Peerlist), 0)
// 	q, err := p.Peers.AddPeer(addr)
// 	assert.Equal(t, len(p.Peers.Peerlist), 1)
// 	assert.Nil(t, err)
// 	q.Private = true
// 	p.RemovePeer(addr)
// 	assert.Equal(t, len(p.Peers.Peerlist), 1)
// 	q.Private = false
// 	p.RemovePeer(addr)
// 	assert.Equal(t, len(p.Peers.Peerlist), 0)
// 	assert.NotPanics(t, func() { p.RemovePeer(addr) })
// }

// func setupPeersShutdown(t *testing.T) *Peers {
// 	cleanupPeers()
// 	fn := "./" + pex.PeerDatabaseFilename
// 	_, err := os.Stat(fn)
// 	if err == nil {
// 		os.Remove(fn)
// 	}
// 	peers := NewPeers(NewPeersConfig())
// 	peers.Init()
// 	_, err = peers.Peers.AddPeer(addr)
// 	assert.Nil(t, err)
// 	return peers
// }

// func confirmPeersShutdown(t *testing.T) {
// 	f, err := os.Open("./" + pex.PeerDatabaseFilename)
// 	assert.Nil(t, err)
// 	if err != nil {
// 		t.Fatalf("Failed to open %s", "./"+pex.PeerDatabaseFilename)
// 	}
// 	b := make([]byte, len(addr)*2)
// 	n, err := f.Read(b)
// 	assert.Nil(t, err)
// 	assert.Equal(t, strings.Split(string(b[:n]), " ")[0], addr)
// }

// func cleanupPeers() {
// 	os.Remove("./" + pex.BlacklistedDatabaseFilename)
// 	os.Remove("./" + pex.PeerDatabaseFilename)
// }
