package daemon

import (
    "github.com/skycoin/pex"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestInitPeers(t *testing.T) {
    Peers = pex.NewPex(maxPeers)
    assert.Panics(t, func() { InitPeers("x") })
    Peers = nil

    // Write dummy peer db
    fn := "./" + pex.PeerDatabaseFilename
    os.Remove(fn)
    os.Remove("./" + pex.BlacklistedDatabaseFilename)
    f, err := os.Create(fn)
    defer os.Remove(fn)
    defer os.Remove("./" + pex.BlacklistedDatabaseFilename)
    assert.Nil(t, err)
    if err != nil {
        t.Fatalf("Error creating %s", fn)
    }
    _, err = f.Write([]byte(addr + "\n"))
    assert.Nil(t, err)
    f.Close()

    assert.NotPanics(t, func() { InitPeers("./") })
    assert.NotNil(t, Peers)
    assert.Equal(t, len(Peers.Peerlist), 1)
    assert.NotNil(t, Peers.Peerlist[addr])
}

func TestShutdownPeers(t *testing.T) {
    SetupPeersShutdown(t)
    ShutdownPeers("./")
    ConfirmPeersShutdown(t)
}

func SetupPeersShutdown(t *testing.T) {
    os.Remove("./" + pex.BlacklistedDatabaseFilename)
    os.Remove("./" + pex.PeerDatabaseFilename)
    fn := "./" + pex.PeerDatabaseFilename
    _, err := os.Stat(fn)
    if err == nil {
        os.Remove(fn)
    }
    Peers = pex.NewPex(maxPeers)
    _, err = Peers.AddPeer(addr)
    assert.Nil(t, err)
}

func ConfirmPeersShutdown(t *testing.T) {
    defer os.Remove("./" + pex.BlacklistedDatabaseFilename)
    defer os.Remove("./" + pex.PeerDatabaseFilename)
    assert.Nil(t, Peers)

    f, err := os.Open("./" + pex.PeerDatabaseFilename)
    assert.Nil(t, err)
    if err != nil {
        t.Fatalf("Failed to open %s", "./"+pex.PeerDatabaseFilename)
    }
    b := make([]byte, len(addr)*2)
    n, err := f.Read(b)
    assert.Nil(t, err)
    assert.Equal(t, string(b[:n]), addr+"\n")
}
