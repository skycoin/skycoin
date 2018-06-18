package api

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

//go:generate go install
//go:generate goautomock -template=testify Gatewayer

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	Spend(wltID string, password []byte, coins uint64, dest cipher.Address) (*coin.Transaction, error)
	CreateTransaction(w wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error)
	GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalance, error)
	GetWallet(wltID string) (*wallet.Wallet, error)
	GetWallets() (wallet.Wallets, error)
	UpdateWalletLabel(wltID, label string) error
	GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error)
	CreateWallet(wltName string, options wallet.Options) (*wallet.Wallet, error)
	NewAddresses(wltID string, password []byte, n uint64) ([]cipher.Address, error)
	GetWalletDir() (string, error)
	IsWalletAPIEnabled() bool
	EncryptWallet(wltID string, password []byte) (*wallet.Wallet, error)
	DecryptWallet(wltID string, password []byte) (*wallet.Wallet, error)
	GetWalletSeed(wltID string, password []byte) (string, error)
	GetSignedBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error)
	GetSignedBlockBySeq(seq uint64) (*coin.SignedBlock, error)
	GetBlocks(start, end uint64) (*visor.ReadableBlocks, error)
	GetLastBlocks(num uint64) (*visor.ReadableBlocks, error)
	GetBuildInfo() visor.BuildInfo
	GetUnspentOutputs(filters ...daemon.OutputsFilter) (*visor.ReadableOutputSet, error)
	GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error)
	GetBlockchainMetadata() (*visor.BlockchainMetadata, error)
	GetBlockchainProgress() (*daemon.BlockchainProgress, error)
	GetConnection(addr string) *daemon.Connection
	GetConnections() *daemon.Connections
	GetDefaultConnections() []string
	GetTrustConnections() []string
	GetExchgConnection() []string
	GetAllUnconfirmedTxns() ([]visor.UnconfirmedTxn, error)
	GetTransaction(txid cipher.SHA256) (*visor.Transaction, error)
	GetTransactions(flts ...visor.TxFilter) ([]visor.Transaction, error)
	InjectBroadcastTransaction(txn coin.Transaction) error
	ResendUnconfirmedTxns() (*daemon.ResendResult, error)
	GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error)
	GetAddrUxOuts(addr []cipher.Address) ([]*historydb.UxOut, error)
	GetTransactionsForAddress(a cipher.Address) ([]daemon.ReadableTransaction, error)
	GetRichlist(includeDistribution bool) (visor.Richlist, error)
	GetAddressCount() (uint64, error)
	GetHealth() (*daemon.Health, error)
	UnloadWallet(id string) error
	VerifyTxnVerbose(txn *coin.Transaction) ([]wallet.UxBalance, bool, error)
}
