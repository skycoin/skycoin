package pex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/utc"
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

	// // empty string
	// require.False(t, validateAddress("", false))
	// // doubled ip:port
	// require.False(t, validateAddress("112.32.32.14:100112.32.32.14:101", false))
	// // requires port
	// require.False(t, validateAddress("112.32.32.14", false))
	// // not ip
	// require.False(t, validateAddress("112", false))
	// require.False(t, validateAddress("112.32", false))
	// require.False(t, validateAddress("112.32.32", false))
	// // bad part
	// require.False(t, validateAddress("112.32.32.14000", false))
	// // large port
	// require.False(t, validateAddress("112.32.32.14:66666", false))
	// // unspecified
	// require.False(t, validateAddress("0.0.0.0:8888", false))
	// // no ip
	// require.False(t, validateAddress(":8888", false))
	// // multicast
	// require.False(t, validateAddress("224.1.1.1:8888", false))
	// // invalid ports
	// require.False(t, validateAddress("112.32.32.14:0", false))
	// require.False(t, validateAddress("112.32.32.14:1", false))
	// require.False(t, validateAddress("112.32.32.14:10", false))
	// require.False(t, validateAddress("112.32.32.14:100", false))
	// require.False(t, validateAddress("112.32.32.14:1000", false))
	// require.False(t, validateAddress("112.32.32.14:1023", false))
	// require.False(t, validateAddress("112.32.32.14:65536", false))
	// // valid ones
	// require.True(t, validateAddress("112.32.32.14:1024", false))
	// require.True(t, validateAddress("112.32.32.14:10000", false))
	// require.True(t, validateAddress("112.32.32.14:65535", false))
	// // localhost is allowed
	// require.True(t, validateAddress("127.0.0.1:8888", true))
	// // localhost is not allowed
	// require.False(t, validateAddress("127.0.0.1:8888", false))
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
				peerlist: newPeerlist(),
				Config:   cfg,
			}

			err = px.load()
			require.NoError(t, err)

			require.Len(t, px.peerlist.peers, tc.expectN)

			psm := make(map[string]Peer)
			for _, p := range tc.expectIN {
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
		{
			"no max",
			testPeers[:2],
			0,
			testPeers[3],
			nil,
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
			_, ok := px.peerlist.peers[tc.peer]
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
		name     string
		peers    []Peer
		n        int
		expectN  int
		expectIN []Peer
	}{
		{
			"n=0 exchangeable=0",
			[]Peer{
				Peer{Addr: testPeers[0], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			0,
			0,
			[]Peer{},
		},
		{
			"n=0 exchangeable=1",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			0,
			1,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
			},
		},
		{
			"n=0 exchangeable=2",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			0,
			2,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true},
			},
		},
		{
			"n=1 exchangeable=0",
			[]Peer{
				Peer{Addr: testPeers[0], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: true, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			1,
			0,
			[]Peer{},
		},
		{
			"n=1 exchangeable=1",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: false},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			1,
			1,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
			},
		},
		{
			"n=1 exchangeable=2",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			1,
			1,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true},
			},
		},
		{
			"n=2 exchangeable=1",
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
				Peer{Addr: testPeers[1], Private: false, HasIncomingPort: true, RetryTimes: 1},
				Peer{Addr: testPeers[2], Private: true, HasIncomingPort: true},
			},
			2,
			1,
			[]Peer{
				Peer{Addr: testPeers[0], Private: false, HasIncomingPort: true},
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

			// map expectIN peers
			psm := make(map[string]Peer)
			for _, p := range tc.expectIN {
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
		name     string
		peers    []Peer
		n        int
		expectN  int
		expectIN []Peer
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
			for _, p := range tc.expectIN {
				psm[p.Addr] = p
			}

			// check if the returned peers are in the expectIN
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
				testPeers[0]: Peer{Addr: testPeers[0], LastSeen: utc.UnixNow(), RetryTimes: 1},
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
				Peer{Addr: testPeers[0], LastSeen: utc.UnixNow(), RetryTimes: 10},
				Peer{Addr: testPeers[1], RetryTimes: 2},
			},
			testPeers[0],
			[]Peer{
				Peer{Addr: testPeers[0], LastSeen: utc.UnixNow()},
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

			err := pex.SetTrusted(tc.peer)
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

			p, ok := pex.GetPeerByAddr(tc.addr)
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
