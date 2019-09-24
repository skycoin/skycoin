package pex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/util/file"
)

var testPeers = []string{
	"112.32.32.14:7200",
	"112.32.32.15:7200",
	"112.32.32.16:7200",
	"112.32.32.17:7200",
}

var wrongPortPeer = "112.32.32.14:1"

/* Peer tests */

func TestNewPeer(t *testing.T) {
	p := NewPeer(testPeers[0])
	require.NotEqual(t, p.LastSeen, 0)
	require.Equal(t, p.Addr, testPeers[0])
	require.False(t, p.Private)
}

func TestPeerSeen(t *testing.T) {
	p := NewPeer(testPeers[0])
	x := p.LastSeen
	time.Sleep(time.Second)
	p.Seen()
	require.NotEqual(t, x, p.LastSeen)
	if p.LastSeen < x {
		t.Fail()
	}
}

func TestPeerString(t *testing.T) {
	p := NewPeer(testPeers[0])
	require.Equal(t, testPeers[0], p.String())
}

func TestLoadPeersFromFile(t *testing.T) {
	tt := []struct {
		name        string
		noFile      bool
		emptyFile   bool
		ps          []string
		expectPeers map[string]*Peer
		err         error
	}{
		{
			"no file",
			true,
			false,
			testPeers[0:0],
			map[string]*Peer{},
			nil,
		},
		{
			"one addr",
			false,
			false,
			testPeers[:1],
			map[string]*Peer{
				testPeers[0]: NewPeer(testPeers[0]),
			},
			nil,
		},
		{
			"two addr",
			false,
			false,
			testPeers[:2],
			map[string]*Peer{
				testPeers[0]: NewPeer(testPeers[0]),
				testPeers[1]: NewPeer(testPeers[1]),
			},
			nil,
		},
		{
			"empty peer list file",
			false,
			true,
			testPeers[0:0],
			map[string]*Peer{},
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			f, removeFile := preparePeerlistFile(t)
			if !tc.emptyFile {
				persistPeers(t, f, tc.ps)
			}

			if tc.noFile {
				// remove file immediately
				removeFile()
			} else {
				defer removeFile()
			}

			peers, err := loadCachedPeersFile(f)
			require.Equal(t, tc.err, err)
			require.Equal(t, len(tc.expectPeers), len(peers))
			for k, v := range tc.expectPeers {
				p, ok := peers[k]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, *v, *p)
			}
		})
	}
}

func TestPeerlistSetPeers(t *testing.T) {
	tt := []struct {
		name   string
		peers  []Peer
		expect map[string]Peer
	}{
		{
			"empty peers",
			[]Peer{},
			map[string]Peer{},
		},
		{
			"one peer",
			[]Peer{*NewPeer(testPeers[0])},
			map[string]Peer{
				testPeers[0]: *NewPeer(testPeers[0]),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pl := newPeerlist()
			pl.setPeers(tc.peers)
			require.Equal(t, len(tc.expect), len(pl.peers))
			for k, v := range tc.expect {
				p, ok := pl.peers[k]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, v, *p)
			}
		})
	}
}

func TestPeerlistAddPeer(t *testing.T) {
	tt := []struct {
		name        string
		initPeers   []Peer
		addPeer     string
		dup         bool
		expectPeers map[string]*Peer
	}{
		{
			"add peer to empty peer list",
			[]Peer{},
			testPeers[0],
			false,
			map[string]*Peer{
				testPeers[0]: NewPeer(testPeers[0]),
			},
		},
		{
			"add peer to none empty peer list",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[1],
			false,
			map[string]*Peer{
				testPeers[0]: NewPeer(testPeers[0]),
				testPeers[1]: NewPeer(testPeers[1]),
			},
		},
		{
			"add dup peer",
			[]Peer{{Addr: testPeers[0]}},
			testPeers[0],
			true,
			map[string]*Peer{
				testPeers[0]: NewPeer(testPeers[0]),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pl := newPeerlist()
			// init the peers
			pl.setPeers(tc.initPeers)

			if tc.dup {
				// sleep a second so that LastSeen is diff
				time.Sleep(time.Second)
			}

			// add peer
			pl.addPeer(tc.addPeer)

			require.Equal(t, len(tc.expectPeers), len(pl.peers))
			for k, v := range tc.expectPeers {
				p, ok := pl.peers[k]
				require.True(t, ok)
				if tc.dup {
					require.True(t, p.LastSeen > v.LastSeen)
					continue
				}

				peersEqualWithSeenAllowedDiff(t, *v, *p)
			}
		})
	}
}

func TestPeerlistAddPeers(t *testing.T) {
	tt := []struct {
		name        string
		initPeers   []Peer
		addPeers    []string
		expectPeers map[string]*Peer
	}{
		{
			"add one peer to empty list",
			[]Peer{},
			testPeers[:1],
			map[string]*Peer{
				testPeers[0]: NewPeer(testPeers[0]),
			},
		},
		{
			"add two peer to none empty list",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[1:3],
			map[string]*Peer{
				testPeers[0]: NewPeer(testPeers[0]),
				testPeers[1]: NewPeer(testPeers[1]),
				testPeers[2]: NewPeer(testPeers[2]),
			},
		},
		{
			"add dup peers",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[:3],
			map[string]*Peer{
				testPeers[0]: NewPeer(testPeers[0]),
				testPeers[1]: NewPeer(testPeers[1]),
				testPeers[2]: NewPeer(testPeers[2]),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pl := newPeerlist()
			// init peers
			pl.setPeers(tc.initPeers)

			// add peers
			pl.addPeers(tc.addPeers)

			require.Equal(t, len(tc.expectPeers), len(pl.peers))
			for k, v := range tc.expectPeers {
				p, ok := pl.peers[k]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, *v, *p)
			}
		})
	}
}

func TestPeerListSetTrusted(t *testing.T) {
	tt := []struct {
		name      string
		initPeers []Peer
		peer      string
		trust     bool
		err       error
	}{
		{
			"set trust true",
			[]Peer{*NewPeer(testPeers[0])},
			testPeers[0],
			true,
			nil,
		},
		{
			"set trust false",
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
			fmt.Errorf("set peer.Trusted failed: %v does not exist in peer list", testPeers[0]),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pl := newPeerlist()

			// init peer
			pl.setPeers(tc.initPeers)

			err := pl.setTrusted(tc.peer, tc.trust)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			p, ok := pl.peers[tc.peer]
			require.True(t, ok)
			require.Equal(t, tc.trust, p.Trusted)
		})
	}
}

func TestPeerlistFindOldestUntrustedPeer(t *testing.T) {
	peer1 := Peer{
		Addr:     "1.1.1.1:6060",
		LastSeen: time.Now().UTC().Unix() - 60*60*24*2,
	}
	peer2 := Peer{
		Addr:     "2.2.2.2:6060",
		LastSeen: time.Now().UTC().Unix() - 60*60*24*7,
	}
	peer3 := Peer{
		Addr:     "3.3.3.3:6060",
		LastSeen: time.Now().UTC().Unix() - 60,
	}
	trustedPeer := Peer{
		Addr:     "4.4.4.4:6060",
		LastSeen: time.Now().UTC().Unix() - 60*60*24*30,
		Trusted:  true,
	}
	privatePeer := Peer{
		Addr:     "5.5.5.5:6060",
		LastSeen: time.Now().UTC().Unix() - 60*60*24*30,
		Private:  true,
	}

	cases := []struct {
		name      string
		initPeers []Peer
		expect    *Peer
	}{
		{
			name:      "empty peerlist",
			initPeers: []Peer{},
			expect:    nil,
		},

		{
			name: "no untrusted public peers",
			initPeers: []Peer{
				trustedPeer,
				privatePeer,
			},
			expect: nil,
		},

		{
			name: "one peer",
			initPeers: []Peer{
				peer1,
			},
			expect: &peer1,
		},

		{
			name: "3 peers ignore trusted",
			initPeers: []Peer{
				peer1,
				trustedPeer,
				peer2,
				peer3,
			},
			expect: &peer2,
		},

		{
			name: "3 peers ignore private",
			initPeers: []Peer{
				peer1,
				privatePeer,
				peer2,
				peer3,
			},
			expect: &peer2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pl := newPeerlist()
			pl.setPeers(tc.initPeers)

			p := pl.findOldestUntrustedPeer()
			require.Equal(t, tc.expect, p)
		})
	}
}

func TestPeerlistClearOld(t *testing.T) {
	tt := []struct {
		name        string
		initPeers   []Peer
		timeAgo     time.Duration
		expectPeers map[string]Peer
	}{
		{
			"no old peers",
			[]Peer{
				{Addr: testPeers[0], LastSeen: time.Now().UTC().Unix() - 100},
			},
			110 * time.Second,
			map[string]Peer{
				testPeers[0]: {Addr: testPeers[0], LastSeen: time.Now().UTC().Unix() - 100},
			},
		},
		{
			"clear one old peer",
			[]Peer{
				{Addr: testPeers[0], LastSeen: time.Now().UTC().Unix() - 100},
				{Addr: testPeers[1], LastSeen: time.Now().UTC().Unix() - 110},
				{Addr: testPeers[2], LastSeen: time.Now().UTC().Unix() - 120},
			},
			111 * time.Second,
			map[string]Peer{
				testPeers[0]: {Addr: testPeers[0], LastSeen: time.Now().UTC().Unix() - 100},
				testPeers[1]: {Addr: testPeers[1], LastSeen: time.Now().UTC().Unix() - 110},
			},
		},
		{
			"clear two old peers",
			[]Peer{
				{Addr: testPeers[0], LastSeen: time.Now().UTC().Unix() - 100},
				{Addr: testPeers[1], LastSeen: time.Now().UTC().Unix() - 110},
				{Addr: testPeers[2], LastSeen: time.Now().UTC().Unix() - 120},
			},
			101 * time.Second,
			map[string]Peer{
				testPeers[0]: {Addr: testPeers[0], LastSeen: time.Now().UTC().Unix() - 100},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pl := newPeerlist()
			pl.setPeers(tc.initPeers)

			pl.clearOld(tc.timeAgo)
			require.Equal(t, len(pl.peers), len(tc.expectPeers))
			for _, p := range tc.expectPeers {
				v, ok := pl.peers[p.Addr]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, *v, p)
			}
		})
	}
}

func TestPeerlistSave(t *testing.T) {
	tt := []struct {
		name   string
		peers  []Peer
		expect map[string]Peer
	}{
		{
			"save all",
			[]Peer{
				{Addr: testPeers[0]},
				{Addr: testPeers[1]},
			},
			map[string]Peer{
				testPeers[0]: {Addr: testPeers[0]},
				testPeers[1]: {Addr: testPeers[1]},
			},
		},
		{
			"save one peer",
			[]Peer{
				{Addr: testPeers[0], RetryTimes: MaxPeerRetryTimes + 1},
				{Addr: testPeers[1]},
			},
			map[string]Peer{
				testPeers[1]: {Addr: testPeers[1]},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pl := newPeerlist()
			pl.setPeers(tc.peers)

			f, removeFile := preparePeerlistFile(t)
			defer removeFile()
			require.NoError(t, pl.save(f))

			psMap, err := loadCachedPeersFile(f)
			require.NoError(t, err)
			for k, v := range tc.expect {
				p, ok := psMap[k]
				require.True(t, ok)
				peersEqualWithSeenAllowedDiff(t, v, *p)
			}
		})
	}
}

func TestPeerCanTry(t *testing.T) {
	testData := []struct {
		LastSeen   int64
		RetryTimes int
		CanTry     bool
	}{
		{
			time.Now().UTC().Add(time.Duration(100) * time.Second * -1).Unix(),
			1,
			true,
		},
	}

	for _, d := range testData {
		p := Peer{
			LastSeen:   d.LastSeen,
			RetryTimes: d.RetryTimes,
		}
		require.Equal(t, d.CanTry, p.CanTry())
	}
}

func TestPeerJSONParsing(t *testing.T) {
	// The serialized peer json format changed,
	// this tests that the old format can still parse.
	oldFormat := `{
        "Addr": "11.22.33.44:6000",
        "LastSeen": "2017-09-24T06:42:18.999999999Z",
        "Private": true,
        "Trusted": true,
        "HasIncomePort": true
    }`

	newFormat := `{
        "Addr": "11.22.33.44:6000",
        "LastSeen": 1506235338,
        "Private": true,
        "Trusted": true,
        "HasIncomingPort": true
    }`

	check := func(p *Peer) {
		require.Equal(t, "11.22.33.44:6000", p.Addr)
		require.True(t, p.Private)
		require.True(t, p.Trusted)
		require.True(t, p.HasIncomingPort)
		require.Equal(t, int64(1506235338), p.LastSeen)
	}

	load := func(s string) PeerJSON {
		var pj PeerJSON
		dec := json.NewDecoder(strings.NewReader(s))
		dec.UseNumber()
		err := dec.Decode(&pj)
		require.NoError(t, err)
		return pj
	}

	pj := load(oldFormat)
	p, err := newPeerFromJSON(pj)
	require.NoError(t, err)
	check(p)

	pj = load(newFormat)
	p, err = newPeerFromJSON(pj)
	require.NoError(t, err)
	check(p)
}

func peersEqualWithSeenAllowedDiff(t *testing.T, expected Peer, actual Peer) {
	require.WithinDuration(t, time.Unix(expected.LastSeen, 0), time.Unix(actual.LastSeen, 0), 1*time.Second)
	expected.LastSeen = actual.LastSeen
	require.Equal(t, expected, actual)
}

// preparePeerlistFile makes peers.json in temporary dir,
func preparePeerlistFile(t *testing.T) (string, func()) {
	f, err := ioutil.TempFile("", PeerCacheFilename)
	require.NoError(t, err)

	return f.Name(), func() {
		os.Remove(f.Name())
	}
}

func preparePeerlistDir(t *testing.T) (string, func()) {
	f, err := ioutil.TempDir("", "peerlist")
	require.NoError(t, err)

	return f, func() {
		os.Remove(f)
	}
}

func persistPeers(t *testing.T, fn string, peers []string) {
	t.Helper()
	peersMap := make(map[string]*Peer, len(peers))
	for _, p := range peers {
		peersMap[p] = NewPeer(p)
	}

	if err := file.SaveJSON(fn, peersMap, 0600); err != nil {
		panic(fmt.Sprintf("save peer list failed: %v", err))
	}
}
