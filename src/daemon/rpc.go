package daemon

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
)

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

// Result of a Spend() operation
type Spend struct {
    RemainingBalance visor.Balance
    Error            string
}

// Arrays must be wrapped in structs to avoid certain javascript exploits

// An array of connections
type Connections struct {
    Connections []*Connection `json:"connections"`
}

// An array of readable blocks.
type ReadableBlocks struct {
    Blocks []visor.ReadableBlock `json:"blocks"`
}

/* Public API
   Requests for data must be synchronized by the DaemonLoop
*/

// Returns a *Connections
func (self *RPC) GetConnections() interface{} {
    self.requests <- func() interface{} { return self.getConnections() }
    r := <-self.responses
    return r
}

// Returns a *Connection
func (self *RPC) GetConnection(addr string) interface{} {
    self.requests <- func() interface{} { return self.getConnection(addr) }
    r := <-self.responses
    return r
}

// Returns a *coin.Balance
func (self *RPC) GetTotalBalance() interface{} {
    self.requests <- func() interface{} { return self.getTotalBalance() }
    r := <-self.responses
    return r
}

// Returns a *coin.Balance
func (self *RPC) GetBalance(a coin.Address) interface{} {
    self.requests <- func() interface{} { return self.getBalance(a) }
    r := <-self.responses
    return r
}

// Returns a *Spend
func (self *RPC) Spend(amt visor.Balance, fee uint64, dest coin.Address) interface{} {
    self.requests <- func() interface{} { return self.spend(amt, fee, dest) }
    r := <-self.responses
    return r
}

// Returns an error
func (self *RPC) SaveWallet() interface{} {
    self.requests <- func() interface{} { return self.saveWallet() }
    r := <-self.responses
    return r
}

// Returns a *visor.ReadableWalletEntry
func (self *RPC) CreateAddress() interface{} {
    self.requests <- func() interface{} { return self.createAddress() }
    r := <-self.responses
    return r
}

// Returns a *visor.ReadableWallet
func (self *RPC) GetWallet() interface{} {
    self.requests <- func() interface{} { return self.getWallet() }
    r := <-self.responses
    return r
}

// Returns a *visor.BlockchainMetadata
func (self *RPC) GetBlockchainMetadata() interface{} {
    self.requests <- func() interface{} { return self.getBlockchainMetadata() }
    r := <-self.responses
    return r
}

// Returns a *coin.Block
func (self *RPC) GetBlock(seq uint64) interface{} {
    self.requests <- func() interface{} { return self.getBlock(seq) }
    r := <-self.responses
    return r
}

// Returns a []coin.Block
func (self *RPC) GetBlocks(start, end uint64) interface{} {
    self.requests <- func() interface{} { return self.getBlocks(start, end) }
    r := <-self.responses
    return r
}

/* Internal API */

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

func (self *RPC) getTotalBalance() *visor.Balance {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    b := self.Daemon.Visor.Visor.TotalBalancePredicted()
    return &b
}

func (self *RPC) getBalance(a coin.Address) *visor.Balance {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    b := self.Daemon.Visor.Visor.BalancePredicted(a)
    return &b
}

func (self *RPC) spend(amt visor.Balance, fee uint64, dest coin.Address) *Spend {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    _, err := self.Daemon.Visor.Spend(amt, fee, dest, self.Daemon.Pool)
    errString := ""
    if err != nil {
        errString = err.Error()
        logger.Error("Failed to make a spend: %v", err)
    }
    return &Spend{
        RemainingBalance: *(self.getTotalBalance()),
        Error:            errString,
    }
}

func (self *RPC) saveWallet() error {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    return self.Daemon.Visor.Visor.SaveWallet()
}

func (self *RPC) createAddress() *visor.ReadableWalletEntry {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    we, err := self.Daemon.Visor.Visor.CreateAddressAndSave()
    if err != nil {
        return nil
    }
    rwe := visor.NewReadableWalletEntry(&we)
    return &rwe
}

func (self *RPC) getWallet() *visor.ReadableWallet {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    return visor.NewReadableWallet(self.Daemon.Visor.Visor.Wallet)
}

func (self *RPC) getBlockchainMetadata() *visor.BlockchainMetadata {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    bm := self.Daemon.Visor.Visor.GetBlockchainMetadata()
    return &bm
}

func (self *RPC) getBlock(seq uint64) *visor.ReadableBlock {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    b, err := self.Daemon.Visor.Visor.GetReadableBlock(seq)
    if err != nil {
        return nil
    }
    return &b
}

func (self *RPC) getBlocks(start, end uint64) *ReadableBlocks {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    blocks := self.Daemon.Visor.Visor.GetReadableBlocks(start, end)
    return &ReadableBlocks{blocks}
}
