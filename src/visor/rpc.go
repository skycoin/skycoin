package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/wallet"
)

// TransactionResult represents transaction result
type TransactionResult struct {
	Status      TransactionStatus   `json:"status"`
	Time        uint64              `json:"time"`
	Transaction ReadableTransaction `json:"txn"`
}

// NewTransactionResult converts Transaction to TransactionResult
func NewTransactionResult(tx *Transaction) (*TransactionResult, error) {
	if tx == nil {
		return nil, nil
	}

	rbTx, err := NewReadableTransaction(tx)
	if err != nil {
		return nil, err
	}

	return &TransactionResult{
		Transaction: *rbTx,
		Status:      tx.Status,
		Time:        tx.Time,
	}, nil
}

// ReadableBlocks an array of readable blocks.
type ReadableBlocks struct {
	Blocks []ReadableBlock `json:"blocks"`
}

// TransactionResults array of transaction results
type TransactionResults struct {
	Txns []TransactionResult `json:"txns"`
}

// NewTransactionResults converts []Transaction to []TransactionResults
func NewTransactionResults(txs []Transaction) (*TransactionResults, error) {
	txRlts := make([]TransactionResult, 0, len(txs))
	for _, tx := range txs {
		rbTx, err := NewReadableTransaction(&tx)
		if err != nil {
			return nil, err
		}

		txRlts = append(txRlts, TransactionResult{
			Transaction: *rbTx,
			Status:      tx.Status,
			Time:        tx.Time,
		})
	}

	return &TransactionResults{
		Txns: txRlts,
	}, nil
}

// RPC is balance check and transaction injection
// separate wallets out of visor
type RPC struct {
	v *Visor
}

// MakeRPC make RPC instance
func MakeRPC(v *Visor) RPC {
	return RPC{
		v: v,
	}
}

// GetBlockchainMetadata get blockchain meta data
func (rpc RPC) GetBlockchainMetadata(v *Visor) *BlockchainMetadata {
	bm := v.GetBlockchainMetadata()
	return &bm
}

// GetUnspent gets unspent
func (rpc RPC) GetUnspent(v *Visor) blockdb.UnspentPool {
	return v.Blockchain.Unspent()
}

// GetUnconfirmedSpends get unconfirmed spents
func (rpc RPC) GetUnconfirmedSpends(v *Visor, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	return v.Unconfirmed.SpendsOfAddresses(addrs, rpc.GetUnspent(v))
}

// GetUnconfirmedReceiving returns unconfirmed
func (rpc RPC) GetUnconfirmedReceiving(v *Visor, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	head, err := v.Blockchain.Head()
	if err != nil {
		return coin.AddressUxOuts{}, err
	}
	return v.Unconfirmed.RecvOfAddresses(head.Head, addrs)
}

// GetUnconfirmedTxns gets unconfirmed transactions
func (rpc RPC) GetUnconfirmedTxns(v *Visor, addresses []cipher.Address) []UnconfirmedTxn {
	return v.GetUnconfirmedTxns(ToAddresses(addresses))
}

// GetBlock gets block
func (rpc RPC) GetBlock(v *Visor, seq uint64) (*coin.SignedBlock, error) {
	return v.GetBlock(seq)
}

// GetBlocks gets blocks
func (rpc RPC) GetBlocks(v *Visor, start, end uint64) []coin.SignedBlock {
	return v.GetBlocks(start, end)
}

// GetLastBlocks returns the last N blocks
func (rpc RPC) GetLastBlocks(v *Visor, num uint64) []coin.SignedBlock {
	return v.GetLastBlocks(num)
}

// GetBlockBySeq get block in depth
func (rpc RPC) GetBlockBySeq(v *Visor, n uint64) (*coin.SignedBlock, error) {
	return v.GetBlockBySeq(n)

}

// GetTransaction gets transaction
func (rpc RPC) GetTransaction(v *Visor, txHash cipher.SHA256) (*Transaction, error) {
	return v.GetTransaction(txHash)
}

// GetAddressTxns get address transactions
func (rpc RPC) GetAddressTxns(v *Visor,
	addr cipher.Address) ([]Transaction, error) {
	return v.GetAddressTxns(addr)
}

// NewWallet creates new wallet
func (rpc *RPC) NewWallet(wltName string, ops ...wallet.Option) (wallet.Wallet, error) {
	return rpc.v.wallets.CreateWallet(wltName, ops...)
}

// NewAddresses generates new addresses in given wallet
func (rpc *RPC) NewAddresses(wltName string, num int) ([]cipher.Address, error) {
	return rpc.v.wallets.NewAddresses(wltName, num)
}

// GetWalletAddresses returns all addresses in given wallet
func (rpc *RPC) GetWalletAddresses(wltID string) ([]cipher.Address, error) {
	return rpc.v.wallets.GetAddresses(wltID)
}

// CreateAndSignTransaction creates and sign transaction from wallet
func (rpc *RPC) CreateAndSignTransaction(wltID string, vld wallet.Validator,
	unspent blockdb.UnspentGetter,
	headTime uint64,
	amt wallet.Balance,
	dest cipher.Address) (*coin.Transaction, error) {
	return rpc.v.wallets.CreateAndSignTransaction(wltID,
		vld,
		unspent,
		headTime,
		amt,
		dest)
}

// UpdateWalletLabel updates wallet label
func (rpc *RPC) UpdateWalletLabel(wltID, label string) error {
	return rpc.v.wallets.UpdateWalletLabel(wltID, label)
}

// GetWallet returns wallet by id
func (rpc *RPC) GetWallet(wltID string) (wallet.Wallet, bool) {
	return rpc.v.wallets.GetWallet(wltID)
}

// GetWallets returns all wallet
func (rpc *RPC) GetWallets() wallet.Wallets {
	return rpc.v.wallets.GetWallets()
}

// ReloadWallets reloads all wallet from files
func (rpc *RPC) ReloadWallets() error {
	return rpc.v.wallets.ReloadWallets()
}

// GetBuildInfo returns node build info, including version, build time, etc.
func (rpc *RPC) GetBuildInfo() BuildInfo {
	return rpc.v.Config.BuildInfo
}
