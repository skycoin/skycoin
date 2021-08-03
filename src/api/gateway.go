package api

import (
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/kvstorage"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

// Gateway bundles daemon.Daemon, Visor, wallet.Service and kvstorage.Manager into a single object
type Gateway struct {
	*daemon.Daemon
	*visor.Visor
	*wallet.Service
	*kvstorage.Manager
}

// NewGateway creates a Gateway
func NewGateway(d *daemon.Daemon, v *visor.Visor, w *wallet.Service, m *kvstorage.Manager) *Gateway {
	return &Gateway{
		Daemon:  d,
		Visor:   v,
		Service: w,
		Manager: m,
	}
}

//go:generate mockery -name Gatewayer -case underscore -inpkg -testonly

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	Daemoner
	Visorer
	Walleter
	Storer
}

// Daemoner interface for daemon.Daemon methods used by the API
type Daemoner interface {
	DaemonConfig() daemon.DaemonConfig
	GetConnection(addr string) (*daemon.Connection, error)
	GetConnections(f func(c daemon.Connection) bool) ([]daemon.Connection, error)
	DisconnectByGnetID(gnetID uint64) error
	GetDefaultConnections() []string
	GetTrustConnections() []string
	GetExchgConnection() []string
	GetBlockchainProgress(headSeq uint64) *daemon.BlockchainProgress
	InjectBroadcastTransaction(txn coin.Transaction) error
	InjectTransaction(txn coin.Transaction) error
}

// Visorer interface for visor.Visor methods used by the API
type Visorer interface {
	VisorConfig() visor.Config
	StartedAt() time.Time
	HeadBkSeq() (uint64, bool, error)
	GetBlockchainMetadata() (*visor.BlockchainMetadata, error)
	ResendUnconfirmedTxns() ([]cipher.SHA256, error)
	GetSignedBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error)
	GetSignedBlockByHashVerbose(hash cipher.SHA256) (*coin.SignedBlock, [][]visor.TransactionInput, error)
	GetSignedBlockBySeq(seq uint64) (*coin.SignedBlock, error)
	GetSignedBlockBySeqVerbose(seq uint64) (*coin.SignedBlock, [][]visor.TransactionInput, error)
	GetBlocks(seqs []uint64) ([]coin.SignedBlock, error)
	GetBlocksVerbose(seqs []uint64) ([]coin.SignedBlock, [][][]visor.TransactionInput, error)
	GetBlocksInRange(start, end uint64) ([]coin.SignedBlock, error)
	GetBlocksInRangeVerbose(start, end uint64) ([]coin.SignedBlock, [][][]visor.TransactionInput, error)
	GetLastBlocks(num uint64) ([]coin.SignedBlock, error)
	GetLastBlocksVerbose(num uint64) ([]coin.SignedBlock, [][][]visor.TransactionInput, error)
	GetUnspentOutputsSummary(filters []visor.OutputsFilter) (*visor.UnspentOutputsSummary, error)
	GetBalanceOfAddresses(addrs []cipher.Address) ([]wallet.BalancePair, error)
	VerifyTxnVerbose(txn *coin.Transaction, signed visor.TxnSignedFlag) ([]visor.TransactionInput, bool, error)
	AddressCount() (uint64, error)
	GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, uint64, error)
	GetSpentOutputsForAddresses(addr []cipher.Address) ([][]historydb.UxOut, uint64, error)
	GetRichlist(includeDistribution bool) (visor.Richlist, error)
	GetAllUnconfirmedTransactions() ([]visor.UnconfirmedTransaction, error)
	GetAllUnconfirmedTransactionsVerbose() ([]visor.UnconfirmedTransaction, [][]visor.TransactionInput, error)
	GetTransaction(txid cipher.SHA256) (*visor.Transaction, error)
	GetTransactionWithInputs(txid cipher.SHA256) (*visor.Transaction, []visor.TransactionInput, error)
	GetTransactions(flts []visor.TxFilter, order visor.SortOrder, page *visor.PageIndex) ([]visor.Transaction, uint64, error)
	GetTransactionsWithInputs(flts []visor.TxFilter, order visor.SortOrder, page *visor.PageIndex) ([]visor.Transaction, [][]visor.TransactionInput, uint64, error)
	GetWalletUnconfirmedTransactions(wltID string) ([]visor.UnconfirmedTransaction, error)
	GetWalletUnconfirmedTransactionsVerbose(wltID string) ([]visor.UnconfirmedTransaction, [][]visor.TransactionInput, error)
	GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalances, error)
	CreateTransaction(p transaction.Params, wp visor.CreateTransactionParams) (*coin.Transaction, []visor.TransactionInput, error)
	WalletCreateTransaction(wltID string, p transaction.Params, wp visor.CreateTransactionParams) (*coin.Transaction, []visor.TransactionInput, error)
	WalletCreateTransactionSigned(wltID string, password []byte, p transaction.Params, wp visor.CreateTransactionParams) (*coin.Transaction, []visor.TransactionInput, error)
	WalletSignTransaction(wltID string, password []byte, txn *coin.Transaction, signIndexes []int) (*coin.Transaction, []visor.TransactionInput, error)
	ScanWalletAddresses(wltID string, password []byte, num uint64) ([]cipher.Address, error)
	TransactionsFinder() wallet.TransactionsFinder
}

// Walleter interface for wallet.Service methods used by the API
type Walleter interface {
	UnloadWallet(wltID string) error
	EncryptWallet(wltID string, password []byte) (wallet.Wallet, error)
	DecryptWallet(wltID string, password []byte) (wallet.Wallet, error)
	GetWalletSeed(wltID string, password []byte) (string, string, error)
	CreateWallet(wltName string, options wallet.Options) (wallet.Wallet, error)
	RecoverWallet(wltID, seed, seedPassphrase string, password []byte) (wallet.Wallet, error)
	NewAddresses(wltID string, password []byte, n uint64, options ...wallet.Option) ([]cipher.Address, error)
	ScanAddresses(wltID string, password []byte, n uint64, tf wallet.TransactionsFinder) ([]cipher.Address, error)
	GetWallet(wltID string) (wallet.Wallet, error)
	GetWallets() (wallet.Wallets, error)
	UpdateWalletLabel(wltID, label string) error
	WalletDir() (string, error)
}

// Storer interface for kvstorage.Manager methods used by the API
type Storer interface {
	GetStorageValue(storageType kvstorage.Type, key string) (string, error)
	GetAllStorageValues(storageType kvstorage.Type) (map[string]string, error)
	AddStorageValue(storageType kvstorage.Type, key, val string) error
	RemoveStorageValue(storageType kvstorage.Type, key string) error
}
