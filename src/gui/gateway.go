package gui

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

//go:generate goautomock -template=testify Gatewayer

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	Spend(wltID string, coins uint64, dest cipher.Address) (*coin.Transaction, error)
	GetWalletBalance(wltID string) (wallet.BalancePair, error)
	GetWallet(wltID string) (wallet.Wallet, error)
	GetWallets() (wallet.Wallets, error)
	UpdateWalletLabel(wltID, label string) error
	GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error)
	CreateWallet(wltName string, options wallet.Options) (wallet.Wallet, error)
	ScanAheadWalletAddresses(wltName string, scanN uint64) (wallet.Wallet, error)
	NewAddresses(wltID string, n uint64) ([]cipher.Address, error)
	GetWalletDir() (string, error)
	GetBlockByHash(hash cipher.SHA256) (block *visor.ReadableBlock, ok bool)
	GetBlockBySeq(seq uint64) (block *visor.ReadableBlock, ok bool)
	GetBlocks(start, end uint64) (*visor.ReadableBlocks, error)
	GetLastBlocks(num uint64) (*visor.ReadableBlocks, error)
	GetBuildInfo() visor.BuildInfo
	GetUnspentOutputs(filters ...daemon.OutputsFilter) (visor.ReadableOutputSet, error)
	GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error)
	GetBlockchainMetadata() interface{}
	GetBlockchainProgress() interface{}
	GetConnection(addr string) interface{}
	GetConnections() interface{}
	GetDefaultConnections() interface{}
	GetTrustConnections() interface{}
	GetExchgConnection() interface{}
	GetAllUnconfirmedTxns() ([]*visor.ReadableUnconfirmedTxn, error)
	GetLastTxs() ([]*visor.TransactionResult, error)
	GetTransaction(txid cipher.SHA256) (*visor.Transaction, *visor.TransactionResult, error)
	GetTransactions(flts ...visor.TxFilter) ([]*visor.TransactionResult, error)
	InjectBroadcastTransaction(txn coin.Transaction) error
	ResendUnconfirmedTxns() *daemon.ResendResult
	GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error)
	GetAddrUxOuts(addr cipher.Address) ([]*historydb.UxOutJSON, error)
	GetAddressTxns(a cipher.Address) (*visor.TransactionResults, error)
	GetRichlist(includeDistribution bool) (visor.Richlist, error)
	GetAddressCount() (uint64, error)
}
