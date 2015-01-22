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
	Requests chan func() interface{}
	// When a request is done processing, it is placed on this channel
	Responses chan interface{}
}

func NewGateway(c GatewayConfig, D *Daemon) *Gateway {
	return &Gateway{
		Config:    c,
		Daemon:    RPC{},
		Visor:     visor.RPC{},
		D:         D,
		V:         D.Visor.Visor,
		Requests:  make(chan func() interface{}, c.BufferSize),
		Responses: make(chan interface{}, c.BufferSize),
	}
}

/* Daemon RPC wrappers */

/* Daemon internal status */

// Returns a *Connections
func (self *Gateway) GetConnections() interface{} {
	self.Requests <- func() interface{} {
		return self.Daemon.GetConnections(self.D)
	}
	r := <-self.Responses
	return r
}

// Returns a *Connection
func (self *Gateway) GetConnection(addr string) interface{} {
	self.Requests <- func() interface{} {
		return self.Daemon.GetConnection(self.D, addr)
	}
	r := <-self.Responses
	return r
}

/* Wallet API */

// Returns a *Spend
/*
func (self *Gateway) Spend(walletID wallet.WalletID, amt wallet.Balance,
	fee uint64, dest cipher.Address) interface{} {
	self.Requests <- func() interface{} {
		return self.Daemon.Spend(self.D.Visor, self.D.Pool, self.Visor,
			walletID, amt, fee, dest)
	}
	r := <-self.Responses
	return r
}
*/

// Returns a *Balance
/*
func (self *Gateway) GetWalletBalance(walletID wallet.WalletID) interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.GetWalletBalance(self.V, walletID)
	}
	r := <-self.Responses
	return r
}
*/

// Returns map[WalletID]error
/*
func (self *Gateway) SaveWallets() interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.SaveWallets(self.V)
	}
	r := <-self.Responses
	return r
}

// Returns error
func (self *Gateway) SaveWallet(walletID wallet.WalletID) interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.SaveWallet(self.V, walletID)
	}
	r := <-self.Responses
	return r
}

// Returns an error
func (self *Gateway) ReloadWallets() interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.ReloadWallets(self.V)
	}
	r := <-self.Responses
	return r
}
*/

// Returns a *visor.ReadableWallet
/*
func (self *Gateway) GetWallet(walletID wallet.WalletID) interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.GetWallet(self.V, walletID)
	}
	r := <-self.Responses
	return r
}
*/

// Returns a *ReadableWallets
/*
func (self *Gateway) GetWallets() interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.GetWallets(self.V)
	}
	r := <-self.Responses
	return r
}
*/

// Returns a *ReadableWallet
// Deprecate
/*
func (self *Gateway) CreateWallet(seed string) interface{} {

	//w := v.CreateWallet()
	return wallet.NewReadableWallet(w)

	//
		self.Requests <- func() interface{} {
			return self.Visor.CreateWallet(self.V)
		}
		r := <-self.Responses
		return r
	//
}
*/
/* Blockchain & Transaction status */

// Returns a *BlockchainProgress
func (self *Gateway) GetBlockchainProgress() interface{} {
	self.Requests <- func() interface{} {
		return self.Daemon.GetBlockchainProgress(self.D.Visor)
	}
	r := <-self.Responses
	return r
}

// Returns a *ResendResult
func (self *Gateway) ResendTransaction(txn cipher.SHA256) interface{} {
	self.Requests <- func() interface{} {
		return self.Daemon.ResendTransaction(self.D.Visor, self.D.Pool, txn)
	}
	r := <-self.Responses
	return r
}

// Returns a *visor.BlockchainMetadata
func (self *Gateway) GetBlockchainMetadata() interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.GetBlockchainMetadata(self.V)
	}
	r := <-self.Responses
	return r
}

// Returns a *visor.ReadableBlock
func (self *Gateway) GetBlock(seq uint64) interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.GetBlock(self.V, seq)
	}
	r := <-self.Responses
	return r
}

// Returns a *visor.ReadableBlocks
func (self *Gateway) GetBlocks(start, end uint64) interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.GetBlocks(self.V, start, end)
	}
	r := <-self.Responses
	return r
}

// Returns a *visor.TransactionResult
func (self *Gateway) GetTransaction(txn cipher.SHA256) interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.GetTransaction(self.V, txn)
	}
	r := <-self.Responses
	return r
}

// Returns a *visor.TransactionResults
func (self *Gateway) GetAddressTransactions(a cipher.Address) interface{} {
	self.Requests <- func() interface{} {
		return self.Visor.GetAddressTransactions(self.V, a)
	}
	r := <-self.Responses
	return r
}
