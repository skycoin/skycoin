package daemon

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/testutil"
)

func getMirrorPort(c *Connections, ip string, mirror uint32) uint16 {
	c.Lock()
	defer c.Unlock()

	logger.Debugf("getMirrorPort ip=%s mirror=%d", ip, mirror)

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
	require.Equal(t, ErrConnectionAlreadyRegistered, err)
	require.Equal(t, 1, conns.IPCount(ip))
	require.Equal(t, 1, conns.OutgoingLen())
	require.Equal(t, 1, conns.PendingLen())
	require.Equal(t, 1, conns.Len())
	require.Empty(t, conns.mirrors)
	require.False(t, c.HasIntroduced())
	require.Equal(t, addr, c.ListenAddr())

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	c, err = conns.connected(addr)
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
	}

	c, err = conns.introduced(addr, m)
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

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	err = conns.remove(addr)
	require.NoError(t, err)

	require.Equal(t, 0, conns.IPCount(ip))
	require.Equal(t, 0, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 0, conns.Len())
	require.Empty(t, conns.mirrors)

	all = conns.all()
	require.Empty(t, all)

	c = conns.get(addr)
	require.Nil(t, c)
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

	c, err := conns.connected(addr)
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

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	m := &IntroductionMessage{
		// use a different port to make sure that we use the self-reported listen port for incoming connections
		ListenPort:      port + 1,
		Mirror:          1111,
		ProtocolVersion: 2,
	}

	c, err = conns.introduced(addr, m)
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

	all = conns.all()
	require.Equal(t, []connection{*c}, all)

	err = conns.remove(addr)
	require.NoError(t, err)

	require.Equal(t, 0, conns.IPCount(ip))
	require.Equal(t, 0, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 0, conns.Len())
	require.Empty(t, conns.mirrors)

	all = conns.all()
	require.Empty(t, all)

	c = conns.get(addr)
	require.Nil(t, c)
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

	_, err = conns.connected(addr1)
	require.NoError(t, err)
	require.Equal(t, 1, conns.PendingLen())

	_, err = conns.connected(addr2)
	require.NoError(t, err)
	require.Equal(t, 0, conns.PendingLen())

	_, err = conns.introduced(addr1, &IntroductionMessage{
		Mirror:          6,
		ListenPort:      6060,
		ProtocolVersion: 2,
	})
	require.NoError(t, err)
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 1, len(conns.mirrors))

	// introduction fails if a base IP + mirror is already in use
	_, err = conns.introduced(addr2, &IntroductionMessage{
		Mirror:          6,
		ListenPort:      6061,
		ProtocolVersion: 2,
	})
	require.Equal(t, ErrConnectionIPMirrorAlreadyRegistered, err)
	require.Equal(t, 0, conns.PendingLen())

	c := conns.get(addr2)
	require.Equal(t, ConnectionStateConnected, c.State)

	_, err = conns.introduced(addr2, &IntroductionMessage{
		Mirror:          7,
		ListenPort:      6061,
		ProtocolVersion: 2,
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(conns.mirrors))
	require.Equal(t, 2, conns.OutgoingLen())
	require.Equal(t, 0, conns.PendingLen())
	require.Equal(t, 2, conns.Len())
	require.Equal(t, 2, conns.IPCount("127.0.0.1"))

	// Add another connection with a different base IP but same mirror value
	addr3 := "127.1.1.1:12345"
	_, err = conns.connected(addr3)
	require.NoError(t, err)

	_, err = conns.introduced(addr3, &IntroductionMessage{
		Mirror:          6,
		ListenPort:      6060,
		ProtocolVersion: 2,
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

	err = conns.remove(addr1)
	require.NoError(t, err)
	err = conns.remove(addr2)
	require.NoError(t, err)
	err = conns.remove(addr3)
	require.NoError(t, err)
	require.Empty(t, conns.mirrors)
	require.Equal(t, 0, conns.Len())
}

func TestConnectionsErrors(t *testing.T) {
	conns := NewConnections()

	_, err := conns.pending("foo")
	testutil.RequireError(t, err, "address foo: missing port in address")

	_, err = conns.connected("foo")
	testutil.RequireError(t, err, "address foo: missing port in address")

	_, err = conns.introduced("foo", nil)
	testutil.RequireError(t, err, "address foo: missing port in address")

	err = conns.remove("foo")
	testutil.RequireError(t, err, "address foo: missing port in address")

	_, err = conns.introduced("127.0.0.1:6060", &IntroductionMessage{})
	require.Equal(t, ErrConnectionNotExist, err)
}

func TestConnectionsSetHeight(t *testing.T) {
	conns := NewConnections()
	addr := "127.0.0.1:6060"
	height := uint64(1010)

	err := conns.SetHeight(addr, height)
	require.Equal(t, ErrConnectionNotExist, err)

	c, err := conns.connected(addr)
	require.NoError(t, err)
	require.Empty(t, c.Height)

	err = conns.SetHeight(addr, height)
	require.NoError(t, err)

	c = conns.get(addr)
	require.NotNil(t, c)
	require.Equal(t, height, c.Height)
}

func TestConnectionsModifyMirrorPanics(t *testing.T) {
	conns := NewConnections()
	addr := "127.0.0.1:6060"

	_, err := conns.connected(addr)
	require.NoError(t, err)

	// modifying mirror value causes panic
	require.Panics(t, func() {
		conns.modify(addr, func(c *ConnectionDetails) { // nolint: errcheck
			c.Mirror++
		})
	})
}
