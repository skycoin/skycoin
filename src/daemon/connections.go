package daemon

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/util/iputil"
	"github.com/SkycoinProject/skycoin/src/util/useragent"
)

// ConnectionState connection state in the state machine
// Connections have three states: "pending", "connected" and "introduced"
// A connection in the "pending" state has been selected to establish a TCP connection,
// but the connection has not been established yet.
// Only outgoing connections will ever be in the "pending" state;
// incoming connections begin at the "connected" state.
// A connection in the "connected" state has established a TCP connection,
// but has not completed the introduction handshake.
// A connection in the "introduced" state has completed the introduction handshake.
type ConnectionState string

const (
	// ConnectionStatePending prior to establishing a connection
	ConnectionStatePending ConnectionState = "pending"
	// ConnectionStateConnected connected, but not introduced
	ConnectionStateConnected ConnectionState = "connected"
	// ConnectionStateIntroduced connection has introduced itself
	ConnectionStateIntroduced ConnectionState = "introduced"
)

var (
	// ErrConnectionNotExist connection does not exist when performing an operation that requires it to exist
	ErrConnectionNotExist = errors.New("Connection does not exist")
	// ErrConnectionExists connection exists in Connections
	ErrConnectionExists = errors.New("Connection exists")
	// ErrConnectionIPMirrorExists connection exists for a given base IP and mirror
	ErrConnectionIPMirrorExists = errors.New("Connection exists with this base IP and mirror")
	// ErrConnectionStateNotConnected connect state is not "connected"
	ErrConnectionStateNotConnected = errors.New("Connection state is not \"connected\"")
	// ErrConnectionGnetIDMismatch gnet ID in argument does not match gnet ID on record
	ErrConnectionGnetIDMismatch = errors.New("Connection gnet ID does not match")
	// ErrConnectionAlreadyIntroduced attempted to make invalid state transition from introduced state
	ErrConnectionAlreadyIntroduced = errors.New("Connection is already in introduced state")
	// ErrConnectionAlreadyConnected attempted to make invalid state transition from connected state
	ErrConnectionAlreadyConnected = errors.New("Connection is already in connected state")
	// ErrInvalidGnetID invalid gnet ID value used as argument
	ErrInvalidGnetID = errors.New("Invalid gnet ID")
)

// ConnectionDetails connection data managed by daemon
type ConnectionDetails struct {
	State                ConnectionState
	Outgoing             bool
	ConnectedAt          time.Time
	Mirror               uint32
	ListenPort           uint16
	ProtocolVersion      int32
	Height               uint64
	UserAgent            useragent.Data
	UnconfirmedVerifyTxn params.VerifyTxn
	GenesisHash          cipher.SHA256
}

// HasIntroduced returns true if the connection has introduced
func (c ConnectionDetails) HasIntroduced() bool {
	switch c.State {
	case ConnectionStateIntroduced:
		return true
	default:
		return false
	}
}

type connection struct {
	Addr string
	ConnectionDetails
	gnetID uint64
}

// ListenAddr returns the addr that connection listens on, if available
func (c *connection) ListenAddr() string {
	if c.ListenPort == 0 {
		return ""
	}

	ip, _, err := iputil.SplitAddr(c.Addr)
	if err != nil {
		logger.Critical().WithError(err).WithField("addr", c.Addr).Error("connection.ListenAddr addr could not be split")
		return ""
	}

	return fmt.Sprintf("%s:%d", ip, c.ListenPort)
}

// Connections manages a collection of Connection
type Connections struct {
	conns       map[string]*connection
	mirrors     map[uint32]map[string]uint16
	ipCounts    map[string]int
	gnetIDs     map[uint64]string
	listenAddrs map[string][]string
	sync.Mutex
}

// NewConnections creates Connections
func NewConnections() *Connections {
	return &Connections{
		conns:       make(map[string]*connection, 32),
		mirrors:     make(map[uint32]map[string]uint16, 32),
		ipCounts:    make(map[string]int, 32),
		gnetIDs:     make(map[uint64]string, 32),
		listenAddrs: make(map[string][]string, 32),
	}
}

// pending adds a new pending outgoing connection
func (c *Connections) pending(addr string) (*connection, error) {
	c.Lock()
	defer c.Unlock()

	ip, port, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Critical().WithField("addr", addr).WithError(err).Error("Connections.pending called with invalid addr")
		return nil, err
	}

	if _, ok := c.conns[addr]; ok {
		return nil, ErrConnectionExists
	}

	c.ipCounts[ip]++

	conn := &connection{
		Addr: addr,
		ConnectionDetails: ConnectionDetails{
			State:      ConnectionStatePending,
			Outgoing:   true,
			ListenPort: port,
		},
	}

	c.conns[addr] = conn
	c.listenAddrs[addr] = append(c.listenAddrs[addr], addr)

	logger.WithField("addr", addr).Debug("Connections.pending")

	return c.conns[addr], nil
}

// connected the connection has connected
func (c *Connections) connected(addr string, gnetID uint64) (*connection, error) {
	c.Lock()
	defer c.Unlock()

	fields := logrus.Fields{
		"addr":   addr,
		"gnetID": gnetID,
	}

	if gnetID == 0 {
		logger.Critical().WithFields(fields).WithError(ErrInvalidGnetID).Error("Connections.connected called with invalid gnetID")
		return nil, ErrInvalidGnetID
	}

	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Critical().WithFields(fields).WithError(err).Error("Connections.connected called with invalid addr")
		return nil, err
	}

	conn := c.conns[addr]

	if conn == nil {
		c.ipCounts[ip]++

		conn = &connection{
			Addr: addr,
		}

		c.conns[addr] = conn
	} else {
		fields := logrus.Fields{
			"addr":       addr,
			"gnetID":     gnetID,
			"state":      conn.State,
			"outgoing":   conn.Outgoing,
			"connGnetID": conn.gnetID,
		}

		switch conn.State {
		case ConnectionStatePending:
		case ConnectionStateConnected:
			logger.Critical().WithFields(fields).Error("Connections.connected called on already connected connection")
			return nil, ErrConnectionAlreadyConnected
		case ConnectionStateIntroduced:
			logger.Critical().WithFields(fields).Error("Connections.connected called on already introduced connection")
			return nil, ErrConnectionAlreadyIntroduced
		default:
			logger.WithFields(fields).Panic("Connection state invalid")
		}
	}

	c.gnetIDs[gnetID] = addr
	conn.gnetID = gnetID
	conn.ConnectedAt = time.Now().UTC()
	conn.State = ConnectionStateConnected

	fields["outgoing"] = conn.Outgoing

	logger.WithFields(fields).Debug("Connections.connected")

	return conn, nil
}

// introduced the connection has introduced itself
func (c *Connections) introduced(addr string, gnetID uint64, m *IntroductionMessage) (*connection, error) {
	c.Lock()
	defer c.Unlock()

	fields := logrus.Fields{
		"addr":   addr,
		"gnetID": gnetID,
	}

	if gnetID == 0 {
		logger.Critical().WithFields(fields).WithError(ErrInvalidGnetID).Error("Connections.introduced called with invalid gnetID")
		return nil, ErrInvalidGnetID
	}

	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Critical().WithFields(fields).WithError(err).Error("Connections.introduced called with invalid addr")
		return nil, err
	}

	conn := c.conns[addr]
	if conn == nil {
		return nil, ErrConnectionNotExist
	}

	fields["outgoing"] = conn.Outgoing

	errorFields := logrus.Fields{
		"state":      conn.State,
		"connGnetID": conn.gnetID,
	}

	switch conn.State {
	case ConnectionStatePending:
		logger.Critical().WithFields(fields).WithFields(errorFields).Error("Connections.introduced called on pending connection")
		return nil, ErrConnectionStateNotConnected
	case ConnectionStateConnected:
		if gnetID != conn.gnetID {
			logger.Critical().WithFields(fields).WithFields(errorFields).Error("Connections.introduced called with different gnet ID")
			return nil, ErrConnectionGnetIDMismatch
		}
	case ConnectionStateIntroduced:
		logger.Critical().WithFields(fields).WithFields(errorFields).Error("Connections.introduced called on already introduced connection")
		return nil, ErrConnectionAlreadyIntroduced
	default:
		logger.WithFields(fields).WithFields(errorFields).Panic("invalid connection state")
	}

	if err := c.canUpdateMirror(ip, m.Mirror); err != nil {
		logger.WithFields(fields).WithFields(errorFields).WithField("mirror", m.Mirror).WithError(err).Debug("canUpdateMirror failed")
		return nil, err
	}

	// For outgoing connections, which are created by pending,
	// the listen port is set from the addr's port number.
	// Since we are connecting to it, it is presumed to be that peer's open listening port.
	// A misbehaving peer could report a different ListenPort in their IntroductionMessage,
	// but it shouldn't affect our records.
	if conn.Outgoing && conn.ListenPort != m.ListenPort {
		logger.Critical().WithFields(fields).WithFields(logrus.Fields{
			"connListenPort":  conn.ListenPort,
			"introListenPort": m.ListenPort,
		}).Warning("Outgoing connection's ListenPort does not match reported IntroductionMessage ListenPort")
	}

	listenPort := conn.ListenPort
	if !conn.Outgoing {
		listenPort = m.ListenPort
	}

	if err := c.updateMirror(ip, m.Mirror, listenPort); err != nil {
		logger.WithFields(fields).WithField("mirror", m.Mirror).WithError(err).Panic("updateMirror failed, but shouldn't")
	}

	conn.State = ConnectionStateIntroduced
	conn.Mirror = m.Mirror
	conn.ProtocolVersion = m.ProtocolVersion
	conn.ListenPort = listenPort
	conn.UserAgent = m.UserAgent
	conn.UnconfirmedVerifyTxn = m.UnconfirmedVerifyTxn
	conn.GenesisHash = m.GenesisHash

	if !conn.Outgoing {
		listenAddr := conn.ListenAddr()
		c.listenAddrs[listenAddr] = append(c.listenAddrs[listenAddr], addr)
	}

	logger.WithFields(fields).Debug("Connections.introduced")

	return conn, nil
}

// get returns a connection by address
func (c *Connections) get(addr string) *connection {
	c.Lock()
	defer c.Unlock()

	return c.conns[addr]
}

func (c *Connections) getByListenAddr(listenAddr string) []*connection {
	c.Lock()
	defer c.Unlock()

	addrs := c.listenAddrs[listenAddr]
	if len(addrs) == 0 {
		return nil
	}

	conns := make([]*connection, len(addrs))
	for i, a := range addrs {
		conns[i] = c.conns[a]
	}

	return conns
}

func (c *Connections) getByGnetID(gnetID uint64) *connection {
	c.Lock()
	defer c.Unlock()

	if gnetID == 0 {
		return nil
	}

	addr := c.gnetIDs[gnetID]
	if addr == "" {
		return nil
	}

	return c.conns[addr]
}

// modify modifies a connection.
// It is unsafe to modify the Mirror value with this method
func (c *Connections) modify(addr string, gnetID uint64, f func(c *ConnectionDetails)) error {
	conn := c.conns[addr]
	if conn == nil {
		return ErrConnectionNotExist
	}

	if conn.gnetID != gnetID {
		return ErrConnectionGnetIDMismatch
	}

	// copy and modify
	cd := conn.ConnectionDetails

	f(&cd)

	// compare to original
	if cd.Mirror != conn.ConnectionDetails.Mirror {
		logger.WithFields(logrus.Fields{
			"addr":   addr,
			"gnetID": gnetID,
		}).Panic("Connections.modify connection Mirror was changed")
	}

	if cd.ListenPort != conn.ConnectionDetails.ListenPort {
		logger.WithFields(logrus.Fields{
			"addr":   addr,
			"gnetID": gnetID,
		}).Panic("Connections.modify connection ListenPort was changed")
	}

	conn.ConnectionDetails = cd

	return nil
}

// SetHeight sets the height for a connection
func (c *Connections) SetHeight(addr string, gnetID uint64, height uint64) error {
	c.Lock()
	defer c.Unlock()

	return c.modify(addr, gnetID, func(c *ConnectionDetails) {
		c.Height = height
	})
}

func (c *Connections) updateMirror(ip string, mirror uint32, port uint16) error {
	x := c.mirrors[mirror]
	if x == nil {
		x = make(map[string]uint16, 2)
	}

	if _, ok := x[ip]; ok {
		return ErrConnectionIPMirrorExists
	}

	x[ip] = port
	c.mirrors[mirror] = x

	return nil
}

// canUpdateMirror returns false if a connection already exists with the same base IP and mirror value.
// This prevents duplicate connections to/from a single client.
func (c *Connections) canUpdateMirror(ip string, mirror uint32) error {
	x := c.mirrors[mirror]
	if x == nil {
		return nil
	}

	if _, ok := x[ip]; ok {
		return ErrConnectionIPMirrorExists
	}

	return nil
}

// IPCount returns the number of connections for a given base IP (without port)
func (c *Connections) IPCount(ip string) int {
	c.Lock()
	defer c.Unlock()
	return c.ipCounts[ip]
}

// Len returns number of connections
func (c *Connections) Len() int {
	c.Lock()
	defer c.Unlock()
	return len(c.conns)
}

// OutgoingLen returns number of outgoing connections
func (c *Connections) OutgoingLen() int {
	c.Lock()
	defer c.Unlock()
	n := 0
	for _, conn := range c.conns {
		if conn.Outgoing {
			n++
		}
	}
	return n
}

// PendingLen returns the number of status pending connections
func (c *Connections) PendingLen() int {
	c.Lock()
	defer c.Unlock()
	n := 0
	for _, conn := range c.conns {
		if conn.State == ConnectionStatePending {
			n++
		}
	}
	return n
}

// remove removes connection. Returns an error if the addr is invalid.
// If a connection with this address does not exist, nothing happens.
func (c *Connections) remove(addr string, gnetID uint64) error {
	c.Lock()
	defer c.Unlock()

	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Critical().WithError(err).Error("Connections.remove called with invalid addr")
		return err
	}

	conn := c.conns[addr]
	if conn == nil {
		return ErrConnectionNotExist
	}

	fields := logrus.Fields{
		"addr":       addr,
		"connGnetID": conn.gnetID,
		"gnetID":     gnetID,
		"listenPort": conn.ListenPort,
	}

	if conn.gnetID != gnetID {
		logger.Critical().WithFields(fields).Warning("Connections.remove gnetID does not match")
		return ErrConnectionGnetIDMismatch
	}

	x, ok := c.mirrors[conn.Mirror]
	if ok {
		if x[ip] != conn.ListenPort {
			logger.Critical().WithFields(fields).Warning("Indexed IP+Mirror value found but the ListenPort doesn't match")
		}

		delete(x, ip)
	}

	if len(x) == 0 {
		delete(c.mirrors, conn.Mirror)
	}

	if c.ipCounts[ip] > 0 {
		c.ipCounts[ip]--
	} else {
		logger.Critical().WithFields(fields).Warning("ipCount was already 0 when removing existing address")
	}

	listenAddr := conn.ListenAddr()
	if listenAddr != "" {
		addrs := c.listenAddrs[listenAddr]
		for i, a := range addrs {
			if a == conn.Addr {
				addrs = append(addrs[:i], addrs[i+1:]...)
				break
			}
		}
		if len(addrs) == 0 {
			delete(c.listenAddrs, listenAddr)
		} else {
			c.listenAddrs[listenAddr] = addrs
		}
	}

	delete(c.gnetIDs, conn.gnetID)
	delete(c.conns, addr)

	return nil
}

// all returns a copy of all connections
func (c *Connections) all() []connection {
	c.Lock()
	defer c.Unlock()

	conns := make([]connection, 0, len(c.conns))
	for _, c := range c.conns {
		conns = append(conns, *c)
	}

	return conns
}
