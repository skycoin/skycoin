package daemon

import (
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/util/iputil"
)

// ConnectionState connection state in the state machine
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
	// ErrConnectionAlreadyRegistered connection already registered in Connections
	ErrConnectionAlreadyRegistered = errors.New("Connection already registered")
	// ErrConnectionIPMirrorAlreadyRegistered connection already registered for a given base IP and mirror
	ErrConnectionIPMirrorAlreadyRegistered = errors.New("Connection already registered with this base IP and mirror")
	// ErrMirrorZero mirror value is 0
	ErrMirrorZero = errors.New("Mirror cannot be 0")
)

// Connection a connection's state within the daemon
type Connection struct {
	GnetID       int
	LastSent     time.Time
	LastReceived time.Time
	Addr         string
	ConnectionDetails
}

// ConnectionDetails extra connection data
type ConnectionDetails struct {
	State           ConnectionState
	Outgoing        bool
	ConnectedAt     time.Time
	Mirror          uint32
	ListenPort      uint16
	ProtocolVersion int32
	Height          uint64
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

func newConnection(c *connection) Connection {
	if c == nil {
		return Connection{}
	}

	conn := Connection{
		Addr:              c.addr,
		ConnectionDetails: c.ConnectionDetails,
	}

	if c.gnetConnection != nil {
		conn.GnetID = c.gnetConnection.ID
		conn.LastSent = c.gnetConnection.LastSent
		conn.LastReceived = c.gnetConnection.LastReceived
	}

	return conn
}

type connection struct {
	gnetConnection *gnet.Connection
	addr           string
	ConnectionDetails
}

// Connections manages a collection of Connection
type Connections struct {
	conns    map[string]*connection
	mirrors  map[uint32]map[string]uint16
	ipCounts map[string]int
	sync.Mutex
}

// NewConnections creates Connections
func NewConnections() *Connections {
	return &Connections{
		conns:    make(map[string]*connection, 32),
		mirrors:  make(map[uint32]map[string]uint16, 32),
		ipCounts: make(map[string]int, 32),
	}
}

// AddPendingOutgoing adds a new pending outgoing connection
func (c *Connections) AddPendingOutgoing(addr string) (Connection, error) {
	c.Lock()
	defer c.Unlock()

	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		return Connection{}, err
	}

	c.ipCounts[ip]++

	c.conns[addr] = &connection{
		addr: addr,
		ConnectionDetails: ConnectionDetails{
			State:    ConnectionStatePending,
			Outgoing: true,
		},
	}

	logger.WithField("addr", addr).Debug("AddPendingOutgoing")

	return newConnection(c.conns[addr]), nil
}

// Connected the connection has connected
func (c *Connections) Connected(gnetConn *gnet.Connection) (Connection, error) {
	c.Lock()
	defer c.Unlock()

	addr := gnetConn.Addr()

	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		return Connection{}, err
	}

	conn := c.conns[addr]

	if conn == nil {
		c.ipCounts[ip]++

		conn = &connection{
			addr: addr,
		}

		c.conns[addr] = conn
	} else {
		if addr != conn.addr {
			err := errors.New("gnet.Connection.Addr does not match recorded Connection address")
			logger.Critical().WithError(err).Error()
			return Connection{}, err
		}

		if conn.State != ConnectionStatePending {
			logger.Critical().WithField("state", conn.State).Warning("Transitioning to State connected but State is not pending")
		}
	}

	conn.gnetConnection = gnetConn
	conn.ConnectedAt = time.Now().UTC()
	conn.State = ConnectionStateConnected

	logger.WithFields(logrus.Fields{
		"addr":     addr,
		"outgoing": conn.Outgoing,
	}).Debug("Connected")

	return newConnection(conn), nil
}

// Introduced the connection has introduced itself
func (c *Connections) Introduced(addr string, m *IntroductionMessage) (Connection, error) {
	c.Lock()
	defer c.Unlock()

	ip, port, err := iputil.SplitAddr(addr)
	if err != nil {
		return Connection{}, err
	}

	if err := c.canUpdateMirror(ip, port, m.Mirror); err != nil {
		return Connection{}, err
	}

	conn := c.conns[addr]
	if conn == nil {
		return Connection{}, ErrConnectionNotExist
	}

	if err := c.updateMirror(ip, port, m.Mirror); err != nil {
		logger.WithError(err).Panic("updateMirror failed, but shouldn't")
	}

	conn.State = ConnectionStateIntroduced
	conn.Mirror = m.Mirror
	conn.ListenPort = m.Port
	conn.ProtocolVersion = m.Version

	logger.WithFields(logrus.Fields{
		"addr":     addr,
		"outgoing": conn.Outgoing,
	}).Debug("Introduced")

	return newConnection(conn), nil
}

// Get returns a connection by address
func (c *Connections) Get(addr string) (Connection, bool) {
	c.Lock()
	defer c.Unlock()

	conn, ok := c.conns[addr]
	return newConnection(conn), ok
}

// modify modifies a connection.
// It is unsafe to modify the Mirror value with this method
func (c *Connections) modify(addr string, f func(c *ConnectionDetails) error) error {
	conn := c.conns[addr]
	if conn == nil {
		return ErrConnectionNotExist
	}

	cd := conn.ConnectionDetails

	if err := f(&cd); err != nil {
		return err
	}

	if cd.Mirror != conn.ConnectionDetails.Mirror {
		logger.Panic("Connections.modify connection mirror value was changed")
	}

	conn.ConnectionDetails = cd

	return nil
}

// SetHeight sets the height for a connection
func (c *Connections) SetHeight(addr string, height uint64) error {
	c.Lock()
	defer c.Unlock()

	return c.modify(addr, func(c *ConnectionDetails) error {
		c.Height = height
		return nil
	})
}

// GetMirrorPort returns the port matching a given IP address (without port) and mirror value
func (c *Connections) GetMirrorPort(ip string, mirror uint32) uint16 {
	c.Lock()
	defer c.Unlock()

	x := c.mirrors[mirror]
	if x == nil {
		return 0
	}

	return x[ip]
}

func (c *Connections) updateMirror(ip string, port uint16, mirror uint32) error {
	if mirror == 0 {
		return ErrMirrorZero
	}

	x := c.mirrors[mirror]
	if x == nil {
		x = make(map[string]uint16, 2)
	}
	if _, ok := x[ip]; ok {
		return ErrConnectionIPMirrorAlreadyRegistered
	}
	x[ip] = port

	return nil
}

func (c *Connections) canUpdateMirror(ip string, port uint16, mirror uint32) error {
	if mirror == 0 {
		return ErrMirrorZero
	}

	x := c.mirrors[mirror]
	if x == nil {
		return nil
	}

	if _, ok := x[ip]; ok {
		return ErrConnectionIPMirrorAlreadyRegistered
	}

	return nil
}

// GetIPCount returns the number of connections for a given base IP (without port)
func (c *Connections) GetIPCount(ip string) int {
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

// Remove removes connection. Returns an error if the addr is invalid.
// If a connection with this address does not exist, nothing happens.
func (c *Connections) Remove(addr string) error {
	c.Lock()
	defer c.Unlock()
	return c.remove(addr)
}

func (c *Connections) remove(addr string) error {
	ip, port, err := iputil.SplitAddr(addr)
	if err != nil {
		return err
	}

	conn := c.conns[addr]
	if conn != nil {
		x, ok := c.mirrors[conn.Mirror]
		if ok {
			if x[ip] != port {
				logger.Critical().WithField("addr", addr).Warning("Indexed IP+Mirror value found but the port doesn't match")
			}

			delete(x, ip)
		}

		if len(x) == 0 {
			delete(c.mirrors, conn.Mirror)
		}

		if c.ipCounts[ip] > 0 {
			c.ipCounts[ip]--
		} else {
			logger.Critical().WithField("addr", addr).Warning("ipCount was already 0 when removing existing address")
		}
	}

	delete(c.conns, addr)

	return nil
}

// RemoveMatchedBy remove connections that match the matchFunc and return them
func (c *Connections) RemoveMatchedBy(f func(c Connection) (bool, error)) ([]string, error) {
	c.Lock()
	defer c.Unlock()

	var addrs []string
	for addr, conn := range c.conns {
		if ok, err := f(newConnection(conn)); err != nil {
			return nil, err
		} else if ok {
			addrs = append(addrs, addr)
		}
	}

	for _, a := range addrs {
		if err := c.remove(a); err != nil {
			logger.WithError(err).Panic("Invalid address stored inside Connections")
		}
	}

	return addrs, nil
}

// All returns a copy of all connections
func (c *Connections) All() []Connection {
	c.Lock()
	defer c.Unlock()

	conns := make([]Connection, 0, len(c.conns))
	for _, c := range c.conns {
		conns = append(conns, newConnection(c))
	}

	return conns
}
