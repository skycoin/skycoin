package pex

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/stretchr/testify/assert"
)

var testPeers = []string{
	"112.32.32.14:7200",
	"112.32.32.15:7200",
	"112.32.32.16:7200",
	"112.32.32.17:7200",
}

var wrongPortPeer = "112.32.32.14:1"

func init() {
	// silence the logger
	logging.Disable()
}

/* Peer tests */

func TestNewPeer(t *testing.T) {
	p := NewPeer(testPeers[0])
	assert.NotEqual(t, p.LastSeen, 0)
	assert.Equal(t, p.Addr, testPeers[0])
	assert.False(t, p.Private)
}

func TestPeerSeen(t *testing.T) {
	p := NewPeer(testPeers[0])
	x := p.LastSeen
	time.Sleep(time.Second)
	p.Seen()
	assert.NotEqual(t, x, p.LastSeen)
	if p.LastSeen < x {
		t.Fail()
	}
}

func TestPeerString(t *testing.T) {
	p := NewPeer(testPeers[0])
	assert.Equal(t, testPeers[0], p.String())
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
			io.EOF,
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

			peers, err := loadPeersFromFile(f)
			require.Equal(t, tc.err, err)
			require.Equal(t, len(tc.expectPeers), len(peers))
			for k, v := range tc.expectPeers {
				p, ok := peers[k]
				require.True(t, ok)
				require.Equal(t, *v, *p)
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
				require.Equal(t, v, *p)
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
			[]Peer{Peer{Addr: testPeers[0]}},
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

				require.Equal(t, *v, *p)
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
				require.Equal(t, *v, *p)
			}
		})
	}
}

func TestPeerlistRemovePeer(t *testing.T) {
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
			pl := newPeerlist()
			// init peers
			pl.setPeers(tc.initPeers)

			pl.RemovePeer(tc.removePeer)

			require.Equal(t, len(tc.expect), len(pl.peers))
			for k, v := range tc.expect {
				p, ok := pl.peers[k]
				require.True(t, ok)
				require.Equal(t, *v, *p)
			}
		})
	}
}

func TestPeerlistSetPrivate(t *testing.T) {
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
			pl := newPeerlist()
			// init peer list
			pl.setPeers(tc.initPeer)

			err := pl.setPrivate(tc.peer, tc.private)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			p, ok := pl.peers[tc.peer]
			require.True(t, ok)

			require.Equal(t, tc.private, p.Private)
		})
	}
}

func TestPeerlistSetTrust(t *testing.T) {
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

func TestPeerlistSetHasIncomingPort(t *testing.T) {
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
			pl := newPeerlist()

			// init peer
			pl.setPeers(tc.initPeers)

			err := pl.setPeerHasIncomingPort(tc.peer, tc.hasIncomingPort)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			p, ok := pl.peers[tc.peer]
			require.True(t, ok)
			require.Equal(t, tc.hasIncomingPort, p.HasIncomingPort)
		})
	}
}

func TestPeerlistCullInvalidPeers(t *testing.T) {
	tt := []struct {
		name        string
		initPeers   Peers
		culledPeers map[string]Peer
	}{
		{
			"no invalid peer",
			Peers{
				Peer{Addr: testPeers[0], Trusted: true},
				Peer{Addr: testPeers[1], Trusted: true},
			},
			map[string]Peer{},
		},
		{
			"cull invalid peer",
			Peers{
				Peer{Addr: testPeers[0], Trusted: false, RetryTimes: maxRetryTimes + 1},
				Peer{Addr: testPeers[1], Trusted: true},
			},
			map[string]Peer{
				testPeers[0]: Peer{Addr: testPeers[0], Trusted: false, RetryTimes: maxRetryTimes + 1},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pl := newPeerlist()
			pl.setPeers(tc.initPeers)

			invalidPeers := pl.cullInvalidPeers()
			require.Equal(t, len(tc.culledPeers), len(invalidPeers))
			for _, p := range invalidPeers {
				v, ok := tc.culledPeers[p.Addr]
				require.True(t, ok)
				require.Equal(t, v, p)
			}
		})
	}
}

func TestPeerlistGetPeerByAddr(t *testing.T) {
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
			pl := newPeerlist()
			// init peer
			pl.setPeers(tc.initPeers)

			p, ok := pl.GetPeerByAddr(tc.addr)
			require.Equal(t, tc.find, ok)
			if ok {
				require.Equal(t, tc.peer, p)
			}
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
				Peer{Addr: testPeers[0], LastSeen: utc.UnixNow() - 100},
			},
			110 * time.Second,
			map[string]Peer{
				testPeers[0]: Peer{Addr: testPeers[0], LastSeen: utc.UnixNow() - 100},
			},
		},
		{
			"clear one old peer",
			[]Peer{
				Peer{Addr: testPeers[0], LastSeen: utc.UnixNow() - 100},
				Peer{Addr: testPeers[1], LastSeen: utc.UnixNow() - 110},
				Peer{Addr: testPeers[2], LastSeen: utc.UnixNow() - 120},
			},
			111 * time.Second,
			map[string]Peer{
				testPeers[0]: Peer{Addr: testPeers[0], LastSeen: utc.UnixNow() - 100},
				testPeers[1]: Peer{Addr: testPeers[1], LastSeen: utc.UnixNow() - 110},
			},
		},
		{
			"clear two old peers",
			[]Peer{
				Peer{Addr: testPeers[0], LastSeen: utc.UnixNow() - 100},
				Peer{Addr: testPeers[1], LastSeen: utc.UnixNow() - 110},
				Peer{Addr: testPeers[2], LastSeen: utc.UnixNow() - 120},
			},
			101 * time.Second,
			map[string]Peer{
				testPeers[0]: Peer{Addr: testPeers[0], LastSeen: utc.UnixNow() - 100},
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
				require.Equal(t, *v, p)
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
				Peer{Addr: testPeers[0]},
				Peer{Addr: testPeers[1]},
			},
			map[string]Peer{
				testPeers[0]: Peer{Addr: testPeers[0]},
				testPeers[1]: Peer{Addr: testPeers[1]},
			},
		},
		{
			"save one peer",
			[]Peer{
				Peer{Addr: testPeers[0], RetryTimes: maxRetryTimes + 1},
				Peer{Addr: testPeers[1]},
			},
			map[string]Peer{
				testPeers[1]: Peer{Addr: testPeers[1]},
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

			psMap, err := loadPeersFromFile(f)
			require.NoError(t, err)
			for k, v := range tc.expect {
				p, ok := psMap[k]
				require.True(t, ok)
				require.Equal(t, v, *p)
			}
		})
	}
}

func TestPeerlistIncreaseRetryTimes(t *testing.T) {
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
			pl := newPeerlist()
			pl.setPeers(tc.peers)

			pl.IncreaseRetryTimes(tc.addr)

			require.Equal(t, len(tc.expect), len(pl.peers))
			for k, v := range tc.expect {
				p, ok := pl.peers[k]
				require.True(t, ok)
				require.Equal(t, v, *p)
			}
		})
	}
}

func TestPeerlistResetRetryTimes(t *testing.T) {
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
			pl := newPeerlist()
			pl.setPeers(tc.peers)

			pl.ResetRetryTimes(tc.addr)

			for _, p := range tc.expect {
				v, ok := pl.peers[p.Addr]
				require.True(t, ok)
				require.Equal(t, p, *v)
			}
		})
	}
}

func TestPeerlistResetAllRetryTimes(t *testing.T) {
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
			pl := newPeerlist()
			pl.setPeers(tc.peers)
			pl.ResetAllRetryTimes()

			for _, p := range tc.expect {
				v, ok := pl.peers[p.Addr]
				require.True(t, ok)
				require.Equal(t, p, *v)
			}
		})
	}
}

func TestGetPeerlistTrust(t *testing.T) {
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
			pl := newPeerlist()
			pl.setPeers(tc.peers)
			peers := pl.Trusted()
			require.Equal(t, len(tc.expect), len(peers))
			pm := make(map[string]Peer)
			for _, p := range peers {
				pm[p.Addr] = p
			}

			for _, p := range tc.expect {
				v, ok := pm[p.Addr]
				require.True(t, ok)
				require.Equal(t, p, v)
			}
		})
	}
}

func TestPeerlistPrivate(t *testing.T) {
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
			pl := newPeerlist()
			pl.setPeers(tc.peers)
			peers := pl.Private()
			require.Equal(t, len(tc.expect), len(peers))
			pm := make(map[string]Peer)
			for _, p := range peers {
				pm[p.Addr] = p
			}

			for _, p := range tc.expect {
				v, ok := pm[p.Addr]
				require.True(t, ok)
				require.Equal(t, p, v)
			}
		})
	}
}

func TestPeerlistTrustPublic(t *testing.T) {
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
			pl := newPeerlist()
			// init peers
			pl.setPeers(tc.peers)

			// get trusted public peers
			peers := pl.TrustedPublic()

			require.Equal(t, len(tc.expect), len(peers))
			pm := make(map[string]Peer)
			for _, p := range peers {
				pm[p.Addr] = p
			}

			for _, p := range tc.expect {
				v, ok := pm[p.Addr]
				require.True(t, ok)
				require.Equal(t, p, v)
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
			pl := newPeerlist()
			// init peers
			pl.setPeers(tc.peers)

			peers := pl.RandomPublic(tc.n)
			require.Len(t, peers, tc.expectN)
		})
	}
}

func TestPeerlistRandomPublic(t *testing.T) {
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
			pl := newPeerlist()
			pl.setPeers(tc.peers)

			// get N random public
			peers := pl.RandomPublic(tc.n)

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
				require.Equal(t, p, v)
			}
		})
	}
}

func TestPeerlistRandomExchangeable(t *testing.T) {
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
			pl := newPeerlist()
			pl.setPeers(tc.peers)

			peers := pl.RandomExchangeable(tc.n)
			require.Len(t, peers, tc.expectN)

			// map expectIN peers
			psm := make(map[string]Peer)
			for _, p := range tc.expectIN {
				psm[p.Addr] = p
			}

			for _, p := range peers {
				v, ok := psm[p.Addr]
				require.True(t, ok)
				require.Equal(t, p, v)
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
			utc.Now().Add(time.Duration(100) * time.Second * -1).Unix(),
			1,
			true,
		},
	}

	for _, d := range testData {
		p := Peer{
			LastSeen:   d.LastSeen,
			RetryTimes: d.RetryTimes,
		}
		assert.Equal(t, d.CanTry, p.CanTry())
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

	check := func(p Peer) {
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

// preparePeerlistFile makes peers.txt in temporary dir,
func preparePeerlistFile(t *testing.T) (string, func()) {
	f, err := ioutil.TempFile("", "peers.txt")
	require.NoError(t, err)

	return f.Name(), func() {
		os.Remove(f.Name())
	}
}

func preparePeerlistDir(t *testing.T) (string, func()) {
	f, err := ioutil.TempDir("", "peerlist")
	if err != nil {
		panic(err)
	}

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
