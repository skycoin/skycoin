package daemon

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
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
	d *Daemon
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

/* Daemon internal status */

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

/* Wallet API */

// Returns a *Spend
func (self *Gateway) Spend(walletID wallet.WalletID, amt wallet.Balance,
	fee uint64, dest cipher.Address) interface{} {
	self.requests <- func() interface{} {
		return self.Daemon.Spend(self.d.Visor, self.d.Pool, self.Visor,
			walletID, amt, fee, dest)
	}
	r := <-self.responses
	return r
}

// Returns a *Balance
/*
func (self *Gateway) GetWalletBalance(walletID wallet.WalletID) interface{} {
	self.requests <- func() interface{} {
		return self.Visor.GetWalletBalance(self.d.Visor.Visor, walletID)
	}
	r := <-self.responses
	return r
}
*/

// Returns map[WalletID]error
/*
func (self *Gateway) SaveWallets() interface{} {
	self.requests <- func() interface{} {
		return self.Visor.SaveWallets(self.d.Visor.Visor)
	}
	r := <-self.responses
	return r
}

// Returns error
func (self *Gateway) SaveWallet(walletID wallet.WalletID) interface{} {
	self.requests <- func() interface{} {
		return self.Visor.SaveWallet(self.d.Visor.Visor, walletID)
	}
	r := <-self.responses
	return r
}

// Returns an error
func (self *Gateway) ReloadWallets() interface{} {
	self.requests <- func() interface{} {
		return self.Visor.ReloadWallets(self.d.Visor.Visor)
	}
	r := <-self.responses
	return r
}
*/

// Returns a *visor.ReadableWallet
/*
func (self *Gateway) GetWallet(walletID wallet.WalletID) interface{} {
	self.requests <- func() interface{} {
		return self.Visor.GetWallet(self.d.Visor.Visor, walletID)
	}
	r := <-self.responses
	return r
}
*/

// Returns a *ReadableWallets
/*
func (self *Gateway) GetWallets() interface{} {
	self.requests <- func() interface{} {
		return self.Visor.GetWallets(self.d.Visor.Visor)
	}
	r := <-self.responses
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
		self.requests <- func() interface{} {
			return self.Visor.CreateWallet(self.d.Visor.Visor)
		}
		r := <-self.responses
		return r
	//
}
*/
/* Blockchain & Transaction status */

// Returns a *BlockchainProgress
func (self *Gateway) GetBlockchainProgress() interface{} {
	self.requests <- func() interface{} {
		return self.Daemon.GetBlockchainProgress(self.d.Visor)
	}
	r := <-self.responses
	return r
}

// Returns a *ResendResult
func (self *Gateway) ResendTransaction(txn cipher.SHA256) interface{} {
	self.requests <- func() interface{} {
		return self.Daemon.ResendTransaction(self.d.Visor, self.d.Pool, txn)
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
func (self *Gateway) GetTransaction(txn cipher.SHA256) interface{} {
	self.requests <- func() interface{} {
		return self.Visor.GetTransaction(self.d.Visor.Visor, txn)
	}
	r := <-self.responses
	return r
}

// Returns a *visor.TransactionResults
func (self *Gateway) GetAddressTransactions(a cipher.Address) interface{} {
	self.requests <- func() interface{} {
		return self.Visor.GetAddressTransactions(self.d.Visor.Visor, a)
	}
	r := <-self.responses
	return r
}
