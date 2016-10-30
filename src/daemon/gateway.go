package daemon

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	//"github.com/skycoin/skycoin/src/wallet"
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

type Request struct {
	Handle   func() interface{}
	Response chan interface{}
}

func makeRequest(f func() interface{}) Request {
	return Request{
		Handle:   f,
		Response: make(chan interface{}),
	}
}

// RPC interface wrapper for daemon state
type Gateway struct {
	Config GatewayConfig
	Daemon RPC
	Visor  visor.RPC

	// Backref to Daemon
	D *Daemon
	// Backref to Visor
	V *visor.Visor
	// Requests are queued on this channel
	Requests chan Request
	// When a request is done processing, it is placed on this channel
	// Responses chan interface{}
}

func NewGateway(c GatewayConfig, D *Daemon) *Gateway {
	return &Gateway{
		Config:   c,
		Daemon:   RPC{},
		Visor:    visor.RPC{},
		D:        D,
		V:        D.Visor.Visor,
		Requests: make(chan Request, c.BufferSize),
		// Responses: make(chan interface{}, c.BufferSize),
	}
}

func (self *Gateway) doRequest(f func() interface{}) chan interface{} {
	req := makeRequest(f)
	self.Requests <- req
	return req.Response
}

/* Daemon RPC wrappers */

/* Daemon internal status */

// Returns a *Connections
func (self *Gateway) GetConnections() interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Daemon.GetConnections(self.D)
	})

	return <-rsp
}

func (self *Gateway) GetDefaultConnections() interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Daemon.GetDefaultConnections(self.D)
	})
	return <-rsp
}

// Returns a *Connection
func (self *Gateway) GetConnection(addr string) interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Daemon.GetConnection(self.D, addr)
	})
	return <-rsp
}

/* Blockchain & Transaction status */
//DEPRECATE

// Returns a *BlockchainProgress
func (self *Gateway) GetBlockchainProgress() interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Daemon.GetBlockchainProgress(self.D.Visor)
	})
	return <-rsp
}

// Returns a *ResendResult
func (self *Gateway) ResendTransaction(txn cipher.SHA256) interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Daemon.ResendTransaction(self.D.Visor, self.D.Pool, txn)
	})

	return <-rsp
}

// Returns a *visor.BlockchainMetadata
func (self *Gateway) GetBlockchainMetadata() interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Visor.GetBlockchainMetadata(self.V)
	})
	return <-rsp
}

// Returns a *visor.ReadableBlocks
func (self *Gateway) GetBlocks(start, end uint64) interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Visor.GetBlocks(self.V, start, end)
	})
	return <-rsp
}

// GetLastBlocks get last N blocks
func (self *Gateway) GetLastBlocks(num uint64) interface{} {
	rsp := self.doRequest(func() interface{} {
		headSeq := self.V.HeadBkSeq()
		var start uint64
		if (headSeq + 1) > num {
			start = headSeq - num + 1
		}

		blocks := self.V.GetBlocks(start, headSeq)
		return blocks
	})
	return <-rsp
}

// Returns a *visor.TransactionResult
func (self *Gateway) GetTransaction(txn cipher.SHA256) interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Visor.GetTransaction(self.V, txn)
	})
	return <-rsp
}

// Returns a *visor.TransactionResults
func (self *Gateway) GetAddressTransactions(a cipher.Address) interface{} {
	rsp := self.doRequest(func() interface{} {
		return self.Visor.GetAddressTransactions(self.V, a)
	})
	return <-rsp
}
