package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/wallet"
)

//go:generate goautomock -template=testify RPCIface
type RPCIface interface {
	GetBlockchainMetadata(v *Visor) *BlockchainMetadata
	GetUnspent(v *Visor) blockdb.UnspentPool
	GetUnconfirmedSpends(v *Visor, addrs []cipher.Address) (coin.AddressUxOuts, error)
	GetUnconfirmedReceiving(v *Visor, addrs []cipher.Address) (coin.AddressUxOuts, error)
	GetUnconfirmedTxns(v *Visor, addresses []cipher.Address) []UnconfirmedTxn
	GetBlock(v *Visor, seq uint64) (*coin.SignedBlock, error)
	GetBlocks(v *Visor, start, end uint64) []coin.SignedBlock
	GetLastBlocks(v *Visor, num uint64) []coin.SignedBlock
	GetBlockBySeq(v *Visor, n uint64) (*coin.SignedBlock, error)
	GetTransaction(v *Visor, txHash cipher.SHA256) (*Transaction, error)
	GetAddressTxns(v *Visor, addr cipher.Address) ([]Transaction, error)
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
