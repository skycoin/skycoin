package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/wallet"
)

//go:generate goautomock -template=testify RPCIface
type RPCIface interface {
	GetBlockchainMetadata(v Visorer) *BlockchainMetadata
	GetUnspent(v Visorer) blockdb.UnspentPool
	GetUnconfirmedSpends(v Visorer, addrs []cipher.Address) (coin.AddressUxOuts, error)
	GetUnconfirmedReceiving(v Visorer, addrs []cipher.Address) (coin.AddressUxOuts, error)
	GetUnconfirmedTxns(v Visorer, addresses []cipher.Address) []UnconfirmedTxn
	GetBlock(v Visorer, seq uint64) (*coin.SignedBlock, error)
	GetBlocks(v Visorer, start, end uint64) []coin.SignedBlock
	GetLastBlocks(v Visorer, num uint64) []coin.SignedBlock
	GetBlockBySeq(v Visorer, n uint64) (*coin.SignedBlock, error)
	GetTransaction(v Visorer, txHash cipher.SHA256) (*Transaction, error)
	GetAddressTxns(v Visorer, addr cipher.Address) ([]Transaction, error)
	CreateWallet(wltName string, options wallet.Options) (wallet.Wallet, error)
	NewAddresses(wltName string, num uint64) ([]cipher.Address, error)
	GetWalletAddresses(wltID string) ([]cipher.Address, error)
	CreateAndSignTransaction(wltID string, vld wallet.Validator, unspent blockdb.UnspentGetter,
		headTime, coins uint64, dest cipher.Address) (*coin.Transaction, error)
	UpdateWalletLabel(wltID, label string) error
	GetWallet(wltID string) (wallet.Wallet, error)
	GetWallets() wallet.Wallets
	ReloadWallets() error
	GetBuildInfo() BuildInfo
}
