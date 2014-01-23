package daemon

// Exposes a read-only api for use by the gui rpc interface

var (
    apiBufferSize = 32
    apiRequests   = make(chan func() interface{}, apiBufferSize)
    apiResponses  = make(chan interface{}, apiBufferSize)
)

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
func GetConnections() interface{} {
    apiRequests <- func() interface{} { return getConnections() }
    r := <-apiResponses
    return r
}

// Returns a Connection struct
func GetConnection(addr string) interface{} {
    apiRequests <- func() interface{} { return getConnection(addr) }
    r := <-apiResponses
    return r
}

/* Internal API */

// Returns a Connection struct
func getConnection(addr string) *Connection {
    if Pool == nil {
        return nil
    }
    c := Pool.Addresses[addr]
    _, expecting := expectingIntroductions[addr]
    return &Connection{
        Id:           c.Id,
        Addr:         addr,
        LastSent:     c.LastSent.Unix(),
        LastReceived: c.LastReceived.Unix(),
        Outgoing:     (outgoingConnections[addr] == nil),
        Introduced:   !expecting,
        Mirror:       connectionMirrors[addr],
        ListenPort:   getListenPort(addr),
    }
}

// Returns a Connections struct
func getConnections() *Connections {
    if Pool == nil {
        return nil
    }
    conns := make([]*Connection, 0, len(Pool.Pool))
    for _, c := range Pool.Pool {
        conns = append(conns, getConnection(c.Addr()))
    }
    return &Connections{Connections: conns}
}
