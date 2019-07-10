package pex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/file"
)

func TestValidateAddress(t *testing.T) {
	cases := []struct {
		addr           string
		allowLocalhost bool
		err            error
		cleanAddr      string
	}{
		{
			addr:           "",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "112.32.32.14:100112.32.32.14:101",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "112.32.32.14",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "112",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "112.32",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "112.32.32",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "112.32.32.14000",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "112.32.32.14:66666",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "0.0.0.0:8888",
			allowLocalhost: false,
			err:            ErrNotExternalIP,
		},
		{
			addr:           ":8888",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "224.1.1.1:8888",
			allowLocalhost: false,
			err:            ErrNotExternalIP,
		},
		{
			addr:           "112.32.32.14:0",
			allowLocalhost: false,
			err:            ErrPortTooLow,
		},
		{
			addr:           "112.32.32.14:1",
			allowLocalhost: false,
			err:            ErrPortTooLow,
		},
		{
			addr:           "112.32.32.14:10",
			allowLocalhost: false,
			err:            ErrPortTooLow,
		},
		{
			addr:           "112.32.32.14:100",
			allowLocalhost: false,
			err:            ErrPortTooLow,
		},
		{
			addr:           "112.32.32.14:1000",
			allowLocalhost: false,
			err:            ErrPortTooLow,
		},
		{
			addr:           "112.32.32.14:1023",
			allowLocalhost: false,
			err:            ErrPortTooLow,
		},
		{
			addr:           "112.32.32.14:65536",
			allowLocalhost: false,
			err:            ErrInvalidAddress,
		},
		{
			addr:           "112.32.32.14:1024",
			allowLocalhost: false,
		},
		{
			addr:           "112.32.32.14:10000",
			allowLocalhost: false,
		},
		{
			addr:           "112.32.32.14:65535",
			allowLocalhost: false,
		},
		{
			addr:           "127.0.0.1:8888",
			allowLocalhost: true,
		},
		{
			addr:           "127.0.0.1:8888",
			allowLocalhost: false,
			err:            ErrNoLocalhost,
		},
		{
			addr:           "11.22.33.44:8080",
			allowLocalhost: false,
		},
		{
			addr:           " 11.22.33.44:8080\n",
			allowLocalhost: false,
			cleanAddr:      "11.22.33.44:8080",
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%+v", tc)
		t.Run(name, func(t *testing.T) {
			cleanAddr, err := validateAddress(tc.addr, tc.allowLocalhost)
			require.Equal(t, tc.err, err)

			if err == nil {
				if tc.cleanAddr == "" {
					require.Equal(t, tc.addr, cleanAddr)
				} else {
					require.Equal(t, tc.cleanAddr, cleanAddr)
				}
			}
		})
	}
}

func TestNewPex(t *testing.T) {
	dir, err := ioutil.TempDir("", "peerlist")
	require.NoError(t, err)
	defer os.Remove(dir)

	config := NewConfig()
	config.DataDirectory = dir
	config.DefaultConnections = testPeers[:]

	// Add a peer and register it as a default peer.
	// It will be marked as "trusted"
	// in the peers cache file, but the next time it is loaded,
	// it will be reset to untrusted because it is not in the
	// cfg.DefaultConnections
	addr := "11.22.33.44:5566"
	config.DefaultConnections = append(config.DefaultConnections, addr)

	_, err = New(config)
	require.NoError(t, err)

	// check if peers are saved to disk
	peers, err := loadCachedPeersFile(filepath.Join(dir, PeerCacheFilename))
	require.NoError(t, err)

	require.Equal(t, len(testPeers)+1, len(peers))

	for _, p := range append(testPeers, addr) {
		v, ok := peers[p]
		require.True(t, ok)
		require.True(t, v.Trusted)
	}

	// Recreate pex with the extra peer removed form DefaultConnections
	config.DefaultConnections = config.DefaultConnections[:len(config.DefaultConnections)-1]
	_, err = New(config)
	require.NoError(t, err)

	peers, err = loadCachedPeersFile(filepath.Join(dir, PeerCacheFilename))
	require.NoError(t, err)

	require.Equal(t, len(testPeers)+1, len(peers))

	for _, p := range testPeers {
		v, ok := peers[p]
		require.True(t, ok)
		require.True(t, v.Trusted)
	}

	v, ok := peers[addr]
	require.True(t, ok)
	require.False(t, v.Trusted)
}

func TestNewPexDisableTrustedPeers(t *testing.T) {
	dir, err := ioutil.TempDir("", "peerlist")
	require.NoError(t, err)
	defer os.Remove(dir)

	config := NewConfig()
	config.DataDirectory = dir
	config.DefaultConnections = testPeers[:]
	config.DisableTrustedPeers = true

	_, err = New(config)
	require.NoError(t, err)

	// check if peers are saved to disk
	peers, err := loadCachedPeersFile(filepath.Join(dir, PeerCacheFilename))
	require.NoError(t, err)

	for _, p := range testPeers {
		v, ok := peers[p]
		require.True(t, ok)
		require.False(t, v.Trusted)
	}
}

func TestNewPexLoadCustomPeers(t *testing.T) {
	dir, err := ioutil.TempDir("", "peerlist")
	require.NoError(t, err)
	defer os.Remove(dir)

	fn, err := os.Create(filepath.Join(dir, "custom-peers.txt"))
	require.NoError(t, err)
	defer fn.Close()

	_, err = fn.Write([]byte(`123.45.67.89:2020
34.34.21.21:12222
`))
	require.NoError(t, err)

	err = fn.Close()
	require.NoError(t, err)

	config := NewConfig()
	config.DataDirectory = dir
	config.DefaultConnections = nil
	config.CustomPeersFile = fn.Name()

	_, err = New(config)
	require.NoError(t, err)

	// check if peers are saved to disk
	peers, err := loadCachedPeersFile(filepath.Join(dir, PeerCacheFilename))
	require.NoError(t, err)

	expectedPeers := []string{
		"123.45.67.89:2020",
		"34.34.21.21:12222",
	}

	for _, p := range expectedPeers {
		v, ok := peers[p]
		require.True(t, ok)
		require.False(t, v.Trusted)
	}

	require.Len(t, peers, len(expectedPeers))
}

func TestPexLoadPeers(t *testing.T) {
	tt := []struct {
		name        string
		filename    string
		peers       []Peer
		max         int
		expectN     int
		expectPeers []Peer
	}{
		{
			name:     "load all",
			filename: PeerCacheFilename,
			peers: []Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			max:     2,
			expectN: 2,
			expectPeers: []Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
		},
		{
			name:     "reach max",
			filename: PeerCacheFilename,
			peers: []Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
				Peer{Addr: testPeers[2]},
			},
			max:     2,
			expectN: 2,
			expectPeers: []Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
				Peer{Addr: testPeers[2]},
			},
		},
		{
			name:     "including invalid addr",
			filename: PeerCacheFilename,
			peers: []Peer{
				Peer{Addr: wrongPortPeer},
				Peer{Addr: testPeers[1]},
				Peer{Addr: testPeers[2]},
			},
			max:     2,
			expectN: 2,
			expectPeers: []Peer{
				Peer{Addr: testPeers[1]},
				Peer{Addr: testPeers[2]},
			},
		},
		{
			name:     "load all, fallback on oldPeerCacheFilename",
			filename: oldPeerCacheFilename,
			peers: []Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			max:     2,
			expectN: 2,
			expectPeers: []Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
		},
		{
			name:     "no peers file",
			filename: "foo.json",
			peers: []Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			max:         2,
			expectN:     0,
			expectPeers: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "peerlist")
			require.NoError(t, err)
			defer os.Remove(dir)

			// write peers to file
			fn := filepath.Join(dir, tc.filename)
			defer os.Remove(fn)

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
				peerlist: newPeerlist(),
				Config:   cfg,
			}

			err = px.loadCache()
			require.NoError(t, err)

			require.Len(t, px.peerlist.peers, tc.expectN)

			psm := make(map[string]Peer)
			for _, p := range tc.expectPeers {
				psm[p.Addr] = p
			}

			for _, p := range px.peerlist.peers {
				v, ok := psm[p.Addr]
				require.True(t, ok)
				require.Equal(t, v, *p)
			}
		})
	}
}

func TestPexAddPeer(t *testing.T) {
	now := time.Now().UTC().Unix()

	tt := []struct {
		name        string
		peers       []Peer
		max         int
		finalLen    int
		peer        string
		err         error
		removedPeer string
		check       func(px *Pex)
	}{
		{
			name: "ok",
			peers: []Peer{
				{
					Addr:     testPeers[0],
					LastSeen: now,
				},
			},
			max:      2,
			peer:     testPeers[1],
			err:      nil,
			finalLen: 2,
		},
		{
			name: "invalid peer",
			peers: []Peer{
				{
					Addr:     testPeers[0],
					LastSeen: now,
				},
			},
			max:      2,
			peer:     wrongPortPeer,
			err:      ErrInvalidAddress,
			finalLen: 1,
		},
		{
			name: "peer known",
			peers: []Peer{
				{
					Addr:     testPeers[0],
					LastSeen: now,
				},
				{
					Addr:     testPeers[1],
					LastSeen: now - 60,
				},
			},
			max:      2,
			peer:     testPeers[1],
			err:      nil,
			finalLen: 2,
			check: func(px *Pex) {
				p := px.peerlist.peers[testPeers[1]]
				require.NotNil(t, p)
				// Peer should have been marked as seen
				require.True(t, p.LastSeen > now-60)
			},
		},
		{
			name: "reach max",
			peers: []Peer{
				{
					Addr:     testPeers[0],
					LastSeen: now,
				},
				{
					Addr:     testPeers[1],
					LastSeen: now,
				},
			},
			max:      2,
			peer:     testPeers[3],
			err:      ErrPeerlistFull,
			finalLen: 2,
		},
		{
			name: "reach max but kicked old peer",
			peers: []Peer{
				{
					Addr:     testPeers[0],
					LastSeen: now,
				},
				{
					Addr:     testPeers[1],
					LastSeen: now - 60*60*24*2,
				},
			},
			max:         2,
			peer:        testPeers[3],
			err:         nil,
			finalLen:    2,
			removedPeer: testPeers[1],
		},
		{
			name: "no max",
			peers: []Peer{
				{
					Addr:     testPeers[0],
					LastSeen: now,
				},
				{
					Addr:     testPeers[1],
					LastSeen: now,
				},
			},
			max:      0,
			peer:     testPeers[3],
			err:      nil,
			finalLen: 3,
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
			cfg.DefaultConnections = []string{}

			// create px instance and load peers
			px, err := New(cfg)
			require.NoError(t, err)

			px.peerlist.setPeers(tc.peers)

			err = px.AddPeer(tc.peer)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.finalLen, len(px.peerlist.peers))

			if tc.check != nil {
				tc.check(px)
			}

			if err != nil {
				return
			}

			// check if the peer is in the peer list
			_, ok := px.peerlist.peers[tc.peer]
			require.True(t, ok)

			if tc.removedPeer != "" {
				_, ok := px.peerlist.peers[tc.removedPeer]
				require.False(t, ok)
			}
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
			cfg.DefaultConnections = tc.peers

			// create px instance and load peers
			px, err := New(cfg)
			require.NoError(t, err)

			n := px.AddPeers(tc.addPeers)
			require.Equal(t, tc.addN, n)

			for _, p := range tc.expectPeers {
				_, ok := px.peerlist.peers[p]
				require.True(t, ok)
			}
		})
	}
}

func TestPexTrustedPublic(t *testing.T) {
	tt := []struct {
		name   string
		peers  []Peer
		expect []Peer
	}{
		{

			"none peer",
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1], Trusted: true, Private: true},
			},
			[]Peer{},
		},
		{

			"one trusted public peer",
			[]Peer{
				Peer{Addr: testPeers[0], Trusted: true, Private: false},
				Peer{Addr: testPeers[1], Trusted: true, Private: true},
			},
			[]Peer{
				Peer{Addr: testPeers[0], Trusted: true, Private: false},
			},
		},
		{

			"all trust peer",
			[]Peer{
				Peer{Addr: testPeers[0], Trusted: true, Private: false},
				Peer{Addr: testPeers[1], Trusted: true, Private: false},
			},
			[]Peer{
				Peer{Addr: testPeers[0], Trusted: true, Private: false},
				Peer{Addr: testPeers[1], Trusted: true, Private: false},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			// get trusted public peers
			peers := pex.TrustedPublic()

			require.Equal(t, len(tc.expect), len(peers))
			pm := make(map[string]Peer)
			for _, p := range peers {
				pm[p.Addr] = p
			}

			for _, p := range tc.expect {
				v, ok := pm[p.Addr]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, p, v)
			}
		})
	}
}

func TestPexRandomExchangeable(t *testing.T) {
	tt := []struct {
		name        string
		peers       []Peer
		n           int
		expectN     int
		expectPeers []Peer
	}{
		{
			name: "n=0 exchangeable=0",
			peers: []Peer{
				Peer{Addr: testPeers[0], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			n:           0,
			expectN:     0,
			expectPeers: []Peer{},
		},
		{
			name: "n=0 exchangeable=1",
			peers: []Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			n:       0,
			expectN: 1,
			expectPeers: []Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
			},
		},
		{
			name: "n=0 exchangeable=2",
			peers: []Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			n:       0,
			expectN: 2,
			expectPeers: []Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true},
			},
		},
		{
			name: "n=1 exchangeable=0",
			peers: []Peer{
				Peer{Addr: testPeers[0], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			n:           1,
			expectN:     0,
			expectPeers: []Peer{},
		},
		{
			name: "n=1 exchangeable=1",
			peers: []Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: false},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			n:       1,
			expectN: 1,
			expectPeers: []Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
			},
		},
		{
			name: "n=1 exchangeable=2",
			peers: []Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			n:       1,
			expectN: 1,
			expectPeers: []Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			peers := pex.RandomExchangeable(tc.n)
			require.Len(t, peers, tc.expectN)

			// map expectPeers peers
			psm := make(map[string]Peer)
			for _, p := range tc.expectPeers {
				psm[p.Addr] = p
			}

			for _, p := range peers {
				v, ok := psm[p.Addr]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, p, v)
			}
		})
	}
}

func TestPeerRandomPublic(t *testing.T) {
	tt := []struct {
		name    string
		peers   []Peer
		n       int
		expectN int
	}{
		{
			"0 peer",
			[]Peer{},
			1,
			0,
		},
		{
			"1 peer",
			[]Peer{
				Peer{Addr: testPeers[0]},
			},
			1,
			1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			peers := pex.RandomPublic(tc.n)
			require.Len(t, peers, tc.expectN)
		})
	}
}

func TestPexRandomPublic(t *testing.T) {
	tt := []struct {
		name        string
		peers       []Peer
		n           int
		expectN     int
		expectPeers []Peer
	}{
		{
			"n=0 public=0",
			[]Peer{
				Peer{Addr: testPeers[0], Private: true},
				Peer{Addr: testPeers[1], Private: true},
				Peer{Addr: testPeers[2], Private: true},
			},
			0,
			0,
			[]Peer{},
		},
		{
			"n=0 public=2",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false},
				Peer{Addr: testPeers[1], Private: false},
				Peer{Addr: testPeers[2], Private: true},
			},
			0,
			2,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false},
				Peer{Addr: testPeers[1], Private: false},
			},
		},
		{
			"n=1 public=0",
			[]Peer{
				Peer{Addr: testPeers[0], Private: true},
				Peer{Addr: testPeers[1], Private: true},
				Peer{Addr: testPeers[2], Private: true},
			},
			1,
			0,
			[]Peer{},
		},
		{
			"n=1 public=2",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false},
				Peer{Addr: testPeers[1], Private: false},
				Peer{Addr: testPeers[2], Private: true},
			},
			1,
			1,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false},
				Peer{Addr: testPeers[1], Private: false},
			},
		},
		{
			"n=2 public=0",
			[]Peer{
				Peer{Addr: testPeers[0], Private: true},
				Peer{Addr: testPeers[1], Private: true},
				Peer{Addr: testPeers[2], Private: true},
			},
			2,
			0,
			[]Peer{},
		},
		{
			"n=2 public=1",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false},
				Peer{Addr: testPeers[1], Private: true},
				Peer{Addr: testPeers[2], Private: true},
			},
			2,
			1,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false},
			},
		},
		{
			"n=2 public=2",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false},
				Peer{Addr: testPeers[1], Private: false},
				Peer{Addr: testPeers[2], Private: true},
			},
			2,
			2,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false},
				Peer{Addr: testPeers[1], Private: false},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			// get N random public
			peers := pex.RandomPublic(tc.n)

			require.Len(t, peers, tc.expectN)

			// map the peers
			psm := make(map[string]Peer)
			for _, p := range tc.expectPeers {
				psm[p.Addr] = p
			}

			// check if the returned peers are in the expectPeers
			for _, p := range peers {
				v, ok := psm[p.Addr]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, p, v)
			}
		})
	}
}

func TestPexTrusted(t *testing.T) {
	tt := []struct {
		name   string
		peers  []Peer
		expect []Peer
	}{
		{

			"no trust peer",
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			[]Peer{},
		},
		{

			"one trust peer",
			[]Peer{
				Peer{Addr: testPeers[0], Trusted: true},
				Peer{Addr: testPeers[1]},
			},
			[]Peer{
				Peer{Addr: testPeers[0], Trusted: true},
			},
		},
		{

			"all trust peer",
			[]Peer{
				Peer{Addr: testPeers[0], Trusted: true},
				Peer{Addr: testPeers[1], Trusted: true},
			},
			[]Peer{
				Peer{Addr: testPeers[0], Trusted: true},
				Peer{Addr: testPeers[1], Trusted: true},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			peers := pex.Trusted()
			require.Equal(t, len(tc.expect), len(peers))

			pm := make(map[string]Peer)
			for _, p := range peers {
				pm[p.Addr] = p
			}

			for _, p := range tc.expect {
				v, ok := pm[p.Addr]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, p, v)
			}
		})
	}
}

func TestPexPrivate(t *testing.T) {
	tt := []struct {
		name   string
		peers  []Peer
		expect []Peer
	}{
		{

			"no private peer",
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			[]Peer{},
		},
		{

			"one private peer",
			[]Peer{
				Peer{Addr: testPeers[0], Private: true},
				Peer{Addr: testPeers[1]},
			},
			[]Peer{
				Peer{Addr: testPeers[0], Private: true},
			},
		},
		{

			"all trust peer",
			[]Peer{
				Peer{Addr: testPeers[0], Private: true},
				Peer{Addr: testPeers[1], Private: true},
			},
			[]Peer{
				Peer{Addr: testPeers[0], Private: true},
				Peer{Addr: testPeers[1], Private: true},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			peers := pex.Private()
			require.Equal(t, len(tc.expect), len(peers))

			pm := make(map[string]Peer)
			for _, p := range peers {
				pm[p.Addr] = p
			}

			for _, p := range tc.expect {
				v, ok := pm[p.Addr]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, p, v)
			}
		})
	}
}

func TestPexResetAllRetryTimes(t *testing.T) {
	tt := []struct {
		name   string
		peers  []Peer
		expect []Peer
	}{
		{
			"all",
			[]Peer{
				Peer{Addr: testPeers[0], RetryTimes: 1},
				Peer{Addr: testPeers[1], RetryTimes: 2},
			},
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			pex.ResetAllRetryTimes()

			for _, p := range tc.expect {
				v, ok := pex.peerlist.peers[p.Addr]
				require.True(t, ok)
				require.Equal(t, p, *v)
			}
		})
	}
}

func TestPexIncreaseRetryTimes(t *testing.T) {
	tt := []struct {
		name   string
		peers  []Peer
		addr   string
		expect map[string]Peer
	}{
		{
			"addr not exist",
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			testPeers[2],
			map[string]Peer{
				testPeers[0]: Peer{Addr: testPeers[0]},
				testPeers[1]: Peer{Addr: testPeers[1]},
			},
		},
		{
			"ok",
			[]Peer{
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			testPeers[0],
			map[string]Peer{
				testPeers[0]: Peer{Addr: testPeers[0], LastSeen: time.Now().UTC().Unix(), RetryTimes: 1},
				testPeers[1]: Peer{Addr: testPeers[1]},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			pex.IncreaseRetryTimes(tc.addr)

			require.Equal(t, len(tc.expect), len(pex.peerlist.peers))
			for k, v := range tc.expect {
				p, ok := pex.peerlist.peers[k]
				require.True(t, ok)
				if p.LastSeen != 0 {
					require.InDelta(t, v.LastSeen, p.LastSeen, 2)
					p.LastSeen = 0
					v.LastSeen = 0
				}
				require.Equal(t, v, *p)
			}
		})
	}
}

func TestPexResetRetryTimes(t *testing.T) {
	tt := []struct {
		name   string
		peers  []Peer
		addr   string
		expect []Peer
	}{
		{
			"no peer need reset",
			[]Peer{*NewPeer(testPeers[0]), *NewPeer(testPeers[1])},
			testPeers[2],
			[]Peer{*NewPeer(testPeers[0]), *NewPeer(testPeers[1])},
		},
		{
			"reset one",
			[]Peer{
				Peer{Addr: testPeers[0], LastSeen: time.Now().UTC().Unix(), RetryTimes: 10},
				Peer{Addr: testPeers[1], RetryTimes: 2},
			},
			testPeers[0],
			[]Peer{
				Peer{Addr: testPeers[0], LastSeen: time.Now().UTC().Unix()},
				Peer{Addr: testPeers[1], RetryTimes: 2},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.peers)

			pex.ResetRetryTimes(tc.addr)

			for _, p := range tc.expect {
				v, ok := pex.peerlist.peers[p.Addr]
				require.True(t, ok)
				require.Equal(t, p, *v)
			}
		})
	}
}

func TestPexRemovePeer(t *testing.T) {
	tt := []struct {
		name       string
		initPeers  []Peer
		removePeer string
		expect     map[string]*Peer
	}{
		{
			"remove from empty peer list",
			[]Peer{},
			testPeers[0],
			map[string]*Peer{},
		},
		{
			"remove one",
			[]Peer{*NewPeer(testPeers[0]), *NewPeer(testPeers[1])},
			testPeers[0],
			map[string]*Peer{
				testPeers[1]: NewPeer(testPeers[1]),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.initPeers)

			pex.RemovePeer(tc.removePeer)

			require.Equal(t, len(tc.expect), len(pex.peerlist.peers))
			for k, v := range tc.expect {
				p, ok := pex.peerlist.peers[k]
				require.True(t, ok)
				require.Equal(t, *v, *p)
			}
		})
	}
}

func TestPexSetPrivate(t *testing.T) {
	tt := []struct {
		name     string
		initPeer []Peer
		peer     string
		private  bool
		err      error
	}{
		{
			"set private true",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[0],
			true,
			nil,
		},
		{
			"set private false",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[0],
			false,
			nil,
		},
		{
			"set failed",
			[]Peer{*NewPeer(testPeers[1])},
			testPeers[0],
			false,
			fmt.Errorf("set peer.Private failed: %v does not exist in peer list", testPeers[0]),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.initPeer)

			err := pex.SetPrivate(tc.peer, tc.private)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			p, ok := pex.peerlist.peers[tc.peer]
			require.True(t, ok)

			require.Equal(t, tc.private, p.Private)
		})
	}
}

func TestPexSetTrusted(t *testing.T) {
	tt := []struct {
		name      string
		initPeers []Peer
		peer      string
		err       error
	}{
		{
			"set trust true",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[0],
			nil,
		},
		{
			"set failed",
			[]Peer{*NewPeer(testPeers[1])},
			testPeers[0],
			fmt.Errorf("set peer.Trusted failed: %v does not exist in peer list", testPeers[0]),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			// init peer
			pex.peerlist.setPeers(tc.initPeers)

			err := pex.setTrusted(tc.peer)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			p, ok := pex.peerlist.peers[tc.peer]
			require.True(t, ok)
			require.True(t, p.Trusted)
		})
	}
}

func TestPexSetHasIncomingPort(t *testing.T) {
	tt := []struct {
		name            string
		initPeers       []Peer
		peer            string
		hasIncomingPort bool
		err             error
	}{
		{
			"set has incoming port true",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[0],
			true,
			nil,
		},
		{
			"set has incoming port false",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[0],
			false,
			nil,
		},
		{
			"set failed",
			[]Peer{*NewPeer(testPeers[1])},
			testPeers[0],
			false,
			fmt.Errorf("set peer.HasIncomingPort failed: %v does not exist in peer list", testPeers[0]),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.initPeers)

			err := pex.SetHasIncomingPort(tc.peer, tc.hasIncomingPort)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			p, ok := pex.peerlist.peers[tc.peer]
			require.True(t, ok)
			require.Equal(t, tc.hasIncomingPort, p.HasIncomingPort)
		})
	}
}

func TestPexGetPeerByAddr(t *testing.T) {
	tt := []struct {
		name      string
		initPeers []Peer
		addr      string
		find      bool
		peer      Peer
	}{
		{
			"ok",
			[]Peer{
				*NewPeer(testPeers[0]),
				*NewPeer(testPeers[1]),
			},
			testPeers[0],
			true,
			*NewPeer(testPeers[0]),
		},
		{
			"not exist",
			[]Peer{
				*NewPeer(testPeers[0]),
				*NewPeer(testPeers[1]),
			},
			testPeers[2],
			false,
			Peer{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pex := &Pex{
				peerlist: newPeerlist(),
			}

			pex.peerlist.setPeers(tc.initPeers)

			p, ok := pex.GetPeer(tc.addr)
			require.Equal(t, tc.find, ok)
			if ok {
				require.Equal(t, tc.peer, p)
			}
		})
	}
}

func TestPexIsFull(t *testing.T) {
	pex := &Pex{
		peerlist: newPeerlist(),
		Config:   Config{Max: 0},
	}

	require.False(t, pex.IsFull())

	err := pex.AddPeer("11.22.33.44:5555")
	require.NoError(t, err)
	require.False(t, pex.IsFull())

	pex.Config.Max = 2
	require.False(t, pex.IsFull())
	err = pex.AddPeer("33.44.55.66:5555")
	require.NoError(t, err)
	require.True(t, pex.IsFull())

	pex.Config.Max = 1
	require.True(t, pex.IsFull())
}

func TestParseRemotePeerList(t *testing.T) {
	body := `11.22.33.44:5555
66.55.44.33:2020
# comment

127.0.0.1:8080
  54.54.32.32:7899
11.33.11.33
22.44.22.44:99
`

	peers := parseRemotePeerList(body)
	require.Len(t, peers, 3)
	require.Equal(t, []string{
		"11.22.33.44:5555",
		"66.55.44.33:2020",
		"54.54.32.32:7899",
	}, peers)
}

func TestParseLocalPeerList(t *testing.T) {
	cases := []struct {
		name           string
		body           string
		peers          []string
		allowLocalhost bool
		err            error
	}{
		{
			name: "valid, no localhost",
			body: `11.22.33.44:5555
66.55.44.33:2020
# comment

  54.54.32.32:7899
`,
			peers: []string{
				"11.22.33.44:5555",
				"66.55.44.33:2020",
				"54.54.32.32:7899",
			},
			allowLocalhost: false,
		},

		{
			name: "valid, localhost",
			body: `11.22.33.44:5555
66.55.44.33:2020
# comment

127.0.0.1:8080
  54.54.32.32:7899
`,
			peers: []string{
				"11.22.33.44:5555",
				"66.55.44.33:2020",
				"127.0.0.1:8080",
				"54.54.32.32:7899",
			},
			allowLocalhost: true,
		},

		{
			name: "invalid, contains localhost but no localhost allowed",
			body: `11.22.33.44:5555
66.55.44.33:2020
# comment

127.0.0.1:8080
  54.54.32.32:7899
`,
			err:            fmt.Errorf("Peers list has invalid address 127.0.0.1:8080: %v", ErrNoLocalhost),
			allowLocalhost: false,
		},

		{
			name: "invalid, bad address",
			body: `11.22.33.44:5555
66.55.44.33:2020
# comment

  54.54.32.32:99
`,
			err:            fmt.Errorf("Peers list has invalid address 54.54.32.32:99: %v", ErrPortTooLow),
			allowLocalhost: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			peers, err := parseLocalPeerList(tc.body, tc.allowLocalhost)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.peers, peers)
		})
	}
}
