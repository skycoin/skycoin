package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"
)

// TransactionResult represents transaction result
type TransactionResult struct {
	Status      TransactionStatus   `json:"status"`
	Time        uint64              `json:"time"`
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
func (rpc RPC) GetUnspent(v *Visor) *blockdb.UnspentPool {
	return v.Blockchain.Unspent()
}

// GetUnconfirmedSpends get unconfirmed spents
func (rpc RPC) GetUnconfirmedSpends(v *Visor, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	unspent := rpc.GetUnspent(v)
	return v.Unconfirmed.SpendsForAddresses(unspent, addrs)
}

// GetUnconfirmedReceiving returns unconfirmed
func (rpc RPC) GetUnconfirmedReceiving(v *Visor, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	return v.Unconfirmed.RecvOfAddresses(v.Blockchain.Head().Head, addrs)
}

// GetUnspentOutputReadables gets unspent output readables
func (rpc RPC) GetUnspentOutputReadables(v *Visor) ([]ReadableOutput, error) {
	return v.GetUnspentOutputReadables()
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
		Time:        txn.Time,
	}, nil
}

// GetAddressTxns get address transactions
func (rpc RPC) GetAddressTxns(v *Visor,
	addr cipher.Address) (*TransactionResults, error) {
	addrTxns, err := v.GetAddressTxns(addr)
	if err != nil {
		return nil, err
	}

	txns := make([]TransactionResult, len(addrTxns))
	for i, tx := range addrTxns {
		txns[i] = TransactionResult{
			Transaction: NewReadableTransaction(&tx),
			Status:      tx.Status,
			Time:        tx.Time,
		}
	}
	return &TransactionResults{
		Txns: txns,
	}, nil
}
