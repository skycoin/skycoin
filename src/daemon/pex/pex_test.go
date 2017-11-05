package pex

import (
	"encoding/json"
	"net"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	if p.LastSeen < x {
		t.Fail()
	}
}

func TestPeerString(t *testing.T) {
	p := NewPeer(address)
	assert.Equal(t, address, p.String())
}

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
