package pex

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestValidateAddress(t *testing.T) {
	// empty string
	assert.False(t, validateAddress("", false))
	// doubled ip:port
	assert.False(t, validateAddress("112.32.32.14:100112.32.32.14:101", false))
	// requires port
	assert.False(t, validateAddress("112.32.32.14", false))
	// not ip
	assert.False(t, validateAddress("112", false))
	assert.False(t, validateAddress("112.32", false))
	assert.False(t, validateAddress("112.32.32", false))
	// bad part
	assert.False(t, validateAddress("112.32.32.14000", false))
	// large port
	assert.False(t, validateAddress("112.32.32.14:66666", false))
	// unspecified
	assert.False(t, validateAddress("0.0.0.0:8888", false))
	// no ip
	assert.False(t, validateAddress(":8888", false))
	// multicast
	assert.False(t, validateAddress("224.1.1.1:8888", false))
	// invalid ports
	assert.False(t, validateAddress("112.32.32.14:0", false))
	assert.False(t, validateAddress("112.32.32.14:1", false))
	assert.False(t, validateAddress("112.32.32.14:10", false))
	assert.False(t, validateAddress("112.32.32.14:100", false))
	assert.False(t, validateAddress("112.32.32.14:1000", false))
	assert.False(t, validateAddress("112.32.32.14:1023", false))
	assert.False(t, validateAddress("112.32.32.14:65536", false))
	// valid ones
	assert.True(t, validateAddress("112.32.32.14:1024", false))
	assert.True(t, validateAddress("112.32.32.14:10000", false))
	assert.True(t, validateAddress("112.32.32.14:65535", false))
	// localhost is allowed
	assert.True(t, validateAddress("127.0.0.1:8888", true))
	// localhost is not allowed
	assert.False(t, validateAddress("127.0.0.1:8888", false))
}

func TestNewPex(t *testing.T) {
	dir, err := ioutil.TempDir("", "peerlist")
	require.NoError(t, err)
	defer os.Remove(dir)

	// defer removeFile()
	config := NewConfig()
	config.DataDirectory = dir

	_, err = New(config, testPeers[:])
	require.NoError(t, err)

	// check if peers are saved to disk
	peers, err := loadPeersFromFile(filepath.Join(dir, PeerDatabaseFilename))
	require.NoError(t, err)

	for _, p := range testPeers {
		v, ok := peers[p]
		require.True(t, ok)
		require.True(t, v.Trusted)
	}
}

func TestPexLoadPeers(t *testing.T) {
	tt := []struct {
		name     string
		peers    []Peer
		max      int
		expectN  int
		expectIN []Peer
	}{
		{
			"load all",
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			2,
			2,
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
		},
		{
			"reach max",
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
				Peer{Addr: testPeers[2]},
			},
			2,
			2,
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
				Peer{Addr: testPeers[2]},
			},
		},
		{
			"including invalid addr",
			[]Peer{
				Peer{Addr: wrongPortPeer},
				Peer{Addr: testPeers[1]},
				Peer{Addr: testPeers[2]},
			},
			2,
			2,
			[]Peer{
				Peer{Addr: testPeers[1]},
				Peer{Addr: testPeers[2]},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "peerlist")
			require.NoError(t, err)
			defer os.Remove(dir)

			// write peers to file
			fn := filepath.Join(dir, PeerDatabaseFilename)

			peersMap := make(map[string]Peer)
			for _, p := range tc.peers {
				peersMap[p.Addr] = p
			}

			err = file.SaveJSON(fn, peersMap, 0600)
			require.NoError(t, err)

			cfg := NewConfig()
			cfg.DataDirectory = dir
			cfg.Max = tc.max

			px := Pex{
				peerlist: *newPeerlist(),
				Config:   cfg,
			}

			px.load()

			require.Len(t, px.peers, tc.expectN)

			psm := make(map[string]Peer)
			for _, p := range tc.expectIN {
				psm[p.Addr] = p
			}

			for _, p := range px.peers {
				v, ok := psm[p.Addr]
				require.True(t, ok)
				require.Equal(t, v, *p)
			}
		})
	}
}

func TestPexAddPeer(t *testing.T) {
	tt := []struct {
		name  string
		peers []string
		max   int
		peer  string
		err   error
	}{
		{
			"ok",
			testPeers[:1],
			2,
			testPeers[1],
			nil,
		},
		{
			"invalid peer",
			testPeers[:1],
			2,
			wrongPortPeer,
			ErrInvalidAddress,
		},
		{
			"reach max",
			testPeers[:2],
			2,
			testPeers[3],
			ErrPeerlistFull,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// create temp peer list file
			dir, removeDir := preparePeerlistDir(t)
			defer removeDir()

			// create px config
			cfg := NewConfig()
			cfg.Max = tc.max
			cfg.DataDirectory = dir

			// create px instance and load peers
			px, err := New(cfg, tc.peers)
			require.NoError(t, err)

			err = px.AddPeer(tc.peer)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			// check if the peer is in the peer list
			_, ok := px.peers[tc.peer]
			require.True(t, ok)
		})
	}
}

func TestPexAddPeers(t *testing.T) {
	tt := []struct {
		name        string
		peers       []string
		max         int
		addPeers    []string
		addN        int
		expectPeers []string
	}{
		{
			"ok",
			testPeers[:1],
			5,
			testPeers[1:3],
			2,
			testPeers[1:3],
		},
		{
			"almost full",
			testPeers[:1],
			2,
			testPeers[1:3],
			1,
			testPeers[1:2],
		},
		{
			"already full",
			testPeers[:2],
			2,
			testPeers[2:3],
			0,
			testPeers[0:0],
		},
		{
			"including invalid address",
			testPeers[:1],
			2,
			[]string{testPeers[1], wrongPortPeer},
			1,
			testPeers[1:2],
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// create temp peer list file
			dir, removeDir := preparePeerlistDir(t)
			defer removeDir()

			// create px config
			cfg := NewConfig()
			cfg.Max = tc.max
			cfg.DataDirectory = dir

			// create px instance and load peers
			px, err := New(cfg, tc.peers)
			require.NoError(t, err)

			n := px.AddPeers(tc.addPeers)
			require.Equal(t, tc.addN, n)

			for _, p := range tc.expectPeers {
				_, ok := px.peers[p]
				require.True(t, ok)
			}
		})
	}
}
