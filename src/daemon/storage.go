package daemon

import (
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/daemon/pex"
)

// base storage struct
type store struct {
	value map[interface{}]interface{}
	lk    sync.Mutex
}

type storeFunc func(*store) error
type matchFunc func(k interface{}, v interface{}) bool

func (s *store) setValue(k interface{}, v interface{}) {
	s.lk.Lock()
	s.value[k] = v
	s.lk.Unlock()
}

func (s *store) getValue(k interface{}) (interface{}, bool) {
	s.lk.Lock()
	defer s.lk.Unlock()
	v, ok := s.value[k]
	return v, ok
}

func (s *store) do(sf storeFunc) error {
	s.lk.Lock()
	defer s.lk.Unlock()
	return sf(s)
}

func (s *store) remove(k interface{}) {
	s.lk.Lock()
	delete(s.value, k)
	s.lk.Unlock()
}

func (s *store) len() int {
	s.lk.Lock()
	defer s.lk.Unlock()
	return len(s.value)
}

// ExpectIntroductions records connections that are expecting introduction msg.
type ExpectIntroductions struct {
	store
}

// CullMatchFunc function for checking if the connection need to be culled
type CullMatchFunc func(addr string, t time.Time) (bool, error)

// NewExpectIntroductions creates a ExpectIntroduction instance
func NewExpectIntroductions() *ExpectIntroductions {
	return &ExpectIntroductions{
		store: store{
			value: make(map[interface{}]interface{}),
		},
	}
}

// Add adds expecting introduction connection
func (ei *ExpectIntroductions) Add(addr string, tm time.Time) {
	ei.setValue(addr, tm)
}

// Remove removes connection
func (ei *ExpectIntroductions) Remove(addr string) {
	ei.remove(addr)
}

// CullInvalidConns cull connections that match the matchFunc
func (ei *ExpectIntroductions) CullInvalidConns(f CullMatchFunc) ([]string, error) {
	var addrs []string
	if err := ei.do(func(s *store) error {
		for k, v := range s.value {
			addr := k.(string)
			t := v.(time.Time)
			ok, err := f(addr, t)
			if err != nil {
				return err
			}

			if ok {
				addrs = append(addrs, addr)
				delete(s.value, k)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return addrs, nil
}

// Get returns the time of speicific address
func (ei *ExpectIntroductions) Get(addr string) (time.Time, bool) {
	if v, ok := ei.getValue(addr); ok {
		return v.(time.Time), ok
	}
	return time.Time{}, false
}

// ConnectionMirrors records mirror for connection
type ConnectionMirrors struct {
	store
}

// NewConnectionMirrors create ConnectionMirrors instance.
func NewConnectionMirrors() *ConnectionMirrors {
	return &ConnectionMirrors{
		store: store{
			value: make(map[interface{}]interface{}),
		},
	}
}

// Add adds connection mirror
func (cm *ConnectionMirrors) Add(addr string, mirror uint32) {
	cm.setValue(addr, mirror)
}

// Get returns the mirror of connection
func (cm *ConnectionMirrors) Get(addr string) (uint32, bool) {
	v, ok := cm.getValue(addr)
	if ok {
		return v.(uint32), ok
	}
	return 0, false
}

// Remove remove connection mirror
func (cm *ConnectionMirrors) Remove(addr string) {
	cm.remove(addr)
}

// OutgoingConnections records the outgoing connections
type OutgoingConnections struct {
	store
}

// NewOutgoingConnections create OutgoingConnection instance
func NewOutgoingConnections(max int) *OutgoingConnections {
	return &OutgoingConnections{
		store: store{
			value: make(map[interface{}]interface{}, max),
		},
	}
}

// Add records connection
func (oc *OutgoingConnections) Add(addr string) {
	oc.setValue(addr, true)
}

// Remove remove connection
func (oc *OutgoingConnections) Remove(addr string) {
	oc.remove(addr)
}

// Get returns if connection is outgoing
func (oc *OutgoingConnections) Get(addr string) bool {
	_, ok := oc.getValue(addr)
	return ok
}

// Len returns the outgoing connections count
func (oc *OutgoingConnections) Len() int {
	return oc.len()
}

// PendingConnections records pending connection peers
type PendingConnections struct {
	store
}

// NewPendingConnections creates new PendingConnections instance
func NewPendingConnections(maxConn int) *PendingConnections {
	return &PendingConnections{
		store: store{
			value: make(map[interface{}]interface{}, maxConn),
		},
	}
}

// Add adds pending connection
func (pc *PendingConnections) Add(addr string, peer pex.Peer) {
	pc.setValue(addr, peer)
}

// Get returns pending connections
func (pc *PendingConnections) Get(addr string) (pex.Peer, bool) {
	v, ok := pc.getValue(addr)
	if ok {
		return v.(pex.Peer), true
	}
	return pex.Peer{}, false
}

// Remove removes pending connection
func (pc *PendingConnections) Remove(addr string) {
	pc.remove(addr)
}

// Len returns pending connection number
func (pc *PendingConnections) Len() int {
	return pc.len()
}

// MirrorConnections records mirror connections
type MirrorConnections struct {
	store
}

// NewMirrorConnections create mirror connection instance
func NewMirrorConnections() *MirrorConnections {
	return &MirrorConnections{
		store: store{
			value: make(map[interface{}]interface{}),
		},
	}
}

// Add adds mirror connection
func (mc *MirrorConnections) Add(mirror uint32, ip string, port uint16) {
	mc.do(func(s *store) error {
		if m, ok := s.value[mirror]; ok {
			m.(map[string]uint16)[ip] = port
			return nil
		}

		m := make(map[string]uint16)
		m[ip] = port
		s.value[mirror] = m
		return nil
	})
}

// Get returns ip port of specific mirror
func (mc *MirrorConnections) Get(mirror uint32, ip string) (uint16, bool) {
	var port uint16
	var exist bool
	mc.do(func(s *store) error {
		if m, ok := s.value[mirror]; ok {
			port, exist = m.(map[string]uint16)[ip]
		}
		return nil
	})
	return port, exist
}

// Remove removes port of ip for specific mirror
func (mc *MirrorConnections) Remove(mirror uint32, ip string) {
	mc.do(func(s *store) error {
		if m, ok := s.value[mirror]; ok {
			delete(m.(map[string]uint16), ip)
		}
		return nil
	})
}

// IPCount records connection number from the same base ip
type IPCount struct {
	store
}

// NewIPCount returns IPCount instance
func NewIPCount() *IPCount {
	return &IPCount{
		store: store{
			value: make(map[interface{}]interface{}),
		},
	}
}

// Set sets ip count
// func (ic *IPCount) Set(ip string, n int) {
// 	ic.setValue(ip, n)
// }

// Increase increases one for specific ip
func (ic *IPCount) Increase(ip string) {
	ic.do(func(s *store) error {
		if v, ok := s.value[ip]; ok {
			c := v.(int)
			c++
			s.value[ip] = c
			return nil
		}

		s.value[ip] = 1
		return nil
	})
}

// Decrease decreases one for specific ip
func (ic *IPCount) Decrease(ip string) {
	ic.do(func(s *store) error {
		if v, ok := s.value[ip]; ok {
			c := v.(int)
			if c <= 1 {
				delete(s.value, ip)
				return nil
			}
			c--
			s.value[ip] = c
		}
		return nil
	})
}

// Get return ip count
func (ic *IPCount) Get(ip string) (int, bool) {
	v, ok := ic.getValue(ip)
	if ok {
		return v.(int), true
	}
	return 0, false
}
