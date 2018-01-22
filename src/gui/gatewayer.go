package gui

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	Spend(wltID string, coins uint64, dest cipher.Address) (*coin.Transaction, error)
	GetWalletBalance(wltID string) (wallet.BalancePair, error)
	GetWallet(wltID string) (wallet.Wallet, error)
	UpdateWalletLabel(wltID, label string) error
	GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error)
	CreateWallet(wltName string, options wallet.Options) (wallet.Wallet, error)
	ScanAheadWalletAddresses(wltName string, scanN uint64) (wallet.Wallet, error)
	NewAddresses(wltID string, n uint64) ([]cipher.Address, error)
	GetWalletDir() (string, error)
	GetBlockByHash(hash cipher.SHA256) (block coin.SignedBlock, ok bool)
	GetBlockBySeq(seq uint64) (block coin.SignedBlock, ok bool)
	GetBlocks(start, end uint64) (*visor.ReadableBlocks, error)
	GetAllUnconfirmedTxns() []visor.UnconfirmedTxn
	GetTransaction(txid cipher.SHA256) (tx *visor.Transaction, err error)
	InjectTransaction(txn coin.Transaction) error
	ResendUnconfirmedTxns() (rlt *daemon.ResendResult)
}

//// Gatewayer interface for Gateway methods
//type Gatewayer interface {
//	Spend(wltID string, coins uint64, dest cipher.Address) (*coin.Transaction, error)
//	GetWalletBalance(wltID string) (wallet.BalancePair, error)
//	GetWallet(wltID string) (wallet.Wallet, error)
//	UpdateWalletLabel(wltID, label string) error
//	GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error)
//	CreateWallet(wltName string, options wallet.Options) (wallet.Wallet, error)
//	ScanAheadWalletAddresses(wltName string, scanN uint64) (wallet.Wallet, error)
//	NewAddresses(wltID string, n uint64) ([]cipher.Address, error)
//	GetWalletDir() (string, error)
//	GetBlockByHash(hash cipher.SHA256) (block coin.SignedBlock, ok bool)
//	GetBlockBySeq(seq uint64) (block coin.SignedBlock, ok bool)
//	GetBlocks(start, end uint64) (*visor.ReadableBlocks, error)
//}
