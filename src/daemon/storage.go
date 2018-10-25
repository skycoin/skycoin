package daemon

import (
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon/pex"
)

// ExpectIntroductions records connections that are expecting introduction msg.
type ExpectIntroductions struct {
	value map[string]time.Time
	sync.Mutex
}

// CullMatchFunc function for checking if the connection need to be culled
type CullMatchFunc func(addr string, t time.Time) (bool, error)

// NewExpectIntroductions creates a ExpectIntroduction instance
func NewExpectIntroductions() *ExpectIntroductions {
	return &ExpectIntroductions{
		value: make(map[string]time.Time),
	}
}

// Add adds expecting introduction connection
func (s *ExpectIntroductions) Add(addr string, tm time.Time) {
	s.Lock()
	defer s.Unlock()
	s.value[addr] = tm
}

// Get returns the time of specific address
func (s *ExpectIntroductions) Get(addr string) (time.Time, bool) {
	s.Lock()
	defer s.Unlock()
	t, ok := s.value[addr]
	return t, ok
}

// Remove removes connection
func (s *ExpectIntroductions) Remove(addr string) {
	s.Lock()
	defer s.Unlock()
	delete(s.value, addr)
}

// CullInvalidConns cull connections that match the matchFunc
func (s *ExpectIntroductions) CullInvalidConns(f CullMatchFunc) ([]string, error) {
	s.Lock()
	defer s.Unlock()

	var addrs []string
	for addr, t := range s.value {
		ok, err := f(addr, t)
		if err != nil {
			return nil, err
		}

		if ok {
			addrs = append(addrs, addr)
		}
	}

	for _, addr := range addrs {
		delete(s.value, addr)
	}

	return addrs, nil
}

// ConnectionMirrors records mirror for connection
type ConnectionMirrors struct {
	value map[string]uint32
	sync.Mutex
}

// NewConnectionMirrors create ConnectionMirrors instance.
func NewConnectionMirrors() *ConnectionMirrors {
	return &ConnectionMirrors{
		value: make(map[string]uint32),
	}
}

// Add adds connection mirror
func (s *ConnectionMirrors) Add(addr string, mirror uint32) {
	s.Lock()
	defer s.Unlock()
	s.value[addr] = mirror
}

// Get returns the mirror of connection
func (s *ConnectionMirrors) Get(addr string) (uint32, bool) {
	s.Lock()
	defer s.Unlock()
	v, ok := s.value[addr]
	return v, ok
}

// Remove remove connection mirror
func (s *ConnectionMirrors) Remove(addr string) {
	s.Lock()
	defer s.Unlock()
	delete(s.value, addr)
}

// StringSet existence set for string values
type StringSet struct {
	value map[string]struct{}
	sync.Mutex
}

// NewStringSet create StringSet
func NewStringSet(size int) *StringSet {
	return &StringSet{
		value: make(map[string]struct{}, size),
	}
}

// Add records connection
func (s *StringSet) Add(addr string) {
	s.Lock()
	defer s.Unlock()
	s.value[addr] = struct{}{}
}

// Remove remove connection
func (s *StringSet) Remove(addr string) {
	s.Lock()
	defer s.Unlock()
	delete(s.value, addr)
}

// Get returns if connection is outgoing
func (s *StringSet) Get(addr string) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.value[addr]
	return ok
}

// Len returns the outgoing connections count
func (s *StringSet) Len() int {
	s.Lock()
	defer s.Unlock()
	return len(s.value)
}

// PendingConnections records pending connection peers
type PendingConnections struct {
	value map[string]pex.Peer
	sync.Mutex
}

// NewPendingConnections creates new PendingConnections instance
func NewPendingConnections(size int) *PendingConnections {
	return &PendingConnections{
		value: make(map[string]pex.Peer, size),
	}
}

// Add adds pending connection
func (s *PendingConnections) Add(peer pex.Peer) {
	s.Lock()
	defer s.Unlock()
	s.value[peer.Addr] = peer
}

// Get returns pending connections
func (s *PendingConnections) Get(addr string) (pex.Peer, bool) {
	s.Lock()
	defer s.Unlock()
	v, ok := s.value[addr]
	return v, ok
}

// Remove removes pending connection
func (s *PendingConnections) Remove(addr string) {
	s.Lock()
	defer s.Unlock()
	delete(s.value, addr)
}

// Len returns pending connection number
func (s *PendingConnections) Len() int {
	s.Lock()
	defer s.Unlock()
	return len(s.value)
}

// MirrorConnections records mirror connections
type MirrorConnections struct {
	value map[uint32]map[string]uint16
	sync.Mutex
}

// NewMirrorConnections create mirror connection instance
func NewMirrorConnections() *MirrorConnections {
	return &MirrorConnections{
		value: make(map[uint32]map[string]uint16),
	}
}

// Add adds mirror connection
func (s *MirrorConnections) Add(mirror uint32, ip string, port uint16) {
	s.Lock()
	defer s.Unlock()

	if m, ok := s.value[mirror]; ok {
		m[ip] = port
		return
	}

	m := make(map[string]uint16)
	m[ip] = port
	s.value[mirror] = m
}

// Get returns ip port of specific mirror
func (s *MirrorConnections) Get(mirror uint32, ip string) (uint16, bool) {
	s.Lock()
	defer s.Unlock()

	m, ok := s.value[mirror]
	if ok {
		port, exist := m[ip]
		return port, exist
	}

	return 0, false
}

// Remove removes port of ip for specific mirror
func (s *MirrorConnections) Remove(mirror uint32, ip string) {
	s.Lock()
	defer s.Unlock()

	m, ok := s.value[mirror]
	if ok {
		delete(m, ip)
	}
}

// IPCount records connection number from the same base ip
type IPCount struct {
	value map[string]int
	sync.Mutex
}

// NewIPCount returns IPCount instance
func NewIPCount() *IPCount {
	return &IPCount{
		value: make(map[string]int),
	}
}

// Increase increases one for specific ip
func (s *IPCount) Increase(ip string) {
	s.Lock()
	defer s.Unlock()

	if c, ok := s.value[ip]; ok {
		c++
		s.value[ip] = c
		return
	}

	s.value[ip] = 1
}

// Decrease decreases one for specific ip
func (s *IPCount) Decrease(ip string) {
	s.Lock()
	defer s.Unlock()

	if c, ok := s.value[ip]; ok {
		if c <= 1 {
			delete(s.value, ip)
			return
		}
		c--
		s.value[ip] = c
	}
}

// Get return ip count
func (s *IPCount) Get(ip string) (int, bool) {
	s.Lock()
	defer s.Unlock()
	v, ok := s.value[ip]
	return v, ok
}

type announcedTxnsCache struct {
	sync.Mutex
	cache map[cipher.SHA256]int64
}

func newAnnouncedTxnsCache() *announcedTxnsCache {
	return &announcedTxnsCache{
		cache: make(map[cipher.SHA256]int64),
	}
}

func (c *announcedTxnsCache) add(txns []cipher.SHA256) {
	c.Lock()
	defer c.Unlock()

	t := time.Now().UTC().UnixNano()
	for _, txn := range txns {
		c.cache[txn] = t
	}
}

func (c *announcedTxnsCache) flush() map[cipher.SHA256]int64 {
	c.Lock()
	defer c.Unlock()

	if len(c.cache) == 0 {
		return nil
	}

	cache := c.cache

	c.cache = make(map[cipher.SHA256]int64)

	return cache
}
