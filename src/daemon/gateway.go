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

func (gw *Gateway) doRequest(f func() interface{}) chan interface{} {
	req := makeRequest(f)
	gw.Requests <- req
	return req.Response
}

/* Daemon RPC wrappers */

/* Daemon internal status */

// Returns a *Connections
func (gw *Gateway) GetConnections() interface{} {
	rsp := gw.doRequest(func() interface{} {
		return gw.Daemon.GetConnections(gw.D)
	})

	return <-rsp
}

func (gw *Gateway) GetDefaultConnections() interface{} {
	rsp := gw.doRequest(func() interface{} {
		return gw.Daemon.GetDefaultConnections(gw.D)
	})
	return <-rsp
}

// Returns a *Connection
func (gw *Gateway) GetConnection(addr string) interface{} {
	rsp := gw.doRequest(func() interface{} {
		return gw.Daemon.GetConnection(gw.D, addr)
	})
	return <-rsp
}

/* Blockchain & Transaction status */
//DEPRECATE

// Returns a *BlockchainProgress
func (gw *Gateway) GetBlockchainProgress() interface{} {
	rsp := gw.doRequest(func() interface{} {
		return gw.Daemon.GetBlockchainProgress(gw.D.Visor)
	})
	return <-rsp
}

// Returns a *ResendResult
func (gw *Gateway) ResendTransaction(txn cipher.SHA256) interface{} {
	rsp := gw.doRequest(func() interface{} {
		return gw.Daemon.ResendTransaction(gw.D.Visor, gw.D.Pool, txn)
	})

	return <-rsp
}

// Returns a *visor.BlockchainMetadata
func (gw *Gateway) GetBlockchainMetadata() interface{} {
	rsp := gw.doRequest(func() interface{} {
		return gw.Visor.GetBlockchainMetadata(gw.V)
	})
	return <-rsp
}

// GetBlocks returns a *visor.ReadableBlocks
func (gw *Gateway) GetBlocks(start, end uint64) *visor.ReadableBlocks {
	rsp := gw.doRequest(func() interface{} {
		return gw.Visor.GetBlocks(gw.V, start, end)
	})
	v := <-rsp
	return v.(*visor.ReadableBlocks)
}

// GetLastBlocks get last N blocks
func (gw *Gateway) GetLastBlocks(num uint64) *visor.ReadableBlocks {
	rsp := gw.doRequest(func() interface{} {
		headSeq := gw.V.HeadBkSeq()
		var start uint64
		if (headSeq + 1) > num {
			start = headSeq - num + 1
		}

		blocks := gw.Visor.GetBlocks(gw.V, start, headSeq)
		return blocks
	})
	v := <-rsp
	return v.(*visor.ReadableBlocks)
}

// Returns a *visor.TransactionResult
func (gw *Gateway) GetTransaction(txn cipher.SHA256) interface{} {
	rsp := gw.doRequest(func() interface{} {
		return gw.Visor.GetTransaction(gw.V, txn)
	})
	return <-rsp
}

// Returns a *visor.TransactionResults
func (gw *Gateway) GetAddressTransactions(a cipher.Address) interface{} {
	rsp := gw.doRequest(func() interface{} {
		return gw.Visor.GetAddressTransactions(gw.V, a)
	})
	return <-rsp
}
