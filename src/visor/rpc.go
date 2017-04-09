package visor

import (
	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
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

func (self RPC) GetUnspent(v *Visor) coin.UnspentPool {
	return v.Blockchain.GetUnspent().Clone()
}

func (self RPC) GetUnconfirmedSpends(v *Visor, addrs map[cipher.Address]byte) coin.AddressUxOuts {
	unspent := self.GetUnspent(v)
	return v.Unconfirmed.SpendsForAddresses(&unspent, addrs)
}

func (self RPC) CreateSpendingTransaction(v *Visor, wlt wallet.Wallet, amt wallet.Balance, dest cipher.Address) (tx coin.Transaction, err error) {
	unspent := self.GetUnspent(v)
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

func (self RPC) GetUnspentOutputReadables(v *Visor) []ReadableOutput {
	ret := v.GetUnspentOutputReadables()
	return ret
}

func (self RPC) GetUnconfirmedTxns(v *Visor, addresses []cipher.Address) []ReadableUnconfirmedTxn {
	ret := v.GetUnconfirmedTxns(addresses)
	rut := make([]ReadableUnconfirmedTxn, len(ret))
	for i := range ret {
		rut[i] = NewReadableUnconfirmedTxn(&ret[i])
	}
	return rut
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
