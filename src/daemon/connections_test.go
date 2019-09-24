package daemon

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/util/useragent"
)

var userAgent = useragent.MustParse("skycoin:0.24.1(foo)")

func getMirrorPort(c *Connections, ip string, mirror uint32) uint16 {
	c.Lock()
	defer c.Unlock()

	x := c.mirrors[mirror]
	if x == nil {
		return 0
	}

	return x[ip]
}

func TestConnectionsOutgoingFlow(t *testing.T) {
	conns := NewConnections()

	ip := "127.0.0.1"
	port := uint16(6060)
	addr := fmt.Sprintf("%s:%d", ip, port)

	all := conns.all()
	require.Empty(t, all)

	require.Equal(t, 0, conns.IPCount(ip))
	require.Equal(t, 0, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 0, conns.Len())
	require.Empty(t, conns.mirrors)

	// Flow: pending, connected, introduced

	c, err := conns.pending(addr)
	require.NoError(t, err)

	require.True(t, c.Outgoing)
	require.Equal(t, addr, c.Addr)
	require.Equal(t, port, c.ListenPort)
	require.Equal(t, ConnectionStatePending, c.State)
	require.Equal(t, 1, conns.IPCount(ip))
	require.Equal(t, 1, conns.OutgoingLen())
	require.Equal(t, 1, conns.PendingLen())
	require.Equal(t, 1, conns.Len())
	require.Empty(t, c.Mirror)
	require.Empty(t, conns.mirrors)
	require.False(t, c.HasIntroduced())
	require.Equal(t, addr, c.ListenAddr())

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	_, err = conns.pending(addr)
	require.Equal(t, ErrConnectionExists, err)
	require.Equal(t, 1, conns.IPCount(ip))
	require.Equal(t, 1, conns.OutgoingLen())
	require.Equal(t, 1, conns.PendingLen())
	require.Equal(t, 1, conns.Len())
	require.Empty(t, conns.mirrors)
	require.False(t, c.HasIntroduced())
	require.Equal(t, addr, c.ListenAddr())

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	_, err = conns.introduced(addr, 1, &IntroductionMessage{
		UserAgent: userAgent,
	})
	require.Equal(t, ErrConnectionStateNotConnected, err)
	require.Equal(t, 1, conns.PendingLen())

	c, err = conns.connected(addr, 1)
	require.NoError(t, err)

	require.True(t, c.Outgoing)
	require.Equal(t, addr, c.Addr)
	require.Equal(t, port, c.ListenPort)
	require.Equal(t, ConnectionStateConnected, c.State)
	require.Equal(t, 1, conns.IPCount(ip))
	require.Equal(t, 1, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 1, conns.Len())
	require.Empty(t, c.Mirror)
	require.Empty(t, conns.mirrors)
	require.False(t, c.HasIntroduced())
	require.Equal(t, addr, c.ListenAddr())

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	m := &IntroductionMessage{
		// use a different port to make sure we don't overwrite the true listen port
		ListenPort:      port + 1,
		Mirror:          1111,
		ProtocolVersion: 2,
		UserAgent:       userAgent,
	}

	c, err = conns.introduced(addr, 1, m)
	require.NoError(t, err)

	require.True(t, c.Outgoing)
	require.Equal(t, addr, c.Addr)
	require.Equal(t, port, c.ListenPort)
	require.Equal(t, ConnectionStateIntroduced, c.State)
	require.Equal(t, 1, conns.IPCount(ip))
	require.Equal(t, 1, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 1, conns.Len())
	require.Equal(t, m.Mirror, c.Mirror)
	require.Equal(t, m.ProtocolVersion, c.ProtocolVersion)
	require.Len(t, conns.mirrors, 1)
	require.Equal(t, port, getMirrorPort(conns, ip, c.Mirror))
	require.True(t, c.HasIntroduced())
	require.Equal(t, addr, c.ListenAddr())
	require.Equal(t, userAgent, c.UserAgent)

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	require.Equal(t, c, conns.get(c.Addr))
	require.Equal(t, c, conns.getByGnetID(c.gnetID))
	require.Equal(t, []*connection{c}, conns.getByListenAddr(c.ListenAddr()))

	err = conns.remove(addr, 1)
	require.NoError(t, err)

	require.Equal(t, 0, conns.IPCount(ip))
	require.Equal(t, 0, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 0, conns.Len())
	require.Empty(t, conns.mirrors)

	all = conns.all()
	require.Empty(t, all)

	require.Nil(t, conns.getByGnetID(c.gnetID))
	require.Nil(t, conns.getByListenAddr(c.ListenAddr()))
	require.Nil(t, conns.get(c.Addr))
}

func TestConnectionsIncomingFlow(t *testing.T) {
	conns := NewConnections()

	ip := "127.0.0.1"
	port := uint16(6060)
	addr := fmt.Sprintf("%s:%d", ip, port)

	all := conns.all()
	require.Empty(t, all)

	require.Equal(t, 0, conns.IPCount(ip))
	require.Equal(t, 0, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 0, conns.Len())
	require.Empty(t, conns.mirrors)

	// Flow: connected, introduced

	c, err := conns.connected(addr, 1)
	require.NoError(t, err)

	require.False(t, c.Outgoing)
	require.Equal(t, addr, c.Addr)
	require.Empty(t, c.ListenPort)
	require.Equal(t, ConnectionStateConnected, c.State)
	require.Equal(t, 1, conns.IPCount(ip))
	require.Equal(t, 0, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 1, conns.Len())
	require.Empty(t, c.Mirror)
	require.Empty(t, conns.mirrors)
	require.False(t, c.HasIntroduced())
	require.Empty(t, c.ListenAddr())
	require.Empty(t, conns.listenAddrs)

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	m := &IntroductionMessage{
		// use a different port to make sure that we use the self-reported listen port for incoming connections
		ListenPort:      port + 1,
		Mirror:          1111,
		ProtocolVersion: 2,
		UserAgent:       userAgent,
		UnconfirmedVerifyTxn: params.VerifyTxn{
			BurnFactor:          4,
			MaxTransactionSize:  1111,
			MaxDropletPrecision: 2,
		},
	}

	c, err = conns.introduced(addr, 1, m)
	require.NoError(t, err)

	require.False(t, c.Outgoing)
	require.Equal(t, addr, c.Addr)
	require.Equal(t, m.ListenPort, c.ListenPort)
	require.Equal(t, ConnectionStateIntroduced, c.State)
	require.Equal(t, 1, conns.IPCount(ip))
	require.Equal(t, 0, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 1, conns.Len())
	require.Equal(t, m.Mirror, c.Mirror)
	require.Equal(t, m.ProtocolVersion, c.ProtocolVersion)
	require.Len(t, conns.mirrors, 1)
	require.Equal(t, m.ListenPort, getMirrorPort(conns, ip, c.Mirror))
	require.True(t, c.HasIntroduced())
	require.Equal(t, fmt.Sprintf("%s:%d", ip, m.ListenPort), c.ListenAddr())
	require.Equal(t, userAgent, c.UserAgent)
	require.Equal(t, uint32(4), c.UnconfirmedVerifyTxn.BurnFactor)
	require.Equal(t, uint32(1111), c.UnconfirmedVerifyTxn.MaxTransactionSize)
	require.Equal(t, uint8(2), c.UnconfirmedVerifyTxn.MaxDropletPrecision)

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	err = conns.remove(addr, 1)
	require.NoError(t, err)

	require.Equal(t, 0, conns.IPCount(ip))
	require.Equal(t, 0, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 0, conns.Len())
	require.Empty(t, conns.mirrors)

	all = conns.all()
	require.Empty(t, all)

	require.Nil(t, conns.getByGnetID(c.gnetID))
	require.Nil(t, conns.getByListenAddr(c.ListenAddr()))
	require.Nil(t, conns.get(c.Addr))
}

func TestConnectionsMultiple(t *testing.T) {
	conns := NewConnections()

	addr1 := "127.0.0.1:6060"
	addr2 := "127.0.0.1:6061"

	_, err := conns.pending(addr1)
	require.NoError(t, err)

	_, err = conns.pending(addr2)
	require.NoError(t, err)

	require.Equal(t, 2, conns.OutgoingLen())
	require.Equal(t, 2, conns.PendingLen())
	require.Equal(t, 2, conns.IPCount("127.0.0.1"))

	_, err = conns.connected(addr1, 1)
	require.NoError(t, err)
	require.Equal(t, 1, conns.PendingLen())

	_, err = conns.connected(addr2, 2)
	require.NoError(t, err)
	require.Equal(t, 0, conns.PendingLen())

	_, err = conns.introduced(addr1, 1, &IntroductionMessage{
		Mirror:          6,
		ListenPort:      6060,
		ProtocolVersion: 2,
		UserAgent:       userAgent,
	})
	require.NoError(t, err)
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 1, len(conns.mirrors))

	// introduction fails if a base IP + mirror is already in use
	_, err = conns.introduced(addr2, 2, &IntroductionMessage{
		Mirror:          6,
		ListenPort:      6061,
		ProtocolVersion: 2,
		UserAgent:       userAgent,
	})
	require.Equal(t, ErrConnectionIPMirrorExists, err)
	require.Equal(t, 0, conns.PendingLen())

	c := conns.get(addr2)
	require.Equal(t, ConnectionStateConnected, c.State)
	require.Equal(t, c, conns.getByGnetID(2))
	require.Equal(t, []*connection{c}, conns.getByListenAddr("127.0.0.1:6061"))

	_, err = conns.introduced(addr2, 2, &IntroductionMessage{
		Mirror:          7,
		ListenPort:      6061,
		ProtocolVersion: 2,
		UserAgent:       userAgent,
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(conns.mirrors))
	require.Equal(t, 2, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 2, conns.Len())
	require.Equal(t, 2, conns.IPCount("127.0.0.1"))
	c = conns.get(addr2)
	require.Equal(t, ConnectionStateIntroduced, c.State)
	require.Equal(t, c, conns.getByGnetID(2))
	require.Equal(t, []*connection{c}, conns.getByListenAddr("127.0.0.1:6061"))

	// Add another connection with a different base IP but same mirror value
	addr3 := "127.1.1.1:12345"
	_, err = conns.connected(addr3, 3)
	require.NoError(t, err)

	_, err = conns.introduced(addr3, 3, &IntroductionMessage{
		Mirror:          6,
		ListenPort:      6060,
		ProtocolVersion: 2,
		UserAgent:       userAgent,
	})
	require.NoError(t, err)

	require.Equal(t, 2, len(conns.mirrors))
	require.Equal(t, 2, len(conns.mirrors[6]))
	require.Equal(t, uint16(6060), getMirrorPort(conns, "127.0.0.1", 6))
	require.Equal(t, uint16(6061), getMirrorPort(conns, "127.0.0.1", 7))
	require.Equal(t, uint16(6060), getMirrorPort(conns, "127.1.1.1", 6))

	require.Equal(t, 2, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 3, conns.Len())
	require.Equal(t, 2, conns.IPCount("127.0.0.1"))
	require.Equal(t, 1, conns.IPCount("127.1.1.1"))

	err = conns.remove(addr1, 2)
	require.Equal(t, ErrConnectionGnetIDMismatch, err)
	require.Equal(t, 3, conns.Len())

	err = conns.remove(addr1, 1)
	require.NoError(t, err)
	err = conns.remove(addr2, 2)
	require.NoError(t, err)
	err = conns.remove(addr3, 3)
	require.NoError(t, err)
	require.Empty(t, conns.mirrors)
	require.Equal(t, 0, conns.Len())

	err = conns.remove(addr1, 1)
	require.Equal(t, ErrConnectionNotExist, err)
}

func TestConnectionsMultipleSameListenPort(t *testing.T) {
	conns := NewConnections()

	addr1 := "127.0.0.1:6060"
	addr2 := "127.0.0.1:51414"

	c, err := conns.pending(addr1)
	require.NoError(t, err)
	require.Equal(t, []*connection{c}, conns.getByListenAddr(addr1))

	c2, err := conns.connected(addr2, 2)
	require.NoError(t, err)
	require.Equal(t, []*connection{c}, conns.getByListenAddr(addr1))

	_, err = conns.introduced(addr2, 2, &IntroductionMessage{
		Mirror:          6,
		ListenPort:      6060,
		ProtocolVersion: 2,
		UserAgent:       userAgent,
	})
	require.NoError(t, err)

	listenAddrConns := conns.getByListenAddr(addr1)
	require.Len(t, listenAddrConns, 2)
	require.True(t, *listenAddrConns[0] == *c || *listenAddrConns[0] == *c2)
	if *listenAddrConns[0] == *c {
		require.Equal(t, c2, listenAddrConns[1])
	} else if *listenAddrConns[0] == *c2 {
		require.Equal(t, c, listenAddrConns[1])
	}

	err = conns.remove(addr1, 0)
	require.NoError(t, err)

	listenAddrConns = conns.getByListenAddr(addr1)
	require.Len(t, listenAddrConns, 1)
	require.Equal(t, c2, listenAddrConns[0])

	err = conns.remove(addr2, 2)
	require.NoError(t, err)
	require.Len(t, conns.getByListenAddr(addr1), 0)

	err = conns.remove(addr2, 2)
	require.Equal(t, ErrConnectionNotExist, err)

	require.Len(t, conns.listenAddrs, 0)
}

func TestConnectionsErrors(t *testing.T) {
	conns := NewConnections()

	_, err := conns.pending("foo")
	testutil.RequireError(t, err, "address foo: missing port in address")

	_, err = conns.connected("foo", 1)
	testutil.RequireError(t, err, "address foo: missing port in address")

	_, err = conns.introduced("foo", 1, &IntroductionMessage{})
	testutil.RequireError(t, err, "address foo: missing port in address")

	err = conns.remove("foo", 0)
	testutil.RequireError(t, err, "address foo: missing port in address")

	_, err = conns.introduced("127.0.0.1:6060", 1, &IntroductionMessage{})
	require.Equal(t, ErrConnectionNotExist, err)

	_, err = conns.connected("127.0.0.1:6060", 0)
	require.Equal(t, ErrInvalidGnetID, err)

	_, err = conns.introduced("127.0.0.1:6060", 0, &IntroductionMessage{})
	require.Equal(t, ErrInvalidGnetID, err)
}

func TestConnectionsSetHeight(t *testing.T) {
	conns := NewConnections()
	addr := "127.0.0.1:6060"
	height := uint64(1010)

	err := conns.SetHeight(addr, 1, height)
	require.Equal(t, ErrConnectionNotExist, err)

	c, err := conns.connected(addr, 1)
	require.NoError(t, err)
	require.Empty(t, c.Height)

	err = conns.SetHeight(addr, 1, height)
	require.NoError(t, err)

	err = conns.SetHeight(addr, 2, height)
	require.Equal(t, ErrConnectionGnetIDMismatch, err)

	c = conns.get(addr)
	require.NotNil(t, c)
	require.Equal(t, height, c.Height)
}

func TestConnectionsModifyMirrorPanics(t *testing.T) {
	conns := NewConnections()
	addr := "127.0.0.1:6060"

	_, err := conns.connected(addr, 1)
	require.NoError(t, err)

	// modifying mirror value causes panic
	require.Panics(t, func() {
		conns.modify(addr, 1, func(c *ConnectionDetails) { //nolint:errcheck
			c.Mirror++
		})
	})

	// modifying ListenPort causes panic
	require.Panics(t, func() {
		conns.modify(addr, 1, func(c *ConnectionDetails) { //nolint:errcheck
			c.ListenPort = 999
		})
	})
}

func TestConnectionsStateTransitionErrors(t *testing.T) {
	conns := NewConnections()
	addr := "127.0.0.1:6060"

	_, err := conns.pending(addr)
	require.NoError(t, err)

	// pending -> pending fails
	_, err = conns.pending(addr)
	require.Equal(t, ErrConnectionExists, err)

	// pending -> introduced fails
	_, err = conns.introduced(addr, 1, &IntroductionMessage{})
	require.Equal(t, ErrConnectionStateNotConnected, err)

	_, err = conns.connected(addr, 1)
	require.NoError(t, err)

	// connected -> connected fails
	_, err = conns.connected(addr, 1)
	require.Equal(t, ErrConnectionAlreadyConnected, err)

	// connected -> introduced fails if gnet ID does not match
	_, err = conns.introduced(addr, 2, &IntroductionMessage{})
	require.Equal(t, ErrConnectionGnetIDMismatch, err)

	_, err = conns.introduced(addr, 1, &IntroductionMessage{})
	require.NoError(t, err)

	// introduced -> connected fails
	_, err = conns.connected(addr, 1)
	require.Equal(t, ErrConnectionAlreadyIntroduced, err)

	// introduced -> pending fails
	_, err = conns.pending(addr)
	require.Equal(t, ErrConnectionExists, err)

	// introduced -> introduced fails
	_, err = conns.introduced(addr, 1, &IntroductionMessage{})
	require.Equal(t, ErrConnectionAlreadyIntroduced, err)
}
