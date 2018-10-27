package daemon

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/util/iputil"
	"github.com/skycoin/skycoin/src/util/useragent"
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
	// ErrConnectionAlreadyRegistered connection already registered in Connections
	ErrConnectionAlreadyRegistered = errors.New("Connection already registered")
	// ErrConnectionIPMirrorAlreadyRegistered connection already registered for a given base IP and mirror
	ErrConnectionIPMirrorAlreadyRegistered = errors.New("Connection already registered with this base IP and mirror")
)

// ConnectionDetails connection data managed by daemon
type ConnectionDetails struct {
	State           ConnectionState
	Outgoing        bool
	ConnectedAt     time.Time
	Mirror          uint32
	ListenPort      uint16
	ProtocolVersion int32
	Height          uint64
	UserAgent       useragent.Data
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

// pending adds a new pending outgoing connection
func (c *Connections) pending(addr string) (*connection, error) {
	c.Lock()
	defer c.Unlock()

	ip, port, err := iputil.SplitAddr(addr)
	if err != nil {
		return nil, err
	}

	if _, ok := c.conns[addr]; ok {
		return nil, ErrConnectionAlreadyRegistered
	}

	c.ipCounts[ip]++

	c.conns[addr] = &connection{
		Addr: addr,
		ConnectionDetails: ConnectionDetails{
			State:      ConnectionStatePending,
			Outgoing:   true,
			ListenPort: port,
		},
	}

	logger.WithField("addr", addr).Debug("Connections.pending")

	return c.conns[addr], nil
}

// connected the connection has connected
func (c *Connections) connected(addr string) (*connection, error) {
	c.Lock()
	defer c.Unlock()

	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
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
		if addr != conn.Addr {
			err := errors.New("gnet.Connection.Addr does not match recorded Connection address")
			logger.Critical().WithError(err).Error()
			return nil, err
		}

		if conn.State != ConnectionStatePending {
			logger.Critical().WithField("state", conn.State).Warningf("Transitioning to State %q but State is not %q", ConnectionStateConnected, ConnectionStatePending)
		}
	}

	conn.ConnectedAt = time.Now().UTC()
	conn.State = ConnectionStateConnected

	logger.WithFields(logrus.Fields{
		"addr":     addr,
		"outgoing": conn.Outgoing,
	}).Debug("Connections.connected")

	return conn, nil
}

// introduced the connection has introduced itself
func (c *Connections) introduced(addr string, m *IntroductionMessage) (*connection, error) {
	c.Lock()
	defer c.Unlock()

	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		return nil, err
	}

	conn := c.conns[addr]
	if conn == nil {
		return nil, ErrConnectionNotExist
	}

	if conn.State != ConnectionStateConnected {
		logger.Critical().WithFields(logrus.Fields{
			"addr":  conn.Addr,
			"state": conn.State,
		}).Warningf("Transitioning to State %q but State is not %q", ConnectionStateIntroduced, ConnectionStateConnected)
	}

	if err := c.canUpdateMirror(ip, m.Mirror); err != nil {
		return nil, err
	}

	// For outgoing connections, which are created by pending,
	// the listen port is set from the addr's port number.
	// Since we are connecting to it, it is presumed to be that peer's open listening port.
	// A misbehaving peer could report a different ListenPort in their IntroductionMessage,
	// but it shouldn't affect our records.
	if conn.Outgoing && conn.ListenPort != m.ListenPort {
		logger.Critical().WithFields(logrus.Fields{
			"addr":              conn.Addr,
			"connListenPort":    conn.ListenPort,
			"messageListenPort": m.ListenPort,
		}).Warning("Outgoing connection's ListenPort does not match reported IntroductionMessage ListenPort")
	}

	listenPort := conn.ListenPort
	if !conn.Outgoing {
		listenPort = m.ListenPort
	}

	if err := c.updateMirror(ip, m.Mirror, listenPort); err != nil {
		logger.WithError(err).Panic("updateMirror failed, but shouldn't")
	}

	conn.State = ConnectionStateIntroduced
	conn.Mirror = m.Mirror
	conn.ProtocolVersion = m.ProtocolVersion
	conn.ListenPort = listenPort
	conn.UserAgent = m.userAgentData

	logger.WithFields(logrus.Fields{
		"addr":     addr,
		"outgoing": conn.Outgoing,
	}).Debug("Connections.introduced")

	return conn, nil
}

// get returns a connection by address
func (c *Connections) get(addr string) *connection {
	c.Lock()
	defer c.Unlock()

	return c.conns[addr]
}

// modify modifies a connection.
// It is unsafe to modify the Mirror value with this method
func (c *Connections) modify(addr string, f func(c *ConnectionDetails)) error {
	conn := c.conns[addr]
	if conn == nil {
		return ErrConnectionNotExist
	}

	// copy and modify
	cd := conn.ConnectionDetails

	f(&cd)

	// compare to original
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

	return c.modify(addr, func(c *ConnectionDetails) {
		c.Height = height
	})
}

func (c *Connections) updateMirror(ip string, mirror uint32, port uint16) error {
	logger.Debugf("updateMirror ip=%s mirror=%d port=%d", ip, mirror, port)

	x := c.mirrors[mirror]
	if x == nil {
		x = make(map[string]uint16, 2)
	}

	if _, ok := x[ip]; ok {
		return ErrConnectionIPMirrorAlreadyRegistered
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
		return ErrConnectionIPMirrorAlreadyRegistered
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
func (c *Connections) remove(addr string) error {
	c.Lock()
	defer c.Unlock()

	ip, port, err := iputil.SplitAddr(addr)
	if err != nil {
		return err
	}

	conn := c.conns[addr]
	if conn != nil {
		x, ok := c.mirrors[conn.Mirror]
		if ok {
			if x[ip] != conn.ListenPort {
				logger.Critical().WithField("addr", addr).Warning("Indexed IP+Mirror value found but the ListenPort doesn't match")
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
