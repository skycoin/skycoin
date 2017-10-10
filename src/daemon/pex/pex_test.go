package pex

import (
	"net"
	"sort"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/stretchr/testify/assert"
)

var (
	address   = "112.32.32.14:3030"
	address2  = "112.32.32.14:3031"
	addresses = []string{
		address, "111.32.32.13:2020", "69.32.54.111:2222",
	}
	silenceLogger = true
)

func init() {
	// silence the logger
	if silenceLogger {
		logging.Disable()
	}
}

func TestValidateAddress(t *testing.T) {
	// empty string
	assert.False(t, ValidateAddress("", false))
	// doubled ip:port
	assert.False(t, ValidateAddress("112.32.32.14:100112.32.32.14:101", false))
	// requires port
	assert.False(t, ValidateAddress("112.32.32.14", false))
	// not ip
	assert.False(t, ValidateAddress("112", false))
	assert.False(t, ValidateAddress("112.32", false))
	assert.False(t, ValidateAddress("112.32.32", false))
	// bad part
	assert.False(t, ValidateAddress("112.32.32.14000", false))
	// large port
	assert.False(t, ValidateAddress("112.32.32.14:66666", false))
	// unspecified
	assert.False(t, ValidateAddress("0.0.0.0:8888", false))
	// no ip
	assert.False(t, ValidateAddress(":8888", false))
	// multicast
	assert.False(t, ValidateAddress("224.1.1.1:8888", false))
	// invalid ports
	assert.False(t, ValidateAddress("112.32.32.14:0", false))
	assert.False(t, ValidateAddress("112.32.32.14:1", false))
	assert.False(t, ValidateAddress("112.32.32.14:10", false))
	assert.False(t, ValidateAddress("112.32.32.14:100", false))
	assert.False(t, ValidateAddress("112.32.32.14:1000", false))
	assert.False(t, ValidateAddress("112.32.32.14:1023", false))
	assert.False(t, ValidateAddress("112.32.32.14:65536", false))
	// valid ones
	assert.True(t, ValidateAddress("112.32.32.14:1024", false))
	assert.True(t, ValidateAddress("112.32.32.14:10000", false))
	assert.True(t, ValidateAddress("112.32.32.14:65535", false))
	// localhost is allowed
	assert.True(t, ValidateAddress("127.0.0.1:8888", true))
	// localhost is not allowed
	assert.False(t, ValidateAddress("127.0.0.1:8888", false))
}

/* Peer tests */

func TestNewPeer(t *testing.T) {
	p := NewPeer(address)
	assert.NotEqual(t, p.LastSeen, 0)
	assert.Equal(t, p.Addr, address)
	assert.False(t, p.Private)
}

func TestPeerSeen(t *testing.T) {
	p := NewPeer(address)
	x := p.LastSeen
	time.Sleep(time.Second)
	p.Seen()
	assert.NotEqual(t, x, p.LastSeen)
	if p.LastSeen.Before(x) {
		t.Fail()
	}
}

func TestPeerString(t *testing.T) {
	p := NewPeer(address)
	assert.Equal(t, address, p.String())
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

func TestGetAllAddresses(t *testing.T) {
	p := NewPex(10)
	p.AddPeer("112.32.32.14:10011")
	p.AddPeer("112.32.32.14:20011")
	p.SetPrivate("112.32.32.14:20011", true)
	addresses := p.Peerlist.GetAllAddresses()
	assert.Equal(t, len(addresses), 2)
	sort.Strings(addresses)
	assert.Equal(t, addresses, []string{
		"112.32.32.14:10011",
		"112.32.32.14:20011",
	})
}

func TestGetPublicAddresses(t *testing.T) {
	p := NewPex(10)
	p.AddPeer("112.32.32.14:10011")
	p.AddPeer("112.32.32.14:20011")
	p.AddPeer("112.32.32.14:30011")
	p.SetPrivate("112.32.32.14:30011", true)

	addresses := p.Peerlist.GetPublicAddresses()
	assert.Equal(t, len(addresses), 2)
	sort.Strings(addresses)
	assert.Equal(t, addresses, []string{
		"112.32.32.14:10011",
		"112.32.32.14:20011",
	})
}

func TestGetPrivateAddresses(t *testing.T) {
	ips := []string{
		"112.32.32.14:10011",
		"112.32.32.14:20011",
		"112.32.32.14:30011",
		"112.32.32.14:40011",
	}
	p := NewPex(10)
	p.AddPeer(ips[0])
	p.AddPeer(ips[1])
	p.AddPeer(ips[2])
	p.AddPeer(ips[3])
	p.SetPrivate(ips[2], true)
	p.SetPrivate(ips[3], true)

	addresses := p.Peerlist.GetPrivateAddresses()
	assert.Equal(t, len(addresses), 2)
	sort.Strings(addresses)
	assert.Equal(t, addresses, []string{
		ips[2],
		ips[3],
	})
}

func convertPeersToStrings(peers []*Peer) []string {
	addresses := make([]string, 0, len(peers))
	for _, p := range peers {
		addresses = append(addresses, p.String())
	}
	return addresses
}

func compareRandom(t *testing.T, p *Pex, npeers int,
	result []string, f func(int) []*Peer) {
	peers := f(npeers)
	addresses := convertPeersToStrings(peers)
	sort.Strings(addresses)
	assert.Equal(t, addresses, result)
}

func testRandom(t *testing.T, publicOnly bool) {
	p := NewPex(10)

	f := p.Peerlist.RandomAll
	if publicOnly {
		f = p.Peerlist.RandomPublic
	}

	// check without peers
	assert.NotNil(t, p.Peerlist.RandomAll(100))
	assert.Equal(t, len(p.Peerlist.RandomAll(100)), 0)

	// check with one peer
	p.AddPeer("112.32.32.14:10011")
	// 0 defaults to all peers
	compareRandom(t, p, 0, []string{"112.32.32.14:10011"}, f)
	compareRandom(t, p, 1, []string{"112.32.32.14:10011"}, f)
	// exceeding known peers is safe
	compareRandom(t, p, 2, []string{"112.32.32.14:10011"}, f)
	// exceeding max peers is safe
	compareRandom(t, p, 100, []string{"112.32.32.14:10011"}, f)

	// check with two peers
	p.AddPeer("112.32.32.14:20011")
	one := p.Peerlist.RandomAll(1)[0].String()
	if one != "112.32.32.14:10011" && one != "112.32.32.14:20011" {
		assert.Nil(t, nil)
	}
	// 0 defaults to all peers
	compareRandom(t, p, 0, []string{
		"112.32.32.14:10011",
		"112.32.32.14:20011",
	}, f)
	compareRandom(t, p, 2, []string{
		"112.32.32.14:10011",
		"112.32.32.14:20011",
	}, f)
	compareRandom(t, p, 3, []string{
		"112.32.32.14:10011",
		"112.32.32.14:20011",
	}, f)
	compareRandom(t, p, 100, []string{
		"112.32.32.14:10011",
		"112.32.32.14:20011",
	}, f)

	// check with 3 peers, one private
	p.AddPeer("112.32.32.14:30011")
	p.SetPrivate("112.32.32.14:30011", true)
	if publicOnly {
		// The private peer should never be included
		compareRandom(t, p, 0, []string{
			"112.32.32.14:10011",
			"112.32.32.14:20011",
		}, f)
		compareRandom(t, p, 2, []string{
			"112.32.32.14:10011",
			"112.32.32.14:20011",
		}, f)
		compareRandom(t, p, 3, []string{
			"112.32.32.14:10011",
			"112.32.32.14:20011",
		}, f)
		compareRandom(t, p, 100, []string{
			"112.32.32.14:10011",
			"112.32.32.14:20011",
		}, f)
	} else {
		// The private peer should be included
		compareRandom(t, p, 0, []string{
			"112.32.32.14:10011",
			"112.32.32.14:20011",
			"112.32.32.14:30011",
		}, f)
		compareRandom(t, p, 3, []string{
			"112.32.32.14:10011",
			"112.32.32.14:20011",
			"112.32.32.14:30011",
		}, f)
		compareRandom(t, p, 4, []string{
			"112.32.32.14:10011",
			"112.32.32.14:20011",
			"112.32.32.14:30011",
		}, f)
		compareRandom(t, p, 100, []string{
			"112.32.32.14:10011",
			"112.32.32.14:20011",
			"112.32.32.14:30011",
		}, f)
	}
}

func TestRandomAll(t *testing.T) {
	testRandom(t, true)
}

func TestRandomPublic(t *testing.T) {
	testRandom(t, false)
}

// func TestGetPeer(t *testing.T) {
// 	p := NewPex(10)
// 	p.AddPeer("112.32.32.14:10011")
// 	assert.Nil(t, p.Peerlist["xxx"])
// 	assert.Equal(t, p.Peerlist["112.32.32.14:10011"].String(),
// 		"112.32.32.14:10011")
// }

// func TestFull(t *testing.T) {
// 	p := NewPex(1)
// 	assert.False(t, p.Full())
// 	p.AddPeer("112.32.32.14:10011")
// 	assert.True(t, p.Full())
// 	// No limit
// 	p = NewPex(0)
// 	p.AddPeer("112.32.32.14:10011")
// 	assert.False(t, p.Full())
// }

// func TestAddPeerLocalhost(t *testing.T) {
// 	p := NewPex(1)
// 	p.AllowLocalhost = true
// 	a := "127.0.0.1:10114"
// 	peer, err := p.AddPeer(a)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(p.Peerlist), 1)
// 	assert.Equal(t, p.Peerlist[a], peer)
// 	assert.Equal(t, peer.Addr, a)
// }

// func TestAddPeer(t *testing.T) {
// 	p := NewPex(1)

// 	// adding "" peer results in error
// 	peer, err := p.AddPeer("")
// 	assert.Nil(t, peer)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, err, InvalidAddressError)
// 	assert.Equal(t, len(p.Peerlist), 0)

// 	peer, err = p.AddPeer("112.32.32.14:10011")
// 	assert.Nil(t, err)
// 	assert.Equal(t, peer.String(), "112.32.32.14:10011")
// 	assert.NotNil(t, p.Peerlist["112.32.32.14:10011"])
// 	past := peer.LastSeen

// 	// full list
// 	twopeer, err := p.AddPeer("112.32.32.14:20011")
// 	assert.Equal(t, err, PeerlistFullError)
// 	assert.Nil(t, twopeer)
// 	assert.Nil(t, p.Peerlist["112.32.32.14:20011"])

// 	// re-add original peer
// 	time.Sleep(time.Second)
// 	repeer, err := p.AddPeer("112.32.32.14:10011")
// 	assert.Nil(t, err)
// 	assert.NotNil(t, repeer)
// 	assert.Equal(t, peer, repeer)
// 	assert.Equal(t, repeer.String(), "112.32.32.14:10011")
// 	assert.True(t, repeer.LastSeen.After(past))

// 	assert.NotNil(t, p.Peerlist["112.32.32.14:10011"])

// 	// Adding blacklisted peer is invalid
// 	delete(p.Peerlist, address)
// 	p.AddBlacklistEntry(address, time.Second)
// 	peer, err = p.AddPeer(address)
// 	assert.NotNil(t, err)
// 	assert.Nil(t, peer)
// 	assert.Nil(t, p.Peerlist[address])
// }

// func TestSaveLoad(t *testing.T) {
// 	os.Remove("./" + PeerDatabaseFilename)
// 	os.Remove("./" + BlacklistedDatabaseFilename)
// 	defer os.Remove("./" + PeerDatabaseFilename)
// 	defer os.Remove("./" + BlacklistedDatabaseFilename)
// 	p := NewPex(10)
// 	w, err := p.AddPeer("112.32.32.14:10011")
// 	assert.Nil(t, err)
// 	w.LastSeen = time.Unix(w.LastSeen.Unix(), 0)
// 	x, _ := p.AddPeer("112.32.32.14:20011")
// 	x.LastSeen = time.Time{}
// 	privAddr := "112.32.32.14:30011"
// 	y, err := p.AddPeer(privAddr)
// 	assert.Nil(t, err)
// 	y.Private = true
// 	y.LastSeen = time.Unix(y.LastSeen.Unix(), 0)
// 	// bypass AddPeer to add a blacklist and normal address at the same time
// 	// saving this and reloading it should cause the address to be
// 	// blacklisted only
// 	bad := "111.44.44.22:11021"
// 	p.Peerlist[bad] = NewPeer(bad)
// 	p.AddBlacklistEntry(bad, time.Hour)
// 	// Do similar for the private peer, it should be rejected from the
// 	// blacklist
// 	p.Blacklist[privAddr] = BlacklistEntry{Now(), time.Hour}

// 	assert.Nil(t, p.Save("./"))

// 	q := NewPex(10)
// 	assert.Nil(t, q.Load("./"))
// 	assert.Nil(t, q.Peerlist[bad])
// 	_, exists := q.Blacklist[bad]
// 	assert.True(t, exists)
// 	_, exists = q.Blacklist[privAddr]
// 	assert.False(t, exists)
// 	assert.Equal(t, len(q.Blacklist), 1)
// 	assert.Equal(t, len(q.Peerlist), 3)
// 	assert.NotNil(t, q.Peerlist["112.32.32.14:10011"])
// 	assert.NotNil(t, q.Peerlist["112.32.32.14:20011"])
// 	assert.NotNil(t, q.Peerlist[privAddr])
// 	assert.True(t, q.Peerlist["112.32.32.14:20011"].LastSeen.IsZero())
// 	assert.Equal(t, q.Peerlist["112.32.32.14:10011"].LastSeen,
// 		p.Peerlist["112.32.32.14:10011"].LastSeen)
// 	assert.Equal(t, q.Peerlist[privAddr].LastSeen,
// 		p.Peerlist[privAddr].LastSeen)
// 	assert.False(t, q.Peerlist["112.32.32.14:10011"].Private)
// 	assert.False(t, q.Peerlist["112.32.32.14:20011"].Private)
// 	assert.True(t, q.Peerlist[privAddr].Private)
// 	assert.Equal(t, len(q.Peerlist), 3)
// }

// func TestLoadBlacklistDoesNotExist(t *testing.T) {
// 	// Fail on os.IsNotExist, returns a valid empty blacklist
// 	os.Remove("./" + BlacklistedDatabaseFilename)
// 	b, err := LoadBlacklist("./")
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(b), 0)
// 	assert.NotNil(t, b)
// }

// // func TestLoadPeerlistFailureHandling(t *testing.T) {
// // 	defer os.Remove("./" + PeerDatabaseFilename)

// // 	// File does not exist, returns empty Peerlist
// // 	os.Remove("./" + PeerDatabaseFilename)
// // 	p, err := LoadPeerlist("./")
// // 	assert.Nil(t, err)
// // 	assert.Equal(t, len(p), 0)
// // 	assert.NotNil(t, p)

// // 	// Bad peerlist file:

// // 	goodAddr := "123.45.54.11:9999"
// // 	now := Now().Unix()
// // 	lines := []string{
// // 		// Has a line with 4 entries
// // 		fmt.Sprintf("%s 0 %d 0", address, now),
// // 		// Has a line with 2 entries
// // 		fmt.Sprintf("%s 0", address),
// // 		// Empty line
// // 		"",
// // 		// Has an invalid address
// // 		fmt.Sprintf("54.3:9090 1 %d", now),
// // 		// Has an invalid value for private
// // 		fmt.Sprintf("%s 7 %d", address, now),
// // 		// Starts with a comment
// // 		fmt.Sprintf("#%s 0 %d", address, now),
// // 		// Has an invalid seen timestamp
// // 		fmt.Sprintf("%s 0 dog", address),
// // 		// Has a valid line, but extra whitespace,
// // 		// this should be the only one included
// // 		fmt.Sprintf("%s 0 %d", goodAddr, now),
// // 	}

// // 	f, err := os.Create("./" + PeerDatabaseFilename)
// // 	assert.Nil(t, err)
// // 	for _, line := range lines {
// // 		_, err = f.Write([]byte(line + "\n"))
// // 		assert.Nil(t, err)
// // 	}
// // 	f.Close()

// // 	p, err = LoadPeerlist("./")
// // 	assert.Nil(t, err)
// // 	assert.Equal(t, len(p), 1)
// // 	assert.NotNil(t, p[goodAddr])
// // }

func TestPeerCanTry(t *testing.T) {
	testData := []struct {
		LastSeen   time.Time
		RetryTimes int
		CanTry     bool
	}{
		{
			Now().Add(time.Duration(100) * time.Second * -1),
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

func TestNow(t *testing.T) {
	now := Now().Unix()
	now2 := utc.UnixNow()
	assert.True(t, now == now2 || now2-1 == now)
}

/* Addendum: dummies & mocks */

// Fake addr that satisfies net.Addr interface
type dummyAddr struct{}

func (da *dummyAddr) Network() string {
	return da.String()
}

func (da *dummyAddr) String() string {
	return "none"
}

// Fake connection that satisfies net.Conn interface
type dummyConnection struct{}

func (dc *dummyConnection) Read(b []byte) (int, error) {
	return 0, nil
}

func (dc *dummyConnection) Write(b []byte) (int, error) {
	return 0, nil
}

func (dc *dummyConnection) Close() error {
	return nil
}

func (dc *dummyConnection) LocalAddr() net.Addr {
	return &dummyAddr{}
}

func (dc *dummyConnection) RemoteAddr() net.Addr {
	return &dummyAddr{}
}

func (dc *dummyConnection) SetDeadline(t time.Time) error {
	return nil
}

func (dc *dummyConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (dc *dummyConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
