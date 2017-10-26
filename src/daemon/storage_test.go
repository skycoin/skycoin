package daemon

import (
	"sync"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/stretchr/testify/assert"
)

func TestStoreAdd(t *testing.T) {
	testData := []struct {
		expect int
		values map[string]string
	}{
		{
			1,
			map[string]string{
				"a": "a_value",
			},
		},
		{
			2,
			map[string]string{
				"a": "a_value",
				"b": "b_value",
			},
		},
		{
			3,
			map[string]string{
				"a": "a_value",
				"b": "b_value",
				"c": "c_value",
			},
		},
	}

	for _, d := range testData {
		s := store{value: make(map[interface{}]interface{})}
		for k, v := range d.values {
			s.setValue(k, v)
		}
		assert.Equal(t, d.expect, len(s.value))
	}
}

func TestStoreRemove(t *testing.T) {
	s := store{value: make(map[interface{}]interface{})}
	s.setValue("a", "a_value")
	s.setValue("b", "b_value")

	s.remove("a")
	assert.Equal(t, 1, len(s.value))
	s.remove("b")
	assert.Equal(t, 0, len(s.value))
}

func TestStoreGet(t *testing.T) {
	s := store{value: make(map[interface{}]interface{})}
	s.setValue("a", "a_value")
	s.setValue("b", "b_value")
	v, ok := s.getValue("a")
	assert.True(t, ok)
	assert.Equal(t, "a_value", v.(string))
}

func TestStoreLen(t *testing.T) {
	s := store{value: make(map[interface{}]interface{})}
	s.setValue("a", "a_value")
	s.setValue("b", "b_value")
	assert.Equal(t, s.len(), len(s.value))
}

func TestNewExpectIntroduction(t *testing.T) {
	ei := NewExpectIntroductions()
	assert.NotNil(t, ei)
	assert.NotNil(t, ei.store.value)
}

func TestExpectIntroAdd(t *testing.T) {
	ei := NewExpectIntroductions()
	now := utc.Now()
	ei.Add("a", now)
	assert.Equal(t, 1, len(ei.store.value))
}

func TestExpectIntroGet(t *testing.T) {
	ei := NewExpectIntroductions()
	now := utc.Now()
	ei.Add("a", now)
	tm, ok := ei.Get("a")
	assert.True(t, ok)
	assert.Equal(t, now, tm)
}

func TestExpectIntroRemove(t *testing.T) {
	ei := NewExpectIntroductions()
	now := utc.Now()
	ei.Add("a", now)
	ei.Add("b", now.Add(1))
	ei.Add("c", now.Add(2))
	assert.Equal(t, 3, len(ei.store.value))
	ei.Remove("a")
	assert.Equal(t, 2, len(ei.store.value))
	_, ok := ei.Get("a")
	assert.False(t, ok)
	bt, ok := ei.Get("b")
	assert.True(t, ok)
	assert.Equal(t, now.Add(1), bt)
	ct, ok := ei.Get("c")
	assert.True(t, ok)
	assert.Equal(t, now.Add(2), ct)
}

func TestCullInvalidConnections(t *testing.T) {
	ei := NewExpectIntroductions()
	now := utc.Now()
	ei.Add("a", now)
	ei.Add("b", now.Add(1))
	ei.Add("c", now.Add(2))
	wg := sync.WaitGroup{}
	vc := make(chan string, 3)
	wg.Add(2)
	go func(w *sync.WaitGroup) {
		as, err := ei.CullInvalidConns(func(addr string, tm time.Time) (bool, error) {
			if addr == "a" || addr == "b" {
				return true, nil
			}
			return false, nil
		})
		assert.Nil(t, err)

		for _, s := range as {
			vc <- s
		}
		w.Done()
	}(&wg)

	go func(w *sync.WaitGroup) {
		// w.Add(1)
		as, err := ei.CullInvalidConns(func(addr string, tm time.Time) (bool, error) {
			if addr == "c" {
				return true, nil
			}
			return false, nil
		})

		assert.Nil(t, err)

		for _, s := range as {
			vc <- s
		}
		w.Done()
	}(&wg)
	wg.Wait()
	assert.Equal(t, 3, len(vc))
}

func TestNewConnectionMirrors(t *testing.T) {
	cm := NewConnectionMirrors()
	assert.NotNil(t, cm)
	assert.NotNil(t, cm.value)
}

func TestAddConnMirrors(t *testing.T) {
	testData := []struct {
		expectNum int
		value     map[string]uint32
	}{
		{
			1,
			map[string]uint32{
				"a": 1,
			},
		},
		{
			2,
			map[string]uint32{
				"a": 1,
				"b": 2,
			},
		},
		{
			3,
			map[string]uint32{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
	}

	for _, data := range testData {
		cm := NewConnectionMirrors()
		for a := range data.value {
			cm.Add(a, data.value[a])
		}
		assert.Equal(t, data.expectNum, len(cm.value))
		for a := range data.value {
			v, ok := cm.Get(a)
			assert.True(t, ok)
			m := cm.value[a].(uint32)
			assert.Equal(t, v, m)
		}
	}
}

func TestConnMirrorsRemove(t *testing.T) {
	cm := NewConnectionMirrors()
	cm.Add("a", 1)
	cm.Remove("a")
	assert.Equal(t, 0, len(cm.value))

	cm.Add("a", 1)
	cm.Add("b", 2)
	cm.Remove("a")
	assert.Equal(t, 1, len(cm.value))
	_, ok := cm.Get("a")
	assert.False(t, ok)
	_, ok = cm.Get("b")
	assert.True(t, ok)
}

func TestNewOutgoingConnections(t *testing.T) {
	oc := NewOutgoingConnections(3)
	assert.NotNil(t, oc)
	assert.NotNil(t, oc.value)
	assert.Equal(t, 0, len(oc.value))
}

func TestOutgoingConnAdd(t *testing.T) {
	oc := NewOutgoingConnections(3)
	oc.Add("a")
	assert.Equal(t, 1, len(oc.value))
	oc.Add("b")
	assert.Equal(t, 2, len(oc.value))
}

func TestOutgoingConnGet(t *testing.T) {
	oc := NewOutgoingConnections(3)
	oc.Add("a")
	oc.Add("b")
	assert.True(t, oc.Get("a"))
	assert.True(t, oc.Get("b"))
	assert.False(t, oc.Get("c"))
}

func TestOutgoingConnLen(t *testing.T) {
	oc := NewOutgoingConnections(3)
	oc.Add("a")
	oc.Add("b")
	assert.Equal(t, oc.Len(), 2)
}

func TestNewPendingConns(t *testing.T) {
	pc := NewPendingConnections(3)
	assert.NotNil(t, pc)
	assert.NotNil(t, pc.value)
	assert.Equal(t, 0, len(pc.value))
}

func TestPendingConnAdd(t *testing.T) {
	pc := NewPendingConnections(3)
	pc.Add("a", pex.Peer{Addr: "a"})
	pc.Add("b", pex.Peer{Addr: "b"})
	assert.Equal(t, 2, len(pc.value))
	a := pc.value["a"].(pex.Peer)
	b := pc.value["b"].(pex.Peer)

	assert.Equal(t, pex.Peer{Addr: "a"}, a)
	assert.Equal(t, pex.Peer{Addr: "b"}, b)
}

func TestPendingConnGet(t *testing.T) {
	pc := NewPendingConnections(3)
	pc.Add("a", pex.Peer{Addr: "a"})
	pc.Add("b", pex.Peer{Addr: "b"})
	v, ok := pc.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "a", v.Addr)

	v, ok = pc.Get("b")
	assert.True(t, ok)
	assert.Equal(t, "b", v.Addr)

}

func TestPendingConnRemove(t *testing.T) {
	pc := NewPendingConnections(3)
	pc.Add("a", pex.Peer{Addr: "a"})
	pc.Add("b", pex.Peer{Addr: "b"})
	assert.Equal(t, 2, len(pc.value))
	pc.Remove("a")
	assert.Equal(t, 1, len(pc.value))
	_, ok := pc.Get("a")
	assert.False(t, ok)
	_, ok = pc.Get("b")
	assert.True(t, ok)
}

func TestPendingConnLen(t *testing.T) {
	pc := NewPendingConnections(3)
	pc.Add("a", pex.Peer{Addr: "a"})
	pc.Add("b", pex.Peer{Addr: "b"})
	assert.Equal(t, 2, pc.Len())
}

func TestNewMirrorConnections(t *testing.T) {
	mc := NewMirrorConnections()
	assert.NotNil(t, mc)
	assert.NotNil(t, mc.value)
	assert.Equal(t, 0, len(mc.value))
}

func TestMirrorConnAdd(t *testing.T) {
	mc := NewMirrorConnections()
	mc.Add(1, "a", 1)
	mc.Add(1, "b", 1)
	assert.Equal(t, 1, len(mc.value))
	assert.Equal(t, 2, len(mc.value[uint32(1)].(map[string]uint16)))
}

func TestMirrorConnGet(t *testing.T) {
	mc := NewMirrorConnections()
	mc.Add(1, "a", 1)
	mc.Add(1, "b", 2)
	p, ok := mc.Get(1, "a")
	assert.True(t, ok)
	assert.Equal(t, uint16(1), p)
	p, ok = mc.Get(1, "b")
	assert.True(t, ok)
	assert.Equal(t, uint16(2), p)
	p, ok = mc.Get(1, "c")
	assert.False(t, ok)
	p, ok = mc.Get(uint32(2), "a")
	assert.False(t, ok)
}

func TestMirrorConnRemove(t *testing.T) {
	mc := NewMirrorConnections()
	mc.Add(1, "a", 1)
	mc.Add(1, "b", 2)
	mc.Add(2, "c", 1)
	mc.Remove(1, "a")
	_, ok := mc.value[uint32(1)].(map[string]uint16)["a"]
	assert.False(t, ok)
	p, ok := mc.value[uint32(1)].(map[string]uint16)["b"]
	assert.True(t, ok)
	assert.Equal(t, uint16(2), p)
}

func TestNewIPCount(t *testing.T) {
	ic := NewIPCount()
	assert.NotNil(t, ic)
	assert.NotNil(t, ic.value)
}

func TestIPCountIncrease(t *testing.T) {
	ic := NewIPCount()
	ic.Increase("a")
	assert.Equal(t, 1, ic.value["a"].(int))
	ic.Increase("a")
	assert.Equal(t, 2, ic.value["a"].(int))
}

func TestIPCountDecrease(t *testing.T) {
	ic := NewIPCount()
	ic.Increase("a")
	assert.Equal(t, 1, ic.value["a"].(int))
	ic.Increase("a")
	assert.Equal(t, 2, ic.value["a"].(int))
	ic.Increase("b")
	assert.Equal(t, 1, ic.value["b"].(int))
	assert.Equal(t, 2, len(ic.value))

	ic.Decrease("a")
	assert.Equal(t, 1, ic.value["a"].(int))
	assert.Equal(t, 2, len(ic.value))
	assert.Equal(t, 1, ic.value["b"].(int))
}
