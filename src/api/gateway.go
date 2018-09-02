package api

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

//go:generate go install
//go:generate mockery -name Gatewayer -case underscore -inpkg -testonly

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	Spend(wltID string, password []byte, coins uint64, dest cipher.Address) (*coin.Transaction, error)
	CreateTransaction(w wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error)
	GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalance, error)
	GetWallet(wltID string) (*wallet.Wallet, error)
	GetWallets() (wallet.Wallets, error)
	UpdateWalletLabel(wltID, label string) error
	GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTransaction, error)
	GetWalletUnconfirmedTxnsVerbose(wltID string) ([]readable.UnconfirmedTransactionVerbose, error)
	CreateWallet(wltName string, options wallet.Options) (*wallet.Wallet, error)
	NewAddresses(wltID string, password []byte, n uint64) ([]cipher.Address, error)
	GetWalletDir() (string, error)
	IsWalletAPIEnabled() bool
	EncryptWallet(wltID string, password []byte) (*wallet.Wallet, error)
	DecryptWallet(wltID string, password []byte) (*wallet.Wallet, error)
	GetWalletSeed(wltID string, password []byte) (string, error)
	GetSignedBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error)
	GetBlockByHashVerbose(hash cipher.SHA256) (*readable.BlockVerbose, error)
	GetSignedBlockBySeq(seq uint64) (*coin.SignedBlock, error)
	GetBlockBySeqVerbose(seq uint64) (*readable.BlockVerbose, error)
	GetBlocks(start, end uint64) (*readable.Blocks, error)
	GetBlocksVerbose(start, end uint64) (*readable.BlocksVerbose, error)
	GetLastBlocks(num uint64) (*readable.Blocks, error)
	GetLastBlocksVerbose(num uint64) (*readable.BlocksVerbose, error)
	GetBuildInfo() visor.BuildInfo
	GetUnspentOutputs(filters ...daemon.OutputsFilter) (*readable.OutputSet, error)
	GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error)
	GetBlockchainMetadata() (*visor.BlockchainMetadata, error)
	GetBlockchainProgress() (*daemon.BlockchainProgress, error)
	GetConnection(addr string) *daemon.Connection
	GetConnections() *daemon.Connections
	GetDefaultConnections() []string
	GetTrustConnections() []string
	GetExchgConnection() []string
	GetAllUnconfirmedTxns() ([]visor.UnconfirmedTransaction, error)
	GetAllUnconfirmedTxnsVerbose() ([]readable.UnconfirmedTransactionVerbose, error)
	GetTransaction(txid cipher.SHA256) (*visor.Transaction, error)
	GetTransactionWithStatus(txid cipher.SHA256) (*readable.TransactionWithStatus, error)
	GetTransactionWithStatusVerbose(txid cipher.SHA256) (*readable.TransactionWithStatusVerbose, error)
	GetTransactionsWithStatus(flts []visor.TxFilter) (*readable.TransactionsWithStatus, error)
	GetTransactionsWithStatusVerbose(flts []visor.TxFilter) (*readable.TransactionsWithStatusVerbose, error)
	InjectBroadcastTransaction(txn coin.Transaction) error
	ResendUnconfirmedTxns() (*daemon.ResendResult, error)
	GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error)
	GetAddrUxOuts(addr []cipher.Address) ([]*historydb.UxOut, error)
	GetVerboseTransactionsForAddress(a cipher.Address) ([]readable.TransactionVerbose, error)
	GetRichlist(includeDistribution bool) (visor.Richlist, error)
	GetAddressCount() (uint64, error)
	GetHealth() (*daemon.Health, error)
	UnloadWallet(id string) error
	VerifyTxnVerbose(txn *coin.Transaction) ([]wallet.UxBalance, bool, error)
}
