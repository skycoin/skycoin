package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

/*
RPC is balance check and transaction injection
- seperate wallets out of visor
*/
type TransactionResult struct {
	Transaction ReadableTransaction `json:"txn"`
	Status      TransactionStatus   `json:"status"`
}

// An array of readable blocks.
type ReadableBlocks struct {
	Blocks []ReadableBlock `json:"blocks"`
}

type TransactionResults struct {
	Txns []TransactionResult `json:"txns"`
}

type RPC struct{}

func (self RPC) GetWalletBalance(v *Visor,
	walletID wallet.WalletID) *wallet.BalancePair {
	bp := v.WalletBalance(walletID)
	return &bp
}

func (self RPC) ReloadWallets(v *Visor) error {
	return v.ReloadWallets()
}

func (self RPC) SaveWallet(v *Visor, walletID wallet.WalletID) error {
	return v.SaveWallet(walletID)
}

func (self RPC) SaveWallets(v *Visor) map[wallet.WalletID]error {
	return v.SaveWallets()
}

func (self RPC) CreateWallet(v *Visor, seed string) *wallet.ReadableWallet {
	w := v.CreateWallet()
	return wallet.NewReadableWallet(w)
}

func (self RPC) GetWallet(v *Visor,
	walletID wallet.WalletID) *wallet.ReadableWallet {
	w := v.Wallets.Get(walletID)
	if w == nil {
		return nil
	} else {
		return wallet.NewReadableWallet(w)
	}
}

func (self RPC) GetWallets(v *Visor) []*wallet.ReadableWallet {
	return v.Wallets.ToPublicReadable()
}

func (self RPC) GetBlockchainMetadata(v *Visor) *BlockchainMetadata {
	bm := v.GetBlockchainMetadata()
	return &bm
}

func (self RPC) GetBlock(v *Visor, seq uint64) *ReadableBlock {
	b, err := v.GetReadableBlock(seq)
	if err != nil {
		return nil
	}
	return &b
}

func (self RPC) GetBlocks(v *Visor, start, end uint64) *ReadableBlocks {
	blocks := v.GetReadableBlocks(start, end)
	return &ReadableBlocks{blocks}
}

func (self RPC) GetTransaction(v *Visor,
	txHash cipher.SHA256) *TransactionResult {
	txn := v.GetTransaction(txHash)
	return &TransactionResult{
		Transaction: NewReadableTransaction(&txn.Txn),
		Status:      txn.Status,
	}
}

func (self RPC) GetAddressTransactions(v *Visor,
	addr cipher.Address) *TransactionResults {
	addrTxns := v.GetAddressTransactions(addr)
	txns := make([]TransactionResult, len(addrTxns))
	for i, tx := range addrTxns {
		txns[i] = TransactionResult{
			Transaction: NewReadableTransaction(&tx.Txn),
			Status:      tx.Status,
		}
	}
	return &TransactionResults{
		Txns: txns,
	}
}
