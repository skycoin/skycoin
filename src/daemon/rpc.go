package daemon

// Exposes a read-only api for use by the gui rpc interface

type RPCConfig struct {
    BufferSize int
}

func NewRPCConfig() RPCConfig {
    return RPCConfig{
        BufferSize: 32,
    }
}

// RPC interface for daemon state
type RPC struct {
    // Backref to Daemon
    Daemon *Daemon
    Config RPCConfig

    // Requests are queued on this channel
    requests chan func() interface{}
    // When a request is done processing, it is placed on this channel
    responses chan interface{}
}

func NewRPC(c RPCConfig, d *Daemon) *RPC {
    return &RPC{
        Config:    c,
        Daemon:    d,
        requests:  make(chan func() interface{}, c.BufferSize),
        responses: make(chan interface{}, c.BufferSize),
    }
}

// A connection's state within the daemon
type Connection struct {
    Id           int    `json:"id"`
    Addr         string `json:"address"`
    LastSent     int64  `json:"last_sent"`
    LastReceived int64  `json:"last_received"`
    // Whether the connection is from us to them (true, outgoing),
    // or from them to us (false, incoming)
    Outgoing bool `json:"outgoing"`
    // Whether the client has identified their version, mirror etc
    Introduced bool   `json:"introduced"`
    Mirror     uint32 `json:"mirror"`
    ListenPort uint16 `json:"listen_port"`
}

// An array of connections
type Connections struct {
    Connections []*Connection `json:"connections"`
}

/* Public API
   Requests for data must be synchronized by the DaemonLoop
*/

// Returns a Connections struct
func (self *RPC) GetConnections() interface{} {
    self.requests <- func() interface{} { return self.getConnections() }
    r := <-self.responses
    return r
}

// Returns a Connection struct
func (self *RPC) GetConnection(addr string) interface{} {
    self.requests <- func() interface{} { return self.getConnection(addr) }
    r := <-self.responses
    return r
}

/* Internal API */

// Returns a Connection struct
func (self *RPC) getConnection(addr string) *Connection {
    if self.Daemon.Pool.Pool == nil {
        return nil
    }
    c := self.Daemon.Pool.Pool.Addresses[addr]
    _, expecting := self.Daemon.expectingIntroductions[addr]
    return &Connection{
        Id:           c.Id,
        Addr:         addr,
        LastSent:     c.LastSent.Unix(),
        LastReceived: c.LastReceived.Unix(),
        Outgoing:     (self.Daemon.outgoingConnections[addr] == nil),
        Introduced:   !expecting,
        Mirror:       self.Daemon.connectionMirrors[addr],
        ListenPort:   self.Daemon.getListenPort(addr),
    }
}

// Returns a Connections struct
func (self *RPC) getConnections() *Connections {
    if self.Daemon.Pool.Pool == nil {
        return nil
    }
    conns := make([]*Connection, 0, len(self.Daemon.Pool.Pool.Pool))
    for _, c := range self.Daemon.Pool.Pool.Pool {
        conns = append(conns, self.getConnection(c.Addr()))
    }
    return &Connections{Connections: conns}
}
