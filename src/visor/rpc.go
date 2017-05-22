package visor

import (
	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
)

// TransactionResult represents transaction result
type TransactionResult struct {
	Status      TransactionStatus   `json:"status"`
	Transaction ReadableTransaction `json:"txn"`
}

// ReadableBlocks an array of readable blocks.
type ReadableBlocks struct {
	Blocks []ReadableBlock `json:"blocks"`
}

// TransactionResults array of transaction results
type TransactionResults struct {
	Txns []TransactionResult `json:"txns"`
}

// RPC is balance check and transaction injection
// separate wallets out of visor
type RPC struct{}

// GetBlockchainMetadata get blockchain meta data
func (rpc RPC) GetBlockchainMetadata(v *Visor) *BlockchainMetadata {
	bm := v.GetBlockchainMetadata()
	return &bm
}

// GetUnspent gets unspent
func (rpc RPC) GetUnspent(v *Visor) coin.UnspentPool {
	return v.Blockchain.GetUnspent().Clone()
}

// GetUnconfirmedSpends get unconfirmed spents
func (rpc RPC) GetUnconfirmedSpends(v *Visor, addrs map[cipher.Address]byte) coin.AddressUxOuts {
	unspent := rpc.GetUnspent(v)
	return v.Unconfirmed.SpendsForAddresses(&unspent, addrs)
}

// CreateSpendingTransaction creates spending transaction
func (rpc RPC) CreateSpendingTransaction(v *Visor, wlt wallet.Wallet, amt wallet.Balance, dest cipher.Address) (tx coin.Transaction, err error) {
	unspent := rpc.GetUnspent(v)
	tm := v.Blockchain.Time()
	tx, err = CreateSpendingTransaction(wlt, v.Unconfirmed, &unspent, tm, amt, dest)
	if err != nil {
		return
	}

	if err := tx.Verify(); err != nil {
		log.Panicf("Invalid transaction, %v", err)
	}

	if err := VerifyTransactionFee(v.Blockchain, &tx); err != nil {
		log.Panicf("Created invalid spending txn: visor fail, %v", err)
	}

	if err := v.Blockchain.VerifyTransaction(tx); err != nil {
		log.Panicf("Created invalid spending txn: blockchain fail, %v", err)
	}
	return
}

// GetUnspentOutputReadables gets unspent output readables
func (rpc RPC) GetUnspentOutputReadables(v *Visor) []ReadableOutput {
	ret := v.GetUnspentOutputReadables()
	return ret
}

// GetUnconfirmedTxns gets unconfirmed transactions
func (rpc RPC) GetUnconfirmedTxns(v *Visor, addresses []cipher.Address) []ReadableUnconfirmedTxn {
	ret := v.GetUnconfirmedTxns(ToAddresses(addresses))
	rut := make([]ReadableUnconfirmedTxn, len(ret))
	for i := range ret {
		rut[i] = NewReadableUnconfirmedTxn(&ret[i])
	}
	return rut
}

// GetBlock gets block
func (rpc RPC) GetBlock(v *Visor, seq uint64) *ReadableBlock {
	b, err := v.GetReadableBlock(seq)
	if err != nil {
		return nil
	}
	return &b
}

// GetBlocks gets blocks
func (rpc RPC) GetBlocks(v *Visor, start, end uint64) *ReadableBlocks {
	blocks := v.GetReadableBlocks(start, end)
	return &ReadableBlocks{blocks}
}

// GetBlockInDepth get block in depth
func (rpc RPC) GetBlockInDepth(v *Visor, n uint64) *ReadableBlock {
	if b := v.GetBlockBySeq(n); b != nil {
		block := NewReadableBlock(b)
		return &block
	}
	return nil
}

// GetTransaction gets transaction
func (rpc RPC) GetTransaction(v *Visor, txHash cipher.SHA256) (*TransactionResult, error) {
	txn, err := v.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}
	if txn == nil {
		return nil, nil
	}

	return &TransactionResult{
		Transaction: NewReadableTransaction(txn),
		Status:      txn.Status,
	}, nil
}

// GetAddressTransactions get address transactions
func (rpc RPC) GetAddressTransactions(v *Visor,
	addr cipher.Address) *TransactionResults {
	addrTxns := v.GetAddressTransactions(addr)
	txns := make([]TransactionResult, len(addrTxns))
	for i, tx := range addrTxns {
		txns[i] = TransactionResult{
			Transaction: NewReadableTransaction(&tx),
			Status:      tx.Status,
		}
	}
	return &TransactionResults{
		Txns: txns,
	}
}
