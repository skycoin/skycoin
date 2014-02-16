package daemon

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
)

// Exposes a read-only api for use by the gui rpc interface

type GatewayConfig struct {
    BufferSize int
}

func NewGatewayConfig() GatewayConfig {
    return GatewayConfig{
        BufferSize: 32,
    }
}

// RPC interface wrapper for daemon state
type Gateway struct {
    Config GatewayConfig
    Daemon RPC
    Visor  visor.RPC

    // Backref to Daemon
    d   *Daemon
    // Requests are queued on this channel
    requests chan func() interface{}
    // When a request is done processing, it is placed on this channel
    responses chan interface{}
}

func NewGateway(c GatewayConfig, d *Daemon) *Gateway {
    return &Gateway{
        Config:    c,
        Daemon:    RPC{},
        Visor:     visor.RPC{},
        d:         d,
        requests:  make(chan func() interface{}, c.BufferSize),
        responses: make(chan interface{}, c.BufferSize),
    }
}

/* Daemon RPC wrappers */

// Returns a *Connections
func (self *Gateway) GetConnections() interface{} {
    self.requests <- func() interface{} {
        return self.Daemon.GetConnections(self.d)
    }
    r := <-self.responses
    return r
}

// Returns a *Connection
func (self *Gateway) GetConnection(addr string) interface{} {
    self.requests <- func() interface{} {
        return self.Daemon.GetConnection(self.d, addr)
    }
    r := <-self.responses
    return r
}

// Returns a *Spend
func (self *Gateway) Spend(amt visor.Balance, fee uint64,
    dest coin.Address) interface{} {
    self.requests <- func() interface{} {
        return self.Daemon.Spend(self.d.Visor, self.d.Pool, self.Visor, amt,
            fee, dest)
    }
    r := <-self.responses
    return r
}

// Returns a *BlockchainProgress
func (self *Gateway) GetBlockchainProgress() interface{} {
    self.requests <- func() interface{} {
        return self.Daemon.GetBlockchainProgress(self.d.Visor)
    }
    r := <-self.responses
    return r
}

// Returns a *ResendResult
func (self *Gateway) ResendTransaction(txn coin.SHA256) interface{} {
    self.requests <- func() interface{} {
        return self.Daemon.ResendTransaction(self.d.Visor, self.d.Pool, txn)
    }
    r := <-self.responses
    return r
}

/* Visor RPC wrappers */

// Returns a *Balance
func (self *Gateway) GetTotalBalance(predicted bool) interface{} {
    self.requests <- func() interface{} {
        return self.Visor.GetTotalBalance(self.d.Visor.Visor, predicted)
    }
    r := <-self.responses
    return r
}

// Returns a *Balance
func (self *Gateway) GetBalance(a coin.Address, predicted bool) interface{} {
    self.requests <- func() interface{} {
        return self.Visor.GetBalance(self.d.Visor.Visor, a, predicted)
    }
    r := <-self.responses
    return r
}

// Returns an error
func (self *Gateway) SaveWallet() interface{} {
    self.requests <- func() interface{} {
        return self.Visor.SaveWallet(self.d.Visor.Visor)
    }
    r := <-self.responses
    return r
}

// Returns a *visor.ReadableWalletEntry
func (self *Gateway) CreateAddress() interface{} {
    self.requests <- func() interface{} {
        return self.Visor.CreateAddress(self.d.Visor.Visor)
    }
    r := <-self.responses
    return r
}

// Returns a *visor.ReadableWallet
func (self *Gateway) GetWallet() interface{} {
    self.requests <- func() interface{} {
        return self.Visor.GetWallet(self.d.Visor.Visor)
    }
    r := <-self.responses
    return r
}

// Returns a *visor.BlockchainMetadata
func (self *Gateway) GetBlockchainMetadata() interface{} {
    self.requests <- func() interface{} {
        return self.Visor.GetBlockchainMetadata(self.d.Visor.Visor)
    }
    r := <-self.responses
    return r
}

// Returns a *visor.ReadableBlock
func (self *Gateway) GetBlock(seq uint64) interface{} {
    self.requests <- func() interface{} {
        return self.Visor.GetBlock(self.d.Visor.Visor, seq)
    }
    r := <-self.responses
    return r
}

// Returns a *visor.ReadableBlocks
func (self *Gateway) GetBlocks(start, end uint64) interface{} {
    self.requests <- func() interface{} {
        return self.Visor.GetBlocks(self.d.Visor.Visor, start, end)
    }
    r := <-self.responses
    return r
}

// Returns a *visor.TransactionResult
func (self *Gateway) GetTransaction(txn coin.SHA256) interface{} {
    self.requests <- func() interface{} {
        return self.Visor.GetTransaction(self.d.Visor.Visor, txn)
    }
    r := <-self.responses
    return r
}

// Returns a *visor.TransactionResults
func (self *Gateway) GetAddressTransactions(a coin.Address) interface{} {
    self.requests <- func() interface{} {
        return self.Visor.GetAddressTransactions(self.d.Visor.Visor, a)
    }
    r := <-self.responses
    return r
}
