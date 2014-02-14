package daemon

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
    "github.com/skycoin/skycoin/src/rpc"
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
    RemainingBalance Balance                   `json:"remaining_balance"`
    Transaction      visor.ReadableTransaction `json:"txn"`
    Error            string                    `json:"error"`
}

type BlockchainProgress struct {
    // Our current blockchain length
    Current uint64 `json:"current"`
    // Our best guess at true blockchain length
    Highest uint64 `json:"highest"`
}

type Balance struct {
    Balance rpc.Balance `json:"balance"`
    // Whether this balance includes unconfirmed txns in its calculation
    Predicted bool `json:"predicted"`
}

type Transaction struct {
    Transaction visor.ReadableTransaction `json:"txn"`
    Status      visor.TransactionStatus   `json:"status"`
}

type ResendResult struct {
    Sent bool `json:"sent"`
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

type Transactions struct {
    Txns []Transaction `json:"txns"`
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

// Returns a *Balance
func (self *RPC) GetTotalBalance(predicted bool) interface{} {
    self.requests <- func() interface{} {
        return self.getTotalBalance(predicted)
    }
    r := <-self.responses
    return r
}

// Returns a *Balance
func (self *RPC) GetBalance(a coin.Address, predicted bool) interface{} {
    self.requests <- func() interface{} {
        return self.getBalance(a, predicted)
    }
    r := <-self.responses
    return r
}

// Returns a *Spend
func (self *RPC) Spend(amt rpc.Balance, fee uint64, dest coin.Address) interface{} {
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

// Returns a *ReadableBlock
func (self *RPC) GetBlock(seq uint64) interface{} {
    self.requests <- func() interface{} { return self.getBlock(seq) }
    r := <-self.responses
    return r
}

// Returns a *ReadableBlocks
func (self *RPC) GetBlocks(start, end uint64) interface{} {
    self.requests <- func() interface{} { return self.getBlocks(start, end) }
    r := <-self.responses
    return r
}

// Returns a *BlockchainProgress
func (self *RPC) GetBlockchainProgress() interface{} {
    self.requests <- func() interface{} { return self.getBlockchainProgress() }
    r := <-self.responses
    return r
}

// Returns a *Transactions
func (self *RPC) GetAddressTransactions(a coin.Address) interface{} {
    self.requests <- func() interface{} {
        return self.getAddressTransactions(a)
    }
    r := <-self.responses
    return r
}

// Returns a *Transaction
func (self *RPC) GetTransaction(txn coin.SHA256) interface{} {
    self.requests <- func() interface{} {
        return self.getTransaction(txn)
    }
    r := <-self.responses
    return r
}

// Returns a *ResendResult
func (self *RPC) ResendTransaction(txn coin.SHA256) interface{} {
    self.requests <- func() interface{} {
        return self.resendTransaction(txn)
    }
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

func (self *RPC) getTotalBalance(predicted bool) *Balance {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    if predicted {
        return nil
    }
    var b rpc.Balance
    // if predicted {
    // b = self.Daemon.Visor.Visor.TotalBalancePredicted()
    // } else {
    b = self.Daemon.Visor.Visor.TotalBalance()
    // }
    return &Balance{
        Balance:   b,
        Predicted: predicted,
    }
}

func (self *RPC) getBalance(a coin.Address, predicted bool) *Balance {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    if predicted {
        // TODO -- prediction is disabled because implementation is not
        // clear
        return nil
    }
    var b rpc.Balance
    // if predicted {
    //     b = self.Daemon.Visor.rpc.BalancePredicted(a)
    // } else {
    b = self.Daemon.Visor.rpc.Balance(a)
    // }
    return &Balance{
        Balance:   b,
        Predicted: predicted,
    }
}

func (self *RPC) spend(amt rpc.Balance, fee uint64, dest coin.Address) *Spend {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    txn, err := self.Daemon.Visor.Spend(amt, fee, dest, self.Daemon.Pool)
    errString := ""
    if err != nil {
        errString = err.Error()
        logger.Error("Failed to make a spend: %v", err)
    }
    return &Spend{
        RemainingBalance: *(self.getTotalBalance(true)),
        Transaction:      visor.NewReadableTransaction(&txn),
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

func (self *RPC) getBlockchainProgress() *BlockchainProgress {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    return &BlockchainProgress{
        Current: self.Daemon.Visor.Visor.MostRecentBkSeq(),
        Highest: self.Daemon.Visor.EstimateBlockchainLength(),
    }
}

func (self *RPC) getTransaction(txHash coin.SHA256) *Transaction {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    txn := self.Daemon.Visor.Visor.GetTransaction(txHash)
    return &Transaction{
        Transaction: visor.NewReadableTransaction(&txn.Txn),
        Status:      txn.Status,
    }
}

func (self *RPC) getAddressTransactions(addr coin.Address) *Transactions {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    addrTxns := self.Daemon.Visor.Visor.GetAddressTransactions(addr)
    txns := make([]Transaction, 0, len(addrTxns))
    for _, tx := range addrTxns {
        txns = append(txns, Transaction{
            Transaction: visor.NewReadableTransaction(&tx.Txn),
            Status:      tx.Status,
        })
    }
    return &Transactions{
        Txns: txns,
    }
}

func (self *RPC) resendTransaction(txHash coin.SHA256) *ResendResult {
    if self.Daemon.Visor.Visor == nil {
        return nil
    }
    sent := self.Daemon.Visor.ResendTransaction(txHash, self.Daemon.Pool)
    return &ResendResult{
        Sent: sent,
    }
}
