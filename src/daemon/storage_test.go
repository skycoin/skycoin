package daemon

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/daemon/gnet"
)

func TestGet(t *testing.T) {
	conns := NewConnections()

	addr := "127.0.0.1:6060"
	details := ConnectionDetails{
		Mirror:      99,
		State:       ConnectionStatePending,
		ConnectedAt: time.Now().UTC(),
		ListenPort:  10101,
		Height:      1111,
	}

	c, err := conns.Add(fakeGnetConn(addr), details)
	require.NoError(t, err)
	require.Equal(t, details, c.ConnectionDetails)
	require.Equal(t, 111, c.GnetID)
	require.NotEmpty(t, c.LastReceived)
	require.NotEmpty(t, c.LastSent)

	c2, ok := conns.Get(addr)
	require.True(t, ok)
	require.Equal(t, c, c2)
}

func TestRemoveMatchedBy(t *testing.T) {
	ei := NewConnections()

	now := time.Now().UTC()

	addr1 := "127.0.0.1:6060"
	addr2 := "127.0.1.1:6061"
	addr3 := "127.1.1.1:6062"

	_, err := ei.Add(fakeGnetConn(addr1), ConnectionDetails{
		ConnectedAt: now,
	})
	require.NoError(t, err)

	_, err = ei.Add(fakeGnetConn(addr2), ConnectionDetails{
		ConnectedAt: now.Add(1),
	})
	require.NoError(t, err)

	_, err = ei.Add(fakeGnetConn(addr3), ConnectionDetails{
		ConnectedAt: now.Add(2),
	})
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	vc := make(chan string, 3)

	wg.Add(2)
	go func() {
		defer wg.Done()
		as, err := ei.RemoveMatchedBy(func(c Connection) (bool, error) {
			if c.Addr == addr1 || c.Addr == addr2 {
				return true, nil
			}
			return false, nil
		})
		require.NoError(t, err)

		for _, s := range as {
			vc <- s
		}
	}()

	go func() {
		defer wg.Done()
		as, err := ei.RemoveMatchedBy(func(c Connection) (bool, error) {
			if c.Addr == addr3 {
				return true, nil
			}
			return false, nil
		})

		require.NoError(t, err)

		for _, s := range as {
			vc <- s
		}
	}()

	wg.Wait()
	require.Equal(t, 3, len(vc))

	_, ok := ei.Get(addr1)
	require.False(t, ok)
	_, ok = ei.Get(addr2)
	require.False(t, ok)
	_, ok = ei.Get(addr3)
	require.False(t, ok)
}

func TestMirrorConnections(t *testing.T) {
	mc := NewConnections()

	localhost := "127.0.0.1"
	addr1 := "127.0.0.1:6060"
	addr2 := "127.0.0.1:6061"
	addr3 := "127.1.1.1:6060"

	c := mc.GetMirrorPort(localhost, 99)
	require.Equal(t, uint16(0), c)

	_, err := mc.Add(fakeGnetConn(addr1), ConnectionDetails{
		Mirror: 99,
	})
	require.NoError(t, err)

	c = mc.GetMirrorPort(localhost, 99)
	require.Equal(t, uint16(6060), c)

	err = mc.Remove(addr1)
	require.NoError(t, err)

	c = mc.GetMirrorPort(localhost, 99)
	require.Equal(t, uint16(0), c)

	_, err = mc.Add(fakeGnetConn(addr2), ConnectionDetails{
		Mirror: 99,
	})
	require.NoError(t, err)

	c = mc.GetMirrorPort(addr2, 99)
	require.Equal(t, uint16(6061), c)

	_, err = mc.Add(fakeGnetConn(addr1), ConnectionDetails{
		Mirror: 99,
	})
	require.Equal(t, ErrConnectionIPMirrorAlreadyRegistered, err)

	c = mc.GetMirrorPort(localhost, 99)
	require.Equal(t, uint16(6061), c)

	_, err = mc.Add(fakeGnetConn(addr1), ConnectionDetails{
		Mirror: 999,
	})

	c = mc.GetMirrorPort(localhost, 99)
	require.Equal(t, uint16(6061), c)
	c = mc.GetMirrorPort(localhost, 999)
	require.Equal(t, uint16(6060), c)

	_, err = mc.Add(fakeGnetConn(addr3), ConnectionDetails{
		Mirror: 99,
	})
	require.NoError(t, err)

	c = mc.GetMirrorPort("127.1.1.1", 99)
	require.Equal(t, uint16(6060), c)

	err = mc.Remove(addr2)
	require.NoError(t, err)

	c = mc.GetMirrorPort(localhost, 99)
	require.Equal(t, uint16(0), c)
	c = mc.GetMirrorPort(localhost, 999)
	require.Equal(t, uint16(6060), c)

	err = mc.Remove(addr1)
	require.NoError(t, err)

	c = mc.GetMirrorPort(localhost, 999)
	require.Equal(t, uint16(0), c)

	err = mc.Remove(addr3)
	require.NoError(t, err)

	c = mc.GetMirrorPort("127.1.1.1", 99)
	require.Equal(t, uint16(0), c)
}

func TestIPCount(t *testing.T) {
	ic := NewConnections()

	localhost := "127.0.0.1"
	addr1 := "127.0.0.1:6060"
	addr2 := "127.0.0.1:6061"

	require.Equal(t, 0, ic.GetIPCount(localhost))
	require.Equal(t, 0, ic.Len())

	_, ok := ic.Get(addr1)
	require.False(t, ok)

	_, err := ic.Add(fakeGnetConn(addr1), ConnectionDetails{
		Mirror: 1,
	})
	require.NoError(t, err)
	require.Equal(t, 1, ic.GetIPCount(localhost))

	_, err = ic.Add(fakeGnetConn(addr2), ConnectionDetails{
		Mirror: 2,
	})
	require.NoError(t, err)
	require.Equal(t, 2, ic.GetIPCount(localhost))

	err = ic.Remove(addr1)
	require.NoError(t, err)
	require.Equal(t, 1, ic.GetIPCount(localhost))

	err = ic.Remove(addr1)
	require.NoError(t, err)
	require.Equal(t, 1, ic.GetIPCount(localhost))

	err = ic.Remove(addr2)
	require.NoError(t, err)
	require.Equal(t, 0, ic.GetIPCount(localhost))
}

func TestConnectionHeightsAll(t *testing.T) {
	p := NewConnections()

	addr1 := "127.0.0.1:1234"
	addr2 := "127.0.0.1:5678"
	addr3 := "127.0.0.1:9999"

	require.Empty(t, p.conns)
	err := p.Remove(addr1)
	require.NoError(t, err)
	require.Empty(t, p.conns)
	require.Empty(t, p.mirrors)
	require.Empty(t, p.ipCounts)

	e := p.EstimateHeight(1)
	require.Equal(t, uint64(1), e)

	e = p.EstimateHeight(13)
	require.Equal(t, uint64(13), e)

	_, err = p.Add(fakeGnetConn(addr1), ConnectionDetails{
		Height: 10,
		Mirror: 1,
	})
	require.NoError(t, err)
	require.Len(t, p.conns, 1)

	records := p.All()
	require.Len(t, records, 1)
	require.Equal(t, addr1, records[0].Addr)
	require.Equal(t, 10, records[0].Height)

	err = p.Modify(addr1, func(c *ConnectionDetails) error {
		c.Height = 11
		return nil
	})
	require.NoError(t, err)
	require.Len(t, p.conns, 1)

	records = p.All()
	require.Len(t, records, 1)
	require.Equal(t, addr1, records[0].Addr)
	require.Equal(t, 11, records[0].Height)

	e = p.EstimateHeight(1)
	require.Equal(t, uint64(11), e)

	e = p.EstimateHeight(13)
	require.Equal(t, uint64(13), e)

	_, err = p.Add(fakeGnetConn(addr2), ConnectionDetails{
		Height: 12,
		Mirror: 2,
	})
	require.NoError(t, err)
	_, err = p.Add(fakeGnetConn(addr3), ConnectionDetails{
		Height: 12,
		Mirror: 3,
	})
	require.NoError(t, err)
	require.Len(t, p.conns, 3)
	require.Equal(t, 3, p.Len())

	records = p.All()
	require.Len(t, records, 3)
	require.Equal(t, addr1, records[0].Addr)
	require.Equal(t, 11, records[0].Height)
	require.Equal(t, addr2, records[1].Addr)
	require.Equal(t, 12, records[1].Height)
	require.Equal(t, addr3, records[2].Addr)
	require.Equal(t, 12, records[2].Height)

	e = p.EstimateHeight(1)
	require.Equal(t, uint64(12), e)

	e = p.EstimateHeight(13)
	require.Equal(t, uint64(13), e)

	_, err = p.Add(fakeGnetConn(addr3), ConnectionDetails{
		Height: 24,
		Mirror: 4,
	})
	require.NoError(t, err)
	e = p.EstimateHeight(13)
	require.Equal(t, uint64(24), e)
}

func fakeGnetConn(addr string) *gnet.Connection {
	return &gnet.Connection{
		ID:           111,
		LastReceived: time.Now().UTC(),
		LastSent:     time.Now().UTC(),
		Conn: fakeNetConn{
			remoteAddr: fakeNetAddr{
				addr: addr,
			},
		},
	}
}

type fakeNetAddr struct {
	net.Addr
	addr string
}

func (f fakeNetAddr) String() string {
	return f.addr
}

type fakeNetConn struct {
	net.Conn
	remoteAddr fakeNetAddr
}

func (f fakeNetConn) RemoteAddr() net.Addr {
	return f.remoteAddr
}
