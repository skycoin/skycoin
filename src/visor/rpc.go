package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	//"github.com/skycoin/skycoin/src/wallet"
)

/*
RPC is balance check and transaction injection
- seperate wallets out of visor
*/
type TransactionResult struct {
	Status      TransactionStatus   `json:"status"`
	Transaction ReadableTransaction `json:"txn"`
}

// An array of readable blocks.
type ReadableBlocks struct {
	Blocks []ReadableBlock `json:"blocks"`
}

type TransactionResults struct {
	Txns []TransactionResult `json:"txns"`
}

type RPC struct{}

func (self RPC) GetBlockchainMetadata(v *Visor) *BlockchainMetadata {
	bm := v.GetBlockchainMetadata()
	return &bm
}

func (self RPC) GetUnspentOutputReadables(v *Visor) []ReadableOutput {
	ret := v.GetUnspentOutputReadables()
	return ret
}

func (self RPC) GetUnconfirmedTxns(v *Visor, addresses []cipher.Address) []ReadableUnconfirmedTxn {
	ret := v.GetUnconfirmedTxns(addresses)
	return ret
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

func (self RPC) GetBlockInDepth(v *Visor, n uint64) *ReadableBlock {
	if b := v.GetBlockBySeq(n); b != nil {
		block := NewReadableBlock(b)
		return &block
	}
	return nil
}

func (self RPC) GetTransaction(v *Visor, txHash cipher.SHA256) (*TransactionResult, error) {
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

func (self RPC) GetAddressTransactions(v *Visor,
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
