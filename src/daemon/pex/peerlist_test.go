package pex

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// var (
// 	address   = "112.32.32.14:3030"
// 	address2  = "112.32.32.14:3031"
// 	addresses = []string{
// 		address, "111.32.32.13:2020", "69.32.54.111:2222",
// 	}
// 	silenceLogger = true
// )

var peers = []string{
	"112.32.32.14:10011",
	"112.32.32.14:20011",
	"112.32.32.14:30011",
	"112.32.32.14:40011",
}

func init() {
	// silence the logger
	logging.Disable()
}

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

/* Peer tests */

func TestNewPeer(t *testing.T) {
	p := NewPeer(peers[0])
	assert.NotEqual(t, p.LastSeen, 0)
	assert.Equal(t, p.Addr, peers[0])
	assert.False(t, p.Private)
}

func TestPeerSeen(t *testing.T) {
	p := NewPeer(peers[0])
	x := p.LastSeen
	time.Sleep(time.Second)
	p.Seen()
	assert.NotEqual(t, x, p.LastSeen)
	if p.LastSeen < x {
		t.Fail()
	}
}

func TestPeerString(t *testing.T) {
	p := NewPeer(peers[0])
	assert.Equal(t, peers[0], p.String())
}

/* BlacklistEntry tests */

// func TestBlacklistEntryExpiresAt(t *testing.T) {
// 	now := utc.Now()
// 	b := BlacklistEntry{Start: now, Duration: time.Second}
// 	assert.Equal(t, now.Add(time.Second), b.ExpiresAt())
// }

/* Blacklist tests */

// func TestBlacklistSaveLoad(t *testing.T) {
// 	// Create and save a blacklist
// 	os.Remove("./" + BlacklistedDatabaseFilename)
// 	b := make(Blacklist)
// 	be := NewBlacklistEntry(time.Minute)
// 	b[address] = be
// 	b[""] = be
// 	b.Save(".")

// 	// Check that the file appears correct
// 	f, err := os.Open("./" + BlacklistedDatabaseFilename)
// 	assert.Nil(t, err)
// 	buf := make([]byte, 1024)
// 	reader := bufio.NewReader(f)
// 	n, err := reader.Read(buf)
// 	assert.Nil(t, err)
// 	buf = buf[:n]
// 	assert.Equal(t, string(buf[:len(address)]), address)
// 	assert.Equal(t, int8(buf[len(buf)-1]), '\n')
// 	f.Close()

// 	// Load the saved blacklist, check the contents match
// 	bb, err := LoadBlacklist(".")
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(bb), len(b)-1)
// 	for k, v := range bb {
// 		assert.Equal(t, v.Start.Unix(), b[k].Start.Unix())
// 		assert.Equal(t, v.Duration, b[k].Duration)
// 	}

// 	// Write a file with bad data
// 	f, err = os.Create("./" + BlacklistedDatabaseFilename)
// 	assert.Nil(t, err)
// 	garbage := []string{
// 		"", // empty line
// 		"#" + address + " 1000 1000", // commented line
// 		"notaddress 1000 1000",       // bad address
// 		address + " xxx 1000",        // bad start time
// 		address + " 1000 xxx",        // bad duration
// 		address + " 1000",            // not enough info
// 		// this one is good, but has extra spaces
// 		address + "  9999999999\t\t1000",
// 	}
// 	w := bufio.NewWriter(f)
// 	data := strings.Join(garbage, "\n") + "\n"
// 	n, err = w.Write([]byte(data))
// 	assert.Nil(t, err)
// 	w.Flush()
// 	f.Close()

// 	// Load the file with bad data and confirm they did not make it
// 	bb, err = LoadBlacklist(".")
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(bb), 1)
// 	assert.NotNil(t, bb[address])
// 	assert.Equal(t, bb[address].Duration, time.Duration(1000)*time.Second)
// }

// func TestBlacklistRefresh(t *testing.T) {
// 	b := make(Blacklist)
// 	be := NewBlacklistEntry(time.Microsecond)
// 	b[address] = be
// 	time.Sleep(time.Microsecond * 500)
// 	assert.Equal(t, len(b), 1)
// 	b.Refresh()
// 	assert.Equal(t, len(b), 0)
// }

// func TestBlacklistGetAddresses(t *testing.T) {
// 	b := make(Blacklist)
// 	for _, a := range addresses {
// 		b[a] = NewBlacklistEntry(time.Second)
// 	}
// 	expect := make([]string, len(addresses))
// 	for i, k := range addresses {
// 		expect[i] = k
// 	}
// 	sort.Strings(expect)
// 	keys := b.GetAddresses()
// 	sort.Strings(keys)
// 	assert.Equal(t, len(keys), len(expect))
// 	for i, v := range keys {
// 		assert.Equal(t, v, expect[i])
// 	}

// }

/* Pex tests */

// func TestNewPex(t *testing.T) {
// 	p := NewPex(10)
// 	assert.NotNil(t, p.Peerlist)
// 	assert.Equal(t, len(p.Peerlist), 0)
// 	assert.NotNil(t, p.Blacklist)
// 	assert.Equal(t, p.maxPeers, 10)
// }

// func TestAddBlacklistEntry(t *testing.T) {
// 	p := NewPex(10)
// 	p.AddPeer(address)
// 	assert.NotNil(t, p.Peerlist[address])
// 	_, exists := p.Blacklist[address]
// 	assert.False(t, exists)
// 	duration := time.Minute * 9
// 	p.AddBlacklistEntry(p.Peerlist[address].Addr, duration)
// 	assert.Nil(t, p.Peerlist[address])
// 	assert.Equal(t, p.Blacklist[address].Duration, duration)
// 	now := time.Now()
// 	assert.True(t, p.Blacklist[address].Start.Before(now))
// 	assert.True(t, p.Blacklist[address].Start.Add(duration).After(now))
// 	// Blacklisting invalid peer triggers logger -- just get the coverage
// 	p.AddBlacklistEntry("xxx", time.Second)
// 	_, exists = p.Blacklist["xxx"]
// 	assert.False(t, exists)
// 	// Blacklisting private peer is prevented
// 	q, err := p.AddPeer(address2)
// 	assert.Nil(t, err)
// 	q.Private = true
// 	p.AddBlacklistEntry(address2, time.Second)
// 	_, exists = p.Blacklist[address2]
// 	assert.False(t, exists)
// 	q = p.Peerlist[address2]
// 	assert.NotNil(t, q)
// 	assert.Equal(t, q.Addr, address2)
// }

// func TestAddPeers(t *testing.T) {
// 	p := NewPex(10)
// 	peers := make([]string, 4)
// 	peers[0] = "112.32.32.14:10011"
// 	peers[1] = "112.32.32.14:20011"
// 	peers[2] = "xxx"
// 	peers[3] = "127.0.0.1:10444"
// 	n := p.AddPeers(peers)
// 	assert.Equal(t, n, 2)
// 	assert.NotNil(t, p.Peerlist[peers[0]])
// 	assert.NotNil(t, p.Peerlist[peers[1]])
// 	assert.Nil(t, p.Peerlist[peers[2]])
// }

// func TestClearOld(t *testing.T) {
// 	p := NewPex(10)
// 	p.AddPeer("112.32.32.14:10011")
// 	q, _ := p.AddPeer("112.32.32.14:20011")
// 	assert.Equal(t, len(p.Peerlist), 2)
// 	p.Peerlist.ClearOld(time.Second * 100)
// 	assert.Equal(t, len(p.Peerlist), 2)
// 	q.LastSeen = q.LastSeen.Add(time.Second * -200)
// 	p.Peerlist.ClearOld(time.Second * 100)
// 	assert.Equal(t, len(p.Peerlist), 1)
// 	assert.Nil(t, p.Peerlist["112.32.32.14:20011"])
// 	assert.NotNil(t, p.Peerlist["112.32.32.14:10011"])

// 	// Should ignore a private peer
// 	assert.Equal(t, len(p.Peerlist), 1)
// 	q, err := p.AddPeer("112.32.32.14:20011")
// 	assert.Nil(t, err)
// 	q.Private = true
// 	assert.Equal(t, len(p.Peerlist), 2)
// 	q.LastSeen = q.LastSeen.Add(time.Second * -200)
// 	p.Peerlist.ClearOld(time.Second * 100)
// 	// Private peer should not be removed
// 	assert.Equal(t, len(p.Peerlist), 2)
// 	assert.NotNil(t, p.Peerlist["112.32.32.14:10011"])
// 	assert.NotNil(t, p.Peerlist["112.32.32.14:20011"])
// }

// func TestGetPublicAddresses(t *testing.T) {
// 	pex, err := NewPex(Config{
// 		MaxPeers: 10,
// 	})
// 	require.NoError(t, err)

// 	pex.AddPeer("112.32.32.14:10011")
// 	pex.AddPeer("112.32.32.14:20011")
// 	pex.AddPeer("112.32.32.14:30011")
// 	pex.SetPrivate("112.32.32.14:30011", true)

// 	addresses := pex.GetPublicAddresses()
// 	assert.Equal(t, len(addresses), 2)
// 	sort.Strings(addresses)
// 	assert.Equal(t, addresses, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	})
// }

func TestPeerlistAddPeer(t *testing.T) {
	l := len(peers)
	pl := newPeerlist(l - 1)
	for i := 0; i < l-1; i++ {
		err := pl.addPeer(peers[i])
		require.NoError(t, err)
	}

	require.Len(t, pl.peers, 3)

	// test peer list full
	require.Error(t, ErrPeerlistFull, pl.addPeer(peers[3]))

	// test dup peer
	require.NoError(t, pl.addPeer(peers[0]))
}

func TestPeerlistAddPeers(t *testing.T) {
	l := len(peers)
	pl := newPeerlist(l)
	verifyF := func(addr string) error {
		if !validateAddress(addr, false) {
			return ErrInvalidAddress
		}
		return nil
	}
	n := pl.addPeers(peers, verifyF)
	ps := append(peers, "localhost:11001")

	n = pl.addPeers(ps, verifyF)
	require.Equal(t, l, n)

	// check peer list full
	ps = append(peers, "112.32.32.14:50011")
	n = pl.addPeers(ps, verifyF)
	require.Equal(t, l, n)
}

func TestPeerlistRemovePeer(t *testing.T) {
	l := len(peers)
	pl := newPeerlist(l)
	for _, p := range peers {
		require.NoError(t, pl.addPeer(p))
	}

	require.Len(t, pl.peers, l)
	pl.RemovePeer(peers[0])
	require.Len(t, pl.peers, l-1)
}

func TestPeerlistGetPublicTrustPeers(t *testing.T) {
	l := len(peers)
	pl := newPeerlist(l)
	for _, p := range peers {
		require.NoError(t, pl.addPeer(p))
	}

}

func TestGetPrivateAddresses(t *testing.T) {
	ips := []string{
		"112.32.32.14:10011",
		"112.32.32.14:20011",
		"112.32.32.14:30011",
		"112.32.32.14:40011",
	}
	pl := newPeerlist(10)
	pl.addPeer(ips[0])
	pl.addPeer(ips[1])
	pl.addPeer(ips[2])
	pl.addPeer(ips[3])
	pl.setPrivate(ips[2], true)
	pl.setPrivate(ips[3], true)

	addresses := pl.GetPrivateAddresses()
	assert.Equal(t, len(addresses), 2)
	sort.Strings(addresses)
	assert.Equal(t, addresses, []string{
		ips[2],
		ips[3],
	})
}

func convertPeersToStrings(peers []Peer) []string {
	addresses := make([]string, 0, len(peers))
	for _, p := range peers {
		addresses = append(addresses, p.String())
	}
	return addresses
}

func compareRandom(t *testing.T, p *Pex, npeers int, result []string, f func(int) []Peer) {
	peers := f(npeers)
	addresses := convertPeersToStrings(peers)
	sort.Strings(addresses)
	assert.Equal(t, addresses, result)
}

// func testRandom(t *testing.T, publicOnly bool) {
// 	p, err := NewPex(Config{
// 		MaxPeers: 10,
// 	})
// 	require.NoError(t, err)

// 	f := p.RandomAll
// 	if publicOnly {
// 		f = p.RandomPublic
// 	}

// 	// check without peers
// 	assert.NotNil(t, p.RandomAll(100))
// 	assert.Equal(t, len(p.RandomAll(100)), 0)

// 	// check with one peer
// 	p.AddPeer("112.32.32.14:10011")
// 	// 0 defaults to all peers
// 	compareRandom(t, p, 0, []string{"112.32.32.14:10011"}, f)
// 	compareRandom(t, p, 1, []string{"112.32.32.14:10011"}, f)
// 	// exceeding known peers is safe
// 	compareRandom(t, p, 2, []string{"112.32.32.14:10011"}, f)
// 	// exceeding max peers is safe
// 	compareRandom(t, p, 100, []string{"112.32.32.14:10011"}, f)

// 	// check with two peers
// 	p.AddPeer("112.32.32.14:20011")
// 	one := p.RandomAll(1)[0].String()
// 	if one != "112.32.32.14:10011" && one != "112.32.32.14:20011" {
// 		assert.Nil(t, nil)
// 	}
// 	// 0 defaults to all peers
// 	compareRandom(t, p, 0, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	}, f)
// 	compareRandom(t, p, 2, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	}, f)
// 	compareRandom(t, p, 3, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	}, f)
// 	compareRandom(t, p, 100, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	}, f)

// 	// check with 3 peers, one private
// 	p.AddPeer("112.32.32.14:30011")
// 	p.SetPrivate("112.32.32.14:30011", true)
// 	if publicOnly {
// 		// The private peer should never be included
// 		compareRandom(t, p, 0, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 		}, f)
// 		compareRandom(t, p, 2, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 		}, f)
// 		compareRandom(t, p, 3, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 		}, f)
// 		compareRandom(t, p, 100, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 		}, f)
// 	} else {
// 		// The private peer should be included
// 		compareRandom(t, p, 0, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 			"112.32.32.14:30011",
// 		}, f)
// 		compareRandom(t, p, 3, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 			"112.32.32.14:30011",
// 		}, f)
// 		compareRandom(t, p, 4, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 			"112.32.32.14:30011",
// 		}, f)
// 		compareRandom(t, p, 100, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 			"112.32.32.14:30011",
// 		}, f)
// 	}
// }

// func TestRandomAll(t *testing.T) {
// 	testRandom(t, true)
// }

// func TestRandomPublic(t *testing.T) {
// 	p, err := NewPex(Config{
// 		MaxPeers: 10,
// 	})
// 	require.NoError(t, err)

// 	f := p.RandomPublic

// 	// check without peers
// 	assert.NotNil(t, p.RandomPublic(100))
// 	assert.Equal(t, len(p.RandomPublic(100)), 0)

// 	// check with one peer
// 	p.AddPeer("112.32.32.14:10011")
// 	// 0 defaults to all peers
// 	compareRandom(t, p, 0, []string{"112.32.32.14:10011"}, f)
// 	compareRandom(t, p, 1, []string{"112.32.32.14:10011"}, f)
// 	// exceeding known peers is safe
// 	compareRandom(t, p, 2, []string{"112.32.32.14:10011"}, f)
// 	// exceeding max peers is safe
// 	compareRandom(t, p, 100, []string{"112.32.32.14:10011"}, f)

// 	// check with two peers
// 	p.AddPeer("112.32.32.14:20011")
// 	one := p.RandomAll(1)[0].String()
// 	if one != "112.32.32.14:10011" && one != "112.32.32.14:20011" {
// 		assert.Nil(t, nil)
// 	}
// 	// 0 defaults to all peers
// 	compareRandom(t, p, 0, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	}, f)
// 	compareRandom(t, p, 2, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	}, f)
// 	compareRandom(t, p, 3, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	}, f)
// 	compareRandom(t, p, 100, []string{
// 		"112.32.32.14:10011",
// 		"112.32.32.14:20011",
// 	}, f)

// 	// check with 3 peers, one private
// 	p.AddPeer("112.32.32.14:30011")
// 	p.SetPrivate("112.32.32.14:30011", true)
// 	if publicOnly {
// 		// The private peer should never be included
// 		compareRandom(t, p, 0, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 		}, f)
// 		compareRandom(t, p, 2, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 		}, f)
// 		compareRandom(t, p, 3, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 		}, f)
// 		compareRandom(t, p, 100, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 		}, f)
// 	} else {
// 		// The private peer should be included
// 		compareRandom(t, p, 0, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 			"112.32.32.14:30011",
// 		}, f)
// 		compareRandom(t, p, 3, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 			"112.32.32.14:30011",
// 		}, f)
// 		compareRandom(t, p, 4, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 			"112.32.32.14:30011",
// 		}, f)
// 		compareRandom(t, p, 100, []string{
// 			"112.32.32.14:10011",
// 			"112.32.32.14:20011",
// 			"112.32.32.14:30011",
// 		}, f)
// 	}
// }

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
	p, err := NewPeerFromJSON(pj)
	require.NoError(t, err)
	check(p)

	pj = load(newFormat)
	p, err = NewPeerFromJSON(pj)
	require.NoError(t, err)
	check(p)
}

/* Addendum: dummies & mocks */

// Fake addr that satisfies net.Addr interface
// type dummyAddr struct{}

// func (da *dummyAddr) Network() string {
// 	return da.String()
// }

// func (da *dummyAddr) String() string {
// 	return "none"
// }

// Fake connection that satisfies net.Conn interface
// type dummyConnection struct{}

// func (dc *dummyConnection) Read(b []byte) (int, error) {
// 	return 0, nil
// }

// func (dc *dummyConnection) Write(b []byte) (int, error) {
// 	return 0, nil
// }

// func (dc *dummyConnection) Close() error {
// 	return nil
// }

// func (dc *dummyConnection) LocalAddr() net.Addr {
// 	return &dummyAddr{}
// }

// func (dc *dummyConnection) RemoteAddr() net.Addr {
// 	return &dummyAddr{}
// }

// func (dc *dummyConnection) SetDeadline(t time.Time) error {
// 	return nil
// }

// func (dc *dummyConnection) SetReadDeadline(t time.Time) error {
// 	return nil
// }

// func (dc *dummyConnection) SetWriteDeadline(t time.Time) error {
// 	return nil
// }
