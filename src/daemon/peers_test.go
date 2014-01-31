package daemon

import (
    "github.com/skycoin/pex"
    "github.com/stretchr/testify/assert"
    "os"
    "strings"
    "testing"
)

func TestInitPeers(t *testing.T) {
    defer cleanupPeers()
    c := NewPeersConfig()
    peers := NewPeers(c)

    // Write dummy peer db
    fn := "./" + pex.PeerDatabaseFilename
    cleanupPeers()
    f, err := os.Create(fn)
    assert.Nil(t, err)
    if err != nil {
        t.Fatalf("Error creating %s", fn)
    }
    _, err = f.Write([]byte(addr + " 0\n"))
    assert.Nil(t, err)
    f.Close()

    peers.Config.DataDirectory = "./"
    assert.NotPanics(t, func() { peers.Init() })
    assert.Equal(t, len(peers.Peers.Peerlist), 1)
    assert.NotNil(t, peers.Peers.Peerlist[addr])
}

func TestShutdownPeers(t *testing.T) {
    defer cleanupPeers()
    peers := setupPeersShutdown(t)
    peers.Shutdown()
    confirmPeersShutdown(t)
}

func setupPeersShutdown(t *testing.T) *Peers {
    cleanupPeers()
    fn := "./" + pex.PeerDatabaseFilename
    _, err := os.Stat(fn)
    if err == nil {
        os.Remove(fn)
    }
    peers := NewPeers(NewPeersConfig())
    peers.Init()
    _, err = peers.Peers.AddPeer(addr)
    assert.Nil(t, err)
    return peers
}

func confirmPeersShutdown(t *testing.T) {
    f, err := os.Open("./" + pex.PeerDatabaseFilename)
    assert.Nil(t, err)
    if err != nil {
        t.Fatalf("Failed to open %s", "./"+pex.PeerDatabaseFilename)
    }
    b := make([]byte, len(addr)*2)
    n, err := f.Read(b)
    assert.Nil(t, err)
    assert.Equal(t, strings.Split(string(b[:n]), " ")[0], addr)
}

func cleanupPeers() {
    os.Remove("./" + pex.BlacklistedDatabaseFilename)
    os.Remove("./" + pex.PeerDatabaseFilename)
}
