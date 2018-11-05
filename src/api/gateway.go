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
//go:generate mockery -name Gatewayer -case underscore -inpkg -testonly

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	Spend(wltID string, password []byte, coins uint64, dest cipher.Address) (*coin.Transaction, error)
	CreateTransaction(w wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error)
	GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalances, error)
	GetWallet(wltID string) (*wallet.Wallet, error)
	GetWallets() (wallet.Wallets, error)
	UpdateWalletLabel(wltID, label string) error
	GetWalletUnconfirmedTransactions(wltID string) ([]visor.UnconfirmedTransaction, error)
	GetWalletUnconfirmedTransactionsVerbose(wltID string) ([]visor.UnconfirmedTransaction, [][]visor.TransactionInput, error)
	CreateWallet(wltName string, options wallet.Options) (*wallet.Wallet, error)
	RecoverWallet(wltID, seed string, password []byte) (*wallet.Wallet, error)
	NewAddresses(wltID string, password []byte, n uint64) ([]cipher.Address, error)
	GetWalletDir() (string, error)
	EncryptWallet(wltID string, password []byte) (*wallet.Wallet, error)
	DecryptWallet(wltID string, password []byte) (*wallet.Wallet, error)
	GetWalletSeed(wltID string, password []byte) (string, error)
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
	GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error)
	GetBlockchainMetadata() (*visor.BlockchainMetadata, error)
	GetBlockchainProgress() (*daemon.BlockchainProgress, error)
	GetConnection(addr string) (*daemon.Connection, error)
	GetConnections(f func(c daemon.Connection) bool) ([]daemon.Connection, error)
	Disconnect(id uint64) error
	GetDefaultConnections() []string
	GetTrustConnections() []string
	GetExchgConnection() []string
	GetAllUnconfirmedTransactions() ([]visor.UnconfirmedTransaction, error)
	GetAllUnconfirmedTransactionsVerbose() ([]visor.UnconfirmedTransaction, [][]visor.TransactionInput, error)
	GetTransaction(txid cipher.SHA256) (*visor.Transaction, error)
	GetTransactionVerbose(txid cipher.SHA256) (*visor.Transaction, []visor.TransactionInput, error)
	GetTransactions(flts []visor.TxFilter) ([]visor.Transaction, error)
	GetTransactionsVerbose(flts []visor.TxFilter) ([]visor.Transaction, [][]visor.TransactionInput, error)
	InjectBroadcastTransaction(txn coin.Transaction) error
	ResendUnconfirmedTxns() ([]cipher.SHA256, error)
	GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error)
	GetSpentOutputsForAddresses(addr []cipher.Address) ([][]historydb.UxOut, error)
	GetVerboseTransactionsForAddress(a cipher.Address) ([]visor.Transaction, [][]visor.TransactionInput, error)
	GetRichlist(includeDistribution bool) (visor.Richlist, error)
	GetAddressCount() (uint64, error)
	GetHealth() (*daemon.Health, error)
	UnloadWallet(id string) error
	VerifyTxnVerbose(txn *coin.Transaction) ([]wallet.UxBalance, bool, error)
}
