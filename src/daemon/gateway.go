package daemon

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	//"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

// Exposes a read-only api for use by the gui rpc interface

// GatewayConfig configuration set of gateway.
type GatewayConfig struct {
	BufferSize int
}

// NewGatewayConfig create and init an GatewayConfig
func NewGatewayConfig() GatewayConfig {
	return GatewayConfig{
		BufferSize: 32,
	}
}

// Gateway RPC interface wrapper for daemon state
type Gateway struct {
	Config GatewayConfig
	Daemon RPC
	Visor  visor.RPC

	// Backref to Daemon
	D *Daemon
	// Backref to Visor
	V *visor.Visor
	// Requests are queued on this channel
	Requests chan func()
	// When a request is done processing, it is placed on this channel
	// Responses chan interface{}
}

// NewGateway create and init an Gateway instance.
func NewGateway(c GatewayConfig, D *Daemon) *Gateway {
	return &Gateway{
		Config:   c,
		Daemon:   RPC{},
		Visor:    visor.RPC{},
		D:        D,
		V:        D.Visor.Visor,
		Requests: make(chan func(), c.BufferSize),
	}
}

// GetConnections returns a *Connections
func (gw *Gateway) GetConnections() interface{} {
	conns := make(chan interface{})
	gw.Requests <- func() {
		conns <- gw.Daemon.GetConnections(gw.D)
	}

	return <-conns
}

// GetDefaultConnections returns default connections
func (gw *Gateway) GetDefaultConnections() interface{} {
	conns := make(chan interface{})
	gw.Requests <- func() {
		conns <- gw.Daemon.GetDefaultConnections(gw.D)
	}
	return <-conns
}

// GetConnection returns a *Connection of specific address
func (gw *Gateway) GetConnection(addr string) interface{} {
	conn := make(chan interface{})
	gw.Requests <- func() {
		conn <- gw.Daemon.GetConnection(gw.D, addr)
	}
	return <-conn
}

/* Blockchain & Transaction status */
//DEPRECATE

// GetBlockchainProgress returns a *BlockchainProgress
func (gw *Gateway) GetBlockchainProgress() interface{} {
	bcp := make(chan interface{})
	gw.Requests <- func() {
		bcp <- gw.Daemon.GetBlockchainProgress(gw.D.Visor)
	}
	return <-bcp
}

// ResendTransaction resent the transaction and return a *ResendResult
func (gw *Gateway) ResendTransaction(txn cipher.SHA256) interface{} {
	result := make(chan interface{})
	gw.Requests <- func() {
		result <- gw.Daemon.ResendTransaction(gw.D.Visor, gw.D.Pool, txn)
	}

	return <-result
}

// GetBlockchainMetadata returns a *visor.BlockchainMetadata
func (gw *Gateway) GetBlockchainMetadata() interface{} {
	bcm := make(chan interface{})
	gw.Requests <- func() {
		bcm <- gw.Visor.GetBlockchainMetadata(gw.V)
	}
	return <-bcm
}

// GetBlocks returns a *visor.ReadableBlocks
func (gw *Gateway) GetBlocks(start, end uint64) *visor.ReadableBlocks {
	blocks := make(chan *visor.ReadableBlocks)
	gw.Requests <- func() {
		blocks <- gw.Visor.GetBlocks(gw.V, start, end)
	}
	return <-blocks
}

// GetLastBlocks get last N blocks
func (gw *Gateway) GetLastBlocks(num uint64) *visor.ReadableBlocks {
	blocks := make(chan *visor.ReadableBlocks)
	gw.Requests <- func() {
		headSeq := gw.V.HeadBkSeq()
		var start uint64
		if (headSeq + 1) > num {
			start = headSeq - num + 1
		}

		blocks <- gw.Visor.GetBlocks(gw.V, start, headSeq)
	}
	return <-blocks
}

// GetUnspentByAddrs gets unspent of specific addresses
func (gw *Gateway) GetUnspentByAddrs(addrs []string) []visor.ReadableOutput {
	outputs := make(chan []visor.ReadableOutput)
	gw.Requests <- func() {
		outs := gw.V.GetUnspentOutputReadables()
		addrMatch := []visor.ReadableOutput{}
		addrMap := make(map[string]bool)
		for _, addr := range addrs {
			addrMap[addr] = true
		}

		for _, u := range outs {
			if _, ok := addrMap[u.Address]; ok {
				addrMatch = append(addrMatch, u)
			}
		}
		outputs <- addrMatch
	}
	return <-outputs
}

// GetUnspentByHashes gets unspent of specific unspent hashes.
func (gw *Gateway) GetUnspentByHashes(hashes []string) []visor.ReadableOutput {
	outputs := make(chan []visor.ReadableOutput)
	gw.Requests <- func() {
		outs := gw.V.GetUnspentOutputReadables()

		hsMatch := []visor.ReadableOutput{}
		hsMap := make(map[string]bool)
		for _, h := range hashes {
			hsMap[h] = true
		}

		for _, u := range outs {
			if _, ok := hsMap[u.Hash]; ok {
				hsMatch = append(hsMatch, u)
			}
		}
		outputs <- hsMatch
	}
	return <-outputs
}

// GetTransaction gets transaction by txid.
func (gw *Gateway) GetTransaction(txid cipher.SHA256) (*visor.TransactionResult, error) {
	var tx *visor.TransactionResult
	var err error
	c := make(chan struct{})
	gw.Requests <- func() {
		tx, err = gw.Visor.GetTransaction(gw.V, txid)
		c <- struct{}{}
	}
	<-c
	return tx, err
}

// InjectTransaction injects transaction
func (gw *Gateway) InjectTransaction(txn coin.Transaction) (coin.Transaction, error) {
	var tx coin.Transaction
	var err error
	c := make(chan struct{})
	gw.Requests <- func() {
		tx, err = gw.D.Visor.InjectTransaction(txn, gw.D.Pool)
		c <- struct{}{}
	}
	<-c
	return tx, err
}

// GetAddressTransactions returns a *visor.TransactionResults
func (gw *Gateway) GetAddressTransactions(a cipher.Address) interface{} {
	tx := make(chan interface{})
	gw.Requests <- func() {
		tx <- gw.Visor.GetAddressTransactions(gw.V, a)
	}
	return <-tx
}

// GetUxOutByID gets UxOut by hash id.
func (gw *Gateway) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	var uxout *historydb.UxOut
	var err error
	c := make(chan struct{})
	gw.Requests <- func() {
		uxout, err = gw.V.GetUxOutByID(id)
		c <- struct{}{}
	}
	<-c
	return uxout, err
}

func (gw *Gateway) GetRecvUxOutOfAddr(addr cipher.Address) ([]*historydb.UxOut, error) {
	var (
		uxouts []*historydb.UxOut
		err    error
	)
	c := make(chan struct{})
	gw.Requests <- func() {
		uxouts, err = gw.V.GetRecvUxOutOfAddr(addr)
		c <- struct{}{}
	}
	<-c
	return uxouts, err
}

func (gw *Gateway) GetSpentUxOutOfAddr(addr cipher.Address) ([]*historydb.UxOut, error) {
	var (
		outputs []*historydb.UxOut
		err     error
	)
	c := make(chan struct{})
	gw.Requests <- func() {
		outputs, err = gw.V.GetSpentUxOutOfAddr(addr)
		c <- struct{}{}
	}
	<-c
	return outputs, err
}
