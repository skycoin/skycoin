package daemon

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/daemon/pex"
)

func TestExpectIntroductions(t *testing.T) {
	ei := NewExpectIntroductions()

	_, ok := ei.Get("foo")
	require.False(t, ok)

	tm := time.Now()
	ei.Add("foo", tm)
	tm2, ok := ei.Get("foo")
	require.True(t, ok)
	require.Equal(t, tm, tm2)

	ei.Remove("foo")
	_, ok = ei.Get("foo")
	require.False(t, ok)
}

func TestExpectIntroductionsCullInvalidConnections(t *testing.T) {
	ei := NewExpectIntroductions()
	now := time.Now().UTC()
	ei.Add("a", now)
	ei.Add("b", now.Add(1))
	ei.Add("c", now.Add(2))

	wg := sync.WaitGroup{}
	vc := make(chan string, 3)

	wg.Add(2)
	go func() {
		defer wg.Done()
		as, err := ei.CullInvalidConns(func(addr string, tm time.Time) (bool, error) {
			if addr == "a" || addr == "b" {
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
		as, err := ei.CullInvalidConns(func(addr string, tm time.Time) (bool, error) {
			if addr == "c" {
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

	_, ok := ei.Get("a")
	require.False(t, ok)
	_, ok = ei.Get("b")
	require.False(t, ok)
	_, ok = ei.Get("c")
	require.False(t, ok)
}

func TestConnectionMirrors(t *testing.T) {
	cm := NewConnectionMirrors()

	_, ok := cm.Get("foo")
	require.False(t, ok)

	cm.Add("foo", 10)
	c, ok := cm.Get("foo")
	require.True(t, ok)
	require.Equal(t, uint32(10), c)

	cm.Remove("foo")
	_, ok = cm.Get("foo")
	require.False(t, ok)
}

func TestStringSet(t *testing.T) {
	oc := NewStringSet(3)

	n := oc.Len()
	require.Equal(t, 0, n)

	ok := oc.Get("foo")
	require.False(t, ok)

	oc.Add("foo")
	ok = oc.Get("foo")
	require.True(t, ok)

	n = oc.Len()
	require.Equal(t, 1, n)

	oc.Add("foo")
	ok = oc.Get("foo")
	require.True(t, ok)

	n = oc.Len()
	require.Equal(t, 1, n)

	oc.Add("foo2")
	ok = oc.Get("foo2")
	require.True(t, ok)

	n = oc.Len()
	require.Equal(t, 2, n)

	oc.Remove("foo")
	ok = oc.Get("foo")
	require.False(t, ok)

	n = oc.Len()
	require.Equal(t, 1, n)
}

func TestPendingConns(t *testing.T) {
	pc := NewPendingConnections(3)

	n := pc.Len()
	require.Equal(t, 0, n)

	_, ok := pc.Get("foo")
	require.False(t, ok)

	pc.Add(pex.Peer{Addr: "foo"})
	p, ok := pc.Get("foo")
	require.Equal(t, pex.Peer{Addr: "foo"}, p)
	require.True(t, ok)

	n = pc.Len()
	require.Equal(t, 1, n)

	pc.Add(pex.Peer{Addr: "foo"})
	p, ok = pc.Get("foo")
	require.Equal(t, pex.Peer{Addr: "foo"}, p)
	require.True(t, ok)

	n = pc.Len()
	require.Equal(t, 1, n)

	pc.Add(pex.Peer{Addr: "foo2"})
	p, ok = pc.Get("foo2")
	require.Equal(t, pex.Peer{Addr: "foo2"}, p)
	require.True(t, ok)

	n = pc.Len()
	require.Equal(t, 2, n)
}

func TestMirrorConnections(t *testing.T) {
	mc := NewMirrorConnections()

	_, ok := mc.Get(99, "foo")
	require.False(t, ok)

	mc.Add(99, "foo", 10)
	c, ok := mc.Get(99, "foo")
	require.True(t, ok)
	require.Equal(t, uint16(10), c)

	mc.Remove(99, "foo")
	_, ok = mc.Get(99, "foo")
	require.False(t, ok)

	mc.Add(99, "foo2", 10)
	c, ok = mc.Get(99, "foo2")
	require.True(t, ok)
	require.Equal(t, uint16(10), c)

	mc.Add(99, "foo", 10)
	c, ok = mc.Get(99, "foo")
	require.True(t, ok)
	require.Equal(t, uint16(10), c)

	mc.Remove(99, "foo2")
	_, ok = mc.Get(99, "foo2")
	require.False(t, ok)

	_, ok = mc.Get(99, "foo")
	require.True(t, ok)
}

func TestIPCount(t *testing.T) {
	ic := NewIPCount()

	_, ok := ic.Get("foo")
	require.False(t, ok)

	ic.Increase("foo")
	n, ok := ic.Get("foo")
	require.Equal(t, 1, n)
	require.True(t, ok)

	for i := 0; i < 3; i++ {
		ic.Decrease("foo")
		n, ok = ic.Get("foo")
		require.Equal(t, 0, n)
		require.False(t, ok)
	}

	for i := 0; i < 5; i++ {
		ic.Increase("foo")
		n, ok := ic.Get("foo")
		require.Equal(t, i+1, n)
		require.True(t, ok)
	}

	n, ok = ic.Get("foo")
	require.Equal(t, 5, n)
	require.True(t, ok)
}
